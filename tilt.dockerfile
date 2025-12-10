FROM alpine:3.23

WORKDIR /

COPY devenv/bin/model-registry /model-registry

ENTRYPOINT ["/model-registry"]
