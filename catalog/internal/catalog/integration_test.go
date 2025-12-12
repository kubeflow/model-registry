package catalog

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNamedQueriesEndToEnd(t *testing.T) {
	// Create temporary YAML file with named queries
	tempDir, err := os.MkdirTemp("", "named-queries-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	yamlContent := `
catalogs: []
labels: []
namedQueries:
  validation-default:
    ttft_p90:
      operator: '<'
      value: 70
    workload_type:
      operator: '='
      value: "Chat"
  high-performance:
    performance_score:
      operator: '>'
      value: 0.95
    tpot_mean:
      operator: '<'
      value: 30
`

	yamlPath := filepath.Join(tempDir, "test-sources.yaml")
	err = os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	require.NoError(t, err)

	// Test configuration loading
	loader := NewLoader(service.Services{}, []string{yamlPath})
	err = loader.parseAndMerge(yamlPath)
	require.NoError(t, err)

	// Verify named queries are loaded
	namedQueries := loader.Sources.GetNamedQueries()
	assert.Len(t, namedQueries, 2)

	// Test validation-default query
	validationQuery := namedQueries["validation-default"]
	assert.Equal(t, "<", validationQuery["ttft_p90"].Operator)
	assert.Equal(t, float64(70), validationQuery["ttft_p90"].Value)
	assert.Equal(t, "=", validationQuery["workload_type"].Operator)
	assert.Equal(t, "Chat", validationQuery["workload_type"].Value)

	// Test high-performance query
	perfQuery := namedQueries["high-performance"]
	assert.Equal(t, ">", perfQuery["performance_score"].Operator)
	assert.Equal(t, float64(0.95), perfQuery["performance_score"].Value)
	assert.Equal(t, "<", perfQuery["tpot_mean"].Operator)
	assert.Equal(t, float64(30), perfQuery["tpot_mean"].Value)

	// Test API response includes named queries
	mockServices := service.Services{
		PropertyOptionsRepository: &mockPropertyRepository{},
	}
	catalog := NewDBCatalog(mockServices, loader.Sources)

	filterOptions, err := catalog.GetFilterOptions(context.Background())
	require.NoError(t, err)
	require.NotNil(t, filterOptions.NamedQueries)

	apiQueries := *filterOptions.NamedQueries
	assert.Len(t, apiQueries, 2)
	assert.Contains(t, apiQueries, "validation-default")
	assert.Contains(t, apiQueries, "high-performance")
}

func TestNamedQueriesValidationErrors(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "named-queries-validation-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test invalid operator
	invalidYaml := `
catalogs: []
namedQueries:
  bad-query:
    field1:
      operator: 'INVALID_OP'
      value: 100
`

	yamlPath := filepath.Join(tempDir, "invalid-sources.yaml")
	err = os.WriteFile(yamlPath, []byte(invalidYaml), 0644)
	require.NoError(t, err)

	loader := NewLoader(service.Services{}, []string{yamlPath})
	err = loader.parseAndMerge(yamlPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported operator 'INVALID_OP'")
}
