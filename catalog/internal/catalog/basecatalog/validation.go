package basecatalog

import (
	"fmt"
	"strings"
)

// supportedOperators defines the valid operators for named query field filters
var supportedOperators = map[string]bool{
	"=":      true,
	"!=":     true,
	">":      true,
	"<":      true,
	">=":     true,
	"<=":     true,
	"LIKE":   true,
	"ILIKE":  true,
	"IN":     true,
	"NOT IN": true,
}

// MCPServerFilterableFields is the set of valid field names that may be used in
// named query filters targeting MCP servers. Field names not in this set are
// rejected at config-parse time to catch typos early.
var MCPServerFilterableFields = map[string]bool{
	// Common Context properties
	"id":                       true,
	"name":                     true,
	"externalId":               true,
	"createTimeSinceEpoch":     true,
	"lastUpdateTimeSinceEpoch": true,
	// MCPServer-specific properties
	"source_id":        true,
	"base_name":        true,
	"description":      true,
	"provider":         true,
	"license":          true,
	"license_link":     true,
	"logo":             true,
	"readme":           true,
	"version":          true,
	"tags":             true,
	"transports":       true,
	"deploymentMode":   true,
	"documentationUrl": true,
	"repositoryUrl":    true,
	"sourceCode":       true,
	"publishedDate":    true,
	"lastUpdated":      true,
	"verifiedSource":   true,
	"secureEndpoint":   true,
	"sast":             true,
	"readOnlyTools":    true,
	"endpoints":        true,
	"artifacts":        true,
	"runtimeMetadata":  true,
}

// ValidateNamedQueries validates the structure and content of named queries
func ValidateNamedQueries(namedQueries map[string]map[string]FieldFilter) error {
	for queryName, fieldFilters := range namedQueries {
		if queryName == "" {
			return fmt.Errorf("named query name cannot be empty")
		}

		if len(fieldFilters) == 0 {
			return fmt.Errorf("named query '%s' must contain at least one field filter", queryName)
		}

		for fieldName, filter := range fieldFilters {
			if fieldName == "" {
				return fmt.Errorf("field name cannot be empty in named query '%s'", queryName)
			}

			if err := validateFieldFilter(queryName, fieldName, filter); err != nil {
				return err
			}
		}
	}

	return nil
}

// ValidateFieldNames checks that every field name used across all named queries
// is present in allowedFields. This catches typos at config-parse time before
// they silently produce empty results at runtime.
func ValidateFieldNames(namedQueries map[string]map[string]FieldFilter, allowedFields map[string]bool) error {
	for queryName, fieldFilters := range namedQueries {
		for fieldName := range fieldFilters {
			if !allowedFields[fieldName] {
				return fmt.Errorf("unknown field %q in named query %q: not a filterable MCP server field", fieldName, queryName)
			}
		}
	}
	return nil
}

// validateFieldFilter validates a single field filter within a named query
func validateFieldFilter(queryName, fieldName string, filter FieldFilter) error {
	if filter.Operator == "" {
		return fmt.Errorf("operator cannot be empty for field '%s' in named query '%s'", fieldName, queryName)
	}

	normalizedOperator := strings.ToUpper(filter.Operator)
	if !supportedOperators[normalizedOperator] {
		return fmt.Errorf("unsupported operator '%s' for field '%s' in named query '%s'", filter.Operator, fieldName, queryName)
	}

	if filter.Value == nil {
		return fmt.Errorf("value cannot be nil for field '%s' in named query '%s'", fieldName, queryName)
	}

	// Additional validation based on operator type
	switch normalizedOperator {
	case "IN", "NOT IN":
		// Value should be an array
		if _, ok := filter.Value.([]any); !ok {
			return fmt.Errorf("operator '%s' requires array value for field '%s' in named query '%s'", filter.Operator, fieldName, queryName)
		}
	}

	return nil
}
