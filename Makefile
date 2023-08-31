gen/grpc: ml_metadata/proto/*.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		ml_metadata/proto/metadata_store.proto ml_metadata/proto/metadata_store_service.proto

clean:
	rm -Rf ml_metadata/*.go

build: gen/grpc
	go build

run/migrate: gen/grpc
	go run main.go migrate

run/server: gen/grpc
	go run main.go serve

run/client: gen/grpc
	python3.9 test/python/test_mlmetadata.py

serve: build
	./ml-metadata-go-server serve --logtostderr=true
