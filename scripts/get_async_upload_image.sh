#!/usr/bin/env bash
# Reads the async-upload job image from the model-registry-operator-parameters
# ConfigMap in the applications namespace (from DSCInitialization).
set -e

APPS_NS=$(kubectl get dscinitializations default-dsci \
  -o jsonpath='{.spec.applicationsNamespace}')

kubectl get configmap model-registry-operator-parameters \
  -n "$APPS_NS" \
  -o jsonpath='{.data.IMAGES_JOBS_ASYNC_UPLOAD}'
