package filter

import (
	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	"github.com/kubeflow/model-registry/internal/db/filter"
)

// CatalogRestEntityType represents catalog-specific REST API entity types
type CatalogRestEntityType string

const (
	RestEntityCatalogModel    CatalogRestEntityType = "CatalogModel"
	RestEntityCatalogArtifact CatalogRestEntityType = "CatalogArtifact"
	RestEntityMCPServer       CatalogRestEntityType = "MCPServer"
	RestEntityMCPServerTool   CatalogRestEntityType = "MCPServerTool"
)

// NewCatalogEntityMappings creates a new instance of catalog entity mappings
func NewCatalogEntityMappings() filter.EntityMappingFunctions {
	return defaultRegistry()
}

// defaultRegistry builds the catalog entity registry with all known entity types.
func defaultRegistry() *basecatalog.CatalogEntityRegistry {
	reg := basecatalog.NewCatalogEntityRegistry()

	reg.Register(filter.RestEntityType(RestEntityCatalogModel), basecatalog.EntityTypeDefinition{
		MLMDEntityType: filter.EntityTypeContext,
		Properties: basecatalog.MergeProperties(
			basecatalog.CommonContextProperties(),
			basecatalog.CommonCatalogContextProperties(),
			catalogModelSpecificProperties,
		),
		RelatedEntityPrefix: "artifacts.",
		RelatedEntityType:   filter.RelatedEntityArtifact,
	})

	reg.Register(filter.RestEntityType(RestEntityCatalogArtifact), basecatalog.EntityTypeDefinition{
		MLMDEntityType: filter.EntityTypeArtifact,
		Properties: basecatalog.MergeProperties(
			basecatalog.CommonArtifactProperties(),
			catalogArtifactSpecificProperties,
		),
	})

	reg.Register(filter.RestEntityType(RestEntityMCPServer), basecatalog.EntityTypeDefinition{
		MLMDEntityType: filter.EntityTypeContext,
		Properties: basecatalog.MergeProperties(
			basecatalog.CommonContextProperties(),
			basecatalog.CommonCatalogContextProperties(),
			mcpServerSpecificProperties,
		),
	})

	reg.Register(filter.RestEntityType(RestEntityMCPServerTool), basecatalog.EntityTypeDefinition{
		MLMDEntityType: filter.EntityTypeArtifact,
		Properties: basecatalog.MergeProperties(
			mcpServerToolEntityProperties,
			mcpServerToolSpecificProperties,
		),
	})

	return reg
}

// --- Entity-specific property maps (only properties NOT in the common sets) ---

// catalogModelSpecificProperties are CatalogModel properties beyond the common context
// and common catalog context sets.
var catalogModelSpecificProperties = map[string]filter.PropertyDefinition{
	"owner":        {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "owner"},
	"state":        {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "state"},
	"language":     {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "language"},
	"library_name": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "library_name"},
	"maturity":     {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "maturity"},
	"tasks":        {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "tasks"},
}

// catalogArtifactSpecificProperties are CatalogArtifact properties beyond the common
// artifact set.
var catalogArtifactSpecificProperties = map[string]filter.PropertyDefinition{
	"artifactType": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "artifactType"},
}

// mcpServerSpecificProperties are MCPServer properties beyond the common context
// and common catalog context sets.
var mcpServerSpecificProperties = map[string]filter.PropertyDefinition{
	// Core
	"base_name": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "base_name"},
	"version":   {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "version"},

	// Filterable dimensions (arrays, bools)
	"tags":           {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "tags"},
	"transports":     {Location: filter.PropertyTable, ValueType: filter.ArrayValueType, Column: "transports"},
	"deploymentMode": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "deploymentMode"},

	// URL fields
	"documentationUrl": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "documentationUrl"},
	"repositoryUrl":    {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "repositoryUrl"},
	"sourceCode":       {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "sourceCode"},

	// Time fields
	"publishedDate": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "publishedDate"},
	"lastUpdated":   {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "lastUpdated"},

	// Security indicators
	"verifiedSource": {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "verifiedSource"},
	"secureEndpoint": {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "secureEndpoint"},
	"sast":           {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "sast"},
	"readOnlyTools":  {Location: filter.PropertyTable, ValueType: filter.BoolValueType, Column: "readOnlyTools"},

	// Complex objects stored as JSON
	"endpoints":       {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "endpoints"},
	"artifacts":       {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "artifacts"},
	"runtimeMetadata": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "runtimeMetadata"},
}

// mcpServerToolEntityProperties are the entity-table properties for MCPServerTool.
// MCPServerTool does not include externalId, uri, or state from the common artifact set.
var mcpServerToolEntityProperties = map[string]filter.PropertyDefinition{
	"id":                       {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "id"},
	"name":                     {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "name"},
	"createTimeSinceEpoch":     {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "create_time_since_epoch"},
	"lastUpdateTimeSinceEpoch": {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "last_update_time_since_epoch"},
}

// mcpServerToolSpecificProperties are MCPServerTool-specific property-table properties.
var mcpServerToolSpecificProperties = map[string]filter.PropertyDefinition{
	"description": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "description"},
	"accessType":  {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "accessType"},
}
