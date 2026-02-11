package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// filterValueType maps a catalog.yaml property type to a filter value type constant.
func filterValueType(specType string) string {
	switch specType {
	case "string":
		return "string_value"
	case "integer", "int":
		return "int_value"
	case "int64":
		return "int_value"
	case "boolean", "bool":
		return "int_value" // bools stored as int (0/1)
	case "number", "float", "double":
		return "double_value"
	default:
		return "string_value"
	}
}

// generateFilterMappings generates the filter entity mappings file.
func generateFilterMappings(config CatalogConfig) error {
	entityName := config.Spec.Entity.Name

	// Built-in fields that are already handled in the template
	builtinFields := map[string]bool{
		"name": true, "externalid": true, "createtimesinceepoch": true,
		"lastupdatetimesinceepoch": true, "id": true,
	}

	// Build property registrations for init() function
	var propRegistrations strings.Builder
	for _, prop := range config.Spec.Entity.Properties {
		if builtinFields[strings.ToLower(prop.Name)] {
			continue
		}
		propRegistrations.WriteString(fmt.Sprintf("\t\t\"%s\": true,\n", prop.Name))
	}

	// Build property definitions for GetPropertyDefinitionForRestEntity switch cases
	var propDefinitions strings.Builder
	for _, prop := range config.Spec.Entity.Properties {
		if builtinFields[strings.ToLower(prop.Name)] {
			continue
		}
		valueType := filterValueType(prop.Type)
		propDefinitions.WriteString(fmt.Sprintf("\tcase \"%s\":\n", prop.Name))
		propDefinitions.WriteString(fmt.Sprintf("\t\treturn filter.PropertyDefinition{Location: filter.PropertyTable, ValueType: \"%s\", Column: \"%s\"}\n", valueType, prop.Name))
	}

	data := map[string]any{
		"EntityName":             entityName,
		"Package":                config.Spec.Package,
		"PropertyRegistrations": propRegistrations.String(),
		"PropertyDefinitions":   propDefinitions.String(),
	}

	serviceDir := filepath.Join("internal", "db", "service")
	if err := ensureDir(serviceDir); err != nil {
		return err
	}

	outputPath := filepath.Join(serviceDir, "filter_mappings.go")
	fmt.Printf("  Generated: internal/db/service/filter_mappings.go\n")
	return executeTemplate(TmplServiceFilterMappings, outputPath, data)
}
