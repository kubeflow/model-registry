IMG ?= quay.io/${USER}/model-registry-storage-initializer:latest

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests.
	go test ./... -coverprofile cover.out

##@ Build

.PHONY: build
build: fmt vet ## Build binary.
	go build -o bin/mr-storage-initializer main.go

.PHONY: run
run: fmt vet ## Run the program
	go run ./main.go $(SOURCE_URI) $(DEST_PATH)

.PHONY: docker-build
docker-build: test ## Build docker image.
	docker build . -f ./Dockerfile -t ${IMG}

.PHONY: docker-push
docker-push: ## Push docker image.
	docker push ${IMG}