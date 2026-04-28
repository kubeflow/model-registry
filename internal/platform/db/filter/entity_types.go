package filter

// RestEntityType represents the specific REST API entity type
type RestEntityType string

// EntityType represents the type of entity for proper query building
type EntityType string

const (
	EntityTypeContext   EntityType = "context"
	EntityTypeArtifact  EntityType = "artifact"
	EntityTypeExecution EntityType = "execution"
)

// EntityMappingFunctions defines the interface for entity type mapping functions.
// This allows different packages (like catalog) to provide their own entity mappings.
type EntityMappingFunctions interface {
	GetMLMDEntityType(restEntityType RestEntityType) EntityType
	GetPropertyDefinitionForRestEntity(restEntityType RestEntityType, propertyName string) PropertyDefinition
	IsChildEntity(entityType RestEntityType) bool
}

// EqualityExpander is an optional interface that EntityMappingFunctions may implement
// to expand equality conditions.
type EqualityExpander interface {
	GetEqualityExpansion(restEntityType RestEntityType, propertyName string, value any) (likeArg any, useExpansion bool)
}
