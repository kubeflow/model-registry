FROM alpine:3.23

WORKDIR /

COPY devenv/bin/manager /manager

ENTRYPOINT ["/manager"]
