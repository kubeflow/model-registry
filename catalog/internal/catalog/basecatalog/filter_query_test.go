package basecatalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldFiltersToFilterQuery(t *testing.T) {
	tests := []struct {
		name     string
		filters  map[string]FieldFilter
		expected string
	}{
		{
			name:     "empty filters",
			filters:  map[string]FieldFilter{},
			expected: "",
		},
		{
			name:     "nil filters",
			filters:  nil,
			expected: "",
		},
		{
			name: "single equality filter with string value",
			filters: map[string]FieldFilter{
				"provider": {Operator: "=", Value: "OpenAI"},
			},
			expected: "provider = 'OpenAI'",
		},
		{
			name: "single equality filter with integer value",
			filters: map[string]FieldFilter{
				"rating": {Operator: ">=", Value: 4},
			},
			expected: "rating >= 4",
		},
		{
			name: "boolean value",
			filters: map[string]FieldFilter{
				"verified": {Operator: "=", Value: true},
			},
			expected: "verified = true",
		},
		{
			name: "false boolean value",
			filters: map[string]FieldFilter{
				"deprecated": {Operator: "=", Value: false},
			},
			expected: "deprecated = false",
		},
		{
			name: "IN operator with string values",
			filters: map[string]FieldFilter{
				"status": {Operator: "IN", Value: []any{"active", "beta"}},
			},
			expected: "status IN ('active', 'beta')",
		},
		{
			name: "IN operator with numeric values",
			filters: map[string]FieldFilter{
				"tier": {Operator: "IN", Value: []any{1, 2, 3}},
			},
			expected: "tier IN (1, 2, 3)",
		},
		{
			name: "NOT IN operator expands to multiple != conditions",
			filters: map[string]FieldFilter{
				"status": {Operator: "NOT IN", Value: []any{"deprecated", "archived"}},
			},
			// NOT IN expands to individual != conditions because the filter
			// parser does not support NOT IN natively.
			expected: "status != 'deprecated' AND status != 'archived'",
		},
		{
			name: "LIKE operator",
			filters: map[string]FieldFilter{
				"name": {Operator: "LIKE", Value: "%assistant%"},
			},
			expected: "name LIKE '%assistant%'",
		},
		{
			name: "multiple filters joined with AND in sorted order",
			filters: map[string]FieldFilter{
				"provider": {Operator: "=", Value: "OpenAI"},
				"verified": {Operator: "=", Value: true},
			},
			// Fields sorted alphabetically: provider before verified
			expected: "provider = 'OpenAI' AND verified = true",
		},
		{
			name: "string with single quote is escaped",
			filters: map[string]FieldFilter{
				"name": {Operator: "=", Value: "O'Reilly"},
			},
			expected: "name = 'O\\'Reilly'",
		},
		{
			name: "ILIKE operator",
			filters: map[string]FieldFilter{
				"description": {Operator: "ILIKE", Value: "%search%"},
			},
			expected: "description ILIKE '%search%'",
		},
		{
			name: "lowercase operator is normalized to uppercase",
			filters: map[string]FieldFilter{
				"status": {Operator: "like", Value: "%active%"},
			},
			expected: "status LIKE '%active%'",
		},
		{
			name: "not-equal operator",
			filters: map[string]FieldFilter{
				"status": {Operator: "!=", Value: "deprecated"},
			},
			expected: "status != 'deprecated'",
		},
		{
			name: "float64 value from YAML number parsing",
			filters: map[string]FieldFilter{
				"rating": {Operator: ">=", Value: float64(4.5)},
			},
			expected: "rating >= 4.5",
		},
	}

	t.Run("error on non-array value for IN operator", func(t *testing.T) {
		filters := map[string]FieldFilter{
			"status": {Operator: "IN", Value: "not-an-array"},
		}
		_, err := FieldFiltersToFilterQuery(filters)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "IN operator requires array value")
	})

	t.Run("error on non-array value for NOT IN operator", func(t *testing.T) {
		filters := map[string]FieldFilter{
			"status": {Operator: "NOT IN", Value: "not-an-array"},
		}
		_, err := FieldFiltersToFilterQuery(filters)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "NOT IN operator requires array value")
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FieldFiltersToFilterQuery(tt.filters)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
