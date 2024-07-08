.PHONY: install
install:
	../../bin/openapi-generator-cli generate -i ../../api/openapi/model-registry.yaml -g python -o src/ --package-name mr_openapi --additional-properties=library=asyncio,generateSourceCodeOnly=true,useOneOfDiscriminatorLookup=true
	mv src/mr_openapi{_,/}README.md
	git apply patches/*
	poetry install

.PHONY: clean
clean:
	rm -rf src/mr_openapi

.PHONY: test
test:
	poetry run pytest -s

.PHONY: lint-check
lint-check:
	poetry run ruff check
	poetry run black src/mr_openapi --check

.PHONY: lint-apply
lint-apply:
	poetry run ruff check --fix --unsafe-fixes || true
	poetry run black src/mr_openapi

.PHONY: build
build: install lint-apply
	poetry build

.PHONY: publish
publish: install
	poetry publish --build -u ${PYPI_USERNAME} -p ${PYPI_PASSWORD}

.PHONY: update
update:
	poetry lock
