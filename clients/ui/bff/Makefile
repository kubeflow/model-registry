CONTAINER_TOOL ?= docker
IMG ?= model-registry-bff:latest

.PHONY: all
all: build

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: clean
clean:
	rm -Rf ./bin

.PHONY: vet
vet:  .
	go vet ./...

.PHONY: test
test:
	go test ./...

.PHONY: build
build: fmt vet test
	go build -o bin/bff cmd/main.go

.PHONY: run
run: fmt vet
	PORT=4000 go run ./cmd/main.go

.PHONY: docker-build
docker-build:
	$(CONTAINER_TOOL) build -t ${IMG} .