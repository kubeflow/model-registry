package basecatalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateNamedQueries(t *testing.T) {
	tests := []struct {
		name          string
		namedQueries  map[string]map[string]FieldFilter
		expectError   bool
		errorContains string
	}{
		{
			name: "valid named queries",
			namedQueries: map[string]map[string]FieldFilter{
				"test-query": {
					"name":     {Operator: "=", Value: "my-server"},
					"provider": {Operator: "=", Value: "Anthropic"},
				},
			},
			expectError: false,
		},
		{
			name: "invalid operator",
			namedQueries: map[string]map[string]FieldFilter{
				"test-query": {
					"name": {Operator: "INVALID", Value: "value"},
				},
			},
			expectError:   true,
			errorContains: "unsupported operator 'INVALID'",
		},
		{
			name: "empty operator",
			namedQueries: map[string]map[string]FieldFilter{
				"test-query": {
					"name": {Operator: "", Value: "value"},
				},
			},
			expectError:   true,
			errorContains: "operator cannot be empty",
		},
		{
			name: "nil value",
			namedQueries: map[string]map[string]FieldFilter{
				"test-query": {
					"name": {Operator: "=", Value: nil},
				},
			},
			expectError:   true,
			errorContains: "value cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNamedQueries(tt.namedQueries)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoaderValidationIntegration(t *testing.T) {
	// Test that the validation logic works correctly
	// (The loader integration is tested in the main loader tests)

	// Test with a valid config
	validConfig := &SourceConfig{
		Catalogs: []ModelSource{},
		NamedQueries: map[string]map[string]FieldFilter{
			"valid-query": {
				"name": {Operator: "=", Value: "my-server"},
			},
		},
	}

	// This should succeed (we're testing the validation logic, not the file I/O)
	err := ValidateNamedQueries(validConfig.NamedQueries)
	assert.NoError(t, err)

	// Test with an invalid config
	invalidConfig := &SourceConfig{
		Catalogs: []ModelSource{},
		NamedQueries: map[string]map[string]FieldFilter{
			"invalid-query": {
				"name": {Operator: "INVALID_OP", Value: "value"},
			},
		},
	}

	// This should fail
	err = ValidateNamedQueries(invalidConfig.NamedQueries)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported operator 'INVALID_OP'")
}
