.PHONY: deploy-mcp-catalog-odh
deploy-mcp-catalog-odh:
	@echo "Deploying MCP catalog config on ODH..."
	bash ../../../scripts/deploy_mcp_catalog_on_odh.sh

.PHONY: undeploy-mcp-catalog-odh
undeploy-mcp-catalog-odh:
	@echo "Undeploying MCP catalog config on ODH..."
	bash ../../../scripts/undeploy_catalog_on_odh.sh

.PHONY: test-e2e-catalog
test-e2e-catalog: deploy-mcp-catalog-odh
	@echo "Running catalog tests..."
	export MC_NAMESPACE=$$(kubectl get datasciencecluster default-dsc -o jsonpath='{.spec.components.modelregistry.registriesNamespace}') && \
	export CATALOG_URL="https://$$(kubectl get route -n "$$MC_NAMESPACE" model-catalog-https -o 'jsonpath={.status.ingress[0].host}')/" && \
	export AUTH_TOKEN=$$(kubectl config view --raw -o jsonpath="{.users[?(@.name==\"$$(kubectl config view -o jsonpath="{.contexts[?(@.name==\"$$(kubectl config current-context)\")].context.user}")\")].user.token}") && \
	export VERIFY_SSL=False && \
	export KIND_CLUSTER=False && \
	poetry install --all-extras && poetry run pytest tests --e2e -svv -rA