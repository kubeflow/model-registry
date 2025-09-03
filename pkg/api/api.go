package api

import "github.com/kubeflow/model-registry/pkg/openapi"

// ListOptions provides options for listing entities with pagination and sorting.
// It includes parameters such as PageSize, OrderBy, SortOrder, and NextPageToken.
type ListOptions struct {
	PageSize      *int32  // The maximum number of entities to be returned per page.
	OrderBy       *string // The field by which entities are ordered.
	SortOrder     *string // The sorting order, which can be "ASC" (ascending) or "DESC" (descending).
	NextPageToken *string // A token to retrieve the next page of entities in a paginated result set.
	FilterQuery   *string // A filter query to restrict results based on entity properties.
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

	// ARTIFACT

	// UpsertModelVersionArtifact create or update an Artifact for a specific ModelVersion, the behavior follows the same
	// approach used by MLMD gRPC api. If Id is provided update the entity otherwise create a new one.
	UpsertModelVersionArtifact(artifact *openapi.Artifact, modelVersionId string) (*openapi.Artifact, error)

	// UpsertArtifact create or update an Artifact, the behavior follows the same
	// approach used by MLMD gRPC api. If Id is provided update the entity otherwise create a new one.
	UpsertArtifact(artifact *openapi.Artifact) (*openapi.Artifact, error)

	// GetArtifactById retrieve Artifact by id
	GetArtifactById(id string) (*openapi.Artifact, error)

	// GetArtifactByParams find Artifact instances that match the provided optional params
	GetArtifactByParams(artifactName *string, parentResourceId *string, externalId *string) (*openapi.Artifact, error)

	// GetArtifacts return all Artifact properly ordered and sized based on listOptions param.
	// if parentResourceId is provided, return all Artifact instances belonging to a specific parent resource
	GetArtifacts(artifactType openapi.ArtifactTypeQueryParam, listOptions ListOptions, parentResourceId *string) (*openapi.ArtifactList, error)

	// MODEL ARTIFACT

	// UpsertModelArtifact creates or inserts an Artifact
	UpsertModelArtifact(modelArtifact *openapi.ModelArtifact) (*openapi.ModelArtifact, error)

	// GetModelArtifactById retrieve ModelArtifact by id
	GetModelArtifactById(id string) (*openapi.ModelArtifact, error)

	// GetModelArtifactByInferenceService retrieve a ModelArtifact by inference service id
	GetModelArtifactByInferenceService(inferenceServiceId string) (*openapi.ModelArtifact, error)

	// GetModelArtifactByParams find ModelArtifact instances that match the provided optional params
	GetModelArtifactByParams(artifactName *string, parentResourceId *string, externalId *string) (*openapi.ModelArtifact, error)

	// GetModelArtifacts return all ModelArtifact properly ordered and sized based on listOptions param.
	// if parentResourceId is provided, return all ModelArtifact instances belonging to a specific parent resource
	GetModelArtifacts(listOptions ListOptions, parentResourceId *string) (*openapi.ModelArtifactList, error)

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
	// if runtime is provided, filter those InferenceService having that runtime
	GetInferenceServices(listOptions ListOptions, servingEnvironmentId *string, runtime *string) (*openapi.InferenceServiceList, error)

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

	// EXPERIMENT
	// UpsertExperiment create or update an experiment, the behavior follows the same
	// approach used by MLMD gRPC api. If Id is provided update the entity otherwise create a new one.
	UpsertExperiment(experiment *openapi.Experiment) (*openapi.Experiment, error)
	// GetExperimentById retrieve Experiment by id
	GetExperimentById(id string) (*openapi.Experiment, error)
	// GetExperimentByParams find Experiment instances that match the provided optional params
	GetExperimentByParams(name *string, externalId *string) (*openapi.Experiment, error)
	// GetExperiments return all Experiment properly ordered and sized based on listOptions param
	GetExperiments(listOptions ListOptions) (*openapi.ExperimentList, error)

	// EXPERIMENT RUN
	// UpsertExperimentRun create or update an experiment run, the behavior follows the same
	// approach used by MLMD gRPC api. If Id is provided update the entity otherwise create a new one.
	// experimentId defines the Experiment to be associated as parent ownership to the newly created ExperimentRun.
	UpsertExperimentRun(experimentRun *openapi.ExperimentRun, experimentId *string) (*openapi.ExperimentRun, error)
	// GetExperimentRunById retrieve ExperimentRun by id
	GetExperimentRunById(id string) (*openapi.ExperimentRun, error)
	// GetExperimentRunByParams find ExperimentRun instances that match the provided optional params
	GetExperimentRunByParams(name *string, experimentId *string, externalId *string) (*openapi.ExperimentRun, error)
	// GetExperimentRuns return all ExperimentRun properly ordered and sized based on listOptions param.
	// if experimentId is provided, return all ExperimentRun instances belonging to a specific Experiment
	GetExperimentRuns(listOptions ListOptions, experimentId *string) (*openapi.ExperimentRunList, error)

	// EXPERIMENT RUN ARTIFACTS
	// UpsertExperimentRunArtifact create or update an Artifact for a specific ExperimentRun, the behavior follows the same
	// approach used by MLMD gRPC api. If Id is provided update the entity otherwise create a new one.
	UpsertExperimentRunArtifact(artifact *openapi.Artifact, experimentRunId string) (*openapi.Artifact, error)
	// GetExperimentRunArtifacts return all Artifact properly ordered and sized based on listOptions param.
	// if experimentRunId is provided, return all Artifact instances belonging to a specific ExperimentRun
	GetExperimentRunArtifacts(artifactType openapi.ArtifactTypeQueryParam, listOptions ListOptions, experimentRunId *string) (*openapi.ArtifactList, error)

	// EXPERIMENT RUN METRIC HISTORY
	// GetExperimentRunMetricHistory return metric history for a specific ExperimentRun properly ordered and sized based on listOptions param.
	// if name is provided, filter metrics by name. if stepIds is provided, filter metrics by step ids
	GetExperimentRunMetricHistory(name *string, stepIds *string, listOptions ListOptions, experimentRunId *string) (*openapi.MetricList, error)
}
