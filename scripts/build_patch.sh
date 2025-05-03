#!/usr/bin/env bash

set -e

PROJECT_ROOT=$(realpath "$(dirname "$0")"/..)

PATCH_FILE_PATH="${PROJECT_ROOT}/patches/type_asserts.patch"

ASSERT_FILE_PATH="${PROJECT_ROOT}/internal/server/openapi/type_asserts.go"

if [[ "$1" == "start" ]]; then
    git add "$ASSERT_FILE_PATH"
    git commit -m "Update type asserts for patch" --signoff
    echo "make changes to the type_asserts.go file now"
elif [[ "$1" == "finish" ]]; then
    git diff "$ASSERT_FILE_PATH" > "$PATCH_FILE_PATH"
else
    echo "Usage: $0 {start|finish}"
    exit 1
fi
