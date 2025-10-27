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
MRCSI_IMG=${MRCSI_IMG:-"ghcr.io/kubeflow/model-registry/storage-initializer:main"}

KSERVE_VERSION=${KSERVE_VERSION:-"0.15"}
MODELREGISTRY_VERSION=${MODELREGISTRY_VERSION:-"v0.3.2"}
MODELREGISTRY_CSI=${MODELREGISTRY_CSI:-"v0.3.2"}

# You can provide a local model registry container image
MR_IMG=${MR_IMG:-"ghcr.io/kubeflow/model-registry/server:$MODELREGISTRY_VERSION"}
# You can provide a local model registry storage initializer container image
MR_CSI_IMG=${MR_CSI_IMG:-"ghcr.io/kubeflow/model-registry/storage-initializer:$MODELREGISTRY_CSI"}

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
MODEL_REGISTRY_SERVICE=model-registry-service
MODEL_REGISTRY_REST_PORT=$(kubectl get svc/$MODEL_REGISTRY_SERVICE -n $NAMESPACE --output jsonpath='{.spec.ports[0].targetPort}')
INGRESS_HOST="localhost:8081"

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
  supportedUriFormats:
    - prefix: model-registry://
EOF

echo "Serving the istio ingress gateway on $INGRESS_HOST ..."

INGRESS_GATEWAY_SERVICE=$(kubectl get svc -n istio-system --selector="app=istio-ingressgateway" --output jsonpath='{.items[0].metadata.name}')
kubectl port-forward -n istio-system svc/${INGRESS_GATEWAY_SERVICE} 8081:80 &
ingress_pf_pid=$!

wait_for_port 8081

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

# Array to track background process PIDs
declare -a scenario_pids=()
declare -a scenario_results=()

# Function to run Scenario 1
run_scenario_1() {
    local scenario_num=1
    local kserve_namespace="kserve-test-${scenario_num}"
    local mr_port=8080
    local mr_hostname="localhost:${mr_port}"

    echo "======== Scenario ${scenario_num} - Testing with default model registry service ========"

    # Create namespace
    if ! kubectl get namespace $kserve_namespace &> /dev/null; then
       kubectl create namespace $kserve_namespace
    fi

    # Port forward for model registry
    kubectl port-forward -n $NAMESPACE svc/$MODEL_REGISTRY_SERVICE "${mr_port}:$MODEL_REGISTRY_REST_PORT" &
    local pf_pid=$!

    wait_for_port $mr_port

    echo "Initializing data into Model Registry for scenario ${scenario_num}..."

    # Create registered model and capture the ID
    local rm_response=$(curl --silent -X 'POST' \
      "$mr_hostname/api/model_registry/v1alpha3/registered_models" \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d '{
      "description": "Iris scikit-learn model",
      "name": "iris-scenario-1"
    }')
    local rm_id=$(echo "$rm_response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    echo "Created registered model with ID: $rm_id"

    # Create model version and capture the ID
    local mv_response=$(curl --silent -X 'POST' \
      "$mr_hostname/api/model_registry/v1alpha3/model_versions" \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d "{
      \"description\": \"Iris model version v1\",
      \"name\": \"v1\",
      \"registeredModelID\": \"$rm_id\"
    }")
    local mv_id=$(echo "$mv_response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    echo "Created model version with ID: $mv_id"

    # Create model artifact
    curl --silent -X 'POST' \
      "$mr_hostname/api/model_registry/v1alpha3/model_versions/$mv_id/artifacts" \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d '{
      "description": "Model artifact for Iris v1",
      "uri": "gs://kfserving-examples/models/sklearn/1.0/model",
      "state": "UNKNOWN",
      "name": "sklearn-iris-v1-scenario-1",
      "modelFormatName": "sklearn",
      "modelFormatVersion": "1",
      "artifactType": "model-artifact"
    }'

    echo "Starting test for scenario ${scenario_num}..."

    kubectl apply -n $kserve_namespace -f - <<EOF
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: "sklearn-iris-scenario-one"
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      storageUri: "model-registry://iris-scenario-1/v1"
EOF

    # wait for pod predictor to be initialized
    repeat_cmd_until "kubectl get pod -n $kserve_namespace --selector='component=predictor' | wc -l" "-gt 0" 60
    local predictor=$(kubectl get pod -n $kserve_namespace --selector="component=predictor" --output jsonpath='{.items[0].metadata.name}')

    kubectl wait --for=condition=Ready pod/$predictor -n $kserve_namespace --timeout=5m

    kubectl wait --for=jsonpath='{.status.url}' inferenceservice/sklearn-iris-scenario-one -n $kserve_namespace --timeout=5m
    sleep 5

    local service_hostname=$(kubectl get inferenceservice sklearn-iris-scenario-one -n $kserve_namespace -o jsonpath='{.status.url}' | cut -d "/" -f 3)
    local result=$(curl -s -H "Host: ${service_hostname}" -H "Content-Type: application/json" "http://${INGRESS_HOST}/v1/models/sklearn-iris-scenario-one:predict" -d @/tmp/iris-input.json)
    echo "Scenario ${scenario_num} received: $result"

    if [ ! "$result" = "{\"predictions\":[1,1]}" ]; then
        echo "Scenario ${scenario_num} - Prediction does not match expectation!"
        echo "Printing some logs for debugging.."
        kubectl logs pod/$predictor -n $kserve_namespace -c storage-initializer || true
        kubectl logs pod/$predictor -n $kserve_namespace -c kserve-container || true
        kill $pf_pid 2>/dev/null || true
        return 1
    else
        echo "Scenario ${scenario_num} - Test succeeded!"
        kubectl logs pod/$predictor -n $kserve_namespace -c storage-initializer || true
        kubectl logs pod/$predictor -n $kserve_namespace -c kserve-container || true
    fi

    echo "Cleaning up scenario ${scenario_num}..."
    kubectl delete inferenceservice sklearn-iris-scenario-one -n $kserve_namespace

    sleep 5
    repeat_cmd_until "kubectl get pod -n $kserve_namespace --selector='component=predictor' --output jsonpath='{.items[*].metadata.name}' | wc -w" "= 0" 60

    kill $pf_pid 2>/dev/null || true

    echo "======== Finished Scenario ${scenario_num} ========"
    return 0
}

# Function to run Scenario 2
run_scenario_2() {
    local scenario_num=2
    local kserve_namespace="kserve-test-${scenario_num}"
    local mr_port=8090
    local mr_hostname="localhost:${mr_port}"

    echo "======== Scenario ${scenario_num} - Testing with default model registry service without model version ========"

    # Create namespace
    if ! kubectl get namespace $kserve_namespace &> /dev/null; then
       kubectl create namespace $kserve_namespace
    fi

    # Port forward for model registry
    kubectl port-forward -n $NAMESPACE svc/$MODEL_REGISTRY_SERVICE "${mr_port}:$MODEL_REGISTRY_REST_PORT" &
    local pf_pid=$!

    wait_for_port $mr_port

    echo "Initializing data into Model Registry for scenario ${scenario_num}..."

    # Create registered model and capture the ID
    local rm_response=$(curl --silent -X 'POST' \
      "$mr_hostname/api/model_registry/v1alpha3/registered_models" \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d '{
      "description": "Iris scikit-learn model",
      "name": "iris-scenario-2"
    }')
    local rm_id=$(echo "$rm_response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    echo "Created registered model with ID: $rm_id"

    # Create model version and capture the ID
    local mv_response=$(curl --silent -X 'POST' \
      "$mr_hostname/api/model_registry/v1alpha3/model_versions" \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d "{
      \"description\": \"Iris model version v1\",
      \"name\": \"v1\",
      \"registeredModelID\": \"$rm_id\"
    }")
    local mv_id=$(echo "$mv_response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    echo "Created model version with ID: $mv_id"

    # Create model artifact
    curl --silent -X 'POST' \
      "$mr_hostname/api/model_registry/v1alpha3/model_versions/$mv_id/artifacts" \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d '{
      "description": "Model artifact for Iris v1",
      "uri": "gs://kfserving-examples/models/sklearn/1.0/model",
      "state": "UNKNOWN",
      "name": "sklearn-iris-v1-scenario-2",
      "modelFormatName": "sklearn",
      "modelFormatVersion": "1",
      "artifactType": "model-artifact"
    }'

    echo "Starting test for scenario ${scenario_num}..."

    kubectl apply -n $kserve_namespace -f - <<EOF
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: "sklearn-iris-scenario-two"
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      storageUri: "model-registry://iris-scenario-2"
EOF

    # wait for pod predictor to be initialized
    repeat_cmd_until "kubectl get pod -n $kserve_namespace --selector='component=predictor' --output jsonpath='{.items[*].metadata.name}' | grep sklearn-iris-scenario-two | wc -l" "-gt 0" 60
    local predictor=$(kubectl get pod -n $kserve_namespace --selector="component=predictor" --output jsonpath='{.items[0].metadata.name}')

    kubectl wait --for=condition=Ready pod/$predictor -n $kserve_namespace --timeout=5m

    kubectl wait --for=jsonpath='{.status.url}' inferenceservice/sklearn-iris-scenario-two -n $kserve_namespace --timeout=5m
    sleep 5

    local service_hostname=$(kubectl get inferenceservice sklearn-iris-scenario-two -n $kserve_namespace -o jsonpath='{.status.url}' | cut -d "/" -f 3)
    local result=$(curl -s -H "Host: ${service_hostname}" -H "Content-Type: application/json" "http://${INGRESS_HOST}/v1/models/sklearn-iris-scenario-two:predict" -d @/tmp/iris-input.json)
    echo "Scenario ${scenario_num} received: $result"

    if [ ! "$result" = "{\"predictions\":[1,1]}" ]; then
        echo "Scenario ${scenario_num} - Prediction does not match expectation!"
        echo "Printing some logs for debugging.."
        kubectl logs pod/$predictor -n $kserve_namespace -c storage-initializer || true
        kubectl logs pod/$predictor -n $kserve_namespace -c kserve-container || true
        kill $pf_pid 2>/dev/null || true
        return 1
    else
        echo "Scenario ${scenario_num} - Test succeeded!"
        kubectl logs pod/$predictor -n $kserve_namespace -c storage-initializer || true
        kubectl logs pod/$predictor -n $kserve_namespace -c kserve-container || true
    fi

    echo "Cleaning up scenario ${scenario_num}..."
    kubectl delete inferenceservice sklearn-iris-scenario-two -n $kserve_namespace

    sleep 5
    repeat_cmd_until "kubectl get pod -n $kserve_namespace --selector='component=predictor' --output jsonpath='{.items[*].metadata.name}' | wc -w" "= 0" 60

    kill $pf_pid 2>/dev/null || true

    echo "======== Finished Scenario ${scenario_num} ========"
    return 0
}

# Function to run Scenario 3
run_scenario_3() {
    local scenario_num=3
    local kserve_namespace="kserve-test-${scenario_num}"
    local mr_port=8082
    local mr_hostname="localhost:${mr_port}"

    echo "======== Scenario ${scenario_num} - Testing with custom model registry service ========"

    # Create namespace
    if ! kubectl get namespace $kserve_namespace &> /dev/null; then
       kubectl create namespace $kserve_namespace
    fi

    # Port forward for model registry in test namespace
    kubectl port-forward -n $TESTNAMESPACE svc/$MODEL_REGISTRY_SERVICE "${mr_port}:$MODEL_REGISTRY_REST_PORT" &
    local pf_pid=$!

    wait_for_port $mr_port

    echo "Initializing data into Model Registry in ${TESTNAMESPACE} namespace for scenario ${scenario_num}..."

    # Create registered model and capture the ID
    local rm_response=$(curl --silent -X 'POST' \
      "$mr_hostname/api/model_registry/v1alpha3/registered_models" \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d '{
      "description": "Iris scikit-learn model",
      "name": "iris-test-scenario-3"
    }')
    local rm_id=$(echo "$rm_response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    echo "Created registered model with ID: $rm_id"

    # Create model version and capture the ID
    local mv_response=$(curl --silent -X 'POST' \
      "$mr_hostname/api/model_registry/v1alpha3/model_versions" \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d "{
      \"description\": \"Iris model version v1\",
      \"name\": \"v1-test\",
      \"registeredModelID\": \"$rm_id\"
    }")
    local mv_id=$(echo "$mv_response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    echo "Created model version with ID: $mv_id"

    # Create model artifact
    curl --silent -X 'POST' \
      "$mr_hostname/api/model_registry/v1alpha3/model_versions/$mv_id/artifacts" \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d '{
      "description": "Model artifact for Iris v1",
      "uri": "gs://kfserving-examples/models/sklearn/1.0/model",
      "state": "UNKNOWN",
      "name": "sklearn-iris-test-v1-scenario-3",
      "modelFormatName": "sklearn",
      "modelFormatVersion": "1",
      "artifactType": "model-artifact"
    }'

    echo "Starting test for scenario ${scenario_num}..."

    kubectl apply -n $kserve_namespace -f - <<EOF
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: "sklearn-iris-scenario-three"
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      storageUri: "model-registry://$MODEL_REGISTRY_SERVICE.${TESTNAMESPACE}.svc.cluster.local:$MODEL_REGISTRY_REST_PORT/iris-test-scenario-3/v1-test"
EOF

    # wait for pod predictor to be initialized
    repeat_cmd_until "kubectl get pod -n $kserve_namespace --selector='component=predictor' --output jsonpath='{.items[*].metadata.name}' | grep sklearn-iris-scenario-three-predictor | wc -l" "-gt 0" 60
    local predictor=$(kubectl get pod -n $kserve_namespace --selector="component=predictor" --output jsonpath='{.items[0].metadata.name}')

    kubectl wait --for=condition=Ready pod/$predictor -n $kserve_namespace --timeout=5m

    kubectl wait --for=jsonpath='{.status.url}' inferenceservice/sklearn-iris-scenario-three -n $kserve_namespace --timeout=5m
    sleep 5

    local service_hostname=$(kubectl get inferenceservice sklearn-iris-scenario-three -n $kserve_namespace -o jsonpath='{.status.url}' | cut -d "/" -f 3)
    local result=$(curl -s -H "Host: ${service_hostname}" -H "Content-Type: application/json" "http://${INGRESS_HOST}/v1/models/sklearn-iris-scenario-three:predict" -d @/tmp/iris-input.json)
    echo "Scenario ${scenario_num} received: $result"

    if [ ! "$result" = "{\"predictions\":[1,1]}" ]; then
        echo "Scenario ${scenario_num} - Prediction does not match expectation!"
        echo "Printing some logs for debugging.."
        kubectl logs pod/$predictor -n $kserve_namespace -c storage-initializer || true
        kubectl logs pod/$predictor -n $kserve_namespace -c kserve-container || true
        kill $pf_pid 2>/dev/null || true
        return 1
    else
        echo "Scenario ${scenario_num} - Test succeeded!"
        kubectl logs pod/$predictor -n $kserve_namespace -c storage-initializer || true
        kubectl logs pod/$predictor -n $kserve_namespace -c kserve-container || true
    fi

    echo "Cleaning up scenario ${scenario_num}..."
    kubectl delete inferenceservice sklearn-iris-scenario-three -n $kserve_namespace

    sleep 5
    repeat_cmd_until "kubectl get pod -n $kserve_namespace --selector='component=predictor' --output jsonpath='{.items[*].metadata.name}' | wc -w" "= 0" 60

    kill $pf_pid 2>/dev/null || true

    echo "======== Finished Scenario ${scenario_num} ========"
    return 0
}

# Function to run Scenario 4
run_scenario_4() {
    local scenario_num=4
    local kserve_namespace="kserve-test-${scenario_num}"
    local mr_port=8092
    local mr_hostname="localhost:${mr_port}"

    echo "======== Scenario ${scenario_num} - Testing with custom model registry service without model version ========"

    # Create namespace
    if ! kubectl get namespace $kserve_namespace &> /dev/null; then
       kubectl create namespace $kserve_namespace
    fi

    # Port forward for model registry in test namespace
    kubectl port-forward -n $TESTNAMESPACE svc/$MODEL_REGISTRY_SERVICE "${mr_port}:$MODEL_REGISTRY_REST_PORT" &
    local pf_pid=$!

    wait_for_port $mr_port

    echo "Initializing data into Model Registry in ${TESTNAMESPACE} namespace for scenario ${scenario_num}..."

    # Create registered model and capture the ID
    local rm_response=$(curl --silent -X 'POST' \
      "$mr_hostname/api/model_registry/v1alpha3/registered_models" \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d '{
      "description": "Iris scikit-learn model",
      "name": "iris-test-scenario-4"
    }')
    local rm_id=$(echo "$rm_response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    echo "Created registered model with ID: $rm_id"

    # Create model version and capture the ID
    local mv_response=$(curl --silent -X 'POST' \
      "$mr_hostname/api/model_registry/v1alpha3/model_versions" \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d "{
      \"description\": \"Iris model version v1\",
      \"name\": \"v1-test\",
      \"registeredModelID\": \"$rm_id\"
    }")
    local mv_id=$(echo "$mv_response" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    echo "Created model version with ID: $mv_id"

    # Create model artifact
    curl --silent -X 'POST' \
      "$mr_hostname/api/model_registry/v1alpha3/model_versions/$mv_id/artifacts" \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d '{
      "description": "Model artifact for Iris v1",
      "uri": "gs://kfserving-examples/models/sklearn/1.0/model",
      "state": "UNKNOWN",
      "name": "sklearn-iris-test-v1-scenario-4",
      "modelFormatName": "sklearn",
      "modelFormatVersion": "1",
      "artifactType": "model-artifact"
    }'

    echo "Starting test for scenario ${scenario_num}..."

    kubectl apply -n $kserve_namespace -f - <<EOF
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: "sklearn-iris-scenario-four"
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      storageUri: "model-registry://$MODEL_REGISTRY_SERVICE.${TESTNAMESPACE}.svc.cluster.local:$MODEL_REGISTRY_REST_PORT/iris-test-scenario-4"
EOF

    # wait for pod predictor to be initialized
    repeat_cmd_until "kubectl get pod -n $kserve_namespace --selector='component=predictor' --output jsonpath='{.items[*].metadata.name}' | grep sklearn-iris-scenario-four-predictor | wc -l" "-gt 0" 60
    local predictor=$(kubectl get pod -n $kserve_namespace --selector="component=predictor" --output jsonpath='{.items[0].metadata.name}')

    kubectl wait --for=condition=Ready pod/$predictor -n $kserve_namespace --timeout=5m

    kubectl wait --for=jsonpath='{.status.url}' inferenceservice/sklearn-iris-scenario-four -n $kserve_namespace --timeout=5m
    sleep 5

    local service_hostname=$(kubectl get inferenceservice sklearn-iris-scenario-four -n $kserve_namespace -o jsonpath='{.status.url}' | cut -d "/" -f 3)
    local result=$(curl -s -H "Host: ${service_hostname}" -H "Content-Type: application/json" "http://${INGRESS_HOST}/v1/models/sklearn-iris-scenario-four:predict" -d @/tmp/iris-input.json)
    echo "Scenario ${scenario_num} received: $result"

    if [ ! "$result" = "{\"predictions\":[1,1]}" ]; then
        echo "Scenario ${scenario_num} - Prediction does not match expectation!"
        echo "Printing some logs for debugging.."
        kubectl logs pod/$predictor -n $kserve_namespace -c storage-initializer || true
        kubectl logs pod/$predictor -n $kserve_namespace -c kserve-container || true
        kill $pf_pid 2>/dev/null || true
        return 1
    else
        echo "Scenario ${scenario_num} - Test succeeded!"
        kubectl logs pod/$predictor -n $kserve_namespace -c storage-initializer || true
        kubectl logs pod/$predictor -n $kserve_namespace -c kserve-container || true
    fi

    echo "Cleaning up scenario ${scenario_num}..."
    kubectl delete inferenceservice sklearn-iris-scenario-four -n $kserve_namespace

    sleep 5
    repeat_cmd_until "kubectl get pod -n $kserve_namespace --selector='component=predictor' --output jsonpath='{.items[*].metadata.name}' | wc -w" "= 0" 60

    kill $pf_pid 2>/dev/null || true

    echo "======== Finished Scenario ${scenario_num} ========"
    return 0
}

# Launch all scenarios in parallel
echo "======== Launching all scenarios in parallel ========"

run_scenario_1 &
scenario_pids[1]=$!

run_scenario_2 &
scenario_pids[2]=$!

run_scenario_3 &
scenario_pids[3]=$!

run_scenario_4 &
scenario_pids[4]=$!

# Wait for all scenarios and collect results
echo "======== Waiting for all scenarios to complete ========"

all_passed=true
for i in 1 2 3 4; do
    if wait ${scenario_pids[$i]}; then
        echo "✅ Scenario $i passed"
        scenario_results[$i]=0
    else
        echo "❌ Scenario $i failed"
        scenario_results[$i]=1
        all_passed=false
    fi
done

# Cleanup ingress port-forward
kill $ingress_pf_pid 2>/dev/null || true

# Print summary
echo "========================================"
echo "Test Summary:"
echo "========================================"
for i in 1 2 3 4; do
    if [ ${scenario_results[$i]} -eq 0 ]; then
        echo "Scenario $i: ✅ PASSED"
    else
        echo "Scenario $i: ❌ FAILED"
    fi
done
echo "========================================"

if [ "$all_passed" = true ]; then
    echo "🎉 All tests passed!"
    exit 0
else
    echo "❌ Some tests failed!"
    exit 1
fi

