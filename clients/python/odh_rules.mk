.PHONY: deploy-mr-odh
deploy-mr-odh:
	cd ../../ && ./scripts/deploy_on_odh.sh

.PHONY: undeploy-mr-odh
undeploy-mr-odh:
	cd ../../ && ./scripts/undeploy_on_odh.sh

.PHONY: test-e2e-odh
test-e2e-odh:
	@echo "Ensuring all extras are installed..."
	poetry install --all-extras
	@echo "Running tests..."
	@set -a; . ../../scripts/manifests/minio/.env; set +a; \
	export MR_NAMESPACE=$$(kubectl get datasciencecluster default-dsc -o jsonpath='{.spec.components.modelregistry.registriesNamespace}' 2>/dev/null) && poetry run pytest --e2e -s -rA \
	&& rm -f ../../scripts/manifests/minio/.env
