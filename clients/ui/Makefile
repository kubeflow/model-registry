DEFAULT_ENV_FILE := .env
ifneq ("$(wildcard $(DEFAULT_ENV_FILE))","")
include ${DEFAULT_ENV_FILE}
export $(shell sed 's/=.*//' ${DEFAULT_ENV_FILE})
endif

DEV_ENV_FILE := .env.development
ifneq ("$(wildcard $(DEV_ENV_FILE))","")
include ${DEV_ENV_FILE}
export $(shell sed 's/=.*//' ${DEV_ENV_FILE})
endif

ENV_FILE := .env.local
ifneq ("$(wildcard $(ENV_FILE))","")
include ${ENV_FILE}
export $(shell sed 's/=.*//' ${ENV_FILE})
endif

.PHONY: all
all: build

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


############ Dev Environment ############

.PHONY: dev-install-dependencies
dev-install-dependencies:
	cd frontend && npm install

.PHONY: dev-bff
dev-bff:
	cd bff && make run PORT=4000 MOCK_K8S_CLIENT=true MOCK_MR_CLIENT=true DEV_MODE=true STANDALONE_MODE=true

.PHONY: dev-frontend
dev-frontend:
	cd frontend && npm run start:dev

.PHONY: dev-start
dev-start: 
	make -j 2 dev-bff dev-frontend

########### Dev Integrated ############
.PHONY: dev-start-kubeflow
dev-start-kubeflow: 
	make -j 2 dev-bff-kubeflow dev-frontend-kubeflow

.PHONY: dev-frontend-kubeflow
dev-frontend-kubeflow:
	DEPLOYMENT_MODE=integrated && cd frontend && npm run start:dev

.PHONY: dev-bff-kubeflow
dev-bff-kubeflow:
	cd bff && make run PORT=4000 MOCK_K8S_CLIENT=false MOCK_MR_CLIENT=false DEV_MODE=true STANDALONE_MODE=false DEV_MODE_PORT=8085

############ Build ############

.PHONY: docker-build
docker-build:
	$(CONTAINER_TOOL) build -t ${IMG_UI} .

.PHONY: docker-build-standalone
docker-build-standalone:
	$(CONTAINER_TOOL) build --build-arg DEPLOYMENT_MODE=standalone -t ${IMG_UI_STANDALONE} .

.PHONY: docker-buildx
docker-buildx:
	docker buildx build --platform ${PLATFORM} -t ${IMG_UI} --push .

.PHONY: docker-buildx-standalone
docker-buildx-standalone:
	docker buildx build --build-arg DEPLOYMENT_MODE=standalone --platform ${PLATFORM} -t ${IMG_UI_STANDALONE} --push .

############ Push ############

.PHONY: docker-push
docker-push:
	${CONTAINER_TOOL} push ${IMG_UI}

.PHONY: docker-push-standalone
docker-push-standalone:
	${CONTAINER_TOOL} push ${IMG_UI_STANDALONE}

############ Deployment ############

.PHONY: kind-deployment
kind-deployment:
	./scripts/deploy_kind_cluster.sh

.PHONY: kubeflow-deployment
kubeflow-deployment:
	./scripts/deploy_kubeflow_cluster.sh

############ Build ############	
.PHONY: frontend-build
frontend-build:
	cd frontend && npm run build:prod

.PHONY: frontend-build-standalone
frontend-build-standalone:
	cd frontend && DEPLOYMENT_MODE=standalone npm run build:prod

.PHONY: bff-build
bff-build:
	cd bff && make build

.PHONY: build
build: frontend-build bff-build

############ Run mocked ########
.PHONY: run-local-mocked
run-local-mocked: frontend-build-standalone bff-build
	rm -r ./bff/static-local-run && cp -r ./frontend/dist/ ./bff/static-local-run/ && cd bff && make run STATIC_ASSETS_DIR=./static-local-run MOCK_K8S_CLIENT=true DEV_MODE=true


	