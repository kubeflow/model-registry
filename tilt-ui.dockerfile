FROM alpine:3.21

WORKDIR /

COPY devenv/bin/bff /bff
COPY devenv/ui-static /static

ENTRYPOINT ["/bff"]
