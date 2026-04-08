#!/usr/bin/env bash

set -e

DSC_NAME="default-dsc"

echo "Check if DataScienceCluster exists"
if kubectl get datasciencecluster "$DSC_NAME" &> /dev/null; then
  echo "DataScienceCluster '$DSC_NAME' exists."
else
  echo "DataScienceCluster '$DSC_NAME' does NOT exist."
  exit 1
fi

echo "Check if Model Registry is enabled in DSC"
MR_STATE=$(kubectl get datasciencecluster "$DSC_NAME" -o jsonpath='{.spec.components.modelregistry.managementState}' 2>/dev/null)
if [ "$MR_STATE" != "Managed" ]; then
  echo "Model Registry is not enabled (managementState='$MR_STATE'). Expected 'Managed'."
  exit 1
fi
echo "Model Registry is enabled (managementState='Managed')."

MR_NAMESPACE=$(kubectl get datasciencecluster "$DSC_NAME" -o jsonpath='{.spec.components.modelregistry.registriesNamespace}' 2>/dev/null)
if [ -z "$MR_NAMESPACE" ]; then
  echo "Could not determine registriesNamespace from DSC."
  exit 1
fi
echo "Model Registry namespace: '$MR_NAMESPACE'"

echo "Looking for catalog sources ConfigMap in namespace '$MR_NAMESPACE'"
if kubectl get configmap mcp-catalog-sources -n "$MR_NAMESPACE" &> /dev/null; then
  CATALOG_CONFIGMAP="mcp-catalog-sources"
  echo "ConfigMap 'mcp-catalog-sources' found."
elif kubectl get configmap model-catalog-sources -n "$MR_NAMESPACE" &> /dev/null; then
  CATALOG_CONFIGMAP="model-catalog-sources"
  echo "ConfigMap 'model-catalog-sources' found (fallback)."
else
  echo "Neither 'mcp-catalog-sources' nor 'model-catalog-sources' ConfigMap found in namespace '$MR_NAMESPACE'."
  exit 1
fi

SCRIPT_DIR="$(dirname "$(realpath "$BASH_SOURCE")")"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
MCP_SERVERS_FILE="${REPO_ROOT}/manifests/kustomize/options/catalog/overlays/e2e/test-mcp-servers.yaml"

if [ ! -f "$MCP_SERVERS_FILE" ]; then
  echo "Required file not found: $MCP_SERVERS_FILE"
  exit 1
fi

echo "Fetching current sources.yaml from ConfigMap"
CURRENT_SOURCES=$(kubectl get configmap "$CATALOG_CONFIGMAP" -n "$MR_NAMESPACE" -o jsonpath='{.data.sources\.yaml}')

if echo "$CURRENT_SOURCES" | grep -q "test_mcp_servers"; then
  echo "test_mcp_servers already present in sources.yaml, skipping patch."
else
  echo "Patching ConfigMap to add mcp_catalogs, labels, and namedQueries"

  # Read the MCP servers test data file content
  MCP_SERVERS_CONTENT=$(cat "$MCP_SERVERS_FILE")

  # Replace empty mcp_catalogs: [] and append test config
  UPDATED_SOURCES=$(echo "$CURRENT_SOURCES" | sed 's/mcp_catalogs: \[\]//')
  UPDATED_SOURCES="${UPDATED_SOURCES}
mcp_catalogs:
  - name: Test MCP Servers
    id: test_mcp_servers
    type: yaml
    enabled: true
    properties:
      yamlCatalogPath: test-mcp-servers.yaml
    labels:
      - Test MCP Servers

labels:
  - name: mcp-label
    assetType: mcp_servers

namedQueries:
  production_ready:
    assetType: mcp_servers
    filters:
      verifiedSource:
        operator: '='
        value: true
  security_focused:
    assetType: mcp_servers
    filters:
      sast:
        operator: '='
        value: true
      readOnlyTools:
        operator: '='
        value: true
"

  # Patch the existing ConfigMap: update sources.yaml and add the MCP servers YAML as a data key
  kubectl patch configmap "$CATALOG_CONFIGMAP" -n "$MR_NAMESPACE" --type merge -p "$(cat <<EOF
{"data": {"sources.yaml": $(echo "$UPDATED_SOURCES" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))'), "test-mcp-servers.yaml": $(echo "$MCP_SERVERS_CONTENT" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))')}}
EOF
)"

  echo "ConfigMap '$CATALOG_CONFIGMAP' patched successfully."

  echo "Restarting model-catalog pod to pick up config changes"
  CATALOG_DEPLOYMENT=$(kubectl get deployment -n "$MR_NAMESPACE" -l app.kubernetes.io/name=model-catalog -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
  if [ -z "$CATALOG_DEPLOYMENT" ]; then
    echo "Could not find model-catalog deployment in namespace '$MR_NAMESPACE'."
    exit 1
  fi
  kubectl delete pod -l app.kubernetes.io/name=model-catalog -n "$MR_NAMESPACE" --wait=true
  kubectl wait --for=condition=Available deployment/"$CATALOG_DEPLOYMENT" -n "$MR_NAMESPACE" --timeout=5m
  echo "model-catalog pod is ready."
fi
