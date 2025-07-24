JOB_IMG_VERSION ?= latest
JOB_IMG_REGISTRY ?= quay.io
JOB_IMG_ORG ?= opendatahub
JOB_IMG_NAME ?= model-registry-job-async-upload
JOB_IMG ?= $(JOB_IMG_REGISTRY)/$(JOB_IMG_ORG)/$(JOB_IMG_NAME):$(JOB_IMG_VERSION)
BUILD_IMAGE ?= true # whether to build the MR server image
CLUSTER_NAME ?= mr-e2e

# MR Server Params
IMG_VERSION ?= latest
# keep IMG consistent with root Makefile and ci/GHA:
IMG ?= quay.io/opendatahub/model-registry

.PHONY: deploy-latest-mr
deploy-latest-mr:
	cd ../../ && \
	$(if $(filter true,$(BUILD_IMAGE)),\
		IMG_VERSION=${IMG_VERSION} IMG=${IMG} make image/build ARGS="--load$(if ${DEV_BUILD}, --target dev-build)" && \
	) \
	LOCAL=1 IMG=$(IMG):$(IMG_VERSION) ./scripts/deploy_on_kind.sh
# TODO RHOAIENG-30453 align consistency ./scripts/deploy_on_kind.sh uses IMG with :tag, Vs, Makefile(s) and ci/GHA we use IMG without trailing :tag
	kubectl port-forward -n kubeflow services/model-registry-service 8080:8080 & echo $$! >> .port-forwards.pid

.PHONY: deploy-test-minio
deploy-test-minio:
	cd ../../ && ./scripts/deploy_minio_on_kind.sh
	kubectl port-forward -n minio svc/minio 9000:9000 & echo $$! >> .port-forwards.pid
	kubectl port-forward -n minio svc/minio 9001:9001 & echo $$! >> .port-forwards.pid

.PHONY: deploy-local-registry
deploy-local-registry:
	cd ../../ && ./scripts/deploy_local_kind_registry.sh
	kubectl port-forward service/distribution-registry-test-service 5001:5001 & echo $$! >> .port-forwards.pid


.PHONY: dev-load-image
dev-load-image:
	docker buildx build --load -t $(JOB_IMG) .
	kind load docker-image $(JOB_IMG) -n $(CLUSTER_NAME)

.PHONY: test
test:
	poetry run pytest -s -x -rA

test-with-cov:
	poetry run pytest --cov=job --cov-report=term-missing tests/

.PHONY: test-e2e
test-e2e: deploy-latest-mr deploy-local-registry deploy-test-minio
	@echo "Starting test-e2e"
	-$(MAKE) test-e2e-run; STATUS=$$?
	$(MAKE) test-e2e-cleanup
	@exit $$STATUS

.PHONY: test-e2e-run
test-e2e-run:
	@echo "Ensuring all extras are installed..."
	poetry install --all-extras
	@echo "Running tests..."
	poetry run pytest --e2e -s -x -rA

.PHONY: test-e2e-cleanup
test-e2e-cleanup:
	@echo "Cleaning up port-forward processes..."
	@if [ -f .port-forwards.pid ]; then \
		kill $$(cat .port-forwards.pid) || true; \
		rm -f .port-forwards.pid; \
	fi

.PHONY: test-integration
test-integration: deploy-latest-mr deploy-local-registry deploy-test-minio dev-load-image
	@echo "Starting test-integration"
	-$(MAKE) test-integration-run; STATUS=$$?
	$(MAKE) test-e2e-cleanup
	@exit $$STATUS

.PHONY: test-integration-run
test-integration-run:
	@echo "Ensuring all extras are installed..."
	poetry install --all-extras --with integration
	@echo "Running integration tests..."
	CONTAINER_IMAGE_URI=$(JOB_IMG) poetry run pytest --integration tests/integration/ -vs

.PHONY: install
install:
	poetry install
