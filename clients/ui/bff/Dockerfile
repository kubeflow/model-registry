# Use the golang image to build the application
FROM golang:1.22.2 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /ui

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the go source files
COPY cmd/ cmd/
COPY api/ api/
COPY config/ config/
COPY data/ data/
COPY integrations/ integrations/
COPY internals/ internals/
COPY validation/ validation/



# Build the Go application
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o bff ./cmd/main.go

# Use distroless as minimal base image to package the application binary
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder ui/bff ./
USER 65532:65532

# Expose port 4000
EXPOSE 4000

# Define environment variables
ENV PORT 4001
ENV ENV development

ENTRYPOINT ["/bff"]
