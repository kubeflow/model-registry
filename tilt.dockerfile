FROM alpine:3.12

WORKDIR /

COPY devenv/bin/model-registry /model-registry

ENTRYPOINT ["/model-registry"]
