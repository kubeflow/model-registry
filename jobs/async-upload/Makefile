IMG_VERSION ?= latest

.PHONY: deploy-latest-mr
deploy-latest-mr:
	cd ../../ && IMG_VERSION=${IMG_VERSION} make image/build ARGS="--load$(if ${DEV_BUILD}, --target dev-build)"  && LOCAL=1 ./scripts/deploy_on_kind.sh
	kubectl port-forward -n kubeflow services/model-registry-service 8080:8080 & echo $$! >> .port-forwards.pid


.PHONY: test
test:
	poetry run pytest -s -rA

test-with-cov:
	poetry run pytest --cov=job --cov-report=term-missing tests/

.PHONY: test-e2e
test-e2e: deploy-latest-mr
	@echo "Starting test-e2e"
	-$(MAKE) test-e2e-run; STATUS=$$?
	$(MAKE) test-e2e-cleanup
	@exit $$STATUS

.PHONY: test-e2e-run
test-e2e-run:
	@echo "Ensuring all extras are installed..."
	poetry install --all-extras
	@echo "Running tests..."
	poetry run pytest -m e2e -s -rA

.PHONY: test-e2e-cleanup
test-e2e-cleanup:
	@echo "Cleaning up port-forward processes..."
	@if [ -f .port-forwards.pid ]; then \
		kill $$(cat .port-forwards.pid) || true; \
		rm -f .port-forwards.pid; \
	fi