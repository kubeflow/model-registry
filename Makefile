model-registry: build

internal/ml_metadata/proto/%.pb.go: api/grpc/ml_metadata/proto/%.proto
	protoc -I./api/grpc --go_out=./internal --go_opt=paths=source_relative \
		--go-grpc_out=./internal --go-grpc_opt=paths=source_relative $<

.PHONY: gen/grpc
gen/grpc: internal/ml_metadata/proto/metadata_store.pb.go internal/ml_metadata/proto/metadata_store_service.pb.go

.PHONY: gen/graph
gen/graph: internal/model/graph/models_gen.go

internal/model/graph/models_gen.go: api/graphql/*.graphqls gqlgen.yml
	go run github.com/99designs/gqlgen generate

.PHONY: vet
vet:
	go vet ./...

.PHONY: clean
clean:
	rm -Rf ./model-registry internal/ml_metadata/proto/*.go internal/model/graph/*.go

.PHONY: deps
deps:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2
	go install github.com/99designs/gqlgen@latest
	go install github.com/searKing/golang/tools/go-enum@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: build
build: gen vet lint
	go build

.PHONY: gen
gen: gen/grpc gen/graph
	go generate ./...

.PHONY: lint
lint: gen
	golangci-lint run main.go
	golangci-lint run cmd/... internal/...

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
