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

func TestEvaluationRecordUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name             string
		jsonData         string
		wantModelID      string
		wantBenchmark    string
		wantCustomProps  map[string]interface{}
		wantErr          bool
		checkCustomProps bool
	}{
		{
			name: "complete evaluation record",
			jsonData: `{
				"model_id": "test-model-123",
				"benchmark": "aime24",
				"score": 63.3333,
				"created_at": 1609459200,
				"updated_at": 1609545600
			}`,
			wantModelID:   "test-model-123",
			wantBenchmark: "aime24",
			wantCustomProps: map[string]interface{}{
				"model_id":   "test-model-123",
				"benchmark":  "aime24",
				"score":      63.3333,
				"created_at": float64(1609459200),
				"updated_at": float64(1609545600),
			},
			wantErr:          false,
			checkCustomProps: true,
		},
		{
			name: "minimal evaluation record with only core fields",
			jsonData: `{
				"model_id": "minimal-model",
				"benchmark": "test-benchmark"
			}`,
			wantModelID:   "minimal-model",
			wantBenchmark: "test-benchmark",
			wantCustomProps: map[string]interface{}{
				"model_id":  "minimal-model",
				"benchmark": "test-benchmark",
			},
			wantErr:          false,
			checkCustomProps: true,
		},
		{
			name: "evaluation record with custom properties",
			jsonData: `{
				"model_id": "custom-model",
				"benchmark": "custom-bench",
				"score": 95.5,
				"custom_field_string": "custom value",
				"custom_field_number": 42,
				"custom_field_float": 3.14,
				"custom_field_bool": true
			}`,
			wantModelID:   "custom-model",
			wantBenchmark: "custom-bench",
			wantCustomProps: map[string]interface{}{
				"model_id":            "custom-model",
				"benchmark":           "custom-bench",
				"score":               95.5,
				"custom_field_string": "custom value",
				"custom_field_number": float64(42),
				"custom_field_float":  3.14,
				"custom_field_bool":   true,
			},
			wantErr:          false,
			checkCustomProps: true,
		},
		{
			name: "evaluation record with nested objects",
			jsonData: `{
				"model_id": "nested-model",
				"benchmark": "nested-bench",
				"custom_object": {
					"nested_key": "nested_value",
					"nested_number": 123
				},
				"custom_array": ["item1", "item2", "item3"]
			}`,
			wantModelID:      "nested-model",
			wantBenchmark:    "nested-bench",
			wantErr:          false,
			checkCustomProps: false, // Don't check deep equality for complex nested structures
		},
		{
			name: "evaluation record with null values",
			jsonData: `{
				"model_id": "null-model",
				"benchmark": "null-bench",
				"null_field": null,
				"score": 50.0
			}`,
			wantModelID:   "null-model",
			wantBenchmark: "null-bench",
			wantCustomProps: map[string]interface{}{
				"model_id":   "null-model",
				"benchmark":  "null-bench",
				"null_field": nil,
				"score":      50.0,
			},
			wantErr:          false,
			checkCustomProps: true,
		},
		{
			name: "evaluation record missing core fields",
			jsonData: `{
				"score": 75.5,
				"created_at": 1609459200
			}`,
			wantModelID:      "",
			wantBenchmark:    "",
			wantErr:          false,
			checkCustomProps: false,
		},
		{
			name: "evaluation record with wrong type for core fields",
			jsonData: `{
				"model_id": 123,
				"benchmark": 456,
				"score": 85.0
			}`,
			wantModelID:      "",
			wantBenchmark:    "",
			wantErr:          false,
			checkCustomProps: false,
		},
		{
			name:             "empty JSON object",
			jsonData:         `{}`,
			wantModelID:      "",
			wantBenchmark:    "",
			wantErr:          false,
			checkCustomProps: false,
		},
		{
			name:             "invalid JSON",
			jsonData:         `{"model_id": "invalid", "benchmark":}`,
			wantErr:          true,
			checkCustomProps: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var er evaluationRecord
			err := er.UnmarshalJSON([]byte(tt.jsonData))

			if (err != nil) != tt.wantErr {
				t.Errorf("evaluationRecord.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return // If we expected an error and got one, we're done
			}

			// Check core fields
			if er.ModelID != tt.wantModelID {
				t.Errorf("ModelID = %v, want %v", er.ModelID, tt.wantModelID)
			}
			if er.Benchmark != tt.wantBenchmark {
				t.Errorf("Benchmark = %v, want %v", er.Benchmark, tt.wantBenchmark)
			}

			// Check CustomProperties
			if er.CustomProperties == nil {
				t.Error("CustomProperties should not be nil")
			}

			// Optionally check custom properties in detail
			if tt.checkCustomProps {
				if len(er.CustomProperties) != len(tt.wantCustomProps) {
					t.Errorf("CustomProperties length = %v, want %v", len(er.CustomProperties), len(tt.wantCustomProps))
				}
				for key, wantValue := range tt.wantCustomProps {
					gotValue, exists := er.CustomProperties[key]
					if !exists {
						t.Errorf("CustomProperties missing key %v", key)
						continue
					}
					if gotValue != wantValue {
						t.Errorf("CustomProperties[%v] = %v (type %T), want %v (type %T)",
							key, gotValue, gotValue, wantValue, wantValue)
					}
				}
			}
		})
	}
}

func TestPerformanceRecordUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name             string
		jsonData         string
		wantID           string
		wantModelID      string
		wantCustomProps  map[string]interface{}
		wantErr          bool
		checkCustomProps bool
	}{
		{
			name: "complete performance record",
			jsonData: `{
				"id": "perf-123",
				"model_id": "test-model-456",
				"throughput": 1000.5,
				"latency_p50": 10.5,
				"latency_p95": 25.3,
				"latency_p99": 50.1,
				"created_at": 1609459200,
				"updated_at": 1609545600
			}`,
			wantID:      "perf-123",
			wantModelID: "test-model-456",
			wantCustomProps: map[string]interface{}{
				"id":          "perf-123",
				"model_id":    "test-model-456",
				"throughput":  1000.5,
				"latency_p50": 10.5,
				"latency_p95": 25.3,
				"latency_p99": 50.1,
				"created_at":  float64(1609459200),
				"updated_at":  float64(1609545600),
			},
			wantErr:          false,
			checkCustomProps: true,
		},
		{
			name: "minimal performance record with only core fields",
			jsonData: `{
				"id": "minimal-perf",
				"model_id": "minimal-model"
			}`,
			wantID:      "minimal-perf",
			wantModelID: "minimal-model",
			wantCustomProps: map[string]interface{}{
				"id":       "minimal-perf",
				"model_id": "minimal-model",
			},
			wantErr:          false,
			checkCustomProps: true,
		},
		{
			name: "performance record with custom properties",
			jsonData: `{
				"id": "custom-perf",
				"model_id": "custom-model",
				"throughput": 500.0,
				"custom_field_string": "custom value",
				"custom_field_number": 42,
				"custom_field_float": 3.14,
				"custom_field_bool": true
			}`,
			wantID:      "custom-perf",
			wantModelID: "custom-model",
			wantCustomProps: map[string]interface{}{
				"id":                  "custom-perf",
				"model_id":            "custom-model",
				"throughput":          500.0,
				"custom_field_string": "custom value",
				"custom_field_number": float64(42),
				"custom_field_float":  3.14,
				"custom_field_bool":   true,
			},
			wantErr:          false,
			checkCustomProps: true,
		},
		{
			name: "performance record with nested objects and arrays",
			jsonData: `{
				"id": "nested-perf",
				"model_id": "nested-model",
				"custom_object": {
					"nested_key": "nested_value",
					"nested_number": 123
				},
				"custom_array": ["item1", "item2", "item3"]
			}`,
			wantID:           "nested-perf",
			wantModelID:      "nested-model",
			wantErr:          false,
			checkCustomProps: false, // Don't check deep equality for complex nested structures
		},
		{
			name: "performance record with null values",
			jsonData: `{
				"id": "null-perf",
				"model_id": "null-model",
				"null_field": null,
				"throughput": 250.0
			}`,
			wantID:      "null-perf",
			wantModelID: "null-model",
			wantCustomProps: map[string]interface{}{
				"id":         "null-perf",
				"model_id":   "null-model",
				"null_field": nil,
				"throughput": 250.0,
			},
			wantErr:          false,
			checkCustomProps: true,
		},
		{
			name: "performance record missing core fields",
			jsonData: `{
				"throughput": 100.0,
				"latency_p50": 5.0
			}`,
			wantID:           "",
			wantModelID:      "",
			wantErr:          false,
			checkCustomProps: false,
		},
		{
			name: "performance record with wrong type for core fields",
			jsonData: `{
				"id": 123,
				"model_id": 456,
				"throughput": 500.0
			}`,
			wantID:           "",
			wantModelID:      "",
			wantErr:          false,
			checkCustomProps: false,
		},
		{
			name: "performance record with zero values",
			jsonData: `{
				"id": "zero-perf",
				"model_id": "zero-model",
				"throughput": 0,
				"latency_p50": 0.0,
				"created_at": 0
			}`,
			wantID:      "zero-perf",
			wantModelID: "zero-model",
			wantCustomProps: map[string]interface{}{
				"id":          "zero-perf",
				"model_id":    "zero-model",
				"throughput":  float64(0),
				"latency_p50": 0.0,
				"created_at":  float64(0),
			},
			wantErr:          false,
			checkCustomProps: true,
		},
		{
			name:             "empty JSON object",
			jsonData:         `{}`,
			wantID:           "",
			wantModelID:      "",
			wantErr:          false,
			checkCustomProps: false,
		},
		{
			name:             "invalid JSON",
			jsonData:         `{"id": "invalid", "model_id":}`,
			wantErr:          true,
			checkCustomProps: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pr performanceRecord
			err := pr.UnmarshalJSON([]byte(tt.jsonData))

			if (err != nil) != tt.wantErr {
				t.Errorf("performanceRecord.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return // If we expected an error and got one, we're done
			}

			// Check core fields
			if pr.ID != tt.wantID {
				t.Errorf("ID = %v, want %v", pr.ID, tt.wantID)
			}
			if pr.ModelID != tt.wantModelID {
				t.Errorf("ModelID = %v, want %v", pr.ModelID, tt.wantModelID)
			}

			// Check CustomProperties
			if pr.CustomProperties == nil {
				t.Error("CustomProperties should not be nil")
			}

			// Optionally check custom properties in detail
			if tt.checkCustomProps {
				if len(pr.CustomProperties) != len(tt.wantCustomProps) {
					t.Errorf("CustomProperties length = %v, want %v", len(pr.CustomProperties), len(tt.wantCustomProps))
				}
				for key, wantValue := range tt.wantCustomProps {
					gotValue, exists := pr.CustomProperties[key]
					if !exists {
						t.Errorf("CustomProperties missing key %v", key)
						continue
					}
					if gotValue != wantValue {
						t.Errorf("CustomProperties[%v] = %v (type %T), want %v (type %T)",
							key, gotValue, gotValue, wantValue, wantValue)
					}
				}
			}
		})
	}
}

func TestEvaluationRecordUnmarshalJSON_CoreFieldsInCustomProperties(t *testing.T) {
	// Test that core fields are included in CustomProperties
	jsonData := `{
		"model_id": "test-model",
		"benchmark": "test-benchmark",
		"score": 90.5
	}`

	var er evaluationRecord
	err := er.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify core fields are in CustomProperties
	if er.CustomProperties["model_id"] != "test-model" {
		t.Errorf("CustomProperties[model_id] = %v, want %v", er.CustomProperties["model_id"], "test-model")
	}
	if er.CustomProperties["benchmark"] != "test-benchmark" {
		t.Errorf("CustomProperties[benchmark] = %v, want %v", er.CustomProperties["benchmark"], "test-benchmark")
	}
	if er.CustomProperties["score"] != 90.5 {
		t.Errorf("CustomProperties[score] = %v, want %v", er.CustomProperties["score"], 90.5)
	}
}

func TestPerformanceRecordUnmarshalJSON_CoreFieldsInCustomProperties(t *testing.T) {
	// Test that core fields are included in CustomProperties
	jsonData := `{
		"id": "perf-id",
		"model_id": "test-model",
		"throughput": 1000.0
	}`

	var pr performanceRecord
	err := pr.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify core fields are in CustomProperties
	if pr.CustomProperties["id"] != "perf-id" {
		t.Errorf("CustomProperties[id] = %v, want %v", pr.CustomProperties["id"], "perf-id")
	}
	if pr.CustomProperties["model_id"] != "test-model" {
		t.Errorf("CustomProperties[model_id] = %v, want %v", pr.CustomProperties["model_id"], "test-model")
	}
	if pr.CustomProperties["throughput"] != 1000.0 {
		t.Errorf("CustomProperties[throughput] = %v, want %v", pr.CustomProperties["throughput"], 1000.0)
	}
}

func TestUnmarshalJSON_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		wantErr  bool
	}{
		{
			name:     "null JSON for evaluationRecord",
			jsonData: `null`,
			wantErr:  false, // null JSON unmarshals to empty map, not an error
		},
		{
			name:     "array instead of object for evaluationRecord",
			jsonData: `["not", "an", "object"]`,
			wantErr:  true,
		},
		{
			name:     "string instead of object for evaluationRecord",
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
		t.Run(tt.name+" (evaluationRecord)", func(t *testing.T) {
			var er evaluationRecord
			err := er.UnmarshalJSON([]byte(tt.jsonData))

			if (err != nil) != tt.wantErr {
				t.Errorf("evaluationRecord.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

		t.Run(tt.name+" (performanceRecord)", func(t *testing.T) {
			var pr performanceRecord
			err := pr.UnmarshalJSON([]byte(tt.jsonData))

			if (err != nil) != tt.wantErr {
				t.Errorf("performanceRecord.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
