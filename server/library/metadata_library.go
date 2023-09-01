package library

import "github.com/dhirajsb/ml-metadata-go-server/ml_metadata/proto"

type MetadataType struct {
	Name        *string                       `yaml:"name,omitempty"`
	Version     *string                       `yaml:"version,omitempty"`
	Description *string                       `yaml:"description,omitempty"`
	ExternalId  *string                       `yaml:"external_id,omitempty"`
	Properties  map[string]proto.PropertyType `yaml:"properties,omitempty"`
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
	ArtifactTypes  []ArtifactType `yaml:"artifact-types,omitempty"`
	ContextTypes   []ArtifactType `yaml:"context-types,omitempty"`
	ExecutionTypes []ArtifactType `yaml:"execution-types,omitempty"`
}
