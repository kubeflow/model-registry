#!/usr/bin/env bash

set -e

PROJECT_ROOT=$(realpath "$(dirname "$0")"/..)

ASSERT_FILE_PATH="${PROJECT_ROOT}/internal/server/openapi/type_asserts.go"
PATCH="${PROJECT_ROOT}/patches/type_asserts.patch"

# Remove the existing file identified by env.ASSERT_FILE_PATH
if [ -f "$ASSERT_FILE_PATH" ]; then
    rm "$ASSERT_FILE_PATH"
fi

# Create an empty file
touch "$ASSERT_FILE_PATH"

INITIAL_CONTENT=$(
    cat <<EOF
/*
 * Model Registry REST API
 *
 * REST API for Model Registry to create and manage ML model metadata
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 *
 */

// File generated by scripts/gen_type_assert.sh - DO NOT EDIT

package openapi

import (
	model "github.com/kubeflow/model-registry/pkg/openapi"
)


EOF
)

# Create the file and initialize it with the specified content
echo -e "$INITIAL_CONTENT" >"$ASSERT_FILE_PATH"

# Iterate over files starting with "model_" in the internal/server/openapi/ folder
for file in "$PROJECT_ROOT"/internal/server/openapi/model_*; do
    # Check if the file is a regular file
    if [ -f "$file" ]; then
        # Ignore first 15 lines containing license, package and imports
        sed -n '13,$p' "$file" >>"$ASSERT_FILE_PATH"

        # Remove the merged file
        rm "$file"
    fi
done

python "$PROJECT_ROOT"/scripts/reorder_type_asserts.py "$PROJECT_ROOT"/internal/server/openapi/type_asserts.go

gofmt -w "$ASSERT_FILE_PATH"

git apply "$PATCH"
