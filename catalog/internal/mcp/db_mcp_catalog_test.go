package mcp

import (
	"testing"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

func TestConvertNamedQueryToFilterQuery(t *testing.T) {
	tests := []struct {
		name           string
		fieldFilters   map[string]model.FieldFilter
		expectedResult string
	}{
		{
			name:           "empty filters",
			fieldFilters:   map[string]model.FieldFilter{},
			expectedResult: "",
		},
		{
			name:           "nil filters",
			fieldFilters:   nil,
			expectedResult: "",
		},
		{
			name: "single string equality",
			fieldFilters: map[string]model.FieldFilter{
				"provider": {Operator: "=", Value: "anthropic"},
			},
			expectedResult: "provider = 'anthropic'",
		},
		{
			name: "single boolean value",
			fieldFilters: map[string]model.FieldFilter{
				"verifiedSource": {Operator: "=", Value: true},
			},
			expectedResult: "verifiedSource = true",
		},
		{
			name: "single integer value",
			fieldFilters: map[string]model.FieldFilter{
				"toolCount": {Operator: ">=", Value: float64(5)},
			},
			expectedResult: "toolCount >= 5",
		},
		{
			name: "IN operator with string array",
			fieldFilters: map[string]model.FieldFilter{
				"tags": {Operator: "IN", Value: []interface{}{"security", "verified"}},
			},
			expectedResult: "tags IN ('security', 'verified')",
		},
		{
			name: "LIKE operator",
			fieldFilters: map[string]model.FieldFilter{
				"name": {Operator: "LIKE", Value: "%github%"},
			},
			expectedResult: "name LIKE '%github%'",
		},
		{
			name: "not equals operator",
			fieldFilters: map[string]model.FieldFilter{
				"status": {Operator: "!=", Value: "deprecated"},
			},
			expectedResult: "status != 'deprecated'",
		},
		{
			name: "ANYOF operator (alias for IN)",
			fieldFilters: map[string]model.FieldFilter{
				"deploymentMode": {Operator: "ANYOF", Value: []interface{}{"remote", "local"}},
			},
			expectedResult: "deploymentMode IN ('remote', 'local')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertNamedQueryToFilterQuery(tt.fieldFilters)
			if result != tt.expectedResult {
				t.Errorf("convertNamedQueryToFilterQuery() = %q, want %q", result, tt.expectedResult)
			}
		})
	}
}

func TestConvertFieldFilterToCondition(t *testing.T) {
	tests := []struct {
		name           string
		fieldName      string
		filter         model.FieldFilter
		expectedResult string
	}{
		{
			name:           "string equality with = operator",
			fieldName:      "provider",
			filter:         model.FieldFilter{Operator: "=", Value: "anthropic"},
			expectedResult: "provider = 'anthropic'",
		},
		{
			name:           "string equality with EQUALS operator",
			fieldName:      "provider",
			filter:         model.FieldFilter{Operator: "EQUALS", Value: "openai"},
			expectedResult: "provider = 'openai'",
		},
		{
			name:           "string not equals with != operator",
			fieldName:      "status",
			filter:         model.FieldFilter{Operator: "!=", Value: "deprecated"},
			expectedResult: "status != 'deprecated'",
		},
		{
			name:           "string not equals with NOT_EQUALS operator",
			fieldName:      "status",
			filter:         model.FieldFilter{Operator: "NOT_EQUALS", Value: "deprecated"},
			expectedResult: "status != 'deprecated'",
		},
		{
			name:           "string not equals with <> operator",
			fieldName:      "status",
			filter:         model.FieldFilter{Operator: "<>", Value: "deprecated"},
			expectedResult: "status != 'deprecated'",
		},
		{
			name:           "boolean true",
			fieldName:      "verifiedSource",
			filter:         model.FieldFilter{Operator: "=", Value: true},
			expectedResult: "verifiedSource = true",
		},
		{
			name:           "boolean false",
			fieldName:      "verifiedSource",
			filter:         model.FieldFilter{Operator: "=", Value: false},
			expectedResult: "verifiedSource = false",
		},
		{
			name:           "greater than integer",
			fieldName:      "toolCount",
			filter:         model.FieldFilter{Operator: ">", Value: float64(10)},
			expectedResult: "toolCount > 10",
		},
		{
			name:           "greater than or equal",
			fieldName:      "toolCount",
			filter:         model.FieldFilter{Operator: ">=", Value: float64(5)},
			expectedResult: "toolCount >= 5",
		},
		{
			name:           "less than",
			fieldName:      "latency",
			filter:         model.FieldFilter{Operator: "<", Value: float64(100)},
			expectedResult: "latency < 100",
		},
		{
			name:           "less than or equal",
			fieldName:      "latency",
			filter:         model.FieldFilter{Operator: "<=", Value: float64(50)},
			expectedResult: "latency <= 50",
		},
		{
			name:           "float value",
			fieldName:      "score",
			filter:         model.FieldFilter{Operator: ">=", Value: float64(0.95)},
			expectedResult: "score >= 0.950000",
		},
		{
			name:      "IN operator with interface array",
			fieldName: "tags",
			filter: model.FieldFilter{
				Operator: "IN",
				Value:    []interface{}{"security", "verified", "production"},
			},
			expectedResult: "tags IN ('security', 'verified', 'production')",
		},
		{
			name:      "IN operator with string array",
			fieldName: "transports",
			filter: model.FieldFilter{
				Operator: "IN",
				Value:    []string{"http", "sse"},
			},
			expectedResult: "transports IN ('http', 'sse')",
		},
		{
			name:           "LIKE operator",
			fieldName:      "name",
			filter:         model.FieldFilter{Operator: "LIKE", Value: "%github%"},
			expectedResult: "name LIKE '%github%'",
		},
		{
			name:           "ILIKE operator (case insensitive)",
			fieldName:      "description",
			filter:         model.FieldFilter{Operator: "ILIKE", Value: "%ai%"},
			expectedResult: "description ILIKE '%ai%'",
		},
		{
			name:           "default to equality for unknown operator",
			fieldName:      "field",
			filter:         model.FieldFilter{Operator: "UNKNOWN", Value: "value"},
			expectedResult: "field = 'value'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertFieldFilterToCondition(tt.fieldName, tt.filter)
			if result != tt.expectedResult {
				t.Errorf("convertFieldFilterToCondition() = %q, want %q", result, tt.expectedResult)
			}
		})
	}
}

func TestEscapeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no special characters",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "single quote",
			input:    "it's a test",
			expected: "it''s a test",
		},
		{
			name:     "multiple single quotes",
			input:    "it's Bob's test",
			expected: "it''s Bob''s test",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeString(tt.input)
			if result != tt.expected {
				t.Errorf("escapeString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatLicenseDisplayName(t *testing.T) {
	tests := []struct {
		name         string
		spdxLicense  string
		expectedName string
	}{
		{
			name:         "apache-2.0",
			spdxLicense:  "apache-2.0",
			expectedName: "Apache 2.0",
		},
		{
			name:         "mit",
			spdxLicense:  "mit",
			expectedName: "MIT",
		},
		{
			name:         "gpl-3.0",
			spdxLicense:  "gpl-3.0",
			expectedName: "GPL 3.0",
		},
		{
			name:         "unknown license",
			spdxLicense:  "custom-license",
			expectedName: "custom-license",
		},
		{
			name:         "llama-3.1",
			spdxLicense:  "llama-3.1",
			expectedName: "Llama 3.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatLicenseDisplayName(tt.spdxLicense)
			if result != tt.expectedName {
				t.Errorf("formatLicenseDisplayName() = %q, want %q", result, tt.expectedName)
			}
		})
	}
}

func TestTransformLicenseInFilterQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no license filter",
			input:    "provider='anthropic'",
			expected: "provider='anthropic'",
		},
		{
			name:     "single license display name",
			input:    "license='Apache 2.0'",
			expected: "license='apache-2.0'",
		},
		{
			name:     "MIT license",
			input:    "license='MIT'",
			expected: "license='mit'",
		},
		{
			name:     "license in complex query",
			input:    "provider='anthropic' AND license='Apache 2.0' AND tags='ai'",
			expected: "provider='anthropic' AND license='apache-2.0' AND tags='ai'",
		},
		{
			name:     "license IN clause",
			input:    "license IN ('Apache 2.0','MIT')",
			expected: "license IN ('apache-2.0','mit')",
		},
		{
			name:     "unknown license unchanged",
			input:    "license='Custom License'",
			expected: "license='Custom License'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformLicenseInFilterQuery(tt.input)
			if result != tt.expected {
				t.Errorf("transformLicenseInFilterQuery() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestMapKeysToInterface(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]struct{}
		expected []interface{}
	}{
		{
			name: "multiple keys",
			input: map[string]struct{}{
				"apple":  {},
				"banana": {},
				"cherry": {},
			},
			expected: []interface{}{"apple", "banana", "cherry"}, // sorted
		},
		{
			name:     "empty map",
			input:    map[string]struct{}{},
			expected: []interface{}{},
		},
		{
			name: "single key",
			input: map[string]struct{}{
				"only": {},
			},
			expected: []interface{}{"only"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapKeysToInterface(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("mapKeysToInterface() returned %d items, want %d", len(result), len(tt.expected))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("mapKeysToInterface()[%d] = %v, want %v", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestConvertStringToTransport(t *testing.T) {
	tests := []struct {
		name      string
		transport string
		expected  model.McpTransportType
	}{
		{
			name:      "http",
			transport: "http",
			expected:  model.MCPTRANSPORTTYPE_HTTP,
		},
		{
			name:      "sse",
			transport: "sse",
			expected:  model.MCPTRANSPORTTYPE_SSE,
		},
		{
			name:      "stdio",
			transport: "stdio",
			expected:  model.MCPTRANSPORTTYPE_STDIO,
		},
		{
			name:      "unknown defaults to stdio",
			transport: "websocket",
			expected:  model.MCPTRANSPORTTYPE_STDIO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertStringToTransport(tt.transport)
			if result != tt.expected {
				t.Errorf("convertStringToTransport() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertStringToAccessType(t *testing.T) {
	tests := []struct {
		name       string
		accessType string
		expected   model.McpToolAccessType
	}{
		{
			name:       "read_only",
			accessType: "read_only",
			expected:   model.MCPTOOLACCESSTYPE_READ_ONLY,
		},
		{
			name:       "read_write",
			accessType: "read_write",
			expected:   model.MCPTOOLACCESSTYPE_READ_WRITE,
		},
		{
			name:       "execute",
			accessType: "execute",
			expected:   model.MCPTOOLACCESSTYPE_EXECUTE,
		},
		{
			name:       "unknown defaults to read_only",
			accessType: "unknown",
			expected:   model.MCPTOOLACCESSTYPE_READ_ONLY,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertStringToAccessType(tt.accessType)
			if result != tt.expected {
				t.Errorf("convertStringToAccessType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertStringToDeploymentMode(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		expected model.McpDeploymentMode
	}{
		{
			name:     "remote",
			mode:     "remote",
			expected: model.MCPDEPLOYMENTMODE_REMOTE,
		},
		{
			name:     "local",
			mode:     "local",
			expected: model.MCPDEPLOYMENTMODE_LOCAL,
		},
		{
			name:     "unknown defaults to local",
			mode:     "hybrid",
			expected: model.MCPDEPLOYMENTMODE_LOCAL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertStringToDeploymentMode(tt.mode)
			if result != tt.expected {
				t.Errorf("convertStringToDeploymentMode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNamedQueryResolver(t *testing.T) {
	// Test that SetNamedQueryResolver and the resolver work correctly
	provider := NewDbMcpCatalogProvider(nil) // nil repository is fine for this test

	if provider.namedQueryResolver != nil {
		t.Error("namedQueryResolver should be nil initially")
	}

	// Set a resolver
	testQueries := map[string]map[string]model.FieldFilter{
		"production_ready": {
			"verifiedSource": {Operator: "=", Value: true},
			"provider":       {Operator: "IN", Value: []interface{}{"anthropic", "openai"}},
		},
	}

	provider.SetNamedQueryResolver(func() map[string]map[string]model.FieldFilter {
		return testQueries
	})

	if provider.namedQueryResolver == nil {
		t.Error("namedQueryResolver should be set after SetNamedQueryResolver")
	}

	// Verify the resolver returns the expected queries
	result := provider.namedQueryResolver()
	if len(result) != 1 {
		t.Errorf("namedQueryResolver() returned %d queries, want 1", len(result))
	}

	if _, exists := result["production_ready"]; !exists {
		t.Error("namedQueryResolver() should contain 'production_ready' query")
	}
}

func TestMultipleConditionsJoin(t *testing.T) {
	// Test that multiple field filters are joined with AND
	fieldFilters := map[string]model.FieldFilter{
		"verifiedSource": {Operator: "=", Value: true},
	}

	result := convertNamedQueryToFilterQuery(fieldFilters)

	// With a single filter, there should be no AND
	if result == "" {
		t.Error("convertNamedQueryToFilterQuery() returned empty string for non-empty filters")
	}

	// The result should contain the condition
	if result != "verifiedSource = true" {
		t.Errorf("convertNamedQueryToFilterQuery() = %q, want %q", result, "verifiedSource = true")
	}
}
