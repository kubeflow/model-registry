all: install tidy

IMG_REGISTRY ?= quay.io
IMG_VERSION ?= latest 

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
	cd ../../ && IMG_VERSION=${IMG_VERSION} make image/build && LOCAL=1 ./scripts/deploy_on_kind.sh
	kubectl port-forward -n kubeflow services/model-registry-service 8080:8080 &

.PHONY: deploy-test-minio
deploy-test-minio:
	cd ../../ && ./scripts/deploy_minio_on_kind.sh
	kubectl port-forward -n minio svc/minio 9000:9000 &

.PHONY: deploy-local-registry
deploy-local-registry:
	cd ../../ && ./scripts/deploy_local_kind_registry.sh

.PHONY: test-e2e
test-e2e: deploy-latest-mr deploy-local-registry deploy-test-minio
	@set -a; . ../../scripts/manifests/minio/.env; set +a; \
	poetry run pytest --e2e -s -rA
	rm -fr ../../scripts/manifests/minio/.env

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
