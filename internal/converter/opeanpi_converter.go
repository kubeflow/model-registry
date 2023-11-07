package converter

import "github.com/opendatahub-io/model-registry/internal/model/openapi"

// goverter:converter
// goverter:output:file ./generated/openapi_converter.gen.go
// goverter:wrapErrors
// goverter:matchIgnoreCase
// goverter:useZeroValueOnPointerInconsistency
type OpenAPIConverter interface {
	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertRegisteredModelCreate(source *openapi.RegisteredModelCreate) (*openapi.RegisteredModel, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name
	ConvertRegisteredModelUpdate(source *openapi.RegisteredModelUpdate) (*openapi.RegisteredModel, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertModelVersionCreate(source *openapi.ModelVersionCreate) (*openapi.ModelVersion, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch Name
	ConvertModelVersionUpdate(source *openapi.ModelVersionUpdate) (*openapi.ModelVersion, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch ArtifactType
	ConvertModelArtifactCreate(source *openapi.ModelArtifactCreate) (*openapi.ModelArtifact, error)

	// goverter:ignore Id CreateTimeSinceEpoch LastUpdateTimeSinceEpoch ArtifactType Name
	ConvertModelArtifactUpdate(source *openapi.ModelArtifactUpdate) (*openapi.ModelArtifact, error)
}
