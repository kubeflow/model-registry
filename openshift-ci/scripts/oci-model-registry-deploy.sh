#!/bin/bash

# Define variables
MODEL_REGISTRY_DEPLOY_MANIFEST="model-registry-operator-deploy.yaml"
DSC_INITIALIZATION_MANIFEST="model-registry-DSCInitialization.yaml"
DATA_SCIENCE_CLUSTER_MANIFEST="model-registry-data-science-cluster.yaml"
TIMEOUT=${DEPLOY_TIMEOUT:-300s}  # Default timeout is 300 seconds, can be overridden by setting DEPLOY_TIMEOUT environment variable

# Function to deploy and wait for deployment
deploy_and_wait() {
    local manifest=$1
    local resource_name=$(basename -s .yaml $manifest)
    
    echo "Deploying $resource_name from $manifest..."
    oc apply -f $manifest

    echo "Waiting for $resource_name to be ready with timeout $TIMEOUT..."
    if ! oc wait --for=condition=Available deployment/$resource_name --timeout=$TIMEOUT; then
        echo "Error: $resource_name deployment failed or timed out."
        exit 1
    fi

    echo "$resource_name deployed successfully!"
}

# Deploy resource and wait for readiness
deploy_resource() {
    local manifest=$1
    deploy_and_wait $manifest
}

# Main function for orchestrating deployments
main() {
    deploy_resource $MODEL_REGISTRY_DEPLOY_MANIFEST
    deploy_resource $DSC_INITIALIZATION_MANIFEST
    deploy_resource $DATA_SCIENCE_CLUSTER_MANIFEST
}

# Execute main function
main