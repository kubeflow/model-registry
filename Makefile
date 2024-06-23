# useful paths
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
PROJECT_BIN := $(PROJECT_PATH)/bin
GO := $(PROJECT_BIN)/go1.21.9

# add tools bin directory
PATH := $(PROJECT_BIN):$(PATH)

MLMD_VERSION ?= 1.14.0

# docker executable
DOCKER ?= docker
# default Dockerfile
DOCKERFILE ?= Dockerfile
# container registry, default to empty (dockerhub) if not explicitly set
IMG_REGISTRY ?= quay.io
# container image organization
IMG_ORG ?= opendatahub
# container image version
IMG_VERSION ?= main
# container image repository
IMG_REPO ?= model-registry
# container image
ifdef IMG_REGISTRY
    IMG := ${IMG_REGISTRY}/${IMG_ORG}/${IMG_REPO}
else
    IMG := ${IMG_ORG}/${IMG_REPO}
endif

model-registry: build

# clean the ml-metadata protos and trigger a fresh new build which downloads
# ml-metadata protos based on specified MLMD_VERSION
.PHONY: update/ml_metadata
update/ml_metadata: clean/ml_metadata clean build

clean/ml_metadata:
	rm -rf api/grpc/ml_metadata/proto/*.proto

api/grpc/ml_metadata/proto/metadata_source.proto:
	mkdir -p api/grpc/ml_metadata/proto/
	cd api/grpc/ml_metadata/proto/ && \
		curl -LO "https://raw.githubusercontent.com/google/ml-metadata/v${MLMD_VERSION}/ml_metadata/proto/metadata_source.proto" && \
		sed -i 's#syntax = "proto[23]";#&\noption go_package = "github.com/kubeflow/model-registry/internal/ml_metadata/proto";#' metadata_source.proto

api/grpc/ml_metadata/proto/metadata_store.proto:
	mkdir -p api/grpc/ml_metadata/proto/
	cd api/grpc/ml_metadata/proto/ && \
		curl -LO "https://raw.githubusercontent.com/google/ml-metadata/v${MLMD_VERSION}/ml_metadata/proto/metadata_store.proto" && \
		sed -i 's#syntax = "proto[23]";#&\noption go_package = "github.com/kubeflow/model-registry/internal/ml_metadata/proto";#' metadata_store.proto

api/grpc/ml_metadata/proto/metadata_store_service.proto:
	mkdir -p api/grpc/ml_metadata/proto/
	cd api/grpc/ml_metadata/proto/ && \
		curl -LO "https://raw.githubusercontent.com/google/ml-metadata/v${MLMD_VERSION}/ml_metadata/proto/metadata_store_service.proto" && \
		sed -i 's#syntax = "proto[23]";#&\noption go_package = "github.com/kubeflow/model-registry/internal/ml_metadata/proto";#' metadata_store_service.proto

internal/ml_metadata/proto/%.pb.go: api/grpc/ml_metadata/proto/%.proto
	bin/protoc -I./api/grpc --go_out=./internal --go_opt=paths=source_relative \
		--go-grpc_out=./internal --go-grpc_opt=paths=source_relative $<

.PHONY: gen/grpc
gen/grpc: internal/ml_metadata/proto/metadata_store.pb.go internal/ml_metadata/proto/metadata_store_service.pb.go

internal/converter/generated/converter.go: internal/converter/*.go
	${GOVERTER} gen github.com/kubeflow/model-registry/internal/converter/

.PHONY: gen/converter
gen/converter: gen/grpc internal/converter/generated/converter.go

# validate the openapi schema
.PHONY: openapi/validate
openapi/validate: bin/openapi-generator-cli
	${OPENAPI_GENERATOR} validate -i api/openapi/model-registry.yaml

# generate the openapi server implementation
.PHONY: gen/openapi-server
gen/openapi-server: bin/openapi-generator-cli openapi/validate
	@if git diff --exit-code --name-only | grep -q "api/openapi/model-registry.yaml" || \
		git diff --exit-code --name-only | grep -q "api/openapi/model-registry.yaml" || \
		[ -n "${FORCE_SERVER_GENERATION}" ]; then \
		ROOT_FOLDER="." ./scripts/gen_openapi_server.sh; \
	else \
		echo "INFO api/openapi/model-registry.yaml is not staged or modified, will not re-generate server"; \
	fi

# generate the openapi schema model and client
.PHONY: gen/openapi
gen/openapi: bin/openapi-generator-cli openapi/validate pkg/openapi/client.go

pkg/openapi/client.go: bin/openapi-generator-cli api/openapi/model-registry.yaml
	rm -rf pkg/openapi
	${OPENAPI_GENERATOR} generate \
		-i api/openapi/model-registry.yaml -g go -o pkg/openapi --package-name openapi \
		--ignore-file-override ./.openapi-generator-ignore --additional-properties=isGoSubmodule=true,enumClassPrefix=true,useOneOfDiscriminatorLookup=true
	gofmt -w pkg/openapi

.PHONY: vet
vet:
	${GO} vet ./...

.PHONY: clean
clean:
	rm -Rf ./model-registry internal/ml_metadata/proto/*.go internal/converter/generated/*.go pkg/openapi

.PHONY: clean/odh
clean/odh:
	rm -Rf ./model-registry

bin/go:
	GOBIN=$(PROJECT_BIN) go install golang.org/dl/go1.21.9@latest
	$(PROJECT_BIN)/go1.21.9 download

bin/protoc:
	./scripts/install_protoc.sh

bin/go-enum:
	GOBIN=$(PROJECT_BIN) ${GO} install github.com/searKing/golang/tools/go-enum@v1.2.97

bin/protoc-gen-go:
	GOBIN=$(PROJECT_BIN) ${GO} install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0

bin/protoc-gen-go-grpc:
	GOBIN=$(PROJECT_BIN) ${GO} install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

GOLANGCI_LINT ?= ${PROJECT_BIN}/golangci-lint
bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v1.54.2

GOVERTER ?= ${PROJECT_BIN}/goverter
bin/goverter:
	GOBIN=$(PROJECT_PATH)/bin ${GO} install github.com/jmattheis/goverter/cmd/goverter@v1.1.1

OPENAPI_GENERATOR ?= ${PROJECT_BIN}/openapi-generator-cli
NPM ?= "$(shell which npm)"
bin/openapi-generator-cli:
ifeq (, $(shell which ${NPM} 2> /dev/null))
	@echo "npm is not available please install it to be able to install openapi-generator"
	exit 1
endif
ifeq (, $(shell which ${PROJECT_BIN}/openapi-generator-cli 2> /dev/null))
	@{ \
	set -e ;\
	mkdir -p ${PROJECT_BIN} ;\
	mkdir -p ${PROJECT_BIN}/openapi-generator-installation ;\
	cd ${PROJECT_BIN} ;\
	${NPM} install -g --prefix ${PROJECT_BIN}/openapi-generator-installation @openapitools/openapi-generator-cli ;\
	ln -s openapi-generator-installation/bin/openapi-generator-cli openapi-generator-cli ;\
	}
endif

.PHONY: clean/deps
clean/deps:
	rm -Rf bin/*

.PHONY: deps
deps: bin/go bin/protoc bin/go-enum bin/protoc-gen-go bin/protoc-gen-go-grpc bin/golangci-lint bin/goverter bin/openapi-generator-cli

.PHONY: vendor
vendor:
	${GO} mod vendor

.PHONY: build
build: gen vet lint
	${GO} build -buildvcs=false

.PHONY: build/odh
build/odh: vet
	${GO} build -buildvcs=false

.PHONY: gen
gen: deps gen/grpc gen/openapi gen/openapi-server gen/converter
	${GO} generate ./...

.PHONY: lint
lint:
	${GOLANGCI_LINT} run main.go
	${GOLANGCI_LINT} run cmd/... internal/... ./pkg/...

.PHONY: test
test: gen
	${GO} test ./internal/... ./pkg/...

.PHONY: test-nocache
test-nocache: gen
	${GO} test ./internal/... ./pkg/... -count=1

.PHONY: test-cover
test-cover: gen
	${GO} test ./internal/... ./pkg/... -coverprofile=coverage.txt
	${GO} tool cover -html=coverage.txt -o coverage.html

.PHONY: run/proxy
run/proxy: gen
	${GO} run main.go proxy --logtostderr=true

.PHONY: proxy
proxy: build
	./model-registry proxy --logtostderr=true

# login to docker
.PHONY: docker/login
docker/login:
ifdef IMG_REGISTRY
	$(DOCKER) login -u "${DOCKER_USER}" -p "${DOCKER_PWD}" "${IMG_REGISTRY}"
else
	$(DOCKER) login -u "${DOCKER_USER}" -p "${DOCKER_PWD}"
endif


# build docker image
.PHONY: image/build
image/build:
	${DOCKER} build . -f ${DOCKERFILE} -t ${IMG}:$(IMG_VERSION)

# push docker image
.PHONY: image/push
image/push:
	${DOCKER} push ${IMG}:$(IMG_VERSION)

all: model-registry
