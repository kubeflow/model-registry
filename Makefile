# useful paths
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
PROJECT_BIN := $(PROJECT_PATH)/bin

# add tools bin directory
PATH := $(PROJECT_BIN):$(PATH)
model-registry: build

internal/ml_metadata/proto/%.pb.go: api/grpc/ml_metadata/proto/%.proto
	protoc -I./api/grpc --go_out=./internal --go_opt=paths=source_relative \
		--go-grpc_out=./internal --go-grpc_opt=paths=source_relative $<

.PHONY: gen/grpc
gen/grpc: internal/ml_metadata/proto/metadata_store.pb.go internal/ml_metadata/proto/metadata_store_service.pb.go

internal/converter/generated/converter.go: internal/converter/*.go
	goverter -packageName generated -output ./internal/converter/generated/converter.go github.com/opendatahub-io/model-registry/internal/converter/

.PHONY: gen/converter
gen/converter: gen/grpc gen/graph internal/converter/generated/converter.go

.PHONY: gen/graph
gen/graph: internal/model/graph/models_gen.go

internal/model/graph/models_gen.go: api/graphql/*.graphqls gqlgen.yml
	gqlgen generate

.PHONY: vet
vet:
	go vet ./...

.PHONY: clean
clean:
	rm -Rf ./model-registry internal/ml_metadata/proto/*.go internal/model/graph/models_gen.go internal/converter/generated/converter.go

bin/go-enum:
	GOBIN=$(PROJECT_BIN) go install github.com/searKing/golang/tools/go-enum@v1.2.97

bin/protoc-gen-go:
	GOBIN=$(PROJECT_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0

bin/protoc-gen-go-grpc:
	GOBIN=$(PROJECT_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

bin/gqlgen:
	GOBIN=$(PROJECT_BIN) go install github.com/99designs/gqlgen@v0.17.36

bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v1.54.2

bin/goverter:
	GOBIN=$(PROJECT_PATH)/bin go install github.com/jmattheis/goverter/cmd/goverter@v0.18.0

.PHONY: deps
deps: bin/go-enum bin/protoc-gen-go bin/protoc-gen-go-grpc bin/gqlgen bin/golangci-lint bin/goverter

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: build
build: gen vet lint
	go build

.PHONY: gen
gen: deps gen/grpc gen/graph gen/converter
	go generate ./...

.PHONY: lint
lint: gen
	golangci-lint run main.go
	golangci-lint run cmd/... internal/...

.PHONY: test
test: gen
	go test ./internal/...

.PHONY: run/migrate
run/migrate: gen
	go run main.go migrate --logtostderr=true -m config/metadata-library

metadata.sqlite.db: run/migrate

.PHONY: run/server
run/server: gen metadata.sqlite.db
	go run main.go serve --logtostderr=true

.PHONY: run/client
run/client: gen
	python test/python/test_mlmetadata.py

.PHONY: serve
serve: build
	./model-registry serve --logtostderr=true

all: model-registry
