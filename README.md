# Error backend for nginx ingress

**Warning:** This repository is _Proof Of Concept_.
 
This is default backend for [Nginx ingress Controller][the-ingress]. 
It is based on the [example error backend][original] with predefined error pages and Go templating.
We have fixed some errors and added few new ones.

We used [haproxy error pages][error-pages] from [Jonathan Rosewood][jonathan] as a base for the error messages. Thank you.

# How to build and test

    # build static image
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o custom-error-pages
    
    # run 
    DEBUG=1 ERROR_FILES_PATH=./rootfs/www ./custom-error-pages
    
    # test in other terminal
    curl localhost:8080 -H 'X-Code: 502' -I
    curl localhost:8081/metrics

There is default HTTP backend on port 8080. You should use this as your [ingress default backend][default-backend].
There is secondary port 8081, which can be used for health checking (on `/healthz` URI) and monitoring (`/metrics`)

# How to build Docker image

    # build static Docker image
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o custom-error-pages
    docker build -t localhost/custom-error -f Dockerfile.devel .

    # run and test the image
    docker run --rm -it -e DEBUG=1 -p 8080:8080 -p 8081:8081 localhost/custom-error:latest

## How to deploy

1. deploy the image to Deployment

2. deply Service named `default-backend` to some namespace ( we use `ingress-nginx` )

3. according to the [ingress parameters][ingress-parameters], update [command line arguments][command-line-args] to set [default backend][default-backend].
   `--default-backend-service=ingress-nginx/default-backend`

4. Add `custom-http-errors: 503,502,403` to [ingress config map][custom-http-errors].

5. Watch metrics `default_http_backend_http_error_count_total > 0`. 

Example manifests are available at [GitHub repository][example-manifests]. You still have to modify ingress config.


# License

Apache 2 (same as from [original][original] example)


[the-ingress]: https://kubernetes.github.io/ingress-nginx/
[custom-errors]: https://kubernetes.github.io/ingress-nginx/user-guide/custom-errors/
[original]: https://github.com/kubernetes/ingress-nginx/tree/master/images/custom-error-pages
[error-pages]: https://github.com/Jonathan-Rosewood/haproxy-custom-errors
[jonathan]: https://github.com/Jonathan-Rosewood
[ingress-config]: https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/
[ingress-parameters]: https://kubernetes.github.io/ingress-nginx/examples/customization/custom-errors/#ingress-controller-configuration
[default-backend]: https://kubernetes.github.io/ingress-nginx/user-guide/default-backend/
[command-line-args]: https://kubernetes.github.io/ingress-nginx/user-guide/cli-arguments/
[custom-http-errors]: https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#custom-http-errors
[example-manifests]: https://github.com/wftech/nginx-ingress-error-backend/tree/master/manifests
