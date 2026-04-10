#!/usr/bin/env bash
# Run a kind cluster with Model Catalog demo overlay (full performance data).
# Prerequisites: Docker (or Colima) running, kind and kubectl installed.

set -e

CATALOG_NAMESPACE="${CATALOG_NAMESPACE:-model-catalog}"
CLUSTER_NAME="${CLUSTER_NAME:-model-registry}"
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

cd "$REPO_ROOT"

echo "=== Creating kind cluster (if needed) ==="
if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
  echo "Cluster ${CLUSTER_NAME} already exists."
  kubectl config use-context "kind-${CLUSTER_NAME}"
else
  kind create cluster --name "$CLUSTER_NAME"
fi

echo "=== Creating namespace ${CATALOG_NAMESPACE} ==="
kubectl create namespace "$CATALOG_NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

echo "=== Deploying Model Catalog with demo overlay (full performance data) ==="
kubectl apply -k manifests/kustomize/options/catalog/overlays/demo -n "$CATALOG_NAMESPACE"

echo "=== Waiting for Model Catalog Postgres to be ready ==="
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=postgres,app.kubernetes.io/part-of=model-catalog -n "$CATALOG_NAMESPACE" --timeout=120s 2>/dev/null || true
echo "=== Waiting for Model Catalog server (with perf data) to be available ==="
kubectl wait --for=condition=available deployment/model-catalog-server -n "$CATALOG_NAMESPACE" --timeout=5m

echo "=== Done ==="
kubectl get pods -n "$CATALOG_NAMESPACE"
echo ""
echo "To access the Model Catalog API (with performance metrics):"
echo "  kubectl port-forward -n $CATALOG_NAMESPACE svc/model-catalog-server 8080:8080"
echo "Then open http://localhost:8080 (or use the API with performance-metrics from /perf-data)."
