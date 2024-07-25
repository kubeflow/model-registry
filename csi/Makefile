# Useful paths
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
PROJECT_BIN := $(PROJECT_PATH)/bin
GO := $(PROJECT_BIN)/go1.21.9

# add tools bin directory
PATH := $(PROJECT_BIN):$(PATH)

# container tool
DOCKER ?= docker
DOCKERFILE ?= Dockerfile
DOCKERFILE_DEV ?= Dockerfile.dev

IMG_REGISTRY ?= quay.io
# container image organization
IMG_ORG ?= opendatahub
# container image version
IMG_VERSION ?= main
# container image repository
IMG_REPO ?= model-registry-storage-initializer
# container image
ifdef IMG_REGISTRY
    IMG := ${IMG_REGISTRY}/${IMG_ORG}/${IMG_REPO}
else
    IMG := ${IMG_ORG}/${IMG_REPO}
endif

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

bin-go: ## Install go tool
	GOBIN=$(PROJECT_BIN) go install golang.org/dl/go1.21.9@latest
	$(PROJECT_BIN)/go1.21.9 download

.PHONY: bin
bin: bin-go ## Install required dependencies

.PHONY: tidy
tidy: ## Run go mod tidy.
	${GO} mod tidy

.PHONY: fmt
fmt: bin ## Run go fmt against code.
	${GO} fmt ./...

.PHONY: vet
vet: bin ## Run go vet against code.
	${GO} vet ./...

.PHONY: test
test: fmt vet bin ## Run tests.
	${GO} test ./... -coverprofile cover.out

##@ Build

.PHONY: build
build: fmt vet bin ## Build binary.
	${GO} build -o bin/mr-storage-initializer main.go

.PHONY: run
run: fmt vet ## Run the program
	${GO} run ./main.go $(SOURCE_URI) $(DEST_PATH)

.PHONY: docker-build
docker-build: ## Build container image.
	${DOCKER} build . -f ./$(DOCKERFILE) -t ${IMG}:$(IMG_VERSION)

.PHONY: docker-build-dev
docker-build-dev: ## Build container image using local model-registry module.
	cd ../ && ${DOCKER} build . -f ./csi/$(DOCKERFILE_DEV) -t ${IMG}:$(IMG_VERSION)

.PHONY: docker-push
docker-push: ## Push container image.
	${DOCKER} push ${IMG}:$(IMG_VERSION)