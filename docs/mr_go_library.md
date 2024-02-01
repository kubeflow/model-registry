# Model Registry Service

The Model Registry Service go library provides a convenient interface for managing and interacting with models, model versions, artifacts, serving environments, inference services, and serve models through the underlying [ML Metadata (MLMD)](https://github.com/google/ml-metadata) service.

## Installation

The recommended way is using `go get`, from your custom project run:
```bash
go get github.com/opendatahub-io/model-registry
```

## Getting Started

Model Registry Service is an high level Go client (or library) for ML Metadata (MLMD) store/service.
It provides model registry metadata capabilities, e.g., store and retrieve ML models metadata and related artifacts, through a custom defined [API](../pkg/api/api.go).

### Prerequisites

* MLMD server, check [ml-metadata doc](https://github.com/google/ml-metadata/blob/f0fef74eae2bdf6650a79ba976b36bea0b777c2e/g3doc/get_started.md#use-mlmd-with-a-remote-grpc-server) for more details on how to startup a MLMD store server.
* Go >= 1.19

### Usage

Assuming that MLMD server is already running at `localhost:9090`, as first step you should setup a gRPC connection to the server:

```go
import (
  "context"
  "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

conn, err := grpc.DialContext(
  context.Background(),
  "localhost:9090",
  grpc.WithReturnConnectionError(),
  grpc.WithBlock(), // optional
  grpc.WithTransportCredentials(insecure.NewCredentials()),
)
if err != nil {
  return fmt.Errorf("error dialing connection to mlmd server localhost:9090: %v", err)
}
defer conn.Close()
```

> NOTE: check [grpc doc](https://pkg.go.dev/google.golang.org/grpc#DialContext) for more details.

Once the gRPC connection is setup, let's create the `ModelRegistryService`:

```go
import (
  "fmt"
  "github.com/opendatahub-io/model-registry/pkg/core"
)

service, err := core.NewModelRegistryService(conn)
if err != nil {
  return fmt.Errorf("error creating model registry core service: %v", err)
}
```

Everything is ready, you can start using the `ModelRegistryService` library!

Here some usage examples:

#### Model Registration

Create a `RegisteredModel`

```go
modelName := "MODEL_NAME"
modelDescription := "MODEL_DESCRIPTION"

// register a new model
registeredModel, err = service.UpsertRegisteredModel(&openapi.RegisteredModel{
  Name:        &modelName,
  Description: &modelDescription,
})
if err != nil {
  return fmt.Errorf("error registering model: %v", err)
}
```

Create a new `ModelVersion` for the previous registered model

```go
versionName := "VERSION_NAME"
versionDescription := "VERSION_DESCRIPTION"
versionScore := 0.83

// register model version
modelVersion, err = service.UpsertModelVersion(&openapi.ModelVersion{
  Name:        &versionName,
  Description: &versionDescription,
  CustomProperties: &map[string]openapi.MetadataValue{
    "score": {
      MetadataDoubleValue: &openapi.MetadataDoubleValue{
        DoubleValue: &versionScore,
      },
    },
  },
}, registeredModel.Id)
if err != nil {
  return fmt.Errorf("error registering model version: %v", err)
}
```

Create a new `ModelArtifact` for the newly created version

```go
artifactName := "ARTIFACT_NAME"
artifactDescription := "ARTIFACT_DESCRIPTION"
artifactUri := "ARTIFACT_URI"

// register model artifact
modelArtifact, err := service.UpsertModelArtifact(&openapi.ModelArtifact{
  Name:        &artifactName,
  Description: &artifactDescription,
  Uri:         &artifactUri,
}, modelVersion.Id)
if err != nil {
  return fmt.Errorf("error creating model artifact: %v", err)
}
```

#### Model Query

Get `RegisteredModel` by name, for now the `name` must match.
```go
modelName := "QUERY_MODEL_NAME"
registeredModel, err := service.GetRegisteredModelByParams(&modelName, nil)
if err != nil {
  log.Printf("unable to find model %s: %v", getModelCfg.RegisteredModelName, err)
  return err
}
```

Get all `ModelVersion` associated to a specific registered model

```go
allVersions, err := service.GetModelVersions(api.ListOptions{}, registeredModel.Id)
if err != nil {
  return fmt.Errorf("error retrieving model versions for model %s: %v", *registeredModel.Id, err)
}
```
