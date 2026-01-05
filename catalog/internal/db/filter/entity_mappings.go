package filter

import (
	"strings"

	"github.com/kubeflow/model-registry/internal/db/filter"
)

// CatalogRestEntityType represents catalog-specific REST API entity types
type CatalogRestEntityType string

const (
	RestEntityCatalogModel    CatalogRestEntityType = "CatalogModel"
	RestEntityCatalogArtifact CatalogRestEntityType = "CatalogArtifact"
	RestEntityMcpServer       CatalogRestEntityType = "McpServer"
	RestEntityMcpServerTool   CatalogRestEntityType = "McpServerTool"
)

// catalogEntityMappings implements EntityMappingFunctions for the catalog package
type catalogEntityMappings struct{}

// NewCatalogEntityMappings creates a new instance of catalog entity mappings
func NewCatalogEntityMappings() filter.EntityMappingFunctions {
	return &catalogEntityMappings{}
}

// GetMLMDEntityType maps catalog REST entity types to their underlying MLMD entity type
func (c *catalogEntityMappings) GetMLMDEntityType(restEntityType filter.RestEntityType) filter.EntityType {
	switch restEntityType {
	case filter.RestEntityType(RestEntityCatalogArtifact):
		return filter.EntityTypeArtifact
	default:
		return filter.EntityTypeContext
	}
}

// GetPropertyDefinitionForRestEntity returns property definition for a catalog REST entity type
func (c *catalogEntityMappings) GetPropertyDefinitionForRestEntity(restEntityType filter.RestEntityType, propertyName string) filter.PropertyDefinition {
	// Check if this is a well-known property for catalog entities
	if restEntityType == filter.RestEntityType(RestEntityCatalogModel) {
		if _, isWellKnown := catalogModelProperties[propertyName]; isWellKnown {
			// Use the well-known property definition
			return catalogModelProperties[propertyName]
		}

		// Check if this is a property path referencing a related artifact
		// Format: artifacts.<propertyName> or artifacts.customProperties.<propertyName>
		if strings.HasPrefix(propertyName, "artifacts.") {
			// Extract the artifact property path (everything after "artifacts.")
			artifactPropertyPath := strings.TrimPrefix(propertyName, "artifacts.")

			// Return a RelatedEntity property definition
			// ValueType is left empty to allow runtime type inference from the value
			return filter.PropertyDefinition{
				Location:          filter.RelatedEntity,
				ValueType:         "", // Empty to enable runtime type inference
				Column:            artifactPropertyPath,
				RelatedEntityType: filter.RelatedEntityArtifact,
				RelatedProperty:   artifactPropertyPath,
				JoinTable:         "Attribution", // Join through Attribution table
			}
		}
	}

	if restEntityType == filter.RestEntityType(RestEntityCatalogArtifact) {
		if _, isWellKnown := catalogArtifactProperties[propertyName]; isWellKnown {
			// Use the well-known property definition
			return catalogArtifactProperties[propertyName]
		}
	}

	if restEntityType == filter.RestEntityType(RestEntityMcpServer) {
		if _, isWellKnown := mcpServerProperties[propertyName]; isWellKnown {
			// Use the well-known property definition
			return mcpServerProperties[propertyName]
		}
	}

	if restEntityType == filter.RestEntityType(RestEntityMcpServerTool) {
		if _, isWellKnown := mcpServerToolProperties[propertyName]; isWellKnown {
			// Use the well-known property definition
			return mcpServerToolProperties[propertyName]
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
	"language":     {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "language"},
	"library_name": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "library_name"},
	"license_link": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "license_link"},
	"license":      {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "license"},
	"logo":         {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "logo"},
	"maturity":     {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "maturity"},
	"provider":     {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "provider"},
	"readme":       {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "readme"},
	"tasks":        {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "tasks"},
}

// catalogArtifactProperties defines the allowed properties for CatalogArtifact entities
var catalogArtifactProperties = map[string]filter.PropertyDefinition{
	// Common Artifact properties
	"id":                       {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "id"},
	"name":                     {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "name"},
	"externalId":               {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "external_id"},
	"createTimeSinceEpoch":     {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "last_update_time_since_epoch"},
	"uri":                      {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "uri"},
	"state":                    {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "state"},

	// Artifact type (stored in type_id but we can filter by string representation)
	"artifactType": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "artifactType"},
}

// mcpServerProperties defines the allowed properties for McpServer entities.
// This follows the same pattern as catalogModelProperties - only properties that are:
// 1. Entity table columns (required for core identity)
// 2. Key filterable dimensions that need explicit type handling (arrays, bools)
// 3. Common properties shared with CatalogModel for consistency
// All other properties can be queried via custom property fallback.
var mcpServerProperties = map[string]filter.PropertyDefinition{
	// Common Context properties (Entity Table - required)
	"id":                       {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "id"},
	"name":                     {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "name"},
	"externalId":               {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "external_id"},
	"createTimeSinceEpoch":     {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "last_update_time_since_epoch"},

	// Core properties matching CatalogModel pattern
	"source_id":    {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "source_id"},
	"description":  {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "description"},
	"provider":     {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "provider"},
	"license":      {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "license"},
	"license_link": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "license_link"},
	"logo":         {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "logo"},
	"readme":       {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "readme"},

	// MCP-specific filterable dimensions (need explicit type for arrays)
	"tags":           {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "tags"},
	"transports":     {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "transports"},
	"deploymentMode": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "deploymentMode"},

	// Security indicators (need explicit type for booleans)
	"verifiedSource":  {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "verifiedSource"},
	"secureEndpoint":  {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "secureEndpoint"},
	"sast":            {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "sast"},
	"readOnlyTools":   {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "readOnlyTools"},
}

// mcpServerToolProperties defines the allowed properties for McpServerTool entities
var mcpServerToolProperties = map[string]filter.PropertyDefinition{
	// Common Artifact properties
	"id":                       {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "id"},
	"name":                     {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "name"},
	"createTimeSinceEpoch":     {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "last_update_time_since_epoch"},

	// McpServerTool-specific properties stored in ArtifactProperty table
	"description": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "description"},
	"accessType":  {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "accessType"},
}
