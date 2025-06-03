FROM alpine:3.22

WORKDIR /

COPY devenv/bin/manager /manager

ENTRYPOINT ["/manager"]
