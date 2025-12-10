package catalog

import (
	"testing"
)

func TestTransformLicenseToHumanReadable(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Empty and edge cases
		{
			name:     "empty license",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: "",
		},

		// Standard Open Source Licenses
		{
			name:     "apache-2.0",
			input:    "apache-2.0",
			expected: "Apache 2.0",
		},
		{
			name:     "mit",
			input:    "mit",
			expected: "MIT",
		},
		{
			name:     "bsd-3-clause",
			input:    "bsd-3-clause",
			expected: "BSD 3-Clause",
		},
		{
			name:     "bsd-2-clause",
			input:    "bsd-2-clause",
			expected: "BSD 2-Clause",
		},
		{
			name:     "gpl-3.0",
			input:    "gpl-3.0",
			expected: "GPL 3.0",
		},
		{
			name:     "gpl-2.0",
			input:    "gpl-2.0",
			expected: "GPL 2.0",
		},
		{
			name:     "lgpl-3.0",
			input:    "lgpl-3.0",
			expected: "LGPL 3.0",
		},
		{
			name:     "lgpl-2.1",
			input:    "lgpl-2.1",
			expected: "LGPL 2.1",
		},

		// Creative Commons
		{
			name:     "cc-by-4.0",
			input:    "cc-by-4.0",
			expected: "Creative Commons Attribution 4.0 International",
		},
		{
			name:     "cc-by-sa-4.0",
			input:    "cc-by-sa-4.0",
			expected: "Creative Commons Attribution Share Alike 4.0 International",
		},
		{
			name:     "cc-by-nc-4.0",
			input:    "cc-by-nc-4.0",
			expected: "Creative Commons Attribution Non Commercial 4.0 International",
		},
		{
			name:     "cc0-1.0",
			input:    "cc0-1.0",
			expected: "Creative Commons Zero v1.0 Universal",
		},

		// Special/Model-Specific Licenses
		{
			name:     "unlicense",
			input:    "unlicense",
			expected: "The Unlicense",
		},
		{
			name:     "llama2",
			input:    "llama2",
			expected: "Llama 2 Community License",
		},
		{
			name:     "llama3",
			input:    "llama3",
			expected: "Llama 3 Community License",
		},
		{
			name:     "llama3.1",
			input:    "llama3.1",
			expected: "Llama 3.1 Community License",
		},
		{
			name:     "llama3.2",
			input:    "llama3.2",
			expected: "Llama 3.2 Community License",
		},
		{
			name:     "llama3.3",
			input:    "llama3.3",
			expected: "Llama 3.3 Community License",
		},
		{
			name:     "llama4",
			input:    "llama4",
			expected: "Llama 4 Community License",
		},
		{
			name:     "bigscience-openrail-m",
			input:    "bigscience-openrail-m",
			expected: "BigScience OpenRAIL-M License",
		},
		{
			name:     "openrail",
			input:    "openrail",
			expected: "OpenRAIL License",
		},
		{
			name:     "gemma",
			input:    "gemma",
			expected: "Gemma License",
		},

		// Case insensitive matching
		{
			name:     "APACHE-2.0 uppercase",
			input:    "APACHE-2.0",
			expected: "Apache 2.0",
		},
		{
			name:     "Mit mixed case",
			input:    "Mit",
			expected: "MIT",
		},
		{
			name:     "GPL-3.0 mixed case",
			input:    "GPL-3.0",
			expected: "GPL 3.0",
		},

		// Whitespace handling
		{
			name:     "apache-2.0 with whitespace",
			input:    "  apache-2.0  ",
			expected: "Apache 2.0",
		},

		// Fallback on unknown licenses
		{
			name:     "unknown license",
			input:    "custom-license-1.0",
			expected: "custom-license-1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformLicenseToHumanReadable(tt.input)
			if result != tt.expected {
				t.Errorf("transformLicenseToHumanReadable(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}
