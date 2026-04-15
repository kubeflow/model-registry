#!/usr/bin/env bash
# Tear down the model-registry dev environment.
# Stops Frontend, BFF, and Tilt. Optionally deletes the Kind cluster.
#
# Usage:
#   ./scripts/dev_teardown.sh                  # Stop everything, delete cluster
#   ./scripts/dev_teardown.sh --keep-cluster   # Stop processes, keep cluster
#
# Environment variables (override default ports):
#   FRONTEND_PORT  (default: 9000)
#   BFF_PORT       (default: 4000)

set -e

FRONTEND_PORT="${FRONTEND_PORT:-9000}"
BFF_PORT="${BFF_PORT:-4000}"
CLUSTER_NAME="${CLUSTER_NAME:-model-registry}"
KEEP_CLUSTER=false
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

for arg in "$@"; do
  case "$arg" in
    --keep-cluster) KEEP_CLUSTER=true ;;
  esac
done

cd "$REPO_ROOT"

echo "=== Stopping Frontend (port ${FRONTEND_PORT}) ==="
lsof -ti:"${FRONTEND_PORT}" | xargs kill -9 2>/dev/null || true

echo "=== Stopping BFF (port ${BFF_PORT}) ==="
lsof -ti:"${BFF_PORT}" | xargs kill -9 2>/dev/null || true

echo "=== Stopping Tilt ==="
cd devenv && make tilt-down 2>/dev/null || true
cd "$REPO_ROOT"

if [ "$KEEP_CLUSTER" = false ]; then
  echo "=== Deleting Kind cluster '${CLUSTER_NAME}' ==="
  kind delete cluster --name "$CLUSTER_NAME"

  if [ $(which colima 2>/dev/null) ]; then
    echo "=== Stopping Colima ==="
    colima stop 2>/dev/null
  fi
else
  echo "=== Keeping Kind cluster running ==="
fi

echo "=== Done ==="
