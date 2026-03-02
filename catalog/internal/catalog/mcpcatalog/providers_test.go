package mcpcatalog

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYamlMCPServerConversion(t *testing.T) {
	yamlServer := &yamlMCPServer{
		Name:        "test-server",
		Description: strPtr("Test MCP server"),
		Provider:    strPtr("Test Provider"),
		Version:     strPtr("1.0.0"),
		License:     strPtr("MIT"),
		Tools: []*yamlMCPTool{
			{Name: "tool1", Description: strPtr("First tool")},
			{Name: "tool2", Description: strPtr("Second tool")},
		},
	}

	record := yamlServer.ToMCPServerProviderRecord()

	assert.NotNil(t, record.Server)
	assert.Nil(t, record.Error)
	assert.Equal(t, 2, len(record.Tools))
	assert.Equal(t, "tool1", record.Tools[0].Name)
	assert.Equal(t, "tool2", record.Tools[1].Name)

	attrs := record.Server.GetAttributes()
	require.NotNil(t, attrs)
	assert.Equal(t, "test-server", *attrs.Name)

	props := record.Server.GetProperties()
	require.NotNil(t, props)

	// Check that basic properties are set
	foundDescription := false
	foundProvider := false
	for _, prop := range *props {
		if prop.Name == "description" && prop.StringValue != nil {
			assert.Equal(t, "Test MCP server", *prop.StringValue)
			foundDescription = true
		}
		if prop.Name == "provider" && prop.StringValue != nil {
			assert.Equal(t, "Test Provider", *prop.StringValue)
			foundProvider = true
		}
	}
	assert.True(t, foundDescription, "description property should be set")
	assert.True(t, foundProvider, "provider property should be set")
}

func TestNewYamlMCPProviderPaths(t *testing.T) {
	// Test with absolute path
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-catalog.yaml")
	err := os.WriteFile(testFile, []byte("mcp_servers: []\n"), 0644)
	require.NoError(t, err)

	source := basecatalog.MCPSource{
		Type: "yaml",
		Properties: map[string]any{
			yamlMCPCatalogPathKey: testFile,
		},
	}

	provider, err := NewYamlMCPProvider(source)
	require.NoError(t, err)
	assert.NotNil(t, provider)

	yamlProvider, ok := provider.(*yamlMCPProvider)
	require.True(t, ok)
	assert.Equal(t, 1, len(yamlProvider.paths))
	assert.Equal(t, testFile, yamlProvider.paths[0])
}

func TestNewYamlMCPProviderRelativePath(t *testing.T) {
	// Test with relative path - should be resolved relative to the config file's directory
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "sources.yaml")

	source := basecatalog.MCPSource{
		Type: "yaml",
		Properties: map[string]any{
			yamlMCPCatalogPathKey: "servers/test.yaml",
		},
		Origin: configFile, // Absolute path to the config file
	}

	provider, err := NewYamlMCPProvider(source)
	require.NoError(t, err)
	assert.NotNil(t, provider)

	yamlProvider, ok := provider.(*yamlMCPProvider)
	require.True(t, ok)
	assert.Equal(t, 1, len(yamlProvider.paths))

	// Should be converted to absolute path
	assert.True(t, filepath.IsAbs(yamlProvider.paths[0]))

	// Should be resolved relative to the config file's directory
	expectedPath := filepath.Join(tmpDir, "servers/test.yaml")
	assert.Equal(t, expectedPath, yamlProvider.paths[0])
}

func TestNewYamlMCPProviderMissingPath(t *testing.T) {
	source := basecatalog.MCPSource{
		Type:       "yaml",
		Properties: map[string]any{},
	}

	_, err := NewYamlMCPProvider(source)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yamlCatalogPath property is required")
}

func TestYamlMCPProviderEmit(t *testing.T) {
	// Create a temporary YAML file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-catalog.yaml")

	yamlContent := `mcp_servers:
  - name: "test-server-1"
    description: "First test server"
    provider: "Test Provider"
    version: "1.0.0"
    tools:
      - name: "test-tool"
        description: "Test tool"
  - name: "test-server-2"
    description: "Second test server"
    provider: "Test Provider"
    version: "2.0.0"
`
	err := os.WriteFile(testFile, []byte(yamlContent), 0644)
	require.NoError(t, err)

	source := basecatalog.MCPSource{
		Type: "yaml",
		Properties: map[string]any{
			yamlMCPCatalogPathKey: testFile,
		},
	}

	provider, err := NewYamlMCPProvider(source)
	require.NoError(t, err)

	ctx := context.Background()
	recordChan := provider.Servers(ctx)

	servers := make([]string, 0)
	for record := range recordChan {
		assert.Nil(t, record.Error)
		assert.NotNil(t, record.Server)
		attrs := record.Server.GetAttributes()
		require.NotNil(t, attrs)
		servers = append(servers, *attrs.Name)
	}

	assert.Equal(t, 2, len(servers))
	assert.Contains(t, servers, "test-server-1")
	assert.Contains(t, servers, "test-server-2")
}

func TestYamlMCPProviderEmitWithCancellation(t *testing.T) {
	// Create a temporary YAML file with many servers
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-catalog.yaml")

	// Create YAML with multiple servers
	var yamlContent strings.Builder
	yamlContent.WriteString("mcp_servers:\n")
	for i := range 10 {
		yamlContent.WriteString(`  - name: "test-server-` + string(rune('0'+i)) + `"
    description: "Test server"
`)
	}
	err := os.WriteFile(testFile, []byte(yamlContent.String()), 0644)
	require.NoError(t, err)

	source := basecatalog.MCPSource{
		Type: "yaml",
		Properties: map[string]any{
			yamlMCPCatalogPathKey: testFile,
		},
	}

	provider, err := NewYamlMCPProvider(source)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	recordChan := provider.Servers(ctx)

	// Cancel after receiving first record
	recordCount := 0
	for range recordChan {
		recordCount++
		if recordCount == 1 {
			cancel()
		}
	}

	// Should have received at least 1 record before cancellation
	assert.GreaterOrEqual(t, recordCount, 1)
}

func TestRegisterMCPProvider(t *testing.T) {
	// Test registering a new provider
	testProviderFunc := func(source basecatalog.MCPSource) (MCPProvider, error) {
		return nil, nil
	}

	err := RegisterMCPProvider("test-provider", testProviderFunc)
	assert.NoError(t, err)
	t.Cleanup(func() { unregisterMCPProvider("test-provider") })

	// Test registering duplicate provider
	err = RegisterMCPProvider("test-provider", testProviderFunc)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")

	// Verify provider can be retrieved
	provider, exists := GetMCPProvider("test-provider")
	assert.True(t, exists)
	assert.NotNil(t, provider)

	// Verify unknown provider returns false
	_, exists = GetMCPProvider("unknown-provider")
	assert.False(t, exists)
}

func TestYamlMCPProviderInvalidFile(t *testing.T) {
	// Test with non-existent file
	source := basecatalog.MCPSource{
		Type: "yaml",
		Properties: map[string]any{
			yamlMCPCatalogPathKey: "/nonexistent/path/catalog.yaml",
		},
	}

	provider, err := NewYamlMCPProvider(source)
	require.NoError(t, err)

	ctx := context.Background()
	recordChan := provider.Servers(ctx)

	// Should emit exactly one error record so the caller can mark the source as errored.
	var records []MCPServerProviderRecord
	for record := range recordChan {
		records = append(records, record)
	}
	require.Len(t, records, 1, "expected one error record for unreadable file")
	assert.NotNil(t, records[0].Error, "error record should have non-nil Error")
	assert.Nil(t, records[0].Server, "error record should have nil Server")
}

// Helper function
func strPtr(s string) *string {
	return &s
}
