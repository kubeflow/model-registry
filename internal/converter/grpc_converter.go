package converter

import (
	"fmt"

	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/db"
	"gorm.io/gorm"
)

var (
	// singleton DB connection for type name lookup
	globalDB *gorm.DB
)

// SetConverterDB must be called before using gRPC converters,
// it uses the singleton DB connection to lookup type names
func SetConverterDB(db *gorm.DB) error {
	if globalDB != nil {
		return fmt.Errorf("converter global DB connection MUST only be set once")
	}
	globalDB = db
	return initTypeNameCache()
}

// goverter:converter
// goverter:output:file ./generated/grpc_converter.gen.go
// goverter:wrapErrors
// goverter:matchIgnoreCase
// goverter:useZeroValueOnPointerInconsistency
type GRPCConverter interface {
	// goverter:map State | ConvertArtifact_State
	// goverter:map . Properties | ConvertProtoArtifactProperties
	// goverter:ignore ArtifactType Attributions Events
	ConvertArtifact(source *proto.Artifact) (*db.Artifact, error)
	// goverter:map State | ConvertToArtifact_State
	// goverter:map ID Id
	// goverter:map ID Type | ConvertTypeIDToName
	// goverter:map Properties Properties | ConvertToProtoArtifactProperties
	// goverter:map Properties CustomProperties | ConvertToProtoArtifactCustomProperties
	// goverter:ignore state sizeCache unknownFields SystemMetadata
	ConvertToArtifact(source *db.Artifact) (*proto.Artifact, error)

	// goverter:map . Properties | ConvertProtoContextProperties
	// goverter:ignore ContextType Attributions Associations Parents Children
	ConvertContext(source *proto.Context) (*db.Context, error)
	// goverter:map ID Id
	// goverter:map ID Type | ConvertTypeIDToName
	// goverter:map Properties | ConvertToProtoContextProperties
	// goverter:map Properties CustomProperties | ConvertToProtoContextCustomProperties
	// goverter:ignore state sizeCache unknownFields SystemMetadata
	ConvertToContext(source *db.Context) (*proto.Context, error)

	// goverter:map LastKnownState | ConvertExecution_State
	// goverter:map . Properties | ConvertProtoExecutionProperties
	// goverter:ignore ExecutionType Associations Events
	ConvertExecution(source *proto.Execution) (*db.Execution, error)
	// goverter:map ID Id
	// goverter:map ID Type | ConvertTypeIDToName
	// goverter:map LastKnownState | ConvertToExecution_State
	// goverter:map Properties | ConvertToProtoExecutionProperties
	// goverter:map Properties CustomProperties | ConvertToProtoExecutionCustomProperties
	// goverter:ignore state sizeCache unknownFields SystemMetadata
	ConvertToExecution(source *db.Execution) (*proto.Execution, error)

	// goverter:map Type | ConvertProtoEventType
	// goverter:map Path PathSteps | ConvertProtoEventPath
	// goverter:ignore ID Artifact Execution
	ConvertEvent(source *proto.Event) (*db.Event, error)
	// goverter:map Type | ConvertToProtoEventType
	// goverter:map PathSteps Path | ConvertToProtoEventPath
	// goverter:ignore state sizeCache unknownFields SystemMetadata
	ConvertToEvent(source *db.Event) (*proto.Event, error)

	// goverter:ignore ID Context Artifact
	ConvertAttribution(source *proto.Attribution) (*db.Attribution, error)
	// goverter:ignore ID state sizeCache unknownFields
	ConvertToAttribution(source *db.Attribution) (*proto.Attribution, error)

	// goverter:ignore ID Context Execution
	ConvertAssociation(source *proto.Association) (*db.Association, error)
	// goverter:ignore state sizeCache unknownFields
	ConvertToAssociation(source *db.Association) (*proto.Association, error)
}
