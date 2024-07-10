# Build the model-registry binary
FROM registry.access.redhat.com/ubi8/go-toolset:1.21 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY ["go.mod", "go.sum", "./"]
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

USER root

# Copy the go source
COPY ["Makefile", "main.go", "./"]

# Copy rest of the source
COPY bin/ bin/
COPY pkg/ pkg/

# Build
USER root
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 make build

# Use distroless as minimal base image to package the model-registry storage initializer binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
WORKDIR /
# copy the storage initializer binary
COPY --from=builder /workspace/bin/mr-storage-initializer .
USER 65532:65532

ENTRYPOINT ["/mr-storage-initializer"]