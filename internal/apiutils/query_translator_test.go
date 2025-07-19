package apiutils

import (
	"testing"
)

func TestTranslateFilterQuery(t *testing.T) {
	tests := []struct {
		name         string
		restQuery    string
		entityType   EntityType
		expectedMLMD string
		expectError  bool
	}{
		{
			name:         "Empty query",
			restQuery:    "",
			entityType:   RegisteredModelEntity,
			expectedMLMD: "",
			expectError:  false,
		},
		{
			name:         "Core MLMD field - RegisteredModel",
			restQuery:    "name = 'my-model'",
			entityType:   RegisteredModelEntity,
			expectedMLMD: "name = \"my-model\"",
			expectError:  false,
		},
		{
			name:         "REST property - RegisteredModel state",
			restQuery:    "state = 'LIVE'",
			entityType:   RegisteredModelEntity,
			expectedMLMD: "properties.state.string_value = \"LIVE\"",
			expectError:  false,
		},
		{
			name:         "REST property - ModelVersion author",
			restQuery:    "author = 'john.doe'",
			entityType:   ModelVersionEntity,
			expectedMLMD: "properties.author.string_value = \"john.doe\"",
			expectError:  false,
		},
		{
			name:         "Custom property with explicit type",
			restQuery:    "accuracy.double_value > 0.95",
			entityType:   ModelVersionEntity,
			expectedMLMD: "custom_properties.accuracy.double_value > 0.95",
			expectError:  false,
		},
		{
			name:         "Custom property with escaped dot in name",
			restQuery:    "model\\.version.string_value = 'v1.0'",
			entityType:   ModelVersionEntity,
			expectedMLMD: "custom_properties.`model.version`.string_value = \"v1.0\"",
			expectError:  false,
		},
		{
			name:         "Custom property with multiple escaped dots",
			restQuery:    "my\\.complex\\.property.int_value = 42",
			entityType:   ExperimentEntity,
			expectedMLMD: "custom_properties.`my.complex.property`.int_value = 42",
			expectError:  false,
		},
		{
			name:         "Custom property with escaped dot but no explicit type",
			restQuery:    "model\\.version = 'v1.0'",
			entityType:   ModelVersionEntity,
			expectedMLMD: "custom_properties.`model.version`.string_value = \"v1.0\"",
			expectError:  false,
		},
		{
			name:         "Custom property with colon (requires quoting)",
			restQuery:    "namespace:key.string_value = 'value'",
			entityType:   ExperimentEntity,
			expectedMLMD: "custom_properties.`namespace:key`.string_value = \"value\"",
			expectError:  false,
		},
		{
			name:         "Custom property without type (assumed string)",
			restQuery:    "custom_field = 'value'",
			entityType:   ExperimentEntity,
			expectedMLMD: "custom_properties.custom_field.string_value = \"value\"",
			expectError:  false,
		},
		{
			name:         "Complex query with AND",
			restQuery:    "state = 'LIVE' AND description LIKE '%test%'",
			entityType:   RegisteredModelEntity,
			expectedMLMD: "properties.state.string_value = \"LIVE\" AND properties.description.string_value LIKE \"%test%\"",
			expectError:  false,
		},
		{
			name:         "Query with parentheses",
			restQuery:    "(state = 'LIVE' OR state = 'ARCHIVED') AND name LIKE 'model%'",
			entityType:   RegisteredModelEntity,
			expectedMLMD: "(properties.state.string_value = \"LIVE\" OR properties.state.string_value = \"ARCHIVED\") AND name LIKE \"model%\"",
			expectError:  false,
		},
		{
			name:         "Invalid custom property value type",
			restQuery:    "field.invalid_type = \"value\"",
			entityType:   RegisteredModelEntity,
			expectedMLMD: "",
			expectError:  true,
		},
		{
			name:         "Unmatched parentheses",
			restQuery:    "state = 'LIVE' AND (description LIKE '%test%'",
			entityType:   RegisteredModelEntity,
			expectedMLMD: "",
			expectError:  true,
		},
		{
			name:         "Real world query with hyphens in value",
			restQuery:    "name ='e2e-test-dbee6a77'",
			entityType:   RegisteredModelEntity,
			expectedMLMD: "name =\"e2e-test-dbee6a77\"",
			expectError:  false,
		},
		{
			name:         "Query with double quotes",
			restQuery:    `name ="e2e-test-dbee6a77"`,
			entityType:   RegisteredModelEntity,
			expectedMLMD: `name ="e2e-test-dbee6a77"`,
			expectError:  false,
		},
		{
			name:         "Query with no space around equals",
			restQuery:    "name='e2e-test-dbee6a77'",
			entityType:   RegisteredModelEntity,
			expectedMLMD: "name=\"e2e-test-dbee6a77\"",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TranslateFilterQuery(tt.restQuery, tt.entityType)

			// Debug output for the problematic query
			if tt.name == "Real world query with hyphens in value" || tt.name == "Query with double quotes" || tt.name == "Query with no space around equals" {
				t.Logf("Input: %q", tt.restQuery)
				t.Logf("Output: %q", result)
				if err != nil {
					t.Logf("Error: %v", err)
				}
			}

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError && result != tt.expectedMLMD {
				t.Errorf("Expected '%s', got '%s'", tt.expectedMLMD, result)
			}
		})
	}
}

func TestValidateQuerySyntax(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		expectError bool
	}{
		{
			name:        "Valid query with balanced parentheses",
			query:       "(state = 'LIVE' AND (description LIKE '%test%'))",
			expectError: false,
		},
		{
			name:        "Unmatched opening parenthesis",
			query:       "state = 'LIVE' AND (description LIKE '%test%'",
			expectError: true,
		},
		{
			name:        "Unmatched closing parenthesis",
			query:       "state = 'LIVE') AND description LIKE '%test%'",
			expectError: true,
		},
		{
			name:        "Multiple unmatched opening parentheses",
			query:       "((state = 'LIVE'",
			expectError: true,
		},
		{
			name:        "No parentheses",
			query:       "state = 'LIVE' AND description = 'test'",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateQuerySyntax(tt.query)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestIsValidValueType(t *testing.T) {
	tests := []struct {
		valueType string
		expected  bool
	}{
		{"string_value", true},
		{"int_value", true},
		{"double_value", true},
		{"bool_value", true},
		{"struct_value", false},
		{"proto_value", false},
		{"invalid_type", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.valueType, func(t *testing.T) {
			result := isValidValueType(tt.valueType)
			if result != tt.expected {
				t.Errorf("Expected %v for '%s', got %v", tt.expected, tt.valueType, result)
			}
		})
	}
}

func TestNormalizeQuotesForMLMD(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple single quotes",
			input:    "name = 'test'",
			expected: "name = \"test\"",
		},
		{
			name:     "escaped single quote within single quotes",
			input:    "name = 'John\\'s Model'",
			expected: "name = \"John's Model\"",
		},
		{
			name:     "double quotes within single quotes",
			input:    "name = 'Best \"ML\" Model'",
			expected: "name = \"Best \\\"ML\\\" Model\"",
		},
		{
			name:     "complex escaped content",
			input:    "name = 'John\\'s \"Best\" Model'",
			expected: "name = \"John's \\\"Best\\\" Model\"",
		},
		{
			name:     "double quotes should remain unchanged",
			input:    "name = \"test\"",
			expected: "name = \"test\"",
		},
		{
			name:     "mixed single and double quotes",
			input:    "name = 'John\\'s' AND description = \"Best Model\"",
			expected: "name = \"John's\" AND description = \"Best Model\"",
		},
		{
			name:     "escaped backslash in single quotes",
			input:    "path = 'C:\\\\models'",
			expected: "path = \"C:\\\\models\"",
		},
		{
			name:     "escaped backslash and quotes",
			input:    "name = 'Model\\\\with\\'quotes'",
			expected: "name = \"Model\\\\with'quotes\"",
		},
		{
			name:     "escaped single quote at end of string",
			input:    "name = 'John\\'s'",
			expected: "name = \"John's\"",
		},
		{
			name:     "empty single quoted string",
			input:    "name = ''",
			expected: "name = \"\"",
		},
		{
			name:     "multiple single quoted strings",
			input:    "name = 'test1' AND description = 'test2'",
			expected: "name = \"test1\" AND description = \"test2\"",
		},
		{
			name:     "escaped double quote in double quoted string",
			input:    "name = \"John\\\"s Model\"",
			expected: "name = \"John\\\"s Model\"",
		},
		{
			name:     "single quote in double quoted string",
			input:    "name = \"John's Model\"",
			expected: "name = \"John's Model\"",
		},
		{
			name:     "no quotes",
			input:    "id = 123 AND active = true",
			expected: "id = 123 AND active = true",
		},
		{
			name:     "complex query with various quote scenarios",
			input:    "name = 'John\\'s \"AI\" Model' AND description = \"Best Model\" AND path = 'C:\\\\data'",
			expected: "name = \"John's \\\"AI\\\" Model\" AND description = \"Best Model\" AND path = \"C:\\\\data\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeQuotesForMLMD(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeQuotesForMLMD() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestFindLastUnescapedDot(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"model.version", 5},                 // simple dot
		{"model\\.version", -1},              // escaped dot
		{"model.version.string_value", 13},   // multiple dots, last one
		{"model\\.version.string_value", 14}, // escaped dot, then real dot (corrected)
		{"model\\.version\\.name", -1},       // all dots escaped
		{"no_dots", -1},                      // no dots
		{"", -1},                             // empty string
		{"model\\\\\\..version", 9},          // escaped backslash followed by dot (corrected)
		{"model\\\\.version", 7},             // double escaped backslash, then dot (corrected)
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := findLastUnescapedDot(tt.input)
			if result != tt.expected {
				t.Errorf("findLastUnescapedDot(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnescapeDots(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"model\\.version", "model.version"},
		{"my\\.complex\\.property", "my.complex.property"},
		{"no_escaped_dots", "no_escaped_dots"},
		{"model\\\\version", "model\\\\version"}, // escaped backslash, not escaped dot
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := unescapeDots(tt.input)
			if result != tt.expected {
				t.Errorf("unescapeDots(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestQuotePropertyNameIfNeeded(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Valid characters [0-9A-Za-z_] - no quoting needed
		{"simple", "simple"},                       // lowercase letters
		{"Simple", "Simple"},                       // uppercase letters
		{"simple123", "simple123"},                 // letters and numbers
		{"Simple_123", "Simple_123"},               // letters, numbers, underscore
		{"already_good", "already_good"},           // underscore is valid
		{"CamelCase123", "CamelCase123"},           // mixed case alphanumeric
		{"_underscore_start", "_underscore_start"}, // underscore at start
		{"123numeric", "123numeric"},               // numbers at start
		{"", ""},                                   // empty string (edge case)

		// Invalid characters - require quoting
		{"model.version", "`model.version`"},     // dot requires quoting
		{"namespace:key", "`namespace:key`"},     // colon requires quoting
		{"my key", "`my key`"},                   // space requires quoting
		{"my-property", "`my-property`"},         // hyphen requires quoting
		{"path/to/file", "`path/to/file`"},       // slash requires quoting
		{"path\\to\\file", "`path\\to\\file`"},   // backslash requires quoting
		{"property@domain", "`property@domain`"}, // @ requires quoting
		{"property+suffix", "`property+suffix`"}, // + requires quoting
		{"property$var", "`property$var`"},       // $ requires quoting
		{"property#tag", "`property#tag`"},       // # requires quoting
		{"property%val", "`property%val`"},       // % requires quoting
		{"property(param)", "`property(param)`"}, // parentheses require quoting
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := quotePropertyNameIfNeeded(tt.input)
			if result != tt.expected {
				t.Errorf("quotePropertyNameIfNeeded(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseCustomPropertyField(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedProperty  string
		expectedValueType string
		expectedIsCustom  bool
	}{
		{
			name:              "simple custom property",
			input:             "accuracy.double_value",
			expectedProperty:  "accuracy",
			expectedValueType: "double_value",
			expectedIsCustom:  true,
		},
		{
			name:              "custom property with escaped dot",
			input:             "model\\.version.string_value",
			expectedProperty:  "model.version",
			expectedValueType: "string_value",
			expectedIsCustom:  true,
		},
		{
			name:              "multiple escaped dots",
			input:             "my\\.complex\\.property.int_value",
			expectedProperty:  "my.complex.property",
			expectedValueType: "int_value",
			expectedIsCustom:  true,
		},
		{
			name:              "no dot - not custom property",
			input:             "simple_field",
			expectedProperty:  "",
			expectedValueType: "",
			expectedIsCustom:  false,
		},
		{
			name:              "invalid value type",
			input:             "field.invalid_type",
			expectedProperty:  "",
			expectedValueType: "",
			expectedIsCustom:  false,
		},
		{
			name:              "all dots escaped - not custom property",
			input:             "model\\.version\\.name",
			expectedProperty:  "",
			expectedValueType: "",
			expectedIsCustom:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			property, valueType, isCustom := parseCustomPropertyField(tt.input)
			if property != tt.expectedProperty {
				t.Errorf("parseCustomPropertyField(%q) property = %q, expected %q", tt.input, property, tt.expectedProperty)
			}
			if valueType != tt.expectedValueType {
				t.Errorf("parseCustomPropertyField(%q) valueType = %q, expected %q", tt.input, valueType, tt.expectedValueType)
			}
			if isCustom != tt.expectedIsCustom {
				t.Errorf("parseCustomPropertyField(%q) isCustom = %v, expected %v", tt.input, isCustom, tt.expectedIsCustom)
			}
		})
	}
}
