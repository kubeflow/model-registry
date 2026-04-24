package mcpcatalog

import (
	"fmt"

	"github.com/kubeflow/hub/catalog/internal/catalog/basecatalog"
)

// ServerFilter encapsulates include/exclude pattern matching for MCP server names.
type ServerFilter = basecatalog.NameFilter

// ValidateServerFilters validates that the includedServers and excludedServers patterns
// are valid (non-empty, compilable, non-conflicting). This is useful for early validation
// at configuration load time without constructing the full ServerFilter.
func ValidateServerFilters(included, excluded []string) error {
	return basecatalog.ValidatePatterns("includedServers", included, "excludedServers", excluded)
}

// NewServerFilter builds a ServerFilter from the provided include/exclude pattern lists.
func NewServerFilter(included, excluded []string) (*ServerFilter, error) {
	return basecatalog.NewNameFilter("includedServers", included, "excludedServers", excluded)
}

// NewServerFilterFromSource composes a ServerFilter using the source-level configuration.
func NewServerFilterFromSource(source *basecatalog.MCPSource) (*ServerFilter, error) {
	if source == nil {
		return nil, fmt.Errorf("source cannot be nil when building server filters")
	}

	filter, err := NewServerFilter(source.IncludedServers, source.ExcludedServers)
	if err != nil {
		return nil, fmt.Errorf("invalid includedServers/excludedServers configuration for source %s: %w", source.ID, err)
	}

	return filter, nil
}
