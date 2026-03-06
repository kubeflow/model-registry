---
name: kind-dev-setup
description: Set up a local Kind cluster with Colima and run the model-registry dev environment. Starts infrastructure (Colima + Tilt), BFF (Go backend), and Frontend (React). Includes optional setups for performance data, MinIO S3 storage, and OCI model transfer jobs. Use when the user asks to set up Kind, start the dev environment, run the local cluster, start the model-registry UI locally, deploy performance data, set up MinIO, or test model transfer jobs.
---

# Kind Dev Environment Setup

Set up a local Kubernetes dev environment for model-registry using Colima, Kind, and Tilt, then run the BFF and Frontend.

## Prerequisites

- Docker, Colima, Kind, Tilt, Go, Node.js installed
- Workspace at the model-registry repo root

## Core Setup (3 Terminals)

### Terminal 1: Infrastructure

1. Start Colima (if not running):

```bash
colima start
```

2. Verify/switch kubectl context to `kind-model-registry`:

```bash
kubectl config use-context kind-model-registry
```

3. Start Tilt:

```bash
cd devenv && ./bin/tilt up
```

If port 10350 is occupied: `lsof -ti:10350 | xargs kill -9`

4. Wait for Tilt resources to become healthy before proceeding.

### Terminal 2: BFF

Run the Go BFF server from `clients/ui/bff`:

**With mock K8s client (no cluster needed for basic UI testing):**

```bash
cd clients/ui/bff
go run ./cmd --port=4000 --dev-mode --dev-mode-model-registry-port=8080 --dev-mode-catalog-port=8082 --deployment-mode=standalone --mock-k8s-client
```

**With real K8s client (requires running Kind cluster):**

```bash
cd clients/ui/bff
go run ./cmd --port=4000 --dev-mode --dev-mode-model-registry-port=8080 --dev-mode-catalog-port=8082 --deployment-mode=standalone
```

If port 4000 is occupied: `lsof -ti:4000 | xargs kill -9`

### Terminal 3: Frontend

```bash
cd clients/ui/frontend
DEPLOYMENT_MODE=standalone STYLE_THEME=patternfly npm run start:dev
```

The frontend dev server proxies API requests to the BFF on port 4000.

### RBAC Setup (real K8s only)

When using the real K8s client, the BFF's namespace registry access check uses SubjectAccessReview with only the `User` field (no groups), so group-based RoleBindings don't take effect. Create a ClusterRoleBinding that directly binds each namespace's default ServiceAccount:

```bash
kubectl apply -f - <<EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: model-registry-all-sa-service-access
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: model-registry-ui-services-reader
subjects:
- kind: ServiceAccount
  name: default
  namespace: default
- kind: ServiceAccount
  name: default
  namespace: kubeflow
- kind: ServiceAccount
  name: default
  namespace: minio
- kind: ServiceAccount
  name: default
  namespace: kube-system
- kind: ServiceAccount
  name: default
  namespace: kube-public
- kind: ServiceAccount
  name: default
  namespace: kube-node-lease
- kind: ServiceAccount
  name: default
  namespace: local-path-storage
EOF
```

Without this, the UI shows "The selected namespace does not have access to this model registry" for every namespace.

### Switching Between Mock and Real K8s

Kill the BFF process in Terminal 2 and restart with or without `--mock-k8s-client`. When switching to real K8s, ensure the RBAC ClusterRoleBinding above is applied.

## Optional: Performance Data (Model Catalog)

Deploy the Model Catalog with demo performance data (evaluations, metrics) using the provided script:

```bash
./scripts/deploy_catalog_demo_on_kind.sh
```

This script:
- Creates/reuses the `model-registry` Kind cluster
- Creates the `model-catalog` namespace
- Deploys the catalog with the `demo` kustomize overlay including perf data
- Waits for Postgres and the catalog server to be ready

The demo overlay loads performance data from `manifests/kustomize/options/catalog/overlays/demo/perf-data/` which contains evaluation and performance ndjson files for certified models.

To access the catalog API after deployment:

```bash
kubectl port-forward -n model-catalog svc/model-catalog-server 8082:8080
```

Then available at `http://localhost:8082`. The BFF connects to this on `--dev-mode-catalog-port=8082`.

## Optional: MinIO (S3 Storage)

Deploy MinIO for testing transfer jobs with S3 sources. All resources must be applied to the `minio` namespace explicitly (the manifests don't specify one):

```bash
kubectl create namespace minio --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -n minio -f scripts/manifests/minio/deployment.yaml
kubectl wait --for=condition=available deployment/minio -n minio --timeout=120s
kubectl apply -n minio -f scripts/manifests/minio/create_bucket.yaml
```

**Important:** MinIO has no persistent volume — data is lost on pod restart. If the `minio-init` job already shows as "Complete" but the bucket is gone, delete and re-run:

```bash
kubectl delete job minio-init -n minio --ignore-not-found
kubectl apply -n minio -f scripts/manifests/minio/create_bucket.yaml
```

| Detail | Value |
|--------|-------|
| Internal endpoint | `http://minio.minio.svc.cluster.local:9000` |
| Bucket | `default` |
| Credentials | `minioadmin` / `minioadmin` |
| Console NodePort | `30091` |
| K8s Secret | `minio-secret` (namespace: `minio`) |

Upload test data:

```bash
kubectl run minio-upload --rm -i --restart=Never -n minio \
  --image=minio/mc --command -- sh -c '
mc --config-dir /tmp alias set local http://minio:9000 minioadmin minioadmin
echo "sample model content" | mc --config-dir /tmp pipe local/default/models/sample-model/model.txt
'
```

## Optional: OCI Model Transfer Jobs

Test S3-to-OCI model transfer jobs end-to-end. Requires MinIO (above) and a destination OCI registry (e.g. quay.io).

No local ARM64 image build is needed — upstream, midstream, and downstream each have their own async-upload images.

### Create a Transfer Job via UI

Use these values when filling the form:

| Field | Value | Notes |
|-------|-------|-------|
| Source type | `s3` | |
| S3 endpoint | `http://minio.minio:9000` | Internal cluster DNS |
| S3 bucket | `default` | |
| S3 key | `models/sample-model/` | Directory prefix, **not** full file path |
| S3 access key | `minioadmin` | |
| S3 secret key | `minioadmin` | |
| Destination type | `oci` | |
| Destination URI | `quay.io/yourorg/yourrepo:tag` | OCI ref format, **no** `https://` |
| Destination registry | `quay.io` | |

### Key Gotchas

- **S3 key must be a directory prefix** (e.g. `models/dir/`), not a full file path. Using the exact file path causes an EBUSY error because `os.path.relpath` resolves to `.`.
- **Destination URI must be an OCI reference** (`quay.io/org/repo:tag`), not a web URL (`https://quay.io/repository/...`). The upload code prepends `docker://`.
- Sample job manifests are in `jobs/async-upload/samples/` (`create_model_example.yaml`, `create_version_example.yaml`, `sample_job_s3_to_oci.yaml`).

## Troubleshooting

| Problem | Fix |
|---------|-----|
| Tilt refuses to start (production context) | `kubectl config use-context kind-model-registry` |
| Port conflict (4000 or 10350) | `lsof -ti:PORT \| xargs kill -9` |
| Frontend proxy `ECONNREFUSED` | BFF not ready yet; wait for it to start |
| `ImagePullBackOff` on async-upload job | Verify the correct image is configured for your environment (upstream/midstream/downstream each have their own) |
| "namespace does not have access to this model registry" | Apply the RBAC ClusterRoleBinding (see above). The BFF's SAR uses `User` only, not groups. |
| MinIO nodePort 30091 already allocated | MinIO already exists in `minio` namespace. Don't apply without `-n minio` or it creates a duplicate in `default`. |
| MinIO bucket missing after pod restart | MinIO has no PV. Delete the old job and re-apply: `kubectl delete job minio-init -n minio && kubectl apply -n minio -f scripts/manifests/minio/create_bucket.yaml` |
| `envtest` port lock error in Go tests | `rm -f ~/Library/Caches/kubebuilder-envtest/port-*` |
| Transfer job S3 download fails with EBUSY | Use directory prefix as source key, not full file path |
| Transfer job OCI push "invalid reference" | Use OCI ref format `quay.io/org/repo:tag`, not web URL |
