package filter

import (
	platformfilter "github.com/kubeflow/hub/internal/platform/db/filter"
)

// Entity type constants and types
type EntityType = platformfilter.EntityType

const (
	EntityTypeContext   = platformfilter.EntityTypeContext
	EntityTypeArtifact  = platformfilter.EntityTypeArtifact
	EntityTypeExecution = platformfilter.EntityTypeExecution
)

// Interfaces
type EntityMappingFunctions = platformfilter.EntityMappingFunctions
type EqualityExpander = platformfilter.EqualityExpander

// QueryBuilder type
type QueryBuilder = platformfilter.QueryBuilder

// NewQueryBuilderForRestEntity creates a new query builder for the specified REST entity type.
// If mappingFuncs is nil, it falls back to the model-registry default entity mappings.
func NewQueryBuilderForRestEntity(restEntityType RestEntityType, mappingFuncs EntityMappingFunctions) *QueryBuilder {
	if mappingFuncs == nil {
		mappingFuncs = &defaultEntityMappings{}
	}
	return platformfilter.NewQueryBuilderForRestEntity(restEntityType, mappingFuncs)
}

// DefaultEntityMappingFuncs returns the default MR entity mapping functions.
func DefaultEntityMappingFuncs() EntityMappingFunctions {
	return &defaultEntityMappings{}
}

// defaultEntityMappings implements EntityMappingFunctions for the model registry
type defaultEntityMappings struct{}

func (d *defaultEntityMappings) GetMLMDEntityType(restEntityType RestEntityType) EntityType {
	return GetMLMDEntityType(restEntityType)
}

func (d *defaultEntityMappings) GetPropertyDefinitionForRestEntity(restEntityType RestEntityType, propertyName string) PropertyDefinition {
	return GetPropertyDefinitionForRestEntity(restEntityType, propertyName)
}

func (d *defaultEntityMappings) IsChildEntity(entityType RestEntityType) bool {
	return isChildEntity(entityType)
}
