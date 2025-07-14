# Build the model-registry binary
FROM --platform=$BUILDPLATFORM registry.access.redhat.com/ubi8/go-toolset:1.24.4 AS common
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY ["go.mod", "go.sum", "./"]
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Download tools
COPY Makefile .
COPY ["main.go", ".openapi-generator-ignore", "openapitools.json", "./"]
COPY cmd/ cmd/
COPY api/ api/
COPY internal/ internal/
COPY scripts/ scripts/
COPY pkg/ pkg/
COPY templates/ templates/
COPY patches/ patches/
COPY catalog/ catalog/

###### Dev stage - start ######
# see: https://github.com/kubeflow/model-registry/pull/984#discussion_r2048732415

FROM common AS dev

USER root

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o model-registry

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest AS dev-build

WORKDIR /
COPY --from=dev /workspace/model-registry .
USER 65532:65532

ENTRYPOINT ["/model-registry"]

###### Dev stage - end ######

FROM common AS builder

USER root
# default NodeJS 14 is not enough for openapi-generator-cli, switch to Node JS currently supported
RUN yum remove -y nodejs npm
RUN yum module -y reset nodejs
RUN yum module -y enable nodejs:18
# install npm and java for openapi-generator-cli
RUN yum install -y nodejs npm java-11 python3

RUN make deps

# NOTE: The two instructions below are effectively equivalent to 'make clean build'
# DO NOT REMOVE THE 'build/prepare' TARGET!!!
# It ensures consitent repeatable Dockerfile builds

# prepare the build in a separate layer
RUN make clean build/prepare
# compile separately to optimize multi-platform builds
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} make build/compile

# Use distroless as minimal base image to package the model-registry binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
WORKDIR /
# copy the registry binary
COPY --from=builder /workspace/model-registry .
USER 65532:65532

ENTRYPOINT ["/model-registry"]
