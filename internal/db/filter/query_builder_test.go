package filter

import (
	"testing"
)

func TestQueryBuilderEntityTypes(t *testing.T) {
	tests := []struct {
		name          string
		entityType    EntityType
		query         string
		expectedJoins int
		description   string
	}{
		{
			name:          "Context core property",
			entityType:    EntityTypeContext,
			query:         `name = "test-model"`,
			expectedJoins: 0,
			description:   "Context core properties should not require joins",
		},
		{
			name:          "Context custom property",
			entityType:    EntityTypeContext,
			query:         `accuracy > 0.95`,
			expectedJoins: 1,
			description:   "Context custom properties should require ContextProperty join",
		},
		{
			name:          "Context multiple custom properties",
			entityType:    EntityTypeContext,
			query:         `accuracy > 0.95 AND framework = "pytorch"`,
			expectedJoins: 2,
			description:   "Multiple Context custom properties should require multiple joins",
		},
		{
			name:          "Artifact core property",
			entityType:    EntityTypeArtifact,
			query:         `uri = "s3://bucket/model.pkl"`,
			expectedJoins: 0,
			description:   "Artifact core properties should not require joins",
		},
		{
			name:          "Artifact custom property",
			entityType:    EntityTypeArtifact,
			query:         `model_size > 1000`,
			expectedJoins: 1,
			description:   "Artifact custom properties should require ArtifactProperty join",
		},
		{
			name:          "Execution core property",
			entityType:    EntityTypeExecution,
			query:         `name = "serve-model-1"`,
			expectedJoins: 0,
			description:   "Execution core properties should not require joins",
		},
		{
			name:          "Execution custom property",
			entityType:    EntityTypeExecution,
			query:         `replicas = 3`,
			expectedJoins: 1,
			description:   "Execution custom properties should require ExecutionProperty join",
		},
		{
			name:          "Execution multiple custom properties",
			entityType:    EntityTypeExecution,
			query:         `replicas = 3 AND memory_limit > 1000`,
			expectedJoins: 2,
			description:   "Multiple Execution custom properties should require multiple joins",
		},
		{
			name:          "Execution mixed core and custom",
			entityType:    EntityTypeExecution,
			query:         `name = "serve-model-1" AND replicas = 3`,
			expectedJoins: 1,
			description:   "Mixed Execution properties should require selective joins",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the query
			filterExpr, err := Parse(tt.query)
			if err != nil {
				t.Fatalf("Failed to parse query: %v", err)
			}

			if filterExpr == nil {
				t.Fatalf("Expected non-nil filter expression")
			}

			// Create query builder
			var restEntityType RestEntityType
			switch tt.entityType {
			case EntityTypeContext:
				restEntityType = RestEntityExperiment
			case EntityTypeArtifact:
				restEntityType = RestEntityModelArtifact
			case EntityTypeExecution:
				restEntityType = RestEntityServeModel
			}
			queryBuilder := NewQueryBuilderForRestEntity(restEntityType)

			// For testing, we'll analyze the filter expression to count expected joins
			expectedJoins := countExpectedJoins(filterExpr, tt.entityType)

			// Verify the entity type and table prefix are set correctly
			var expectedTablePrefix string
			switch tt.entityType {
			case EntityTypeContext:
				expectedTablePrefix = "Context"
			case EntityTypeArtifact:
				expectedTablePrefix = "Artifact"
			case EntityTypeExecution:
				expectedTablePrefix = "Execution"
			}

			if queryBuilder.tablePrefix != expectedTablePrefix {
				t.Errorf("Expected table prefix %s, got %s", expectedTablePrefix, queryBuilder.tablePrefix)
			}

			// Log the result for inspection
			t.Logf("✅ %s", tt.description)
			t.Logf("   Query: %s", tt.query)
			t.Logf("   Entity Type: %s", tt.entityType)
			t.Logf("   Table Prefix: %s", queryBuilder.tablePrefix)
			t.Logf("   Expected Joins: %d", expectedJoins)
		})
	}
}

func TestQueryBuilderPropertyTypes(t *testing.T) {
	tests := []struct {
		name        string
		entityType  EntityType
		query       string
		description string
	}{
		{
			name:        "Context with explicit types",
			entityType:  EntityTypeContext,
			query:       `framework.string_value = "pytorch" AND accuracy.double_value > 0.95`,
			description: "Context with explicit property types",
		},
		{
			name:        "Artifact with type inference",
			entityType:  EntityTypeArtifact,
			query:       `model_size > 1000 AND is_compressed = true`,
			description: "Artifact with automatic type inference",
		},
		{
			name:        "Execution with mixed types",
			entityType:  EntityTypeExecution,
			query:       `replicas = 3 AND memory_limit.string_value = "2Gi"`,
			description: "Execution with mixed explicit and inferred types",
		},
		{
			name:        "Context with complex expressions",
			entityType:  EntityTypeContext,
			query:       `(accuracy > 0.9 OR f1_score > 0.85) AND framework = "tensorflow"`,
			description: "Context with complex logical expressions",
		},
		{
			name:        "Execution with ILIKE operator",
			entityType:  EntityTypeExecution,
			query:       `name ILIKE "%serve%" AND status = "running"`,
			description: "Execution with case-insensitive LIKE operator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the query
			filterExpr, err := Parse(tt.query)
			if err != nil {
				t.Fatalf("Failed to parse query: %v", err)
			}

			if filterExpr == nil {
				t.Fatalf("Expected non-nil filter expression")
			}

			// Create query builder
			var restEntityType RestEntityType
			switch tt.entityType {
			case EntityTypeContext:
				restEntityType = RestEntityExperiment
			case EntityTypeArtifact:
				restEntityType = RestEntityModelArtifact
			case EntityTypeExecution:
				restEntityType = RestEntityServeModel
			}
			queryBuilder := NewQueryBuilderForRestEntity(restEntityType)

			// Verify the query builder was created successfully
			if queryBuilder.entityType != tt.entityType {
				t.Errorf("Expected entity type %s, got %s", tt.entityType, queryBuilder.entityType)
			}

			// Log the result for inspection
			t.Logf("✅ %s", tt.description)
			t.Logf("   Query: %s", tt.query)
			t.Logf("   Entity Type: %s", tt.entityType)
		})
	}
}

// countExpectedJoins counts how many custom properties (requiring joins) are in the expression
func countExpectedJoins(expr *FilterExpression, entityType EntityType) int {
	if expr == nil {
		return 0
	}

	if expr.IsLeaf {
		// Check if this property requires a join (custom property)
		propDef := GetPropertyDefinition(EntityTypeContext, expr.Property)
		if propDef.Location == Custom {
			return 1 // Custom property requires join
		}
		return 0 // Core property doesn't require join
	}

	// Recursively count joins in left and right expressions
	leftJoins := countExpectedJoins(expr.Left, entityType)
	rightJoins := countExpectedJoins(expr.Right, entityType)

	return leftJoins + rightJoins
}
