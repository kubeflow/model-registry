#! /bin/bash
set -euxo pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

VERSION="24.3"
OS="linux"
if [[ "$OSTYPE" == "darwin"* ]]; then
  # Mac OSX
  OS="osx"
fi
ARCH=$(uname -m)
if [[ "$ARCH" == "arm"* ]]; then
  ARCH="aarch_64"
elif [[ "$ARCH" == "s390x" ]]; then
  ARCH="s390_64"
elif [[ "$ARCH" == "ppc64le" ]] ; then
  ARCH="ppcle_64"
fi

mkdir -p ${SCRIPT_DIR}/../bin

wget -q https://github.com/protocolbuffers/protobuf/releases/download/v${VERSION}/protoc-${VERSION}-${OS}-${ARCH}.zip -O ${SCRIPT_DIR}/../protoc.zip && \
  unzip -qo ${SCRIPT_DIR}/../protoc.zip -d ${SCRIPT_DIR}/.. && \
  bin/protoc --version && \
  rm ${SCRIPT_DIR}/../protoc.zip
