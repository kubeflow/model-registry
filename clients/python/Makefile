all: install tidy

IMG_REGISTRY ?= quay.io
IMG_VERSION ?= latest
IMG ?= quay.io/opendatahub/model-registry:${IMG_VERSION}
BUILD_IMAGE ?= true # whether to build the MR server image
include odh_rules.mk

.PHONY: install
install:
	../../bin/openapi-generator-cli generate -i ../../api/openapi/model-registry.yaml -g python -o src/ --package-name mr_openapi --additional-properties=library=asyncio,generateSourceCodeOnly=true,useOneOfDiscriminatorLookup=true --template-dir templates/
	mv src/mr_openapi_README.md src/mr_openapi/README.md
	git apply patches/*
	poetry install

.PHONY: clean
clean:
	rm -rf src/mr_openapi

.PHONY: deploy-latest-mr
deploy-latest-mr:
	cd ../../ && \
	$(if $(filter true,$(BUILD_IMAGE)),\
		IMG_VERSION=${IMG_VERSION} IMG=${IMG} make image/build ARGS="--load$(if ${DEV_BUILD}, --target dev-build)" && \
	) \
	LOCAL=1 ./scripts/deploy_on_kind.sh
	kubectl port-forward -n kubeflow services/model-registry-service 8080:8080 & echo $$! >> .port-forwards.pid

.PHONY: deploy-test-minio
deploy-test-minio:
	cd ../../ && ./scripts/deploy_minio_on_kind.sh
	kubectl port-forward -n minio svc/minio 9000:9000 & echo $$! >> .port-forwards.pid

.PHONY: deploy-local-registry
deploy-local-registry:
	cd ../../ && ./scripts/deploy_local_kind_registry.sh
	kubectl port-forward service/distribution-registry-test-service 5001:5001 & echo $$! >> .port-forwards.pid

.PHONY: test-e2e
test-e2e: deploy-latest-mr deploy-local-registry deploy-test-minio
	@echo "Starting test-e2e"
	-$(MAKE) test-e2e-run; STATUS=$$?
	$(MAKE) test-e2e-cleanup
	@exit $$STATUS

.PHONY: test-fuzz
test-fuzz: deploy-latest-mr deploy-local-registry deploy-test-minio
	@echo "Starting test-fuzz"
	poetry install --all-extras
	@set -a; . ../../scripts/manifests/minio/.env; set +a; \
	poetry run pytest --fuzz -v -s --hypothesis-show-statistics
	@rm -f ../../scripts/manifests/minio/.env
	$(MAKE) test-e2e-cleanup
	@exit $$STATUS

.PHONY: test-e2e-run
test-e2e-run:
	@echo "Ensuring all extras are installed..."
	poetry install --all-extras
	@echo "Running tests..."
	@set -a; . ../../scripts/manifests/minio/.env; set +a; \
	poetry run pytest --e2e -s -rA --html=report.html --self-contained-html
	@rm -f ../../scripts/manifests/minio/.env

.PHONY: test-e2e-cleanup
test-e2e-cleanup:
	@echo "Cleaning up database..."
	cd ../../ && ./scripts/cleanup.sh
	@echo "Cleaning up port-forward processes..."
	@if [ -f .port-forwards.pid ]; then \
		kill $$(cat .port-forwards.pid) || true; \
		rm -f .port-forwards.pid; \
	fi

.PHONY: undeploy-minio
undeploy-minio:
	@echo "Undeploy Minio..."
	cd ../../ && ./scripts/undeploy_minio.sh

.PHONY: undeploy-local-kind-registry
undeploy-local-kind-registry:
	@echo "Undeploy Local Kind Registry..."
	cd ../../ && ./scripts/undeploy_local_registry.sh

.PHONY: test
test:
	poetry run pytest -s -rA

.PHONY: lint
lint:
	poetry run ruff check src/model_registry

.PHONY: tidy
tidy:
	poetry run ruff check --fix --unsafe-fixes || true
	poetry run black src/mr_openapi

.PHONY: build
build: install tidy
	poetry build

.PHONY: publish
publish: build
	poetry publish --build -u ${PYPI_USERNAME} -p ${PYPI_PASSWORD}

.PHONY: update
update:
	poetry lock
