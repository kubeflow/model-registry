package catalog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYamlModelToModelProviderRecord(t *testing.T) {
	tests := []struct {
		name         string
		yamlModel    yamlModel
		expectError  bool
		validateFunc func(t *testing.T, record ModelProviderRecord)
	}{
		{
			name: "complete model with all properties",
			yamlModel: yamlModel{
				CatalogModel: model.CatalogModel{
					Name:                     "test-model",
					Description:              apiutils.Of("Test model description"),
					Readme:                   apiutils.Of("# Test Model\nThis is a test model."),
					Maturity:                 apiutils.Of("Generally Available"),
					Language:                 []string{"en", "fr"},
					Tasks:                    []string{"text-generation", "nlp"},
					Provider:                 apiutils.Of("IBM"),
					Logo:                     apiutils.Of("https://example.com/logo.png"),
					License:                  apiutils.Of("apache-2.0"),
					LicenseLink:              apiutils.Of("https://www.apache.org/licenses/LICENSE-2.0"),
					LibraryName:              apiutils.Of("transformers"),
					SourceId:                 apiutils.Of("test-source"),
					CreateTimeSinceEpoch:     apiutils.Of("1678886400000"),
					LastUpdateTimeSinceEpoch: apiutils.Of("1681564800000"),
					CustomProperties: &map[string]model.MetadataValue{
						"custom_key": {
							MetadataStringValue: &model.MetadataStringValue{
								StringValue:  "custom_value",
								MetadataType: "MetadataStringValue",
							},
						},
					},
				},
				Artifacts: []*yamlArtifact{
					{
						CatalogArtifact: model.CatalogArtifact{
							CatalogModelArtifact: &model.CatalogModelArtifact{
								ArtifactType:             "model-artifact",
								Uri:                      "https://example.com/model.tar.gz",
								CreateTimeSinceEpoch:     apiutils.Of("1678886400000"),
								LastUpdateTimeSinceEpoch: apiutils.Of("1681564800000"),
								CustomProperties: &map[string]model.MetadataValue{
									"model_size": {
										MetadataStringValue: &model.MetadataStringValue{
											StringValue:  "2GB",
											MetadataType: "MetadataStringValue",
										},
									},
									"accuracy": {
										MetadataDoubleValue: &model.MetadataDoubleValue{
											DoubleValue:  0.95,
											MetadataType: "MetadataDoubleValue",
										},
									},
								},
							},
						},
					},
					{
						CatalogArtifact: model.CatalogArtifact{
							CatalogMetricsArtifact: &model.CatalogMetricsArtifact{
								ArtifactType:             "metrics-artifact",
								MetricsType:              "evaluation-metrics",
								CreateTimeSinceEpoch:     apiutils.Of("1678886400000"),
								LastUpdateTimeSinceEpoch: apiutils.Of("1681564800000"),
								CustomProperties: &map[string]model.MetadataValue{
									"framework": {
										MetadataStringValue: &model.MetadataStringValue{
											StringValue:  "scikit-learn",
											MetadataType: "MetadataStringValue",
										},
									},
								},
							},
						},
					},
				},
			},
			validateFunc: func(t *testing.T, record ModelProviderRecord) {
				require.NotNil(t, record.Model)

				attrs := record.Model.GetAttributes()
				require.NotNil(t, attrs)
				assert.Equal(t, int64(1678886400000), *attrs.CreateTimeSinceEpoch)
				assert.Equal(t, int64(1681564800000), *attrs.LastUpdateTimeSinceEpoch)

				// Check regular properties (spec-defined properties)
				regularProps := record.Model.GetProperties()
				require.NotNil(t, regularProps)

				regularPropMap := make(map[string]models.Properties)
				for _, prop := range *regularProps {
					regularPropMap[prop.Name] = prop
				}

				// Validate spec-defined properties are regular properties
				assert.Contains(t, regularPropMap, "description")
				assert.Equal(t, "Test model description", *regularPropMap["description"].StringValue)
				assert.False(t, regularPropMap["description"].IsCustomProperty)

				assert.Contains(t, regularPropMap, "readme")
				assert.Equal(t, "# Test Model\nThis is a test model.", *regularPropMap["readme"].StringValue)
				assert.False(t, regularPropMap["readme"].IsCustomProperty)

				assert.Contains(t, regularPropMap, "maturity")
				assert.Equal(t, "Generally Available", *regularPropMap["maturity"].StringValue)
				assert.False(t, regularPropMap["maturity"].IsCustomProperty)

				assert.Contains(t, regularPropMap, "provider")
				assert.Equal(t, "IBM", *regularPropMap["provider"].StringValue)
				assert.False(t, regularPropMap["provider"].IsCustomProperty)

				assert.Contains(t, regularPropMap, "logo")
				assert.Equal(t, "https://example.com/logo.png", *regularPropMap["logo"].StringValue)
				assert.False(t, regularPropMap["logo"].IsCustomProperty)

				assert.Contains(t, regularPropMap, "license")
				assert.Equal(t, "apache-2.0", *regularPropMap["license"].StringValue)
				assert.False(t, regularPropMap["license"].IsCustomProperty)

				assert.Contains(t, regularPropMap, "license_link")
				assert.Equal(t, "https://www.apache.org/licenses/LICENSE-2.0", *regularPropMap["license_link"].StringValue)
				assert.False(t, regularPropMap["license_link"].IsCustomProperty)

				assert.Contains(t, regularPropMap, "library_name")
				assert.Equal(t, "transformers", *regularPropMap["library_name"].StringValue)
				assert.False(t, regularPropMap["library_name"].IsCustomProperty)

				assert.Contains(t, regularPropMap, "source_id")
				assert.Equal(t, "test-source", *regularPropMap["source_id"].StringValue)
				assert.False(t, regularPropMap["source_id"].IsCustomProperty)

				// Validate array properties are JSON encoded as regular properties
				assert.Contains(t, regularPropMap, "language")
				var languages []string
				err := json.Unmarshal([]byte(*regularPropMap["language"].StringValue), &languages)
				require.NoError(t, err)
				assert.Equal(t, []string{"en", "fr"}, languages)
				assert.False(t, regularPropMap["language"].IsCustomProperty)

				assert.Contains(t, regularPropMap, "tasks")
				var tasks []string
				err = json.Unmarshal([]byte(*regularPropMap["tasks"].StringValue), &tasks)
				require.NoError(t, err)
				assert.Equal(t, []string{"text-generation", "nlp"}, tasks)
				assert.False(t, regularPropMap["tasks"].IsCustomProperty)

				// Check custom properties
				customProps := record.Model.GetCustomProperties()
				require.NotNil(t, customProps)

				customPropMap := make(map[string]models.Properties)
				for _, prop := range *customProps {
					customPropMap[prop.Name] = prop
				}

				// Validate truly custom properties
				assert.Contains(t, customPropMap, "custom_key")
				assert.Equal(t, "custom_value", *customPropMap["custom_key"].StringValue)
				assert.True(t, customPropMap["custom_key"].IsCustomProperty)

				// Validate artifacts
				assert.Len(t, record.Artifacts, 2)

				// Validate ModelArtifact
				modelArtifact := record.Artifacts[0]
				require.NotNil(t, modelArtifact.CatalogModelArtifact)
				assert.Nil(t, modelArtifact.CatalogMetricsArtifact)

				// Check CatalogModelArtifact attributes
				modelAttrs := modelArtifact.CatalogModelArtifact.GetAttributes()
				assert.Equal(t, "https://example.com/model.tar.gz", *modelAttrs.URI)
				assert.Equal(t, int64(1678886400000), *modelAttrs.CreateTimeSinceEpoch)
				assert.Equal(t, int64(1681564800000), *modelAttrs.LastUpdateTimeSinceEpoch)

				// Check CatalogModelArtifact regular properties
				modelArtifactProps := modelArtifact.CatalogModelArtifact.GetProperties()
				require.NotNil(t, modelArtifactProps)
				modelArtifactPropMap := make(map[string]models.Properties)
				for _, prop := range *modelArtifactProps {
					modelArtifactPropMap[prop.Name] = prop
				}
				assert.Contains(t, modelArtifactPropMap, "uri")
				assert.Equal(t, "https://example.com/model.tar.gz", *modelArtifactPropMap["uri"].StringValue)
				assert.False(t, modelArtifactPropMap["uri"].IsCustomProperty)

				// Check CatalogModelArtifact custom properties
				modelArtifactCustomProps := modelArtifact.CatalogModelArtifact.GetCustomProperties()
				require.NotNil(t, modelArtifactCustomProps)
				modelArtifactCustomPropMap := make(map[string]models.Properties)
				for _, prop := range *modelArtifactCustomProps {
					modelArtifactCustomPropMap[prop.Name] = prop
				}
				assert.Contains(t, modelArtifactCustomPropMap, "model_size")
				assert.Equal(t, "2GB", *modelArtifactCustomPropMap["model_size"].StringValue)
				assert.True(t, modelArtifactCustomPropMap["model_size"].IsCustomProperty)
				assert.Contains(t, modelArtifactCustomPropMap, "accuracy")
				assert.Equal(t, 0.95, *modelArtifactCustomPropMap["accuracy"].DoubleValue)
				assert.True(t, modelArtifactCustomPropMap["accuracy"].IsCustomProperty)

				// Validate CatalogMetricsArtifact
				metricsArtifact := record.Artifacts[1]
				require.NotNil(t, metricsArtifact.CatalogMetricsArtifact)
				assert.Nil(t, metricsArtifact.CatalogModelArtifact)

				// Check CatalogMetricsArtifact attributes
				metricsAttrs := metricsArtifact.CatalogMetricsArtifact.GetAttributes()
				assert.Equal(t, "evaluation-metrics", string(metricsAttrs.MetricsType))
				assert.Equal(t, int64(1678886400000), *metricsAttrs.CreateTimeSinceEpoch)
				assert.Equal(t, int64(1681564800000), *metricsAttrs.LastUpdateTimeSinceEpoch)

				// Check CatalogMetricsArtifact regular properties
				metricsArtifactProps := metricsArtifact.CatalogMetricsArtifact.GetProperties()
				require.NotNil(t, metricsArtifactProps)
				metricsArtifactPropMap := make(map[string]models.Properties)
				for _, prop := range *metricsArtifactProps {
					metricsArtifactPropMap[prop.Name] = prop
				}
				assert.Contains(t, metricsArtifactPropMap, "metricsType")
				assert.Equal(t, "evaluation-metrics", *metricsArtifactPropMap["metricsType"].StringValue)
				assert.False(t, metricsArtifactPropMap["metricsType"].IsCustomProperty)

				// Check CatalogMetricsArtifact custom properties
				metricsArtifactCustomProps := metricsArtifact.CatalogMetricsArtifact.GetCustomProperties()
				require.NotNil(t, metricsArtifactCustomProps)
				metricsArtifactCustomPropMap := make(map[string]models.Properties)
				for _, prop := range *metricsArtifactCustomProps {
					metricsArtifactCustomPropMap[prop.Name] = prop
				}
				assert.Contains(t, metricsArtifactCustomPropMap, "framework")
				assert.Equal(t, "scikit-learn", *metricsArtifactCustomPropMap["framework"].StringValue)
				assert.True(t, metricsArtifactCustomPropMap["framework"].IsCustomProperty)
			},
		},
		{
			name: "minimal model with only required fields",
			yamlModel: yamlModel{
				CatalogModel: model.CatalogModel{
					Name: "minimal-model",
				},
			},
			validateFunc: func(t *testing.T, record ModelProviderRecord) {
				require.NotNil(t, record.Model)

				attrs := record.Model.GetAttributes()
				require.NotNil(t, attrs)
				assert.Nil(t, attrs.CreateTimeSinceEpoch)
				assert.Nil(t, attrs.LastUpdateTimeSinceEpoch)

				// Should have no regular properties for minimal model
				regularProps := record.Model.GetProperties()
				if regularProps != nil {
					*regularProps = slices.DeleteFunc(*regularProps, func(p models.Properties) bool {
						switch p.Name {
						case "language", "tasks":
							return true
						}
						return false
					})
					assert.Empty(t, *regularProps)
				}

				// Should have no custom properties for minimal model
				customProps := record.Model.GetCustomProperties()
				if customProps != nil {
					assert.Empty(t, *customProps)
				}

				// Should have no artifacts for minimal model
				assert.Empty(t, record.Artifacts)
			},
		},
		{
			name: "model with only ModelArtifact",
			yamlModel: yamlModel{
				CatalogModel: model.CatalogModel{
					Name: "model-with-artifact",
				},
				Artifacts: []*yamlArtifact{
					{
						CatalogArtifact: model.CatalogArtifact{
							CatalogModelArtifact: &model.CatalogModelArtifact{
								ArtifactType: "model-artifact",
								Uri:          "s3://bucket/model.bin",
							},
						},
					},
				},
			},
			validateFunc: func(t *testing.T, record ModelProviderRecord) {
				require.NotNil(t, record.Model)
				assert.Len(t, record.Artifacts, 1)

				artifact := record.Artifacts[0]
				require.NotNil(t, artifact.CatalogModelArtifact)
				assert.Nil(t, artifact.CatalogMetricsArtifact)

				attrs := artifact.CatalogModelArtifact.GetAttributes()
				assert.Equal(t, "s3://bucket/model.bin", *attrs.URI)
				assert.Nil(t, attrs.CreateTimeSinceEpoch)
				assert.Nil(t, attrs.LastUpdateTimeSinceEpoch)

				// Check regular properties
				props := artifact.CatalogModelArtifact.GetProperties()
				require.NotNil(t, props)
				assert.Len(t, *props, 1)
				assert.Equal(t, "uri", (*props)[0].Name)
				assert.Equal(t, "s3://bucket/model.bin", *(*props)[0].StringValue)
				assert.False(t, (*props)[0].IsCustomProperty)

				// Should have no custom properties
				customProps := artifact.CatalogModelArtifact.GetCustomProperties()
				if customProps != nil {
					assert.Empty(t, *customProps)
				}
			},
		},
		{
			name: "model with only MetricsArtifact",
			yamlModel: yamlModel{
				CatalogModel: model.CatalogModel{
					Name: "model-with-metrics",
				},
				Artifacts: []*yamlArtifact{
					{
						CatalogArtifact: model.CatalogArtifact{
							CatalogMetricsArtifact: &model.CatalogMetricsArtifact{
								ArtifactType: "metrics-artifact",
								MetricsType:  "performance-metrics",
							},
						},
					},
				},
			},
			validateFunc: func(t *testing.T, record ModelProviderRecord) {
				require.NotNil(t, record.Model)
				assert.Len(t, record.Artifacts, 1)

				artifact := record.Artifacts[0]
				assert.Nil(t, artifact.CatalogModelArtifact)
				require.NotNil(t, artifact.CatalogMetricsArtifact)

				attrs := artifact.CatalogMetricsArtifact.GetAttributes()
				assert.Equal(t, "performance-metrics", string(attrs.MetricsType))
				assert.Nil(t, attrs.CreateTimeSinceEpoch)
				assert.Nil(t, attrs.LastUpdateTimeSinceEpoch)

				// Check regular properties
				props := artifact.CatalogMetricsArtifact.GetProperties()
				require.NotNil(t, props)
				assert.Len(t, *props, 1)
				assert.Equal(t, "metricsType", (*props)[0].Name)
				assert.Equal(t, "performance-metrics", *(*props)[0].StringValue)
				assert.False(t, (*props)[0].IsCustomProperty)

				// Should have no custom properties
				customProps := artifact.CatalogMetricsArtifact.GetCustomProperties()
				if customProps != nil {
					assert.Empty(t, *customProps)
				}
			},
		},
		{
			name: "artifacts with invalid timestamps",
			yamlModel: yamlModel{
				CatalogModel: model.CatalogModel{
					Name: "model-with-invalid-artifact-timestamps",
				},
				Artifacts: []*yamlArtifact{
					{
						CatalogArtifact: model.CatalogArtifact{
							CatalogModelArtifact: &model.CatalogModelArtifact{
								ArtifactType:             "model-artifact",
								Uri:                      "https://example.com/model.bin",
								CreateTimeSinceEpoch:     apiutils.Of("invalid-timestamp"),
								LastUpdateTimeSinceEpoch: apiutils.Of("also-invalid"),
							},
						},
					},
				},
			},
			validateFunc: func(t *testing.T, record ModelProviderRecord) {
				require.NotNil(t, record.Model)
				assert.Len(t, record.Artifacts, 1)

				artifact := record.Artifacts[0]
				require.NotNil(t, artifact.CatalogModelArtifact)

				attrs := artifact.CatalogModelArtifact.GetAttributes()
				assert.Equal(t, "https://example.com/model.bin", *attrs.URI)
				// Invalid timestamps should be ignored (not set)
				assert.Nil(t, attrs.CreateTimeSinceEpoch)
				assert.Nil(t, attrs.LastUpdateTimeSinceEpoch)
			},
		},
		{
			name: "model with invalid timestamps",
			yamlModel: yamlModel{
				CatalogModel: model.CatalogModel{
					Name:                     "invalid-timestamp-model",
					CreateTimeSinceEpoch:     apiutils.Of("invalid-timestamp"),
					LastUpdateTimeSinceEpoch: apiutils.Of("also-invalid"),
				},
			},
			validateFunc: func(t *testing.T, record ModelProviderRecord) {
				require.NotNil(t, record.Model)

				attrs := record.Model.GetAttributes()
				require.NotNil(t, attrs)
				assert.Equal(t, "invalid-timestamp-model", *attrs.Name)
				// Invalid timestamps should be ignored (not set)
				assert.Nil(t, attrs.CreateTimeSinceEpoch)
				assert.Nil(t, attrs.LastUpdateTimeSinceEpoch)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := tt.yamlModel.ToModelProviderRecord()
			tt.validateFunc(t, record)
		})
	}
}

func TestNewYamlModelProviderAbsolutePath(t *testing.T) {
	// Create a temporary YAML file
	tempDir, err := os.MkdirTemp("", "yaml_catalog_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	yamlContent := `
source: test
models:
  - name: test-model
    description: A test model
`

	absolutePath := filepath.Join(tempDir, "catalog.yaml")
	err = os.WriteFile(absolutePath, []byte(yamlContent), 0644)
	require.NoError(t, err)

	// Test that absolute paths work correctly (bug fix)
	t.Run("absolute path works correctly", func(t *testing.T) {
		source := &Source{
			Properties: map[string]interface{}{
				yamlCatalogPathKey: absolutePath, // Use absolute path to existing file
			},
		}

		// Use a different reldir - this shouldn't matter for absolute paths
		reldir := "/some/other/directory"

		_, err := newYamlModelProvider(t.Context(), source, reldir)

		// This should succeed because we provided a valid absolute path
		require.NoError(t, err, "Absolute path should work without reldir being prepended")
	})

	// Test that relative paths work correctly
	t.Run("relative path works correctly", func(t *testing.T) {
		// Create the expected structure for relative path
		reldir := tempDir
		relativePath := "catalog.yaml"

		source := &Source{
			Properties: map[string]interface{}{
				yamlCatalogPathKey: relativePath,
			},
		}

		_, err := newYamlModelProvider(t.Context(), source, reldir)

		// This should work because filepath.Join(reldir, relativePath) points to our file
		require.NoError(t, err, "Relative paths should work correctly")
	})
}
