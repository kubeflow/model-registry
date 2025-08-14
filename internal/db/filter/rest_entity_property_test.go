package filter

import (
	"testing"
)

func TestRestEntityPropertyTypeDistinction(t *testing.T) {
	tests := []struct {
		name              string
		restEntityType    RestEntityType
		propertyName      string
		expectedLocation  PropertyLocation
		expectedValueType string
		description       string
	}{
		// Well-known properties should use their defined types/locations
		{
			name:              "RegisteredModel well-known property",
			restEntityType:    RestEntityRegisteredModel,
			propertyName:      "name",
			expectedLocation:  EntityTable,
			expectedValueType: StringValueType,
			description:       "RegisteredModel name should be EntityTable/string_value",
		},
		{
			name:              "RegisteredModel well-known state property",
			restEntityType:    RestEntityRegisteredModel,
			propertyName:      "state",
			expectedLocation:  PropertyTable,
			expectedValueType: StringValueType,
			description:       "RegisteredModel state should be PropertyTable/string_value",
		},
		{
			name:              "ServeModel well-known property",
			restEntityType:    RestEntityServeModel,
			propertyName:      "lastKnownState",
			expectedLocation:  EntityTable,
			expectedValueType: IntValueType,
			description:       "ServeModel lastKnownState should be EntityTable/int_value",
		},
		{
			name:              "Metric well-known property",
			restEntityType:    RestEntityMetric,
			propertyName:      "step",
			expectedLocation:  PropertyTable,
			expectedValueType: IntValueType,
			description:       "Metric step should be PropertyTable/int_value",
		},

		// Custom properties should always be Custom/string_value (default)
		{
			name:              "RegisteredModel custom property",
			restEntityType:    RestEntityRegisteredModel,
			propertyName:      "customProperty123",
			expectedLocation:  Custom,
			expectedValueType: StringValueType,
			description:       "Custom properties should be Custom/string_value",
		},
		{
			name:              "RegisteredModel using experimentId as custom",
			restEntityType:    RestEntityRegisteredModel,
			propertyName:      "experimentId",
			expectedLocation:  Custom,
			expectedValueType: StringValueType,
			description:       "Properties from other entities should be treated as custom",
		},
		{
			name:              "Metric using modelFormatName as custom",
			restEntityType:    RestEntityMetric,
			propertyName:      "modelFormatName",
			expectedLocation:  Custom,
			expectedValueType: StringValueType,
			description:       "Properties from other entities should be treated as custom",
		},
		{
			name:              "ServeModel custom property",
			restEntityType:    RestEntityServeModel,
			propertyName:      "myCustomField",
			expectedLocation:  Custom,
			expectedValueType: StringValueType,
			description:       "Custom properties should be Custom/string_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			propDef := GetPropertyDefinitionForRestEntity(tt.restEntityType, tt.propertyName)

			if propDef.Location != tt.expectedLocation {
				t.Errorf("Expected location %v, got %v for property '%s' on entity '%s'",
					tt.expectedLocation, propDef.Location, tt.propertyName, tt.restEntityType)
			}

			if propDef.ValueType != tt.expectedValueType {
				t.Errorf("Expected value type %s, got %s for property '%s' on entity '%s'",
					tt.expectedValueType, propDef.ValueType, tt.propertyName, tt.restEntityType)
			}

			t.Logf("✅ %s: property '%s' on entity '%s' correctly mapped to Location=%v, ValueType=%s",
				tt.description, tt.propertyName, tt.restEntityType, propDef.Location, propDef.ValueType)
		})
	}
}

func TestMLMDEntityTypeMapping(t *testing.T) {
	tests := []struct {
		name             string
		restEntityType   RestEntityType
		expectedMLMDType EntityType
		description      string
	}{
		{
			name:             "RegisteredModel maps to Context",
			restEntityType:   RestEntityRegisteredModel,
			expectedMLMDType: EntityTypeContext,
			description:      "RegisteredModel should map to MLMD Context type",
		},
		{
			name:             "ModelArtifact maps to Artifact",
			restEntityType:   RestEntityModelArtifact,
			expectedMLMDType: EntityTypeArtifact,
			description:      "ModelArtifact should map to MLMD Artifact type",
		},
		{
			name:             "ServeModel maps to Execution",
			restEntityType:   RestEntityServeModel,
			expectedMLMDType: EntityTypeExecution,
			description:      "ServeModel should map to MLMD Execution type",
		},
		{
			name:             "Experiment maps to Context",
			restEntityType:   RestEntityExperiment,
			expectedMLMDType: EntityTypeContext,
			description:      "Experiment should map to MLMD Context type",
		},
		{
			name:             "Metric maps to Artifact",
			restEntityType:   RestEntityMetric,
			expectedMLMDType: EntityTypeArtifact,
			description:      "Metric should map to MLMD Artifact type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mlmdType := GetMLMDEntityType(tt.restEntityType)

			if mlmdType != tt.expectedMLMDType {
				t.Errorf("Expected MLMD type %s, got %s for REST entity '%s'",
					tt.expectedMLMDType, mlmdType, tt.restEntityType)
			} else {
				t.Logf("✅ %s: %s correctly maps to %s",
					tt.description, tt.restEntityType, mlmdType)
			}
		})
	}
}

func TestQueryBuilderWithPropertyTypeDistinction(t *testing.T) {
	tests := []struct {
		name           string
		restEntityType RestEntityType
		description    string
	}{
		{
			name:           "RegisteredModel query builder",
			restEntityType: RestEntityRegisteredModel,
			description:    "Should create query builder with proper REST entity type",
		},
		{
			name:           "ServeModel query builder",
			restEntityType: RestEntityServeModel,
			description:    "Should create query builder with proper REST entity type",
		},
		{
			name:           "Metric query builder",
			restEntityType: RestEntityMetric,
			description:    "Should create query builder with proper REST entity type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create REST entity-aware query builder
			qb := NewQueryBuilderForRestEntity(tt.restEntityType)

			// Verify the query builder was created with correct types
			if qb.restEntityType != tt.restEntityType {
				t.Errorf("Expected REST entity type %s, got %s", tt.restEntityType, qb.restEntityType)
			}

			expectedMLMDType := GetMLMDEntityType(tt.restEntityType)
			if qb.entityType != expectedMLMDType {
				t.Errorf("Expected MLMD entity type %s, got %s", expectedMLMDType, qb.entityType)
			}

			t.Logf("✅ %s: Query builder created successfully with REST entity type %s (MLMD: %s)",
				tt.description, qb.restEntityType, qb.entityType)
		})
	}
}
