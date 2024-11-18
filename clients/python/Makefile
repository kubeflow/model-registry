all: install tidy

IMG_REGISTRY ?= quay.io
IMG_VERSION ?= latest 

.PHONY: install
install:
	../../bin/openapi-generator-cli generate -i ../../api/openapi/model-registry.yaml -g python -o src/ --package-name mr_openapi --additional-properties=library=asyncio,generateSourceCodeOnly=true,useOneOfDiscriminatorLookup=true
	mv src/mr_openapi{_,/}README.md
	git apply patches/*
	poetry install

.PHONY: clean
clean:
	rm -rf src/mr_openapi

.PHONY: deploy-latest-mr
deploy-latest-mr:
	cd ../../ && IMG_VERSION=${IMG_VERSION} make image/build && LOCAL=1 ./scripts/deploy_on_kind.sh
	kubectl port-forward -n kubeflow services/model-registry-service 8080:8080 &

.PHONY: test-e2e
test-e2e: deploy-latest-mr
	poetry run pytest --e2e -s -rA

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
