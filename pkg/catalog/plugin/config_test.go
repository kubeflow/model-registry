package plugin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	yaml := `
apiVersion: catalog/v1alpha1
kind: CatalogSources
catalogs:
  models:
    sources:
      - id: "source-1"
        name: "Source One"
        type: "yaml"
        enabled: true
        properties:
          yamlCatalogPath: "./data/models.yaml"
      - id: "source-2"
        name: "Source Two"
        type: "hf"
        properties:
          allowedOrganization: "redhat"
    labels:
      - name: "production"
        color: "green"
    namedQueries:
      large-models:
        parameters:
          operator: "gt"
          value: 7000000000
  datasets:
    sources:
      - id: "internal-datasets"
        type: "yaml"
        properties:
          yamlCatalogPath: "./data/datasets.yaml"
`

	cfg, err := ParseConfig([]byte(yaml), "/test/path/sources.yaml")
	require.NoError(t, err)

	assert.Equal(t, "catalog/v1alpha1", cfg.APIVersion)
	assert.Equal(t, "CatalogSources", cfg.Kind)

	// Check models catalog
	models, ok := cfg.Catalogs["models"]
	assert.True(t, ok)
	assert.Equal(t, 2, len(models.Sources))
	assert.Equal(t, "source-1", models.Sources[0].ID)
	assert.Equal(t, "Source One", models.Sources[0].Name)
	assert.Equal(t, "yaml", models.Sources[0].Type)
	assert.True(t, *models.Sources[0].Enabled)
	assert.Equal(t, "/test/path/sources.yaml", models.Sources[0].Origin)

	assert.Equal(t, 1, len(models.Labels))
	assert.Equal(t, "production", models.Labels[0]["name"])

	assert.NotNil(t, models.NamedQueries)
	assert.Contains(t, models.NamedQueries, "large-models")

	// Check datasets catalog
	datasets, ok := cfg.Catalogs["datasets"]
	assert.True(t, ok)
	assert.Equal(t, 1, len(datasets.Sources))
}

func TestLoadConfig(t *testing.T) {
	// Create a temp config file
	yaml := `
apiVersion: catalog/v1alpha1
kind: CatalogSources
catalogs:
  models:
    sources:
      - id: "test-source"
        name: "Test Source"
        type: "yaml"
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "sources.yaml")
	err := os.WriteFile(configPath, []byte(yaml), 0644)
	require.NoError(t, err)

	cfg, err := LoadConfig(configPath)
	require.NoError(t, err)

	assert.Equal(t, "catalog/v1alpha1", cfg.APIVersion)
	assert.Contains(t, cfg.Catalogs, "models")
	assert.Equal(t, configPath, cfg.Catalogs["models"].Sources[0].Origin)
}

func TestMergeConfigs(t *testing.T) {
	base := &CatalogSourcesConfig{
		APIVersion: "v1",
		Kind:       "CatalogSources",
		Catalogs: map[string]CatalogSection{
			"models": {
				Sources: []SourceConfig{
					{ID: "source-1", Name: "Base Source", Type: "yaml", Enabled: boolPtr(false)},
				},
			},
		},
	}

	override := &CatalogSourcesConfig{
		Catalogs: map[string]CatalogSection{
			"models": {
				Sources: []SourceConfig{
					{ID: "source-1", Enabled: boolPtr(true)}, // Enable the source
					{ID: "source-2", Name: "New Source", Type: "hf"},
				},
			},
			"datasets": {
				Sources: []SourceConfig{
					{ID: "ds-1", Name: "Dataset Source", Type: "yaml"},
				},
			},
		},
	}

	result := MergeConfigs(base, override)

	assert.Equal(t, "v1", result.APIVersion)

	// Models should have merged sources
	models := result.Catalogs["models"]
	assert.Equal(t, 2, len(models.Sources))

	// Find source-1 and verify it was merged
	var source1 *SourceConfig
	for i := range models.Sources {
		if models.Sources[i].ID == "source-1" {
			source1 = &models.Sources[i]
			break
		}
	}
	require.NotNil(t, source1)
	assert.Equal(t, "Base Source", source1.Name) // Inherited from base
	assert.Equal(t, "yaml", source1.Type)        // Inherited from base
	assert.True(t, *source1.Enabled)             // Overridden

	// Datasets should be added
	datasets := result.Catalogs["datasets"]
	assert.Equal(t, 1, len(datasets.Sources))
	assert.Equal(t, "ds-1", datasets.Sources[0].ID)
}

func TestSourceConfigIsEnabled(t *testing.T) {
	// Default (nil) should be enabled
	s := SourceConfig{}
	assert.True(t, s.IsEnabled())

	// Explicitly enabled
	s.Enabled = boolPtr(true)
	assert.True(t, s.IsEnabled())

	// Explicitly disabled
	s.Enabled = boolPtr(false)
	assert.False(t, s.IsEnabled())
}

func boolPtr(b bool) *bool {
	return &b
}
