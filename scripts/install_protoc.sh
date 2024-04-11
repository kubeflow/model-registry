#! /bin/bash
set -euxo pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

VERSION="24.3"
OS="linux"
if [[ "$OSTYPE" == "darwin"* ]]; then
  # Mac OSX
  OS="osx"
fi
ARCH="x86_64"
if [[ "$(uname -m)" == "arm"* ]]; then
  ARCH="arm64"
fi

mkdir -p ${SCRIPT_DIR}/../bin

wget -q https://github.com/protocolbuffers/protobuf/releases/download/v${VERSION}/protoc-${VERSION}-${OS}-${ARCH}.zip -O ${SCRIPT_DIR}/../protoc.zip && \
  unzip -qo ${SCRIPT_DIR}/../protoc.zip -d ${SCRIPT_DIR}/.. && \
  bin/protoc --version && \
  rm ${SCRIPT_DIR}/../protoc.zip
