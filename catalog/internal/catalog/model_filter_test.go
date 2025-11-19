package catalog

import (
	"testing"

	apimodels "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelFilterAllows(t *testing.T) {
	filter, err := NewModelFilter([]string{"Granite/*"}, []string{"Granite/beta-*"})
	require.NoError(t, err)

	assert.True(t, filter.Allows("Granite/3-1-instruct"))
	assert.False(t, filter.Allows("Granite/beta-release"))
	assert.False(t, filter.Allows("Other/model"))

	allowAll, err := NewModelFilter([]string{"*"}, nil)
	require.NoError(t, err)
	assert.True(t, allowAll.Allows("anything/goes"))
}

func TestModelFilterConflictsAndValidation(t *testing.T) {
	_, err := NewModelFilter([]string{"Granite/*"}, []string{"Granite/*"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pattern \"Granite/*\"")

	_, err = NewModelFilter([]string{""}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pattern cannot be empty")
}

func TestNewModelFilterFromSourceMergesLegacy(t *testing.T) {
	source := &Source{
		CatalogSource: apimodels.CatalogSource{
			Id:             "test",
			Name:           "Test source",
			Labels:         []string{},
			IncludedModels: []string{"Granite/*"},
		},
	}

	filter, err := NewModelFilterFromSource(source, nil, []string{"Legacy/*"})
	require.NoError(t, err)

	assert.True(t, filter.Allows("Granite/model"))
	assert.False(t, filter.Allows("Legacy/model"))
}

func TestModelFilterWithWildcardInMiddle(t *testing.T) {
	// Test that wildcards match across the entire name
	filter, err := NewModelFilter(nil, []string{"*deprecated*", "*old*"})
	require.NoError(t, err)

	assert.True(t, filter.Allows("Granite/empty-stable"))
	assert.False(t, filter.Allows("Mistral/empty-deprecated"))
	assert.False(t, filter.Allows("DeepSeek/empty-old-v1"))
	assert.False(t, filter.Allows("Foo/old"))
	assert.False(t, filter.Allows("Bar/deprecated"))
	
	// Test that */pattern* requires the pattern immediately after /
	filter2, err := NewModelFilter(nil, []string{"*/deprecated", "*/old*"})
	require.NoError(t, err)
	
	assert.True(t, filter2.Allows("Mistral/empty-deprecated")) // doesn't match */deprecated (no immediate match after /)
	assert.False(t, filter2.Allows("Foo/deprecated")) // matches */deprecated
	assert.False(t, filter2.Allows("Bar/old-model")) // matches */old*
}
