PROJECT_BIN := $(CURDIR)/../bin
OPENAPI_GENERATOR := $(PROJECT_BIN)/openapi-generator-cli
GENQLIENT_BIN ?= $(PROJECT_BIN)/genqlient
GENQLIENT_CONFIG := internal/catalog/genqlient/genqlient.yaml
GENQLIENT_OUTPUT := internal/catalog/genqlient/generated.go
GENQLIENT_SOURCES := $(wildcard internal/catalog/genqlient/queries/*.graphql)
GRAPHQL_SCHEMA := internal/catalog/genqlient/queries/schema.graphql

.PHONY: gen/openapi-server
gen/openapi-server: internal/server/openapi/api_model_catalog_service.go

internal/server/openapi/api_model_catalog_service.go: ../api/openapi/catalog.yaml
	OPENAPI_GENERATOR="$(OPENAPI_GENERATOR)" ./scripts/gen_openapi_server.sh

.PHONY: gen/openapi
gen/openapi: pkg/openapi/client.go

pkg/openapi/client.go: ../api/openapi/catalog.yaml
	${OPENAPI_GENERATOR} generate \
		-i ../api/openapi/catalog.yaml -g go -o pkg/openapi --package-name openapi \
		--ignore-file-override ./.openapi-generator-ignore --additional-properties=isGoSubmodule=true,enumClassPrefix=true,useOneOfDiscriminatorLookup=true
	gofmt -w pkg/openapi

.PHONY: gen/graphql
gen/graphql: $(GENQLIENT_OUTPUT)

$(GENQLIENT_OUTPUT): $(GENQLIENT_CONFIG) $(GENQLIENT_SOURCES) $(PROJECT_BIN)/genqlient
	$(GENQLIENT_BIN) --config $(GENQLIENT_CONFIG)

.PHONY: download/graphql-schema
download/graphql-schema:
	npx get-graphql-schema https://catalog.redhat.com/api/containers/graphql/ > $(GRAPHQL_SCHEMA)

.PHONY: clean-pkg-openapi
clean-pkg-openapi:
	while IFS= read -r file; do rm -f "pkg/openapi/$$file"; done < pkg/openapi/.openapi-generator/FILES

.PHONY: clean-graphql
clean-graphql:
	rm -f $(GENQLIENT_OUTPUT)

.PHONY: clean-internal-server-openapi
clean-internal-server-openapi:
	while IFS= read -r file; do rm -f "internal/server/openapi/$$file"; done < internal/server/openapi/.openapi-generator/FILES
