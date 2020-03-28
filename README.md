# Error backend for nginx ingress

Warning: This repository is Work In Progress.
 
This is error backend for [Nginx ingress Controller][the-ingress]. 

It is based on the [example error backend][original] with predefined error pages and Go templating. 

We used [haproxy error pages][error-pages] from [Jonatahan Rosewood][jonathan] as a base for the error messages.

# How to build and test

    # build static image
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o custom-error-pages
    
    # run 
    DEBUG=1 ERROR_FILES_PATH=./rootfs/www ./custom-error-pages
    
    # test in other terminal
    curl localhost:8080 -H 'X-Code: 502' -I

# How to do build Docker image


    # build static Docker image
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o custom-error-pages
    docker build -t localhost/custom-error -f Dockerfile.devel .

    # run the image
    docker run --rm -it -e DEBUG=1 -p 8080:8080 localhost/custom-error:latest

## How to deploy

1. deploy the image to Deployment

2. deply Service named `default-backend` to some namespace ( we use `ingress-nginx` )

3. according to the [ingress parameters][ingress-parameters], update [command line arguments][command-line-args] to set [default backend][default-backend].
   `--default-backend-service=ingress-nginx/default-backend`

4. Add `custom-http-errors: 503,502,403` to [ingress config map][custom-http-errors].


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
