#!/bin/bash

set -e
set -o xtrace

# This test assumes there is a Kubernetes environment up and running.
# It could be either a remote one or a local one (e.g., using KinD or minikube).

DIR="$(dirname "$0")"

source ./${DIR}/test_utils.sh

KUBECTL=${KUBECTL:-"kubectl"}

# You can provide a local version of the model registry storage initializer
# In that case, assure that is visible to the local k8s env, e.g., using 
# `kind load docker-image $MRCSI_IMG`
MRCSI_IMG=${MRCSI_IMG:-"kubeflow/model-registry-storage-initializer:main"}

KSERVE_VERSION=${KSERVE_VERSION:-"0.12"}
MODELREGISTRY_VERSION=${MODELREGISTRY_VERSION:-"v0.2.10"}
MODELREGISTRY_CSI=${MODELREGISTRY_CSI:-"v0.2.10"}

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
TESTNAMESPACE=${TESTNAMESPACE:-"test"}
MR_HOSTNAME=localhost:8080
MR_TEST_HOSTNAME=localhost:8082
MODEL_REGISTRY_SERVICE=model-registry-service
MODEL_REGISTRY_REST_PORT=$(kubectl get svc/$MODEL_REGISTRY_SERVICE -n $NAMESPACE --output jsonpath='{.spec.ports[0].targetPort}')
INGRESS_HOST="localhost:8081"
KSERVE_TEST_NAMESPACE=kserve-test

echo "======== Preparing test environment ========"

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
    imagePullPolicy: IfNotPresent
    env:
    - name: MODEL_REGISTRY_BASE_URL
      value: "$MODEL_REGISTRY_SERVICE.$NAMESPACE.svc.cluster.local:$MODEL_REGISTRY_REST_PORT"
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

echo "Serving the istio ingress gateway on $INGRESS_HOST ..."

INGRESS_GATEWAY_SERVICE=$(kubectl get svc -n istio-system --selector="app=istio-ingressgateway" --output jsonpath='{.items[0].metadata.name}')
kubectl port-forward -n istio-system svc/${INGRESS_GATEWAY_SERVICE} 8081:80 &
pf_pid=$!

wait_for_port 8081

echo "Creating $KSERVE_TEST_NAMESPACE namespace ..."

if ! kubectl get namespace $KSERVE_TEST_NAMESPACE &> /dev/null; then
   kubectl create namespace $KSERVE_TEST_NAMESPACE
fi

echo "Creating dummy input data for testing ..."

cat <<EOF > "/tmp/iris-input.json"
{
  "instances": [
    [6.8,  2.8,  4.8,  1.4],
    [6.0,  3.4,  4.5,  1.6]
  ]
}
EOF

echo "======== Finished preparing test environment ========"

echo "======== Scenario 1 - Testing with default model registry service ========"

kubectl port-forward -n $NAMESPACE svc/$MODEL_REGISTRY_SERVICE "8080:$MODEL_REGISTRY_REST_PORT" &
pf_pid=$!

wait_for_port 8080

echo "Initializing data into Model Registry in ${NAMESPACE} namespace..."

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

echo "Starting test ..."

kubectl apply -n $KSERVE_TEST_NAMESPACE -f - <<EOF
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: "sklearn-iris-scenario-one"
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      storageUri: "model-registry://iris/v1"
EOF

# wait for pod predictor to be initialized
repeat_cmd_until "kubectl get pod -n $KSERVE_TEST_NAMESPACE --selector='component=predictor' | wc -l" "-gt 0" 60
predictor_one=$(kubectl get pod -n $KSERVE_TEST_NAMESPACE --selector="component=predictor" --output jsonpath='{.items[0].metadata.name}')

kubectl wait --for=condition=Ready pod/$predictor_one -n $KSERVE_TEST_NAMESPACE --timeout=5m

kubectl wait --for=jsonpath='{.status.url}' inferenceservice/sklearn-iris-scenario-one -n $KSERVE_TEST_NAMESPACE --timeout=5m
sleep 5

SERVICE_HOSTNAME=$(kubectl get inferenceservice sklearn-iris-scenario-one -n $KSERVE_TEST_NAMESPACE -o jsonpath='{.status.url}' | cut -d "/" -f 3)
res_one=$(curl -s -H "Host: ${SERVICE_HOSTNAME}" -H "Content-Type: application/json" "http://${INGRESS_HOST}/v1/models/sklearn-iris-scenario-one:predict" -d @/tmp/iris-input.json)
echo "Received: $res_one"

if [ ! "$res_one" = "{\"predictions\":[1,1]}" ]; then
    echo "Prediction does not match expectation!"
    echo "Printing some logs for debugging.."
    kubectl logs pod/$predictor_one -n $KSERVE_TEST_NAMESPACE -c storage-initializer
    kubectl logs pod/$predictor_one -n $KSERVE_TEST_NAMESPACE -c kserve-container
    exit 1
else
    echo "Scenario 1 - Test succeeded!"
fi

echo "Cleaning up inferenceservice sklearn-iris-scenario-one ..."

kubectl delete inferenceservice sklearn-iris-scenario-one -n $KSERVE_TEST_NAMESPACE

echo "======== Finished Scenario 1 ========"

echo "======== Scenario 2 - Testing with default model registry service without model version ========"

echo "Starting test ..."

kubectl apply -n $KSERVE_TEST_NAMESPACE -f - <<EOF
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: "sklearn-iris-scenario-two"
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      storageUri: "model-registry://iris"
EOF

# wait for pod predictor to be initialized
repeat_cmd_until "kubectl get pod -n $KSERVE_TEST_NAMESPACE --selector='component=predictor' --output jsonpath='{.items[*].metadata.name}' | grep sklearn-iris-scenario-two | wc -l" "-gt 0" 60
predictor_two=$(kubectl get pod -n $KSERVE_TEST_NAMESPACE --selector="component=predictor" --output jsonpath='{.items[1].metadata.name}')

kubectl wait --for=condition=Ready pod/$predictor_two -n $KSERVE_TEST_NAMESPACE --timeout=5m

kubectl wait --for=jsonpath='{.status.url}' inferenceservice/sklearn-iris-scenario-two -n $KSERVE_TEST_NAMESPACE --timeout=5m
sleep 5

SERVICE_HOSTNAME_TEST=$(kubectl get inferenceservice sklearn-iris-scenario-two -n $KSERVE_TEST_NAMESPACE -o jsonpath='{.status.url}' | cut -d "/" -f 3)
res_two=$(curl -s -H "Host: ${SERVICE_HOSTNAME_TEST}" -H "Content-Type: application/json" "http://${INGRESS_HOST}/v1/models/sklearn-iris-scenario-two:predict" -d @/tmp/iris-input.json)
echo "Received: $res_two"

if [ ! "$res_two" = "{\"predictions\":[1,1]}" ]; then
    echo "Prediction does not match expectation!"
    echo "Printing some logs for debugging.."
    kubectl logs pod/$predictor_two -n $KSERVE_TEST_NAMESPACE -c storage-initializer
    kubectl logs pod/$predictor_two -n $KSERVE_TEST_NAMESPACE -c kserve-container
    exit 1
else
    echo "Scenario 2 - Test succeeded!"
fi

echo "Cleaning up inferenceservice sklearn-iris-scenario-two ..."

kubectl delete inferenceservice sklearn-iris-scenario-two -n $KSERVE_TEST_NAMESPACE

echo "======== Finished Scenario 2 ========"

echo "======== Scenario 3 - Testing with custom model registry service ========"

kubectl port-forward -n $TESTNAMESPACE svc/$MODEL_REGISTRY_SERVICE "8082:$MODEL_REGISTRY_REST_PORT" &
pf_pid=$!

wait_for_port 8082

echo "Initializing data into Model Registry in ${TESTNAMESPACE} namespace..."

curl --silent -X 'POST' \
  "$MR_TEST_HOSTNAME/api/model_registry/v1alpha3/registered_models" \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "description": "Iris scikit-learn model",
  "name": "iris-test"
}'

curl --silent -X 'POST' \
  "$MR_TEST_HOSTNAME/api/model_registry/v1alpha3/model_versions" \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "description": "Iris model version v1",
  "name": "v1-test",
  "registeredModelID": "1"
}'

curl --silent -X 'POST' \
  "$MR_TEST_HOSTNAME/api/model_registry/v1alpha3/model_versions/2/artifacts" \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "description": "Model artifact for Iris v1",
  "uri": "gs://kfserving-examples/models/sklearn/1.0/model",
  "state": "UNKNOWN",
  "name": "sklearn-iris-test-v1",
  "modelFormatName": "sklearn",
  "modelFormatVersion": "1",
  "artifactType": "model-artifact"
}'

echo "Starting test ..."

kubectl apply -n $KSERVE_TEST_NAMESPACE -f - <<EOF
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: "sklearn-iris-scenario-three"
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      storageUri: "model-registry://$MODEL_REGISTRY_SERVICE.${TESTNAMESPACE}.svc.cluster.local:$MODEL_REGISTRY_REST_PORT/iris-test/v1-test"
EOF

# wait for pod predictor to be initialized
repeat_cmd_until "kubectl get pod -n $KSERVE_TEST_NAMESPACE --selector='component=predictor' --output jsonpath='{.items[*].metadata.name}' | grep sklearn-iris-scenario-three-predictor | wc -l" "-gt 0" 60
predictor_three=$(kubectl get pod -n $KSERVE_TEST_NAMESPACE --selector="component=predictor" --output jsonpath='{.items[1].metadata.name}')

kubectl wait --for=condition=Ready pod/$predictor_three -n $KSERVE_TEST_NAMESPACE --timeout=5m

kubectl wait --for=jsonpath='{.status.url}' inferenceservice/sklearn-iris-scenario-three -n $KSERVE_TEST_NAMESPACE --timeout=5m
sleep 5

SERVICE_HOSTNAME_TEST=$(kubectl get inferenceservice sklearn-iris-scenario-three -n $KSERVE_TEST_NAMESPACE -o jsonpath='{.status.url}' | cut -d "/" -f 3)
res_three=$(curl -s -H "Host: ${SERVICE_HOSTNAME_TEST}" -H "Content-Type: application/json" "http://${INGRESS_HOST}/v1/models/sklearn-iris-scenario-three:predict" -d @/tmp/iris-input.json)
echo "Received: $res_three"

if [ ! "$res_three" = "{\"predictions\":[1,1]}" ]; then
    echo "Prediction does not match expectation!"
    echo "Printing some logs for debugging.."
    kubectl logs pod/$predictor_three -n $KSERVE_TEST_NAMESPACE -c storage-initializer
    kubectl logs pod/$predictor_three -n $KSERVE_TEST_NAMESPACE -c kserve-container
    exit 1
else
    echo "Scenario 3 - Test succeeded!"
fi

echo "Cleaning up inferenceservice sklearn-iris-scenario-three ..."

kubectl delete inferenceservice sklearn-iris-scenario-three -n $KSERVE_TEST_NAMESPACE

echo "======== Finished Scenario 3 ========"

echo "======== Scenario 4 - Testing with custom model registry service without model version ========"

echo "Starting test ..."

kubectl apply -n $KSERVE_TEST_NAMESPACE -f - <<EOF
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: "sklearn-iris-scenario-four"
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      storageUri: "model-registry://$MODEL_REGISTRY_SERVICE.${TESTNAMESPACE}.svc.cluster.local:$MODEL_REGISTRY_REST_PORT/iris-test"
EOF

# wait for pod predictor to be initialized
repeat_cmd_until "kubectl get pod -n $KSERVE_TEST_NAMESPACE --selector='component=predictor' --output jsonpath='{.items[*].metadata.name}' | grep sklearn-iris-scenario-four-predictor | wc -l" "-gt 0" 60
predictor_four=$(kubectl get pod -n $KSERVE_TEST_NAMESPACE --selector="component=predictor" --output jsonpath='{.items[1].metadata.name}')

kubectl wait --for=condition=Ready pod/$predictor_four -n $KSERVE_TEST_NAMESPACE --timeout=5m

kubectl wait --for=jsonpath='{.status.url}' inferenceservice/sklearn-iris-scenario-four -n $KSERVE_TEST_NAMESPACE --timeout=5m
sleep 5

SERVICE_HOSTNAME_TEST=$(kubectl get inferenceservice sklearn-iris-scenario-four -n $KSERVE_TEST_NAMESPACE -o jsonpath='{.status.url}' | cut -d "/" -f 3)
res_four=$(curl -s -H "Host: ${SERVICE_HOSTNAME_TEST}" -H "Content-Type: application/json" "http://${INGRESS_HOST}/v1/models/sklearn-iris-scenario-four:predict" -d @/tmp/iris-input.json)
echo "Received: $res_four"

if [ ! "$res_four" = "{\"predictions\":[1,1]}" ]; then
    echo "Prediction does not match expectation!"
    echo "Printing some logs for debugging.."
    kubectl logs pod/$predictor_four -n $KSERVE_TEST_NAMESPACE -c storage-initializer
    kubectl logs pod/$predictor_four -n $KSERVE_TEST_NAMESPACE -c kserve-container
    exit 1
else
    echo "Scenario 4 - Test succeeded!"
fi

echo "Cleaning up inferenceservice sklearn-iris-scenario-four ..."

kubectl delete inferenceservice sklearn-iris-scenario-four -n $KSERVE_TEST_NAMESPACE

echo "======== Finished Scenario 4 ========"

echo "All tests passed!"
