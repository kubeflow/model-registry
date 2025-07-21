// Package apiutils provides query translation for REST API filter queries to MLMD filter queries.
//
// MLMD filter queries use SQL-like syntax but with MLMD-specific extensions:
// - Triple-dot notation: custom_properties.field.value_type
// - MLMD aliases: contexts_a.id, parent_contexts_a.name
// - Property access: properties.field.value_type
//
// Since this isn't pure SQL, we use a custom parser optimized for MLMD's hybrid syntax.
package apiutils

import (
	"fmt"
	"regexp"
	"strings"
)

// EntityType represents the type of entity for field mapping
type EntityType string

const (
	RegisteredModelEntity EntityType = "registered_model"
	ModelVersionEntity    EntityType = "model_version"
	ExperimentEntity      EntityType = "experiment"
	ExperimentRunEntity   EntityType = "experiment_run"
)

// FieldType represents the type of field (core MLMD field, property, or custom property)
type FieldType int

const (
	MLMDField      FieldType = iota // Core MLMD fields like id, name, external_id
	Property                        // Properties explicitly defined in REST API schemas
	CustomProperty                  // User-defined custom properties
)

// FieldMapping represents a mapping from REST API field to MLMD field
type FieldMapping struct {
	MLMDField    string
	FieldType    FieldType
	PropertyType string // For properties: "string_value", "int_value", "double_value", "bool_value"
}

// getFieldMappings returns the field mappings for each entity type
func getFieldMappings() map[EntityType]map[string]FieldMapping {
	return map[EntityType]map[string]FieldMapping{
		RegisteredModelEntity: {
			// Core MLMD fields
			"id":                       {MLMDField: "id", FieldType: MLMDField},
			"name":                     {MLMDField: "name", FieldType: MLMDField},
			"externalId":               {MLMDField: "external_id", FieldType: MLMDField},
			"createTimeSinceEpoch":     {MLMDField: "create_time_since_epoch", FieldType: MLMDField},
			"lastUpdateTimeSinceEpoch": {MLMDField: "last_update_time_since_epoch", FieldType: MLMDField},

			// Properties defined in RegisteredModel schema
			"description": {MLMDField: "properties.description.string_value", FieldType: Property},
			"state":       {MLMDField: "properties.state.string_value", FieldType: Property},
			"owner":       {MLMDField: "properties.owner.string_value", FieldType: Property},
		},
		ModelVersionEntity: {
			// Core MLMD fields
			"id":                       {MLMDField: "id", FieldType: MLMDField},
			"name":                     {MLMDField: "name", FieldType: MLMDField},
			"externalId":               {MLMDField: "external_id", FieldType: MLMDField},
			"createTimeSinceEpoch":     {MLMDField: "create_time_since_epoch", FieldType: MLMDField},
			"lastUpdateTimeSinceEpoch": {MLMDField: "last_update_time_since_epoch", FieldType: MLMDField},

			// Properties defined in ModelVersion schema
			"description":       {MLMDField: "properties.description.string_value", FieldType: Property},
			"state":             {MLMDField: "properties.state.string_value", FieldType: Property},
			"author":            {MLMDField: "properties.author.string_value", FieldType: Property},
			"registeredModelId": {MLMDField: "properties.registered_model_id.string_value", FieldType: Property},
		},
		ExperimentEntity: {
			// Core MLMD fields
			"id":                       {MLMDField: "id", FieldType: MLMDField},
			"name":                     {MLMDField: "name", FieldType: MLMDField},
			"externalId":               {MLMDField: "external_id", FieldType: MLMDField},
			"createTimeSinceEpoch":     {MLMDField: "create_time_since_epoch", FieldType: MLMDField},
			"lastUpdateTimeSinceEpoch": {MLMDField: "last_update_time_since_epoch", FieldType: MLMDField},

			// Properties defined in Experiment schema
			"description": {MLMDField: "properties.description.string_value", FieldType: Property},
			"state":       {MLMDField: "properties.state.string_value", FieldType: Property},
		},
		ExperimentRunEntity: {
			// Core MLMD fields
			"id":                       {MLMDField: "id", FieldType: MLMDField},
			"name":                     {MLMDField: "name", FieldType: MLMDField},
			"externalId":               {MLMDField: "external_id", FieldType: MLMDField},
			"createTimeSinceEpoch":     {MLMDField: "create_time_since_epoch", FieldType: MLMDField},
			"lastUpdateTimeSinceEpoch": {MLMDField: "last_update_time_since_epoch", FieldType: MLMDField},

			// Properties defined in ExperimentRun schema
			"description":  {MLMDField: "properties.description.string_value", FieldType: Property},
			"state":        {MLMDField: "properties.state.string_value", FieldType: Property},
			"experimentId": {MLMDField: "properties.experiment_id.string_value", FieldType: Property},
		},
	}
}

// isValidValueType checks if the provided value type is valid for MLMD filter queries
func isValidValueType(valueType string) bool {
	validTypes := map[string]bool{
		"string_value": true,
		"int_value":    true,
		"double_value": true,
		"bool_value":   true,
	}
	return validTypes[valueType]
}

// unescapeDots converts escaped dots (\.) back to regular dots (.) in property names
func unescapeDots(s string) string {
	return strings.ReplaceAll(s, "\\.", ".")
}

// findLastUnescapedDot finds the position of the last unescaped dot in a string
// Returns -1 if no unescaped dot is found
func findLastUnescapedDot(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			// Check if this dot is escaped by counting preceding backslashes
			backslashCount := 0
			for j := i - 1; j >= 0 && s[j] == '\\'; j-- {
				backslashCount++
			}
			// If even number of backslashes (including 0), the dot is not escaped
			if backslashCount%2 == 0 {
				return i
			}
		}
	}
	return -1
}

// parseCustomPropertyField parses a field name that may contain escaped dots
// Returns (propertyName, valueType, isCustomProperty)
// For "model\.version.string_value" returns ("model.version", "string_value", true)
// For "name" returns ("", "", false)
func parseCustomPropertyField(fieldName string) (string, string, bool) {
	// Find the last unescaped dot
	dotPos := findLastUnescapedDot(fieldName)
	if dotPos == -1 {
		// No unescaped dot found - this is not a custom property with value type
		return "", "", false
	}

	propertyName := fieldName[:dotPos]
	valueType := fieldName[dotPos+1:]

	// Validate the value type
	if !isValidValueType(valueType) {
		return "", "", false
	}

	// Unescape dots in the property name
	unescapedPropertyName := unescapeDots(propertyName)

	return unescapedPropertyName, valueType, true
}

// quotePropertyNameIfNeeded wraps property names in backticks if they contain special characters
// MLMD requires backticks around property names that contain characters other than [0-9A-Za-z_]
// Based on MLMD proto documentation: "other than [0-9A-Za-z_], then the name need to be backquoted"
func quotePropertyNameIfNeeded(propertyName string) string {
	// Check if property name contains only valid characters [0-9A-Za-z_]
	// This is the most efficient approach using direct ASCII range checks
	for i := 0; i < len(propertyName); i++ {
		char := propertyName[i]
		// Check if character is NOT in valid ranges (using De Morgan's law)
		if (char < '0' || char > '9') && // NOT 0-9
			(char < 'A' || char > 'Z') && // NOT A-Z
			(char < 'a' || char > 'z') && // NOT a-z
			char != '_' { // NOT underscore
			// Found invalid character, needs quoting
			return fmt.Sprintf("`%s`", propertyName)
		}
	}

	// All characters are valid, no quoting needed
	return propertyName
}

// TranslateFilterQuery converts REST API field names to MLMD field names
func TranslateFilterQuery(restQuery string, entityType EntityType) (string, error) {
	if restQuery == "" {
		return "", nil
	}

	// Validate syntax including parentheses
	if err := validateQuerySyntax(restQuery); err != nil {
		return "", err
	}

	// Validate field names and types
	if err := validateRestFilterQuery(restQuery, entityType); err != nil {
		return "", err
	}

	// Get field mappings for this entity type
	fieldMappings := getFieldMappings()[entityType]

	// Translate field names using a simpler approach
	translatedQuery := restQuery

	// Match field names followed by operators - exclude parentheses from field names
	// This matches field names that can contain special characters but excludes parentheses and whitespace
	// Pattern explanation:
	// - ([^\\s()]+) captures field names without spaces or parentheses
	// - \\s* matches optional whitespace
	// - (=|!=|<|>|<=|>=|LIKE|IN|IS) captures operators
	fieldOperatorRegex := regexp.MustCompile(`([^\s()]+)\s*(=|!=|<|>|<=|>=|LIKE|IN|IS)`)

	translatedQuery = fieldOperatorRegex.ReplaceAllStringFunc(translatedQuery, func(match string) string {
		// Extract field name and operator from the match
		parts := fieldOperatorRegex.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}
		fieldName := parts[1]

		// Skip operators and keywords
		if isOperatorOrKeyword(fieldName) {
			return match
		}

		// Try to parse as custom property with escaped dots
		if propertyName, valueType, isCustomProperty := parseCustomPropertyField(fieldName); isCustomProperty {
			// Custom properties are not in field mappings (check unescaped name)
			if _, exists := fieldMappings[propertyName]; !exists {
				// Wrap property name in backticks if it contains special characters
				quotedPropertyName := quotePropertyNameIfNeeded(propertyName)
				replacement := fmt.Sprintf("custom_properties.%s.%s", quotedPropertyName, valueType)
				return strings.Replace(match, fieldName, replacement, 1)
			}
		}

		// Regular field mapping - context-sensitive translation
		if mapping, exists := fieldMappings[fieldName]; exists {
			return strings.Replace(match, fieldName, mapping.MLMDField, 1)
		}

		// Unknown field - assume it's a custom property with string value
		// Unescape dots in the field name for the MLMD property name
		unescapedFieldName := unescapeDots(fieldName)
		// Wrap property name in backticks if it contains special characters
		quotedPropertyName := quotePropertyNameIfNeeded(unescapedFieldName)
		replacement := fmt.Sprintf("custom_properties.%s.string_value", quotedPropertyName)
		return strings.Replace(match, fieldName, replacement, 1)
	})

	// Normalize quotes: Convert single quotes to double quotes for MLMD compatibility
	// MLMD expects double quotes for string literals
	translatedQuery = normalizeQuotesForMLMD(translatedQuery)

	return translatedQuery, nil
}

// normalizeQuotesForMLMD converts single quotes to double quotes for string literals
// MLMD requires double quotes for string literals in filter queries
// This function properly handles escaped quotes within string literals
func normalizeQuotesForMLMD(query string) string {
	var result strings.Builder
	i := 0

	for i < len(query) {
		char := query[i]

		// Handle single-quoted strings
		if char == '\'' {
			result.WriteByte('"') // Start with double quote
			i++                   // Skip the opening single quote

			// Process the content inside the single quotes
			for i < len(query) {
				char = query[i]

				if char == '\\' && i+1 < len(query) {
					// Handle escaped characters
					nextChar := query[i+1]
					switch nextChar {
					case '\'':
						// Escaped single quote inside single-quoted string -> unescaped single quote in double-quoted string
						result.WriteByte('\'')
						i += 2
					case '"':
						// Escaped double quote -> keep it escaped in double-quoted string
						result.WriteString("\\\"")
						i += 2
					case '\\':
						// Escaped backslash -> keep it escaped
						result.WriteString("\\\\")
						i += 2
					default:
						// Other escaped characters -> keep as is
						result.WriteByte('\\')
						result.WriteByte(nextChar)
						i += 2
					}
				} else if char == '\'' {
					// Unescaped single quote -> end of string
					result.WriteByte('"') // End with double quote
					i++
					break
				} else if char == '"' {
					// Unescaped double quote inside single-quoted string -> escape it
					result.WriteString("\\\"")
					i++
				} else {
					// Regular character
					result.WriteByte(char)
					i++
				}
			}
		} else if char == '"' {
			// Handle double-quoted strings - keep as is but process escaped characters
			result.WriteByte('"')
			i++

			for i < len(query) {
				char = query[i]

				if char == '\\' && i+1 < len(query) {
					// Escaped character - copy as is
					result.WriteByte('\\')
					result.WriteByte(query[i+1])
					i += 2
				} else if char == '"' {
					// End of double-quoted string
					result.WriteByte('"')
					i++
					break
				} else {
					// Regular character
					result.WriteByte(char)
					i++
				}
			}
		} else {
			// Regular character outside of quotes
			result.WriteByte(char)
			i++
		}
	}

	return result.String()
}

// validateQuerySyntax validates the basic syntax of the query
func validateQuerySyntax(query string) error {
	// Check balanced parentheses
	stack := 0
	for i, char := range query {
		switch char {
		case '(':
			stack++
		case ')':
			stack--
			if stack < 0 {
				return fmt.Errorf("unmatched closing parenthesis at position %d", i+1)
			}
		}
	}

	if stack > 0 {
		return fmt.Errorf("unmatched opening parentheses (%d)", stack)
	}

	return nil
}

// validateRestFilterQuery validates field names and custom property syntax
func validateRestFilterQuery(query string, entityType EntityType) error {
	// Use the same approach as translation - only validate field names before operators
	fieldOperatorRegex := regexp.MustCompile(`([^\s()]+)\s*(=|!=|<|>|<=|>=|LIKE|IN|IS)`)

	matches := fieldOperatorRegex.FindAllStringSubmatch(query, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		fieldName := match[1]

		// Skip operators and keywords
		if isOperatorOrKeyword(fieldName) {
			continue
		}

		// Check if field looks like a custom property (contains a dot) and validate value type
		if strings.Contains(fieldName, ".") {
			// Split on the last unescaped dot to check if it has a value type suffix
			dotPos := findLastUnescapedDot(fieldName)
			if dotPos != -1 {
				valueType := fieldName[dotPos+1:]
				// If the part after the dot is not a valid value type, it's an error
				// This catches cases like "field.invalid_type" where "invalid_type" is not valid
				if !isValidValueType(valueType) {
					return fmt.Errorf("invalid custom property value type '%s'. Valid types: string_value, int_value, double_value, bool_value", valueType)
				}
			}
		}

		// For regular fields, we allow unknown fields as they might be custom properties
		// The validation is lenient to support user-defined custom properties
	}

	return nil
}

// isOperatorOrKeyword checks if a string is an operator or keyword
func isOperatorOrKeyword(s string) bool {
	keywords := map[string]bool{
		"AND": true, "OR": true, "NOT": true,
		"and": true, "or": true, "not": true,
		"IS": true, "NULL": true, "LIKE": true,
		"is": true, "null": true, "like": true,
		"IN": true, "in": true,
	}
	return keywords[s]
}
