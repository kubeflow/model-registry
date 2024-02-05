package converter

import (
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/pkg/openapi"
)

type OpenAPIModelWrapper[
	M openapi.RegisteredModel |
		openapi.ModelVersion |
		openapi.ModelArtifact |
		openapi.ServingEnvironment |
		openapi.InferenceService |
		openapi.ServeModel,
] struct {
	TypeId           int64
	Model            *M
	ParentResourceId *string // optional parent id
	ModelName        *string // optional registered model name
}

// goverter:converter
// goverter:output:file ./generated/openapi_mlmd_converter.gen.go
// goverter:wrapErrors
// goverter:matchIgnoreCase
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend Int64ToString
// goverter:extend StringToInt64
// goverter:extend MapOpenAPICustomProperties
type OpenAPIToMLMDConverter interface {
	// goverter:autoMap Model
	// goverter:map Model Type | MapRegisteredModelType
	// goverter:map Model Properties | MapRegisteredModelProperties
	// goverter:ignore state sizeCache unknownFields SystemMetadata CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertRegisteredModel(source *OpenAPIModelWrapper[openapi.RegisteredModel]) (*proto.Context, error)

	// goverter:autoMap Model
	// goverter:map . Name | MapModelVersionName
	// goverter:map Model Type | MapModelVersionType
	// goverter:map . Properties | MapModelVersionProperties
	// goverter:ignore state sizeCache unknownFields SystemMetadata CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertModelVersion(source *OpenAPIModelWrapper[openapi.ModelVersion]) (*proto.Context, error)

	// goverter:autoMap Model
	// goverter:map . Name | MapModelArtifactName
	// goverter:map Model Type | MapModelArtifactType
	// goverter:map Model Properties | MapModelArtifactProperties
	// goverter:map Model.State State | MapOpenAPIArtifactState
	// goverter:ignore state sizeCache unknownFields SystemMetadata CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertModelArtifact(source *OpenAPIModelWrapper[openapi.ModelArtifact]) (*proto.Artifact, error)

	// goverter:autoMap Model
	// goverter:map Model Type | MapServingEnvironmentType
	// goverter:map Model Properties | MapServingEnvironmentProperties
	// goverter:ignore state sizeCache unknownFields SystemMetadata CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertServingEnvironment(source *OpenAPIModelWrapper[openapi.ServingEnvironment]) (*proto.Context, error)

	// goverter:autoMap Model
	// goverter:map . Name | MapInferenceServiceName
	// goverter:map Model Type | MapInferenceServiceType
	// goverter:map Model Properties | MapInferenceServiceProperties
	// goverter:ignore state sizeCache unknownFields SystemMetadata CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertInferenceService(source *OpenAPIModelWrapper[openapi.InferenceService]) (*proto.Context, error)

	// goverter:autoMap Model
	// goverter:map . Name | MapServeModelName
	// goverter:map Model Type | MapServeModelType
	// goverter:map Model Properties | MapServeModelProperties
	// goverter:map Model.LastKnownState LastKnownState | MapLastKnownState
	// goverter:ignore state sizeCache unknownFields SystemMetadata CreateTimeSinceEpoch LastUpdateTimeSinceEpoch
	ConvertServeModel(source *OpenAPIModelWrapper[openapi.ServeModel]) (*proto.Execution, error)
}
