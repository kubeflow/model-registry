FROM alpine:3.12

WORKDIR /

COPY devenv/bin/bff /bff
COPY devenv/ui-static /static

ENTRYPOINT ["/bff"]
