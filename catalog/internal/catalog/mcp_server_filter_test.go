package catalog

import (
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMcpServerFilterAllows(t *testing.T) {
	filter, err := NewMcpServerFilter([]string{"github/*"}, []string{"github/beta-*"})
	require.NoError(t, err)

	assert.True(t, filter.Allows("github/filesystem"))
	assert.True(t, filter.Allows("github/fetch"))
	assert.False(t, filter.Allows("github/beta-server"))
	assert.False(t, filter.Allows("other/mcp-server"))

	// Test case-insensitive matching
	assert.True(t, filter.Allows("GitHub/filesystem"))
	assert.True(t, filter.Allows("GITHUB/filesystem"))
	assert.False(t, filter.Allows("GITHUB/BETA-server"))

	allowAll, err := NewMcpServerFilter([]string{"*"}, nil)
	require.NoError(t, err)
	assert.True(t, allowAll.Allows("anything/goes"))
}

func TestMcpServerFilterNilReturnsNil(t *testing.T) {
	// When no filters are provided, NewMcpServerFilter returns nil, nil
	filter, err := NewMcpServerFilter(nil, nil)
	require.NoError(t, err)
	assert.Nil(t, filter)

	// A nil filter should allow everything
	assert.True(t, (*McpServerFilter)(nil).Allows("any-server"))
}

func TestMcpServerFilterExcludeOnly(t *testing.T) {
	// Test with only exclusions - should allow all except excluded patterns
	filter, err := NewMcpServerFilter(nil, []string{"*-deprecated", "*-legacy"})
	require.NoError(t, err)

	assert.True(t, filter.Allows("github-filesystem"))
	assert.True(t, filter.Allows("slack-mcp"))
	assert.False(t, filter.Allows("old-server-deprecated"))
	assert.False(t, filter.Allows("v1-legacy"))
}

func TestMcpServerFilterIncludeOnly(t *testing.T) {
	// Test with only inclusions - should only allow matching patterns
	filter, err := NewMcpServerFilter([]string{"github-*", "slack-*"}, nil)
	require.NoError(t, err)

	assert.True(t, filter.Allows("github-filesystem"))
	assert.True(t, filter.Allows("slack-mcp"))
	assert.False(t, filter.Allows("other-server"))
	assert.False(t, filter.Allows("custom-mcp"))
}

func TestMcpServerFilterConflictsAndValidation(t *testing.T) {
	_, err := NewMcpServerFilter([]string{"github/*"}, []string{"github/*"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pattern \"github/*\"")

	_, err = NewMcpServerFilter([]string{""}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pattern cannot be empty")
}

func TestMcpServerFilterFromSource(t *testing.T) {
	t.Run("source with filters", func(t *testing.T) {
		source := &mcp.McpSource{
			Id:              "test-source",
			Name:            "Test MCP Source",
			IncludedServers: []string{"github-*", "slack-*"},
			ExcludedServers: []string{"*-deprecated"},
		}

		filter, err := NewMcpServerFilterFromSource(source)
		require.NoError(t, err)
		require.NotNil(t, filter)

		assert.True(t, filter.Allows("github-filesystem"))
		assert.True(t, filter.Allows("slack-mcp"))
		assert.False(t, filter.Allows("other-server"))
		assert.False(t, filter.Allows("github-deprecated"))
	})

	t.Run("source without filters", func(t *testing.T) {
		source := &mcp.McpSource{
			Id:   "test-source",
			Name: "Test MCP Source",
		}

		filter, err := NewMcpServerFilterFromSource(source)
		require.NoError(t, err)
		assert.Nil(t, filter) // No filters means nil filter (allow all)
	})

	t.Run("nil source", func(t *testing.T) {
		_, err := NewMcpServerFilterFromSource(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "source cannot be nil")
	})

	t.Run("source with invalid patterns", func(t *testing.T) {
		source := &mcp.McpSource{
			Id:              "test-source",
			Name:            "Test MCP Source",
			IncludedServers: []string{""},
		}

		_, err := NewMcpServerFilterFromSource(source)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid include/exclude configuration")
	})
}

func TestMcpServerFilterWithWildcardInMiddle(t *testing.T) {
	// Test that wildcards match across the entire name
	filter, err := NewMcpServerFilter(nil, []string{"*deprecated*", "*old*"})
	require.NoError(t, err)

	assert.True(t, filter.Allows("github-filesystem"))
	assert.False(t, filter.Allows("server-deprecated-v1"))
	assert.False(t, filter.Allows("old-server"))
	assert.False(t, filter.Allows("my-old-mcp"))

	// Test that */pattern* requires the pattern immediately after /
	filter2, err := NewMcpServerFilter(nil, []string{"*/deprecated", "*/old*"})
	require.NoError(t, err)

	assert.True(t, filter2.Allows("foo/bar-deprecated")) // doesn't match */deprecated
	assert.False(t, filter2.Allows("foo/deprecated"))    // matches */deprecated
	assert.False(t, filter2.Allows("bar/old-server"))    // matches */old*
}

func TestValidateMcpServerSourceFilters(t *testing.T) {
	t.Run("no filters", func(t *testing.T) {
		err := ValidateMcpServerSourceFilters(nil, nil)
		assert.NoError(t, err)
	})

	t.Run("valid patterns", func(t *testing.T) {
		err := ValidateMcpServerSourceFilters([]string{"github-*", "slack-*"}, []string{"*-beta"})
		assert.NoError(t, err)
	})

	t.Run("conflicting patterns", func(t *testing.T) {
		err := ValidateMcpServerSourceFilters([]string{"github-*"}, []string{"github-*"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "github-*")
	})

	t.Run("empty pattern in includedServers", func(t *testing.T) {
		err := ValidateMcpServerSourceFilters([]string{"github-*", ""}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pattern cannot be empty")
	})

	t.Run("whitespace-only pattern", func(t *testing.T) {
		err := ValidateMcpServerSourceFilters([]string{"   "}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pattern cannot be empty")
	})

	t.Run("valid glob patterns", func(t *testing.T) {
		err := ValidateMcpServerSourceFilters([]string{"valid-*"}, nil)
		assert.NoError(t, err)
	})
}
