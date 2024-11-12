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
	cd bff && make run PORT=4000 MOCK_K8S_CLIENT=true MOCK_MR_CLIENT=true

.PHONY: dev-frontend
dev-frontend:
	cd frontend && npm run start:dev

.PHONY: dev-start
dev-start: 
	make -j 2 dev-bff dev-frontend

############ Build ############

.PHONY: build-bff
build-bff:
	$(CONTAINER_TOOL) build -t ${IMG_BFF} ./bff

.PHONY: build-frontend
build-frontend:
	$(CONTAINER_TOOL) build -t ${IMG_FRONTEND} ./frontend

.PHONY: build
build: build-bff build-frontend

############ Push ############

.PHONY: push-bff
push-bff:
	${CONTAINER_TOOL} push ${IMG_BFF}

.PHONY: push-frontend
push-frontend:
	${CONTAINER_TOOL} push ${IMG_FRONTEND}

.PHONY: push
push: push-bff push-frontend

############ Deployment ############

.PHONY: docker-compose
docker-compose:
	$(CONTAINER_TOOL) compose -f docker-compose.yaml up

.PHONY: kind-deployment
kind-deployment:
	./scripts/deploy_kind_cluster.sh
