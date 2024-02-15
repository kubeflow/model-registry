#!/bin/bash

set -e

echo "Generating the OpenAPI server"

ROOT_FOLDER="${ROOT_FOLDER:-..}"

openapi-generator-cli generate \
		-i $ROOT_FOLDER/api/openapi/model-registry.yaml -g go-server -o $ROOT_FOLDER/internal/server/openapi --package-name openapi --global-property models,apis \
		--ignore-file-override $ROOT_FOLDER/.openapi-generator-ignore --additional-properties=outputAsLibrary=true,enumClassPrefix=true,router=chi,sourceFolder=,onlyInterfaces=true,isGoSubmodule=true,enumClassPrefix=true,useOneOfDiscriminatorLookup=true \
		--template-dir $ROOT_FOLDER/templates/go-server

sed -i 's/, orderByParam/, model.OrderByField(orderByParam)/g' $ROOT_FOLDER/internal/server/openapi/api_model_registry_service.go
sed -i 's/, sortOrderParam/, model.SortOrder(sortOrderParam)/g' $ROOT_FOLDER/internal/server/openapi/api_model_registry_service.go

echo "Assembling type_assert Go file"
./scripts/gen_type_asserts.sh

gofmt -w $ROOT_FOLDER/internal/server/openapi

echo "OpenAPI server generation completed"
