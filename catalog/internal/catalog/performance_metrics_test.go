package catalog

import (
	"testing"
)

func TestParseMetadataJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     metadataJSON
		wantErr  bool
	}{
		{
			name: "complete metadata with all core fields",
			jsonData: `{
				"id": "test-model-123",
				"description": "A test model for unit testing",
				"readme": "# Test Model\nThis is a test model.",
				"maturity": "stable",
				"languages": ["python", "go"],
				"tasks": ["classification", "regression"],
				"provider_name": "test-provider",
				"logo": "https://example.com/logo.png",
				"license": "MIT",
				"license_link": "https://opensource.org/licenses/MIT",
				"library_name": "test-library",
				"created_at": 1609459200,
				"updated_at": 1609545600
			}`,
			want: metadataJSON{
				ID: "test-model-123",
			},
			wantErr: false,
		},
		{
			name: "minimal metadata with only required fields",
			jsonData: `{
				"id": "minimal-model"
			}`,
			want: metadataJSON{
				ID: "minimal-model",
			},
			wantErr: false,
		},
		{
			name: "metadata with custom properties",
			jsonData: `{
				"id": "custom-model",
				"description": "Model with custom properties",
				"custom_field_string": "custom value",
				"custom_field_number": 42,
				"custom_field_float": 3.14,
				"custom_field_bool": true,
				"custom_field_array": ["item1", "item2"],
				"custom_field_object": {"nested": "value"}
			}`,
			want: metadataJSON{
				ID: "custom-model",
			},
			wantErr: false,
		},
		{
			name: "metadata with mixed core and custom fields",
			jsonData: `{
				"id": "mixed-model",
				"description": "Mixed fields model",
				"languages": ["python"],
				"custom_version": "1.0.0",
				"custom_tags": ["ml", "ai"],
				"custom_config": {
					"batch_size": 32,
					"learning_rate": 0.001
				}
			}`,
			want: metadataJSON{
				ID: "mixed-model",
			},
			wantErr: false,
		},
		{
			name: "empty arrays and objects",
			jsonData: `{
				"id": "empty-arrays-model",
				"languages": [],
				"tasks": [],
				"custom_empty_array": [],
				"custom_empty_object": {}
			}`,
			want: metadataJSON{
				ID: "empty-arrays-model",
			},
			wantErr: false,
		},
		{
			name: "zero timestamps",
			jsonData: `{
				"id": "zero-timestamps-model",
				"created_at": 0,
				"updated_at": 0
			}`,
			want: metadataJSON{
				ID: "zero-timestamps-model",
			},
			wantErr: false,
		},
		{
			name: "null values in custom properties",
			jsonData: `{
				"id": "null-values-model",
				"custom_null_field": null,
				"custom_string": "not null"
			}`,
			want: metadataJSON{
				ID: "null-values-model",
			},
			wantErr: false,
		},
		{
			name:     "invalid JSON",
			jsonData: `{"id": "invalid-json", "description":}`,
			want:     metadataJSON{},
			wantErr:  true,
		},
		{
			name:     "empty JSON object",
			jsonData: `{}`,
			want:     metadataJSON{},
			wantErr:  true, // Should error because ID is required
		},
		{
			name:     "missing ID field",
			jsonData: `{"description": "has description but no id"}`,
			want:     metadataJSON{},
			wantErr:  true, // Should error because ID is required
		},
		{
			name: "JSON with type mismatches should fail",
			jsonData: `{
				"id": 123,
				"languages": "not-an-array",
				"created_at": "not-a-number"
			}`,
			want:    metadataJSON{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMetadataJSON([]byte(tt.jsonData))

			if (err != nil) != tt.wantErr {
				t.Errorf("parseMetadataJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return // If we expected an error and got one, we're done
			}

			// Compare all fields
			if got.ID != tt.want.ID {
				t.Errorf("parseMetadataJSON() ID = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

func TestParseMetadataJSON_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name:     "null JSON",
			jsonData: `null`,
			wantErr:  true, // Should error because ID will be empty
		},
		{
			name:     "array instead of object",
			jsonData: `["not", "an", "object"]`,
			wantErr:  true,
		},
		{
			name:     "string instead of object",
			jsonData: `"not an object"`,
			wantErr:  true,
		},
		{
			name:     "number instead of object",
			jsonData: `42`,
			wantErr:  true,
		},
		{
			name:     "boolean instead of object",
			jsonData: `true`,
			wantErr:  true,
		},
		{
			name:     "empty string",
			jsonData: ``,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseMetadataJSON([]byte(tt.jsonData))

			if (err != nil) != tt.wantErr {
				t.Errorf("parseMetadataJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseMetadataJSON_OnlyIDMatters(t *testing.T) {
	// Test that only the ID field is extracted, other fields are ignored
	jsonData := `{
		"id": "test-id",
		"custom_field": "ignored"
	}`

	metadata, err := parseMetadataJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify that only the ID field is populated
	if metadata.ID != "test-id" {
		t.Errorf("ID = %v, want %v", metadata.ID, "test-id")
	}
}
