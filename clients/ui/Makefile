CONTAINER_TOOL ?= docker

.PHONY: all
all: build

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

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

.PHONY: docker-compose
docker-compose:
	$(CONTAINER_TOOL) compose -f docker-compose.yaml up

.PHONY: kind-deployment
kind-deployment:
	./scripts/deploy_kind_cluster.sh
