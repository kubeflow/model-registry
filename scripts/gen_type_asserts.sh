#!/usr/bin/env bash

set -e

PROJECT_ROOT=$(realpath "$(dirname "$0")"/..)

ASSERT_FILE_PATH="${PROJECT_ROOT}/internal/server/openapi/type_asserts.go"
PATCH="${PROJECT_ROOT}/patches/type_asserts.patch"

python3 "${PROJECT_ROOT}/scripts/gen_type_asserts.py" >"$ASSERT_FILE_PATH"

gofmt -w "$ASSERT_FILE_PATH"

git apply "$PATCH"
