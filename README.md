# model-registry
A go based server that implements a gRPC interface for [ml_metadata](https://github.com/google/ml-metadata/) library.
It adds other features on top of the functionality offered by the gRPC interface.
## Pre-requisites:
- go >= 1.19
- protoc - [Protocol buffer compiler](https://grpc.io/docs/protoc-installation/).
- go tools - Installed with the following commands:
```
go install github.com/99designs/gqlgen@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```
- gRPC go plugins - Installed with the following commands:
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
- python 3.9
## Building
Run the following command to build the server binary:
```
make build
```
The generated binary uses spf13 cmdline args. More information on using the server can be obtained by running the command:
```
./model-registry --help
```
## Creating/Migrating Server DB
The server uses a local SQLite DB file (`metadata.sqlite.db` by default), which can be configured using the `-d` cmdline option.
The following command creates the DB:
```
./model-registry migrate
```
### Loading metadata library
Run the following command to load a metadata library:
```
./model-registry migrate -m config/metadata-library
```
Note that currently no duplicate detection is done as the implementation is a WIP proof of concept. 
Running this command multiple times will create duplicate metadata types. 
To clear the DB simply delete the SQLite DB file `metadata.sqlite.db`. 

### Running Server
Run the following command to start the server:
```
make run/server &
```
### Running Python ml-metadata test client
Before running the test client, install the required Python libraries (using a python venv, if using one) 
using the command:
```
pip install ml_metadata grpcio
```
Run the following command to run the ml-metadata Python test client:
```
make run/client
```
### Running GraphQL Playground
This project includes support for a GraphiQL playground, which supports interactive query design. 
It can be reached by opening the following URL in a web browser:
```
http://localhost:8080/
```
Where, 8080 is the default port that the server listens on. This port can be changed with the `-p` option.  
### Clean
Run the following command to clean the DB file, generated gRPC and GraphQL models, etc.:
```
make clean
```