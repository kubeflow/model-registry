ml-metadata-go-server: build

ml_metadata/proto/%.pb.go: api/grpc/ml_metadata/proto/%.proto
	protoc -I./api/grpc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative $<

.PHONY: gen/grpc
gen/grpc: ml_metadata/proto/metadata_store.pb.go ml_metadata/proto/metadata_store_service.pb.go

.PHONY: clean
clean:
	rm -Rf ml_metadata/proto/*.go ./ml-metadata-go-server

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: build
build: gen/grpc
	go build

.PHONY: gen
gen: gen/grpc

.PHONY: run/migrate
run/migrate: gen/grpc
	go run main.go migrate --logtostderr=true

.PHONY: run/server
run/server: gen/grpc
	go run main.go serve --logtostderr=true

.PHONY: run/client
run/client: gen/grpc
	python3.9 test/python/test_mlmetadata.py

.PHONY: serve
serve: build
	./ml-metadata-go-server serve --logtostderr=true

all: ml-metadata-go-server
