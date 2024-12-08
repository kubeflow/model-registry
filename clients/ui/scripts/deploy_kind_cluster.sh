#!/usr/bin/env bash

# Check for required tools
command -v docker >/dev/null 2>&1 || { echo >&2 "Docker is required but it's not installed. Aborting."; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo >&2 "kubectl is required but it's not installed. Aborting."; exit 1; }
command -v kind >/dev/null 2>&1 || { echo >&2 "kind is required but it's not installed. Aborting."; exit 1; }

if kubectl get deployment model-registry-deployment -n kubeflow >/dev/null 2>&1; then
  echo "Model Registry deployment already exists. Skipping to step 4."
else
    # Step 1: Create a kind cluster
    echo "Creating kind cluster..."
    kind create cluster

    # Verify cluster creation
    echo "Verifying cluster..."
    kubectl cluster-info

    # Step 2: Create kubeflow namespace
    echo "Creating kubeflow namespace..."
    kubectl create namespace kubeflow

    # Step 3: Deploy Model Registry to cluster
    echo "Deploying Model Registry to cluster..."
    kubectl apply -k "https://github.com/alexcreasy/model-registry/manifests/kustomize/overlays/db?ref=kind"

    # Wait for deployment to be available
    echo "Waiting for Model Registry deployment to be available..."
    kubectl wait --for=condition=available -n kubeflow deployment/model-registry-deployment --timeout=1m

    # Verify deployment
    echo "Verifying deployment..."
    kubectl get pods -n kubeflow
fi

pushd  ./manifests/base
kustomize edit set namespace kubeflow
kustomize edit set image model-registry-ui-image=${IMG_FRONTEND}
kustomize edit set image model-registry-bff-image=${IMG_BFF}

# Step 4: Deploy model registry UI
echo "Deploying Model Registry UI..."
kubectl apply -n kubeflow -k .

# Wait for deployment to be available
echo "Waiting Model Registry UI to be available..."
kubectl wait --for=condition=available -n kubeflow deployment/model-registry-ui --timeout=1m

pushd  ../user-rbac
# Step 5: Apply admin user service account in the cluster
echo "Applying admin user service account and rolebinding..."
kubectl apply -k .

# Step 6: Generate token for admin user and display it
echo "In your browser, you will need to inject your requests with a kubeflow-userid header for authorization purposes."
echo "For example, you can use the Header Editor - https://chromewebstore.google.com/detail/eningockdidmgiojffjmkdblpjocbhgh extension in Chrome to set the kubeflow-userid header to user@example.com."

# Step 5: Port-forward the service
echo "Port-forwarding Model Registry UI..."
echo -e "\033[32mDashboard available in http://localhost:8080\033[0m"
kubectl port-forward svc/model-registry-ui-service -n kubeflow 8080:8080
