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
	assert.False(t, filter2.Allows("Foo/deprecated"))          // matches */deprecated
	assert.False(t, filter2.Allows("Bar/old-model"))           // matches */old*
}

func TestValidateSourceFilters(t *testing.T) {
	t.Run("valid source with no filters", func(t *testing.T) {
		source := &Source{
			CatalogSource: apimodels.CatalogSource{
				Id:   "test",
				Name: "Test source",
			},
		}
		err := ValidateSourceFilters(source)
		assert.NoError(t, err)
	})

	t.Run("valid source with patterns", func(t *testing.T) {
		source := &Source{
			CatalogSource: apimodels.CatalogSource{
				Id:             "test",
				Name:           "Test source",
				IncludedModels: []string{"Granite/*", "Meta/*"},
				ExcludedModels: []string{"*-beta"},
			},
		}
		err := ValidateSourceFilters(source)
		assert.NoError(t, err)
	})

	t.Run("conflicting patterns", func(t *testing.T) {
		source := &Source{
			CatalogSource: apimodels.CatalogSource{
				Id:             "test",
				Name:           "Test source",
				IncludedModels: []string{"Granite/*"},
				ExcludedModels: []string{"Granite/*"},
			},
		}
		err := ValidateSourceFilters(source)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "source test")
		assert.Contains(t, err.Error(), "Granite/*")
	})

	t.Run("empty pattern in includedModels", func(t *testing.T) {
		source := &Source{
			CatalogSource: apimodels.CatalogSource{
				Id:             "test",
				Name:           "Test source",
				IncludedModels: []string{"Granite/*", ""},
			},
		}
		err := ValidateSourceFilters(source)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "source test")
		assert.Contains(t, err.Error(), "pattern cannot be empty")
	})

	t.Run("invalid regex pattern", func(t *testing.T) {
		// This shouldn't happen with our glob-to-regex conversion,
		// but validates the error path exists
		source := &Source{
			CatalogSource: apimodels.CatalogSource{
				Id:             "test",
				Name:           "Test source",
				IncludedModels: []string{"valid/*"},
			},
		}
		err := ValidateSourceFilters(source)
		assert.NoError(t, err) // Our conversion always produces valid regex
	})

	t.Run("nil source", func(t *testing.T) {
		err := ValidateSourceFilters(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "source cannot be nil")
	})
}
