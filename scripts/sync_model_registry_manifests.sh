#!/bin/bash

# Define root folder of Model Registry repository
ROOT_FOLDER="${ROOT_FOLDER:-..}"

# Define paths to Model Registry and Kubeflow release repositories
MODEL_REGISTRY_REPO="${ROOT_FOLDER}/model-registry"
KUBEFLOW_MANIFESTS_PATH="${MODEL_REGISTRY_REPO}/manifests/kustomize"

# Ensure the Model Registry repository path exists
if [ ! -d "$MODEL_REGISTRY_REPO" ]; then
    echo "Error: Model Registry repository path does not exist."
    exit 1
fi

# Run the script to generate OpenAPI server
"${MODEL_REGISTRY_REPO}/gen_openapi_server.sh"

# Sync manifests
rsync -av --delete "$MODEL_REGISTRY_REPO/manifests/kustomize/" "$KUBEFLOW_MANIFESTS_PATH/model_registry/"

echo "Model Registry manifests synced successfully."
