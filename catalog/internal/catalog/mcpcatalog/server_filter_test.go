package mcpcatalog

import (
	"testing"

	"github.com/kubeflow/hub/catalog/internal/catalog/basecatalog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServerFilterFromSource_NilSource(t *testing.T) {
	_, err := NewServerFilterFromSource(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source cannot be nil")
}

func TestNewServerFilterFromSource_NoFilters(t *testing.T) {
	source := &basecatalog.MCPSource{ID: "test"}
	filter, err := NewServerFilterFromSource(source)
	require.NoError(t, err)
	assert.Nil(t, filter, "nil filter expected when no patterns are set")
}

func TestNewServerFilterFromSource_IncludeOnly(t *testing.T) {
	source := &basecatalog.MCPSource{
		ID:              "test",
		IncludedServers: []string{"github*", "slack*"},
	}
	filter, err := NewServerFilterFromSource(source)
	require.NoError(t, err)
	require.NotNil(t, filter)

	assert.True(t, filter.Allows("github-mcp"))
	assert.True(t, filter.Allows("slack-mcp"))
	assert.False(t, filter.Allows("jira-mcp"))
}

func TestNewServerFilterFromSource_ExcludeOnly(t *testing.T) {
	source := &basecatalog.MCPSource{
		ID:              "test",
		ExcludedServers: []string{"internal*"},
	}
	filter, err := NewServerFilterFromSource(source)
	require.NoError(t, err)
	require.NotNil(t, filter)

	assert.True(t, filter.Allows("github-mcp"))
	assert.False(t, filter.Allows("internal-mcp"))
	assert.False(t, filter.Allows("internal-tools"))
}

func TestNewServerFilterFromSource_IncludeAndExclude(t *testing.T) {
	source := &basecatalog.MCPSource{
		ID:              "test",
		IncludedServers: []string{"github*"},
		ExcludedServers: []string{"github-internal"},
	}
	filter, err := NewServerFilterFromSource(source)
	require.NoError(t, err)
	require.NotNil(t, filter)

	assert.True(t, filter.Allows("github-public"))
	assert.False(t, filter.Allows("github-internal"))
	assert.False(t, filter.Allows("slack-mcp"))
}

func TestNewServerFilterFromSource_InvalidPattern(t *testing.T) {
	source := &basecatalog.MCPSource{
		ID:              "test",
		IncludedServers: []string{""},
	}
	_, err := NewServerFilterFromSource(source)
	require.Error(t, err)
}

func TestNewServerFilter_NilSafe(t *testing.T) {
	// A nil *ServerFilter must allow everything (NameFilter.Allows is nil-safe)
	var filter *ServerFilter
	assert.True(t, filter.Allows("anything"))
}

func TestValidateServerFilters_Valid(t *testing.T) {
	err := ValidateServerFilters([]string{"github*"}, []string{"*internal"})
	require.NoError(t, err)
}

func TestValidateServerFilters_EmptyPattern(t *testing.T) {
	err := ValidateServerFilters([]string{""}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "includedServers")
}
