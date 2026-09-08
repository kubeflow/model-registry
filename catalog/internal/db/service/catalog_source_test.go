package service_test

import (
	"errors"
	"testing"

	"github.com/kubeflow/hub/catalog/internal/db/models"
	"github.com/kubeflow/hub/catalog/internal/db/service"
	"github.com/kubeflow/hub/catalog/internal/testhelpers"
	dbmodels "github.com/kubeflow/hub/internal/platform/db/entity"
	"github.com/kubeflow/hub/internal/platform/db/schema"
	"github.com/kubeflow/hub/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCatalogSourceRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupPostgresWithMigrations(t, testhelpers.MustDatastoreSpec(t))
	defer cleanup()

	// Get the CatalogSource type ID
	typeID := getCatalogSourceTypeID(t, sharedDB)
	repo := service.NewCatalogSourceRepository(sharedDB, typeID)

	t.Run("TestSave_Create", func(t *testing.T) {
		// Test creating a new catalog source
		source := &models.CatalogSourceImpl{
			Attributes: &models.CatalogSourceAttributes{
				Name: new("test-source-create"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "status",
					StringValue: new("available"),
				},
			},
		}

		saved, err := repo.Save(source)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-source-create", *saved.GetAttributes().Name)

		// Verify timestamps were set
		attrs := saved.GetAttributes()
		require.NotNil(t, attrs.CreateTimeSinceEpoch)
		require.NotNil(t, attrs.LastUpdateTimeSinceEpoch)
		assert.Greater(t, *attrs.CreateTimeSinceEpoch, int64(0))
	})

	t.Run("TestSave_Update", func(t *testing.T) {
		// First create a source
		source := &models.CatalogSourceImpl{
			Attributes: &models.CatalogSourceAttributes{
				Name: new("test-source-update"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "status",
					StringValue: new("available"),
				},
			},
		}

		saved, err := repo.Save(source)
		require.NoError(t, err)
		originalCreateTime := *saved.GetAttributes().CreateTimeSinceEpoch

		// Update the same source (by name)
		updatedSource := &models.CatalogSourceImpl{
			Attributes: &models.CatalogSourceAttributes{
				Name: new("test-source-update"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "status",
					StringValue: new("error"),
				},
				{
					Name:        "error",
					StringValue: new("connection failed"),
				},
			},
		}

		updated, err := repo.Save(updatedSource)
		require.NoError(t, err)
		require.NotNil(t, updated)

		// Verify ID is the same (update, not create)
		assert.Equal(t, *saved.GetID(), *updated.GetID())

		// Verify CreateTimeSinceEpoch is preserved
		assert.Equal(t, originalCreateTime, *updated.GetAttributes().CreateTimeSinceEpoch)

		// Verify properties were updated
		require.NotNil(t, updated.GetProperties())
		assert.Len(t, *updated.GetProperties(), 2)
	})

	t.Run("TestSave_MissingName", func(t *testing.T) {
		// Test saving without a name should fail
		source := &models.CatalogSourceImpl{
			Attributes: &models.CatalogSourceAttributes{
				// No Name set
			},
		}

		_, err := repo.Save(source)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "source ID (name) is required")
	})

	t.Run("TestGetBySourceID_Found", func(t *testing.T) {
		// First create a source
		source := &models.CatalogSourceImpl{
			Attributes: &models.CatalogSourceAttributes{
				Name: new("test-source-get"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "status",
					StringValue: new("available"),
				},
			},
		}

		saved, err := repo.Save(source)
		require.NoError(t, err)

		// Retrieve by source ID
		retrieved, err := repo.GetBySourceID("test-source-get")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "test-source-get", *retrieved.GetAttributes().Name)

		// Verify properties were retrieved
		require.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 1)
	})

	t.Run("TestGetBySourceID_NotFound", func(t *testing.T) {
		_, err := repo.GetBySourceID("non-existent-source")
		require.Error(t, err)
		assert.True(t, errors.Is(err, service.ErrCatalogSourceNotFound))
	})

	t.Run("TestDelete_Success", func(t *testing.T) {
		// First create a source
		source := &models.CatalogSourceImpl{
			Attributes: &models.CatalogSourceAttributes{
				Name: new("test-source-delete"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "status",
					StringValue: new("available"),
				},
			},
		}

		_, err := repo.Save(source)
		require.NoError(t, err)

		// Verify it exists
		_, err = repo.GetBySourceID("test-source-delete")
		require.NoError(t, err)

		// Delete it
		err = repo.Delete("test-source-delete")
		require.NoError(t, err)

		// Verify it's gone
		_, err = repo.GetBySourceID("test-source-delete")
		require.Error(t, err)
		assert.True(t, errors.Is(err, service.ErrCatalogSourceNotFound))
	})

	t.Run("TestDelete_NonExistent", func(t *testing.T) {
		// Deleting a non-existent source should not error
		err := repo.Delete("non-existent-source-to-delete")
		require.NoError(t, err)
	})

	t.Run("TestGetAll_Empty", func(t *testing.T) {
		// Clear all existing sources first
		existingSources, err := repo.GetAll()
		require.NoError(t, err)
		for _, s := range existingSources {
			if attrs := s.GetAttributes(); attrs != nil && attrs.Name != nil {
				err := repo.Delete(*attrs.Name)
				require.NoError(t, err)
			}
		}

		// Now test empty case
		sources, err := repo.GetAll()
		require.NoError(t, err)
		assert.Empty(t, sources)
	})

	t.Run("TestGetAll_Populated", func(t *testing.T) {
		// Create multiple sources
		sources := []string{"getall-source-1", "getall-source-2", "getall-source-3"}
		for _, name := range sources {
			source := &models.CatalogSourceImpl{
				Attributes: &models.CatalogSourceAttributes{
					Name: new(name),
				},
				Properties: &[]dbmodels.Properties{
					{
						Name:        "status",
						StringValue: new("available"),
					},
				},
			}
			_, err := repo.Save(source)
			require.NoError(t, err)
		}

		// Get all sources
		allSources, err := repo.GetAll()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(allSources), 3)

		// Verify our sources are in the result
		foundNames := make(map[string]bool)
		for _, s := range allSources {
			if attrs := s.GetAttributes(); attrs != nil && attrs.Name != nil {
				foundNames[*attrs.Name] = true
			}
		}

		for _, name := range sources {
			assert.True(t, foundNames[name], "Source %s should be in GetAll results", name)
		}
	})

	t.Run("TestGetAllStatuses_Empty", func(t *testing.T) {
		// Clear all existing sources first
		existingSources, err := repo.GetAll()
		require.NoError(t, err)
		for _, s := range existingSources {
			if attrs := s.GetAttributes(); attrs != nil && attrs.Name != nil {
				err := repo.Delete(*attrs.Name)
				require.NoError(t, err)
			}
		}

		// Now test empty case
		statuses, err := repo.GetAllStatuses()
		require.NoError(t, err)
		assert.Empty(t, statuses)
	})

	t.Run("TestGetAllStatuses_WithStatusAndError", func(t *testing.T) {
		// Create sources with different statuses
		testCases := []struct {
			name     string
			status   string
			errorMsg string
			hasError bool
		}{
			{"status-available", "available", "", false},
			{"status-partially-available", "partially-available", "some models failed to load", true},
			{"status-error", "error", "connection timeout", true},
			{"status-disabled", "disabled", "", false},
		}

		for _, tc := range testCases {
			props := []dbmodels.Properties{
				{
					Name:        "status",
					StringValue: new(tc.status),
				},
			}
			if tc.hasError {
				props = append(props, dbmodels.Properties{
					Name:        "error",
					StringValue: new(tc.errorMsg),
				})
			}

			source := &models.CatalogSourceImpl{
				Attributes: &models.CatalogSourceAttributes{
					Name: new(tc.name),
				},
				Properties: &props,
			}
			_, err := repo.Save(source)
			require.NoError(t, err)
		}

		// Get all statuses
		statuses, err := repo.GetAllStatuses()
		require.NoError(t, err)

		// Verify each status
		for _, tc := range testCases {
			status, ok := statuses[tc.name]
			require.True(t, ok, "Status for %s should exist", tc.name)
			assert.Equal(t, tc.status, status.Status, "Status mismatch for %s", tc.name)
			assert.Equal(t, tc.errorMsg, status.Error, "Error mismatch for %s", tc.name)
		}
	})

	t.Run("TestGetAllStatuses_StatusExtraction", func(t *testing.T) {
		// Test that status is correctly extracted from properties
		source := &models.CatalogSourceImpl{
			Attributes: &models.CatalogSourceAttributes{
				Name: new("status-extraction-test"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "status",
					StringValue: new("error"),
				},
				{
					Name:        "error",
					StringValue: new("detailed error message here"),
				},
				{
					Name:        "other_prop",
					StringValue: new("should be ignored"),
				},
			},
		}

		_, err := repo.Save(source)
		require.NoError(t, err)

		statuses, err := repo.GetAllStatuses()
		require.NoError(t, err)

		status, ok := statuses["status-extraction-test"]
		require.True(t, ok)
		assert.Equal(t, "error", status.Status)
		assert.Equal(t, "detailed error message here", status.Error)
	})

	t.Run("TestGetAllStatuses_WithCredentials", func(t *testing.T) {
		// Test that has_api_key and authenticated are correctly extracted
		trueVal := true
		falseVal := false

		hfUser := "jdoe"

		// Source with API key, successful auth, and HF username
		source1 := &models.CatalogSourceImpl{
			Attributes: &models.CatalogSourceAttributes{
				Name: new("cred-source-authed"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "status", StringValue: new("available")},
				{Name: "has_api_key", BoolValue: &trueVal},
				{Name: "authenticated", BoolValue: &trueVal},
				{Name: "hf_username", StringValue: &hfUser},
			},
		}
		_, err := repo.Save(source1)
		require.NoError(t, err)

		// Source with API key and failed auth (no username)
		source2 := &models.CatalogSourceImpl{
			Attributes: &models.CatalogSourceAttributes{
				Name: new("cred-source-unauthed"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "status", StringValue: new("error")},
				{Name: "error", StringValue: new("invalid credentials")},
				{Name: "has_api_key", BoolValue: &trueVal},
				{Name: "authenticated", BoolValue: &falseVal},
			},
		}
		_, err = repo.Save(source2)
		require.NoError(t, err)

		// Source without API key
		source3 := &models.CatalogSourceImpl{
			Attributes: &models.CatalogSourceAttributes{
				Name: new("cred-source-nokey"),
			},
			Properties: &[]dbmodels.Properties{
				{Name: "status", StringValue: new("available")},
				{Name: "has_api_key", BoolValue: &falseVal},
			},
		}
		_, err = repo.Save(source3)
		require.NoError(t, err)

		statuses, err := repo.GetAllStatuses()
		require.NoError(t, err)

		// Verify source with valid key and username
		s1, ok := statuses["cred-source-authed"]
		require.True(t, ok)
		require.NotNil(t, s1.HasApiKey)
		assert.True(t, *s1.HasApiKey)
		require.NotNil(t, s1.Authenticated)
		assert.True(t, *s1.Authenticated)
		require.NotNil(t, s1.HfUsername)
		assert.Equal(t, "jdoe", *s1.HfUsername)

		// Verify source with invalid key (no username)
		s2, ok := statuses["cred-source-unauthed"]
		require.True(t, ok)
		require.NotNil(t, s2.HasApiKey)
		assert.True(t, *s2.HasApiKey)
		require.NotNil(t, s2.Authenticated)
		assert.False(t, *s2.Authenticated)
		assert.Nil(t, s2.HfUsername)

		// Verify source without key
		s3, ok := statuses["cred-source-nokey"]
		require.True(t, ok)
		require.NotNil(t, s3.HasApiKey)
		assert.False(t, *s3.HasApiKey)
		assert.Nil(t, s3.Authenticated)
		assert.Nil(t, s3.HfUsername)
	})

	t.Run("TestPropertiesPreserved", func(t *testing.T) {
		// Test that all properties are correctly saved and retrieved
		source := &models.CatalogSourceImpl{
			Attributes: &models.CatalogSourceAttributes{
				Name: new("props-test-source"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "status",
					StringValue: new("available"),
				},
				{
					Name:     "count",
					IntValue: new(int32(42)),
				},
				{
					Name:      "enabled",
					BoolValue: new(true),
				},
			},
		}

		_, err := repo.Save(source)
		require.NoError(t, err)

		retrieved, err := repo.GetBySourceID("props-test-source")
		require.NoError(t, err)

		props := retrieved.GetProperties()
		require.NotNil(t, props)
		assert.Len(t, *props, 3)

		// Verify properties by name
		propMap := make(map[string]dbmodels.Properties)
		for _, p := range *props {
			propMap[p.Name] = p
		}

		assert.Equal(t, "available", *propMap["status"].StringValue)
		assert.Equal(t, int32(42), *propMap["count"].IntValue)
		assert.Equal(t, true, *propMap["enabled"].BoolValue)
	})
}

// Helper function to get or create CatalogSource type ID
func getCatalogSourceTypeID(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogSourceTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogSource type")
	}

	return typeRecord.ID
}
