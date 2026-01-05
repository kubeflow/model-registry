package catalog

import (
	"github.com/kubeflow/model-registry/catalog/internal/mcp"
)

// mergeMcpSources performs field-level merging of two McpSource structs.
// Fields from 'override' take precedence over 'base' when they are explicitly set.
// A field is considered "set" if:
// - For strings: non-empty
// - For pointers: non-nil
// - For slices: non-nil (empty slice is considered explicitly set to "no items")
// - For maps: non-nil (empty map is considered explicitly set)
//
// This follows the same pattern as mergeSources() for model catalog sources.
func mergeMcpSources(base, override mcp.McpSource) mcp.McpSource {
	result := base

	// Id is always taken from override (it's the key)
	result.Id = override.Id

	// Name: override if non-empty
	if override.Name != "" {
		result.Name = override.Name
	}

	// Enabled: override if non-nil
	if override.Enabled != nil {
		result.Enabled = override.Enabled
	}

	// Labels: override if non-nil (empty slice means "explicitly no labels")
	if override.Labels != nil {
		result.Labels = override.Labels
	}

	// IncludedServers: override if non-nil (empty slice means "no inclusions")
	if override.IncludedServers != nil {
		result.IncludedServers = override.IncludedServers
	}

	// ExcludedServers: override if non-nil (empty slice means "no exclusions")
	if override.ExcludedServers != nil {
		result.ExcludedServers = override.ExcludedServers
	}

	// Type: override if non-empty
	if override.Type != "" {
		result.Type = override.Type
	}

	// Properties: override if non-nil (complete replacement, not deep merge)
	if override.Properties != nil {
		result.Properties = override.Properties
	}

	// Origin: use override's origin if Properties are overridden (since relative
	// paths in Properties should resolve relative to where they were defined).
	// Otherwise, keep base origin (where Type and original Properties came from).
	if override.Properties != nil && override.Origin != "" {
		result.Origin = override.Origin
	}

	return result
}

// applyMcpSourceDefaults applies default values to an McpSource for fields that are not set.
func applyMcpSourceDefaults(source mcp.McpSource) mcp.McpSource {
	// Default Enabled to true if not set
	if source.Enabled == nil {
		enabled := true
		source.Enabled = &enabled
	}

	// Default Labels to empty slice if not set
	if source.Labels == nil {
		source.Labels = []string{}
	}

	return source
}

// MergeMcpSourcesFromPaths reads MCP source configurations from multiple paths
// and merges sources with the same ID using field-level merging.
// Sources from later paths override fields from earlier paths (priority order).
func MergeMcpSourcesFromPaths(paths []string, readFunc func(path string) ([]mcp.McpSource, error)) (map[string]mcp.McpSource, error) {
	// Collect sources by origin in priority order
	type originEntry struct {
		origin  string
		sources map[string]mcp.McpSource
	}

	entries := make([]originEntry, 0, len(paths))

	for _, path := range paths {
		sources, err := readFunc(path)
		if err != nil {
			// Return error for invalid config files (caller handles warnings)
			continue
		}

		sourceMap := make(map[string]mcp.McpSource, len(sources))
		for _, source := range sources {
			source.Origin = path
			sourceMap[source.Id] = source
		}

		entries = append(entries, originEntry{origin: path, sources: sourceMap})
	}

	// Merge sources with field-level priority
	result := make(map[string]mcp.McpSource)

	for _, entry := range entries {
		for id, source := range entry.sources {
			if existing, ok := result[id]; ok {
				// Field-level merge: existing is base, source is override
				result[id] = mergeMcpSources(existing, source)
			} else {
				result[id] = source
			}
		}
	}

	// Apply defaults to all merged sources
	for id, source := range result {
		result[id] = applyMcpSourceDefaults(source)
	}

	return result, nil
}
