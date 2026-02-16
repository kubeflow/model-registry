package service_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/apiutils"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMCPServerRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupPostgresWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	// Create or get the MCPServer type ID
	typeID := getMCPServerTypeID(t, sharedDB)
	repo := service.NewMCPServerRepository(sharedDB, typeID)

	t.Run("TestSave_Create", func(t *testing.T) {
		// Test creating a new MCP server
		mcpServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name:       apiutils.Of("test-mcp-server"),
				ExternalID: apiutils.Of("mcp-ext-123"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("test-source"),
				},
				{
					Name:        "description",
					StringValue: apiutils.Of("Test MCP server description"),
				},
				{
					Name:        "provider",
					StringValue: apiutils.Of("test-provider"),
				},
				{
					Name:        "version",
					StringValue: apiutils.Of("1.0.0"),
				},
				{
					Name:       "verifiedSource",
					BoolValue:  apiutils.Of(true),
				},
			},
			CustomProperties: &[]dbmodels.Properties{
				{
					Name:        "custom-prop",
					StringValue: apiutils.Of("custom-value"),
				},
			},
		}

		saved, err := repo.Save(mcpServer)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-mcp-server", *saved.GetAttributes().Name)
		assert.Equal(t, "mcp-ext-123", *saved.GetAttributes().ExternalID)
	})

	t.Run("TestSave_Update", func(t *testing.T) {
		// Create initial server
		mcpServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("update-test-server"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("update-source"),
				},
				{
					Name:        "version",
					StringValue: apiutils.Of("1.0.0"),
				},
			},
		}

		saved, err := repo.Save(mcpServer)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Update the server
		updateServer := &models.MCPServerImpl{
			ID: saved.GetID(),
			Attributes: &models.MCPServerAttributes{
				Name:                     apiutils.Of("update-test-server"),
				CreateTimeSinceEpoch:     saved.GetAttributes().CreateTimeSinceEpoch,
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("update-source"),
				},
				{
					Name:        "version",
					StringValue: apiutils.Of("2.0.0"), // Updated version
				},
				{
					Name:        "description",
					StringValue: apiutils.Of("Updated description"),
				},
			},
		}

		updated, err := repo.Save(updateServer)
		require.NoError(t, err)
		assert.Equal(t, *saved.GetID(), *updated.GetID())

		// Verify properties were updated
		require.NotNil(t, updated.GetProperties())
		versionFound := false
		descriptionFound := false
		for _, prop := range *updated.GetProperties() {
			if prop.Name == "version" && prop.StringValue != nil && *prop.StringValue == "2.0.0" {
				versionFound = true
			}
			if prop.Name == "description" && prop.StringValue != nil {
				descriptionFound = true
			}
		}
		assert.True(t, versionFound, "Version should be updated")
		assert.True(t, descriptionFound, "Description should be added")
	})

	t.Run("TestSave_UpsertByNameAndVersion_SameVersion", func(t *testing.T) {
		// Create first server with name and version
		server1 := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("upsert-server"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("upsert-source"),
				},
				{
					Name:        "version",
					StringValue: apiutils.Of("1.0.0"),
				},
				{
					Name:        "description",
					StringValue: apiutils.Of("Initial description"),
				},
			},
		}

		saved1, err := repo.Save(server1)
		require.NoError(t, err)
		require.NotNil(t, saved1.GetID())

		// Save another server with same name and version (without ID)
		// This should UPDATE the existing server
		server2 := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("upsert-server"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("upsert-source"),
				},
				{
					Name:        "version",
					StringValue: apiutils.Of("1.0.0"),
				},
				{
					Name:        "description",
					StringValue: apiutils.Of("Updated description"),
				},
			},
		}

		saved2, err := repo.Save(server2)
		require.NoError(t, err)
		assert.Equal(t, *saved1.GetID(), *saved2.GetID(), "Should update existing server with same name and version")

		// Verify the description was updated
		retrieved, err := repo.GetByNameAndVersion("upsert-server", "1.0.0")
		require.NoError(t, err)
		descFound := false
		for _, prop := range *retrieved.GetProperties() {
			if prop.Name == "description" && prop.StringValue != nil && *prop.StringValue == "Updated description" {
				descFound = true
				break
			}
		}
		assert.True(t, descFound, "Description should be updated")
	})

	t.Run("TestSave_UpsertByNameAndVersion_DifferentVersions", func(t *testing.T) {
		// Create first server with version 1.0.0
		server1 := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("multi-version-server"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("multi-source"),
				},
				{
					Name:        "version",
					StringValue: apiutils.Of("1.0.0"),
				},
			},
		}

		saved1, err := repo.Save(server1)
		require.NoError(t, err)
		require.NotNil(t, saved1.GetID())

		// Save server with same name but different version (2.0.0)
		// This should CREATE a new server
		server2 := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("multi-version-server"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("multi-source"),
				},
				{
					Name:        "version",
					StringValue: apiutils.Of("2.0.0"),
				},
			},
		}

		saved2, err := repo.Save(server2)
		require.NoError(t, err)
		assert.NotEqual(t, *saved1.GetID(), *saved2.GetID(), "Should create new server with different version")

		// Verify both versions exist
		v1, err := repo.GetByNameAndVersion("multi-version-server", "1.0.0")
		require.NoError(t, err)
		assert.Equal(t, *saved1.GetID(), *v1.GetID())

		v2, err := repo.GetByNameAndVersion("multi-version-server", "2.0.0")
		require.NoError(t, err)
		assert.Equal(t, *saved2.GetID(), *v2.GetID())
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// Create a server
		mcpServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("get-by-id-test"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("get-source"),
				},
			},
		}

		saved, err := repo.Save(mcpServer)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Retrieve by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-by-id-test", *retrieved.GetAttributes().Name)

		// Test non-existent ID
		_, err = repo.GetByID(99999)
		assert.ErrorIs(t, err, service.ErrMCPServerNotFound)
	})

	t.Run("TestGetByNameAndVersion_WithVersion", func(t *testing.T) {
		// Create a server with version
		mcpServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("versioned-server"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("version-source"),
				},
				{
					Name:        "version",
					StringValue: apiutils.Of("1.0.0"),
				},
			},
		}

		saved, err := repo.Save(mcpServer)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Retrieve by name and version
		retrieved, err := repo.GetByNameAndVersion("versioned-server", "1.0.0")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "versioned-server", *retrieved.GetAttributes().Name)

		// Test non-existent version
		_, err = repo.GetByNameAndVersion("versioned-server", "2.0.0")
		assert.ErrorIs(t, err, service.ErrMCPServerNotFound)

		// Test non-existent name
		_, err = repo.GetByNameAndVersion("non-existent-server", "1.0.0")
		assert.ErrorIs(t, err, service.ErrMCPServerNotFound)
	})

	t.Run("TestGetByNameAndVersion_NoVersion", func(t *testing.T) {
		// Create a server without version
		mcpServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("unversioned-server"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("no-version-source"),
				},
				{
					Name:        "description",
					StringValue: apiutils.Of("Server without version"),
				},
			},
		}

		saved, err := repo.Save(mcpServer)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Retrieve by name with empty version
		retrieved, err := repo.GetByNameAndVersion("unversioned-server", "")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "unversioned-server", *retrieved.GetAttributes().Name)

		// Test non-existent name with no version
		_, err = repo.GetByNameAndVersion("non-existent-unversioned", "")
		assert.ErrorIs(t, err, service.ErrMCPServerNotFound)
	})

	t.Run("TestList_Basic", func(t *testing.T) {
		// Create multiple servers
		testServers := []*models.MCPServerImpl{
			{
				Attributes: &models.MCPServerAttributes{
					Name: apiutils.Of("list-server-1"),
				},
				Properties: &[]dbmodels.Properties{
					{
						Name:        "source_id",
						StringValue: apiutils.Of("list-source-1"),
					},
					{
						Name:        "provider",
						StringValue: apiutils.Of("provider-a"),
					},
				},
			},
			{
				Attributes: &models.MCPServerAttributes{
					Name: apiutils.Of("list-server-2"),
				},
				Properties: &[]dbmodels.Properties{
					{
						Name:        "source_id",
						StringValue: apiutils.Of("list-source-2"),
					},
					{
						Name:        "provider",
						StringValue: apiutils.Of("provider-b"),
					},
				},
			},
		}

		for _, server := range testServers {
			_, err := repo.Save(server)
			require.NoError(t, err)
		}

		// List all
		listOptions := models.MCPServerListOptions{}
		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 2)
	})

	t.Run("TestList_FilterByName", func(t *testing.T) {
		// Create a unique server
		mcpServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("filter-by-name-unique"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("filter-source"),
				},
			},
		}

		_, err := repo.Save(mcpServer)
		require.NoError(t, err)

		// Filter by name
		nameFilter := "filter-by-name-unique"
		listOptions := models.MCPServerListOptions{
			Name: &nameFilter,
		}
		result, err := repo.List(listOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(result.Items))
		assert.Equal(t, "filter-by-name-unique", *result.Items[0].GetAttributes().Name)
	})

	t.Run("TestList_FilterByQuery", func(t *testing.T) {
		// Create servers with searchable content
		testServers := []*models.MCPServerImpl{
			{
				Attributes: &models.MCPServerAttributes{
					Name: apiutils.Of("query-test-server-xyz"),
				},
				Properties: &[]dbmodels.Properties{
					{
						Name:        "source_id",
						StringValue: apiutils.Of("query-source"),
					},
					{
						Name:        "description",
						StringValue: apiutils.Of("This is a special server"),
					},
				},
			},
			{
				Attributes: &models.MCPServerAttributes{
					Name: apiutils.Of("another-server"),
				},
				Properties: &[]dbmodels.Properties{
					{
						Name:        "source_id",
						StringValue: apiutils.Of("query-source-2"),
					},
					{
						Name:        "provider",
						StringValue: apiutils.Of("special provider"),
					},
				},
			},
		}

		for _, server := range testServers {
			_, err := repo.Save(server)
			require.NoError(t, err)
		}

		// Search for "special" (should find both: one in description, one in provider)
		query := "special"
		listOptions := models.MCPServerListOptions{
			Query: &query,
		}
		result, err := repo.List(listOptions)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result.Items), 2)
	})

	t.Run("TestList_FilterBySourceIDs", func(t *testing.T) {
		// Create servers with different source IDs
		testServers := []*models.MCPServerImpl{
			{
				Attributes: &models.MCPServerAttributes{
					Name: apiutils.Of("source-filter-server-1"),
				},
				Properties: &[]dbmodels.Properties{
					{
						Name:        "source_id",
						StringValue: apiutils.Of("source-alpha"),
					},
				},
			},
			{
				Attributes: &models.MCPServerAttributes{
					Name: apiutils.Of("source-filter-server-2"),
				},
				Properties: &[]dbmodels.Properties{
					{
						Name:        "source_id",
						StringValue: apiutils.Of("source-beta"),
					},
				},
			},
			{
				Attributes: &models.MCPServerAttributes{
					Name: apiutils.Of("source-filter-server-3"),
				},
				Properties: &[]dbmodels.Properties{
					{
						Name:        "source_id",
						StringValue: apiutils.Of("source-gamma"),
					},
				},
			},
		}

		for _, server := range testServers {
			_, err := repo.Save(server)
			require.NoError(t, err)
		}

		// Filter by specific source IDs
		sourceIDs := []string{"source-alpha", "source-beta"}
		listOptions := models.MCPServerListOptions{
			SourceIDs: &sourceIDs,
		}
		result, err := repo.List(listOptions)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result.Items), 2)

		// Verify all results have one of the filtered source IDs
		for _, item := range result.Items {
			sourceIDFound := false
			if item.GetProperties() != nil {
				for _, prop := range *item.GetProperties() {
					if prop.Name == "source_id" && prop.StringValue != nil {
						if *prop.StringValue == "source-alpha" || *prop.StringValue == "source-beta" {
							sourceIDFound = true
							break
						}
					}
				}
			}
			assert.True(t, sourceIDFound, "All results should have filtered source IDs")
		}
	})

	t.Run("TestDeleteBySource", func(t *testing.T) {
		// Create servers from the same source
		sourceID := "delete-test-source"
		testServers := []*models.MCPServerImpl{
			{
				Attributes: &models.MCPServerAttributes{
					Name: apiutils.Of("delete-source-server-1"),
				},
				Properties: &[]dbmodels.Properties{
					{
						Name:        "source_id",
						StringValue: apiutils.Of(sourceID),
					},
				},
			},
			{
				Attributes: &models.MCPServerAttributes{
					Name: apiutils.Of("delete-source-server-2"),
				},
				Properties: &[]dbmodels.Properties{
					{
						Name:        "source_id",
						StringValue: apiutils.Of(sourceID),
					},
				},
			},
		}

		for _, server := range testServers {
			_, err := repo.Save(server)
			require.NoError(t, err)
		}

		// Delete all servers from this source
		err := repo.DeleteBySource(sourceID)
		require.NoError(t, err)

		// Verify they're deleted
		sourceIDs := []string{sourceID}
		listOptions := models.MCPServerListOptions{
			SourceIDs: &sourceIDs,
		}
		result, err := repo.List(listOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(result.Items), "All servers from source should be deleted")
	})

	t.Run("TestDeleteByID", func(t *testing.T) {
		// Create a server
		mcpServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("delete-by-id-test"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("delete-id-source"),
				},
			},
		}

		saved, err := repo.Save(mcpServer)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Delete by ID
		err = repo.DeleteByID(*saved.GetID())
		require.NoError(t, err)

		// Verify it's deleted
		_, err = repo.GetByID(*saved.GetID())
		assert.ErrorIs(t, err, service.ErrMCPServerNotFound)

		// Test deleting non-existent ID
		err = repo.DeleteByID(99999)
		assert.Error(t, err)
	})

	t.Run("TestGetDistinctSourceIDs", func(t *testing.T) {
		// Create servers with different source IDs
		testServers := []*models.MCPServerImpl{
			{
				Attributes: &models.MCPServerAttributes{
					Name: apiutils.Of("distinct-source-1"),
				},
				Properties: &[]dbmodels.Properties{
					{
						Name:        "source_id",
						StringValue: apiutils.Of("distinct-source-alpha"),
					},
				},
			},
			{
				Attributes: &models.MCPServerAttributes{
					Name: apiutils.Of("distinct-source-2"),
				},
				Properties: &[]dbmodels.Properties{
					{
						Name:        "source_id",
						StringValue: apiutils.Of("distinct-source-beta"),
					},
				},
			},
			{
				Attributes: &models.MCPServerAttributes{
					Name: apiutils.Of("distinct-source-3"),
				},
				Properties: &[]dbmodels.Properties{
					{
						Name:        "source_id",
						StringValue: apiutils.Of("distinct-source-alpha"), // Duplicate
					},
				},
			},
		}

		for _, server := range testServers {
			_, err := repo.Save(server)
			require.NoError(t, err)
		}

		// Get distinct source IDs
		sourceIDs, err := repo.GetDistinctSourceIDs()
		require.NoError(t, err)
		require.NotNil(t, sourceIDs)

		// Verify we get unique source IDs
		assert.Contains(t, sourceIDs, "distinct-source-alpha")
		assert.Contains(t, sourceIDs, "distinct-source-beta")

		// Count occurrences to ensure no duplicates
		alphaCount := 0
		for _, sid := range sourceIDs {
			if sid == "distinct-source-alpha" {
				alphaCount++
			}
		}
		assert.Equal(t, 1, alphaCount, "Should only have one occurrence of each source ID")
	})

	t.Run("TestPagination", func(t *testing.T) {
		// Create multiple servers for pagination testing
		for i := 0; i < 5; i++ {
			mcpServer := &models.MCPServerImpl{
				Attributes: &models.MCPServerAttributes{
					Name: apiutils.Of(fmt.Sprintf("pagination-server-%d", i)),
				},
				Properties: &[]dbmodels.Properties{
					{
						Name:        "source_id",
						StringValue: apiutils.Of("pagination-source"),
					},
				},
			}
			_, err := repo.Save(mcpServer)
			require.NoError(t, err)
		}

		// Test first page
		pageSize := int32(2)
		listOptions := models.MCPServerListOptions{
			Pagination: dbmodels.Pagination{
				PageSize: &pageSize,
			},
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(result.Items), 2)
		assert.NotEmpty(t, result.NextPageToken, "Should have next page token")

		// Test second page
		listOptions.Pagination.NextPageToken = &result.NextPageToken
		result2, err := repo.List(listOptions)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(result2.Items), 2)
	})

	t.Run("TestValidation_BaseNameContainsAtSymbol", func(t *testing.T) {
		// Test that base_name cannot contain @ symbol
		server := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("invalid@name"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("test-source"),
				},
			},
		}
		_, err := repo.Save(server)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrBaseNameContainsAtSign)
		assert.Contains(t, err.Error(), "@")
	})

	t.Run("TestValidation_EmptyBaseName", func(t *testing.T) {
		// Test that base_name cannot be empty
		server := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of(""),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("test-source"),
				},
			},
		}
		_, err := repo.Save(server)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrBaseNameEmpty)
	})

	t.Run("TestValidation_WhitespaceOnlyBaseName", func(t *testing.T) {
		// Test that base_name with only whitespace is treated as empty
		server := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("   "),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("test-source"),
				},
			},
		}
		_, err := repo.Save(server)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrBaseNameEmpty)
	})

	t.Run("TestValidation_BaseNameTooLong", func(t *testing.T) {
		// Test that base_name exceeding 255 characters is rejected
		longName := strings.Repeat("a", 256)
		server := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of(longName),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("test-source"),
				},
			},
		}
		_, err := repo.Save(server)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrBaseNameTooLong)
	})

	t.Run("TestValidation_VersionTooLong", func(t *testing.T) {
		// Test that version exceeding 100 characters is rejected
		longVersion := strings.Repeat("1", 101)
		server := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("test-server"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("test-source"),
				},
				{
					Name:        "version",
					StringValue: apiutils.Of(longVersion),
				},
			},
		}
		_, err := repo.Save(server)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrVersionTooLong)
	})

	t.Run("TestValidation_ValidBaseNameWithSpecialChars", func(t *testing.T) {
		// Test that base_name with other special characters (not @) is allowed
		server := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("test-server_v1.0"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("test-source"),
				},
			},
		}
		saved, err := repo.Save(server)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-server_v1.0", *saved.GetAttributes().Name)
	})

	t.Run("TestValidation_BaseNameTrimming", func(t *testing.T) {
		// Test that base_name is trimmed of leading/trailing whitespace
		server := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("  trimmed-server  "),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("test-source"),
				},
			},
		}
		saved, err := repo.Save(server)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Retrieve and verify the name was trimmed
		retrieved, err := repo.GetByNameAndVersion("trimmed-server", "")
		require.NoError(t, err)
		assert.Equal(t, "trimmed-server", *retrieved.GetAttributes().Name)
	})

	t.Run("TestValidation_VersionContainsAtSymbol", func(t *testing.T) {
		// Test that version cannot contain @ symbol
		server := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("test-server"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: apiutils.Of("test-source"),
				},
				{
					Name:        "version",
					StringValue: apiutils.Of("v1.0@beta"),
				},
			},
		}
		_, err := repo.Save(server)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrVersionContainsAtSign)
		assert.Contains(t, err.Error(), "@")
	})
}

// Helper function to get or create the MCPServer type ID
func getMCPServerTypeID(t *testing.T, db *gorm.DB) int32 {
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

// ==============================================================================
// Converter Unit Tests
// ==============================================================================

func TestConvertOpenapiMCPServerToDb(t *testing.T) {
	t.Run("Minimal_RequiredFieldsOnly", func(t *testing.T) {
		// Test with only required fields
		openapiServer := &openapi.MCPServer{
			Name:      "minimal-server",
			ToolCount: 5,
		}

		dbServer := service.ConvertOpenapiMCPServerToDb(openapiServer)
		require.NotNil(t, dbServer)

		attrs := dbServer.GetAttributes()
		require.NotNil(t, attrs)
		assert.Equal(t, "minimal-server", *attrs.Name)

		props := dbServer.GetProperties()
		require.NotNil(t, props)
		// Minimal server should have no additional properties beyond base_name
	})

	t.Run("AllSimpleFields", func(t *testing.T) {
		// Test with all simple string fields
		openapiServer := &openapi.MCPServer{
			Name:             "full-server",
			ToolCount:        10,
			SourceId:         apiutils.Of("test-source"),
			Provider:         apiutils.Of("Acme Corp"),
			Logo:             apiutils.Of("https://example.com/logo.png"),
			Version:          apiutils.Of("1.0.0"),
			License:          apiutils.Of("MIT"),
			LicenseLink:      apiutils.Of("https://opensource.org/licenses/MIT"),
			Readme:           apiutils.Of("# Test Server\nThis is a test"),
			DeploymentMode:   apiutils.Of("local"),
			DocumentationUrl: apiutils.Of("https://docs.example.com"),
			RepositoryUrl:    apiutils.Of("https://github.com/example/repo"),
			SourceCode:       apiutils.Of("https://github.com/example/repo/tree/main"),
			Description:      apiutils.Of("A test MCP server"),
		}

		dbServer := service.ConvertOpenapiMCPServerToDb(openapiServer)
		props := dbServer.GetProperties()
		require.NotNil(t, props)

		// Verify all simple string fields are in properties
		propMap := make(map[string]string)
		for _, prop := range *props {
			if prop.StringValue != nil {
				propMap[prop.Name] = *prop.StringValue
			}
		}

		assert.Equal(t, "test-source", propMap["source_id"])
		assert.Equal(t, "Acme Corp", propMap["provider"])
		assert.Equal(t, "https://example.com/logo.png", propMap["logo"])
		assert.Equal(t, "1.0.0", propMap["version"])
		assert.Equal(t, "MIT", propMap["license"])
		assert.Equal(t, "https://opensource.org/licenses/MIT", propMap["license_link"])
		assert.Equal(t, "# Test Server\nThis is a test", propMap["readme"])
		assert.Equal(t, "local", propMap["deploymentMode"])
		assert.Equal(t, "https://docs.example.com", propMap["documentationUrl"])
		assert.Equal(t, "https://github.com/example/repo", propMap["repositoryUrl"])
		assert.Equal(t, "https://github.com/example/repo/tree/main", propMap["sourceCode"])
		assert.Equal(t, "A test MCP server", propMap["description"])
	})

	t.Run("ArrayFields", func(t *testing.T) {
		// Test array fields (tags, transports)
		openapiServer := &openapi.MCPServer{
			Name:       "array-server",
			ToolCount:  3,
			Tags:       []string{"monitoring", "observability", "apm"},
			Transports: []string{"stdio", "http"},
		}

		dbServer := service.ConvertOpenapiMCPServerToDb(openapiServer)
		props := dbServer.GetProperties()
		require.NotNil(t, props)

		// Find and verify array properties
		var tagsJSON, transportsJSON string
		for _, prop := range *props {
			if prop.Name == "tags" && prop.StringValue != nil {
				tagsJSON = *prop.StringValue
			}
			if prop.Name == "transports" && prop.StringValue != nil {
				transportsJSON = *prop.StringValue
			}
		}

		assert.Contains(t, tagsJSON, "monitoring")
		assert.Contains(t, tagsJSON, "observability")
		assert.Contains(t, tagsJSON, "apm")
		assert.Contains(t, transportsJSON, "stdio")
		assert.Contains(t, transportsJSON, "http")
	})

	t.Run("TimeFields", func(t *testing.T) {
		// Test time field conversions
		// NOTE: Time fields are stored as ISO 8601 strings for consistency with other timestamps
		publishedDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		lastUpdated := time.Date(2024, 1, 20, 12, 30, 0, 0, time.UTC)

		openapiServer := &openapi.MCPServer{
			Name:          "time-server",
			ToolCount:     1,
			PublishedDate: &publishedDate,
			LastUpdated:   &lastUpdated,
		}

		dbServer := service.ConvertOpenapiMCPServerToDb(openapiServer)
		props := dbServer.GetProperties()
		require.NotNil(t, props)

		// Find time properties
		var publishedStr, updatedStr string
		for _, prop := range *props {
			if prop.Name == "publishedDate" && prop.StringValue != nil {
				publishedStr = *prop.StringValue
			}
			if prop.Name == "lastUpdated" && prop.StringValue != nil {
				updatedStr = *prop.StringValue
			}
		}

		// Verify conversion to ISO 8601 strings
		assert.NotEmpty(t, publishedStr)
		assert.NotEmpty(t, updatedStr)
		assert.Equal(t, "2024-01-10T00:00:00Z", publishedStr)
		assert.Equal(t, "2024-01-20T12:30:00Z", updatedStr)
	})

	t.Run("SecurityIndicators", func(t *testing.T) {
		// Test SecurityIndicators expansion to 4 boolean properties
		openapiServer := &openapi.MCPServer{
			Name:      "secure-server",
			ToolCount: 2,
			SecurityIndicators: &openapi.MCPSecurityIndicator{
				VerifiedSource: apiutils.Of(true),
				SecureEndpoint: apiutils.Of(true),
				Sast:           apiutils.Of(false),
				ReadOnlyTools:  apiutils.Of(true),
			},
		}

		dbServer := service.ConvertOpenapiMCPServerToDb(openapiServer)
		props := dbServer.GetProperties()
		require.NotNil(t, props)

		// Find boolean properties
		boolProps := make(map[string]bool)
		for _, prop := range *props {
			if prop.BoolValue != nil {
				boolProps[prop.Name] = *prop.BoolValue
			}
		}

		assert.True(t, boolProps["verifiedSource"])
		assert.True(t, boolProps["secureEndpoint"])
		assert.False(t, boolProps["sast"])
		assert.True(t, boolProps["readOnlyTools"])
	})

	t.Run("ComplexObjects_JSON", func(t *testing.T) {
		// Test complex objects stored as JSON
		openapiServer := &openapi.MCPServer{
			Name:      "json-server",
			ToolCount: 1,
			Endpoints: &openapi.MCPEndpoints{
				Http: apiutils.Of("http://localhost:8080"),
				Sse:  apiutils.Of("http://localhost:8080/events"),
			},
			Artifacts: []openapi.MCPArtifact{
				{
					Uri: "oci://registry.example.com/server:v1.0.0",
				},
			},
			RuntimeMetadata: &openapi.MCPRuntimeMetadata{
				DefaultPort: apiutils.Of(int32(8080)),
				DefaultArgs: []string{"--log-level", "info"},
			},
		}

		dbServer := service.ConvertOpenapiMCPServerToDb(openapiServer)
		props := dbServer.GetProperties()
		require.NotNil(t, props)

		// Find JSON properties
		jsonProps := make(map[string]string)
		for _, prop := range *props {
			if prop.StringValue != nil {
				jsonProps[prop.Name] = *prop.StringValue
			}
		}

		// Verify JSON encoding
		assert.Contains(t, jsonProps["endpoints"], "localhost:8080")
		assert.Contains(t, jsonProps["artifacts"], "oci://registry.example.com")
		assert.Contains(t, jsonProps["runtimeMetadata"], "8080")
		assert.Contains(t, jsonProps["runtimeMetadata"], "--log-level")
	})
}

func TestConvertDbMCPServerToOpenapi(t *testing.T) {
	t.Run("BasicFields", func(t *testing.T) {
		// Test basic field conversion
		dbServer := &models.MCPServerImpl{
			ID: apiutils.Of(int32(123)),
			Attributes: &models.MCPServerAttributes{
				Name:                     apiutils.Of("test-server"),
				ExternalID:               apiutils.Of("ext-123"),
				CreateTimeSinceEpoch:     apiutils.Of(int64(1704067200000)),
				LastUpdateTimeSinceEpoch: apiutils.Of(int64(1704153600000)),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "source_id", StringValue: apiutils.Of("test-source")},
				{Name: "provider", StringValue: apiutils.Of("Test Provider")},
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
				{Name: "description", StringValue: apiutils.Of("Test description")},
			},
		}

		openapiServer := service.ConvertDbMCPServerToOpenapi(dbServer)
		require.NotNil(t, openapiServer)

		assert.Equal(t, "test-server", openapiServer.Name)
		assert.Equal(t, "123", *openapiServer.Id)
		assert.Equal(t, "ext-123", *openapiServer.ExternalId)
		assert.Equal(t, "1704067200000", *openapiServer.CreateTimeSinceEpoch)
		assert.Equal(t, "1704153600000", *openapiServer.LastUpdateTimeSinceEpoch)
		assert.Equal(t, "test-source", *openapiServer.SourceId)
		assert.Equal(t, "Test Provider", *openapiServer.Provider)
		assert.Equal(t, "1.0.0", *openapiServer.Version)
		assert.Equal(t, "Test description", *openapiServer.Description)
		assert.Equal(t, int32(0), openapiServer.ToolCount) // Default, not computed in converter
	})

	t.Run("ArrayFields", func(t *testing.T) {
		// Test array field decoding from JSON
		dbServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("array-server"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "tags", StringValue: apiutils.Of(`["monitoring","observability"]`)},
				{Name: "transports", StringValue: apiutils.Of(`["stdio","http","sse"]`)},
			},
		}

		openapiServer := service.ConvertDbMCPServerToOpenapi(dbServer)
		require.NotNil(t, openapiServer)

		require.Len(t, openapiServer.Tags, 2)
		assert.Equal(t, "monitoring", openapiServer.Tags[0])
		assert.Equal(t, "observability", openapiServer.Tags[1])

		require.Len(t, openapiServer.Transports, 3)
		assert.Equal(t, "stdio", openapiServer.Transports[0])
		assert.Equal(t, "http", openapiServer.Transports[1])
		assert.Equal(t, "sse", openapiServer.Transports[2])
	})

	t.Run("TimeFields", func(t *testing.T) {
		// Test time field conversion from int32 milliseconds
		// NOTE: Time fields stored as ISO 8601 strings for consistency
		publishedStr := "2024-01-10T00:00:00Z"
		updatedStr := "2024-01-20T12:30:00Z"

		dbServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("time-server"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "publishedDate", StringValue: &publishedStr},
				{Name: "lastUpdated", StringValue: &updatedStr},
			},
		}

		openapiServer := service.ConvertDbMCPServerToOpenapi(dbServer)
		require.NotNil(t, openapiServer)

		require.NotNil(t, openapiServer.PublishedDate)
		require.NotNil(t, openapiServer.LastUpdated)

		// Verify the dates are parsed correctly
		assert.Equal(t, "2024-01-10T00:00:00Z", openapiServer.PublishedDate.Format(time.RFC3339))
		assert.Equal(t, "2024-01-20T12:30:00Z", openapiServer.LastUpdated.Format(time.RFC3339))
		assert.True(t, openapiServer.LastUpdated.After(*openapiServer.PublishedDate))
	})

	t.Run("SecurityIndicators", func(t *testing.T) {
		// Test SecurityIndicators reconstruction from 4 boolean properties
		dbServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("secure-server"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "verifiedSource", BoolValue: apiutils.Of(true)},
				{Name: "secureEndpoint", BoolValue: apiutils.Of(true)},
				{Name: "sast", BoolValue: apiutils.Of(false)},
				{Name: "readOnlyTools", BoolValue: apiutils.Of(true)},
			},
		}

		openapiServer := service.ConvertDbMCPServerToOpenapi(dbServer)
		require.NotNil(t, openapiServer)
		require.NotNil(t, openapiServer.SecurityIndicators)

		assert.True(t, *openapiServer.SecurityIndicators.VerifiedSource)
		assert.True(t, *openapiServer.SecurityIndicators.SecureEndpoint)
		assert.False(t, *openapiServer.SecurityIndicators.Sast)
		assert.True(t, *openapiServer.SecurityIndicators.ReadOnlyTools)
	})

	t.Run("ComplexObjects_JSON", func(t *testing.T) {
		// Test complex object decoding from JSON
		endpointsJSON := `{"http":"http://localhost:8080","sse":"http://localhost:8080/events"}`
		artifactsJSON := `[{"uri":"oci://registry.example.com/server:v1.0.0"}]`
		runtimeJSON := `{"defaultPort":8080,"defaultArgs":["--log-level","info"]}`

		dbServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("json-server"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "endpoints", StringValue: &endpointsJSON},
				{Name: "artifacts", StringValue: &artifactsJSON},
				{Name: "runtimeMetadata", StringValue: &runtimeJSON},
			},
		}

		openapiServer := service.ConvertDbMCPServerToOpenapi(dbServer)
		require.NotNil(t, openapiServer)

		// Verify endpoints
		require.NotNil(t, openapiServer.Endpoints)
		assert.Equal(t, "http://localhost:8080", *openapiServer.Endpoints.Http)
		assert.Equal(t, "http://localhost:8080/events", *openapiServer.Endpoints.Sse)

		// Verify artifacts
		require.Len(t, openapiServer.Artifacts, 1)
		assert.Equal(t, "oci://registry.example.com/server:v1.0.0", openapiServer.Artifacts[0].Uri)

		// Verify runtime metadata
		require.NotNil(t, openapiServer.RuntimeMetadata)
		assert.Equal(t, int32(8080), *openapiServer.RuntimeMetadata.DefaultPort)
		require.Len(t, openapiServer.RuntimeMetadata.DefaultArgs, 2)
		assert.Equal(t, "--log-level", openapiServer.RuntimeMetadata.DefaultArgs[0])
	})

	t.Run("NoSecurityIndicators_WhenAllFalse", func(t *testing.T) {
		// Test that SecurityIndicators is nil when no properties exist
		dbServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("no-security-server"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}

		openapiServer := service.ConvertDbMCPServerToOpenapi(dbServer)
		require.NotNil(t, openapiServer)
		assert.Nil(t, openapiServer.SecurityIndicators)
	})
}

func TestRoundTrip_OpenapiToDbToOpenapi(t *testing.T) {
	t.Run("FullServer_RoundTrip", func(t *testing.T) {
		// Create a fully populated OpenAPI server
		// Using realistic dates - timestamps stored as ISO 8601 strings
		publishedDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		lastUpdated := time.Date(2024, 1, 20, 12, 30, 0, 0, time.UTC)

		original := &openapi.MCPServer{
			Name:             "roundtrip-server",
			ToolCount:        15,
			SourceId:         apiutils.Of("source-123"),
			Provider:         apiutils.Of("RoundTrip Corp"),
			Logo:             apiutils.Of("https://example.com/logo.png"),
			Version:          apiutils.Of("2.5.0"),
			License:          apiutils.Of("Apache-2.0"),
			LicenseLink:      apiutils.Of("https://www.apache.org/licenses/LICENSE-2.0"),
			Readme:           apiutils.Of("# RoundTrip Server\nFull featured test"),
			DeploymentMode:   apiutils.Of("remote"),
			DocumentationUrl: apiutils.Of("https://docs.example.com/roundtrip"),
			RepositoryUrl:    apiutils.Of("https://github.com/example/roundtrip"),
			SourceCode:       apiutils.Of("https://github.com/example/roundtrip/tree/v2.5.0"),
			Description:      apiutils.Of("A comprehensive roundtrip test server"),
			Tags:             []string{"test", "roundtrip", "validation"},
			Transports:       []string{"stdio", "http", "sse"},
			PublishedDate:    &publishedDate,
			LastUpdated:      &lastUpdated,
			SecurityIndicators: &openapi.MCPSecurityIndicator{
				VerifiedSource: apiutils.Of(true),
				SecureEndpoint: apiutils.Of(true),
				Sast:           apiutils.Of(true),
				ReadOnlyTools:  apiutils.Of(false),
			},
			Endpoints: &openapi.MCPEndpoints{
				Http: apiutils.Of("https://api.example.com"),
				Sse:  apiutils.Of("https://api.example.com/events"),
			},
		}

		// Convert to DB
		dbServer := service.ConvertOpenapiMCPServerToDb(original)
		require.NotNil(t, dbServer)

		// Convert back to OpenAPI
		result := service.ConvertDbMCPServerToOpenapi(dbServer)
		require.NotNil(t, result)

		// Verify all fields preserved (except toolCount which is computed)
		assert.Equal(t, original.Name, result.Name)
		assert.Equal(t, original.SourceId, result.SourceId)
		assert.Equal(t, original.Provider, result.Provider)
		assert.Equal(t, original.Logo, result.Logo)
		assert.Equal(t, original.Version, result.Version)
		assert.Equal(t, original.License, result.License)
		assert.Equal(t, original.LicenseLink, result.LicenseLink)
		assert.Equal(t, original.Readme, result.Readme)
		assert.Equal(t, original.DeploymentMode, result.DeploymentMode)
		assert.Equal(t, original.DocumentationUrl, result.DocumentationUrl)
		assert.Equal(t, original.RepositoryUrl, result.RepositoryUrl)
		assert.Equal(t, original.SourceCode, result.SourceCode)
		assert.Equal(t, original.Description, result.Description)

		// Verify arrays
		assert.Equal(t, original.Tags, result.Tags)
		assert.Equal(t, original.Transports, result.Transports)

		// Verify security indicators
		require.NotNil(t, result.SecurityIndicators)
		assert.Equal(t, *original.SecurityIndicators.VerifiedSource, *result.SecurityIndicators.VerifiedSource)
		assert.Equal(t, *original.SecurityIndicators.SecureEndpoint, *result.SecurityIndicators.SecureEndpoint)
		assert.Equal(t, *original.SecurityIndicators.Sast, *result.SecurityIndicators.Sast)
		assert.Equal(t, *original.SecurityIndicators.ReadOnlyTools, *result.SecurityIndicators.ReadOnlyTools)

		// Verify endpoints
		require.NotNil(t, result.Endpoints)
		assert.Equal(t, *original.Endpoints.Http, *result.Endpoints.Http)
		assert.Equal(t, *original.Endpoints.Sse, *result.Endpoints.Sse)
	})

	t.Run("MinimalServer_RoundTrip", func(t *testing.T) {
		// Test with minimal fields to ensure nil handling
		original := &openapi.MCPServer{
			Name:      "minimal",
			ToolCount: 0,
		}

		dbServer := service.ConvertOpenapiMCPServerToDb(original)
		result := service.ConvertDbMCPServerToOpenapi(dbServer)

		assert.Equal(t, original.Name, result.Name)
		assert.Nil(t, result.SourceId)
		assert.Nil(t, result.Provider)
		assert.Nil(t, result.SecurityIndicators)
		assert.Nil(t, result.Endpoints)
	})
}

func TestConvertOpenapiMCPToolToDb(t *testing.T) {
	t.Run("Minimal_RequiredFieldsOnly", func(t *testing.T) {
		// Test with only required fields: name and accessType
		openapiTool := &openapi.MCPTool{
			Name:       "test-tool",
			AccessType: "read-only",
		}

		dbTool := service.ConvertOpenapiMCPToolToDb(openapiTool)
		require.NotNil(t, dbTool)

		attr := dbTool.GetAttributes()
		require.NotNil(t, attr)
		require.NotNil(t, attr.Name)
		assert.Equal(t, "test-tool", *attr.Name)

		// Verify accessType property
		props := dbTool.GetProperties()
		require.NotNil(t, props)

		var accessType string
		for _, prop := range *props {
			if prop.Name == "accessType" && prop.StringValue != nil {
				accessType = *prop.StringValue
			}
		}
		assert.Equal(t, "read-only", accessType)
	})

	t.Run("WithDescription", func(t *testing.T) {
		// Test with description field
		description := "A tool that reads data"
		openapiTool := &openapi.MCPTool{
			Name:        "read-tool",
			AccessType:  "read-only",
			Description: &description,
		}

		dbTool := service.ConvertOpenapiMCPToolToDb(openapiTool)
		props := dbTool.GetProperties()
		require.NotNil(t, props)

		var descValue string
		for _, prop := range *props {
			if prop.Name == "description" && prop.StringValue != nil {
				descValue = *prop.StringValue
			}
		}
		assert.Equal(t, "A tool that reads data", descValue)
	})

	t.Run("WithParameters", func(t *testing.T) {
		// Test with parameters array
		openapiTool := &openapi.MCPTool{
			Name:       "param-tool",
			AccessType: "read-write",
			Parameters: []openapi.MCPToolParameter{
				{
					Name:     "input",
					Type:     "string",
					Required: true,
				},
				{
					Name:     "count",
					Type:     "integer",
					Required: false,
				},
			},
		}

		dbTool := service.ConvertOpenapiMCPToolToDb(openapiTool)
		props := dbTool.GetProperties()
		require.NotNil(t, props)

		// Find and verify parameters property (stored as JSON)
		var paramsJSON string
		for _, prop := range *props {
			if prop.Name == "parameters" && prop.StringValue != nil {
				paramsJSON = *prop.StringValue
			}
		}
		assert.NotEmpty(t, paramsJSON)

		// Verify JSON can be parsed back
		var params []openapi.MCPToolParameter
		err := json.Unmarshal([]byte(paramsJSON), &params)
		require.NoError(t, err)
		assert.Len(t, params, 2)
		assert.Equal(t, "input", params[0].Name)
		assert.Equal(t, "count", params[1].Name)
	})

	t.Run("WithExternalId", func(t *testing.T) {
		// Test with ExternalId field
		externalId := "ext-tool-123"
		openapiTool := &openapi.MCPTool{
			Name:       "external-tool",
			AccessType: "read-only",
			ExternalId: &externalId,
		}

		dbTool := service.ConvertOpenapiMCPToolToDb(openapiTool)
		props := dbTool.GetProperties()
		require.NotNil(t, props)

		var extId string
		for _, prop := range *props {
			if prop.Name == "externalId" && prop.StringValue != nil {
				extId = *prop.StringValue
			}
		}
		assert.Equal(t, "ext-tool-123", extId)
	})
}

func TestConvertDbMCPToolToOpenapi(t *testing.T) {
	t.Run("BasicFields", func(t *testing.T) {
		// Test basic conversion from DB to OpenAPI
		dbTool := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{
				Name: apiutils.Of("db-tool"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-write")},
			},
		}

		// Set ID
		id := int32(42)
		dbTool.ID = &id

		openapiTool := service.ConvertDbMCPToolToOpenapi(dbTool)
		require.NotNil(t, openapiTool)
		assert.Equal(t, "db-tool", openapiTool.Name)
		assert.Equal(t, "read-write", openapiTool.AccessType)
		assert.NotNil(t, openapiTool.Id)
		assert.Equal(t, "42", *openapiTool.Id)
	})

	t.Run("WithDescription", func(t *testing.T) {
		dbTool := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{
				Name: apiutils.Of("desc-tool"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-only")},
				{Name: "description", StringValue: apiutils.Of("Test description")},
			},
		}

		openapiTool := service.ConvertDbMCPToolToOpenapi(dbTool)
		require.NotNil(t, openapiTool)
		require.NotNil(t, openapiTool.Description)
		assert.Equal(t, "Test description", *openapiTool.Description)
	})

	t.Run("WithParameters", func(t *testing.T) {
		// Store parameters as JSON
		paramsJSON := `[{"name":"input","type":"string","required":true},{"name":"count","type":"integer"}]`
		dbTool := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{
				Name: apiutils.Of("param-tool"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-write")},
				{Name: "parameters", StringValue: &paramsJSON},
			},
		}

		openapiTool := service.ConvertDbMCPToolToOpenapi(dbTool)
		require.NotNil(t, openapiTool)
		require.NotNil(t, openapiTool.Parameters)
		assert.Len(t, openapiTool.Parameters, 2)
		assert.Equal(t, "input", openapiTool.Parameters[0].Name)
		assert.Equal(t, "count", openapiTool.Parameters[1].Name)
	})

	t.Run("WithTimestamps", func(t *testing.T) {
		createTime := int64(1704067200000)
		updateTime := int64(1704153600000)

		dbTool := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{
				Name:                     apiutils.Of("time-tool"),
				CreateTimeSinceEpoch:     &createTime,
				LastUpdateTimeSinceEpoch: &updateTime,
			},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-only")},
			},
		}

		openapiTool := service.ConvertDbMCPToolToOpenapi(dbTool)
		require.NotNil(t, openapiTool)
		require.NotNil(t, openapiTool.CreateTimeSinceEpoch)
		require.NotNil(t, openapiTool.LastUpdateTimeSinceEpoch)
		assert.Equal(t, "1704067200000", *openapiTool.CreateTimeSinceEpoch)
		assert.Equal(t, "1704153600000", *openapiTool.LastUpdateTimeSinceEpoch)
	})
}

func TestRoundTrip_OpenapiToolToDbToOpenapi(t *testing.T) {
	t.Run("FullTool_RoundTrip", func(t *testing.T) {
		// Create a fully populated tool
		description := "Full test tool"
		externalId := "ext-123"
		original := &openapi.MCPTool{
			Name:        "roundtrip-tool",
			AccessType:  "read-write",
			Description: &description,
			ExternalId:  &externalId,
			Parameters: []openapi.MCPToolParameter{
				{
					Name:        "param1",
					Type:        "string",
					Required:    true,
					Description: apiutils.Of("First parameter"),
				},
			},
		}

		// Convert to DB and back
		dbTool := service.ConvertOpenapiMCPToolToDb(original)
		result := service.ConvertDbMCPToolToOpenapi(dbTool)

		// Verify all fields match
		assert.Equal(t, original.Name, result.Name)
		assert.Equal(t, original.AccessType, result.AccessType)
		require.NotNil(t, result.Description)
		assert.Equal(t, *original.Description, *result.Description)
		require.NotNil(t, result.ExternalId)
		assert.Equal(t, *original.ExternalId, *result.ExternalId)
		require.NotNil(t, result.Parameters)
		assert.Len(t, result.Parameters, 1)
		assert.Equal(t, "param1", result.Parameters[0].Name)
	})

	t.Run("MinimalTool_RoundTrip", func(t *testing.T) {
		// Test with minimal required fields
		original := &openapi.MCPTool{
			Name:       "minimal-tool",
			AccessType: "read-only",
		}

		dbTool := service.ConvertOpenapiMCPToolToDb(original)
		result := service.ConvertDbMCPToolToOpenapi(dbTool)

		assert.Equal(t, original.Name, result.Name)
		assert.Equal(t, original.AccessType, result.AccessType)
		assert.Nil(t, result.Description)
		assert.Nil(t, result.ExternalId)
		assert.Empty(t, result.Parameters)
	})
}

func TestConvertDbMCPServerWithToolsToOpenapi(t *testing.T) {
	t.Run("WithTools_CorrectToolCountAndToolsArray", func(t *testing.T) {
		dbServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("test-server"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
				{Name: "description", StringValue: apiutils.Of("Test server")},
			},
		}

		// Create test tools
		tools := []models.MCPServerTool{
			&models.MCPServerToolImpl{
				Attributes: &models.MCPServerToolAttributes{
					Name: apiutils.Of("read-tool"),
				},
				Properties: &[]dbmodels.Properties{
					{Name: "accessType", StringValue: apiutils.Of("read-only")},
					{Name: "description", StringValue: apiutils.Of("Reads data")},
				},
			},
			&models.MCPServerToolImpl{
				Attributes: &models.MCPServerToolAttributes{
					Name: apiutils.Of("write-tool"),
				},
				Properties: &[]dbmodels.Properties{
					{Name: "accessType", StringValue: apiutils.Of("read-write")},
					{Name: "description", StringValue: apiutils.Of("Writes data")},
				},
			},
			&models.MCPServerToolImpl{
				Attributes: &models.MCPServerToolAttributes{
					Name: apiutils.Of("admin-tool"),
				},
				Properties: &[]dbmodels.Properties{
					{Name: "accessType", StringValue: apiutils.Of("admin")},
				},
			},
		}

		openapiServer := service.ConvertDbMCPServerWithToolsToOpenapi(dbServer, tools)

		require.NotNil(t, openapiServer)
		assert.Equal(t, "test-server", openapiServer.Name)
		assert.Equal(t, int32(3), openapiServer.ToolCount)
		require.NotNil(t, openapiServer.Tools)
		assert.Len(t, openapiServer.Tools, 3)

		// Verify tool details
		assert.Equal(t, "read-tool", openapiServer.Tools[0].Name)
		assert.Equal(t, "read-only", openapiServer.Tools[0].AccessType)
		assert.Equal(t, "Reads data", *openapiServer.Tools[0].Description)

		assert.Equal(t, "write-tool", openapiServer.Tools[1].Name)
		assert.Equal(t, "read-write", openapiServer.Tools[1].AccessType)
		assert.Equal(t, "Writes data", *openapiServer.Tools[1].Description)

		assert.Equal(t, "admin-tool", openapiServer.Tools[2].Name)
		assert.Equal(t, "admin", openapiServer.Tools[2].AccessType)
	})

	t.Run("WithNilTools_ZeroToolCount", func(t *testing.T) {
		dbServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("empty-server"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}

		openapiServer := service.ConvertDbMCPServerWithToolsToOpenapi(dbServer, nil)

		require.NotNil(t, openapiServer)
		assert.Equal(t, "empty-server", openapiServer.Name)
		assert.Equal(t, int32(0), openapiServer.ToolCount)
		assert.Nil(t, openapiServer.Tools)
	})

	t.Run("WithEmptyToolsSlice_ZeroToolCount", func(t *testing.T) {
		dbServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("empty-server"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}

		openapiServer := service.ConvertDbMCPServerWithToolsToOpenapi(dbServer, []models.MCPServerTool{})

		require.NotNil(t, openapiServer)
		assert.Equal(t, "empty-server", openapiServer.Name)
		assert.Equal(t, int32(0), openapiServer.ToolCount)
		assert.Nil(t, openapiServer.Tools) // Empty slice not set
	})

	t.Run("WithSingleTool_ToolCountOne", func(t *testing.T) {
		dbServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("single-tool-server"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}

		tools := []models.MCPServerTool{
			&models.MCPServerToolImpl{
				Attributes: &models.MCPServerToolAttributes{
					Name: apiutils.Of("only-tool"),
				},
				Properties: &[]dbmodels.Properties{
					{Name: "accessType", StringValue: apiutils.Of("read-only")},
				},
			},
		}

		openapiServer := service.ConvertDbMCPServerWithToolsToOpenapi(dbServer, tools)

		require.NotNil(t, openapiServer)
		assert.Equal(t, int32(1), openapiServer.ToolCount)
		require.NotNil(t, openapiServer.Tools)
		assert.Len(t, openapiServer.Tools, 1)
		assert.Equal(t, "only-tool", openapiServer.Tools[0].Name)
	})

	t.Run("SimpleConverter_ReturnsZeroToolCount", func(t *testing.T) {
		// Test that the simple ConvertDbMCPServerToOpenapi (without tools) returns 0
		dbServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("test-server"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}

		openapiServer := service.ConvertDbMCPServerToOpenapi(dbServer)

		require.NotNil(t, openapiServer)
		assert.Equal(t, "test-server", openapiServer.Name)
		assert.Equal(t, int32(0), openapiServer.ToolCount)
		assert.Nil(t, openapiServer.Tools)
	})

	t.Run("ToolConversionSkipsNilTools", func(t *testing.T) {
		// Test that nil tools in the slice are skipped
		dbServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("test-server"),
			},
		}

		// Create a tool with missing required name (will return nil from converter)
		tools := []models.MCPServerTool{
			&models.MCPServerToolImpl{
				Attributes: &models.MCPServerToolAttributes{
					Name: nil, // Missing required field
				},
			},
			&models.MCPServerToolImpl{
				Attributes: &models.MCPServerToolAttributes{
					Name: apiutils.Of("valid-tool"),
				},
				Properties: &[]dbmodels.Properties{
					{Name: "accessType", StringValue: apiutils.Of("read-only")},
				},
			},
		}

		openapiServer := service.ConvertDbMCPServerWithToolsToOpenapi(dbServer, tools)

		require.NotNil(t, openapiServer)
		// Only the valid tool should be included
		assert.Equal(t, int32(1), openapiServer.ToolCount)
		require.Len(t, openapiServer.Tools, 1)
		assert.Equal(t, "valid-tool", openapiServer.Tools[0].Name)
	})
}
