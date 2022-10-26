/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// FormatHeader name of the header used to extract the format
	FormatHeader = "X-Format"

	// CodeHeader name of the header used as source of the HTTP status code to return
	CodeHeader = "X-Code"

	// ContentType name of the header that defines the format of the reply
	ContentType = "Content-Type"

	// OriginalURI name of the header with the original URL from NGINX
	OriginalURI = "X-Original-URI"

	// Namespace name of the header that contains information about the Ingress namespace
	Namespace = "X-Namespace"

	// IngressName name of the header that contains the matched Ingress
	IngressName = "X-Ingress-Name"

	// ServiceName name of the header that contains the matched Service in the Ingress
	ServiceName = "X-Service-Name"

	// ServicePort name of the header that contains the matched Service port in the Ingress
	ServicePort = "X-Service-Port"

	// RequestId is a unique ID that identifies the request - same as for backend service
	RequestId = "X-Request-ID"

	// ErrFilesPathVar is the name of the environment variable indicating
	// the location on disk of files served by the handler.
	ErrFilesPathVar = "ERROR_FILES_PATH"

	// CustomErrFilesPathVar is the name of the environment variable indicating
	// the location on disk of override files served by the handler.
	CustomErrFilesPathVar = "CUSTOM_ERROR_FILES_PATH"	
)

func main() {
	errFilesPath := "/www"
	if os.Getenv(ErrFilesPathVar) != "" {
		errFilesPath = os.Getenv(ErrFilesPathVar)
	}

	var customErrFilesPath string
	if os.Getenv(CustomErrFilesPathVar) != "" {
		customErrFilesPath = os.Getenv(CustomErrFilesPathVar)
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	errFilesMux := http.NewServeMux()
	errFilesMux.HandleFunc("/", errorHandler(errFilesPath, customErrFilesPath))

	go func() {
		errFilesServer := http.Server{
			Addr:    ":8080",
			Handler: errFilesMux,
		}
		errFilesServer.ListenAndServe()
		wg.Done()
	}()

	promMux := http.NewServeMux()
	promMux.Handle("/metrics", promhttp.Handler())

	promMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	go func() {
		promServer := http.Server{
			Addr:    ":8081",
			Handler: promMux,
		}
		promServer.ListenAndServe()
		wg.Done()
	}()

	wg.Wait()
}

func errorHandler(path, customPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ext := "html"

		if os.Getenv("DEBUG") != "" {
			w.Header().Set(FormatHeader, r.Header.Get(FormatHeader))
			w.Header().Set(CodeHeader, r.Header.Get(CodeHeader))
			w.Header().Set(ContentType, r.Header.Get(ContentType))
			w.Header().Set(OriginalURI, r.Header.Get(OriginalURI))
			w.Header().Set(Namespace, r.Header.Get(Namespace))
			w.Header().Set(IngressName, r.Header.Get(IngressName))
			w.Header().Set(ServiceName, r.Header.Get(ServiceName))
			w.Header().Set(ServicePort, r.Header.Get(ServicePort))
			w.Header().Set(RequestId, r.Header.Get(RequestId))
		}

		format := r.Header.Get(FormatHeader)

		switch format {
		case "application/json":
			format = "application/json"
			ext = "json"
		default:
			format = "text/html"
			ext = "html"
		}
		w.Header().Set(ContentType, format)

		var code int
		var status_label string
		errCode := r.Header.Get(CodeHeader)
		if errCode == "" {
			code = 404
			status_label = "default"
		} else {
			status_label = errCode
			a_code, err := strconv.Atoi(errCode)
			if err != nil {
				code = 404
				log.Printf("unexpected error reading return code: %v. Using %v", err, code)
			} else {
				code = a_code
			}
		}
		w.WriteHeader(code)

		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}

		var filename string
		if customPath != "" {
			filename = fmt.Sprintf("%v/%v%v", customPath, code, ext)
			_, err := os.Stat(filename)
			if err != nil {
				filename = ""
				if os.Getenv("DEBUG") != "" {
					log.Printf("custom error file for %v code not found: %v", code, err)
				}
			} 
		}

		if filename == "" {
			filename = fmt.Sprintf("%v/%v%v", path, code, ext)
			_, err := os.Stat(filename)
			if err != nil {
				log.Printf("unexpected error opening file: %v", err)
				scode := strconv.Itoa(code)
				filename := fmt.Sprintf("%v/%cxx%v", path, scode[0], ext)
				_, err = os.Stat(filename)
				if err != nil {
					log.Printf("unexpected error opening file: %v", err)
					errorCount.WithLabelValues(status_label).Inc()
					http.NotFound(w, r)
					return
				}
			}
		}

		tmpl := template.Must(template.ParseFiles(filename))
		data := struct {
			RequestId string
		}{
			r.Header.Get(RequestId),
		}

		tmpl.Execute(w, data)

		proto := strconv.Itoa(r.ProtoMajor)
		proto = fmt.Sprintf("%s.%s", proto, strconv.Itoa(r.ProtoMinor))
		requestCount.WithLabelValues(proto, status_label).Inc()
		duration := time.Now().Sub(start).Seconds()
		requestDuration.WithLabelValues(proto, status_label).Observe(duration)
	}
}
