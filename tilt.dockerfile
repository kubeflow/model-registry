FROM alpine:3.21

WORKDIR /

COPY devenv/bin/model-registry /model-registry

ENTRYPOINT ["/model-registry"]
