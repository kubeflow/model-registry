package converter

import (
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// goverter:converter
// goverter:output:file ./generated/mlmd_openapi_converter.gen.go
// goverter:wrapErrors
// goverter:matchIgnoreCase
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend Int64ToString
// goverter:extend StringToInt64
// goverter:extend MapMLMDCustomProperties
type MLMDToOpenAPIConverter interface {
	// goverter:map Properties Description | MapDescription
	// goverter:map Properties Owner | MapOwner
	// goverter:map Properties State | MapRegisteredModelState
	ConvertRegisteredModel(source *proto.Context) (*openapi.RegisteredModel, error)

	// goverter:map Name | MapName
	// goverter:map Name RegisteredModelId | MapRegisteredModelIdFromOwned
	// goverter:map Properties Description | MapDescription
	// goverter:map Properties State | MapModelVersionState
	// goverter:map Properties Author | MapPropertyAuthor
	ConvertModelVersion(source *proto.Context) (*openapi.ModelVersion, error)

	// goverter:map Name | MapNameFromOwned
	// goverter:map . ArtifactType | MapArtifactType
	// goverter:map State | MapMLMDArtifactState
	// goverter:map Properties Description | MapDescription
	// goverter:map Properties ModelFormatName | MapModelArtifactFormatName
	// goverter:map Properties ModelFormatVersion | MapModelArtifactFormatVersion
	// goverter:map Properties StorageKey | MapModelArtifactStorageKey
	// goverter:map Properties StoragePath | MapModelArtifactStoragePath
	// goverter:map Properties ServiceAccountName | MapModelArtifactServiceAccountName
	// goverter:map Properties ModelSourceKind | MapModelArtifactModelSourceKind
	// goverter:map Properties ModelSourceClass | MapModelArtifactModelSourceClass
	// goverter:map Properties ModelSourceGroup | MapModelArtifactModelSourceGroup
	// goverter:map Properties ModelSourceId | MapModelArtifactModelSourceId
	// goverter:map Properties ModelSourceName | MapModelArtifactModelSourceName
	ConvertModelArtifact(source *proto.Artifact) (*openapi.ModelArtifact, error)

	// goverter:map Name | MapNameFromOwned
	// goverter:map . ArtifactType | MapArtifactType
	// goverter:map State | MapMLMDArtifactState
	// goverter:map Properties Description | MapDescription
	ConvertDocArtifact(source *proto.Artifact) (*openapi.DocArtifact, error)

	// goverter:map Name | MapNameFromOwned
	// goverter:map Properties Description | MapDescription
	ConvertServingEnvironment(source *proto.Context) (*openapi.ServingEnvironment, error)

	// goverter:map Name | MapNameFromOwned
	// goverter:map Properties Description | MapDescription
	// goverter:map Properties Runtime | MapPropertyRuntime
	// goverter:map Properties ModelVersionId | MapPropertyModelVersionId
	// goverter:map Properties RegisteredModelId | MapPropertyRegisteredModelId
	// goverter:map Properties ServingEnvironmentId | MapPropertyServingEnvironmentId
	// goverter:map Properties DesiredState | MapInferenceServiceDesiredState
	ConvertInferenceService(source *proto.Context) (*openapi.InferenceService, error)

	// goverter:map Name | MapNameFromOwned
	// goverter:map Properties Description | MapDescription
	// goverter:map Properties ModelVersionId | MapPropertyModelVersionIdAsValue
	// goverter:map LastKnownState | MapMLMDServeModelLastKnownState
	ConvertServeModel(source *proto.Execution) (*openapi.ServeModel, error)
}
