package service

import (
	"encoding/json"
	"testing"

	"github.com/kubeflow/hub/catalog/internal/converter"
	"github.com/kubeflow/hub/catalog/pkg/openapi"
	"github.com/kubeflow/hub/internal/apiutils"
	"github.com/kubeflow/hub/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPServerPrerequisites(t *testing.T) {
	sharedDB, cleanup := testutils.SetupPostgresWithMigrations(t, testDatastoreSpec())
	defer cleanup()

	typeID := getMCPServerTypeID(t, sharedDB)
	repo := NewMCPServerRepository(sharedDB, typeID)

	t.Run("CreateFullPrerequisites", func(t *testing.T) {
		// Test case: MCP server with ServiceAccount, Secrets, and ConfigMaps
		mcpServer := &openapi.MCPServer{
			Name:      "kubernetes-mcp",
			Version:   apiutils.Of("1.0.0"),
			SourceId:  apiutils.Of("k8s-source"),
			ToolCount: 5,
			RuntimeMetadata: &openapi.MCPRuntimeMetadata{
				DefaultPort: apiutils.Of(int32(8080)),
				McpPath:     apiutils.Of("/mcp"),
				Prerequisites: &openapi.MCPPrerequisites{
					ServiceAccount: &openapi.MCPServiceAccountRequirement{
						Required:      apiutils.Of(true),
						Hint:          apiutils.Of("Needs 'view' ClusterRole for read-only K8s access"),
						SuggestedName: apiutils.Of("mcp-viewer"),
					},
					Secrets: []openapi.MCPSecretRequirement{
						{
							Name:        "openai-credentials",
							Description: "kubectl create secret generic openai-credentials --from-literal=api-key=sk-...",
							Keys: []openapi.MCPSecretKey{
								{
									Key:         "api-key",
									Description: "OpenAI API key for authentication",
									EnvVarName:  apiutils.Of("OPENAI_API_KEY"),
									Required:    apiutils.Of(true),
								},
							},
							MountAsFile: apiutils.Of(false),
						},
					},
					ConfigMaps: []openapi.MCPConfigMapRequirement{
						{
							Name:        "mcp-server-config",
							Description: "Configuration files for MCP server",
							Keys: []openapi.MCPConfigMapKey{
								{
									Key:            "config.toml",
									Description:    "Main configuration file",
									DefaultContent: apiutils.Of("log_level = \"info\"\nport = 8080\n"),
									Required:       apiutils.Of(true),
								},
							},
							MountAsFile: apiutils.Of(true),
							MountPath:   apiutils.Of("/etc/mcp-config"),
						},
					},
				},
			},
		}

		// Convert to DB model
		dbServer := converter.ConvertOpenapiMCPServerToDb(mcpServer)

		// Save to database
		saved, err := repo.Save(dbServer)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Retrieve from database
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		// Convert back to OpenAPI
		result := converter.ConvertDbMCPServerToOpenapi(retrieved)
		require.NotNil(t, result)
		require.NotNil(t, result.RuntimeMetadata)
		require.NotNil(t, result.RuntimeMetadata.Prerequisites)

		// Verify mcpPath
		require.NotNil(t, result.RuntimeMetadata.McpPath)
		assert.Equal(t, "/mcp", *result.RuntimeMetadata.McpPath)

		// Verify ServiceAccount
		require.NotNil(t, result.RuntimeMetadata.Prerequisites.ServiceAccount)
		assert.True(t, *result.RuntimeMetadata.Prerequisites.ServiceAccount.Required)
		assert.Equal(t, "Needs 'view' ClusterRole for read-only K8s access",
			*result.RuntimeMetadata.Prerequisites.ServiceAccount.Hint)
		assert.Equal(t, "mcp-viewer",
			*result.RuntimeMetadata.Prerequisites.ServiceAccount.SuggestedName)

		// Verify Secrets
		require.Len(t, result.RuntimeMetadata.Prerequisites.Secrets, 1)
		secret := result.RuntimeMetadata.Prerequisites.Secrets[0]
		assert.Equal(t, "openai-credentials", secret.Name)
		require.Len(t, secret.Keys, 1)
		assert.Equal(t, "api-key", secret.Keys[0].Key)
		assert.Equal(t, "OPENAI_API_KEY", *secret.Keys[0].EnvVarName)

		// Verify ConfigMaps
		require.Len(t, result.RuntimeMetadata.Prerequisites.ConfigMaps, 1)
		configMap := result.RuntimeMetadata.Prerequisites.ConfigMaps[0]
		assert.Equal(t, "mcp-server-config", configMap.Name)
		assert.Equal(t, "/etc/mcp-config", *configMap.MountPath)
		require.Len(t, configMap.Keys, 1)
		assert.Contains(t, *configMap.Keys[0].DefaultContent, "log_level")
	})

	t.Run("CreatePartialPrerequisitesOnlySecrets", func(t *testing.T) {
		// Test case: Only Secrets, no ServiceAccount or ConfigMaps
		mcpServer := &openapi.MCPServer{
			Name:      "partial-prereqs",
			Version:   apiutils.Of("1.0.0"),
			SourceId:  apiutils.Of("partial-source"),
			ToolCount: 1,
			RuntimeMetadata: &openapi.MCPRuntimeMetadata{
				DefaultPort: apiutils.Of(int32(3000)),
				Prerequisites: &openapi.MCPPrerequisites{
					Secrets: []openapi.MCPSecretRequirement{
						{
							Name:        "api-token",
							Description: "API token for external service",
							Keys: []openapi.MCPSecretKey{
								{
									Key:         "token",
									Description: "Bearer token",
									Required:    apiutils.Of(true),
								},
							},
						},
					},
				},
			},
		}

		dbServer := converter.ConvertOpenapiMCPServerToDb(mcpServer)
		saved, err := repo.Save(dbServer)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		result := converter.ConvertDbMCPServerToOpenapi(retrieved)
		require.NotNil(t, result.RuntimeMetadata)
		require.NotNil(t, result.RuntimeMetadata.Prerequisites)

		// Verify only Secrets present
		assert.Nil(t, result.RuntimeMetadata.Prerequisites.ServiceAccount)
		assert.Len(t, result.RuntimeMetadata.Prerequisites.Secrets, 1)
		assert.Len(t, result.RuntimeMetadata.Prerequisites.ConfigMaps, 0)
	})

	t.Run("BackwardCompatibilityNoPrerequisites", func(t *testing.T) {
		// Test case: Old-style MCP server without prerequisites field
		mcpServer := &openapi.MCPServer{
			Name:      "legacy-server",
			Version:   apiutils.Of("1.0.0"),
			SourceId:  apiutils.Of("legacy-source"),
			ToolCount: 2,
			RuntimeMetadata: &openapi.MCPRuntimeMetadata{
				DefaultPort: apiutils.Of(int32(8080)),
				// No prerequisites field
			},
		}

		dbServer := converter.ConvertOpenapiMCPServerToDb(mcpServer)
		saved, err := repo.Save(dbServer)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		result := converter.ConvertDbMCPServerToOpenapi(retrieved)
		require.NotNil(t, result)
		// Prerequisites should be nil for backward compatibility
		if result.RuntimeMetadata != nil {
			assert.Nil(t, result.RuntimeMetadata.Prerequisites)
		}
	})

	t.Run("McpPathCustomValue", func(t *testing.T) {
		// Test case: mcpPath with custom value
		mcpServer := &openapi.MCPServer{
			Name:      "custom-path-server",
			Version:   apiutils.Of("1.0.0"),
			SourceId:  apiutils.Of("custom-source"),
			ToolCount: 3,
			RuntimeMetadata: &openapi.MCPRuntimeMetadata{
				DefaultPort: apiutils.Of(int32(8080)),
				McpPath:     apiutils.Of("/api/v1/mcp"),
			},
		}

		dbServer := converter.ConvertOpenapiMCPServerToDb(mcpServer)
		saved, err := repo.Save(dbServer)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		result := converter.ConvertDbMCPServerToOpenapi(retrieved)

		// Verify custom mcpPath is preserved
		require.NotNil(t, result.RuntimeMetadata)
		require.NotNil(t, result.RuntimeMetadata.McpPath)
		assert.Equal(t, "/api/v1/mcp", *result.RuntimeMetadata.McpPath)
	})

	t.Run("JSONRoundTripComplexPrerequisites", func(t *testing.T) {
		// Test case: Verify JSON serialization preserves all nested fields
		prerequisites := &openapi.MCPPrerequisites{
			ServiceAccount: &openapi.MCPServiceAccountRequirement{
				Required:      apiutils.Of(true),
				Hint:          apiutils.Of("Needs cluster admin"),
				SuggestedName: apiutils.Of("mcp-admin"),
			},
			Secrets: []openapi.MCPSecretRequirement{
				{
					Name:        "secret1",
					Description: "First secret",
					Keys: []openapi.MCPSecretKey{
						{Key: "key1", Description: "Key 1"},
						{Key: "key2", Description: "Key 2"},
					},
				},
			},
			ConfigMaps: []openapi.MCPConfigMapRequirement{
				{
					Name:        "config1",
					Description: "First config",
					Keys: []openapi.MCPConfigMapKey{
						{
							Key:            "config.yaml",
							Description:    "Main config",
							DefaultContent: apiutils.Of("key: value\n"),
						},
					},
					MountAsFile: apiutils.Of(true),
					MountPath:   apiutils.Of("/config"),
				},
			},
			EnvironmentVariables: []openapi.MCPEnvVarMetadata{
				{
					Name:        "ENV_VAR_1",
					Description: "Environment variable 1",
					Required:    apiutils.Of(true),
				},
			},
			CustomResources: []string{"PersistentVolumeClaim: data-pvc"},
		}

		// Marshal to JSON
		jsonBytes, err := json.Marshal(prerequisites)
		require.NoError(t, err)

		// Unmarshal back
		var unmarshaled openapi.MCPPrerequisites
		err = json.Unmarshal(jsonBytes, &unmarshaled)
		require.NoError(t, err)

		// Verify all fields preserved
		assert.Equal(t, *prerequisites.ServiceAccount.Required,
			*unmarshaled.ServiceAccount.Required)
		assert.Equal(t, len(prerequisites.Secrets), len(unmarshaled.Secrets))
		assert.Equal(t, len(prerequisites.ConfigMaps), len(unmarshaled.ConfigMaps))
		assert.Equal(t, len(prerequisites.EnvironmentVariables),
			len(unmarshaled.EnvironmentVariables))
		assert.Equal(t, prerequisites.CustomResources, unmarshaled.CustomResources)
	})

	t.Run("PrerequisitesWithEnvironmentVariables", func(t *testing.T) {
		// Test case: Prerequisites with environment variables persisted through database
		mcpServer := &openapi.MCPServer{
			Name:      "env-vars-server",
			Version:   apiutils.Of("1.0.0"),
			SourceId:  apiutils.Of("env-source"),
			ToolCount: 1,
			RuntimeMetadata: &openapi.MCPRuntimeMetadata{
				DefaultPort: apiutils.Of(int32(8080)),
				Prerequisites: &openapi.MCPPrerequisites{
					EnvironmentVariables: []openapi.MCPEnvVarMetadata{
						{
							Name:        "API_ENDPOINT",
							Description: "Service API endpoint URL",
							Required:    apiutils.Of(true),
							Type:        apiutils.Of("string"),
						},
						{
							Name:         "LOG_LEVEL",
							Description:  "Logging verbosity level",
							Required:     apiutils.Of(false),
							DefaultValue: apiutils.Of("info"),
							Type:         apiutils.Of("string"),
						},
					},
				},
			},
		}

		dbServer := converter.ConvertOpenapiMCPServerToDb(mcpServer)
		saved, err := repo.Save(dbServer)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		result := converter.ConvertDbMCPServerToOpenapi(retrieved)
		require.NotNil(t, result.RuntimeMetadata)
		require.NotNil(t, result.RuntimeMetadata.Prerequisites)

		// Verify environment variables are preserved
		require.Len(t, result.RuntimeMetadata.Prerequisites.EnvironmentVariables, 2)
		assert.Equal(t, "API_ENDPOINT", result.RuntimeMetadata.Prerequisites.EnvironmentVariables[0].Name)
		assert.True(t, *result.RuntimeMetadata.Prerequisites.EnvironmentVariables[0].Required)
		assert.Equal(t, "LOG_LEVEL", result.RuntimeMetadata.Prerequisites.EnvironmentVariables[1].Name)
		assert.Equal(t, "info", *result.RuntimeMetadata.Prerequisites.EnvironmentVariables[1].DefaultValue)
	})

	t.Run("PrerequisitesEmptyArrays", func(t *testing.T) {
		// Test case: Prerequisites with empty arrays
		mcpServer := &openapi.MCPServer{
			Name:      "empty-arrays-server",
			Version:   apiutils.Of("1.0.0"),
			SourceId:  apiutils.Of("empty-source"),
			ToolCount: 1,
			RuntimeMetadata: &openapi.MCPRuntimeMetadata{
				DefaultPort: apiutils.Of(int32(8080)),
				Prerequisites: &openapi.MCPPrerequisites{
					Secrets:              []openapi.MCPSecretRequirement{},
					ConfigMaps:           []openapi.MCPConfigMapRequirement{},
					EnvironmentVariables: []openapi.MCPEnvVarMetadata{},
					CustomResources:      []string{},
				},
			},
		}

		dbServer := converter.ConvertOpenapiMCPServerToDb(mcpServer)
		saved, err := repo.Save(dbServer)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		result := converter.ConvertDbMCPServerToOpenapi(retrieved)
		require.NotNil(t, result.RuntimeMetadata)
		require.NotNil(t, result.RuntimeMetadata.Prerequisites)

		// Verify empty arrays are preserved
		assert.Len(t, result.RuntimeMetadata.Prerequisites.Secrets, 0)
		assert.Len(t, result.RuntimeMetadata.Prerequisites.ConfigMaps, 0)
		assert.Len(t, result.RuntimeMetadata.Prerequisites.EnvironmentVariables, 0)
		assert.Len(t, result.RuntimeMetadata.Prerequisites.CustomResources, 0)
	})

	t.Run("PrerequisitesSpecialCharactersInDescriptions", func(t *testing.T) {
		// Test case: Special characters and unicode in hints/descriptions
		mcpServer := &openapi.MCPServer{
			Name:      "special-chars-server",
			Version:   apiutils.Of("1.0.0"),
			SourceId:  apiutils.Of("special-source"),
			ToolCount: 1,
			RuntimeMetadata: &openapi.MCPRuntimeMetadata{
				DefaultPort: apiutils.Of(int32(8080)),
				Prerequisites: &openapi.MCPPrerequisites{
					ServiceAccount: &openapi.MCPServiceAccountRequirement{
						Required:      apiutils.Of(true),
						Hint:          apiutils.Of("Requires RBAC: pods/list, pods/get, pods/watch → 'view' ClusterRole ✓"),
						SuggestedName: apiutils.Of("mcp-server-sa"),
					},
					Secrets: []openapi.MCPSecretRequirement{
						{
							Name:        "special-secret",
							Description: "kubectl create secret generic special-secret --from-literal=key='value with spaces & \"quotes\"'",
							Keys: []openapi.MCPSecretKey{
								{
									Key:         "special-key",
									Description: "API key with special chars: !@#$%^&*()_+-=[]{}|;':\"<>?,./",
									Required:    apiutils.Of(true),
								},
							},
						},
					},
				},
			},
		}

		dbServer := converter.ConvertOpenapiMCPServerToDb(mcpServer)
		saved, err := repo.Save(dbServer)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		result := converter.ConvertDbMCPServerToOpenapi(retrieved)
		require.NotNil(t, result.RuntimeMetadata)
		require.NotNil(t, result.RuntimeMetadata.Prerequisites)

		// Verify special characters are preserved
		assert.Contains(t, *result.RuntimeMetadata.Prerequisites.ServiceAccount.Hint, "→")
		assert.Contains(t, *result.RuntimeMetadata.Prerequisites.ServiceAccount.Hint, "✓")
		assert.Contains(t, result.RuntimeMetadata.Prerequisites.Secrets[0].Description, "\"quotes\"")
		assert.Contains(t, result.RuntimeMetadata.Prerequisites.Secrets[0].Keys[0].Description, "!@#$%")
	})

	t.Run("PrerequisitesCustomResources", func(t *testing.T) {
		// Test case: Custom resources hints
		mcpServer := &openapi.MCPServer{
			Name:      "custom-resources-server",
			Version:   apiutils.Of("1.0.0"),
			SourceId:  apiutils.Of("custom-source"),
			ToolCount: 1,
			RuntimeMetadata: &openapi.MCPRuntimeMetadata{
				DefaultPort: apiutils.Of(int32(8080)),
				Prerequisites: &openapi.MCPPrerequisites{
					CustomResources: []string{
						"PersistentVolumeClaim: mcp-data-storage (100Gi)",
						"NetworkPolicy: allow-mcp-traffic",
						"ServiceMonitor: mcp-metrics",
					},
				},
			},
		}

		dbServer := converter.ConvertOpenapiMCPServerToDb(mcpServer)
		saved, err := repo.Save(dbServer)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		result := converter.ConvertDbMCPServerToOpenapi(retrieved)
		require.NotNil(t, result.RuntimeMetadata)
		require.NotNil(t, result.RuntimeMetadata.Prerequisites)

		// Verify custom resources are preserved
		require.Len(t, result.RuntimeMetadata.Prerequisites.CustomResources, 3)
		assert.Contains(t, result.RuntimeMetadata.Prerequisites.CustomResources[0], "PersistentVolumeClaim")
		assert.Contains(t, result.RuntimeMetadata.Prerequisites.CustomResources[1], "NetworkPolicy")
		assert.Contains(t, result.RuntimeMetadata.Prerequisites.CustomResources[2], "ServiceMonitor")
	})
}
