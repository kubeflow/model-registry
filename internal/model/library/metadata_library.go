package library

import (
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
)

//go:generate go-enum -type=PropertyType

type PropertyType int32

const (
	UNKNOWN PropertyType = iota
	INT
	DOUBLE
	STRING
	STRUCT
	PROTO
	BOOLEAN
)

type MetadataType struct {
	Name        *string                 `yaml:"name,omitempty"`
	Version     *string                 `yaml:"version,omitempty"`
	Description *string                 `yaml:"description,omitempty"`
	ExternalId  *string                 `yaml:"external_id,omitempty"`
	Properties  map[string]PropertyType `yaml:"properties,omitempty"`
}

type ArtifactType struct {
	MetadataType `yaml:",inline"`
	// TODO add support for base type enum
	//BaseType *ArtifactType_SystemDefinedBaseType `yaml:"base_type,omitempty"`
}

type ContextType struct {
	MetadataType `yaml:",inline"`
}

type ExecutionType struct {
	MetadataType `yaml:",inline"`
	//InputType  *ArtifactStructType                  `yaml:"input_type,omitempty"`
	//OutputType *ArtifactStructType                  `yaml:"output_type,omitempty"`
	//BaseType   *ExecutionType_SystemDefinedBaseType `yaml:"base_type,omitempty"`
}

type MetadataLibrary struct {
	ArtifactTypes  []ArtifactType  `yaml:"artifact-types,omitempty"`
	ContextTypes   []ContextType   `yaml:"context-types,omitempty"`
	ExecutionTypes []ExecutionType `yaml:"execution-types,omitempty"`
}

func ToProtoProperties(props map[string]PropertyType) map[string]proto.PropertyType {
	result := make(map[string]proto.PropertyType)
	for name, prop := range props {
		result[name] = proto.PropertyType(prop)
	}
	return result
}
