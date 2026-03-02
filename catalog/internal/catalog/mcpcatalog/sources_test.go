package mcpcatalog

import (
	"reflect"
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	"github.com/kubeflow/model-registry/internal/apiutils"
)

func TestMCPSourceCollection_MergeOverride(t *testing.T) {
	tests := []struct {
		name          string
		originOrder   []string
		mergeSequence []struct {
			origin  string
			sources map[string]basecatalog.MCPSource
		}
		expectedSources map[string]basecatalog.MCPSource
	}{
		{
			name:        "later origin overrides earlier origin",
			originOrder: []string{"community.yaml", "org.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]basecatalog.MCPSource
			}{
				{
					origin: "community.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "GitHub MCP Server",
							Type:    "sse",
							Enabled: apiutils.Of(true),
							Labels:  []string{"community"},
							Properties: map[string]any{
								"url": "https://community.example.com/github",
							},
						},
					},
				},
				{
					origin: "org.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "GitHub MCP Server Enterprise",
							Type:    "sse",
							Enabled: apiutils.Of(true),
							Labels:  []string{"enterprise"},
							Properties: map[string]any{
								"url": "https://org.example.com/github",
							},
						},
					},
				},
			},
			expectedSources: map[string]basecatalog.MCPSource{
				"mcp_github": {
					ID:      "mcp_github",
					Name:    "GitHub MCP Server Enterprise",
					Type:    "sse",
					Enabled: apiutils.Of(true),
					Labels:  []string{"enterprise"},
					Properties: map[string]any{
						"url": "https://org.example.com/github",
					},
				},
			},
		},
		{
			name:        "single origin no merge needed",
			originOrder: []string{"community.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]basecatalog.MCPSource
			}{
				{
					origin: "community.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_slack": {
							ID:      "mcp_slack",
							Name:    "Slack MCP Server",
							Type:    "stdio",
							Enabled: apiutils.Of(true),
							Labels:  []string{"community"},
							Properties: map[string]any{
								"command": "slack-mcp",
							},
						},
					},
				},
			},
			expectedSources: map[string]basecatalog.MCPSource{
				"mcp_slack": {
					ID:      "mcp_slack",
					Name:    "Slack MCP Server",
					Type:    "stdio",
					Enabled: apiutils.Of(true),
					Labels:  []string{"community"},
					Properties: map[string]any{
						"command": "slack-mcp",
					},
				},
			},
		},
		{
			name:        "multiple sources only one overridden",
			originOrder: []string{"community.yaml", "org.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]basecatalog.MCPSource
			}{
				{
					origin: "community.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "GitHub MCP Server",
							Type:    "sse",
							Enabled: apiutils.Of(true),
							Labels:  []string{"community"},
						},
						"mcp_slack": {
							ID:      "mcp_slack",
							Name:    "Slack MCP Server",
							Type:    "stdio",
							Enabled: apiutils.Of(true),
							Labels:  []string{"community"},
						},
					},
				},
				{
					origin: "org.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "GitHub MCP Server Enterprise",
							Type:    "sse",
							Enabled: apiutils.Of(true),
							Labels:  []string{"enterprise"},
						},
					},
				},
			},
			expectedSources: map[string]basecatalog.MCPSource{
				"mcp_github": {
					ID:      "mcp_github",
					Name:    "GitHub MCP Server Enterprise",
					Type:    "sse",
					Enabled: apiutils.Of(true),
					Labels:  []string{"enterprise"},
				},
				"mcp_slack": {
					ID:      "mcp_slack",
					Name:    "Slack MCP Server",
					Type:    "stdio",
					Enabled: apiutils.Of(true),
					Labels:  []string{"community"},
				},
			},
		},
		{
			name:        "three-origin cascading override community to org to team",
			originOrder: []string{"community.yaml", "org.yaml", "team.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]basecatalog.MCPSource
			}{
				{
					origin: "community.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "Community GitHub",
							Type:    "sse",
							Enabled: apiutils.Of(true),
							Labels:  []string{"community"},
						},
					},
				},
				{
					origin: "org.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "Org GitHub",
							Type:    "sse",
							Enabled: apiutils.Of(true),
							Labels:  []string{"org"},
						},
					},
				},
				{
					origin: "team.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "Team GitHub",
							Type:    "sse",
							Enabled: apiutils.Of(true),
							Labels:  []string{"team"},
						},
					},
				},
			},
			expectedSources: map[string]basecatalog.MCPSource{
				"mcp_github": {
					ID:      "mcp_github",
					Name:    "Team GitHub",
					Type:    "sse",
					Enabled: apiutils.Of(true),
					Labels:  []string{"team"},
				},
			},
		},
		{
			name:        "override disables an enabled source",
			originOrder: []string{"community.yaml", "org.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]basecatalog.MCPSource
			}{
				{
					origin: "community.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "GitHub MCP Server",
							Type:    "sse",
							Enabled: apiutils.Of(true),
							Labels:  []string{},
						},
					},
				},
				{
					origin: "org.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "GitHub MCP Server",
							Type:    "sse",
							Enabled: apiutils.Of(false), // Disable it
							Labels:  []string{},
						},
					},
				},
			},
			expectedSources: map[string]basecatalog.MCPSource{
				"mcp_github": {
					ID:      "mcp_github",
					Name:    "GitHub MCP Server",
					Type:    "sse",
					Enabled: apiutils.Of(false),
					Labels:  []string{},
				},
			},
		},
		{
			name:        "sparse override: only id and enabled set, all other fields inherited",
			originOrder: []string{"community.yaml", "org.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]basecatalog.MCPSource
			}{
				{
					origin: "community.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "GitHub MCP Server",
							Type:    "sse",
							Enabled: apiutils.Of(false), // Disabled in base
							Labels:  []string{"community", "certified"},
							Properties: map[string]any{
								"url": "https://community.example.com/github",
							},
						},
					},
				},
				{
					origin: "org.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Enabled: apiutils.Of(true), // Only enable it
							// Name, Type, Labels, Properties all empty/nil -> inherited
						},
					},
				},
			},
			expectedSources: map[string]basecatalog.MCPSource{
				"mcp_github": {
					ID:      "mcp_github",
					Name:    "GitHub MCP Server",                // Inherited
					Type:    "sse",                              // Inherited
					Enabled: apiutils.Of(true),                  // Overridden
					Labels:  []string{"community", "certified"}, // Inherited
					Properties: map[string]any{ // Inherited
						"url": "https://community.example.com/github",
					},
				},
			},
		},
		{
			name:        "sparse override: only labels changed",
			originOrder: []string{"community.yaml", "org.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]basecatalog.MCPSource
			}{
				{
					origin: "community.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "GitHub MCP Server",
							Type:    "sse",
							Enabled: apiutils.Of(true),
							Labels:  []string{"community"},
						},
					},
				},
				{
					origin: "org.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:     "mcp_github",
							Labels: []string{"enterprise", "validated"}, // Override only labels
						},
					},
				},
			},
			expectedSources: map[string]basecatalog.MCPSource{
				"mcp_github": {
					ID:      "mcp_github",
					Name:    "GitHub MCP Server",                 // Inherited
					Type:    "sse",                               // Inherited
					Enabled: apiutils.Of(true),                   // Inherited
					Labels:  []string{"enterprise", "validated"}, // Overridden
				},
			},
		},
		{
			name:        "sparse override: empty slice clears labels",
			originOrder: []string{"community.yaml", "org.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]basecatalog.MCPSource
			}{
				{
					origin: "community.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "GitHub MCP Server",
							Type:    "sse",
							Enabled: apiutils.Of(true),
							Labels:  []string{"community", "public"},
						},
					},
				},
				{
					origin: "org.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:     "mcp_github",
							Labels: []string{}, // Explicitly clear labels
						},
					},
				},
			},
			expectedSources: map[string]basecatalog.MCPSource{
				"mcp_github": {
					ID:      "mcp_github",
					Name:    "GitHub MCP Server", // Inherited
					Type:    "sse",               // Inherited
					Enabled: apiutils.Of(true),   // Inherited
					Labels:  []string{},          // Overridden to empty
				},
			},
		},
		{
			name:        "defaults applied: nil Enabled defaults to true, nil Labels defaults to empty slice",
			originOrder: []string{"community.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]basecatalog.MCPSource
			}{
				{
					origin: "community.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:   "mcp_github",
							Name: "GitHub MCP Server",
							Type: "sse",
							// Enabled and Labels are nil - defaults should be applied
						},
					},
				},
			},
			expectedSources: map[string]basecatalog.MCPSource{
				"mcp_github": {
					ID:      "mcp_github",
					Name:    "GitHub MCP Server",
					Type:    "sse",
					Enabled: apiutils.Of(true), // Default applied
					Labels:  []string{},        // Default applied
				},
			},
		},
		{
			name:        "type and properties inherited when not overridden",
			originOrder: []string{"community.yaml", "org.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]basecatalog.MCPSource
			}{
				{
					origin: "community.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "GitHub MCP Server",
							Type:    "sse",
							Enabled: apiutils.Of(false),
							Labels:  []string{"community"},
							Properties: map[string]any{
								"url":     "https://example.com/mcp",
								"timeout": 30,
							},
						},
					},
				},
				{
					origin: "org.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Enabled: apiutils.Of(true), // Only enable it; Type and Properties not set
						},
					},
				},
			},
			expectedSources: map[string]basecatalog.MCPSource{
				"mcp_github": {
					ID:      "mcp_github",
					Name:    "GitHub MCP Server",   // Inherited
					Type:    "sse",                 // Inherited
					Enabled: apiutils.Of(true),     // Overridden
					Labels:  []string{"community"}, // Inherited
					Properties: map[string]any{ // Inherited
						"url":     "https://example.com/mcp",
						"timeout": 30,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msc := NewMCPSourceCollection(tt.originOrder...)

			for _, merge := range tt.mergeSequence {
				err := msc.Merge(merge.origin, merge.sources)
				if err != nil {
					t.Fatalf("Merge(%s) failed: %v", merge.origin, err)
				}
			}

			result := msc.AllSources()

			if len(result) != len(tt.expectedSources) {
				t.Errorf("AllSources() returned %d sources, want %d", len(result), len(tt.expectedSources))
			}

			for id, expected := range tt.expectedSources {
				got, ok := result[id]
				if !ok {
					t.Errorf("source %s not found in result", id)
					continue
				}
				if got.ID != expected.ID {
					t.Errorf("source %s: ID = %s, want %s", id, got.ID, expected.ID)
				}
				if got.Name != expected.Name {
					t.Errorf("source %s: Name = %s, want %s", id, got.Name, expected.Name)
				}
				if got.Type != expected.Type {
					t.Errorf("source %s: Type = %s, want %s", id, got.Type, expected.Type)
				}
				if got.Enabled == nil || expected.Enabled == nil {
					if got.Enabled != expected.Enabled {
						t.Errorf("source %s: Enabled = %v, want %v", id, got.Enabled, expected.Enabled)
					}
				} else if *got.Enabled != *expected.Enabled {
					t.Errorf("source %s: Enabled = %v, want %v", id, *got.Enabled, *expected.Enabled)
				}
				if !reflect.DeepEqual(got.Labels, expected.Labels) {
					t.Errorf("source %s: Labels = %v, want %v", id, got.Labels, expected.Labels)
				}
				if !reflect.DeepEqual(got.Properties, expected.Properties) {
					t.Errorf("source %s: Properties = %v, want %v", id, got.Properties, expected.Properties)
				}
			}
		})
	}
}

func TestMCPSourceCollection_MergeOverride_Origin(t *testing.T) {
	tests := []struct {
		name          string
		originOrder   []string
		mergeSequence []struct {
			origin  string
			sources map[string]basecatalog.MCPSource
		}
		expectedOrigins map[string]string // sourceId -> expected origin
	}{
		{
			name:        "origin preserved from base when Properties NOT overridden",
			originOrder: []string{"/config/community.yaml", "/config/org.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]basecatalog.MCPSource
			}{
				{
					origin: "/config/community.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "GitHub MCP Server",
							Type:    "sse",
							Enabled: apiutils.Of(true),
							Properties: map[string]any{
								"url": "https://community.example.com/github",
							},
							Origin: "/config/community.yaml",
						},
					},
				},
				{
					origin: "/config/org.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_github": {
							ID:      "mcp_github",
							Name:    "GitHub MCP Server Custom",
							Enabled: apiutils.Of(true),
							// Properties not set -> Origin should stay with base
							Origin: "/config/org.yaml",
						},
					},
				},
			},
			expectedOrigins: map[string]string{
				"mcp_github": "/config/community.yaml", // Base origin preserved
			},
		},
		{
			name:        "origin changes when Properties ARE overridden",
			originOrder: []string{"/config/community.yaml", "/config/org.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]basecatalog.MCPSource
			}{
				{
					origin: "/config/community.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_local": {
							ID:      "mcp_local",
							Name:    "Local MCP Server",
							Type:    "stdio",
							Enabled: apiutils.Of(true),
							Properties: map[string]any{
								"command": "community-mcp",
							},
							Origin: "/config/community.yaml",
						},
					},
				},
				{
					origin: "/config/org.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_local": {
							ID:      "mcp_local",
							Enabled: apiutils.Of(true),
							// Override Properties -> Origin should change
							Properties: map[string]any{
								"command": "org-mcp",
							},
							Origin: "/config/org.yaml",
						},
					},
				},
			},
			expectedOrigins: map[string]string{
				"mcp_local": "/config/org.yaml", // Origin changed because Properties overridden
			},
		},
		{
			name:        "multiple sources from different origins keep their own Origins",
			originOrder: []string{"/admin/sources.yaml", "/user/sources.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]basecatalog.MCPSource
			}{
				{
					origin: "/admin/sources.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_admin": {
							ID:      "mcp_admin",
							Name:    "Admin MCP Server",
							Type:    "sse",
							Enabled: apiutils.Of(true),
							Properties: map[string]any{
								"url": "https://admin.example.com/mcp",
							},
							Origin: "/admin/sources.yaml",
						},
					},
				},
				{
					origin: "/user/sources.yaml",
					sources: map[string]basecatalog.MCPSource{
						"mcp_user": {
							ID:      "mcp_user",
							Name:    "User MCP Server",
							Type:    "stdio",
							Enabled: apiutils.Of(true),
							Properties: map[string]any{
								"command": "user-mcp",
							},
							Origin: "/user/sources.yaml",
						},
					},
				},
			},
			expectedOrigins: map[string]string{
				"mcp_admin": "/admin/sources.yaml",
				"mcp_user":  "/user/sources.yaml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msc := NewMCPSourceCollection(tt.originOrder...)

			for _, merge := range tt.mergeSequence {
				err := msc.Merge(merge.origin, merge.sources)
				if err != nil {
					t.Fatalf("Merge(%s) failed: %v", merge.origin, err)
				}
			}

			sources := msc.AllSources()

			for id, expectedOrigin := range tt.expectedOrigins {
				source, ok := sources[id]
				if !ok {
					t.Errorf("source %s not found in AllSources()", id)
					continue
				}
				if source.Origin != expectedOrigin {
					t.Errorf("source %s: Origin = %s, want %s", id, source.Origin, expectedOrigin)
				}
			}
		})
	}
}

func TestMCPSourceCollection_MergeOverride_DynamicOrigin(t *testing.T) {
	// Test that origins not in the initial originOrder are appended dynamically
	// when Merge() is called with an unknown origin.
	msc := NewMCPSourceCollection("community.yaml")

	err := msc.Merge("community.yaml", map[string]basecatalog.MCPSource{
		"mcp_github": {ID: "mcp_github", Name: "Community", Enabled: apiutils.Of(true), Labels: []string{}},
	})
	if err != nil {
		t.Fatalf("Merge(community.yaml) failed: %v", err)
	}

	// Dynamic origin not in initial order - should be appended and take precedence
	err = msc.Merge("extra.yaml", map[string]basecatalog.MCPSource{
		"mcp_github": {ID: "mcp_github", Name: "Extra Override", Enabled: apiutils.Of(true), Labels: []string{}},
	})
	if err != nil {
		t.Fatalf("Merge(extra.yaml) failed: %v", err)
	}

	result := msc.AllSources()
	source, ok := result["mcp_github"]
	if !ok {
		t.Fatal("AllSources() should return mcp_github")
	}

	if source.Name != "Extra Override" {
		t.Errorf("dynamically added origin should override earlier origins, got Name = %s, want 'Extra Override'", source.Name)
	}
}
