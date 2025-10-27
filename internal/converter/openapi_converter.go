package converter

import (
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// NOTE: methods must follow these patterns, otherwise tests could not find possible issues:
// Converters createEntity to entity: Convert<ENTITY>Create
// Converters updateEntity to entity: Convert<ENTITY>Update
// Converters override fields entity: OverrideNotEditableFor<ENTITY>

type OpenAPIModelWrapper[
	M OpenAPIModel,
] struct {
	Model            *M
	ParentResourceId *string
	ModelName        *string
	TypeId           int32
}

// goverter:converter
// goverter:output:file ./generated/openapi_converter.gen.go
// goverter:wrapErrors
// goverter:enum:unknown @error
// goverter:matchIgnoreCase
// goverter:useZeroValueOnPointerInconsistency
type OpenAPIConverter interface {
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertRegisteredModelCreate(source *openapi.RegisteredModelCreate) (*openapi.RegisteredModel, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name
	ConvertRegisteredModelUpdate(source *openapi.RegisteredModelUpdate) (*openapi.RegisteredModel, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertModelVersionCreate(source *openapi.ModelVersionCreate) (*openapi.ModelVersion, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name RegisteredModelId
	ConvertModelVersionUpdate(source *openapi.ModelVersionUpdate) (*openapi.ModelVersion, error)

	// goverter:map DocArtifactCreate DocArtifact
	// goverter:map ModelArtifactCreate ModelArtifact
	// goverter:map DataSetCreate DataSet
	// goverter:map MetricCreate Metric
	// goverter:map ParameterCreate Parameter
	ConvertArtifactCreate(source *openapi.ArtifactCreate) (*openapi.Artifact, error)

	// goverter:map DocArtifactUpdate DocArtifact
	// goverter:map ModelArtifactUpdate ModelArtifact
	// goverter:map DataSetUpdate DataSet
	// goverter:map MetricUpdate Metric
	// goverter:map ParameterUpdate Parameter
	ConvertArtifactUpdate(source *openapi.ArtifactUpdate) (*openapi.Artifact, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch ArtifactType ExperimentId ExperimentRunId
	ConvertDocArtifactCreate(source *openapi.DocArtifactCreate) (*openapi.DocArtifact, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch ArtifactType Name ExperimentId ExperimentRunId
	ConvertDocArtifactUpdate(source *openapi.DocArtifactUpdate) (*openapi.DocArtifact, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch ArtifactType ExperimentId ExperimentRunId
	ConvertModelArtifactCreate(source *openapi.ModelArtifactCreate) (*openapi.ModelArtifact, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch ArtifactType Name ExperimentId ExperimentRunId
	ConvertModelArtifactUpdate(source *openapi.ModelArtifactUpdate) (*openapi.ModelArtifact, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch ArtifactType ExperimentId ExperimentRunId
	ConvertDataSetCreate(source *openapi.DataSetCreate) (*openapi.DataSet, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch ArtifactType Name ExperimentId ExperimentRunId
	ConvertDataSetUpdate(source *openapi.DataSetUpdate) (*openapi.DataSet, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch ArtifactType ExperimentId ExperimentRunId
	ConvertMetricCreate(source *openapi.MetricCreate) (*openapi.Metric, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch ArtifactType Name ExperimentId ExperimentRunId
	ConvertMetricUpdate(source *openapi.MetricUpdate) (*openapi.Metric, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch ArtifactType ExperimentId ExperimentRunId
	ConvertParameterCreate(source *openapi.ParameterCreate) (*openapi.Parameter, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch ArtifactType Name ExperimentId ExperimentRunId
	ConvertParameterUpdate(source *openapi.ParameterUpdate) (*openapi.Parameter, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertServingEnvironmentCreate(source *openapi.ServingEnvironmentCreate) (*openapi.ServingEnvironment, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name
	ConvertServingEnvironmentUpdate(source *openapi.ServingEnvironmentUpdate) (*openapi.ServingEnvironment, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertInferenceServiceCreate(source *openapi.InferenceServiceCreate) (*openapi.InferenceService, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name RegisteredModelId ServingEnvironmentId
	ConvertInferenceServiceUpdate(source *openapi.InferenceServiceUpdate) (*openapi.InferenceService, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertServeModelCreate(source *openapi.ServeModelCreate) (*openapi.ServeModel, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name ModelVersionId
	ConvertServeModelUpdate(source *openapi.ServeModelUpdate) (*openapi.ServeModel, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertExperimentCreate(source *openapi.ExperimentCreate) (*openapi.Experiment, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name
	ConvertExperimentUpdate(source *openapi.ExperimentUpdate) (*openapi.Experiment, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertExperimentRunCreate(source *openapi.ExperimentRunCreate) (*openapi.ExperimentRun, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name ExperimentId StartTimeSinceEpoch
	ConvertExperimentRunUpdate(source *openapi.ExperimentRunUpdate) (*openapi.ExperimentRun, error)

	// Ignore all fields that ARE editable
	// goverter:default InitWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalId CustomProperties State Owner Readme Maturity Language Tasks Provider Logo License LicenseLink LibraryName
	OverrideNotEditableForRegisteredModel(source OpenapiUpdateWrapper[openapi.RegisteredModel]) (openapi.RegisteredModel, error)

	// Ignore all fields that ARE editable
	// goverter:default InitWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalId CustomProperties State Author
	OverrideNotEditableForModelVersion(source OpenapiUpdateWrapper[openapi.ModelVersion]) (openapi.ModelVersion, error)

	// Ignore all fields that ARE editable
	// goverter:default InitWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore DocArtifact ModelArtifact DataSet Metric Parameter
	OverrideNotEditableForArtifact(source OpenapiUpdateWrapper[openapi.Artifact]) (openapi.Artifact, error)

	// Ignore all fields that ARE editable
	// goverter:default InitWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id Name ArtifactType CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalId CustomProperties Uri State
	OverrideNotEditableForDocArtifact(source OpenapiUpdateWrapper[openapi.DocArtifact]) (openapi.DocArtifact, error)

	// Ignore all fields that ARE editable
	// goverter:default InitWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalId CustomProperties Uri State ServiceAccountName ModelFormatName ModelFormatVersion StorageKey StoragePath ModelSourceKind ModelSourceClass ModelSourceGroup ModelSourceId ModelSourceName
	OverrideNotEditableForModelArtifact(source OpenapiUpdateWrapper[openapi.ModelArtifact]) (openapi.ModelArtifact, error)

	// Ignore all fields that ARE editable
	// goverter:default InitWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalId CustomProperties Uri State Digest SourceType Source Schema Profile
	OverrideNotEditableForDataSet(source OpenapiUpdateWrapper[openapi.DataSet]) (openapi.DataSet, error)

	// Ignore all fields that ARE editable
	// goverter:default InitWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalId CustomProperties State Value Timestamp Step
	OverrideNotEditableForMetric(source OpenapiUpdateWrapper[openapi.Metric]) (openapi.Metric, error)

	// Ignore all fields that ARE editable
	// goverter:default InitWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalId CustomProperties State Value ParameterType
	OverrideNotEditableForParameter(source OpenapiUpdateWrapper[openapi.Parameter]) (openapi.Parameter, error)

	// Ignore all fields that ARE editable
	// goverter:default InitWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalId CustomProperties
	OverrideNotEditableForServingEnvironment(source OpenapiUpdateWrapper[openapi.ServingEnvironment]) (openapi.ServingEnvironment, error)

	// Ignore all fields that ARE editable
	// goverter:default InitWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalId CustomProperties ModelVersionId Runtime DesiredState
	OverrideNotEditableForInferenceService(source OpenapiUpdateWrapper[openapi.InferenceService]) (openapi.InferenceService, error)

	// Ignore all fields that ARE editable
	// goverter:default InitWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalId CustomProperties LastKnownState
	OverrideNotEditableForServeModel(source OpenapiUpdateWrapper[openapi.ServeModel]) (openapi.ServeModel, error)

	// Ignore all fields that ARE editable for Experiment
	// goverter:default InitWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalId CustomProperties State Owner
	OverrideNotEditableForExperiment(source OpenapiUpdateWrapper[openapi.Experiment]) (openapi.Experiment, error)

	// Ignore all fields that ARE editable for ExperimentRun
	// goverter:default InitWithUpdate
	// goverter:autoMap Existing
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Description ExternalId CustomProperties State Owner Status StartTimeSinceEpoch EndTimeSinceEpoch
	OverrideNotEditableForExperimentRun(source OpenapiUpdateWrapper[openapi.ExperimentRun]) (openapi.ExperimentRun, error)
}
