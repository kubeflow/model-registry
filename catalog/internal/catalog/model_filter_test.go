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

	// Test case-insensitive matching
	assert.True(t, filter.Allows("granite/3-1-instruct"))
	assert.True(t, filter.Allows("GRANITE/3-1-instruct"))
	assert.False(t, filter.Allows("granite/beta-release"))

	allowAll, err := NewModelFilter([]string{"*"}, nil)
	require.NoError(t, err)
	assert.True(t, allowAll.Allows("anything/goes"))
}

func TestModelFilterConflictsAndValidation(t *testing.T) {
	_, err := NewModelFilter([]string{"meta-llama/Llama-3.2-1B"}, []string{"meta-llama/Llama-3.2-1B"})
	require.NoError(t, err)

	_, err = NewModelFilter([]string{"Granite/*"}, []string{"Granite/*"})
	require.NoError(t, err)

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
	assert.False(t, filter2.Allows("Foo/deprecated"))          // matches */deprecated
	assert.False(t, filter2.Allows("Bar/old-model"))           // matches */old*
}

func TestValidateSourceFilters(t *testing.T) {
	t.Run("no filters", func(t *testing.T) {
		err := ValidateSourceFilters(nil, nil)
		assert.NoError(t, err)
	})

	t.Run("valid patterns", func(t *testing.T) {
		err := ValidateSourceFilters([]string{"Granite/*", "Meta/*"}, []string{"*-beta"})
		assert.NoError(t, err)
	})

	t.Run("conflicting patterns", func(t *testing.T) {
		err := ValidateSourceFilters([]string{"Granite/*"}, []string{"Granite/*"})
		require.NoError(t, err)
	})

	t.Run("conflicting patterns", func(t *testing.T) {
		err := ValidateSourceFilters([]string{"Granite/model-2"}, []string{"Granite/model-2"})
		require.NoError(t, err)
	})

	t.Run("empty pattern in includedModels", func(t *testing.T) {
		err := ValidateSourceFilters([]string{"Granite/*", ""}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pattern cannot be empty")
	})

	t.Run("whitespace-only pattern", func(t *testing.T) {
		err := ValidateSourceFilters([]string{"   "}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pattern cannot be empty")
	})

	t.Run("valid glob patterns", func(t *testing.T) {
		err := ValidateSourceFilters([]string{"valid/*"}, nil)
		assert.NoError(t, err) // Our conversion always produces valid regex
	})
}
