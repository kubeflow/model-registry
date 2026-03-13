# Kind Cluster Setup

Set up the local Kind dev environment for model-registry UI development. Execute the setup based on the parameters provided after this command. Run each terminal as a separate backgrounded shell process.

## Arguments

Prerequisites: Docker, Colima, Kind, kubectl, Go >= 1.25.7, Node.js >= 22.0.0 (Tilt is auto-downloaded by `make tilt-up`).

Options (space-separated after the command). If none specified, run core setup only with real K8s client.

- **catalog** — Deploy Model Catalog with demo overlay (performance data included)
- **perf-data** — Same as `catalog` (demo overlay includes performance data)
- **minio** — Deploy MinIO S3 storage for testing transfer jobs
- **mock-k8s** — Use mock K8s client instead of real (add `--mock-k8s-client` to BFF). Skips RBAC setup.
- **all** — Everything: core + catalog + minio (with real K8s)

## Core Setup (always runs)

### Terminal 1 - Infrastructure
```
colima start
kind get clusters | grep -q '^model-registry$' || kind create cluster --name model-registry
kubectl config use-context kind-model-registry
cd devenv && make tilt-up
```

### Terminal 2 - BFF
```
cd clients/ui/bff
go run ./cmd --port=4000 --dev-mode --dev-mode-model-registry-port=8080 --dev-mode-catalog-port=8082 --deployment-mode=standalone
```
If `mock-k8s` is specified, add `--mock-k8s-client`.

### RBAC Setup (real K8s only, skip if `mock-k8s`)

The BFF's namespace registry access check uses SubjectAccessReview with only the `User` field (no groups), so group-based RoleBindings don't take effect. Create a ClusterRoleBinding that directly binds each namespace's default ServiceAccount:

```
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
Without this, the UI shows "The selected namespace does not have access to this model registry."

### Terminal 3 - Frontend
```
cd clients/ui/frontend
DEPLOYMENT_MODE=standalone STYLE_THEME=patternfly npm run start:dev
```

## Optional: catalog / perf-data
```
kubectl create namespace model-catalog --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -k manifests/kustomize/options/catalog/overlays/demo -n model-catalog
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=postgres,app.kubernetes.io/part-of=model-catalog -n model-catalog --timeout=120s || true
kubectl wait --for=condition=available deployment/model-catalog-server -n model-catalog --timeout=5m
kubectl port-forward -n model-catalog svc/model-catalog-server 8082:8080
```

## Optional: minio

Create namespace and deploy (all resources must go to the `minio` namespace):
```
kubectl create namespace minio --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -n minio -f scripts/manifests/minio/deployment.yaml
kubectl wait --for=condition=available deployment/minio -n minio --timeout=120s
kubectl apply -n minio -f scripts/manifests/minio/create_bucket.yaml
```

MinIO has no persistent volume — data is lost on pod restart. If the `minio-init` job already exists but the bucket is gone, delete and re-run:
```
kubectl delete job minio-init -n minio --ignore-not-found
kubectl apply -n minio -f scripts/manifests/minio/create_bucket.yaml
```

Upload test data:
```
kubectl run minio-upload --rm -i --restart=Never -n minio --image=minio/mc --command -- sh -c 'mc --config-dir /tmp alias set local http://minio:9000 minioadmin minioadmin && echo "sample model content" | mc --config-dir /tmp pipe local/default/models/sample-model/model.txt'
```

## Examples

- `/kind-cluster-setup` — Core only (real K8s)
- `/kind-cluster-setup mock-k8s` — Core with mock K8s (no RBAC needed)
- `/kind-cluster-setup catalog` — Core + Model Catalog
- `/kind-cluster-setup minio` — Core + MinIO for transfer jobs
- `/kind-cluster-setup all` — Everything
