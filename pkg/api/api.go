package api

import "github.com/opendatahub-io/model-registry/internal/model/openapi"

type ListOptions struct {
	PageSize      *int32
	OrderBy       *string
	SortOrder     *string
	NextPageToken *string
}

// ModelRegistryApi defines the external API for the Model Registry library
type ModelRegistryApi interface {
	// REGISTERED MODEL

	// UpsertRegisteredModel create or update a registered model, the behavior follows the same
	// approach used by MLMD gRPC api. If Id is provided update the entity otherwise create a new one.
	UpsertRegisteredModel(registeredModel *openapi.RegisteredModel) (*openapi.RegisteredModel, error)

	// GetRegisteredModelById retrieve RegisteredModel by id
	GetRegisteredModelById(id string) (*openapi.RegisteredModel, error)

	// GetRegisteredModelByInferenceService retrieve a RegisteredModel by inference service id
	GetRegisteredModelByInferenceService(inferenceServiceId string) (*openapi.RegisteredModel, error)

	// GetRegisteredModelByParams find RegisteredModel instances that match the provided optional params
	GetRegisteredModelByParams(name *string, externalId *string) (*openapi.RegisteredModel, error)

	// GetRegisteredModels return all ModelArtifact properly ordered and sized based on listOptions param.
	GetRegisteredModels(listOptions ListOptions) (*openapi.RegisteredModelList, error)

	// MODEL VERSION

	// UpsertModelVersion create a new Model Version or update a Model Version associated to a
	// specific RegisteredModel identified by registeredModelId parameter
	UpsertModelVersion(modelVersion *openapi.ModelVersion, registeredModelId *string) (*openapi.ModelVersion, error)

	// GetModelVersionById retrieve ModelVersion by id
	GetModelVersionById(id string) (*openapi.ModelVersion, error)

	// GetModelVersionByInferenceService retrieve a ModelVersion by inference service id
	GetModelVersionByInferenceService(inferenceServiceId string) (*openapi.ModelVersion, error)

	// GetModelVersionByParams find ModelVersion instances that match the provided optional params
	GetModelVersionByParams(versionName *string, registeredModelId *string, externalId *string) (*openapi.ModelVersion, error)

	// GetModelVersions return all ModelArtifact properly ordered and sized based on listOptions param.
	// if registeredModelId is provided, return all ModelVersion instances belonging to a specific RegisteredModel
	GetModelVersions(listOptions ListOptions, registeredModelId *string) (*openapi.ModelVersionList, error)

	// MODEL ARTIFACT

	// UpsertModelArtifact create a new Artifact or update an Artifact associated to a specific
	// ModelVersion identified by modelVersionId parameter
	UpsertModelArtifact(modelArtifact *openapi.ModelArtifact, modelVersionId *string) (*openapi.ModelArtifact, error)

	// GetModelArtifactById retrieve ModelArtifact by id
	GetModelArtifactById(id string) (*openapi.ModelArtifact, error)

	// GetModelArtifactByParams find ModelArtifact instances that match the provided optional params
	GetModelArtifactByParams(artifactName *string, modelVersionId *string, externalId *string) (*openapi.ModelArtifact, error)

	// GetModelArtifacts return all ModelArtifact properly ordered and sized based on listOptions param.
	// if modelVersionId is provided, return all ModelArtifact instances belonging to a specific ModelVersion
	GetModelArtifacts(listOptions ListOptions, modelVersionId *string) (*openapi.ModelArtifactList, error)

	// SERVING ENVIRONMENT

	// UpsertServingEnvironment create or update a serving environmet, the behavior follows the same
	// approach used by MLMD gRPC api. If Id is provided update the entity otherwise create a new one.
	UpsertServingEnvironment(registeredModel *openapi.ServingEnvironment) (*openapi.ServingEnvironment, error)

	// GetInferenceServiceById retrieve ServingEnvironment by id
	GetServingEnvironmentById(id string) (*openapi.ServingEnvironment, error)

	// GetServingEnvironmentByParams find ServingEnvironment instances that match the provided optional params
	GetServingEnvironmentByParams(name *string, externalId *string) (*openapi.ServingEnvironment, error)

	// GetServingEnvironments return all ServingEnvironment properly ordered and sized based on listOptions param
	GetServingEnvironments(listOptions ListOptions) (*openapi.ServingEnvironmentList, error)

	// INFERENCE SERVICE

	// UpsertInferenceService create or update an inference service, the behavior follows the same
	// approach used by MLMD gRPC api. If Id is provided update the entity otherwise create a new one.
	// inferenceService.servingEnvironmentId defines the ServingEnvironment to be associated as parent ownership
	// to the newly created InferenceService.
	UpsertInferenceService(inferenceService *openapi.InferenceService) (*openapi.InferenceService, error)

	// GetInferenceServiceById retrieve InferenceService by id
	GetInferenceServiceById(id string) (*openapi.InferenceService, error)

	// GetInferenceServiceByParams find InferenceService instances that match the provided optional params
	GetInferenceServiceByParams(name *string, parentResourceId *string, externalId *string) (*openapi.InferenceService, error)

	// GetInferenceServices return all InferenceService properly ordered and sized based on listOptions param
	// if servingEnvironmentId is provided, return all InferenceService instances belonging to a specific ServingEnvironment
	GetInferenceServices(listOptions ListOptions, servingEnvironmentId *string) (*openapi.InferenceServiceList, error)

	// SERVE MODEL

	// UpsertServeModel create or update a serve model, the behavior follows the same
	// approach used by MLMD gRPC api. If Id is provided update the entity otherwise create a new one.
	// inferenceServiceId defines the InferenceService to be linked to the newly created ServeModel.
	UpsertServeModel(serveModel *openapi.ServeModel, inferenceServiceId *string) (*openapi.ServeModel, error)

	// GetServeModelById retrieve ServeModel by id
	GetServeModelById(id string) (*openapi.ServeModel, error)

	// GetServeModels get all ServeModel objects properly ordered and sized based on listOptions param.
	// if inferenceServiceId is provided, return all ServeModel instances belonging to a specific InferenceService
	GetServeModels(listOptions ListOptions, inferenceServiceId *string) (*openapi.ServeModelList, error)
}
