FROM alpine:3.22

WORKDIR /

COPY devenv/bin/model-registry /model-registry

ENTRYPOINT ["/model-registry"]
