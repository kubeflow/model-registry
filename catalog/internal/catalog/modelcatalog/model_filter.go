package modelcatalog

import (
	"fmt"

	"github.com/kubeflow/hub/catalog/internal/catalog/basecatalog"
)

// ModelFilter encapsulates include/exclude pattern matching for model names.
type ModelFilter = basecatalog.NameFilter

// ValidateSourceFilters validates that the includedModels and excludedModels patterns
// are valid (non-empty, compilable, non-conflicting). This is useful for early validation
// at configuration load time without constructing the full ModelFilter.
func ValidateSourceFilters(included, excluded []string) error {
	return basecatalog.ValidatePatterns("includedModels", included, "excludedModels", excluded)
}

// NewModelFilter builds a ModelFilter from the provided include/exclude pattern lists.
func NewModelFilter(included, excluded []string) (*ModelFilter, error) {
	return basecatalog.NewNameFilter("includedModels", included, "excludedModels", excluded)
}

// NewModelFilterFromSource composes a ModelFilter using the source-level configuration and any legacy additions.
func NewModelFilterFromSource(source *basecatalog.ModelSource, extraIncluded, extraExcluded []string) (*ModelFilter, error) {
	if source == nil {
		return nil, fmt.Errorf("source cannot be nil when building filters")
	}

	included := append([]string{}, source.IncludedModels...)
	if len(extraIncluded) > 0 {
		included = append(included, extraIncluded...)
	}

	excluded := append([]string{}, source.ExcludedModels...)
	if len(extraExcluded) > 0 {
		excluded = append(excluded, extraExcluded...)
	}

	filter, err := NewModelFilter(included, excluded)
	if err != nil {
		return nil, fmt.Errorf("invalid include/exclude configuration for source %s: %w", source.Id, err)
	}

	return filter, nil
}
