package catalog

import (
	"fmt"
	"regexp"
	"strings"
)

// ItemFilter provides include/exclude pattern matching for item names.
// It uses glob-style patterns with '*' as the only supported wildcard.
type ItemFilter struct {
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
	b.WriteString("(?i)^") // case insensitive
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

// ValidatePatterns validates that the included and excluded patterns
// are valid (non-empty, compilable). This is useful for early validation
// at configuration load time without constructing the full ItemFilter.
func ValidatePatterns(included, excluded []string) error {
	if _, err := compilePatterns("included", included); err != nil {
		return err
	}

	if _, err := compilePatterns("excluded", excluded); err != nil {
		return err
	}

	return nil
}

// NewItemFilter builds an ItemFilter from the provided include/exclude pattern lists.
// Patterns support glob-style wildcards ('*').
//
// Include logic:
//   - If included is non-empty, items must match at least one pattern to be allowed.
//   - If included is empty, all items are allowed (subject to exclusions).
//
// Exclude logic:
//   - Items matching any excluded pattern are rejected, even if they match an include.
//
// Returns nil if both lists are empty (no filtering needed).
func NewItemFilter(included, excluded []string) (*ItemFilter, error) {
	if err := ValidatePatterns(included, excluded); err != nil {
		return nil, err
	}

	inc, err := compilePatterns("included", included)
	if err != nil {
		return nil, err
	}

	exc, err := compilePatterns("excluded", excluded)
	if err != nil {
		return nil, err
	}

	if len(inc) == 0 && len(exc) == 0 {
		return nil, nil
	}

	return &ItemFilter{
		included: inc,
		excluded: exc,
	}, nil
}

// NewItemFilterFromSource creates an ItemFilter from a Source's configuration.
// Additional patterns can be appended via extraIncluded and extraExcluded.
func NewItemFilterFromSource(source *Source, extraIncluded, extraExcluded []string) (*ItemFilter, error) {
	if source == nil {
		return nil, fmt.Errorf("source cannot be nil when building filters")
	}

	included := append([]string{}, source.IncludedItems...)
	if len(extraIncluded) > 0 {
		included = append(included, extraIncluded...)
	}

	excluded := append([]string{}, source.ExcludedItems...)
	if len(extraExcluded) > 0 {
		excluded = append(excluded, extraExcluded...)
	}

	filter, err := NewItemFilter(included, excluded)
	if err != nil {
		return nil, fmt.Errorf("invalid include/exclude configuration for source %s: %w", source.ID, err)
	}

	return filter, nil
}

// Allows returns true if the provided item name passes the include/exclude rules.
// A nil filter allows everything.
func (f *ItemFilter) Allows(name string) bool {
	if f == nil {
		return true
	}

	// Check include patterns first
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

	// Check exclude patterns
	for _, pattern := range f.excluded {
		if pattern.re.MatchString(name) {
			return false
		}
	}

	return true
}

// HasPatterns returns true if the filter has any include or exclude patterns.
func (f *ItemFilter) HasPatterns() bool {
	if f == nil {
		return false
	}
	return len(f.included) > 0 || len(f.excluded) > 0
}
