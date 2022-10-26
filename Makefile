.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o custom-error-pages

.PHONY: test-run
test-run: build
	DEBUG=1 ERROR_FILES_PATH=./rootfs/www CUSTOM_ERROR_FILES_PATH=./rootfs/www.override ./custom-error-pages