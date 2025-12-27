package catalog

import (
	"fmt"
	"strings"

	"github.com/kubeflow/model-registry/catalog/internal/mcp"
)

// McpServerFilter encapsulates include/exclude pattern matching for MCP server names.
// It reuses the same compiledPattern type from model_filter.go for consistency.
type McpServerFilter struct {
	included []*compiledPattern
	excluded []*compiledPattern
}

// ValidateMcpServerSourceFilters validates that the includedServers and excludedServers patterns
// are valid (non-empty, compilable, non-conflicting). This is useful for early validation
// at configuration load time without constructing the full McpServerFilter.
func ValidateMcpServerSourceFilters(included, excluded []string) error {
	if err := detectConflictingMcpServerPatterns(included, excluded); err != nil {
		return err
	}

	if _, err := compilePatterns("includedServers", included); err != nil {
		return err
	}

	if _, err := compilePatterns("excludedServers", excluded); err != nil {
		return err
	}

	return nil
}

// NewMcpServerFilter builds a McpServerFilter from the provided include/exclude pattern lists.
func NewMcpServerFilter(included, excluded []string) (*McpServerFilter, error) {
	if err := ValidateMcpServerSourceFilters(included, excluded); err != nil {
		return nil, err
	}

	inc, err := compilePatterns("includedServers", included)
	if err != nil {
		return nil, err
	}

	exc, err := compilePatterns("excludedServers", excluded)
	if err != nil {
		return nil, err
	}

	if len(inc) == 0 && len(exc) == 0 {
		return nil, nil
	}

	return &McpServerFilter{
		included: inc,
		excluded: exc,
	}, nil
}

// detectConflictingMcpServerPatterns checks if any pattern appears in both included and excluded lists.
func detectConflictingMcpServerPatterns(included, excluded []string) error {
	if len(included) == 0 || len(excluded) == 0 {
		return nil
	}

	includedIdx := make(map[string]int, len(included))
	for i, pattern := range included {
		value := strings.TrimSpace(pattern)
		includedIdx[value] = i
	}

	for j, pattern := range excluded {
		value := strings.TrimSpace(pattern)
		if i, exists := includedIdx[value]; exists {
			return fmt.Errorf("pattern %q is defined in both includedServers[%d] and excludedServers[%d]", value, i, j)
		}
	}
	return nil
}

// Allows returns true if the provided MCP server name passes the include/exclude rules.
func (f *McpServerFilter) Allows(name string) bool {
	if f == nil {
		return true
	}

	if len(f.included) > 0 {
		matched := false
		for _, pattern := range f.included {
			if pattern.re.MatchString(name) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	for _, pattern := range f.excluded {
		if pattern.re.MatchString(name) {
			return false
		}
	}

	return true
}

// NewMcpServerFilterFromSource composes a McpServerFilter using the source-level configuration.
func NewMcpServerFilterFromSource(source *mcp.McpSource) (*McpServerFilter, error) {
	if source == nil {
		return nil, fmt.Errorf("source cannot be nil when building filters")
	}

	filter, err := NewMcpServerFilter(source.IncludedServers, source.ExcludedServers)
	if err != nil {
		return nil, fmt.Errorf("invalid include/exclude configuration for MCP source %s: %w", source.Id, err)
	}

	return filter, nil
}

