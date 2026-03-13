package basecatalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNameFilterAllows(t *testing.T) {
	filter, err := NewNameFilter("included", []string{"Granite/*"}, "excluded", []string{"Granite/beta-*"})
	require.NoError(t, err)

	assert.True(t, filter.Allows("Granite/3-1-instruct"))
	assert.False(t, filter.Allows("Granite/beta-release"))
	assert.False(t, filter.Allows("Other/model"))

	// Test case-insensitive matching
	assert.True(t, filter.Allows("granite/3-1-instruct"))
	assert.True(t, filter.Allows("GRANITE/3-1-instruct"))
	assert.False(t, filter.Allows("granite/beta-release"))

	allowAll, err := NewNameFilter("included", []string{"*"}, "excluded", nil)
	require.NoError(t, err)
	assert.True(t, allowAll.Allows("anything/goes"))
}

func TestNameFilterConflictsAndValidation(t *testing.T) {
	_, err := NewNameFilter("included", []string{"meta-llama/Llama-3.2-1B"}, "excluded", []string{"meta-llama/Llama-3.2-1B"})
	require.NoError(t, err)

	_, err = NewNameFilter("included", []string{"Granite/*"}, "excluded", []string{"Granite/*"})
	require.NoError(t, err)

	_, err = NewNameFilter("included", []string{""}, "excluded", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pattern cannot be empty")
}

func TestNameFilterWithWildcardInMiddle(t *testing.T) {
	// Test that wildcards match across the entire name
	filter, err := NewNameFilter("included", nil, "excluded", []string{"*deprecated*", "*old*"})
	require.NoError(t, err)

	assert.True(t, filter.Allows("Granite/empty-stable"))
	assert.False(t, filter.Allows("Mistral/empty-deprecated"))
	assert.False(t, filter.Allows("DeepSeek/empty-old-v1"))
	assert.False(t, filter.Allows("Foo/old"))
	assert.False(t, filter.Allows("Bar/deprecated"))

	// Test that */pattern* requires the pattern immediately after /
	filter2, err := NewNameFilter("included", nil, "excluded", []string{"*/deprecated", "*/old*"})
	require.NoError(t, err)

	assert.True(t, filter2.Allows("Mistral/empty-deprecated")) // doesn't match */deprecated (no immediate match after /)
	assert.False(t, filter2.Allows("Foo/deprecated"))          // matches */deprecated
	assert.False(t, filter2.Allows("Bar/old-model"))           // matches */old*
}

func TestValidatePatterns(t *testing.T) {
	t.Run("no filters", func(t *testing.T) {
		err := ValidatePatterns("includedModels", nil, "excludedModels", nil)
		assert.NoError(t, err)
	})

	t.Run("valid patterns", func(t *testing.T) {
		err := ValidatePatterns("includedModels", []string{"Granite/*", "Meta/*"}, "excludedModels", []string{"*-beta"})
		assert.NoError(t, err)
	})

	t.Run("conflicting patterns", func(t *testing.T) {
		err := ValidatePatterns("includedModels", []string{"Granite/*"}, "excludedModels", []string{"Granite/*"})
		require.NoError(t, err)
	})

	t.Run("conflicting patterns", func(t *testing.T) {
		err := ValidatePatterns("includedModels", []string{"Granite/model-2"}, "excludedModels", []string{"Granite/model-2"})
		require.NoError(t, err)
	})

	t.Run("empty pattern in includedModels", func(t *testing.T) {
		err := ValidatePatterns("includedModels", []string{"Granite/*", ""}, "excludedModels", nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pattern cannot be empty")
	})

	t.Run("whitespace-only pattern", func(t *testing.T) {
		err := ValidatePatterns("includedModels", []string{"   "}, "excludedModels", nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pattern cannot be empty")
	})

	t.Run("valid glob patterns", func(t *testing.T) {
		err := ValidatePatterns("includedModels", []string{"valid/*"}, "excludedModels", nil)
		assert.NoError(t, err) // Our conversion always produces valid regex
	})
}
