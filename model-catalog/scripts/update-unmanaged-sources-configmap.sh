#!/bin/bash

# Check if required tools are installed
for cmd in yq jq; do
    if ! command -v $cmd >/dev/null 2>&1; then
        echo "Error: $cmd is not installed. Please install it first."
        exit 1
    fi
done

# Check if logged into a cluster
if ! oc whoami >/dev/null 2>&1; then
    echo "Error: Not logged into a cluster. Please run 'oc login' first."
    exit 1
fi

# Check if input file is provided
if [ $# -lt 1 ] || [ $# -gt 2 ]; then
    echo "Usage: $0 <input-yaml-file> [namespace]"
    echo "  namespace: Optional. Defaults to redhat-ods-applications. Can be changed to opendatahub if needed."
    exit 1
fi

INPUT_YAML=$1
# Default to redhat-ods-applications, can be overridden by second argument
NAMESPACE=${2:-redhat-ods-applications}

# Check if namespace exists
if ! oc get namespace "$NAMESPACE" >/dev/null 2>&1; then
    echo "Error: Namespace '$NAMESPACE' does not exist in the cluster"
    exit 1
fi

# Check if input file exists
if [ ! -f "$INPUT_YAML" ]; then
    echo "Error: Input file $INPUT_YAML does not exist"
    exit 1
fi

# Convert input YAML to JSON and wrap it in a sources array
_model_catalog_sources_json_string=$(yq -o=json "$INPUT_YAML" | jq -c '{sources: [.]}')

# Grab the existing configmap and update the modelCatalogSources field with the new content
mkdir tmp
oc get configmap model-catalog-unmanaged-sources -n "$NAMESPACE" -o yaml > tmp/model-catalog-unmanaged-sources.yaml

# Write the new JSON string to a temporary file
echo "$_model_catalog_sources_json_string" > tmp/new_sources_content.json

# Use yq to load the string content from the temporary file
# and assign it to the modelCatalogSources field.
yq -i '.data.modelCatalogSources = load_str("tmp/new_sources_content.json")' tmp/model-catalog-unmanaged-sources.yaml

# Update the configmap with the new content
oc apply -f tmp/model-catalog-unmanaged-sources.yaml -n "$NAMESPACE"

# Clean up
rm tmp/model-catalog-unmanaged-sources.yaml
rm tmp/new_sources_content.json # Remove the temporary content file
rmdir tmp

echo "Success"