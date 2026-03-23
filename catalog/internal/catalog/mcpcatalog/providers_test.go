package mcpcatalog

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	apimodels "github.com/kubeflow/model-registry/catalog/pkg/openapi"
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

func TestYamlMCPServerLicenseTransformation(t *testing.T) {
	tests := []struct {
		name            string
		inputLicense    string
		expectedLicense string
	}{
		{
			name:            "apache-2.0 SPDX transforms to human-readable",
			inputLicense:    "apache-2.0",
			expectedLicense: "Apache 2.0",
		},
		{
			name:            "MIT SPDX transforms to human-readable",
			inputLicense:    "mit",
			expectedLicense: "MIT",
		},
		{
			name:            "GPL-3.0 SPDX transforms to human-readable",
			inputLicense:    "gpl-3.0",
			expectedLicense: "GPL 3.0",
		},
		{
			name:            "BSD-3-Clause SPDX transforms to human-readable",
			inputLicense:    "bsd-3-clause",
			expectedLicense: "BSD 3-Clause",
		},
		{
			name:            "Llama3.1 custom license transforms",
			inputLicense:    "llama3.1",
			expectedLicense: "Llama 3.1 Community License",
		},
		{
			name:            "Llama3 custom license transforms",
			inputLicense:    "llama3",
			expectedLicense: "Llama 3 Community License",
		},
		{
			name:            "Gemma custom license transforms",
			inputLicense:    "gemma",
			expectedLicense: "Gemma License",
		},
		{
			name:            "Creative Commons CC-BY-4.0 transforms",
			inputLicense:    "cc-by-4.0",
			expectedLicense: "Creative Commons Attribution 4.0 International",
		},
		{
			name:            "Creative Commons CC0 transforms",
			inputLicense:    "cc0-1.0",
			expectedLicense: "Creative Commons Zero v1.0 Universal",
		},
		{
			name:            "OpenRAIL license transforms",
			inputLicense:    "openrail",
			expectedLicense: "OpenRAIL License",
		},
		{
			name:            "BigScience OpenRAIL-M transforms",
			inputLicense:    "bigscience-openrail-m",
			expectedLicense: "BigScience OpenRAIL-M License",
		},
		{
			name:            "Uppercase APACHE-2.0 normalizes and transforms",
			inputLicense:    "APACHE-2.0",
			expectedLicense: "Apache 2.0",
		},
		{
			name:            "Mixed case Mit normalizes and transforms",
			inputLicense:    "Mit",
			expectedLicense: "MIT",
		},
		{
			name:            "GPL-2.0 SPDX transforms to human-readable",
			inputLicense:    "gpl-2.0",
			expectedLicense: "GPL 2.0",
		},
		{
			name:            "LGPL-3.0 SPDX transforms to human-readable",
			inputLicense:    "lgpl-3.0",
			expectedLicense: "LGPL 3.0",
		},
		{
			name:            "Unknown license passes through unchanged",
			inputLicense:    "custom-proprietary-license-1.0",
			expectedLicense: "custom-proprietary-license-1.0",
		},
		{
			name:            "Empty license results in no license property",
			inputLicense:    "",
			expectedLicense: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var yamlServer *yamlMCPServer
			if tt.inputLicense == "" {
				yamlServer = &yamlMCPServer{
					Name:     "test-server",
					License:  nil, // No license field
					Provider: strPtr("Test Provider"),
				}
			} else {
				yamlServer = &yamlMCPServer{
					Name:     "test-server",
					License:  &tt.inputLicense,
					Provider: strPtr("Test Provider"),
				}
			}

			record := yamlServer.ToMCPServerProviderRecord()

			assert.Nil(t, record.Error, "expected no error converting server")
			assert.NotNil(t, record.Server, "expected server to be created")

			props := record.Server.GetProperties()
			require.NotNil(t, props, "expected properties to be set")

			// Find the license property
			var foundLicense *string
			for _, prop := range *props {
				if prop.Name == "license" && prop.StringValue != nil {
					foundLicense = prop.StringValue
					break
				}
			}

			if tt.expectedLicense == "" {
				assert.Nil(t, foundLicense, "expected no license property for empty input")
			} else {
				require.NotNil(t, foundLicense, "expected license property to be set")
				assert.Equal(t, tt.expectedLicense, *foundLicense,
					"license should be transformed from %q to %q", tt.inputLicense, tt.expectedLicense)
			}
		})
	}
}

func TestYamlMCPServerLicenseConsistencyWithModels(t *testing.T) {
	// This test verifies that MCP servers apply the same license transformation
	// as models to ensure API consistency
	testLicenses := []string{
		"apache-2.0",
		"mit",
		"llama3.1",
		"cc-by-4.0",
		"gpl-3.0",
	}

	for _, license := range testLicenses {
		t.Run(license, func(t *testing.T) {
			// Get the expected human-readable transformation
			expectedHumanReadable := basecatalog.TransformLicenseToHumanReadable(license)

			// Create MCP server with this license
			yamlServer := &yamlMCPServer{
				Name:    "test-server",
				License: &license,
			}

			record := yamlServer.ToMCPServerProviderRecord()
			require.Nil(t, record.Error)
			require.NotNil(t, record.Server)

			// Extract license from properties
			props := record.Server.GetProperties()
			require.NotNil(t, props)

			var actualLicense *string
			for _, prop := range *props {
				if prop.Name == "license" && prop.StringValue != nil {
					actualLicense = prop.StringValue
					break
				}
			}

			require.NotNil(t, actualLicense, "license property should be set")
			assert.Equal(t, expectedHumanReadable, *actualLicense,
				"MCP server license transformation should match model catalog transformation")
		})
	}
}

func TestYamlMCPToolAccessTypePreserved(t *testing.T) {
	accessType := "read_write"
	paramDesc := "Project key"
	yamlServer := &yamlMCPServer{
		Name: "test-server",
		Tools: []*yamlMCPTool{
			{
				Name:       "create_issue",
				AccessType: &accessType,
				Parameters: []yamlMCPParameter{
					{Name: "project", Type: "string", Description: &paramDesc, Required: true},
					{Name: "summary", Type: "string", Required: true},
				},
			},
			{
				Name: "search_issues",
				// no AccessType set - should be nil in the record
			},
		},
	}

	record := yamlServer.ToMCPServerProviderRecord()

	require.Len(t, record.Tools, 2)

	createTool := record.Tools[0]
	require.NotNil(t, createTool.AccessType, "accessType should be preserved")
	assert.Equal(t, "read_write", *createTool.AccessType)
	require.Len(t, createTool.Parameters, 2)
	assert.Equal(t, "project", createTool.Parameters[0].Name)
	assert.Equal(t, "string", createTool.Parameters[0].Type)
	assert.True(t, createTool.Parameters[0].Required)
	require.NotNil(t, createTool.Parameters[0].Description)
	assert.Equal(t, "Project key", *createTool.Parameters[0].Description)
	assert.Equal(t, "summary", createTool.Parameters[1].Name)
	assert.True(t, createTool.Parameters[1].Required)

	searchTool := record.Tools[1]
	assert.Nil(t, searchTool.AccessType, "nil accessType should stay nil")
	assert.Empty(t, searchTool.Parameters)
}

func TestYamlMCPServerRuntimeMetadataConversion(t *testing.T) {
	defaultPort := int32(8080)
	mcpPath := "/mcp"
	saRequired := true
	saHint := "Needs 'view' ClusterRole"
	saName := "mcp-viewer"
	cmMountAsFile := true
	cmMountPath := "/etc/mcp-config"
	cmKeyRequired := false

	yamlServer := &yamlMCPServer{
		Name: "runtime-meta-server",
		RuntimeMetadata: &apimodels.MCPRuntimeMetadata{
			DefaultPort: &defaultPort,
			McpPath:     &mcpPath,
			Prerequisites: &apimodels.MCPPrerequisites{
				ServiceAccount: &apimodels.MCPServiceAccountRequirement{
					Required:      &saRequired,
					Hint:          &saHint,
					SuggestedName: &saName,
				},
				ConfigMaps: []apimodels.MCPConfigMapRequirement{
					{
						Name:        "server-config",
						Description: "Config files",
						MountAsFile: &cmMountAsFile,
						MountPath:   &cmMountPath,
						Keys: []apimodels.MCPConfigMapKey{
							{
								Key:         "config.toml",
								Description: "Main config",
								Required:    &cmKeyRequired,
							},
						},
					},
				},
			},
		},
	}

	record := yamlServer.ToMCPServerProviderRecord()
	assert.Nil(t, record.Error)
	assert.NotNil(t, record.Server)

	props := record.Server.GetProperties()
	require.NotNil(t, props)

	var runtimeJSON *string
	for _, prop := range *props {
		if prop.Name == "runtimeMetadata" && prop.StringValue != nil {
			runtimeJSON = prop.StringValue
			break
		}
	}
	require.NotNil(t, runtimeJSON, "runtimeMetadata property should be set")

	var parsed apimodels.MCPRuntimeMetadata
	err := json.Unmarshal([]byte(*runtimeJSON), &parsed)
	require.NoError(t, err)

	assert.Equal(t, int32(8080), *parsed.DefaultPort)
	assert.Equal(t, "/mcp", *parsed.McpPath)
	require.NotNil(t, parsed.Prerequisites)
	require.NotNil(t, parsed.Prerequisites.ServiceAccount)
	assert.True(t, *parsed.Prerequisites.ServiceAccount.Required)
	assert.Equal(t, "mcp-viewer", *parsed.Prerequisites.ServiceAccount.SuggestedName)
	require.Len(t, parsed.Prerequisites.ConfigMaps, 1)
	assert.Equal(t, "server-config", parsed.Prerequisites.ConfigMaps[0].Name)
}

func TestYamlMCPServerRuntimeMetadataNil(t *testing.T) {
	yamlServer := &yamlMCPServer{
		Name: "no-runtime-meta",
	}

	record := yamlServer.ToMCPServerProviderRecord()
	assert.Nil(t, record.Error)

	props := record.Server.GetProperties()
	require.NotNil(t, props)

	for _, prop := range *props {
		assert.NotEqual(t, "runtimeMetadata", prop.Name,
			"runtimeMetadata property should not be set when RuntimeMetadata is nil")
	}
}

func TestYamlMCPProviderEmitWithRuntimeMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-catalog.yaml")

	yamlContent := `mcp_servers:
  - name: "k8s-mcp"
    provider: "Test"
    version: "1.0.0"
    runtimeMetadata:
      defaultPort: 8080
      mcpPath: /mcp
      prerequisites:
        serviceAccount:
          required: true
          hint: "Needs view ClusterRole"
          suggestedName: mcp-viewer
        secrets:
          - name: api-creds
            description: "API credentials"
            keys:
              - key: api-key
                description: "API key"
                envVarName: API_KEY
                required: true
            mountAsFile: false
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

	var record MCPServerProviderRecord
	for r := range recordChan {
		record = r
	}

	assert.Nil(t, record.Error)
	assert.NotNil(t, record.Server)

	props := record.Server.GetProperties()
	require.NotNil(t, props)

	var runtimeJSON *string
	for _, prop := range *props {
		if prop.Name == "runtimeMetadata" && prop.StringValue != nil {
			runtimeJSON = prop.StringValue
			break
		}
	}
	require.NotNil(t, runtimeJSON, "runtimeMetadata should be parsed from YAML file")

	var parsed apimodels.MCPRuntimeMetadata
	err = json.Unmarshal([]byte(*runtimeJSON), &parsed)
	require.NoError(t, err)

	assert.Equal(t, int32(8080), *parsed.DefaultPort)
	assert.Equal(t, "/mcp", *parsed.McpPath)
	require.NotNil(t, parsed.Prerequisites)
	require.NotNil(t, parsed.Prerequisites.ServiceAccount)
	assert.True(t, *parsed.Prerequisites.ServiceAccount.Required)
	assert.Equal(t, "mcp-viewer", *parsed.Prerequisites.ServiceAccount.SuggestedName)
	require.Len(t, parsed.Prerequisites.Secrets, 1)
	assert.Equal(t, "api-creds", parsed.Prerequisites.Secrets[0].Name)
	require.Len(t, parsed.Prerequisites.Secrets[0].Keys, 1)
	assert.Equal(t, "API_KEY", *parsed.Prerequisites.Secrets[0].Keys[0].EnvVarName)
}

// Helper function
func strPtr(s string) *string {
	return &s
}
