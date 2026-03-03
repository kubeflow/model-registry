package mcpcatalog

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog/modelcatalog"
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	os.Exit(testutils.TestMainPostgresHelper(m))
}

func getMCPServerTypeIDForTest(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.MCPServerTypeName).First(&typeRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create the type if it doesn't exist
			typeRecord = schema.Type{
				Name: service.MCPServerTypeName,
			}
			err = db.Create(&typeRecord).Error
			require.NoError(t, err)
		} else {
			require.NoError(t, err)
		}
	}
	return typeRecord.ID
}

func getMCPServerToolTypeIDForTest(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.MCPServerToolTypeName).First(&typeRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create the type if it doesn't exist
			typeRecord = schema.Type{
				Name: service.MCPServerToolTypeName,
			}
			err = db.Create(&typeRecord).Error
			require.NoError(t, err)
		} else {
			require.NoError(t, err)
		}
	}
	return typeRecord.ID
}

func setupMCPLoaderTest(t *testing.T) (*gorm.DB, service.Services, func()) {
	sharedDB, cleanup := testutils.SetupPostgresWithMigrations(t, service.DatastoreSpec())

	// Get type IDs
	catalogModelTypeID := modelcatalog.GetCatalogModelTypeIDForDBTest(t, sharedDB)
	modelArtifactTypeID := modelcatalog.GetCatalogModelArtifactTypeIDForDBTest(t, sharedDB)
	metricsArtifactTypeID := modelcatalog.GetCatalogMetricsArtifactTypeIDForDBTest(t, sharedDB)
	catalogSourceTypeID := modelcatalog.GetCatalogSourceTypeIDForDBTest(t, sharedDB)
	mcpServerTypeID := getMCPServerTypeIDForTest(t, sharedDB)
	mcpServerToolTypeID := getMCPServerToolTypeIDForTest(t, sharedDB)

	// Create repositories
	catalogModelRepo := service.NewCatalogModelRepository(sharedDB, catalogModelTypeID)
	catalogArtifactRepo := service.NewCatalogArtifactRepository(sharedDB, map[string]int32{
		service.CatalogModelArtifactTypeName:   modelArtifactTypeID,
		service.CatalogMetricsArtifactTypeName: metricsArtifactTypeID,
	})
	modelArtifactRepo := service.NewCatalogModelArtifactRepository(sharedDB, modelArtifactTypeID)
	metricsArtifactRepo := service.NewCatalogMetricsArtifactRepository(sharedDB, metricsArtifactTypeID)
	catalogSourceRepo := service.NewCatalogSourceRepository(sharedDB, catalogSourceTypeID)
	mcpServerRepo := service.NewMCPServerRepository(sharedDB, mcpServerTypeID)
	mcpServerToolRepo := service.NewMCPServerToolRepository(sharedDB, mcpServerToolTypeID)

	services := service.NewServices(
		catalogModelRepo,
		catalogArtifactRepo,
		modelArtifactRepo,
		metricsArtifactRepo,
		catalogSourceRepo,
		service.NewPropertyOptionsRepository(sharedDB),
		mcpServerRepo,
		mcpServerToolRepo,
	)

	return sharedDB, services, cleanup
}

func TestMCPLoaderBasicLoad(t *testing.T) {
	_, services, cleanup := setupMCPLoaderTest(t)
	defer cleanup()

	// Create temporary MCP catalog files
	tmpDir := t.TempDir()

	serversFile := filepath.Join(tmpDir, "servers.yaml")
	err := os.WriteFile(serversFile, []byte(`mcp_servers:
  - name: "test-server"
    description: "Test MCP server"
    provider: "Test Provider"
    version: "1.0.0"
    tools:
      - name: "test-tool"
        description: "Test tool"
`), 0644)
	require.NoError(t, err)

	sourcesFile := filepath.Join(tmpDir, "sources.yaml")
	err = os.WriteFile(sourcesFile, []byte(`mcp_catalogs:
  - name: "Test MCP Catalog"
    id: test_mcp_catalog
    type: yaml
    enabled: true
    properties:
      yamlCatalogPath: `+serversFile+`
`), 0644)
	require.NoError(t, err)

	// Create and start the loader
	baseLoader := basecatalog.NewBaseLoader([]string{sourcesFile})
	loader := NewMCPLoaderWithState(services, baseLoader)
	ctx := context.Background()

	// Parse configs
	err = loader.ParseAllConfigs()
	require.NoError(t, err)

	// Set leader mode and perform leader operations
	baseLoader.SetLeader(true)

	// Perform leader operations in background
	leaderDone := make(chan error, 1)
	go func() {
		leaderDone <- loader.PerformLeaderOperations(ctx, mapset.NewSet[string]())
	}()

	// Wait for leader operations to complete
	select {
	case err := <-leaderDone:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for leader operations")
	}

	// Wait for all database writes to complete
	baseLoader.WaitForInflightWrites(5 * time.Second)

	// Verify server was loaded
	server, err := services.MCPServerRepository.GetByNameAndVersion("test-server", "1.0.0")
	require.NoError(t, err)
	assert.NotNil(t, server)

	attrs := server.GetAttributes()
	require.NotNil(t, attrs)
	assert.Equal(t, "test-server", *attrs.Name)

	// Verify properties
	props := server.GetProperties()
	require.NotNil(t, props)

	foundDescription := false
	foundSourceID := false
	for _, prop := range *props {
		if prop.Name == "description" && prop.StringValue != nil {
			assert.Equal(t, "Test MCP server", *prop.StringValue)
			foundDescription = true
		}
		if prop.Name == "source_id" && prop.StringValue != nil {
			assert.Equal(t, "test_mcp_catalog", *prop.StringValue)
			foundSourceID = true
		}
	}
	assert.True(t, foundDescription, "description property should be set")
	assert.True(t, foundSourceID, "source_id property should be set")

	// Verify tools were persisted to the database
	require.NotNil(t, server.GetID())
	tools, err := services.MCPServerToolRepository.List(models.MCPServerToolListOptions{ParentID: *server.GetID()})
	require.NoError(t, err)
	require.Len(t, tools, 1)
	toolAttrs := tools[0].GetAttributes()
	require.NotNil(t, toolAttrs)
	assert.Equal(t, "test-server@1.0.0:test-tool", *toolAttrs.Name)
}

func TestMCPLoaderDisabledSource(t *testing.T) {
	_, services, cleanup := setupMCPLoaderTest(t)
	defer cleanup()

	// Create temporary MCP catalog files
	tmpDir := t.TempDir()

	serversFile := filepath.Join(tmpDir, "servers.yaml")
	err := os.WriteFile(serversFile, []byte(`mcp_servers:
  - name: "disabled-server"
    description: "This server should not be loaded"
`), 0644)
	require.NoError(t, err)

	sourcesFile := filepath.Join(tmpDir, "sources.yaml")
	err = os.WriteFile(sourcesFile, []byte(`mcp_catalogs:
  - name: "Disabled MCP Catalog"
    id: disabled_mcp_catalog
    type: yaml
    enabled: false
    properties:
      yamlCatalogPath: `+serversFile+`
`), 0644)
	require.NoError(t, err)

	// Create and start the loader
	baseLoader := basecatalog.NewBaseLoader([]string{sourcesFile})
	loader := NewMCPLoaderWithState(services, baseLoader)
	ctx := context.Background()

	// Parse configs
	err = loader.ParseAllConfigs()
	require.NoError(t, err)

	// Set leader mode and perform leader operations
	baseLoader.SetLeader(true)

	// Perform leader operations in background
	leaderDone := make(chan error, 1)
	go func() {
		leaderDone <- loader.PerformLeaderOperations(ctx, mapset.NewSet[string]())
	}()

	// Wait for leader operations to complete
	select {
	case err := <-leaderDone:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for leader operations")
	}

	// Wait for all database writes to complete
	baseLoader.WaitForInflightWrites(5 * time.Second)

	// Verify server was NOT loaded
	_, err = services.MCPServerRepository.GetByNameAndVersion("disabled-server", "")
	assert.Error(t, err)

}

func TestMCPLoaderEventHandlers(t *testing.T) {
	_, services, cleanup := setupMCPLoaderTest(t)
	defer cleanup()

	// Create temporary MCP catalog files
	tmpDir := t.TempDir()

	serversFile := filepath.Join(tmpDir, "servers.yaml")
	err := os.WriteFile(serversFile, []byte(`mcp_servers:
  - name: "event-test-server"
    description: "Server for testing event handlers"
    tools:
      - name: "handler-tool"
`), 0644)
	require.NoError(t, err)

	sourcesFile := filepath.Join(tmpDir, "sources.yaml")
	err = os.WriteFile(sourcesFile, []byte(`mcp_catalogs:
  - name: "Event Test Catalog"
    id: event_test_catalog
    type: yaml
    enabled: true
    properties:
      yamlCatalogPath: `+serversFile+`
`), 0644)
	require.NoError(t, err)

	// Create loader and register event handler
	baseLoader := basecatalog.NewBaseLoader([]string{sourcesFile})
	loader := NewMCPLoaderWithState(services, baseLoader)

	handlerCalled := false
	var handlerRecord MCPServerProviderRecord

	loader.RegisterEventHandler(func(ctx context.Context, record MCPServerProviderRecord) error {
		handlerCalled = true
		handlerRecord = record
		return nil
	})

	ctx := context.Background()

	// Parse configs
	err = loader.ParseAllConfigs()
	require.NoError(t, err)

	// Set leader mode and perform leader operations
	baseLoader.SetLeader(true)

	// Perform leader operations in background
	leaderDone := make(chan error, 1)
	go func() {
		leaderDone <- loader.PerformLeaderOperations(ctx, mapset.NewSet[string]())
	}()

	// Wait for leader operations to complete
	select {
	case err := <-leaderDone:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for leader operations")
	}

	// Wait for all database writes to complete
	baseLoader.WaitForInflightWrites(5 * time.Second)

	// Verify handler was called
	assert.True(t, handlerCalled)
	assert.NotNil(t, handlerRecord.Server)
	require.Equal(t, 1, len(handlerRecord.Tools))
	assert.Equal(t, "handler-tool", handlerRecord.Tools[0].Name)

}

func TestMCPLoaderRemoveOrphans(t *testing.T) {
	_, services, cleanup := setupMCPLoaderTest(t)
	defer cleanup()

	// Create temporary MCP catalog files
	tmpDir := t.TempDir()

	serversFile := filepath.Join(tmpDir, "servers.yaml")

	// First load with two servers
	err := os.WriteFile(serversFile, []byte(`mcp_servers:
  - name: "server-1"
    description: "First server"
  - name: "server-2"
    description: "Second server"
`), 0644)
	require.NoError(t, err)

	sourcesFile := filepath.Join(tmpDir, "sources.yaml")
	err = os.WriteFile(sourcesFile, []byte(`mcp_catalogs:
  - name: "Orphan Test Catalog"
    id: orphan_test_catalog
    type: yaml
    enabled: true
    properties:
      yamlCatalogPath: `+serversFile+`
`), 0644)
	require.NoError(t, err)

	baseLoader := basecatalog.NewBaseLoader([]string{sourcesFile})
	loader := NewMCPLoaderWithState(services, baseLoader)

	ctx := context.Background()

	// Parse configs
	err = loader.ParseAllConfigs()
	require.NoError(t, err)

	// Set leader mode and perform leader operations
	baseLoader.SetLeader(true)

	// Perform leader operations in background
	leaderDone := make(chan error, 1)
	go func() {
		leaderDone <- loader.PerformLeaderOperations(ctx, mapset.NewSet[string]())
	}()

	// Wait for leader operations to complete
	select {
	case err := <-leaderDone:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for leader operations")
	}

	// Wait for all database writes to complete
	baseLoader.WaitForInflightWrites(5 * time.Second)

	// Verify both servers are loaded
	_, err = services.MCPServerRepository.GetByNameAndVersion("server-1", "")
	require.NoError(t, err)
	_, err = services.MCPServerRepository.GetByNameAndVersion("server-2", "")
	require.NoError(t, err)

	// Update the servers file to remove server-2
	err = os.WriteFile(serversFile, []byte(`mcp_servers:
  - name: "server-1"
    description: "First server (updated)"
`), 0644)
	require.NoError(t, err)

	// Reload
	baseLoader2 := basecatalog.NewBaseLoader([]string{sourcesFile})
	loader2 := NewMCPLoaderWithState(services, baseLoader2)

	// Parse configs
	err = loader2.ParseAllConfigs()
	require.NoError(t, err)

	// Set leader mode and perform leader operations
	baseLoader2.SetLeader(true)

	// Perform leader operations in background
	leaderDone2 := make(chan error, 1)
	go func() {
		leaderDone2 <- loader2.PerformLeaderOperations(ctx, mapset.NewSet[string]())
	}()

	// Wait for leader operations to complete
	select {
	case err := <-leaderDone2:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for second leader operations")
	}

	// Wait for all database writes to complete
	baseLoader2.WaitForInflightWrites(5 * time.Second)

	// Verify server-1 still exists
	_, err = services.MCPServerRepository.GetByNameAndVersion("server-1", "")
	require.NoError(t, err)

	// Verify server-2 was removed (orphaned)
	_, err = services.MCPServerRepository.GetByNameAndVersion("server-2", "")
	assert.Error(t, err)

}

func TestMCPLoaderRespectsContextCancellation(t *testing.T) {
	_, services, cleanup := setupMCPLoaderTest(t)
	defer cleanup()

	// Create temporary MCP catalog files with 3 servers
	tmpDir := t.TempDir()

	serversFile := filepath.Join(tmpDir, "servers.yaml")
	err := os.WriteFile(serversFile, []byte(`mcp_servers:
  - name: "cancel-server-1"
    description: "First server"
  - name: "cancel-server-2"
    description: "Second server"
  - name: "cancel-server-3"
    description: "Third server"
`), 0644)
	require.NoError(t, err)

	sourcesFile := filepath.Join(tmpDir, "sources.yaml")
	err = os.WriteFile(sourcesFile, []byte(`mcp_catalogs:
  - name: "Cancel Test Catalog"
    id: cancel_test_catalog
    type: yaml
    enabled: true
    properties:
      yamlCatalogPath: `+serversFile+`
`), 0644)
	require.NoError(t, err)

	// Create loader with a cancellable context
	baseLoader := basecatalog.NewBaseLoader([]string{sourcesFile})
	loader := NewMCPLoaderWithState(services, baseLoader)

	ctx, cancel := context.WithCancel(context.Background())

	// Track how many servers were processed via event handler
	serverCount := 0
	loader.RegisterEventHandler(func(ctx context.Context, record MCPServerProviderRecord) error {
		serverCount++
		// Cancel context after the first server is processed
		if serverCount == 1 {
			cancel()
		}
		return nil
	})

	// Parse configs
	err = loader.ParseAllConfigs()
	require.NoError(t, err)

	// Set leader mode and perform leader operations
	baseLoader.SetLeader(true)

	// Perform leader operations in background
	leaderDone := make(chan error, 1)
	go func() {
		leaderDone <- loader.PerformLeaderOperations(ctx, mapset.NewSet[string]())
	}()

	// Wait for leader operations to complete
	select {
	case <-leaderDone:
		// Error may or may not propagate depending on how loadAllServers handles it
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for leader operations")
	}

	// Wait for all database writes to complete
	baseLoader.WaitForInflightWrites(5 * time.Second)

	// Verify that not all 3 servers were processed
	assert.LessOrEqual(t, serverCount, 2, "context cancellation should stop processing before all servers are handled")
}

func TestMCPLoaderSavesSourceStatus(t *testing.T) {
	_, services, cleanup := setupMCPLoaderTest(t)
	defer cleanup()

	tmpDir := t.TempDir()

	serversFile := filepath.Join(tmpDir, "servers.yaml")
	err := os.WriteFile(serversFile, []byte(`mcp_servers:
  - name: "status-test-server"
    description: "Server for status test"
`), 0644)
	require.NoError(t, err)

	sourcesFile := filepath.Join(tmpDir, "sources.yaml")
	err = os.WriteFile(sourcesFile, []byte(`mcp_catalogs:
  - name: "Status Test Catalog"
    id: status_test_catalog
    type: yaml
    enabled: true
    properties:
      yamlCatalogPath: `+serversFile+`
`), 0644)
	require.NoError(t, err)

	baseLoader := basecatalog.NewBaseLoader([]string{sourcesFile})
	loader := NewMCPLoaderWithState(services, baseLoader)
	ctx := context.Background()

	err = loader.ParseAllConfigs()
	require.NoError(t, err)

	baseLoader.SetLeader(true)

	leaderDone := make(chan error, 1)
	go func() {
		leaderDone <- loader.PerformLeaderOperations(ctx, mapset.NewSet[string]())
	}()

	select {
	case err := <-leaderDone:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for leader operations")
	}

	baseLoader.WaitForInflightWrites(5 * time.Second)

	// Verify source status was saved as "available"
	source, err := services.CatalogSourceRepository.GetBySourceID("status_test_catalog")
	require.NoError(t, err)
	require.NotNil(t, source)

	props := source.GetProperties()
	require.NotNil(t, props)

	foundStatus := false
	for _, prop := range *props {
		if prop.Name == "status" && prop.StringValue != nil {
			assert.Equal(t, "available", *prop.StringValue)
			foundStatus = true
		}
	}
	assert.True(t, foundStatus, "status property should be set to 'available'")
}

func TestMCPLoaderSavesDisabledSourceStatus(t *testing.T) {
	_, services, cleanup := setupMCPLoaderTest(t)
	defer cleanup()

	tmpDir := t.TempDir()

	serversFile := filepath.Join(tmpDir, "servers.yaml")
	err := os.WriteFile(serversFile, []byte(`mcp_servers:
  - name: "disabled-status-server"
    description: "This server should not be loaded"
`), 0644)
	require.NoError(t, err)

	sourcesFile := filepath.Join(tmpDir, "sources.yaml")
	err = os.WriteFile(sourcesFile, []byte(`mcp_catalogs:
  - name: "Disabled Status Catalog"
    id: disabled_status_catalog
    type: yaml
    enabled: false
    properties:
      yamlCatalogPath: `+serversFile+`
`), 0644)
	require.NoError(t, err)

	baseLoader := basecatalog.NewBaseLoader([]string{sourcesFile})
	loader := NewMCPLoaderWithState(services, baseLoader)
	ctx := context.Background()

	err = loader.ParseAllConfigs()
	require.NoError(t, err)

	baseLoader.SetLeader(true)

	leaderDone := make(chan error, 1)
	go func() {
		leaderDone <- loader.PerformLeaderOperations(ctx, mapset.NewSet[string]())
	}()

	select {
	case err := <-leaderDone:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for leader operations")
	}

	baseLoader.WaitForInflightWrites(5 * time.Second)

	// Verify source status was saved as "disabled"
	source, err := services.CatalogSourceRepository.GetBySourceID("disabled_status_catalog")
	require.NoError(t, err)
	require.NotNil(t, source)

	props := source.GetProperties()
	require.NotNil(t, props)

	foundStatus := false
	for _, prop := range *props {
		if prop.Name == "status" && prop.StringValue != nil {
			assert.Equal(t, "disabled", *prop.StringValue)
			foundStatus = true
		}
	}
	assert.True(t, foundStatus, "status property should be set to 'disabled'")
}

func TestMCPLoaderSavesErrorSourceStatus(t *testing.T) {
	_, services, cleanup := setupMCPLoaderTest(t)
	defer cleanup()

	tmpDir := t.TempDir()

	sourcesFile := filepath.Join(tmpDir, "sources.yaml")
	err := os.WriteFile(sourcesFile, []byte(`mcp_catalogs:
  - name: "Error Status Catalog"
    id: error_status_catalog
    type: unknown_provider_type
    enabled: true
    properties:
      someProp: someValue
`), 0644)
	require.NoError(t, err)

	baseLoader := basecatalog.NewBaseLoader([]string{sourcesFile})
	loader := NewMCPLoaderWithState(services, baseLoader)
	ctx := context.Background()

	err = loader.ParseAllConfigs()
	require.NoError(t, err)

	baseLoader.SetLeader(true)

	leaderDone := make(chan error, 1)
	go func() {
		leaderDone <- loader.PerformLeaderOperations(ctx, mapset.NewSet[string]())
	}()

	select {
	case err := <-leaderDone:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for leader operations")
	}

	baseLoader.WaitForInflightWrites(5 * time.Second)

	// Verify source status was saved as "error"
	source, err := services.CatalogSourceRepository.GetBySourceID("error_status_catalog")
	require.NoError(t, err)
	require.NotNil(t, source)

	props := source.GetProperties()
	require.NotNil(t, props)

	foundStatus := false
	foundError := false
	for _, prop := range *props {
		if prop.Name == "status" && prop.StringValue != nil {
			assert.Equal(t, "error", *prop.StringValue)
			foundStatus = true
		}
		if prop.Name == "error" && prop.StringValue != nil {
			assert.Contains(t, *prop.StringValue, "unknown MCP provider type")
			foundError = true
		}
	}
	assert.True(t, foundStatus, "status property should be set to 'error'")
	assert.True(t, foundError, "error property should be set")
}
