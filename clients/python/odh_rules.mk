.PHONY: deploy-mr-odh
deploy-mr-odh:
	cd ../../ && ./scripts/deploy_on_odh.sh

.PHONY: undeploy-mr-odh
undeploy-mr-odh:
	cd ../../ && ./scripts/undeploy_on_odh.sh

.PHONY: test-e2e-odh
test-e2e-odh:
	@echo "Running tests..."
	@set -a; . ../../scripts/manifests/minio/.env; set +a; \
	mkdir -p ../../results; \
	export AUTH_TOKEN=$$(kubectl config view --raw -o jsonpath="{.users[?(@.name==\"$$(kubectl config view -o jsonpath="{.contexts[?(@.name==\"$$(kubectl config current-context)\")].context.user}")\")].user.token}") && \
	export VERIFY_SSL=False && \
	export MR_NAMESPACE=$$(kubectl get datasciencecluster default-dsc -o jsonpath='{.spec.components.modelregistry.registriesNamespace}') && \
	export MR_URL="https://$$(kubectl get service -n "$$MR_NAMESPACE" model-registry -o jsonpath='{.metadata.annotations.routing\.opendatahub\.io\/external-address-rest}')" && \
	poetry install --all-extras && poetry run pytest --e2e -svvv -rA --html=../../results/report.html --junit-xml=../../results/xunit_report.xml --self-contained-html && \
	rm -f ../../scripts/manifests/minio/.env

.PHONY: test-e2e-port-cleanup
test-e2e-port-cleanup:
	@echo "Cleaning up port-forward processes..."
	@if [ -f .port-forwards.pid ]; then \
		kill $$(cat .port-forwards.pid) || true; \
		rm -f .port-forwards.pid; \
	fi
