#!/usr/bin/env bash

set -e

ASSERT_FILE_PATH="$1/type_asserts.go"

PROJECT_ROOT=$(realpath "$(dirname "$0")"/..)
PATCH="${PROJECT_ROOT}/patches/type_asserts.patch"

# AssertMetadataValueRequired from this file generates with the incorrect logic.
rm -f $1/model_metadata_value.go

python3 "${PROJECT_ROOT}/scripts/gen_type_asserts.py" $1 >"$ASSERT_FILE_PATH"

gofmt -w "$ASSERT_FILE_PATH"
