package basecatalog

import (
	"fmt"
	"sort"
	"strings"
)

// FieldFiltersToFilterQuery converts a map of FieldFilter (from a named query) into
// a filterQuery string that the existing filter.Parse() infrastructure can process.
// Conditions are joined with AND in sorted field-name order for deterministic output.
//
// Supported operators and their output:
//   - =, !=, >, <, >=, <=, LIKE, ILIKE: `field op 'value'` (strings quoted, numerics/bools unquoted)
//   - IN: `field IN (val1, val2, ...)`
//   - NOT IN: expanded as `field != val1 AND field != val2 ...` (parser does not support NOT IN natively)
func FieldFiltersToFilterQuery(filters map[string]FieldFilter) (string, error) {
	if len(filters) == 0 {
		return "", nil
	}

	// Sort field names for deterministic output.
	fields := make([]string, 0, len(filters))
	for f := range filters {
		fields = append(fields, f)
	}
	sort.Strings(fields)

	var parts []string
	for _, field := range fields {
		ff := filters[field]
		op := strings.ToUpper(ff.Operator)

		switch op {
		case "IN":
			vals, ok := ff.Value.([]any)
			if !ok {
				return "", fmt.Errorf("field %q: IN operator requires array value, got %T", field, ff.Value)
			}
			items := make([]string, 0, len(vals))
			for _, v := range vals {
				items = append(items, formatValue(v))
			}
			parts = append(parts, fmt.Sprintf("%s IN (%s)", field, strings.Join(items, ", ")))
		case "NOT IN":
			// The filter parser does not support NOT IN natively, so expand as multiple != conditions.
			vals, ok := ff.Value.([]any)
			if !ok {
				return "", fmt.Errorf("field %q: NOT IN operator requires array value, got %T", field, ff.Value)
			}
			for _, v := range vals {
				parts = append(parts, fmt.Sprintf("%s != %s", field, formatValue(v)))
			}
		default:
			parts = append(parts, fmt.Sprintf("%s %s %s", field, op, formatValue(ff.Value)))
		}
	}

	return strings.Join(parts, " AND "), nil
}

// formatValue formats a single filter value for inclusion in a filterQuery string.
func formatValue(v any) string {
	switch val := v.(type) {
	case string:
		// Escape single quotes inside string values.
		escaped := strings.ReplaceAll(val, "'", "\\'")
		return fmt.Sprintf("'%s'", escaped)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", val)
	}
}
