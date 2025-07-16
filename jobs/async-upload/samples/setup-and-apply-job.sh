#!/bin/bash

set -euo pipefail

MR_HOST_URL="${MR_HOST_URL:-http://localhost:8080}"
export CONTAINER_IMAGE_URI="${CONTAINER_IMAGE_URI:-ghcr.io/kubeflow/model-registry/job/async-upload:latest}"
export JOB_NAME="${JOB_NAME:-my-async-upload-job}"
echo "Received environment variables:"
echo "  MR_HOST_URL: $MR_HOST_URL"
echo "  CONTAINER_IMAGE_URI: $CONTAINER_IMAGE_URI"
echo "  JOB_NAME: $JOB_NAME"

echo "Creating top-level RegisteredModel..."
export MODEL_ID=$(curl --silent --request POST \
--url "$MR_HOST_URL/api/model_registry/v1alpha3/registered_models" \
--header 'content-type: application/json' \
--data "{\"name\": \"$(openssl rand -hex 4)\"}" | jq -r '.id')
echo "  MODEL_ID: $MODEL_ID"

echo "Creating ModelVersion associated with the RegisteredModel..."
export MODEL_VERSION_ID=$(curl --silent --request POST \
--url "$MR_HOST_URL/api/model_registry/v1alpha3/model_versions" \
--header 'content-type: application/json' \
--data "{\"name\": \"$(openssl rand -hex 4)\", \"registeredModelId\": \"$MODEL_ID\"}" | jq -r '.id')
echo "  MODEL_VERSION_ID: $MODEL_VERSION_ID"

echo "Creating placeholder ModelArtifact associated with the ModelVersion..."
export MODEL_ARTIFACT_ID=$(curl --silent --request POST \
--url "$MR_HOST_URL/api/model_registry/v1alpha3/model_versions/$MODEL_VERSION_ID/artifacts" \
--header 'content-type: application/json' \
--data "{\"uri\": \"PLACEHOLDER\", \"artifactType\": \"ModelArtifact\", \"state\": \"PENDING\"}" | jq -r '.id')
echo "  MODEL_ARTIFACT_ID: $MODEL_ARTIFACT_ID"


TEMP_DIR=$(mktemp -d)

# Clean up function
cleanup() {
    rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

# Download an mnist onnx file. Size ~25KB.
echo "Downloading mnist-8.onnx from GitHub..."
curl --silent -o $TEMP_DIR/mnist-8.onnx https://github.com/onnx/models/raw/refs/heads/main/validated/vision/classification/mnist/model/mnist-8.onnx

# Upload the onnx file to the built-in minio s3 bucket
echo "Uploading mnist-8.onnx to minio..."
S3_BUCKET=default
S3_KEY=my-model/mnist-8.onnx
S3_RESOURCE_PATH=/${S3_BUCKET}/${S3_KEY}
date=`date -R`
_signature="PUT\n\napplication/octet-stream\n${date}\n${S3_RESOURCE_PATH}"
signature=`echo -en ${_signature} | openssl sha1 -hmac minioadmin -binary | base64`
curl -X PUT -T $TEMP_DIR/mnist-8.onnx \
          -H "Host: localhost:9000" \
          -H "Date: ${date}" \
          -H "Content-Type: application/octet-stream" \
          -H "Authorization: AWS minioadmin:${signature}" \
          http://localhost:9000${S3_RESOURCE_PATH}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SAMPLES_DIR="$SCRIPT_DIR"
TEMP_PATCH="$TEMP_DIR/job-values.yaml"
TEMP_KUSTOMIZATION="$TEMP_DIR/kustomization.yaml"
TEMP_JOB_YAML="$TEMP_DIR/sample_job_s3_to_oci.yaml"

echo "Copying job yaml file to temp directory..."

# Copy the job yaml file to temp directory (kustomize security requirement)
cp "$SAMPLES_DIR/sample_job_s3_to_oci.yaml" "$TEMP_JOB_YAML"

# Substitute the model ids into the patch file
echo "Substituting model ids into patch file..."
envsubst < "$SAMPLES_DIR/patches/job-values.yaml" > "$TEMP_PATCH"

# Create temporary kustomization file using the generated patch file
echo "Creating temporary kustomization file..."
cat > "$TEMP_KUSTOMIZATION" << EOF
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- sample_job_s3_to_oci.yaml

patches:
- target:
    kind: Job
    name: my-async-upload-job
  path: job-values.yaml
EOF

# Clean up any existing job first (Jobs are immutable)
echo "Deleting existing job (if it exists)..."
kubectl delete job $JOB_NAME --ignore-not-found=true

# Apply the patched kustomization
echo "Applying job..."
kubectl apply -k $TEMP_DIR

# Wait for job completion and show logs
echo "Waiting for job completion..."
kubectl wait --for=condition=complete job/$JOB_NAME --timeout=10m

echo "Job completed successfully!"
