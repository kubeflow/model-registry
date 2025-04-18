FROM alpine:3.21

WORKDIR /

COPY devenv/bin/manager /manager

ENTRYPOINT ["/manager"]
