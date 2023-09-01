gen/grpc: api/grpc/ml_metadata/proto/*.proto
	protoc -I./api/grpc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		ml_metadata/proto/metadata_store.proto ml_metadata/proto/metadata_store_service.proto

clean:
	rm -Rf ml_metadata/proto/*.go

build: gen/grpc
	go build

run/migrate: gen/grpc
	go run main.go migrate --logtostderr=true

run/server: gen/grpc
	go run main.go serve --logtostderr=true

run/client: gen/grpc
	python3.9 test/python/test_mlmetadata.py

serve: build
	./ml-metadata-go-server serve --logtostderr=true
