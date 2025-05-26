#!/bin/bash

set -e

cd "$(pwd)/$(dirname "$0")/.."

if [ -z "$YQ" ]; then
  if [ -e "bin/yq" ]; then
    YQ="$(realpath "bin/yq")"
  else
    echo "Error: YQ is not set and bin/yq does not exist" >&2
    exit 1
  fi
fi

usage() {
    echo "Usage: $0 [--check] <basename.yaml>"
    echo "  --check: Check for differences in the generated OpenAPI spec."
    exit 0
}

CHECK=false
BASENAME=""
while [[ $# -gt 0 ]]; do
    case "$1" in
        --check)
            CHECK=true
            shift
            ;;
        -h|--help)
            usage
            ;;
        *)
            if [[ "${1#-}" != "$1" ]]; then
                echo "Unknown option: $1"
                usage
            fi
            if [[ "$BASENAME" != "" ]]; then
                usage
            fi

            BASENAME=$1
            shift
            ;;
    esac
done

if [[ "$BASENAME" == "" ]]; then
    usage
fi

BASENAME=$(basename $BASENAME)
SOURCE_FILE="api/openapi/src/$BASENAME"
if [[ ! -f "$SOURCE_FILE" ]]; then
    echo "No source file at $SOURCE_FILE"
    exit 1
fi

OUT_FILE="api/openapi/$BASENAME"
if [[ "$CHECK" == "true" ]]; then
    OUT_FILE="$(mktemp -t modelregistry_openapi_tempXXXXXX).yaml"
    trap 'rm "$OUT_FILE"' EXIT
fi

# Merge the src files together.
$YQ eval-all '. as $item ireduce ({}; . * $item )' $SOURCE_FILE api/openapi/src/lib/*.yaml >"$OUT_FILE"

# Re-order the keys in the generated file.
$YQ eval -i '
    {
        "openapi": .openapi,
        "info": .info,
        "servers": .servers,
        "paths": .paths,
        "components": .components,
        "security": .security,
        "tags": .tags
    } |
        sort_keys(.paths) |
        sort_keys(.components.schemas) |
        sort_keys(.components.responses)
' "$OUT_FILE"

if [[ "$CHECK" == "true" ]]; then
    exec diff -u api/openapi/$BASENAME $OUT_FILE
fi
