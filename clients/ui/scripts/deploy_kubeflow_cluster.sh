#!/usr/bin/env bash

# Check for required tools
command -v docker >/dev/null 2>&1 || { echo >&2 "Docker is required but it's not installed. Aborting."; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo >&2 "kubectl is required but it's not installed. Aborting."; exit 1; }
command -v kind >/dev/null 2>&1 || { echo >&2 "kind is required but it's not installed. Aborting."; exit 1; }

echo -e "\033[33mWARNING: You must have access to a cluster with kubeflow installed.\033[0m"

# Step 1: Deploy Model Registry UI to cluster
pushd  ../../manifests/kustomize/options/ui/overlays/istio
echo -e "\033[32mDeploying Model Registry UI...\033[0m"
kubectl apply -n kubeflow -k .

# Step 2: Edit the centraldashboard-config ConfigMap
echo -e "\033[32mEditing centraldashboard-config ConfigMap...\033[0m"
kubectl get configmap centraldashboard-config -n kubeflow -o json | jq '.data.links |= (fromjson | .menuLinks += [{"icon": "assignment", "link": "/model-registry/", "text": "Model Registry", "type": "item"}] | tojson)' | kubectl apply -f -

# Wait for deployment to be available
echo -e "\033[32mWaiting Model Registry UI to be available...\033[0m"
kubectl wait --for=condition=available -n kubeflow deployment/model-registry-ui --timeout=1m

# Step 5: Port-forward the service
echo "\033[32mPort-forwarding Kubeflow Central Dashboard...\033[0m"
echo -e "\033[32mDashboard available in http://localhost:8080\033[0m"
kubectl port-forward svc/istio-ingressgateway -n istio-system 8080:80
