package filter

import (
	"strings"
	"testing"
)

func TestCrossDatabaseILIKE(t *testing.T) {
	tests := []struct {
		name        string
		filterQuery string
		description string
	}{
		{
			name:        "ILIKE with string value",
			filterQuery: `name ILIKE "%Test%"`,
			description: "Should convert ILIKE to UPPER() comparison for cross-database compatibility",
		},
		{
			name:        "ILIKE with custom property",
			filterQuery: `framework ILIKE "%PyTorch%"`,
			description: "Should handle ILIKE for custom properties using property table subqueries",
		},
		{
			name:        "Regular LIKE unchanged",
			filterQuery: `name LIKE "%test%"`,
			description: "Regular LIKE should remain unchanged",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the filter query
			filterExpr, err := Parse(tt.filterQuery)
			if err != nil {
				t.Fatalf("Failed to parse filter query: %v", err)
			}

			// Create query builder
			queryBuilder := NewQueryBuilderForRestEntity(RestEntityExperiment)

			// Generate condition string
			conditionResult := queryBuilder.buildConditionString(filterExpr)

			t.Logf("Query: %s", tt.filterQuery)
			t.Logf("Generated condition: %s", conditionResult.condition)
			t.Logf("Generated args: %v", conditionResult.args)

			// Verify UPPER() usage based on operator type
			if strings.Contains(tt.filterQuery, "ILIKE") {
				// For ILIKE conditions, verify UPPER() is used
				if !strings.Contains(conditionResult.condition, "UPPER(") {
					t.Errorf("Expected ILIKE to use UPPER() function, but got: %s", conditionResult.condition)
				}
			} else if strings.Contains(tt.filterQuery, "LIKE") {
				// For regular LIKE, verify UPPER() is NOT used
				if strings.Contains(conditionResult.condition, "UPPER(") {
					t.Errorf("Expected regular LIKE to not use UPPER() function, but got: %s", conditionResult.condition)
				}
			}
		})
	}
}
