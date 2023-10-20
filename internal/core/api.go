package core

import "github.com/opendatahub-io/model-registry/internal/model/openapi"

// Note: for convention, here we are keeping here adherence to the mlmd side
type BaseResourceId int64

type ListOptions struct {
	PageSize      *int32
	OrderBy       *string
	SortOrder     *string
	NextPageToken *string
}

type ModelRegistryApi interface {
	// REGISTERED MODEL

	// UpsertRegisteredModel create or update a registered model, the behavior follows the same
	// approach used by MLMD gRPC api. If Id is provided update the entity otherwise create a new one.
	UpsertRegisteredModel(registeredModel *openapi.RegisteredModel) (*openapi.RegisteredModel, error)

	GetRegisteredModelById(id *BaseResourceId) (*openapi.RegisteredModel, error)
	GetRegisteredModelByParams(name *string, externalId *string) (*openapi.RegisteredModel, error)
	GetRegisteredModels(listOptions ListOptions) (*openapi.RegisteredModelList, error)

	// MODEL VERSION

	// Create a new Model Version
	// or update a Model Version associated to a specific RegisteredModel identified by registeredModelId parameter
	UpsertModelVersion(modelVersion *openapi.ModelVersion, registeredModelId *BaseResourceId) (*openapi.ModelVersion, error)

	GetModelVersionById(id *BaseResourceId) (*openapi.ModelVersion, error)
	GetModelVersionByParams(name *string, externalId *string) (*openapi.ModelVersion, error)
	GetModelVersions(listOptions ListOptions, registeredModelId *BaseResourceId) (*openapi.ModelVersionList, error)

	// MODEL ARTIFACT

	// Create or update a Model Artifact associated to a specific ModelVersion
	// identified by ModelArtifact.ModelVersionId
	UpsertModelArtifact(modelArtifact *openapi.ModelArtifact, modelVersionId *BaseResourceId) (*openapi.ModelArtifact, error)

	GetModelArtifactById(id *BaseResourceId) (*openapi.ModelArtifact, error)
	GetModelArtifactByParams(name *string, externalId *string) (*openapi.ModelArtifact, error)
	GetModelArtifacts(listOptions ListOptions, modelVersionId *BaseResourceId) (*openapi.ModelArtifactList, error)
}
