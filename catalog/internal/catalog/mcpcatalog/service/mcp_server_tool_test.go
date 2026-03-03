package service

import (
	"fmt"
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/catalog/mcpcatalog/models"
	"github.com/kubeflow/model-registry/internal/apiutils"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==============================================================================
// MCPServerTool Database Persistence Tests
// ==============================================================================

func TestMCPServerToolRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupPostgresWithMigrations(t, testDatastoreSpec())
	defer cleanup()

	mcpServerTypeID := getMCPServerTypeID(t, sharedDB)
	mcpToolTypeID := getMCPServerToolTypeID(t, sharedDB)
	serverRepo := NewMCPServerRepository(sharedDB, mcpServerTypeID)
	toolRepo := NewMCPServerToolRepository(sharedDB, mcpToolTypeID)

	t.Run("TestSaveTool_Create", func(t *testing.T) {
		// Create parent MCP server
		parentServer := &models.MCPServerImpl{
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
					StringValue: apiutils.Of("1.0.0"),
				},
			},
		}

		savedServer, err := serverRepo.Save(parentServer)
		require.NoError(t, err)
		require.NotNil(t, savedServer.GetID())
		parentID := *savedServer.GetID()

		// Create tool for the server
		tool := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{
				Name: apiutils.Of("test-tool"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "accessType",
					StringValue: apiutils.Of("read-write"),
				},
				{
					Name:        "description",
					StringValue: apiutils.Of("A test tool"),
				},
			},
		}

		savedTool, err := toolRepo.Save(tool, &parentID)
		require.NoError(t, err)
		require.NotNil(t, savedTool.GetID())
		assert.Equal(t, "test-tool", *savedTool.GetAttributes().Name)

		// Verify properties
		props := savedTool.GetProperties()
		require.NotNil(t, props)
		accessTypeFound := false
		for _, prop := range *props {
			if prop.Name == "accessType" && prop.StringValue != nil && *prop.StringValue == "read-write" {
				accessTypeFound = true
			}
		}
		assert.True(t, accessTypeFound, "accessType property should be saved")
	})

	t.Run("TestGetToolByID", func(t *testing.T) {
		// Create parent server
		parentServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("parent-server-get"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "source_id", StringValue: apiutils.Of("source-get")},
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}
		savedServer, err := serverRepo.Save(parentServer)
		require.NoError(t, err)
		parentID := *savedServer.GetID()

		// Create tool
		tool := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{
				Name: apiutils.Of("get-test-tool"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-only")},
			},
		}
		savedTool, err := toolRepo.Save(tool, &parentID)
		require.NoError(t, err)

		// Retrieve by ID
		retrieved, err := toolRepo.GetByID(*savedTool.GetID())
		require.NoError(t, err)
		assert.Equal(t, *savedTool.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-test-tool", *retrieved.GetAttributes().Name)

		// Test non-existent ID
		_, err = toolRepo.GetByID(99999)
		assert.ErrorIs(t, err, ErrMCPServerToolNotFound)
	})

	t.Run("TestListToolsByParent", func(t *testing.T) {
		// Create parent server
		parentServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{
				Name: apiutils.Of("server-with-tools-list"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "source_id", StringValue: apiutils.Of("source-list")},
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}
		savedServer, err := serverRepo.Save(parentServer)
		require.NoError(t, err)
		parentID := *savedServer.GetID()

		// Create multiple tools
		tool1 := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{Name: apiutils.Of("tool-1")},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-only")},
			},
		}
		tool2 := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{Name: apiutils.Of("tool-2")},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-write")},
			},
		}
		tool3 := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{Name: apiutils.Of("tool-3")},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("write-only")},
			},
		}

		_, err = toolRepo.Save(tool1, &parentID)
		require.NoError(t, err)
		_, err = toolRepo.Save(tool2, &parentID)
		require.NoError(t, err)
		_, err = toolRepo.Save(tool3, &parentID)
		require.NoError(t, err)

		// List all tools for parent
		tools, err := toolRepo.List(models.MCPServerToolListOptions{ParentID: parentID})
		require.NoError(t, err)
		assert.Len(t, tools, 3)

		// Verify tool names
		toolNames := make([]string, len(tools))
		for i, tool := range tools {
			toolNames[i] = *tool.GetAttributes().Name
		}
		assert.Contains(t, toolNames, "tool-1")
		assert.Contains(t, toolNames, "tool-2")
		assert.Contains(t, toolNames, "tool-3")
	})

	t.Run("TestUpdateTool", func(t *testing.T) {
		// Create parent server
		parentServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{Name: apiutils.Of("parent-update")},
			Properties: &[]dbmodels.Properties{
				{Name: "source_id", StringValue: apiutils.Of("source-update")},
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}
		savedServer, err := serverRepo.Save(parentServer)
		require.NoError(t, err)
		parentID := *savedServer.GetID()

		// Create tool
		tool := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{Name: apiutils.Of("update-tool")},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-only")},
				{Name: "description", StringValue: apiutils.Of("Original description")},
			},
		}
		saved, err := toolRepo.Save(tool, &parentID)
		require.NoError(t, err)

		// Update the tool
		updateTool := &models.MCPServerToolImpl{
			ID:         saved.GetID(),
			Attributes: &models.MCPServerToolAttributes{Name: apiutils.Of("update-tool")},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-write")},
				{Name: "description", StringValue: apiutils.Of("Updated description")},
			},
		}
		updated, err := toolRepo.Save(updateTool, &parentID)
		require.NoError(t, err)
		assert.Equal(t, *saved.GetID(), *updated.GetID())

		// Verify updates
		props := updated.GetProperties()
		require.NotNil(t, props)
		accessTypeFound := false
		descriptionFound := false
		for _, prop := range *props {
			if prop.Name == "accessType" && prop.StringValue != nil && *prop.StringValue == "read-write" {
				accessTypeFound = true
			}
			if prop.Name == "description" && prop.StringValue != nil && *prop.StringValue == "Updated description" {
				descriptionFound = true
			}
		}
		assert.True(t, accessTypeFound, "accessType should be updated to read-write")
		assert.True(t, descriptionFound, "description should be updated")
	})

	t.Run("TestDeleteToolByID", func(t *testing.T) {
		// Create parent server
		parentServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{Name: apiutils.Of("parent-delete")},
			Properties: &[]dbmodels.Properties{
				{Name: "source_id", StringValue: apiutils.Of("source-delete")},
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}
		savedServer, err := serverRepo.Save(parentServer)
		require.NoError(t, err)
		parentID := *savedServer.GetID()

		// Create tool
		tool := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{Name: apiutils.Of("delete-tool")},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-only")},
			},
		}
		saved, err := toolRepo.Save(tool, &parentID)
		require.NoError(t, err)
		toolID := *saved.GetID()

		// Delete the tool
		err = toolRepo.DeleteByID(toolID)
		require.NoError(t, err)

		// Verify deletion
		_, err = toolRepo.GetByID(toolID)
		assert.ErrorIs(t, err, ErrMCPServerToolNotFound)

		// Attempt to delete non-existent tool
		err = toolRepo.DeleteByID(99999)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrMCPServerToolNotFound)
	})

	t.Run("TestDeleteToolsByParentID", func(t *testing.T) {
		// Create parent server
		parentServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{Name: apiutils.Of("parent-to-delete-all")},
			Properties: &[]dbmodels.Properties{
				{Name: "source_id", StringValue: apiutils.Of("source-delete-all")},
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}
		savedServer, err := serverRepo.Save(parentServer)
		require.NoError(t, err)
		parentID := *savedServer.GetID()

		// Create multiple tools
		for i := 1; i <= 5; i++ {
			tool := &models.MCPServerToolImpl{
				Attributes: &models.MCPServerToolAttributes{
					Name: apiutils.Of(fmt.Sprintf("delete-all-tool-%d", i)),
				},
				Properties: &[]dbmodels.Properties{
					{Name: "accessType", StringValue: apiutils.Of("read-only")},
				},
			}
			_, err := toolRepo.Save(tool, &parentID)
			require.NoError(t, err)
		}

		// Verify 5 tools exist
		tools, err := toolRepo.List(models.MCPServerToolListOptions{ParentID: parentID})
		require.NoError(t, err)
		assert.Len(t, tools, 5)

		// Delete all tools by parent ID
		err = toolRepo.DeleteByParentID(parentID)
		require.NoError(t, err)

		// Verify all tools are deleted
		tools, err = toolRepo.List(models.MCPServerToolListOptions{ParentID: parentID})
		require.NoError(t, err)
		assert.Empty(t, tools)
	})

	t.Run("TestCascadeDeleteToolsWithParent", func(t *testing.T) {
		// Create parent server
		parentServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{Name: apiutils.Of("cascade-parent")},
			Properties: &[]dbmodels.Properties{
				{Name: "source_id", StringValue: apiutils.Of("source-cascade")},
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}
		savedServer, err := serverRepo.Save(parentServer)
		require.NoError(t, err)
		parentID := *savedServer.GetID()

		// Create tools
		tool1 := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{Name: apiutils.Of("cascade-tool-1")},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-only")},
			},
		}
		tool2 := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{Name: apiutils.Of("cascade-tool-2")},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-write")},
			},
		}

		savedTool1, err := toolRepo.Save(tool1, &parentID)
		require.NoError(t, err)
		savedTool2, err := toolRepo.Save(tool2, &parentID)
		require.NoError(t, err)

		// Verify tools exist
		tools, err := toolRepo.List(models.MCPServerToolListOptions{ParentID: parentID})
		require.NoError(t, err)
		assert.Len(t, tools, 2)

		// Delete tools first, then delete parent server
		// Note: ML Metadata schema doesn't cascade delete Executions when Context is deleted
		err = toolRepo.DeleteByParentID(parentID)
		require.NoError(t, err)

		// Verify tools are deleted
		_, err = toolRepo.GetByID(*savedTool1.GetID())
		assert.ErrorIs(t, err, ErrMCPServerToolNotFound)
		_, err = toolRepo.GetByID(*savedTool2.GetID())
		assert.ErrorIs(t, err, ErrMCPServerToolNotFound)

		// Now delete parent server
		err = serverRepo.DeleteByID(parentID)
		require.NoError(t, err)
	})

	t.Run("TestToolWithCustomProperties", func(t *testing.T) {
		// Create parent server
		parentServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{Name: apiutils.Of("parent-custom")},
			Properties: &[]dbmodels.Properties{
				{Name: "source_id", StringValue: apiutils.Of("source-custom")},
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}
		savedServer, err := serverRepo.Save(parentServer)
		require.NoError(t, err)
		parentID := *savedServer.GetID()

		// Create tool with custom properties
		tool := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{Name: apiutils.Of("custom-tool")},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-only")},
			},
			CustomProperties: &[]dbmodels.Properties{
				{Name: "customField1", StringValue: apiutils.Of("custom value 1")},
				{Name: "customField2", IntValue: apiutils.Of(int32(42))},
			},
		}

		saved, err := toolRepo.Save(tool, &parentID)
		require.NoError(t, err)

		// Retrieve and verify custom properties
		retrieved, err := toolRepo.GetByID(*saved.GetID())
		require.NoError(t, err)

		customProps := retrieved.GetCustomProperties()
		require.NotNil(t, customProps)
		assert.Len(t, *customProps, 2)

		// Verify custom property values
		customMap := make(map[string]dbmodels.Properties)
		for _, prop := range *customProps {
			customMap[prop.Name] = prop
		}

		assert.Equal(t, "custom value 1", *customMap["customField1"].StringValue)
		assert.Equal(t, int32(42), *customMap["customField2"].IntValue)
	})

	// ==============================================================================
	// Edge Case and Validation Tests
	// ==============================================================================

	t.Run("TestSaveToolWithNilParentID", func(t *testing.T) {
		// Attempt to save a tool without a parent ID
		tool := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{Name: apiutils.Of("orphan-tool")},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-only")},
			},
		}

		// Save should succeed - GenericRepository handles nil parentID
		saved, err := toolRepo.Save(tool, nil)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Verify the tool was saved
		retrieved, err := toolRepo.GetByID(*saved.GetID())
		require.NoError(t, err)
		assert.Equal(t, "orphan-tool", *retrieved.GetAttributes().Name)
	})

	t.Run("TestSaveToolWithInvalidParentID", func(t *testing.T) {
		// Attempt to save a tool with a non-existent parent ID
		invalidParentID := int32(999999)
		tool := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{Name: apiutils.Of("invalid-parent-tool")},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-only")},
			},
		}

		// Save should fail with foreign key constraint error
		_, err := toolRepo.Save(tool, &invalidParentID)
		assert.Error(t, err)
	})

	t.Run("TestSaveToolWithNilName", func(t *testing.T) {
		// Create parent server
		parentServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{Name: apiutils.Of("parent-nil-name")},
			Properties: &[]dbmodels.Properties{
				{Name: "source_id", StringValue: apiutils.Of("source-nil-name")},
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}
		savedServer, err := serverRepo.Save(parentServer)
		require.NoError(t, err)
		parentID := *savedServer.GetID()

		// Attempt to save a tool with nil name
		tool := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{Name: nil},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-only")},
			},
		}

		_, err = toolRepo.Save(tool, &parentID)
		assert.ErrorIs(t, err, ErrMCPServerToolNameEmpty)
	})

	t.Run("TestSaveToolWithEmptyName", func(t *testing.T) {
		// Create parent server
		parentServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{Name: apiutils.Of("parent-empty-name")},
			Properties: &[]dbmodels.Properties{
				{Name: "source_id", StringValue: apiutils.Of("source-empty-name")},
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}
		savedServer, err := serverRepo.Save(parentServer)
		require.NoError(t, err)
		parentID := *savedServer.GetID()

		// Attempt to save a tool with empty name
		emptyName := ""
		tool := &models.MCPServerToolImpl{
			Attributes: &models.MCPServerToolAttributes{Name: &emptyName},
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-only")},
			},
		}

		_, err = toolRepo.Save(tool, &parentID)
		assert.ErrorIs(t, err, ErrMCPServerToolNameEmpty)
	})

	t.Run("TestSaveToolWithNilAttributes", func(t *testing.T) {
		// Create parent server
		parentServer := &models.MCPServerImpl{
			Attributes: &models.MCPServerAttributes{Name: apiutils.Of("parent-nil-attrs")},
			Properties: &[]dbmodels.Properties{
				{Name: "source_id", StringValue: apiutils.Of("source-nil-attrs")},
				{Name: "version", StringValue: apiutils.Of("1.0.0")},
			},
		}
		savedServer, err := serverRepo.Save(parentServer)
		require.NoError(t, err)
		parentID := *savedServer.GetID()

		// Attempt to save a tool with nil attributes
		tool := &models.MCPServerToolImpl{
			Attributes: nil,
			Properties: &[]dbmodels.Properties{
				{Name: "accessType", StringValue: apiutils.Of("read-only")},
			},
		}

		_, err = toolRepo.Save(tool, &parentID)
		assert.ErrorIs(t, err, ErrMCPServerToolNameEmpty)
	})

	t.Run("TestListToolsForNonExistentParent", func(t *testing.T) {
		// List tools for a non-existent parent ID
		nonExistentParentID := int32(999999)
		tools, err := toolRepo.List(models.MCPServerToolListOptions{ParentID: nonExistentParentID})
		require.NoError(t, err)
		assert.Empty(t, tools, "Should return empty slice for non-existent parent")
	})

	t.Run("TestDeleteToolsByNonExistentParentID", func(t *testing.T) {
		// Delete tools by non-existent parent ID (should be idempotent)
		nonExistentParentID := int32(999999)
		err := toolRepo.DeleteByParentID(nonExistentParentID)
		require.NoError(t, err, "DeleteByParentID should be idempotent")
	})
}
