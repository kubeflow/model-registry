package filter

import (
	"github.com/kubeflow/model-registry/internal/db/filter"
)

// CatalogRestEntityType represents catalog-specific REST API entity types
type CatalogRestEntityType string

const (
	RestEntityCatalogModel CatalogRestEntityType = "CatalogModel"
)

// catalogEntityMappings implements EntityMappingFunctions for the catalog package
type catalogEntityMappings struct{}

// NewCatalogEntityMappings creates a new instance of catalog entity mappings
func NewCatalogEntityMappings() filter.EntityMappingFunctions {
	return &catalogEntityMappings{}
}

// GetMLMDEntityType maps catalog REST entity types to their underlying MLMD entity type
func (c *catalogEntityMappings) GetMLMDEntityType(restEntityType filter.RestEntityType) filter.EntityType {
	return filter.EntityTypeContext
}

// GetPropertyDefinitionForRestEntity returns property definition for a catalog REST entity type
func (c *catalogEntityMappings) GetPropertyDefinitionForRestEntity(restEntityType filter.RestEntityType, propertyName string) filter.PropertyDefinition {
	// Check if this is a well-known property for catalog entities
	if restEntityType == filter.RestEntityType(RestEntityCatalogModel) {
		if _, isWellKnown := catalogModelProperties[propertyName]; isWellKnown {
			// Use the well-known property definition
			return catalogModelProperties[propertyName]
		}
	}

	// Not a well-known property for this entity type, treat as custom
	return filter.PropertyDefinition{
		Location:  filter.Custom,
		ValueType: filter.StringValueType, // Default, will be inferred at runtime
		Column:    propertyName,           // Use the property name as-is for custom properties
	}
}

// IsChildEntity returns true if the catalog REST entity type uses prefixed names (parentId:name)
func (c *catalogEntityMappings) IsChildEntity(entityType filter.RestEntityType) bool {
	return false
}

// catalogModelProperties defines the allowed properties for CatalogModel entities
var catalogModelProperties = map[string]filter.PropertyDefinition{
	// Common Context properties
	"id":                       {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "id"},
	"name":                     {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "name"},
	"externalId":               {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "external_id"},
	"createTimeSinceEpoch":     {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "last_update_time_since_epoch"},

	// CatalogModel-specific properties stored in ContextProperty table
	"source_id":    {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "source_id"},
	"description":  {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "description"},
	"owner":        {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "owner"},
	"state":        {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "state"},
	"language":     {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "language"},
	"library_name": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "library_name"},
	"license_link": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "license_link"},
	"license":      {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "license"},
	"logo":         {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "logo"},
	"maturity":     {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "maturity"},
	"provider":     {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "provider"},
	"readme":       {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "readme"},
	"tasks":        {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "tasks"},
}
