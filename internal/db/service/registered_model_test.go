package service_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	cont_mysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	mysqlContainer, err := cont_mysql.Run(
		ctx,
		"mysql:5.7",
		cont_mysql.WithUsername("root"),
		cont_mysql.WithPassword("root"),
		cont_mysql.WithDatabase("test"),
		cont_mysql.WithConfigFile(filepath.Join("testdata", "testdb.cnf")),
	)
	require.NoError(t, err)

	dbConnector := mysql.NewMySQLDBConnector(mysqlContainer.MustConnectionString(ctx))
	require.NoError(t, err)

	db, err := dbConnector.Connect()
	require.NoError(t, err)

	// Run migrations
	migrator, err := mysql.NewMySQLMigrator(db)
	require.NoError(t, err)
	err = migrator.Migrate()
	require.NoError(t, err)

	// Return cleanup function
	cleanup := func() {
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
		err = testcontainers.TerminateContainer(mysqlContainer)
		require.NoError(t, err)
	}

	return db, cleanup
}

func getRegisteredModelTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.RegisteredModelTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find RegisteredModel type")
	return int64(typeRecord.ID)
}

func TestRegisteredModelRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Get the actual RegisteredModel type ID from the database
	typeID := getRegisteredModelTypeID(t, db)
	repo := service.NewRegisteredModelRepository(db, typeID)

	t.Run("TestSave", func(t *testing.T) {
		// Test creating a new registered model
		registeredModel := &models.RegisteredModelImpl{
			TypeID: int32Ptr(int32(typeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name:       stringPtr("test-model"),
				ExternalID: stringPtr("ext-123"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: stringPtr("Test model description"),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:        "custom-prop",
					StringValue: stringPtr("custom-value"),
				},
			},
		}

		saved, err := repo.Save(registeredModel)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-model", *saved.GetAttributes().Name)
		assert.Equal(t, "ext-123", *saved.GetAttributes().ExternalID)

		// Test updating the same model
		registeredModel.ID = saved.GetID()
		registeredModel.GetAttributes().Name = stringPtr("updated-model")

		updated, err := repo.Save(registeredModel)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, "updated-model", *updated.GetAttributes().Name)
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// First create a model to retrieve
		registeredModel := &models.RegisteredModelImpl{
			TypeID: int32Ptr(int32(typeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name:       stringPtr("get-test-model"),
				ExternalID: stringPtr("get-ext-123"),
			},
		}

		saved, err := repo.Save(registeredModel)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-test-model", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-ext-123", *retrieved.GetAttributes().ExternalID)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.Error(t, err)
	})

	t.Run("TestList", func(t *testing.T) {
		// Create multiple models for listing
		testModels := []*models.RegisteredModelImpl{
			{
				TypeID: int32Ptr(int32(typeID)),
				Attributes: &models.RegisteredModelAttributes{
					Name:       stringPtr("list-model-1"),
					ExternalID: stringPtr("list-ext-1"),
				},
			},
			{
				TypeID: int32Ptr(int32(typeID)),
				Attributes: &models.RegisteredModelAttributes{
					Name:       stringPtr("list-model-2"),
					ExternalID: stringPtr("list-ext-2"),
				},
			},
			{
				TypeID: int32Ptr(int32(typeID)),
				Attributes: &models.RegisteredModelAttributes{
					Name:       stringPtr("list-model-3"),
					ExternalID: stringPtr("list-ext-3"),
				},
			},
		}

		for _, model := range testModels {
			_, err := repo.Save(model)
			require.NoError(t, err)
		}

		// Test listing all models with basic pagination
		pageSize := int32(10)
		listOptions := models.RegisteredModelListOptions{}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // At least our 3 test models

		// Test listing by name
		listOptions = models.RegisteredModelListOptions{
			Name: stringPtr("list-model-1"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-model-1", *result.Items[0].GetAttributes().Name)
		}

		// Test listing by external ID
		listOptions = models.RegisteredModelListOptions{
			ExternalID: stringPtr("list-ext-2"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-ext-2", *result.Items[0].GetAttributes().ExternalID)
		}

		// Test ordering by ID (deterministic)
		listOptions = models.RegisteredModelListOptions{
			Pagination: models.Pagination{
				OrderBy: stringPtr("ID"),
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
		// Create models sequentially with time delays to ensure deterministic ordering
		model1 := &models.RegisteredModelImpl{
			TypeID: int32Ptr(int32(typeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: stringPtr("time-test-model-1"),
			},
		}
		saved1, err := repo.Save(model1)
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		model2 := &models.RegisteredModelImpl{
			TypeID: int32Ptr(int32(typeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: stringPtr("time-test-model-2"),
			},
		}
		saved2, err := repo.Save(model2)
		require.NoError(t, err)

		// Test ordering by CREATE_TIME
		pageSize := int32(10)
		listOptions := models.RegisteredModelListOptions{
			Pagination: models.Pagination{
				OrderBy: stringPtr("CREATE_TIME"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test models in the results
		var foundModel1, foundModel2 models.RegisteredModel
		var index1, index2 int = -1, -1

		for i, item := range result.Items {
			if *item.GetID() == *saved1.GetID() {
				foundModel1 = item
				index1 = i
			}
			if *item.GetID() == *saved2.GetID() {
				foundModel2 = item
				index2 = i
			}
		}

		// Verify both models were found and model1 comes before model2 (ascending order)
		require.NotEqual(t, -1, index1, "Model 1 should be found in results")
		require.NotEqual(t, -1, index2, "Model 2 should be found in results")
		assert.Less(t, index1, index2, "Model 1 should come before Model 2 when ordered by CREATE_TIME")
		assert.Less(t, *foundModel1.GetAttributes().CreateTimeSinceEpoch, *foundModel2.GetAttributes().CreateTimeSinceEpoch, "Model 1 should have earlier create time")
	})

	t.Run("TestSaveWithProperties", func(t *testing.T) {
		registeredModel := &models.RegisteredModelImpl{
			TypeID: int32Ptr(int32(typeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: stringPtr("props-test-model"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: stringPtr("Model with properties"),
				},
				{
					Name:     "version",
					IntValue: int32Ptr(1),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:        "team",
					StringValue: stringPtr("ml-team"),
				},
				{
					Name:     "priority",
					IntValue: int32Ptr(5),
				},
			},
		}

		saved, err := repo.Save(registeredModel)
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

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}
