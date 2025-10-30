#!/usr/bin/env bash

set -e

echo "Generating the OpenAPI server"

OPENAPI_GENERATOR=${OPENAPI_GENERATOR:-openapi-generator-cli}

PROJECT_ROOT=$(realpath "$(dirname "$0")"/..)
SRC="$PROJECT_ROOT/${1:-../api/openapi/catalog.yaml}"
DST="$PROJECT_ROOT/${2:-internal/server/openapi}"

"$OPENAPI_GENERATOR" generate \
    -i "$SRC" -g go-server -o "$DST" --package-name openapi \
    --ignore-file-override "$PROJECT_ROOT"/.openapi-generator-ignore --additional-properties=outputAsLibrary=true,enumClassPrefix=true,router=chi,sourceFolder=,onlyInterfaces=true,isGoSubmodule=true,enumClassPrefix=true,useOneOfDiscriminatorLookup=true,featureCORS=true \
    --template-dir "$PROJECT_ROOT"/../templates/go-server

# Python-based regex replace function
# Usage: py-re-replace <count> <pattern> <replacement> <file1> [file2...]
# count=0: replace all occurrences (like sed with /g flag)
# count=1: replace first occurrence only (like sed without /g flag)
# count=N: replace first N occurrences
py-re-replace() {
  python3 -c "
import fileinput, re, sys
count, pattern, replacement, filepaths = int(sys.argv[1]), sys.argv[2], sys.argv[3], sys.argv[4:]
for filepath in filepaths:
    for line in fileinput.FileInput(filepath, inplace=True, backup=''):
        sys.stdout.write(re.sub(pattern, replacement, line, count=count))
" "$@"
}

py-re-replace 0 ', orderByParam' ', model.OrderByField(orderByParam)' "$PROJECT_ROOT"/internal/server/openapi/api_model_catalog_service.go
py-re-replace 0 ', sortOrderParam' ', model.SortOrder(sortOrderParam)' "$PROJECT_ROOT"/internal/server/openapi/api_model_catalog_service.go

py-re-replace 0 'model\.\[\]ArtifactTypeQueryParam' '[]model.ArtifactTypeQueryParam' "$PROJECT_ROOT"/internal/server/openapi/api.go
py-re-replace 0 'model\.\[\]ArtifactType2QueryParam' '[]model.ArtifactTypeQueryParam' "$PROJECT_ROOT"/internal/server/openapi/api.go

py-re-replace 1 'github\.com/kubeflow/model-registry/pkg/openapi' 'github.com/kubeflow/model-registry/catalog/pkg/openapi' \
    "$PROJECT_ROOT"/internal/server/openapi/api_model_catalog_service.go \
    "$PROJECT_ROOT"/internal/server/openapi/api.go

py-re-replace 1 '\{model_name\+\}|model_name\+' '*' "$PROJECT_ROOT"/internal/server/openapi/api_model_catalog_service.go

echo "Applying patches to generated code"
(
    cd "$PROJECT_ROOT/.."
    git apply patches/api_model_catalog_service.patch
)

echo "Assembling type_assert Go file"
./scripts/gen_type_asserts.sh "$DST"

$PROJECT_ROOT/../bin/goimports -w "$DST"

echo "OpenAPI server generation completed"
