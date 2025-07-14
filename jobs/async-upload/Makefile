IMG_VERSION ?= latest
IMG_REGISTRY ?= ghcr.io
IMG_ORG ?= kubeflow
IMG_NAME ?= model-registry/job/async-upload
IMG ?= $(IMG_REGISTRY)/$(IMG_ORG)/$(IMG_NAME):$(IMG_VERSION)
BUILD_IMAGE ?= true # whether to build the MR server image

.PHONY: deploy-latest-mr
deploy-latest-mr:
	cd ../../ && \
	$(if $(filter true,$(BUILD_IMAGE)),\
		IMG_VERSION=${IMG_VERSION} make image/build ARGS="--load$(if ${DEV_BUILD}, --target dev-build)" && \
	) \
	LOCAL=1 ./scripts/deploy_on_kind.sh
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
	docker buildx build --load -t $(IMG) .
	kind load docker-image $(IMG) -n mr-e2e

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

.PHONY: install
install:
	poetry install
