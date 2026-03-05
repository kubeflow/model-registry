package basecatalog

import (
	"maps"

	"github.com/kubeflow/model-registry/internal/db/filter"
)

// CommonContextProperties returns the entity-table properties shared by all Context entities.
func CommonContextProperties() map[string]filter.PropertyDefinition {
	return map[string]filter.PropertyDefinition{
		"id":                       {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "id"},
		"name":                     {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "name"},
		"externalId":               {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "external_id"},
		"createTimeSinceEpoch":     {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "create_time_since_epoch"},
		"lastUpdateTimeSinceEpoch": {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "last_update_time_since_epoch"},
	}
}

// CommonArtifactProperties returns the entity-table properties shared by all Artifact entities.
func CommonArtifactProperties() map[string]filter.PropertyDefinition {
	return map[string]filter.PropertyDefinition{
		"id":                       {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "id"},
		"name":                     {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "name"},
		"externalId":               {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "external_id"},
		"createTimeSinceEpoch":     {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "create_time_since_epoch"},
		"lastUpdateTimeSinceEpoch": {Location: filter.EntityTable, ValueType: filter.IntValueType, Column: "last_update_time_since_epoch"},
		"uri":                      {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "uri"},
		"state":                    {Location: filter.EntityTable, ValueType: filter.StringValueType, Column: "state"},
	}
}

// CommonCatalogContextProperties returns the property-table properties shared across catalog Context entities
// (e.g. CatalogModel, MCPServer).
func CommonCatalogContextProperties() map[string]filter.PropertyDefinition {
	return map[string]filter.PropertyDefinition{
		"source_id":    {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "source_id"},
		"description":  {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "description"},
		"provider":     {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "provider"},
		"license":      {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "license"},
		"license_link": {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "license_link"},
		"logo":         {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "logo"},
		"readme":       {Location: filter.PropertyTable, ValueType: filter.StringValueType, Column: "readme"},
	}
}

// MergeProperties merges multiple property maps. Later maps win on conflicts.
func MergeProperties(propMaps ...map[string]filter.PropertyDefinition) map[string]filter.PropertyDefinition {
	result := make(map[string]filter.PropertyDefinition)
	for _, m := range propMaps {
		maps.Copy(result, m)
	}
	return result
}
