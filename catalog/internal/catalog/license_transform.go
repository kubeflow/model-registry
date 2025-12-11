package catalog

import (
	"strings"
)

// Generate file for spdxToHumanReadableMap
//go:generate ../../../scripts/gen_license_names.sh

var spdxOverrides = map[string]string{
	"apache-2.0":            "Apache 2.0",
	"mit":                   "MIT",
	"bsd-3-clause":          "BSD 3-Clause",
	"bsd-2-clause":          "BSD 2-Clause",
	"gpl-3.0":               "GPL 3.0",
	"gpl-2.0":               "GPL 2.0",
	"lgpl-3.0":              "LGPL 3.0",
	"lgpl-2.1":              "LGPL 2.1",
	"llama2":                "Llama 2 Community License",
	"llama3":                "Llama 3 Community License",
	"llama3.1":              "Llama 3.1 Community License",
	"llama3.2":              "Llama 3.2 Community License",
	"llama3.3":              "Llama 3.3 Community License",
	"llama4":                "Llama 4 Community License",
	"bigscience-openrail-m": "BigScience OpenRAIL-M License",
	"openrail":              "OpenRAIL License",
	"gemma":                 "Gemma License",
}

// transformLicenseToHumanReadable converts license identifiers to human-readable names
func transformLicenseToHumanReadable(license string) string {
	license = strings.TrimSpace(license)
	if license == "" {
		return ""
	}

	normalized := strings.ToLower(license)

	// Check if we have a direct mapping
	if humanReadable, exists := spdxOverrides[normalized]; exists {
		return humanReadable
	}

	// Try the SPDX license name
	if humanReadable, exists := spdxToHumanReadableMap[normalized]; exists {
		return humanReadable
	}

	// Fallback to the name we were given
	return license
}
