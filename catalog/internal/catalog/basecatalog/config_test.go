package basecatalog

import (
	"testing"

	apimodels "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/stretchr/testify/assert"
)

func TestSourceConfig_GetModelCatalogs(t *testing.T) {
	tests := []struct {
		name     string
		config   *SourceConfig
		expected int
	}{
		{
			name: "only model_catalogs",
			config: &SourceConfig{
				ModelCatalogs: []ModelSource{
					{CatalogSource: apimodels.CatalogSource{Id: "model1", Name: "Model 1"}},
					{CatalogSource: apimodels.CatalogSource{Id: "model2", Name: "Model 2"}},
				},
			},
			expected: 2,
		},
		{
			name: "only deprecated catalogs",
			config: &SourceConfig{
				Catalogs: []ModelSource{
					{CatalogSource: apimodels.CatalogSource{Id: "cat1", Name: "Catalog 1"}},
				},
			},
			expected: 1,
		},
		{
			name: "both fields - model_catalogs takes precedence on conflict",
			config: &SourceConfig{
				ModelCatalogs: []ModelSource{
					{CatalogSource: apimodels.CatalogSource{Id: "shared", Name: "New Name"}},
				},
				Catalogs: []ModelSource{
					{CatalogSource: apimodels.CatalogSource{Id: "shared", Name: "Old Name"}},
					{CatalogSource: apimodels.CatalogSource{Id: "unique", Name: "Unique"}},
				},
			},
			expected: 2, // "shared" from model_catalogs + "unique" from catalogs
		},
		{
			name:     "empty config",
			config:   &SourceConfig{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetModelCatalogs()
			assert.Len(t, result, tt.expected)

			// For the precedence test, verify model_catalogs wins
			if tt.name == "both fields - model_catalogs takes precedence on conflict" {
				var sharedSource *ModelSource
				for _, s := range result {
					if s.Id == "shared" {
						sharedSource = &s
						break
					}
				}
				assert.NotNil(t, sharedSource)
				assert.Equal(t, "New Name", sharedSource.Name)
			}
		})
	}
}

func TestSourceConfig_HasDeprecatedCatalogs(t *testing.T) {
	tests := []struct {
		name     string
		config   *SourceConfig
		expected bool
	}{
		{
			name: "has deprecated catalogs",
			config: &SourceConfig{
				Catalogs: []ModelSource{
					{CatalogSource: apimodels.CatalogSource{Id: "cat1", Name: "Catalog 1"}},
				},
			},
			expected: true,
		},
		{
			name: "no deprecated catalogs",
			config: &SourceConfig{
				ModelCatalogs: []ModelSource{
					{CatalogSource: apimodels.CatalogSource{Id: "model1", Name: "Model 1"}},
				},
			},
			expected: false,
		},
		{
			name:     "empty config",
			config:   &SourceConfig{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.HasDeprecatedCatalogs()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSourceConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    *SourceConfig
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid config with model catalogs",
			config: &SourceConfig{
				ModelCatalogs: []ModelSource{
					{CatalogSource: apimodels.CatalogSource{Id: "model1", Name: "Model 1"}},
				},
			},
			expectErr: false,
		},
		{
			name: "valid config with MCP catalogs",
			config: &SourceConfig{
				MCPCatalogs: []MCPSource{
					{ID: "mcp1", Name: "MCP 1", Type: "yaml"},
				},
			},
			expectErr: false,
		},
		{
			name: "duplicate model catalog IDs",
			config: &SourceConfig{
				ModelCatalogs: []ModelSource{
					{CatalogSource: apimodels.CatalogSource{Id: "dup", Name: "Model 1"}},
					{CatalogSource: apimodels.CatalogSource{Id: "dup", Name: "Model 2"}},
				},
			},
			expectErr: true,
			errMsg:    "duplicate model catalog id: dup",
		},
		{
			name: "duplicate MCP catalog IDs",
			config: &SourceConfig{
				MCPCatalogs: []MCPSource{
					{ID: "dup", Name: "MCP 1", Type: "yaml"},
					{ID: "dup", Name: "MCP 2", Type: "yaml"},
				},
			},
			expectErr: true,
			errMsg:    "duplicate MCP catalog id: dup",
		},
		{
			name: "missing model catalog ID",
			config: &SourceConfig{
				ModelCatalogs: []ModelSource{
					{CatalogSource: apimodels.CatalogSource{Name: "Model 1"}},
				},
			},
			expectErr: true,
			errMsg:    "model catalog source missing id",
		},
		{
			name: "missing MCP catalog ID",
			config: &SourceConfig{
				MCPCatalogs: []MCPSource{
					{Name: "MCP 1", Type: "yaml"},
				},
			},
			expectErr: true,
			errMsg:    "MCP catalog source missing id",
		},
		{
			name: "cross-type ID collision between model and MCP catalogs",
			config: &SourceConfig{
				ModelCatalogs: []ModelSource{
					{CatalogSource: apimodels.CatalogSource{Id: "shared", Name: "Model 1"}},
				},
				MCPCatalogs: []MCPSource{
					{ID: "shared", Name: "MCP 1", Type: "yaml"},
				},
			},
			expectErr: true,
			errMsg:    `id "shared" used in both model_catalogs and mcp_catalogs`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectErr {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
