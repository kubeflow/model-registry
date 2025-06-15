package converter

import "github.com/kubeflow/model-registry/pkg/openapi"

// NOTE: methods must follow these patterns, otherwise tests could not find possible issues:
// Converters patch fields entity: UpdateExisting<ENTITY>

// goverter:converter
// goverter:output:file ./generated/openapi_reconciler.gen.go
// goverter:wrapErrors
// goverter:enum:unknown @error
// goverter:matchIgnoreCase
// goverter:useZeroValueOnPointerInconsistency
type OpenAPIReconciler interface {
	// Ignore all fields that can't be updated
	// goverter:default InitWithExisting
	// goverter:autoMap Update
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name
	UpdateExistingRegisteredModel(source OpenapiUpdateWrapper[openapi.RegisteredModel]) (openapi.RegisteredModel, error)

	// Ignore all fields that can't be updated
	// goverter:default InitWithExisting
	// goverter:autoMap Update
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name RegisteredModelId
	UpdateExistingModelVersion(source OpenapiUpdateWrapper[openapi.ModelVersion]) (openapi.ModelVersion, error)

	// Ignore all fields that can't be updated
	// goverter:default InitWithExisting
	// goverter:autoMap Update
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name ArtifactType
	UpdateExistingDocArtifact(source OpenapiUpdateWrapper[openapi.DocArtifact]) (openapi.DocArtifact, error)

	// Ignore all fields that can't be updated
	// goverter:default InitWithExisting
	// goverter:autoMap Update
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name ArtifactType
	UpdateExistingModelArtifact(source OpenapiUpdateWrapper[openapi.ModelArtifact]) (openapi.ModelArtifact, error)

	// Ignore all fields that can't be updated
	// goverter:default InitWithExisting
	// goverter:autoMap Update
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name ArtifactType
	UpdateExistingDataSet(source OpenapiUpdateWrapper[openapi.DataSet]) (openapi.DataSet, error)

	// Ignore all fields that can't be updated
	// goverter:default InitWithExisting
	// goverter:autoMap Update
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name ArtifactType
	UpdateExistingMetric(source OpenapiUpdateWrapper[openapi.Metric]) (openapi.Metric, error)

	// Ignore all fields that can't be updated
	// goverter:default InitWithExisting
	// goverter:autoMap Update
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name ArtifactType
	UpdateExistingParameter(source OpenapiUpdateWrapper[openapi.Parameter]) (openapi.Parameter, error)

	// Ignore all fields that can't be updated
	// goverter:default InitWithExisting
	// goverter:autoMap Update
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name
	UpdateExistingServingEnvironment(source OpenapiUpdateWrapper[openapi.ServingEnvironment]) (openapi.ServingEnvironment, error)

	// Ignore all fields that can't be updated
	// goverter:default InitWithExisting
	// goverter:autoMap Update
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name RegisteredModelId ServingEnvironmentId
	UpdateExistingInferenceService(source OpenapiUpdateWrapper[openapi.InferenceService]) (openapi.InferenceService, error)

	// Ignore all fields that can't be updated
	// goverter:default InitWithExisting
	// goverter:autoMap Update
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name ModelVersionId
	UpdateExistingServeModel(source OpenapiUpdateWrapper[openapi.ServeModel]) (openapi.ServeModel, error)

	// Ignore all fields that can't be updated
	// goverter:default InitWithExisting
	// goverter:autoMap Update
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name
	UpdateExistingExperiment(source OpenapiUpdateWrapper[openapi.Experiment]) (openapi.Experiment, error)

	// Ignore all fields that can't be updated
	// goverter:default InitWithExisting
	// goverter:autoMap Update
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name ExperimentId
	UpdateExistingExperimentRun(source OpenapiUpdateWrapper[openapi.ExperimentRun]) (openapi.ExperimentRun, error)
}
