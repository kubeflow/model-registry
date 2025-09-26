package service_test

import (
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServingEnvironmentRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupMySQLWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	// Get the actual ServingEnvironment type ID from the database
	typeID := getServingEnvironmentTypeID(t, sharedDB)
	repo := service.NewServingEnvironmentRepository(sharedDB, typeID)

	t.Run("TestSave", func(t *testing.T) {
		// Test creating a new serving environment
		servingEnvironment := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name:       apiutils.Of("test-serving-env"),
				ExternalID: apiutils.Of("serving-ext-123"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Test serving environment description"),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:        "custom-serving-prop",
					StringValue: apiutils.Of("custom-serving-value"),
				},
			},
		}

		saved, err := repo.Save(servingEnvironment)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-serving-env", *saved.GetAttributes().Name)
		assert.Equal(t, "serving-ext-123", *saved.GetAttributes().ExternalID)

		// Test updating the same serving environment
		servingEnvironment.ID = saved.GetID()
		servingEnvironment.GetAttributes().Name = apiutils.Of("updated-serving-env")
		// Preserve CreateTimeSinceEpoch from the saved entity (simulating what OpenAPI converter would do)
		servingEnvironment.GetAttributes().CreateTimeSinceEpoch = saved.GetAttributes().CreateTimeSinceEpoch

		updated, err := repo.Save(servingEnvironment)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, "updated-serving-env", *updated.GetAttributes().Name)
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// First create a serving environment to retrieve
		servingEnvironment := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name:       apiutils.Of("get-test-serving-env"),
				ExternalID: apiutils.Of("get-serving-ext-123"),
			},
		}

		saved, err := repo.Save(servingEnvironment)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-test-serving-env", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-serving-ext-123", *retrieved.GetAttributes().ExternalID)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.Error(t, err)
	})

	t.Run("TestList", func(t *testing.T) {
		// Create multiple serving environments for listing
		testEnvironments := []*models.ServingEnvironmentImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ServingEnvironmentAttributes{
					Name:       apiutils.Of("list-serving-env-1"),
					ExternalID: apiutils.Of("list-serving-ext-1"),
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ServingEnvironmentAttributes{
					Name:       apiutils.Of("list-serving-env-2"),
					ExternalID: apiutils.Of("list-serving-ext-2"),
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ServingEnvironmentAttributes{
					Name:       apiutils.Of("list-serving-env-3"),
					ExternalID: apiutils.Of("list-serving-ext-3"),
				},
			},
		}

		for _, env := range testEnvironments {
			_, err := repo.Save(env)
			require.NoError(t, err)
		}

		// Test listing all environments with basic pagination
		pageSize := int32(10)
		listOptions := models.ServingEnvironmentListOptions{}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // At least our 3 test environments

		// Test listing by name
		listOptions = models.ServingEnvironmentListOptions{
			Name: apiutils.Of("list-serving-env-1"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-serving-env-1", *result.Items[0].GetAttributes().Name)
		}

		// Test listing by external ID
		listOptions = models.ServingEnvironmentListOptions{
			ExternalID: apiutils.Of("list-serving-ext-2"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-serving-ext-2", *result.Items[0].GetAttributes().ExternalID)
		}

		// Test ordering by ID (deterministic)
		listOptions = models.ServingEnvironmentListOptions{
			Pagination: models.Pagination{
				OrderBy: apiutils.Of("ID"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		// Verify we get results back and they are ordered by ID
		assert.GreaterOrEqual(t, len(result.Items), 1)
		if len(result.Items) > 1 {
			// Verify ascending ID order
			firstID := *result.Items[0].GetID()
			secondID := *result.Items[1].GetID()
			assert.Less(t, firstID, secondID, "Results should be ordered by ID ascending")
		}
	})

	t.Run("TestListOrdering", func(t *testing.T) {
		// Create environments sequentially with time delays to ensure deterministic ordering
		env1 := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("time-test-serving-env-1"),
			},
		}
		saved1, err := repo.Save(env1)
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		env2 := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("time-test-serving-env-2"),
			},
		}
		saved2, err := repo.Save(env2)
		require.NoError(t, err)

		// Test ordering by CREATE_TIME
		pageSize := int32(100) // Increased page size to ensure all test entities are included
		listOptions := models.ServingEnvironmentListOptions{
			Pagination: models.Pagination{
				OrderBy: apiutils.Of("CREATE_TIME"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test environments in the results
		var foundEnv1, foundEnv2 models.ServingEnvironment
		var index1, index2 = -1, -1

		for i, item := range result.Items {
			if *item.GetID() == *saved1.GetID() {
				foundEnv1 = item
				index1 = i
			}
			if *item.GetID() == *saved2.GetID() {
				foundEnv2 = item
				index2 = i
			}
		}

		// Verify both environments were found and env1 comes before env2 (ascending order)
		require.NotEqual(t, -1, index1, "Environment 1 should be found in results")
		require.NotEqual(t, -1, index2, "Environment 2 should be found in results")
		assert.Less(t, index1, index2, "Environment 1 should come before Environment 2 when ordered by CREATE_TIME")
		assert.Less(t, *foundEnv1.GetAttributes().CreateTimeSinceEpoch, *foundEnv2.GetAttributes().CreateTimeSinceEpoch, "Environment 1 should have earlier create time")
	})

	t.Run("TestSaveWithProperties", func(t *testing.T) {
		servingEnvironment := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("props-test-serving-env"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Environment with properties"),
				},
				{
					Name:     "capacity",
					IntValue: apiutils.Of(int32(100)),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:        "team",
					StringValue: apiutils.Of("ml-team"),
				},
				{
					Name:     "priority",
					IntValue: apiutils.Of(int32(5)),
				},
			},
		}

		saved, err := repo.Save(servingEnvironment)
		require.NoError(t, err)
		require.NotNil(t, saved)

		// Verify properties were saved
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		assert.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 2)

		assert.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 2)
	})
}
