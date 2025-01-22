#!/usr/bin/env bash

set -e

echo "Generating the OpenAPI server"

PROJECT_ROOT=$(realpath "$(dirname "$0")"/..)

openapi-generator-cli generate \
    -i "$PROJECT_ROOT"/api/openapi/model-registry.yaml -g go-server -o "$PROJECT_ROOT"/internal/server/openapi --package-name openapi --global-property models,apis \
    --ignore-file-override "$PROJECT_ROOT"/.openapi-generator-ignore --additional-properties=outputAsLibrary=true,enumClassPrefix=true,router=chi,sourceFolder=,onlyInterfaces=true,isGoSubmodule=true,enumClassPrefix=true,useOneOfDiscriminatorLookup=true \
    --template-dir "$PROJECT_ROOT"/templates/go-server

if [[ $(uname) == "Darwin" ]]; then
    # introduce -i parameter for Mac OSX sed compatibility
    sed -i '' 's/, orderByParam/, model.OrderByField(orderByParam)/g' "$PROJECT_ROOT"/internal/server/openapi/api_model_registry_service.go
    sed -i '' 's/, sortOrderParam/, model.SortOrder(sortOrderParam)/g' "$PROJECT_ROOT"/internal/server/openapi/api_model_registry_service.go
    sed -i '' 's/, orderByParam/, model.OrderByField(orderByParam)/g' "$PROJECT_ROOT"/internal/server/openapi/api_model_catalog_service.go
    sed -i '' 's/, sortOrderParam/, model.SortOrder(sortOrderParam)/g' "$PROJECT_ROOT"/internal/server/openapi/api_model_catalog_service.go
else
    sed -i 's/, orderByParam/, model.OrderByField(orderByParam)/g' "$PROJECT_ROOT"/internal/server/openapi/api_model_registry_service.go
    sed -i 's/, sortOrderParam/, model.SortOrder(sortOrderParam)/g' "$PROJECT_ROOT"/internal/server/openapi/api_model_registry_service.go
    sed -i 's/, orderByParam/, model.OrderByField(orderByParam)/g' "$PROJECT_ROOT"/internal/server/openapi/api_model_catalog_service.go
    sed -i 's/, sortOrderParam/, model.SortOrder(sortOrderParam)/g' "$PROJECT_ROOT"/internal/server/openapi/api_model_catalog_service.go
fi

echo "Assembling type_assert Go file"
./scripts/gen_type_asserts.sh

gofmt -w "$PROJECT_ROOT"/internal/server/openapi

echo "OpenAPI server generation completed"
