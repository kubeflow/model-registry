package catalog

import (
	"testing"
)

func TestItemFilter(t *testing.T) {
	tests := []struct {
		name     string
		included []string
		excluded []string
		items    map[string]bool // item -> expected result
	}{
		{
			name:     "no patterns allows everything",
			included: nil,
			excluded: nil,
			items: map[string]bool{
				"anything": true,
				"foo":      true,
			},
		},
		{
			name:     "include pattern only",
			included: []string{"foo*"},
			excluded: nil,
			items: map[string]bool{
				"foo":     true,
				"foobar":  true,
				"bar":     false,
				"barfoo":  false,
				"FOO":     true, // case insensitive
				"FOOBAR":  true,
			},
		},
		{
			name:     "exclude pattern only",
			included: nil,
			excluded: []string{"*test*"},
			items: map[string]bool{
				"foo":        true,
				"test":       false,
				"testing":    false,
				"mytest":     false,
				"mytesting":  false,
				"production": true,
			},
		},
		{
			name:     "include and exclude combined",
			included: []string{"model-*"},
			excluded: []string{"*-test"},
			items: map[string]bool{
				"model-a":      true,
				"model-b":      true,
				"model-test":   false, // excluded
				"other-model":  false, // not included
				"model-a-test": false, // excluded
			},
		},
		{
			name:     "exact match",
			included: []string{"mymodel"},
			excluded: nil,
			items: map[string]bool{
				"mymodel":   true,
				"mymodel2":  false,
				"themymodel": false,
				"MYMODEL":   true, // case insensitive
			},
		},
		{
			name:     "wildcard at start",
			included: []string{"*-v1"},
			excluded: nil,
			items: map[string]bool{
				"model-v1":  true,
				"other-v1":  true,
				"model-v2":  false,
				"v1":        false, // doesn't end with -v1
				"-v1":       true,
			},
		},
		{
			name:     "multiple includes",
			included: []string{"foo*", "bar*"},
			excluded: nil,
			items: map[string]bool{
				"foo":   true,
				"fooX":  true,
				"bar":   true,
				"barY":  true,
				"baz":   false,
				"other": false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := NewItemFilter(tt.included, tt.excluded)
			if err != nil {
				t.Fatalf("NewItemFilter failed: %v", err)
			}

			for item, expected := range tt.items {
				result := filter.Allows(item)
				if result != expected {
					t.Errorf("Allows(%q) = %v, want %v", item, result, expected)
				}
			}
		})
	}
}

func TestItemFilterNil(t *testing.T) {
	// nil filter allows everything
	var filter *ItemFilter
	if !filter.Allows("anything") {
		t.Error("nil filter should allow everything")
	}
}

func TestItemFilterEmptyReturnNil(t *testing.T) {
	filter, err := NewItemFilter(nil, nil)
	if err != nil {
		t.Fatalf("NewItemFilter failed: %v", err)
	}
	if filter != nil {
		t.Error("Expected nil filter when no patterns provided")
	}
}

func TestValidatePatterns(t *testing.T) {
	// Valid patterns
	if err := ValidatePatterns([]string{"foo*"}, []string{"bar*"}); err != nil {
		t.Errorf("ValidatePatterns failed for valid patterns: %v", err)
	}

	// Empty pattern should fail
	if err := ValidatePatterns([]string{""}, nil); err == nil {
		t.Error("Expected error for empty pattern")
	}

	// Whitespace-only pattern should fail
	if err := ValidatePatterns([]string{"   "}, nil); err == nil {
		t.Error("Expected error for whitespace-only pattern")
	}
}

func TestNewItemFilterFromSource(t *testing.T) {
	source := &Source{
		ID:            "test",
		IncludedItems: []string{"model-*"},
		ExcludedItems: []string{"*-test"},
	}

	filter, err := NewItemFilterFromSource(source, nil, nil)
	if err != nil {
		t.Fatalf("NewItemFilterFromSource failed: %v", err)
	}

	if !filter.Allows("model-a") {
		t.Error("Expected model-a to be allowed")
	}
	if filter.Allows("model-test") {
		t.Error("Expected model-test to be excluded")
	}
}

func TestNewItemFilterFromSourceWithExtras(t *testing.T) {
	source := &Source{
		ID:            "test",
		IncludedItems: []string{"model-*"},
	}

	// Add extra excluded patterns
	filter, err := NewItemFilterFromSource(source, nil, []string{"*-deprecated"})
	if err != nil {
		t.Fatalf("NewItemFilterFromSource failed: %v", err)
	}

	if !filter.Allows("model-a") {
		t.Error("Expected model-a to be allowed")
	}
	if filter.Allows("model-deprecated") {
		t.Error("Expected model-deprecated to be excluded by extra pattern")
	}
}

func TestItemFilterHasPatterns(t *testing.T) {
	var nilFilter *ItemFilter
	if nilFilter.HasPatterns() {
		t.Error("nil filter should not have patterns")
	}

	filter, _ := NewItemFilter([]string{"foo*"}, nil)
	if !filter.HasPatterns() {
		t.Error("filter with include patterns should have patterns")
	}

	filter, _ = NewItemFilter(nil, []string{"bar*"})
	if !filter.HasPatterns() {
		t.Error("filter with exclude patterns should have patterns")
	}
}
