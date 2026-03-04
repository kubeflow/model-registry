package mcpcatalog

import (
	"context"
	"encoding/json"
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

func TestYamlMCPEndpointsValidation(t *testing.T) {
	tests := []struct {
		name        string
		server      *yamlMCPServer
		wantErr     bool
		errContains string
	}{
		{
			name: "remote server with valid HTTP endpoint passes",
			server: &yamlMCPServer{
				Name:           "remote-http",
				DeploymentMode: strPtr("remote"),
				Endpoints:      &yamlMCPEndpoints{HTTP: strPtr("https://api.example.com/mcp")},
			},
		},
		{
			name: "remote server with valid SSE endpoint passes",
			server: &yamlMCPServer{
				Name:           "remote-sse",
				DeploymentMode: strPtr("remote"),
				Endpoints:      &yamlMCPEndpoints{SSE: strPtr("https://api.example.com/mcp/events")},
			},
		},
		{
			name: "remote server with both endpoints passes",
			server: &yamlMCPServer{
				Name:           "remote-both",
				DeploymentMode: strPtr("remote"),
				Endpoints: &yamlMCPEndpoints{
					HTTP: strPtr("https://api.example.com/mcp"),
					SSE:  strPtr("https://api.example.com/mcp/events"),
				},
			},
		},
		{
			name: "remote server with no endpoints is rejected",
			server: &yamlMCPServer{
				Name:           "remote-no-endpoints",
				DeploymentMode: strPtr("remote"),
			},
			wantErr:     true,
			errContains: "no endpoints defined",
		},
		{
			name: "remote server with empty endpoints struct is rejected",
			server: &yamlMCPServer{
				Name:           "remote-empty-endpoints",
				DeploymentMode: strPtr("remote"),
				Endpoints:      &yamlMCPEndpoints{},
			},
			wantErr:     true,
			errContains: "no endpoints defined",
		},
		{
			name: "local server without endpoints passes",
			server: &yamlMCPServer{
				Name:           "local-no-endpoints",
				DeploymentMode: strPtr("local"),
			},
		},
		{
			name: "local server with endpoints passes (hybrid)",
			server: &yamlMCPServer{
				Name:           "local-with-endpoints",
				DeploymentMode: strPtr("local"),
				Endpoints:      &yamlMCPEndpoints{HTTP: strPtr("https://api.example.com/mcp")},
			},
		},
		{
			name: "non-HTTP URL in HTTP endpoint is rejected",
			server: &yamlMCPServer{
				Name:           "bad-url-scheme",
				DeploymentMode: strPtr("remote"),
				Endpoints:      &yamlMCPEndpoints{HTTP: strPtr("ftp://example.com/mcp")},
			},
			wantErr:     true,
			errContains: "must be a valid URL",
		},
		{
			name: "malformed URL is rejected",
			server: &yamlMCPServer{
				Name:           "malformed-url",
				DeploymentMode: strPtr("remote"),
				Endpoints:      &yamlMCPEndpoints{HTTP: strPtr("not-a-url-at-all")},
			},
			wantErr:     true,
			errContains: "must be a valid URL",
		},
		{
			name: "plain http URL is accepted",
			server: &yamlMCPServer{
				Name:           "plain-http",
				DeploymentMode: strPtr("remote"),
				Endpoints:      &yamlMCPEndpoints{HTTP: strPtr("http://api.example.com/mcp")},
			},
		},
		{
			name: "websocket-only remote server passes",
			server: &yamlMCPServer{
				Name:           "remote-ws",
				DeploymentMode: strPtr("remote"),
				Endpoints:      &yamlMCPEndpoints{WebSocket: strPtr("wss://api.example.com/ws")},
			},
		},
		{
			name: "ws scheme accepted for websocket endpoint",
			server: &yamlMCPServer{
				Name:           "remote-ws-plain",
				DeploymentMode: strPtr("remote"),
				Endpoints:      &yamlMCPEndpoints{WebSocket: strPtr("ws://api.example.com/ws")},
			},
		},
		{
			name: "ftp scheme rejected for websocket endpoint",
			server: &yamlMCPServer{
				Name:           "bad-ws-scheme",
				DeploymentMode: strPtr("remote"),
				Endpoints:      &yamlMCPEndpoints{WebSocket: strPtr("ftp://example.com/ws")},
			},
			wantErr:     true,
			errContains: "must be a valid URL",
		},
		{
			name: "server with no deploymentMode and endpoints passes",
			server: &yamlMCPServer{
				Name:      "no-mode-with-endpoints",
				Endpoints: &yamlMCPEndpoints{HTTP: strPtr("https://api.example.com/mcp")},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			record := tc.server.ToMCPServerProviderRecord()
			if tc.wantErr {
				assert.NotNil(t, record.Error, "expected an error but got none")
				if tc.errContains != "" {
					assert.Contains(t, record.Error.Error(), tc.errContains)
				}
			} else {
				assert.Nil(t, record.Error, "expected no error but got: %v", record.Error)
				assert.NotNil(t, record.Server)
			}
		})
	}
}

func TestYamlMCPEndpointsJSONKeys(t *testing.T) {
	endpoints := &yamlMCPEndpoints{
		HTTP: strPtr("https://api.example.com/mcp"),
		SSE:  strPtr("https://api.example.com/mcp/events"),
	}

	jsonBytes, err := json.Marshal(endpoints)
	require.NoError(t, err)

	jsonStr := string(jsonBytes)
	assert.Contains(t, jsonStr, `"http":`, "JSON key should be lowercase 'http'")
	assert.Contains(t, jsonStr, `"sse":`, "JSON key should be lowercase 'sse'")
	assert.NotContains(t, jsonStr, `"HTTP":`, "JSON key must not be uppercase 'HTTP'")
	assert.NotContains(t, jsonStr, `"SSE":`, "JSON key must not be uppercase 'SSE'")
}

// Helper function
func strPtr(s string) *string {
	return &s
}
