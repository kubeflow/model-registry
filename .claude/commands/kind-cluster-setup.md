---
description: Set up local Kind cluster for model-registry UI development with optional catalog, MinIO, and OCI transfer support
allowed-tools: Bash, Read, Write
argument-hint: [all|catalog|perf-data|minio|oci-transfer|real-k8s]
---

# Kind Cluster Setup

Set up the local Kind dev environment for model-registry UI development. Execute the setup based on the arguments provided.

## Arguments

`$ARGUMENTS` specifies which components to set up. If empty, run core setup only with mock K8s.

Available options (space-separated):
- **catalog** — Deploy Model Catalog with demo overlay via `./scripts/deploy_catalog_demo_on_kind.sh`
- **perf-data** — Same as `catalog` (demo overlay includes performance data)
- **minio** — Deploy MinIO S3 storage for testing transfer jobs
- **oci-transfer** — Full OCI transfer job setup (includes MinIO + ARM64 image build + Kind image load)
- **real-k8s** — Use real K8s client instead of mock (omit `--mock-k8s-client` from BFF)
- **all** — Everything: core + catalog + minio + oci-transfer + real-k8s

## Core Setup (always runs)

Run these in separate terminals:

### Terminal 1 - Infrastructure
```
colima start
kubectl config use-context kind-model-registry
cd devenv && ./bin/tilt up
```

### Terminal 2 - BFF
```
cd clients/ui/bff
go run ./cmd --port=4000 --dev-mode --dev-mode-model-registry-port=8080 --dev-mode-catalog-port=8082 --deployment-mode=standalone --mock-k8s-client
```
If `real-k8s` or `all` is in `$ARGUMENTS`, omit `--mock-k8s-client`.

### Terminal 3 - Frontend
```
cd clients/ui/frontend
DEPLOYMENT_MODE=standalone STYLE_THEME=patternfly npm run start:dev
```

## Optional: catalog / perf-data
```
./scripts/deploy_catalog_demo_on_kind.sh
kubectl port-forward -n model-catalog svc/model-catalog-server 8082:8080
```

## Optional: minio
```
kubectl apply -f scripts/manifests/minio/deployment.yaml
kubectl wait --for=condition=available deployment/minio -n minio --timeout=120s
kubectl apply -f scripts/manifests/minio/create_bucket.yaml
```
Then upload test data:
```
kubectl run minio-upload --rm -i --restart=Never -n minio --image=minio/mc --command -- sh -c 'mc --config-dir /tmp alias set local http://minio:9000 minioadmin minioadmin && echo "sample model content" | mc --config-dir /tmp pipe local/default/models/sample-model/model.txt'
```

## Optional: oci-transfer (implies minio)
Run minio setup above first, then:
```
cd jobs/async-upload
docker build --platform linux/arm64 -t quay.io/opendatahub/model-registry-job-async-upload:latest .
kind load docker-image quay.io/opendatahub/model-registry-job-async-upload:latest --name model-registry
```

## Examples

- `/kind-cluster-setup` — Core only
- `/kind-cluster-setup catalog` — Core + Model Catalog
- `/kind-cluster-setup minio oci-transfer` — Core + MinIO + OCI transfer
- `/kind-cluster-setup all` — Everything
