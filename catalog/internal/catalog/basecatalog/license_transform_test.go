package basecatalog

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
			result := TransformLicenseToHumanReadable(tt.input)
			if result != tt.expected {
				t.Errorf("TransformLicenseToHumanReadable(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTransformLicenseNamesInFilterQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Empty and edge cases
		{
			name:     "empty filter query",
			input:    "",
			expected: "",
		},

		// Single license with single quotes
		{
			name:     "Apache 2.0 single quotes",
			input:    "license='Apache 2.0'",
			expected: "license='apache-2.0'",
		},
		{
			name:     "MIT single quotes",
			input:    "license='MIT'",
			expected: "license='mit'",
		},
		{
			name:     "BSD 3-Clause single quotes",
			input:    "license='BSD 3-Clause'",
			expected: "license='bsd-3-clause'",
		},
		{
			name:     "GPL 3.0 single quotes",
			input:    "license='GPL 3.0'",
			expected: "license='gpl-3.0'",
		},

		// Single license with double quotes
		{
			name:     "Apache 2.0 double quotes",
			input:    `license="Apache 2.0"`,
			expected: `license="apache-2.0"`,
		},
		{
			name:     "MIT double quotes",
			input:    `license="MIT"`,
			expected: `license="mit"`,
		},

		// Already SPDX identifier (pass-through)
		{
			name:     "SPDX ID single quotes (apache-2.0)",
			input:    "license='apache-2.0'",
			expected: "license='apache-2.0'",
		},
		{
			name:     "SPDX ID single quotes (mit)",
			input:    "license='mit'",
			expected: "license='mit'",
		},

		// IN clause with multiple licenses
		{
			name:     "IN clause with display names",
			input:    "license IN ['Apache 2.0', 'MIT', 'BSD 3-Clause']",
			expected: "license IN ['apache-2.0', 'mit', 'bsd-3-clause']",
		},
		{
			name:     "IN clause with mixed quotes",
			input:    `license IN ["Apache 2.0", "MIT"]`,
			expected: `license IN ["apache-2.0", "mit"]`,
		},

		// Combined filters
		{
			name:     "license AND other field",
			input:    "license='Apache 2.0' AND verifiedSource=true",
			expected: "license='apache-2.0' AND verifiedSource=true",
		},
		{
			name:     "complex query with multiple conditions",
			input:    "(license='MIT' OR license='Apache 2.0') AND provider='Acme'",
			expected: "(license='mit' OR license='apache-2.0') AND provider='Acme'",
		},

		// ML model licenses
		{
			name:     "Llama 3.1 Community License",
			input:    "license='Llama 3.1 Community License'",
			expected: "license='llama3.1'",
		},
		{
			name:     "Gemma License",
			input:    "license='Gemma License'",
			expected: "license='gemma'",
		},
		{
			name:     "NVIDIA Open Model License",
			input:    "license='NVIDIA Open Model License'",
			expected: "license='nvidia-open-model-license'",
		},

		// Unknown license (pass-through)
		{
			name:     "unknown license name",
			input:    "license='Custom License 1.0'",
			expected: "license='Custom License 1.0'",
		},

		// Query without license field (unchanged)
		{
			name:     "query without license",
			input:    "provider='Acme' AND verifiedSource=true",
			expected: "provider='Acme' AND verifiedSource=true",
		},

		// Mixed SPDX and display names
		{
			name:     "mixed SPDX and display names",
			input:    "license IN ['apache-2.0', 'MIT', 'bsd-3-clause']",
			expected: "license IN ['apache-2.0', 'mit', 'bsd-3-clause']",
		},

		// Edge cases: substring collision safety
		{
			name:     "license name in description field should not transform",
			input:    "description='SUBMITTED by MIT team' AND license='Apache 2.0'",
			expected: "description='SUBMITTED by MIT team' AND license='apache-2.0'",
		},
		{
			name:     "license in field name should not transform",
			input:    "license_link='http://example.com' AND license='MIT'",
			expected: "license_link='http://example.com' AND license='mit'",
		},
		{
			name:     "license name as part of URL should not transform",
			input:    "url='https://mit.edu' AND license='Apache 2.0'",
			expected: "url='https://mit.edu' AND license='apache-2.0'",
		},

		// Edge cases: multiple occurrences
		{
			name:     "multiple occurrences of same license",
			input:    "(license='MIT' OR license='MIT') AND provider='Acme'",
			expected: "(license='mit' OR license='mit') AND provider='Acme'",
		},
		{
			name:     "multiple different licenses in complex query",
			input:    "(license='MIT' OR license='Apache 2.0') AND (license='BSD 3-Clause' OR license='GPL 3.0')",
			expected: "(license='mit' OR license='apache-2.0') AND (license='bsd-3-clause' OR license='gpl-3.0')",
		},

		// Edge cases: nested quotes and special characters
		{
			name:     "query with mixed quote types",
			input:    `license='Apache 2.0' AND description="MIT-style license"`,
			expected: `license='apache-2.0' AND description="MIT-style license"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TransformLicenseNamesInFilterQuery(tt.input)
			if result != tt.expected {
				t.Errorf("TransformLicenseNamesInFilterQuery(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTransformDisplayNameToSpdx(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Empty and edge cases
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: "",
		},

		// Common open source licenses
		{
			name:     "Apache 2.0 to SPDX",
			input:    "Apache 2.0",
			expected: "apache-2.0",
		},
		{
			name:     "MIT to SPDX",
			input:    "MIT",
			expected: "mit",
		},
		{
			name:     "BSD 3-Clause to SPDX",
			input:    "BSD 3-Clause",
			expected: "bsd-3-clause",
		},
		{
			name:     "GPL 3.0 to SPDX",
			input:    "GPL 3.0",
			expected: "gpl-3.0",
		},

		// Case insensitive
		{
			name:     "apache 2.0 lowercase",
			input:    "apache 2.0",
			expected: "apache-2.0",
		},
		{
			name:     "APACHE 2.0 uppercase",
			input:    "APACHE 2.0",
			expected: "apache-2.0",
		},

		// Already SPDX (pass-through)
		{
			name:     "already SPDX apache-2.0",
			input:    "apache-2.0",
			expected: "apache-2.0",
		},
		{
			name:     "already SPDX mit",
			input:    "mit",
			expected: "mit",
		},

		// Unknown display name (pass-through)
		{
			name:     "unknown display name",
			input:    "Custom License 1.0",
			expected: "Custom License 1.0",
		},

		// ML model licenses
		{
			name:     "Llama 3.1 to SPDX",
			input:    "Llama 3.1 Community License",
			expected: "llama3.1",
		},
		{
			name:     "Gemma to SPDX",
			input:    "Gemma License",
			expected: "gemma",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TransformDisplayNameToSpdx(tt.input)
			if result != tt.expected {
				t.Errorf("TransformDisplayNameToSpdx(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}
