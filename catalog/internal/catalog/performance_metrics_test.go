package catalog

import (
	"reflect"
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
				ID:               "test-model-123",
				Description:      "A test model for unit testing",
				Readme:           "# Test Model\nThis is a test model.",
				Maturity:         "stable",
				Languages:        []string{"python", "go"},
				Tasks:            []string{"classification", "regression"},
				Provider:         "test-provider",
				Logo:             "https://example.com/logo.png",
				License:          "MIT",
				LicenseLink:      "https://opensource.org/licenses/MIT",
				LibraryName:      "test-library",
				CreatedAt:        1609459200,
				UpdatedAt:        1609545600,
				CustomProperties: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "minimal metadata with only required fields",
			jsonData: `{
				"id": "minimal-model"
			}`,
			want: metadataJSON{
				ID:               "minimal-model",
				Description:      "",
				Readme:           "",
				Maturity:         "",
				Languages:        nil,
				Tasks:            nil,
				Provider:         "",
				Logo:             "",
				License:          "",
				LicenseLink:      "",
				LibraryName:      "",
				CreatedAt:        0,
				UpdatedAt:        0,
				CustomProperties: map[string]interface{}{},
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
				ID:          "custom-model",
				Description: "Model with custom properties",
				CustomProperties: map[string]interface{}{
					"custom_field_string": "custom value",
					"custom_field_number": float64(42), // JSON numbers are parsed as float64
					"custom_field_float":  3.14,
					"custom_field_bool":   true,
					"custom_field_array":  []interface{}{"item1", "item2"},
					"custom_field_object": map[string]interface{}{"nested": "value"},
				},
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
				ID:          "mixed-model",
				Description: "Mixed fields model",
				Languages:   []string{"python"},
				CustomProperties: map[string]interface{}{
					"custom_version": "1.0.0",
					"custom_tags":    []interface{}{"ml", "ai"},
					"custom_config": map[string]interface{}{
						"batch_size":    float64(32),
						"learning_rate": 0.001,
					},
				},
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
				ID:        "empty-arrays-model",
				Languages: []string{},
				Tasks:     []string{},
				CustomProperties: map[string]interface{}{
					"custom_empty_array":  []interface{}{},
					"custom_empty_object": map[string]interface{}{},
				},
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
				ID:               "zero-timestamps-model",
				CreatedAt:        0,
				UpdatedAt:        0,
				CustomProperties: map[string]interface{}{},
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
				CustomProperties: map[string]interface{}{
					"custom_null_field": nil,
					"custom_string":     "not null",
				},
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
			want: metadataJSON{
				CustomProperties: map[string]interface{}{},
			},
			wantErr: false,
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
			if got.Description != tt.want.Description {
				t.Errorf("parseMetadataJSON() Description = %v, want %v", got.Description, tt.want.Description)
			}
			if got.Readme != tt.want.Readme {
				t.Errorf("parseMetadataJSON() Readme = %v, want %v", got.Readme, tt.want.Readme)
			}
			if got.Maturity != tt.want.Maturity {
				t.Errorf("parseMetadataJSON() Maturity = %v, want %v", got.Maturity, tt.want.Maturity)
			}
			if !reflect.DeepEqual(got.Languages, tt.want.Languages) {
				t.Errorf("parseMetadataJSON() Languages = %v, want %v", got.Languages, tt.want.Languages)
			}
			if !reflect.DeepEqual(got.Tasks, tt.want.Tasks) {
				t.Errorf("parseMetadataJSON() Tasks = %v, want %v", got.Tasks, tt.want.Tasks)
			}
			if got.Provider != tt.want.Provider {
				t.Errorf("parseMetadataJSON() Provider = %v, want %v", got.Provider, tt.want.Provider)
			}
			if got.Logo != tt.want.Logo {
				t.Errorf("parseMetadataJSON() Logo = %v, want %v", got.Logo, tt.want.Logo)
			}
			if got.License != tt.want.License {
				t.Errorf("parseMetadataJSON() License = %v, want %v", got.License, tt.want.License)
			}
			if got.LicenseLink != tt.want.LicenseLink {
				t.Errorf("parseMetadataJSON() LicenseLink = %v, want %v", got.LicenseLink, tt.want.LicenseLink)
			}
			if got.LibraryName != tt.want.LibraryName {
				t.Errorf("parseMetadataJSON() LibraryName = %v, want %v", got.LibraryName, tt.want.LibraryName)
			}
			if got.CreatedAt != tt.want.CreatedAt {
				t.Errorf("parseMetadataJSON() CreatedAt = %v, want %v", got.CreatedAt, tt.want.CreatedAt)
			}
			if got.UpdatedAt != tt.want.UpdatedAt {
				t.Errorf("parseMetadataJSON() UpdatedAt = %v, want %v", got.UpdatedAt, tt.want.UpdatedAt)
			}

			// Compare custom properties
			if !reflect.DeepEqual(got.CustomProperties, tt.want.CustomProperties) {
				t.Errorf("parseMetadataJSON() CustomProperties = %v, want %v", got.CustomProperties, tt.want.CustomProperties)
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
			wantErr:  false, // null JSON actually unmarshals to zero value struct
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

func TestParseMetadataJSON_KnownFieldsExclusion(t *testing.T) {
	// Test that all known core fields are properly excluded from custom properties
	jsonData := `{
		"id": "test-id",
		"description": "test description",
		"readme": "test readme",
		"maturity": "test maturity",
		"languages": ["python"],
		"tasks": ["classification"],
		"provider_name": "test provider",
		"logo": "test logo",
		"license": "test license",
		"license_link": "test license link",
		"library_name": "test library",
		"created_at": 1609459200,
		"updated_at": 1609545600,
		"custom_field": "should be in custom properties"
	}`

	metadata, err := parseMetadataJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify that only the custom field is in CustomProperties
	expectedCustomProps := map[string]interface{}{
		"custom_field": "should be in custom properties",
	}

	if !reflect.DeepEqual(metadata.CustomProperties, expectedCustomProps) {
		t.Errorf("CustomProperties = %v, want %v", metadata.CustomProperties, expectedCustomProps)
	}

	// Verify that all core fields are populated correctly
	if metadata.ID != "test-id" {
		t.Errorf("ID = %v, want %v", metadata.ID, "test-id")
	}
	if metadata.Description != "test description" {
		t.Errorf("Description = %v, want %v", metadata.Description, "test description")
	}
	// Add more assertions for other core fields as needed
}

func TestParseMetadataJSON_CustomPropertiesInitialization(t *testing.T) {
	// Test that CustomProperties map is always initialized, even with no custom fields
	jsonData := `{
		"id": "test-id",
		"description": "only core fields"
	}`

	metadata, err := parseMetadataJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if metadata.CustomProperties == nil {
		t.Error("CustomProperties should be initialized, got nil")
	}

	if len(metadata.CustomProperties) != 0 {
		t.Errorf("CustomProperties should be empty, got %v", metadata.CustomProperties)
	}
}
