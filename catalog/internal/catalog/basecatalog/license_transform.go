package basecatalog

import (
	"fmt"
	"regexp"
	"strings"
)

// Generate file for spdxToHumanReadableMap
//go:generate ../../../../scripts/gen_license_names.sh

var spdxOverrides = map[string]string{
	"apache-2.0":                "Apache 2.0",
	"bigscience-openrail-m":     "BigScience OpenRAIL-M License",
	"bsd-2-clause":              "BSD 2-Clause",
	"bsd-3-clause":              "BSD 3-Clause",
	"gemma":                     "Gemma License",
	"gpl-2.0":                   "GPL 2.0",
	"gpl-3.0":                   "GPL 3.0",
	"lgpl-2.1":                  "LGPL 2.1",
	"lgpl-3.0":                  "LGPL 3.0",
	"llama2":                    "Llama 2 Community License",
	"llama3.1":                  "Llama 3.1 Community License",
	"llama3.2":                  "Llama 3.2 Community License",
	"llama3.3":                  "Llama 3.3 Community License",
	"llama3":                    "Llama 3 Community License",
	"llama4":                    "Llama 4 Community License",
	"mit":                       "MIT",
	"nvidia-open-model-license": "NVIDIA Open Model License",
	"openrail":                  "OpenRAIL License",
}

// displayNameToSpdx is the reverse mapping from display names to SPDX identifiers.
// Built at init time from both spdxToHumanReadableMap and spdxOverrides.
// Keys are normalized to lowercase for case-insensitive matching.
var displayNameToSpdx map[string]string

// spdxToDisplay is the forward mapping from SPDX identifiers to display names.
// Built at init time and cached to avoid repeated allocations in hot paths.
var spdxToDisplay map[string]string

// singleQuotePatterns contains pre-compiled regex patterns for single-quoted license names.
// Built at init time to avoid regex compilation on hot path (filter queries).
var singleQuotePatterns map[string]*regexp.Regexp

// doubleQuotePatterns contains pre-compiled regex patterns for double-quoted license names.
// Built at init time to avoid regex compilation on hot path (filter queries).
var doubleQuotePatterns map[string]*regexp.Regexp

func init() {
	// Build reverse map: display name → SPDX identifier (case-insensitive)
	displayNameToSpdx = make(map[string]string)

	// Build forward map: SPDX → display name (cached for performance)
	spdxToDisplay = make(map[string]string)

	// Add entries from auto-generated SPDX map
	for spdx, display := range spdxToHumanReadableMap {
		// Store lowercase versions for case-insensitive matching
		displayNameToSpdx[strings.ToLower(display)] = spdx
		spdxToDisplay[spdx] = display
	}

	// Add entries from overrides (these take precedence)
	for spdx, display := range spdxOverrides {
		displayNameToSpdx[strings.ToLower(display)] = spdx
		spdxToDisplay[spdx] = display
	}

	// Pre-compile regex patterns for license transformation to avoid compilation overhead
	// on hot path (filter queries). This provides ~100x performance improvement.
	singleQuotePatterns = make(map[string]*regexp.Regexp, len(spdxToDisplay))
	doubleQuotePatterns = make(map[string]*regexp.Regexp, len(spdxToDisplay))

	for spdx, display := range spdxToDisplay {
		escapedDisplay := regexp.QuoteMeta(display)
		singleQuotePatterns[spdx] = regexp.MustCompile(fmt.Sprintf(`'%s'`, escapedDisplay))
		doubleQuotePatterns[spdx] = regexp.MustCompile(fmt.Sprintf(`"%s"`, escapedDisplay))
	}
}

// TransformLicenseToHumanReadable converts license identifiers to human-readable names
func TransformLicenseToHumanReadable(license string) string {
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

// TransformDisplayNameToSpdx converts a human-readable display name back to SPDX identifier.
// This is useful for filter query transformation, allowing users to filter by display name.
// Returns the original value if no mapping exists (could already be an SPDX ID).
func TransformDisplayNameToSpdx(displayName string) string {
	displayName = strings.TrimSpace(displayName)
	if displayName == "" {
		return ""
	}

	normalized := strings.ToLower(displayName)
	if spdx, exists := displayNameToSpdx[normalized]; exists {
		return spdx
	}

	// Could already be an SPDX ID, return as-is
	return displayName
}

// TransformLicenseNamesInFilterQuery transforms license display names to SPDX identifiers
// in filter query strings. This enables users to filter by either human-readable names or
// SPDX IDs across all AI asset types (models, MCP servers, future assets).
//
// The function handles both single and double quotes to ensure compatibility with various
// query formats and works with IN clauses for multiple license filtering. Uses pre-compiled
// regex patterns for optimal performance on the hot path (filter queries).
//
// Examples:
//   - license='Apache 2.0' → license='apache-2.0'
//   - license="MIT" → license="mit"
//   - license IN ['Apache 2.0', 'MIT'] → license IN ['apache-2.0', 'mit']
//   - license='apache-2.0' → license='apache-2.0' (already SPDX, unchanged)
//
// This function is used by:
//   - Model catalog (catalog/internal/catalog/modelcatalog/db_catalog.go)
//   - MCP catalog (catalog/internal/catalog/mcpcatalog/db_mcp.go)
//   - Future AI asset catalogs
func TransformLicenseNamesInFilterQuery(filterQuery string) string {
	// Early exit for empty queries or queries without license field
	if filterQuery == "" || !strings.Contains(filterQuery, "license") {
		return filterQuery
	}

	result := filterQuery

	// Use pre-compiled regex patterns to avoid compilation overhead on hot path
	// This provides ~100x performance improvement over compiling patterns on each call
	for spdx := range spdxToDisplay {
		// Single quotes: license='Apache 2.0' → license='apache-2.0'
		result = singleQuotePatterns[spdx].ReplaceAllString(result, "'"+spdx+"'")

		// Double quotes: license="Apache 2.0" → license="apache-2.0"
		result = doubleQuotePatterns[spdx].ReplaceAllString(result, `"`+spdx+`"`)
	}

	return result
}
