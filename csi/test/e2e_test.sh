#!/bin/bash

set -e
set -o xtrace

# This test assumes there is a Kubernetes environment up and running.
# It could be either a remote one or a local one (e.g., using KinD or minikube).

# Function to check if the port is ready
wait_for_port() {
  local port=$1
  while ! nc -z localhost $port; do
    sleep 0.1
  done
}

DIR="$(dirname "$0")"

KUBECTL=${KUBECTL:-"kubectl"}

# You can provide a local version of the model registry storage initializer
# In that case, assure that is visible to the local k8s env, e.g., using 
# `kind load docker-image $MRCSI_IMG`
MRCSI_IMG=${MRCSI_IMG:-"kubeflow/model-registry-storage-initializer:main"}

KSERVE_VERSION=${KSERVE_VERSION:-"0.12"}
MODELREGISTRY_VERSION=${MODELREGISTRY_VERSION:-"v0.2.2-alpha"}
MODELREGISTRY_CSI=${MODELREGISTRY_CSI:-"v0.2.2-alpha"}

# You can provide a local model registry container image
MR_IMG=${MR_IMG:-"kubeflow/model-registry:$MODELREGISTRY_VERSION"}
# You can provide a local model registry storage initializer container image
MR_CSI_IMG=${MR_CSI_IMG:-"kubeflow/model-registry-storage-initializer:$MODELREGISTRY_CSI"}

# Check if KUBECTL is a valid command
if [ ! command -v "$KUBECTL" > /dev/null 2>&1 ]; then
    echo "KUBECTL command not found at: $KUBECTL"
    exit 1
fi

if [ ! "$KUBECTL" cluster-info > /dev/null 2>&1 ]; then
    echo "Cluster not available!"
    exit 1
fi

# Setup the environment
./${DIR}/setup_test_env.sh

# Apply the port forward to access the model registry
NAMESPACE=${NAMESPACE:-"kubeflow"}
MR_HOSTNAME=localhost:8080
MODEL_REGISTRY_SERVICE=model-registry-service

MODEL_REGISTRY_REST_PORT=$(kubectl get svc/$MODEL_REGISTRY_SERVICE -n $NAMESPACE --output jsonpath='{.spec.ports[0].targetPort}')

kubectl port-forward -n $NAMESPACE svc/$MODEL_REGISTRY_SERVICE "8080:$MODEL_REGISTRY_REST_PORT" &
pf_pid=$!

wait_for_port 8080

echo "Initializing data into Model Registry ..."

curl --silent -X 'POST' \
  "$MR_HOSTNAME/api/model_registry/v1alpha3/registered_models" \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "description": "Iris scikit-learn model",
  "name": "iris"
}'

curl --silent -X 'POST' \
  "$MR_HOSTNAME/api/model_registry/v1alpha3/model_versions" \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "description": "Iris model version v1",
  "name": "v1",
  "registeredModelID": "1"
}'

curl --silent -X 'POST' \
  "$MR_HOSTNAME/api/model_registry/v1alpha3/model_versions/2/artifacts" \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "description": "Model artifact for Iris v1",
  "uri": "gs://kfserving-examples/models/sklearn/1.0/model",
  "state": "UNKNOWN",
  "name": "sklearn-iris-v1",
  "modelFormatName": "sklearn",
  "modelFormatVersion": "1",
  "artifactType": "model-artifact"
}'

echo "======== Model Registry populated ========"

echo "Applying Model Registry custom storage initializer ..."

kubectl apply -f - <<EOF
apiVersion: "serving.kserve.io/v1alpha1"
kind: ClusterStorageContainer
metadata:
  name: mr-initializer
spec:
  container:
    name: storage-initializer
    image: $MR_CSI_IMG
    env:
    - name: MODEL_REGISTRY_BASE_URL
      value: "$MODEL_REGISTRY_SERVICE.kubeflow.svc.cluster.local:$MODEL_REGISTRY_REST_PORT"
    - name: MODEL_REGISTRY_SCHEME
      value: "http"
    resources:
      requests:
        memory: 100Mi
        cpu: 100m
      limits:
        memory: 1Gi
        cpu: "1"
  supportedUriFormats:
    - prefix: model-registry://
EOF

echo "======== Custom storage initializer applied ========"

echo "Starting test ..."

KSERVE_TEST_NAMESPACE=kserve-test
if ! kubectl get namespace $KSERVE_TEST_NAMESPACE &> /dev/null; then
   kubectl create namespace $KSERVE_TEST_NAMESPACE
fi

kubectl apply -n $KSERVE_TEST_NAMESPACE -f - <<EOF
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: "sklearn-iris"
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      storageUri: "model-registry://iris/v1"
EOF

# wait for pod predictor to be initialized
sleep 2
predictor=$(kubectl get pod -n $KSERVE_TEST_NAMESPACE --selector="component=predictor" --output jsonpath='{.items[0].metadata.name}')
kubectl wait --for=condition=Ready pod/$predictor -n $KSERVE_TEST_NAMESPACE --timeout=5m

INGRESS_GATEWAY_SERVICE=$(kubectl get svc -n istio-system --selector="app=istio-ingressgateway" --output jsonpath='{.items[0].metadata.name}')
kubectl port-forward -n istio-system svc/${INGRESS_GATEWAY_SERVICE} 8081:80 &
pf_pid=$!

wait_for_port 8081

INGRESS_HOST="localhost:8081"

cat <<EOF > "/tmp/iris-input.json"
{
  "instances": [
    [6.8,  2.8,  4.8,  1.4],
    [6.0,  3.4,  4.5,  1.6]
  ]
}
EOF

# kubectl wait --for=condition=Ready inferenceservice/sklearn-iris -n $KSERVE_TEST_NAMESPACE --timeout=5m
kubectl wait --for=jsonpath='{.status.url}' inferenceservice/sklearn-iris -n $KSERVE_TEST_NAMESPACE --timeout=5m
sleep 5

SERVICE_HOSTNAME=$(kubectl get inferenceservice sklearn-iris -n $KSERVE_TEST_NAMESPACE -o jsonpath='{.status.url}' | cut -d "/" -f 3)
res=$(curl -s -H "Host: ${SERVICE_HOSTNAME}" -H "Content-Type: application/json" "http://${INGRESS_HOST}/v1/models/sklearn-iris:predict" -d @/tmp/iris-input.json)
echo "Received: $res"

if [ ! "$res" = "{\"predictions\":[1,1]}" ]; then
    echo "Prediction does not match expectation!"
    echo "Printing some logs for debugging.."
    kubectl logs pod/$predictor -n $KSERVE_TEST_NAMESPACE -c storage-initializer
    kubectl logs pod/$predictor -n $KSERVE_TEST_NAMESPACE -c kserve-container
    exit 1
else
    echo "Test succeeded!"
fi