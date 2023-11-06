package converter

import (
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
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
	ConvertRegisteredModel(source *proto.Context) (*openapi.RegisteredModel, error)

	// goverter:map Name | MapNameFromOwned
	// goverter:map Properties Description | MapDescription
	ConvertModelVersion(source *proto.Context) (*openapi.ModelVersion, error)

	// TODO: map actually ignored properties from Artifact.Properties
	// goverter:map Name | MapNameFromOwned
	// goverter:map . ArtifactType | MapArtifactType
	// goverter:map State | MapMLMDModelArtifactState
	// goverter:map Properties Description | MapDescription
	// goverter:map Properties Runtime | MapModelArtifactRuntime
	// goverter:map Properties ModelFormatName | MapModelArtifactFormatName
	// goverter:map Properties ModelFormatVersion | MapModelArtifactFormatVersion
	// goverter:map Properties StorageKey | MapModelArtifactStorageKey
	// goverter:map Properties StoragePath | MapModelArtifactStoragePath
	// goverter:map Properties ServiceAccountName | MapModelArtifactServiceAccountName
	ConvertModelArtifact(source *proto.Artifact) (*openapi.ModelArtifact, error)
}
