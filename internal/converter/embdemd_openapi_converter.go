package converter

import (
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// goverter:converter
// goverter:output:file ./generated/embedmd_openapi_converter.gen.go
// goverter:wrapErrors
// goverter:matchIgnoreCase
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend Int64ToString
// goverter:extend Int32ToString
// goverter:extend StringToInt64
// goverter:extend MapEmbedMDCustomProperties
type EmbedMDToOpenAPIConverter interface {
	// goverter:map Properties Description | MapEmbedMDDescription
	// goverter:map Properties Owner | MapEmbedMDOwner
	// goverter:map Properties Language | MapEmbedMDPropertyLanguage
	// goverter:map Properties LibraryName | MapEmbedMDPropertyLibraryName
	// goverter:map Properties LicenseLink | MapEmbedMDPropertyLicenseLink
	// goverter:map Properties License | MapEmbedMDPropertyLicense
	// goverter:map Properties Logo | MapEmbedMDPropertyLogo
	// goverter:map Properties Maturity | MapEmbedMDPropertyMaturity
	// goverter:map Properties Provider | MapEmbedMDPropertyProvider
	// goverter:map Properties Readme | MapEmbedMDPropertyReadme
	// goverter:map Properties Tasks | MapEmbedMDPropertyTasks
	// goverter:map Properties State | MapEmbedMDStateRegisteredModel
	// goverter:map Attributes ExternalId | MapEmbedMDExternalIDRegisteredModel
	// goverter:map Attributes Name | MapEmbedMDNameRegisteredModel
	// goverter:map Attributes CreateTimeSinceEpoch | MapEmbedMDCreateTimeSinceEpochRegisteredModel
	// goverter:map Attributes LastUpdateTimeSinceEpoch | MapEmbedMDLastUpdateTimeSinceEpochRegisteredModel
	ConvertRegisteredModel(source *models.RegisteredModelImpl) (*openapi.RegisteredModel, error)

	// goverter:map Properties Description | MapEmbedMDDescription
	// goverter:map Properties Author | MapEmbedMDAuthor
	// goverter:map Properties State | MapEmbedMDStateModelVersion
	// goverter:map Properties RegisteredModelId | MapEmbedMDPropertyRegisteredModelId
	// goverter:map Attributes ExternalId | MapEmbedMDExternalIDModelVersion
	// goverter:map Attributes Name | MapEmbedMDNameModelVersion
	// goverter:map Attributes CreateTimeSinceEpoch | MapEmbedMDCreateTimeSinceEpochModelVersion
	// goverter:map Attributes LastUpdateTimeSinceEpoch | MapEmbedMDLastUpdateTimeSinceEpochModelVersion
	ConvertModelVersion(source *models.ModelVersionImpl) (*openapi.ModelVersion, error)

	// goverter:map Properties Description | MapEmbedMDDescription
	// goverter:map Properties ModelFormatName | MapEmbedMDPropertyModelFormatName
	// goverter:map Properties ModelFormatVersion | MapEmbedMDPropertyModelFormatVersion
	// goverter:map Properties StorageKey | MapEmbedMDPropertyStorageKey
	// goverter:map Properties StoragePath | MapEmbedMDPropertyStoragePath
	// goverter:map Properties ServiceAccountName | MapEmbedMDPropertyServiceAccountName
	// goverter:map Properties ModelSourceKind | MapEmbedMDPropertyModelSourceKind
	// goverter:map Properties ModelSourceClass | MapEmbedMDPropertyModelSourceClass
	// goverter:map Properties ModelSourceGroup | MapEmbedMDPropertyModelSourceGroup
	// goverter:map Properties ModelSourceId | MapEmbedMDPropertyModelSourceId
	// goverter:map Properties ModelSourceName | MapEmbedMDPropertyModelSourceName
	// goverter:map Attributes ExternalId | MapEmbedMDExternalIDModelArtifact
	// goverter:map Attributes Name | MapEmbedMDNameModelArtifact
	// goverter:map Attributes Uri | MapEmbedMDURIModelArtifact
	// goverter:map Attributes State | MapEmbedMDStateModelArtifact
	// goverter:map Attributes ArtifactType | MapEmbedMDArtifactTypeModelArtifact
	// goverter:map Attributes CreateTimeSinceEpoch | MapEmbedMDCreateTimeSinceEpochModelArtifact
	// goverter:map Attributes LastUpdateTimeSinceEpoch | MapEmbedMDLastUpdateTimeSinceEpochModelArtifact
	ConvertModelArtifact(source *models.ModelArtifactImpl) (*openapi.ModelArtifact, error)

	// goverter:map Properties Description | MapEmbedMDDescription
	// goverter:map Attributes ExternalId | MapEmbedMDExternalIDDocArtifact
	// goverter:map Attributes Name | MapEmbedMDNameDocArtifact
	// goverter:map Attributes Uri | MapEmbedMDURIDocArtifact
	// goverter:map Attributes State | MapEmbedMDStateDocArtifact
	// goverter:map Attributes ArtifactType | MapEmbedMDArtifactTypeDocArtifact
	// goverter:map Attributes CreateTimeSinceEpoch | MapEmbedMDCreateTimeSinceEpochDocArtifact
	// goverter:map Attributes LastUpdateTimeSinceEpoch | MapEmbedMDLastUpdateTimeSinceEpochDocArtifact
	ConvertDocArtifact(source *models.DocArtifactImpl) (*openapi.DocArtifact, error)

	// goverter:map Properties Description | MapEmbedMDDescription
	// goverter:map Attributes ExternalId | MapEmbedMDExternalIDServingEnvironment
	// goverter:map Attributes Name | MapEmbedMDNameServingEnvironment
	// goverter:map Attributes CreateTimeSinceEpoch | MapEmbedMDCreateTimeSinceEpochServingEnvironment
	// goverter:map Attributes LastUpdateTimeSinceEpoch | MapEmbedMDLastUpdateTimeSinceEpochServingEnvironment
	ConvertServingEnvironment(source *models.ServingEnvironmentImpl) (*openapi.ServingEnvironment, error)

	// goverter:map Properties Description | MapEmbedMDDescription
	// goverter:map Properties Runtime | MapEmbedMDPropertyRuntime
	// goverter:map Properties DesiredState | MapEmbedMDPropertyDesiredStateInferenceService
	// goverter:map Properties ModelVersionId | MapEmbedMDPropertyModelVersionId
	// goverter:map Properties RegisteredModelId | MapEmbedMDPropertyRegisteredModelId
	// goverter:map Properties ServingEnvironmentId | MapEmbedMDPropertyServingEnvironmentId
	// goverter:map Attributes ExternalId | MapEmbedMDExternalIDInferenceService
	// goverter:map Attributes Name | MapEmbedMDNameInferenceService
	// goverter:map Attributes CreateTimeSinceEpoch | MapEmbedMDCreateTimeSinceEpochInferenceService
	// goverter:map Attributes LastUpdateTimeSinceEpoch | MapEmbedMDLastUpdateTimeSinceEpochInferenceService
	ConvertInferenceService(source *models.InferenceServiceImpl) (*openapi.InferenceService, error)

	// goverter:map Properties Description | MapEmbedMDDescription
	// goverter:map Properties ModelVersionId | MapEmbedMDPropertyModelVersionIdServeModel
	// goverter:map Attributes ExternalId | MapEmbedMDExternalIDServeModel
	// goverter:map Attributes Name | MapEmbedMDNameServeModel
	// goverter:map Attributes LastKnownState | MapEmbedMDLastKnownStateServeModel
	// goverter:map Attributes CreateTimeSinceEpoch | MapEmbedMDCreateTimeSinceEpochServeModel
	// goverter:map Attributes LastUpdateTimeSinceEpoch | MapEmbedMDLastUpdateTimeSinceEpochServeModel
	ConvertServeModel(source *models.ServeModelImpl) (*openapi.ServeModel, error)

	// goverter:map Properties Description | MapEmbedMDDescription
	// goverter:map Properties Owner | MapEmbedMDOwner
	// goverter:map Properties State | MapEmbedMDStateExperiment
	// goverter:map Attributes ExternalId | MapEmbedMDExternalIDExperiment
	// goverter:map Attributes Name | MapEmbedMDNameExperiment
	// goverter:map Attributes CreateTimeSinceEpoch | MapEmbedMDCreateTimeSinceEpochExperiment
	// goverter:map Attributes LastUpdateTimeSinceEpoch | MapEmbedMDLastUpdateTimeSinceEpochExperiment
	ConvertExperiment(source *models.ExperimentImpl) (*openapi.Experiment, error)

	// goverter:map Properties Description | MapEmbedMDDescription
	// goverter:map Properties Owner | MapEmbedMDOwner
	// goverter:map Properties State | MapEmbedMDStateExperimentRun
	// goverter:map Properties Status | MapEmbedMDPropertyStatusExperimentRun
	// goverter:map Properties StartTimeSinceEpoch | MapEmbedMDPropertyStartTimeSinceEpochExperimentRun
	// goverter:map Properties EndTimeSinceEpoch | MapEmbedMDPropertyEndTimeSinceEpochExperimentRun
	// goverter:map Properties ExperimentId | MapEmbedMDPropertyExperimentIdExperimentRun
	// goverter:map Attributes ExternalId | MapEmbedMDExternalIDExperimentRun
	// goverter:map Attributes Name | MapEmbedMDNameExperimentRun
	// goverter:map Attributes CreateTimeSinceEpoch | MapEmbedMDCreateTimeSinceEpochExperimentRun
	// goverter:map Attributes LastUpdateTimeSinceEpoch | MapEmbedMDLastUpdateTimeSinceEpochExperimentRun
	ConvertExperimentRun(source *models.ExperimentRunImpl) (*openapi.ExperimentRun, error)

	// goverter:map Properties Description | MapEmbedMDDescription
	// goverter:map Properties Digest | MapEmbedMDPropertyDigest
	// goverter:map Properties SourceType | MapEmbedMDPropertySourceType
	// goverter:map Properties Source | MapEmbedMDPropertySource
	// goverter:map Properties Schema | MapEmbedMDPropertySchema
	// goverter:map Properties Profile | MapEmbedMDPropertyProfile
	// goverter:map Attributes ExternalId | MapEmbedMDExternalIDDataSet
	// goverter:map Attributes Name | MapEmbedMDNameDataSet
	// goverter:map Attributes Uri | MapEmbedMDURIDataSet
	// goverter:map Attributes State | MapEmbedMDStateDataSet
	// goverter:map Attributes ArtifactType | MapEmbedMDArtifactTypeDataSet
	// goverter:map Attributes CreateTimeSinceEpoch | MapEmbedMDCreateTimeSinceEpochDataSet
	// goverter:map Attributes LastUpdateTimeSinceEpoch | MapEmbedMDLastUpdateTimeSinceEpochDataSet
	ConvertDataSet(source *models.DataSetImpl) (*openapi.DataSet, error)

	// goverter:map Properties Description | MapEmbedMDDescription
	// goverter:map Properties Value | MapEmbedMDPropertyValueMetric
	// goverter:map Properties Timestamp | MapEmbedMDPropertyTimestampMetric
	// goverter:map Properties Step | MapEmbedMDPropertyStepMetric
	// goverter:map Attributes ExternalId | MapEmbedMDExternalIDMetric
	// goverter:map Attributes Name | MapEmbedMDNameMetric
	// goverter:map Attributes State | MapEmbedMDStateMetric
	// goverter:map Attributes ArtifactType | MapEmbedMDArtifactTypeMetric
	// goverter:map Attributes CreateTimeSinceEpoch | MapEmbedMDCreateTimeSinceEpochMetric
	// goverter:map Attributes LastUpdateTimeSinceEpoch | MapEmbedMDLastUpdateTimeSinceEpochMetric
	ConvertMetric(source *models.MetricImpl) (*openapi.Metric, error)

	// goverter:map Properties Description | MapEmbedMDDescription
	// goverter:map Properties Value | MapEmbedMDPropertyValueParameter
	// goverter:map Properties ParameterType | MapEmbedMDPropertyParameterTypeParameter
	// goverter:map Attributes ExternalId | MapEmbedMDExternalIDParameter
	// goverter:map Attributes Name | MapEmbedMDNameParameter
	// goverter:map Attributes State | MapEmbedMDStateParameter
	// goverter:map Attributes ArtifactType | MapEmbedMDArtifactTypeParameter
	// goverter:map Attributes CreateTimeSinceEpoch | MapEmbedMDCreateTimeSinceEpochParameter
	// goverter:map Attributes LastUpdateTimeSinceEpoch | MapEmbedMDLastUpdateTimeSinceEpochParameter
	ConvertParameter(source *models.ParameterImpl) (*openapi.Parameter, error)
}
