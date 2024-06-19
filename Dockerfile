# Build the model-registry binary
FROM registry.access.redhat.com/ubi8/go-toolset:1.21 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY ["go.mod", "go.sum", "./"]
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

USER root
# default NodeJS 14 is not enough for openapi-generator-cli, switch to Node JS currently supported
RUN yum remove -y nodejs npm
RUN yum module -y reset nodejs
RUN yum module -y enable nodejs:18
# install npm and java for openapi-generator-cli
RUN yum install -y nodejs npm java-11

# Copy the go source
COPY ["Makefile", "main.go", ".openapi-generator-ignore", "openapitools.json", "./"]

# Copy rest of the source
COPY .git/ .git/
COPY cmd/ cmd/
COPY api/ api/
COPY internal/ internal/
COPY scripts/ scripts/
COPY pkg/ pkg/
COPY patches/ patches/
COPY templates/ templates/

# Download tools
RUN make deps

# Build
USER root
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 make clean model-registry

# Use distroless as minimal base image to package the model-registry binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
WORKDIR /
# copy the registry binary
COPY --from=builder /workspace/model-registry .
USER 65532:65532

ENTRYPOINT ["/model-registry"]
