# Common ODH environment variables for test targets.
# Usage: . scripts/odh_env.sh
#
# Exported variables:
#   AUTH_TOKEN, VERIFY_SSL, MR_NAMESPACE, MR_HOST_URL, MR_ENDPOINT,
#   MODEL_SYNC_REGISTRY_SERVER_ADDRESS, MODEL_SYNC_REGISTRY_PORT,
#   MODEL_SYNC_REGISTRY_IS_SECURE, MODEL_SYNC_REGISTRY_USER_TOKEN,
#   CONTAINER_IMAGE_URI

SCRIPT_DIR="$(dirname "$(realpath "${BASH_SOURCE[0]}")")"

AUTH_TOKEN=$(kubectl config view --raw -o jsonpath="{.users[?(@.name==\"$(kubectl config view -o jsonpath="{.contexts[?(@.name==\"$(kubectl config current-context)\")].context.user}")\")].user.token}")
export AUTH_TOKEN

export VERIFY_SSL=False

MR_NAMESPACE=$(kubectl get datasciencecluster default-dsc -o jsonpath='{.spec.components.modelregistry.registriesNamespace}')
export MR_NAMESPACE

MR_ENDPOINT=$(kubectl get service -n "${MR_NAMESPACE}" model-registry -o jsonpath='{.metadata.annotations.routing\.opendatahub\.io\/external-address-rest}')
export MR_HOST_URL="https://${MR_ENDPOINT}"
MR_ENDPOINT="${MR_ENDPOINT%%:*}"
export MR_ENDPOINT

export MODEL_SYNC_REGISTRY_SERVER_ADDRESS="https://${MR_ENDPOINT}"
export MODEL_SYNC_REGISTRY_PORT="443"
export MODEL_SYNC_REGISTRY_IS_SECURE="false"
export MODEL_SYNC_REGISTRY_USER_TOKEN="${AUTH_TOKEN}"
CONTAINER_IMAGE_URI=$("${SCRIPT_DIR}/get_async_upload_image.sh")
export CONTAINER_IMAGE_URI
