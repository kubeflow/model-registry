package catalog

import (
	"testing"
)

func TestSourceIsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		enabled  *bool
		expected bool
	}{
		{
			name:     "nil enabled defaults to true",
			enabled:  nil,
			expected: true,
		},
		{
			name:     "explicitly enabled",
			enabled:  boolPtr(true),
			expected: true,
		},
		{
			name:     "explicitly disabled",
			enabled:  boolPtr(false),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Source{Enabled: tt.enabled}
			if got := s.IsEnabled(); got != tt.expected {
				t.Errorf("IsEnabled() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSourceCollection(t *testing.T) {
	sc := NewSourceCollection("default.yaml", "user.yaml")

	// Add default sources
	defaultSources := map[string]Source{
		"source1": {
			ID:      "source1",
			Name:    "Default Source 1",
			Type:    "yaml",
			Enabled: boolPtr(false),
		},
		"source2": {
			ID:   "source2",
			Name: "Default Source 2",
			Type: "http",
		},
	}
	if err := sc.Merge("default.yaml", defaultSources); err != nil {
		t.Fatalf("Merge default failed: %v", err)
	}

	// Add user sources (should override)
	userSources := map[string]Source{
		"source1": {
			ID:      "source1",
			Enabled: boolPtr(true), // Enable the disabled source
		},
	}
	if err := sc.Merge("user.yaml", userSources); err != nil {
		t.Fatalf("Merge user failed: %v", err)
	}

	// Test AllSources
	all := sc.AllSources()
	if len(all) != 2 {
		t.Errorf("Expected 2 sources, got %d", len(all))
	}

	// Test field-level merging: source1 should have merged fields
	s1, ok := all["source1"]
	if !ok {
		t.Fatal("Expected source1 in AllSources")
	}
	if s1.Name != "Default Source 1" {
		t.Errorf("Expected name 'Default Source 1', got '%s'", s1.Name)
	}
	if s1.Type != "yaml" {
		t.Errorf("Expected type 'yaml', got '%s'", s1.Type)
	}
	if !s1.IsEnabled() {
		t.Error("Expected source1 to be enabled after merge")
	}

	// Test Get (only returns enabled sources)
	_, ok = sc.Get("source1")
	if !ok {
		t.Error("Expected to find enabled source1")
	}

	// Test EnabledSources
	enabled := sc.EnabledSources()
	if len(enabled) != 2 {
		t.Errorf("Expected 2 enabled sources, got %d", len(enabled))
	}

	// Test IDs
	ids := sc.IDs()
	if len(ids) != 2 {
		t.Errorf("Expected 2 IDs, got %d", len(ids))
	}
}

func TestSourceCollectionByLabel(t *testing.T) {
	sc := NewSourceCollection()

	sources := map[string]Source{
		"source1": {
			ID:     "source1",
			Labels: []string{"prod", "ml"},
		},
		"source2": {
			ID:     "source2",
			Labels: []string{"dev"},
		},
		"source3": {
			ID:     "source3",
			Labels: nil, // No labels
		},
	}
	sc.Merge("test.yaml", sources)

	// Test finding by label
	result := sc.ByLabel([]string{"prod"})
	if len(result) != 1 {
		t.Errorf("Expected 1 source with 'prod' label, got %d", len(result))
	}

	// Test case insensitive
	result = sc.ByLabel([]string{"PROD"})
	if len(result) != 1 {
		t.Errorf("Expected 1 source with 'PROD' label (case insensitive), got %d", len(result))
	}

	// Test "null" for sources without labels
	result = sc.ByLabel([]string{"null"})
	if len(result) != 1 {
		t.Errorf("Expected 1 source without labels, got %d", len(result))
	}

	// Test multiple labels (OR logic)
	result = sc.ByLabel([]string{"prod", "dev"})
	if len(result) != 2 {
		t.Errorf("Expected 2 sources with 'prod' or 'dev' labels, got %d", len(result))
	}
}

func TestSourceCollectionDynamicOrigin(t *testing.T) {
	// Test that origins can be added dynamically
	sc := NewSourceCollection()

	sources := map[string]Source{
		"source1": {ID: "source1", Type: "yaml"},
	}
	if err := sc.Merge("dynamic.yaml", sources); err != nil {
		t.Fatalf("Dynamic merge failed: %v", err)
	}

	all := sc.AllSources()
	if len(all) != 1 {
		t.Errorf("Expected 1 source after dynamic add, got %d", len(all))
	}
}

func boolPtr(b bool) *bool {
	return &b
}
