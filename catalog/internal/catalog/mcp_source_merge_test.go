package catalog

import (
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/mcp"
	"github.com/stretchr/testify/assert"
)

func ptrBool(b bool) *bool {
	return &b
}

func TestMergeMcpSources(t *testing.T) {
	tests := []struct {
		name     string
		base     mcp.McpSource
		override mcp.McpSource
		expected mcp.McpSource
	}{
		{
			name: "override takes precedence for all set fields",
			base: mcp.McpSource{
				Id:              "test-source",
				Name:            "Base Name",
				Type:            "yaml",
				Enabled:         ptrBool(true),
				Labels:          []string{"base-label"},
				IncludedServers: []string{"base-*"},
				ExcludedServers: []string{"*-deprecated"},
				Properties:      map[string]any{"yamlCatalogPath": "base.yaml"},
				Origin:          "/base/path",
			},
			override: mcp.McpSource{
				Id:              "test-source",
				Name:            "Override Name",
				Type:            "remote",
				Enabled:         ptrBool(false),
				Labels:          []string{"override-label"},
				IncludedServers: []string{"override-*"},
				ExcludedServers: []string{"*-alpha"},
				Properties:      map[string]any{"yamlCatalogPath": "override.yaml"},
				Origin:          "/override/path",
			},
			expected: mcp.McpSource{
				Id:              "test-source",
				Name:            "Override Name",
				Type:            "remote",
				Enabled:         ptrBool(false),
				Labels:          []string{"override-label"},
				IncludedServers: []string{"override-*"},
				ExcludedServers: []string{"*-alpha"},
				Properties:      map[string]any{"yamlCatalogPath": "override.yaml"},
				Origin:          "/override/path", // Changed because Properties was overridden
			},
		},
		{
			name: "empty override preserves base values",
			base: mcp.McpSource{
				Id:              "test-source",
				Name:            "Base Name",
				Type:            "yaml",
				Enabled:         ptrBool(true),
				Labels:          []string{"base-label"},
				IncludedServers: []string{"base-*"},
				ExcludedServers: []string{"*-deprecated"},
				Properties:      map[string]any{"yamlCatalogPath": "base.yaml"},
				Origin:          "/base/path",
			},
			override: mcp.McpSource{
				Id: "test-source",
				// All other fields are empty/nil
			},
			expected: mcp.McpSource{
				Id:              "test-source",
				Name:            "Base Name",
				Type:            "yaml",
				Enabled:         ptrBool(true),
				Labels:          []string{"base-label"},
				IncludedServers: []string{"base-*"},
				ExcludedServers: []string{"*-deprecated"},
				Properties:      map[string]any{"yamlCatalogPath": "base.yaml"},
				Origin:          "/base/path",
			},
		},
		{
			name: "empty slice in override clears base slices",
			base: mcp.McpSource{
				Id:              "test-source",
				Labels:          []string{"label1", "label2"},
				IncludedServers: []string{"include-*"},
				ExcludedServers: []string{"exclude-*"},
			},
			override: mcp.McpSource{
				Id:              "test-source",
				Labels:          []string{}, // Explicitly empty
				IncludedServers: []string{}, // Explicitly empty
				ExcludedServers: []string{}, // Explicitly empty
			},
			expected: mcp.McpSource{
				Id:              "test-source",
				Labels:          []string{},
				IncludedServers: []string{},
				ExcludedServers: []string{},
			},
		},
		{
			name: "partial override preserves non-overridden fields",
			base: mcp.McpSource{
				Id:              "test-source",
				Name:            "Base Name",
				Type:            "yaml",
				Enabled:         ptrBool(true),
				Labels:          []string{"community"},
				IncludedServers: nil,
				ExcludedServers: nil,
				Properties:      map[string]any{"yamlCatalogPath": "base.yaml"},
				Origin:          "/base/path",
			},
			override: mcp.McpSource{
				Id:              "test-source",
				Labels:          []string{"enterprise", "validated"}, // Only override labels
				ExcludedServers: []string{"*-alpha"},                 // Add exclusions
			},
			expected: mcp.McpSource{
				Id:              "test-source",
				Name:            "Base Name",
				Type:            "yaml",
				Enabled:         ptrBool(true),
				Labels:          []string{"enterprise", "validated"},
				IncludedServers: nil, // Preserved from base
				ExcludedServers: []string{"*-alpha"},
				Properties:      map[string]any{"yamlCatalogPath": "base.yaml"},
				Origin:          "/base/path",
			},
		},
		{
			name: "origin preserved when properties not overridden",
			base: mcp.McpSource{
				Id:         "test-source",
				Properties: map[string]any{"yamlCatalogPath": "base.yaml"},
				Origin:     "/base/path",
			},
			override: mcp.McpSource{
				Id:     "test-source",
				Labels: []string{"new-label"},
				Origin: "/override/path",
				// Properties is nil
			},
			expected: mcp.McpSource{
				Id:         "test-source",
				Labels:     []string{"new-label"},
				Properties: map[string]any{"yamlCatalogPath": "base.yaml"},
				Origin:     "/base/path", // Preserved because Properties wasn't overridden
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeMcpSources(tt.base, tt.override)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestApplyMcpSourceDefaults(t *testing.T) {
	tests := []struct {
		name     string
		source   mcp.McpSource
		expected mcp.McpSource
	}{
		{
			name: "defaults applied when fields are nil",
			source: mcp.McpSource{
				Id:   "test-source",
				Name: "Test",
				// Enabled and Labels are nil
			},
			expected: mcp.McpSource{
				Id:      "test-source",
				Name:    "Test",
				Enabled: ptrBool(true),
				Labels:  []string{},
			},
		},
		{
			name: "no changes when fields already set",
			source: mcp.McpSource{
				Id:      "test-source",
				Enabled: ptrBool(false),
				Labels:  []string{"custom-label"},
			},
			expected: mcp.McpSource{
				Id:      "test-source",
				Enabled: ptrBool(false),
				Labels:  []string{"custom-label"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyMcpSourceDefaults(tt.source)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMergeMcpSourcesFromPaths(t *testing.T) {
	tests := []struct {
		name     string
		paths    []string
		sources  map[string][]mcp.McpSource // map path -> sources
		expected map[string]mcp.McpSource
	}{
		{
			name:  "single path returns sources with defaults",
			paths: []string{"/path/one.yaml"},
			sources: map[string][]mcp.McpSource{
				"/path/one.yaml": {
					{Id: "source-a", Name: "Source A", Type: "yaml"},
					{Id: "source-b", Name: "Source B", Type: "yaml"},
				},
			},
			expected: map[string]mcp.McpSource{
				"source-a": {Id: "source-a", Name: "Source A", Type: "yaml", Origin: "/path/one.yaml", Enabled: ptrBool(true), Labels: []string{}},
				"source-b": {Id: "source-b", Name: "Source B", Type: "yaml", Origin: "/path/one.yaml", Enabled: ptrBool(true), Labels: []string{}},
			},
		},
		{
			name:  "later paths override earlier paths",
			paths: []string{"/path/base.yaml", "/path/override.yaml"},
			sources: map[string][]mcp.McpSource{
				"/path/base.yaml": {
					{
						Id:     "shared-source",
						Name:   "Community Source",
						Type:   "yaml",
						Labels: []string{"community"},
						Properties: map[string]any{
							"yamlCatalogPath": "community.yaml",
						},
					},
				},
				"/path/override.yaml": {
					{
						Id:              "shared-source",
						Labels:          []string{"enterprise", "validated"},
						ExcludedServers: []string{"*-alpha", "*-beta"},
						// Type and Properties not set - should be inherited
					},
				},
			},
			expected: map[string]mcp.McpSource{
				"shared-source": {
					Id:              "shared-source",
					Name:            "Community Source", // Preserved from base
					Type:            "yaml",             // Preserved from base
					Labels:          []string{"enterprise", "validated"},
					ExcludedServers: []string{"*-alpha", "*-beta"},
					Properties: map[string]any{
						"yamlCatalogPath": "community.yaml",
					},
					Origin:  "/path/base.yaml", // Preserved (Properties not overridden)
					Enabled: ptrBool(true),
				},
			},
		},
		{
			name:  "three-way merge with priority order",
			paths: []string{"/default.yaml", "/team.yaml", "/user.yaml"},
			sources: map[string][]mcp.McpSource{
				"/default.yaml": {
					{
						Id:     "k8s-mcp",
						Name:   "Kubernetes MCP",
						Type:   "yaml",
						Labels: []string{"default"},
						Properties: map[string]any{
							"yamlCatalogPath": "k8s-mcp.yaml",
						},
					},
				},
				"/team.yaml": {
					{
						Id:              "k8s-mcp",
						Labels:          []string{"team", "validated"},
						IncludedServers: []string{"k8s-*"},
					},
				},
				"/user.yaml": {
					{
						Id:              "k8s-mcp",
						ExcludedServers: []string{"k8s-deprecated-*"},
					},
				},
			},
			expected: map[string]mcp.McpSource{
				"k8s-mcp": {
					Id:     "k8s-mcp",
					Name:   "Kubernetes MCP",
					Type:   "yaml",
					Labels: []string{"team", "validated"}, // From /team.yaml
					Properties: map[string]any{
						"yamlCatalogPath": "k8s-mcp.yaml",
					},
					IncludedServers: []string{"k8s-*"},            // From /team.yaml
					ExcludedServers: []string{"k8s-deprecated-*"}, // From /user.yaml
					Origin:          "/default.yaml",              // Properties from default
					Enabled:         ptrBool(true),
				},
			},
		},
		{
			name:  "disable source in override",
			paths: []string{"/base.yaml", "/override.yaml"},
			sources: map[string][]mcp.McpSource{
				"/base.yaml": {
					{Id: "source-a", Name: "Source A", Type: "yaml", Enabled: ptrBool(true)},
				},
				"/override.yaml": {
					{Id: "source-a", Enabled: ptrBool(false)},
				},
			},
			expected: map[string]mcp.McpSource{
				"source-a": {Id: "source-a", Name: "Source A", Type: "yaml", Enabled: ptrBool(false), Labels: []string{}, Origin: "/base.yaml"},
			},
		},
		{
			name:  "multiple independent sources from different paths",
			paths: []string{"/path/one.yaml", "/path/two.yaml"},
			sources: map[string][]mcp.McpSource{
				"/path/one.yaml": {
					{Id: "source-a", Name: "Source A", Type: "yaml"},
				},
				"/path/two.yaml": {
					{Id: "source-b", Name: "Source B", Type: "yaml"},
				},
			},
			expected: map[string]mcp.McpSource{
				"source-a": {Id: "source-a", Name: "Source A", Type: "yaml", Origin: "/path/one.yaml", Enabled: ptrBool(true), Labels: []string{}},
				"source-b": {Id: "source-b", Name: "Source B", Type: "yaml", Origin: "/path/two.yaml", Enabled: ptrBool(true), Labels: []string{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReadFunc := func(path string) ([]mcp.McpSource, error) {
				if sources, ok := tt.sources[path]; ok {
					return sources, nil
				}
				return nil, nil
			}

			result, err := MergeMcpSourcesFromPaths(tt.paths, mockReadFunc)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMergeMcpSourcesFromPathsWithMissingFiles(t *testing.T) {
	paths := []string{"/missing.yaml", "/valid.yaml"}
	sources := map[string][]mcp.McpSource{
		"/valid.yaml": {
			{Id: "source-a", Name: "Source A", Type: "yaml"},
		},
	}

	mockReadFunc := func(path string) ([]mcp.McpSource, error) {
		if s, ok := sources[path]; ok {
			return s, nil
		}
		return nil, nil // Missing file returns nil
	}

	result, err := MergeMcpSourcesFromPaths(paths, mockReadFunc)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "source-a", result["source-a"].Id)
}
