# Error backend for nginx ingress

This is error backend for [Nginx ingress Controller][the-ingress].

It is based on the [example error backend][original] with predefined error pages and Go templating. 

We used [haproxy error pages][error-pages] from [Jonatahan Rosewood][jonathan]. Thanks.

# How to build

TBD

# How to deploy

TBD

# License

Apache 2 (from [original][original] source])


[the-ingress]: https://kubernetes.github.io/ingress-nginx/
[custom-errors]: https://kubernetes.github.io/ingress-nginx/user-guide/custom-errors/
[original]: https://github.com/kubernetes/ingress-nginx/tree/master/images/custom-error-pages
[error-pages]: https://github.com/Jonathan-Rosewood/haproxy-custom-errors
[jonathan]: https://github.com/Jonathan-Rosewood
