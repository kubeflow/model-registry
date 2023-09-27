# model-registry
A go based server that implements a gRPC interface for [ml_metadata](https://github.com/google/ml-metadata/) library.
It adds other features on top of the functionality offered by the gRPC interface.
## Pre-requisites:
- go >= 1.19
- protoc v24.3 - [Protocol Buffers v24.3 Release](https://github.com/protocolbuffers/protobuf/releases/tag/v24.3)
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
## Server
### Creating/Migrating Server DB
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

### Starting the Server
Run the following command to start the server:
```
make run/server &
```
## Clients
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
Run the following command to clean the server binary, generated gRPC and GraphQL models, etc.:
```
make clean
```
## Docker Image
### Building the Docker Image
The following command builds a docker image for the server with the tag `model-registry``:
```shell
docker build -t model-registry .
```
Note that the first build will be longer as it downloads the build tool dependencies. 
Subsequent builds will re-use the cached tools layer. 
### Creating/Migrating Server DB
The following command migrates or creates a DB for the server:
```shell
docker run -it --user <uid>:<gid> -v <host-path>:/var/db model-registry migrate -d /var/db/metadata.sqlite.db -m /config/metadata-library
```
Where, `<uid>` and `<gid>` are local user and group ids on the host machine to allow volume mapping for the DB files. 
And, `<host-path>` is the path on the local directory writable by the `<uid>:<gid>` user. 
### Running the Server
The following command starts the server:
```shell
docker run -d -p <hostname>:<port>:8080 --user <uid>:<gid> -v <host-path>:/var/db --name server model-registry serve -n 0.0.0.0 -d /var/db/metadata.sqlite.db
```
Where, `<uid>`, `<gid>`, and `<host-path>` are the same as in the migrate command above. 
And `<hostname>` and `<port>` are the local ip and port to use to expose the container's default `8080` listening port. 
The server listens on `localhost` by default, hence the `-n 0.0.0.0` option allows the server port to be exposed. 

Once the server has started, test clients and playground can be used as described in the above sections. 
