#!/bin/bash

set -e

# quay.io credentials
QUAY_REGISTRY=quay.io
QUAY_ORG="${QUAY_ORG:-opendatahub}"
QUAY_IMG_REPO="${QUAY_IMG_REPO:-model-registry}"
QUAY_USERNAME="${QUAY_USERNAME}"
QUAY_PASSWORD="${QUAY_PASSWORD}"

# image version
HASH="$(git rev-parse --short=7 HEAD)"
VERSION="${VERSION:-$HASH}"

# if set to 0 skip image build
# otherwise build it
BUILD_IMAGE="${BUILD_IMAGE:-true}"

# if set to 0 skip push to registry
# otherwise push it
PUSH_IMAGE="${PUSH_IMAGE:-false}"

# skip if image already existing on registry
SKIP_IF_EXISTING="${SKIP_IF_EXISTING:-false}"

# assure docker exists
docker -v foo >/dev/null 2>&1 || { echo >&2 "::error:: Docker is required.  Aborting."; exit 1; }

# skip if image already existing
if [[ "${SKIP_IF_EXISTING,,}" == "true" ]]; then
  TAGS=$(curl --request GET "https://$QUAY_REGISTRY/api/v1/repository/${QUAY_ORG}/${QUAY_IMG_REPO}/tag/?specificTag=${VERSION}")
  LATEST_TAG_HAS_END_TS=$(echo $TAGS | jq .tags - | jq 'sort_by(.start_ts) | reverse' | jq '.[0].end_ts')
  NOT_EMPTY=$(echo ${TAGS} | jq .tags - | jq any)

  # Image only exists if there is a tag that does not have "end_ts" (i.e. it is still present).
  if [[ "$NOT_EMPTY" == "true" && $LATEST_TAG_HAS_END_TS == "null" ]]; then
      echo "::error:: The image ${QUAY_ORG}/${QUAY_IMG_REPO}:${VERSION} already exists"
      exit 1
  else
      echo "Image does not exist...proceeding with build & push."
  fi
fi

# build docker image, login is not required at this step
if [[ "${BUILD_IMAGE,,}" == "true" ]]; then
  echo "Building container image.."
  make \
    IMG_REGISTRY="${QUAY_REGISTRY}" \
    IMG_ORG="${QUAY_ORG}" \
    IMG_REPO="${QUAY_IMG_REPO}" \
    IMG_VERSION="${VERSION}" \
    image/build
else
  echo "Skip container image build."
fi

# push container image to registry, requires login
if [[ "${PUSH_IMAGE,,}" == "true" ]]; then
  echo "Pushing container image.."
  make \
    IMG_REGISTRY="${QUAY_REGISTRY}" \
    IMG_ORG="${QUAY_ORG}" \
    IMG_REPO="${QUAY_IMG_REPO}" \
    IMG_VERSION="${VERSION}" \
    DOCKER_USER="${QUAY_USERNAME}"\
    DOCKER_PWD="${QUAY_PASSWORD}" \
    docker/login \
    image/push
else
  echo "Skip container image push."
fi
