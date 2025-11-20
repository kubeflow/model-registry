package catalog

import (
	"fmt"
	"regexp"
	"strings"
)

// ModelFilter encapsulates include/exclude pattern matching for model names.
type ModelFilter struct {
	included []*compiledPattern
	excluded []*compiledPattern
}

type compiledPattern struct {
	raw string
	re  *regexp.Regexp
}

func newCompiledPattern(field string, idx int, raw string) (*compiledPattern, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil, fmt.Errorf("%s[%d]: pattern cannot be empty", field, idx)
	}

	// Convert a simple glob (only supporting '*') into a regexp.
	var b strings.Builder
	b.WriteString("(?i)^")
	for _, r := range value {
		if r == '*' {
			b.WriteString(".*")
			continue
		}
		b.WriteString(regexp.QuoteMeta(string(r)))
	}
	b.WriteString("$")

	re, err := regexp.Compile(b.String())
	if err != nil {
		return nil, fmt.Errorf("%s[%d]: invalid pattern %q: %w", field, idx, value, err)
	}

	return &compiledPattern{
		raw: value,
		re:  re,
	}, nil
}

func compilePatterns(field string, patterns []string) ([]*compiledPattern, error) {
	if len(patterns) == 0 {
		return nil, nil
	}

	compiled := make([]*compiledPattern, 0, len(patterns))
	for i, pattern := range patterns {
		cp, err := newCompiledPattern(field, i, pattern)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, cp)
	}
	return compiled, nil
}

// ValidateSourceFilters validates that the includedModels and excludedModels patterns
// are valid (non-empty, compilable, non-conflicting). This is useful for early validation
// at configuration load time without constructing the full ModelFilter.
func ValidateSourceFilters(included, excluded []string) error {
	if err := detectConflictingPatterns(included, excluded); err != nil {
		return err
	}

	if _, err := compilePatterns("includedModels", included); err != nil {
		return err
	}

	if _, err := compilePatterns("excludedModels", excluded); err != nil {
		return err
	}

	return nil
}

// NewModelFilter builds a ModelFilter from the provided include/exclude pattern lists.
func NewModelFilter(included, excluded []string) (*ModelFilter, error) {
	if err := ValidateSourceFilters(included, excluded); err != nil {
		return nil, err
	}

	inc, err := compilePatterns("includedModels", included)
	if err != nil {
		return nil, err
	}

	exc, err := compilePatterns("excludedModels", excluded)
	if err != nil {
		return nil, err
	}

	if len(inc) == 0 && len(exc) == 0 {
		return nil, nil
	}

	return &ModelFilter{
		included: inc,
		excluded: exc,
	}, nil
}

func detectConflictingPatterns(included, excluded []string) error {
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
			return fmt.Errorf("pattern %q is defined in both includedModels[%d] and excludedModels[%d]", value, i, j)
		}
	}
	return nil
}

// Allows returns true if the provided model name passes the include/exclude rules.
func (f *ModelFilter) Allows(name string) bool {
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

// NewModelFilterFromSource composes a ModelFilter using the source-level configuration and any legacy additions.
func NewModelFilterFromSource(source *Source, extraIncluded, extraExcluded []string) (*ModelFilter, error) {
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
