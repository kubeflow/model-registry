#!/usr/bin/env bash
set -euxo pipefail

OSTYPE=$(uname -s)
OS="linux"
if [[ $OSTYPE =~ [Dd]arwin ]]; then
    OS="osx"
fi

ARCH=$(uname -m)
if [[ "$ARCH" == "arm"* || "$ARCH" == "aarch64" ]]; then
  ARCH="aarch_64"
elif [[ "$ARCH" == "s390x" ]]; then
  ARCH="s390_64"
elif [[ "$ARCH" == "ppc64le" ]] ; then
  ARCH="ppcle_64"
fi

PROJECT_ROOT=$(realpath "$(dirname "$0")"/..)

VERSION="24.3"
URL=https://github.com/protocolbuffers/protobuf/releases/download/v${VERSION}/protoc-${VERSION}-${OS}-${ARCH}.zip
wget -qx "$URL" -O "$PROJECT_ROOT"/protoc.zip &&
    unzip -qo "$PROJECT_ROOT"/protoc.zip -d "$PROJECT_ROOT" &&
    "$PROJECT_ROOT"/bin/protoc --version &&
    rm "$PROJECT_ROOT"/protoc.zip
