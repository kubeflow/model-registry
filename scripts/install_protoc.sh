#! /bin/bash
set -euxo pipefail

SCRIPT_PARENT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )/.." &> /dev/null && pwd )

VERSION="24.3"
OS="linux"
if [[ "$OSTYPE" == "darwin"* ]]; then
  # Mac OSX
  OS="osx"
fi
ARCH="x86_64"
if [[ "$(uname -m)" == "arm"* ]]; then
  ARCH="aarch_64"
fi

mkdir -p ${SCRIPT_PARENT_DIR}/bin

wget -q https://github.com/protocolbuffers/protobuf/releases/download/v${VERSION}/protoc-${VERSION}-${OS}-${ARCH}.zip -O ${SCRIPT_PARENT_DIR}/protoc.zip && \
  unzip -qo ${SCRIPT_PARENT_DIR}/protoc.zip -d ${SCRIPT_PARENT_DIR} && \
  bin/protoc --version && \
  rm ${SCRIPT_PARENT_DIR}/protoc.zip
