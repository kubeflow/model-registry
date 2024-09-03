all: install tidy

IMG_REGISTRY ?= docker.io
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

.PHONY: build-mr
build-mr:
	cd ../../ && IMG_REGISTRY=${IMG_REGISTRY} IMG_VERSION=${IMG_VERSION} make image/build

.PHONY: test-e2e
test-e2e: build-mr
	poetry run pytest --e2e -s

.PHONY: test
test:
	poetry run pytest -s

.PHONY: lint
lint:
	poetry run ruff check
	poetry run black src/mr_openapi --check

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
