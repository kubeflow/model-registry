.PHONY: deploy-mr-odh
deploy-mr-odh:
	cd ../../ && ./scripts/deploy_on_odh.sh

.PHONY: undeploy-mr-odh
undeploy-mr-odh:
	cd ../../ && ./scripts/undeploy_on_odh.sh

.PHONY: test-e2e-odh
test-e2e-odh: deploy-mr-odh deploy-local-registry deploy-test-minio
	@echo "Running e2e tests..."
	@set -a; . ../../scripts/manifests/minio/.env; set +a; \
	mkdir -p ../../results; \
	export AUTH_TOKEN=$$(kubectl config view --raw -o jsonpath="{.users[?(@.name==\"$$(kubectl config view -o jsonpath="{.contexts[?(@.name==\"$$(kubectl config current-context)\")].context.user}")\")].user.token}") && \
	export VERIFY_SSL=False && \
	export MR_NAMESPACE=$$(kubectl get datasciencecluster default-dsc -o jsonpath='{.spec.components.modelregistry.registriesNamespace}') && \
	export MR_HOST_URL="https://$$(kubectl get service -n "$$MR_NAMESPACE" model-registry -o jsonpath='{.metadata.annotations.routing\.opendatahub\.io\/external-address-rest}')" && \
	export MR_ENDPOINT=$$(kubectl get service -n "$$MR_NAMESPACE" model-registry -o jsonpath='{.metadata.annotations.routing\.opendatahub\.io\/external-address-rest}' | cut -d: -f1) && \
	export MODEL_SYNC_REGISTRY_SERVER_ADDRESS="https://$$MR_ENDPOINT" && \
	export MODEL_SYNC_REGISTRY_PORT="443" && \
	export MODEL_SYNC_REGISTRY_IS_SECURE="false" && \
	export MODEL_SYNC_REGISTRY_USER_TOKEN="$$AUTH_TOKEN" && \
	export CONTAINER_IMAGE_URI=$$(../../scripts/get_async_upload_image.sh) && \
	poetry install --all-extras --with integration && poetry run pytest --e2e tests/integration/ -svvv -rA --html=../../results/report.html --junit-xml=../../results/xunit_report.xml --self-contained-html -o junit_suite_name=odh-async-upload && \
	rm -f ../../scripts/manifests/minio/.env

.PHONY: test-e2e-port-cleanup
test-e2e-odh-cleanup: undeploy-mr-odh undeploy-minio undeploy-local-kind-registry
	@echo "Cleaning up port-forward processes..."
	@if [ -f .port-forwards.pid ]; then \
		kill $$(cat .port-forwards.pid) || true; \
		rm -f .port-forwards.pid; \
	fi
