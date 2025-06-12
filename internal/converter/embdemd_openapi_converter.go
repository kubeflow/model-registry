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
}
