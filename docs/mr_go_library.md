# Model Registry

> ⚠️  NOTE: UNSTABLE API  ⚠️
> This library is provided as a convenience for Kubeflow Model Registry developers.
> If you are not actively involved in the development of Model Registry, please prefer the [REST API](https://editor.swagger.io/?url=https://raw.githubusercontent.com/kubeflow/model-registry/main/api/openapi/model-registry.yaml).

## Getting Started

Model Registry is a high level Go library client for a remote [ML Metadata (MLMD)](https://github.com/google/ml-metadata)
server that provides a [high-level model metadata registry API](../pkg/api/api.go).

You can use Model Registry to interact with MLMD using [convenient type definitions](../pkg/openapi/), instead of managing them yourself.
This includes models, model versions, artifacts, serving environments, inference services, and serve models.

### Prerequisites

* [MLMD server](https://github.com/google/ml-metadata/blob/f0fef74eae2bdf6650a79ba976b36bea0b777c2e/g3doc/get_started.md#use-mlmd-with-a-remote-grpc-server)
* Go >= 1.24

Install it using:

```bash
go get github.com/kubeflow/model-registry
```

Assuming that an MLMD server is running at `localhost:9090`, you can connect with it using a gRPC connection:

<!-- TODO: https://github.com/kubeflow/model-registry/issues/194: drop DialContext from this example -->

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

> NOTE: Check the [go gRPC documentation](https://pkg.go.dev/google.golang.org/grpc#DialContext) for more details.

Once the gRPC connection is established, you can create a `ModelRegistryService`:

```go
import (
  "fmt"
  "github.com/kubeflow/model-registry/pkg/core"
)

service, err := core.NewModelRegistryService(conn)
if err != nil {
  return fmt.Errorf("error creating model registry core service: %v", err)
}
```

### Example usage

After setting up your Model Registry Service, you can try some of the following:

#### Creating models

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

Create a `ModelVersion` for the registered model

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

Create a `ModelArtifact` for the version

```go
artifactName := "ARTIFACT_NAME"
artifactDescription := "ARTIFACT_DESCRIPTION"
artifactUri := "ARTIFACT_URI"

// register model artifact
modelArtifact, err := service.UpsertModelVersionArtifact(&openapi.ModelArtifact{
  Name:        &artifactName,
  Description: &artifactDescription,
  Uri:         &artifactUri,
}, *modelVersion.Id)
if err != nil {
  return fmt.Errorf("error creating model artifact: %v", err)
}
```

#### Updating models

The Model Registry Service provides `Upsert*` methods for all supported models, meaning that you can use them to
in**sert** a new model, or **up**date it:

```go
artifactName := "ARTIFACT_NAME"
artifactDescription := "ARTIFACT_DESCRIPTION"
artifactUri := "ARTIFACT_URI"

// register model artifact
modelArtifact, err := service.UpsertModelVersionArtifact(&openapi.ModelArtifact{
  Name:        &artifactName,
  Description: &artifactDescription,
  Uri:         &artifactUri,
}, *modelVersion.Id)
if err != nil {
  return fmt.Errorf("error creating model artifact: %v", err)
}

newDescription := "update it!"
modelArtifact.Description = &newDescription

modelArtifact, err = service.UpsertModelVersionArtifact(modelArtifact, *modelVersion.Id)
if err != nil {
  return fmt.Errorf("error updating model artifact: %v", err)
}
```


#### Querying models

Get a `RegisteredModel` by name:

```go
modelName := "MODEL_NAME"
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

## API documentation

Check out the [Model Registry Service interface](../pkg/api/api.go) and the [core layer implementation](../pkg/core/) for additional details.
