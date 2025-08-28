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
	export AUTH_TOKEN=$$(kubectl config view --raw -o jsonpath="{.users[?(@.name==\"$$(kubectl config view -o jsonpath="{.contexts[?(@.name==\"$$(kubectl config current-context 2>/dev/null)\")].context.user}" 2>/dev/null)\")].user.token}" 2>/dev/null) && \
	export VERIFY_SSL=False && \
	export MR_NAMESPACE=$$(kubectl get datasciencecluster default-dsc -o jsonpath='{.spec.components.modelregistry.registriesNamespace}' 2>/dev/null) && \
    export MR_URL="https://$$(kubectl get service -n "$$MR_NAMESPACE" model-registry -o jsonpath='{.metadata.annotations.routing\.opendatahub\.io\/external-address-rest}' 2>/dev/null)" && poetry run pytest --e2e -s -rA \
	&& rm -f ../../scripts/manifests/minio/.env

.PHONY: test-e2e-port-cleanup
test-e2e-port-cleanup:
	@echo "Cleaning up port-forward processes..."
	@if [ -f .port-forwards.pid ]; then \
		kill $$(cat .port-forwards.pid) || true; \
		rm -f .port-forwards.pid; \
	fi
