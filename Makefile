model-registry: build

pkg/ml_metadata/proto/%.pb.go: api/grpc/ml_metadata/proto/%.proto
	protoc -I./api/grpc --go_out=./pkg --go_opt=paths=source_relative \
		--go-grpc_out=./pkg --go-grpc_opt=paths=source_relative $<

.PHONY: gen/grpc
gen/grpc: pkg/ml_metadata/proto/metadata_store.pb.go pkg/ml_metadata/proto/metadata_store_service.pb.go

.PHONY: gen/graph
gen/graph: pkg/model/graph/models_gen.go

pkg/model/graph/models_gen.go: api/graphql/*.graphqls gqlgen.yml
	go run github.com/99designs/gqlgen generate

.PHONY: vet
vet:
	go vet ./...

.PHONY: clean
clean:
	rm -Rf ./model-registry pkg/ml_metadata/proto/*.go pkg/model/graph/*.go

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: build
build: gen vet lint
	go build

.PHONY: gen
gen: gen/grpc gen/graph

.PHONY: lint
lint: gen
	golangci-lint run cmd/... pkg/...

.PHONY: run/migrate
run/migrate: gen
	go run main.go migrate --logtostderr=true

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
