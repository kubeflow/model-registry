# useful paths
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
PROJECT_BIN := $(PROJECT_PATH)/bin
GO ?= "$(shell which go)"
UI_PATH := $(PROJECT_PATH)/clients/ui
CSI_PATH := $(PROJECT_PATH)/cmd/csi
CONTROLLER_PATH := $(PROJECT_PATH)/cmd/controller

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.29
ENVTEST ?= $(PROJECT_BIN)/setup-envtest

# add tools bin directory
PATH := $(PROJECT_BIN):$(PATH)

MLMD_VERSION ?= 1.14.0

# docker executable
DOCKER ?= docker
# default Dockerfile
DOCKERFILE ?= Dockerfile
# container registry, default to github container registry
IMG_REGISTRY ?= ghcr.io
# container image organization
IMG_ORG ?= kubeflow
# container image version
IMG_VERSION ?= main
# container image repository
IMG_REPO ?= model-registry/server
# container image build path
BUILD_PATH ?= .
# container image
ifdef IMG_REGISTRY
    IMG := ${IMG_REGISTRY}/${IMG_ORG}/${IMG_REPO}
else
    IMG := ${IMG_ORG}/${IMG_REPO}
endif

# Change Dockerfile path depending on IMG_REPO
ifeq ($(IMG_REPO),model-registry/ui)
    DOCKERFILE := $(UI_PATH)/Dockerfile
	BUILD_PATH := $(UI_PATH)
endif

# The BUILD_PATH is still the root
ifeq ($(IMG_REPO),model-registry/storage-initializer)
    DOCKERFILE := $(CSI_PATH)/Dockerfile.csi
endif

# The BUILD_PATH is still the root
ifeq ($(IMG_REPO),model-registry/controller)
    DOCKERFILE := $(CONTROLLER_PATH)/Dockerfile.controller
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

YQ_EXPR := 'sort_keys(.components.schemas) | sort_keys(.paths) | sort_keys(.components.responses)'

.PHONY: fmt/openapi
fmt/openapi: api/openapi/model-registry.yaml bin/yq
	@$(YQ) -i $(YQ_EXPR) $<

# validate the openapi schema
.PHONY: openapi/validate
openapi/validate: api/openapi/model-registry.yaml bin/openapi-generator-cli bin/yq
	@$(YQ) $(YQ_EXPR) $< | diff -u $< - || (echo "$< is incorrectly formatted. Run 'make fmt/openapi' to fix it."; exit 1)
	$(OPENAPI_GENERATOR) validate -i $<

# generate the openapi server implementation
.PHONY: gen/openapi-server
gen/openapi-server: bin/openapi-generator-cli openapi/validate internal/server/openapi/api_model_registry_service.go

internal/server/openapi/api_model_registry_service.go: bin/openapi-generator-cli api/openapi/model-registry.yaml
	ROOT_FOLDER=${PROJECT_PATH} ./scripts/gen_openapi_server.sh

# generate the openapi schema model and client
.PHONY: gen/openapi
gen/openapi: bin/openapi-generator-cli openapi/validate pkg/openapi/client.go

pkg/openapi/client.go: bin/openapi-generator-cli api/openapi/model-registry.yaml clean-pkg-openapi
	${OPENAPI_GENERATOR} generate \
		-i api/openapi/model-registry.yaml -g go -o pkg/openapi --package-name openapi \
		--ignore-file-override ./.openapi-generator-ignore --additional-properties=isGoSubmodule=true,enumClassPrefix=true,useOneOfDiscriminatorLookup=true
	gofmt -w pkg/openapi

.PHONY: vet
vet:
	${GO} vet ./...

.PHONY: clean/csi
clean/csi:
	rm -Rf ./mr-storage-initializer

.PHONY: clean-pkg-openapi
clean-pkg-openapi:
	while IFS= read -r file; do rm -f "pkg/openapi/$$file"; done < pkg/openapi/.openapi-generator/FILES

.PHONY: clean-internal-server-openapi
clean-internal-server-openapi:
	while IFS= read -r file; do rm -f "internal/server/openapi/$$file"; done < internal/server/openapi/.openapi-generator/FILES

.PHONY: clean
clean: clean-pkg-openapi clean-internal-server-openapi clean/csi
	rm -Rf ./model-registry internal/ml_metadata/proto/*.go internal/converter/generated/*.go

.PHONY: clean/odh
clean/odh:
	rm -Rf ./model-registry

bin/protoc:
	./scripts/install_protoc.sh

bin/go-enum:
	GOBIN=$(PROJECT_BIN) ${GO} install github.com/searKing/golang/tools/go-enum@v1.2.97

bin/protoc-gen-go:
	GOBIN=$(PROJECT_BIN) ${GO} install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0

bin/protoc-gen-go-grpc:
	GOBIN=$(PROJECT_BIN) ${GO} install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

bin/envtest:
	GOBIN=$(PROJECT_BIN) ${GO} install sigs.k8s.io/controller-runtime/tools/setup-envtest@v0.0.0-20240320141353-395cfc7486e6

GOLANGCI_LINT ?= ${PROJECT_BIN}/golangci-lint
bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v2.0.2

GOVERTER ?= ${PROJECT_BIN}/goverter
bin/goverter:
	GOBIN=$(PROJECT_PATH)/bin ${GO} install github.com/jmattheis/goverter/cmd/goverter@v1.8.1

YQ ?= ${PROJECT_BIN}/yq
bin/yq:
	GOBIN=$(PROJECT_PATH)/bin ${GO} install github.com/mikefarah/yq/v4@v4.45.1

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
deps: bin/protoc bin/go-enum bin/protoc-gen-go bin/protoc-gen-go-grpc bin/golangci-lint bin/goverter bin/openapi-generator-cli bin/envtest

.PHONY: vendor
vendor:
	${GO} mod vendor

# WARNING: DO NOT DELETE THIS TARGET, USED BY Dockerfile!!!
.PHONY: build/prepare
build/prepare: gen vet lint

# WARNING: DO NOT DELETE THIS TARGET, USED BY Dockerfile!!!
.PHONY: build/compile
build/compile:
	${GO} build -buildvcs=false

# WARNING: DO NOT EDIT THIS TARGET DIRECTLY!!!
# Use build/prepare to add build prerequisites
# Use build/compile to add/edit go source compilation
# WARNING: Editing this target directly WILL affect the Dockerfile image build!!!
.PHONY: build
build: build/prepare build/compile

.PHONY: build/odh
build/odh: vet
	${GO} build -buildvcs=false

.PHONY: build/prepare/csi
build/prepare/csi: build/prepare lint/csi

.PHONY: build/compile/csi
build/compile/csi:
	${GO} build -buildvcs=false -o mr-storage-initializer ${CSI_PATH}/main.go

.PHONY: build/csi
build/csi: build/prepare/csi build/compile/csi

.PHONY: gen
gen: deps gen/grpc gen/openapi gen/openapi-server gen/converter
	${GO} generate ./...

.PHONY: lint
lint: bin/golangci-lint
	${GOLANGCI_LINT} run main.go  --timeout 3m
	${GOLANGCI_LINT} run cmd/... internal/... ./pkg/...  --timeout 3m

.PHONY: lint/csi
lint/csi: bin/golangci-lint
	${GOLANGCI_LINT} run ${CSI_PATH}/main.go
	${GOLANGCI_LINT} run internal/csi/...

.PHONY: test
test: gen bin/envtest
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" ${GO} test ./internal/... ./pkg/...

.PHONY: test-nocache
test-nocache: gen bin/envtest
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" ${GO} test ./internal/... ./pkg/... -count=1

.PHONY: test-cover
test-cover: gen bin/envtest
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" ${GO} test ./internal/... ./pkg/... -coverprofile=coverage.txt
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
	${DOCKER} build ${BUILD_PATH} -f ${DOCKERFILE} -t ${IMG}:$(IMG_VERSION) $(ARGS)

# build docker image using buildx
# PLATFORMS defines the target platforms for the model registry image be built to provide support to multiple
# architectures. (i.e. make docker-buildx). To use this option you need to:
# - be able to use docker buildx. More info: https://docs.docker.com/build/buildx/
# - have enabled BuildKit. More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image to your registry (i.e. if you do not set a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To adequately provide solutions that are compatible with multiple platforms, you should consider using this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: image/buildx
image/buildx:
ifeq ($(DOCKER),docker)
	# docker uses builder containers
	- $(DOCKER) buildx rm model-registry-builder
	$(DOCKER) buildx create --use --name model-registry-builder --platform=$(PLATFORMS)
	$(DOCKER) buildx build --push --platform=$(PLATFORMS) --tag ${IMG}:$(IMG_VERSION) -f ${DOCKERFILE} .
	$(DOCKER) buildx rm model-registry-builder
else ifeq ($(DOCKER),podman)
	# podman uses image manifests
	$(DOCKER) manifest create -a ${IMG}
	$(DOCKER) buildx build --platform=$(PLATFORMS) --manifest ${IMG}:$(IMG_VERSION) -f ${DOCKERFILE} .
	$(DOCKER) manifest push ${IMG}
	$(DOCKER) manifest rm ${IMG}
else
	$(error Unsupported container tool $(DOCKER))
endif

# push docker image
.PHONY: image/push
image/push:
	${DOCKER} push ${IMG}:$(IMG_VERSION)

all: model-registry

##  ------------------------------- ##
##  ----  Controller Targets   ---- ##
##  ------------------------------- ##

##@ Development

.PHONY: controller/manifests
controller/manifests: bin/controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=model-registry-manager-role crd webhook paths="{./cmd/controller/..., ./internal/controller/...}" output:crd:artifacts:config=manifests/options/controller/crd/bases output:rbac:dir=manifests/kustomize/options/controller/rbac

.PHONY: controller/generate
controller/generate: bin/controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="./cmd/controller/hack/boilerplate.go.txt" paths="{./cmd/controller/..., ./internal/controller/...}"

.PHONY: controller/fmt
controller/fmt: ## Run go fmt against code.
	go fmt ./cmd/controller/... ./internal/controller/...

.PHONY: controller/vet
controller/vet: ## Run go vet against code.
	go vet ./cmd/controller/... ./internal/controller/...

.PHONY: controller/test
controller/test: controller/manifests controller/generate controller/fmt controller/vet bin/envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(PROJECT_BIN) -p path)" go test $$(go list ./internal/controller/... | grep -v /e2e) -coverprofile cover.out

##@ Build

.PHONY: controller/build
controller/build: controller/manifests controller/generate controller/fmt controller/vet ## Build manager binary.
	go build -o bin/manager cmd/controller/main.go

.PHONY: controller/run
controller/run: controller/manifests controller/generate controller/fmt controller/vet ## Run a controller from your host.
	go run ./cmd/controller/main.go

# If you wish to build the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: controller/docker-build
controller/docker-build: ## Build docker image with the manager.
	$(DOCKER) build -t ${IMG} -f ./cmd/controller/Dockerfile.controller .

.PHONY: controller/docker-push
controller/docker-push: ## Push docker image with the manager.
	$(DOCKER) push ${IMG}

# PLATFORMS defines the target platforms for the manager image be built to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/mypoperator:0.0.1). To use this option you need to:
# - be able to use docker buildx. More info: https://docs.docker.com/build/buildx/
# - have enabled BuildKit. More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image to your registry (i.e. if you do not set a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To adequately provide solutions that are compatible with multiple platforms, you should consider using this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: controller/docker-buildx
controller/docker-buildx: ## Build and push docker image for the manager for cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' ./cmd/controller/Dockerfile.controller > Dockerfile.cross
	- $(DOCKER) buildx create --name controller-builder
	$(DOCKER) buildx use controller-builder
	- $(DOCKER) buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross .
	- $(DOCKER) buildx rm controller-builder
	rm Dockerfile.cross

.PHONY: controller/build-installer
controller/build-installer: controller/manifests controller/generate bin/kustomize ## Generate a consolidated YAML with CRDs and deployment.
	mkdir -p dist
	cd manifests/kustomize/options/controller/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build manifests/kustomize/options/controller/default > dist/install.yaml

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: controller/install
controller/install: controller/manifests bin/kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build manifests/kustomize/options/controller/crd | $(KUBECTL) apply -f -

.PHONY: controller/uninstall
controller/uninstall: controller/manifests bin/kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build manifests/kustomize/options/controller/crd | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: controller/deploy
controller/deploy: controller/manifests bin/kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd manifests/kustomize/options/controller/manager && $(KUSTOMIZE) edit set image ghcr.io/kubeflow/model-registry/controller=${IMG}:${IMG_VERSION}
	$(KUSTOMIZE) build manifests/kustomize/options/controller/overlays/base | $(KUBECTL) apply -f -

.PHONY: controller/undeploy
controller/undeploy: bin/kustomize ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build manifests/kustomize/options/controller/overlays/base | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

##@ Tools

KUBECTL ?= kubectl
CONTROLLER_GEN ?= $(PROJECT_BIN)/controller-gen
KUSTOMIZE ?= $(PROJECT_BIN)/kustomize
CONTROLLER_TOOLS_VERSION ?= v0.16.4
KUSTOMIZE_VERSION ?= v5.5.0

.PHONY: bin/kustomize
bin/kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(PROJECT_BIN)
	$(call go-install-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v5,$(KUSTOMIZE_VERSION))

.PHONY: bin/controller-gen
bin/controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(PROJECT_BIN)
	$(call go-install-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen,$(CONTROLLER_TOOLS_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(PROJECT_BIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef
