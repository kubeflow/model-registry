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

function sed_inplace() {
    if [[ $(uname) == "Darwin" ]]; then
        # introduce -i parameter for Mac OSX sed compatibility
        sed -i '' "$@"
    else
        sed -i "$@"
    fi
}

sed_inplace 's/, orderByParam/, model.OrderByField(orderByParam)/g' "$PROJECT_ROOT"/internal/server/openapi/api_model_catalog_service.go
sed_inplace 's/, sortOrderParam/, model.SortOrder(sortOrderParam)/g' "$PROJECT_ROOT"/internal/server/openapi/api_model_catalog_service.go

sed_inplace 's/"encoding\/json"//' "$PROJECT_ROOT"/internal/server/openapi/api_model_catalog_service.go

sed_inplace 's/github.com\/kubeflow\/model-registry\/pkg\/openapi/github.com\/kubeflow\/model-registry\/catalog\/pkg\/openapi/' \
    "$PROJECT_ROOT"/internal/server/openapi/api_model_catalog_service.go \
    "$PROJECT_ROOT"/internal/server/openapi/api.go
sed_inplace 's/{\?model_name+}\?/*/' "$PROJECT_ROOT"/internal/server/openapi/api_model_catalog_service.go

echo "Assembling type_assert Go file"
./scripts/gen_type_asserts.sh "$DST"

gofmt -w "$DST"

echo "OpenAPI server generation completed"
