package service_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExperimentRunRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Get the actual ExperimentRun type ID from the database
	typeID := getExperimentRunTypeID(t, db)
	repo := service.NewExperimentRunRepository(db, typeID)

	// Also get experiment type for creating parent experiments
	experimentTypeID := getExperimentTypeID(t, db)
	experimentRepo := service.NewExperimentRepository(db, experimentTypeID)

	t.Run("TestSave", func(t *testing.T) {
		// First create a parent experiment
		experiment := &models.ExperimentImpl{
			TypeID: apiutils.Of(int32(experimentTypeID)),
			Attributes: &models.ExperimentAttributes{
				Name:       apiutils.Of("parent-experiment"),
				ExternalID: apiutils.Of("parent-experiment-ext-123"),
			},
		}
		savedExperiment, err := experimentRepo.Save(experiment)
		require.NoError(t, err)

		// Test creating a new experiment run
		experimentRun := &models.ExperimentRunImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ExperimentRunAttributes{
				Name:       apiutils.Of(fmt.Sprintf("%d:test-experiment-run", *savedExperiment.GetID())),
				ExternalID: apiutils.Of("experiment-run-ext-123"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Test experiment run description"),
				},
				{
					Name:        "owner",
					StringValue: apiutils.Of("test-user"),
				},
				{
					Name:        "experiment_id",
					StringValue: apiutils.Of(fmt.Sprintf("%d", *savedExperiment.GetID())),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "custom-run-prop",
					StringValue:      apiutils.Of("custom-run-value"),
					IsCustomProperty: true,
				},
			},
		}

		saved, err := repo.Save(experimentRun, savedExperiment.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, fmt.Sprintf("%d:test-experiment-run", *savedExperiment.GetID()), *saved.GetAttributes().Name)
		assert.Equal(t, "experiment-run-ext-123", *saved.GetAttributes().ExternalID)
		assert.NotNil(t, saved.GetAttributes().CreateTimeSinceEpoch)
		assert.NotNil(t, saved.GetAttributes().LastUpdateTimeSinceEpoch)

		// Test updating the same experiment run
		experimentRun.ID = saved.GetID()
		experimentRun.GetAttributes().Name = apiutils.Of(fmt.Sprintf("%d:updated-experiment-run", *savedExperiment.GetID()))
		experimentRun.GetAttributes().ExternalID = apiutils.Of("updated-experiment-run-ext-123")
		// Preserve CreateTimeSinceEpoch from the saved entity (simulating what OpenAPI converter would do)
		experimentRun.GetAttributes().CreateTimeSinceEpoch = saved.GetAttributes().CreateTimeSinceEpoch

		updated, err := repo.Save(experimentRun, savedExperiment.GetID())
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, fmt.Sprintf("%d:updated-experiment-run", *savedExperiment.GetID()), *updated.GetAttributes().Name)
		assert.Equal(t, "updated-experiment-run-ext-123", *updated.GetAttributes().ExternalID)
		assert.Equal(t, *saved.GetAttributes().CreateTimeSinceEpoch, *updated.GetAttributes().CreateTimeSinceEpoch)
		assert.Greater(t, *updated.GetAttributes().LastUpdateTimeSinceEpoch, *saved.GetAttributes().LastUpdateTimeSinceEpoch)
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// First create a parent experiment
		experiment := &models.ExperimentImpl{
			TypeID: apiutils.Of(int32(experimentTypeID)),
			Attributes: &models.ExperimentAttributes{
				Name:       apiutils.Of("parent-experiment-get"),
				ExternalID: apiutils.Of("parent-experiment-get-ext-123"),
			},
		}
		savedExperiment, err := experimentRepo.Save(experiment)
		require.NoError(t, err)

		// Create an experiment run to retrieve
		experimentRun := &models.ExperimentRunImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ExperimentRunAttributes{
				Name:       apiutils.Of(fmt.Sprintf("%d:get-test-experiment-run", *savedExperiment.GetID())),
				ExternalID: apiutils.Of("get-experiment-run-ext-123"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Experiment run for get test"),
				},
				{
					Name:        "experiment_id",
					StringValue: apiutils.Of(fmt.Sprintf("%d", *savedExperiment.GetID())),
				},
			},
		}

		saved, err := repo.Save(experimentRun, savedExperiment.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, fmt.Sprintf("%d:get-test-experiment-run", *savedExperiment.GetID()), *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-experiment-run-ext-123", *retrieved.GetAttributes().ExternalID)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "experiment run by id not found")
	})

	t.Run("TestList", func(t *testing.T) {
		// First create a parent experiment
		experiment := &models.ExperimentImpl{
			TypeID: apiutils.Of(int32(experimentTypeID)),
			Attributes: &models.ExperimentAttributes{
				Name:       apiutils.Of("parent-experiment-list"),
				ExternalID: apiutils.Of("parent-experiment-list-ext-123"),
			},
		}
		savedExperiment, err := experimentRepo.Save(experiment)
		require.NoError(t, err)

		// Create multiple experiment runs for listing
		testExperimentRuns := []*models.ExperimentRunImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ExperimentRunAttributes{
					Name:       apiutils.Of(fmt.Sprintf("%d:list-experiment-run-1", *savedExperiment.GetID())),
					ExternalID: apiutils.Of("list-experiment-run-ext-1"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("First experiment run"),
					},
					{
						Name:        "experiment_id",
						StringValue: apiutils.Of(fmt.Sprintf("%d", *savedExperiment.GetID())),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ExperimentRunAttributes{
					Name:       apiutils.Of(fmt.Sprintf("%d:list-experiment-run-2", *savedExperiment.GetID())),
					ExternalID: apiutils.Of("list-experiment-run-ext-2"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("Second experiment run"),
					},
					{
						Name:        "experiment_id",
						StringValue: apiutils.Of(fmt.Sprintf("%d", *savedExperiment.GetID())),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ExperimentRunAttributes{
					Name:       apiutils.Of(fmt.Sprintf("%d:list-experiment-run-3", *savedExperiment.GetID())),
					ExternalID: apiutils.Of("list-experiment-run-ext-3"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("Third experiment run"),
					},
					{
						Name:        "experiment_id",
						StringValue: apiutils.Of(fmt.Sprintf("%d", *savedExperiment.GetID())),
					},
				},
			},
		}

		for _, experimentRun := range testExperimentRuns {
			_, err := repo.Save(experimentRun, savedExperiment.GetID())
			require.NoError(t, err)
		}

		// Test listing all experiment runs with basic pagination
		pageSize := int32(10)
		listOptions := models.ExperimentRunListOptions{}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // At least our 3 test experiment runs

		// Test listing by name
		listOptions = models.ExperimentRunListOptions{
			Name: apiutils.Of("list-experiment-run-1"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, fmt.Sprintf("%d:list-experiment-run-1", *savedExperiment.GetID()), *result.Items[0].GetAttributes().Name)
		}

		// Test listing by external ID
		listOptions = models.ExperimentRunListOptions{
			ExternalID: apiutils.Of("list-experiment-run-ext-2"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-experiment-run-ext-2", *result.Items[0].GetAttributes().ExternalID)
		}

		// Test listing by experiment ID
		listOptions = models.ExperimentRunListOptions{
			ExperimentID: savedExperiment.GetID(),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // Should find our 3 test experiment runs

		// Test ordering by ID (deterministic)
		listOptions = models.ExperimentRunListOptions{
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
		// First create a parent experiment
		experiment := &models.ExperimentImpl{
			TypeID: apiutils.Of(int32(experimentTypeID)),
			Attributes: &models.ExperimentAttributes{
				Name:       apiutils.Of("parent-experiment-ordering"),
				ExternalID: apiutils.Of("parent-experiment-ordering-ext-123"),
			},
		}
		savedExperiment, err := experimentRepo.Save(experiment)
		require.NoError(t, err)

		// Create experiment runs sequentially with time delays to ensure deterministic ordering
		experimentRun1 := &models.ExperimentRunImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ExperimentRunAttributes{
				Name:       apiutils.Of(fmt.Sprintf("%d:time-test-experiment-run-1", *savedExperiment.GetID())),
				ExternalID: apiutils.Of("time-experiment-run-ext-1"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "experiment_id",
					StringValue: apiutils.Of(fmt.Sprintf("%d", *savedExperiment.GetID())),
				},
			},
		}
		saved1, err := repo.Save(experimentRun1, savedExperiment.GetID())
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		experimentRun2 := &models.ExperimentRunImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ExperimentRunAttributes{
				Name:       apiutils.Of(fmt.Sprintf("%d:time-test-experiment-run-2", *savedExperiment.GetID())),
				ExternalID: apiutils.Of("time-experiment-run-ext-2"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "experiment_id",
					StringValue: apiutils.Of(fmt.Sprintf("%d", *savedExperiment.GetID())),
				},
			},
		}
		saved2, err := repo.Save(experimentRun2, savedExperiment.GetID())
		require.NoError(t, err)

		// Test ordering by CREATE_TIME
		pageSize := int32(10)
		listOptions := models.ExperimentRunListOptions{
			Pagination: models.Pagination{
				OrderBy: apiutils.Of("CREATE_TIME"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test experiment runs in the results
		var foundExperimentRun1, foundExperimentRun2 models.ExperimentRun
		var index1, index2 = -1, -1

		for i, item := range result.Items {
			if *item.GetID() == *saved1.GetID() {
				foundExperimentRun1 = item
				index1 = i
			}
			if *item.GetID() == *saved2.GetID() {
				foundExperimentRun2 = item
				index2 = i
			}
		}

		// Verify both experiment runs were found and experimentRun1 comes before experimentRun2 (ascending order)
		require.NotEqual(t, -1, index1, "Experiment Run 1 should be found in results")
		require.NotEqual(t, -1, index2, "Experiment Run 2 should be found in results")
		assert.Less(t, index1, index2, "Experiment Run 1 should come before Experiment Run 2 when ordered by CREATE_TIME")
		assert.Less(t, *foundExperimentRun1.GetAttributes().CreateTimeSinceEpoch, *foundExperimentRun2.GetAttributes().CreateTimeSinceEpoch, "Experiment Run 1 should have earlier create time")
	})

	t.Run("TestSaveWithProperties", func(t *testing.T) {
		// First create a parent experiment
		experiment := &models.ExperimentImpl{
			TypeID: apiutils.Of(int32(experimentTypeID)),
			Attributes: &models.ExperimentAttributes{
				Name:       apiutils.Of("parent-experiment-props"),
				ExternalID: apiutils.Of("parent-experiment-props-ext-123"),
			},
		}
		savedExperiment, err := experimentRepo.Save(experiment)
		require.NoError(t, err)

		experimentRun := &models.ExperimentRunImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ExperimentRunAttributes{
				Name:       apiutils.Of(fmt.Sprintf("%d:props-test-experiment-run", *savedExperiment.GetID())),
				ExternalID: apiutils.Of("props-experiment-run-ext-123"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Experiment run with properties"),
				},
				{
					Name:        "owner",
					StringValue: apiutils.Of("test-user"),
				},
				{
					Name:        "status",
					StringValue: apiutils.Of("RUNNING"),
				},
				{
					Name:        "experiment_id",
					StringValue: apiutils.Of(fmt.Sprintf("%d", *savedExperiment.GetID())),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "team",
					StringValue:      apiutils.Of("ml-team"),
					IsCustomProperty: true,
				},
				{
					Name:             "priority",
					IntValue:         apiutils.Of(int32(5)),
					IsCustomProperty: true,
				},
			},
		}

		saved, err := repo.Save(experimentRun, savedExperiment.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)

		// Verify properties were saved
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		assert.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 4) // description, owner, status, experiment_id

		assert.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 2) // team, priority

		// Verify specific properties exist
		foundDescription := false
		foundOwner := false
		foundStatus := false
		foundExperimentID := false
		for _, prop := range *retrieved.GetProperties() {
			switch prop.Name {
			case "description":
				foundDescription = true
				assert.Equal(t, "Experiment run with properties", *prop.StringValue)
			case "owner":
				foundOwner = true
				assert.Equal(t, "test-user", *prop.StringValue)
			case "status":
				foundStatus = true
				assert.Equal(t, "RUNNING", *prop.StringValue)
			case "experiment_id":
				foundExperimentID = true
				assert.Equal(t, fmt.Sprintf("%d", *savedExperiment.GetID()), *prop.StringValue)
			}
		}
		assert.True(t, foundDescription, "description property should exist")
		assert.True(t, foundOwner, "owner property should exist")
		assert.True(t, foundStatus, "status property should exist")
		assert.True(t, foundExperimentID, "experiment_id property should exist")

		// Verify custom properties
		foundTeam := false
		foundPriority := false
		for _, prop := range *retrieved.GetCustomProperties() {
			switch prop.Name {
			case "team":
				foundTeam = true
				assert.Equal(t, "ml-team", *prop.StringValue)
			case "priority":
				foundPriority = true
				assert.Equal(t, int32(5), *prop.IntValue)
			}
		}
		assert.True(t, foundTeam, "team custom property should exist")
		assert.True(t, foundPriority, "priority custom property should exist")
	})

	t.Run("TestSaveWithoutExperiment", func(t *testing.T) {
		// Test creating an experiment run without parent experiment (should still work)
		experimentRun := &models.ExperimentRunImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ExperimentRunAttributes{
				Name:       apiutils.Of("standalone-experiment-run"),
				ExternalID: apiutils.Of("standalone-experiment-run-ext-123"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Standalone experiment run without parent experiment"),
				},
			},
		}

		saved, err := repo.Save(experimentRun, nil)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "standalone-experiment-run", *saved.GetAttributes().Name)
		assert.Equal(t, "standalone-experiment-run-ext-123", *saved.GetAttributes().ExternalID)

		// Verify it can be retrieved
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		assert.Equal(t, "standalone-experiment-run", *retrieved.GetAttributes().Name)
	})

	t.Run("TestPagination", func(t *testing.T) {
		// First create a parent experiment
		experiment := &models.ExperimentImpl{
			TypeID: apiutils.Of(int32(experimentTypeID)),
			Attributes: &models.ExperimentAttributes{
				Name:       apiutils.Of("parent-experiment-pagination"),
				ExternalID: apiutils.Of("parent-experiment-pagination-ext-123"),
			},
		}
		savedExperiment, err := experimentRepo.Save(experiment)
		require.NoError(t, err)

		// Create multiple experiment runs for pagination testing
		for i := 0; i < 5; i++ {
			experimentRun := &models.ExperimentRunImpl{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ExperimentRunAttributes{
					Name:       apiutils.Of(fmt.Sprintf("page-experiment-run-%d", i)),
					ExternalID: apiutils.Of(fmt.Sprintf("page-experiment-run-ext-%d", i)),
				},
				Properties: &[]models.Properties{
					{
						Name:        "experiment_id",
						StringValue: apiutils.Of(fmt.Sprintf("%d", *savedExperiment.GetID())),
					},
				},
			}
			_, err := repo.Save(experimentRun, savedExperiment.GetID())
			require.NoError(t, err)
		}

		// Test pagination with page size 2
		pageSize := int32(2)
		listOptions := models.ExperimentRunListOptions{
			Pagination: models.Pagination{
				PageSize: &pageSize,
			},
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Should return at most 2 items
		assert.LessOrEqual(t, len(result.Items), 2)
		assert.Equal(t, pageSize, result.PageSize)
		assert.Equal(t, int32(len(result.Items)), result.Size)

		// If we have more items, there should be a next page token
		if len(result.Items) == 2 {
			assert.NotEmpty(t, result.NextPageToken)
		}
	})

	t.Run("TestEmptyResults", func(t *testing.T) {
		// Test listing with filter that returns no results
		listOptions := models.ExperimentRunListOptions{
			Name: apiutils.Of("non-existent-experiment-run"),
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, len(result.Items))
		assert.Equal(t, int32(0), result.Size)
		assert.Empty(t, result.NextPageToken)
	})
}

func TestExperimentRunRepository_FilterQuery(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Get the actual type IDs from the database
	experimentRunTypeID := getExperimentRunTypeID(t, db)
	experimentTypeID := getExperimentTypeID(t, db)

	experimentRunRepo := service.NewExperimentRunRepository(db, experimentRunTypeID)
	experimentRepo := service.NewExperimentRepository(db, experimentTypeID)

	// Create a parent experiment
	experiment := &models.ExperimentImpl{
		TypeID: apiutils.Of(int32(experimentTypeID)),
		Attributes: &models.ExperimentAttributes{
			Name: apiutils.Of("test-parent-experiment"),
		},
	}
	savedExperiment, err := experimentRepo.Save(experiment)
	require.NoError(t, err)
	experimentID := savedExperiment.GetID()

	// Create multiple experiment runs with different properties for filtering
	experimentRun1 := &models.ExperimentRunImpl{
		TypeID: apiutils.Of(int32(experimentRunTypeID)),
		Attributes: &models.ExperimentRunAttributes{
			Name: apiutils.Of(fmt.Sprintf("%d:pytorch-experiment-run", *experimentID)),
		},
		Properties: &[]models.Properties{
			{
				Name:        "status",
				StringValue: apiutils.Of("RUNNING"),
			},
			{
				Name:     "startTimeSinceEpoch",
				IntValue: apiutils.Of(int32(1640995200)), // 2022-01-01
			},
		},
		CustomProperties: &[]models.Properties{
			{
				Name:             "framework",
				StringValue:      apiutils.Of("pytorch"),
				IsCustomProperty: true,
			},
			{
				Name:             "epochs",
				IntValue:         apiutils.Of(int32(100)),
				IsCustomProperty: true,
			},
			{
				Name:             "learning_rate",
				DoubleValue:      apiutils.Of(0.001),
				IsCustomProperty: true,
			},
		},
	}
	_, err = experimentRunRepo.Save(experimentRun1, savedExperiment.GetID())
	require.NoError(t, err)

	experimentRun2 := &models.ExperimentRunImpl{
		TypeID: apiutils.Of(int32(experimentRunTypeID)),
		Attributes: &models.ExperimentRunAttributes{
			Name: apiutils.Of(fmt.Sprintf("%d:tensorflow-experiment-run", *experimentID)),
		},
		Properties: &[]models.Properties{
			{
				Name:        "status",
				StringValue: apiutils.Of("COMPLETED"),
			},
			{
				Name:     "startTimeSinceEpoch",
				IntValue: apiutils.Of(int32(1641081600)), // 2022-01-02
			},
		},
		CustomProperties: &[]models.Properties{
			{
				Name:             "framework",
				StringValue:      apiutils.Of("tensorflow"),
				IsCustomProperty: true,
			},
			{
				Name:             "epochs",
				IntValue:         apiutils.Of(int32(50)),
				IsCustomProperty: true,
			},
			{
				Name:             "learning_rate",
				DoubleValue:      apiutils.Of(0.01),
				IsCustomProperty: true,
			},
		},
	}
	_, err = experimentRunRepo.Save(experimentRun2, savedExperiment.GetID())
	require.NoError(t, err)

	experimentRun3 := &models.ExperimentRunImpl{
		TypeID: apiutils.Of(int32(experimentRunTypeID)),
		Attributes: &models.ExperimentRunAttributes{
			Name: apiutils.Of(fmt.Sprintf("%d:sklearn-experiment-run", *experimentID)),
		},
		Properties: &[]models.Properties{
			{
				Name:        "status",
				StringValue: apiutils.Of("FAILED"),
			},
			{
				Name:     "startTimeSinceEpoch",
				IntValue: apiutils.Of(int32(1641168000)), // 2022-01-03
			},
		},
		CustomProperties: &[]models.Properties{
			{
				Name:             "framework",
				StringValue:      apiutils.Of("sklearn"),
				IsCustomProperty: true,
			},
			{
				Name:             "epochs",
				IntValue:         apiutils.Of(int32(10)),
				IsCustomProperty: true,
			},
			{
				Name:             "learning_rate",
				DoubleValue:      apiutils.Of(0.1),
				IsCustomProperty: true,
			},
		},
	}
	_, err = experimentRunRepo.Save(experimentRun3, savedExperiment.GetID())
	require.NoError(t, err)

	// Test core property filtering
	t.Run("CorePropertyFilter", func(t *testing.T) {
		filterQuery := `name = "pytorch-experiment-run"`
		pageSize := int32(10)
		listOptions := models.ExperimentRunListOptions{
			Pagination: models.Pagination{
				PageSize:    &pageSize,
				FilterQuery: &filterQuery,
			},
		}

		result, err := experimentRunRepo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items))
		assert.Equal(t, fmt.Sprintf("%d:pytorch-experiment-run", *experimentID), *result.Items[0].GetAttributes().Name)
	})

	// Test custom property filtering
	t.Run("CustomPropertyFilter", func(t *testing.T) {
		filterQuery := `framework = "tensorflow"`
		pageSize := int32(10)
		listOptions := models.ExperimentRunListOptions{
			Pagination: models.Pagination{
				PageSize:    &pageSize,
				FilterQuery: &filterQuery,
			},
		}

		result, err := experimentRunRepo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items))
		assert.Equal(t, fmt.Sprintf("%d:tensorflow-experiment-run", *experimentID), *result.Items[0].GetAttributes().Name)
	})

	// Test numeric custom property filtering
	t.Run("NumericCustomPropertyFilter", func(t *testing.T) {
		filterQuery := `epochs >= 50`
		pageSize := int32(10)
		listOptions := models.ExperimentRunListOptions{
			Pagination: models.Pagination{
				PageSize:    &pageSize,
				FilterQuery: &filterQuery,
			},
		}

		result, err := experimentRunRepo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items)) // pytorch (100) and tensorflow (50)
	})

	// Test double custom property filtering
	t.Run("DoubleCustomPropertyFilter", func(t *testing.T) {
		filterQuery := `learning_rate <= 0.01`
		pageSize := int32(10)
		listOptions := models.ExperimentRunListOptions{
			Pagination: models.Pagination{
				PageSize:    &pageSize,
				FilterQuery: &filterQuery,
			},
		}

		result, err := experimentRunRepo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items)) // pytorch (0.001) and tensorflow (0.01)
	})

	// Test standard property filtering
	t.Run("StandardPropertyFilter", func(t *testing.T) {
		filterQuery := `status = "COMPLETED"`
		pageSize := int32(10)
		listOptions := models.ExperimentRunListOptions{
			Pagination: models.Pagination{
				PageSize:    &pageSize,
				FilterQuery: &filterQuery,
			},
		}

		result, err := experimentRunRepo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items))
		assert.Equal(t, fmt.Sprintf("%d:tensorflow-experiment-run", *experimentID), *result.Items[0].GetAttributes().Name)
	})

	// Test complex AND filter
	t.Run("ComplexANDFilter", func(t *testing.T) {
		filterQuery := `framework = "pytorch" AND epochs = 100`
		pageSize := int32(10)
		listOptions := models.ExperimentRunListOptions{
			Pagination: models.Pagination{
				PageSize:    &pageSize,
				FilterQuery: &filterQuery,
			},
		}

		result, err := experimentRunRepo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items))
		assert.Equal(t, fmt.Sprintf("%d:pytorch-experiment-run", *experimentID), *result.Items[0].GetAttributes().Name)
	})

	// Test complex OR filter
	t.Run("ComplexORFilter", func(t *testing.T) {
		filterQuery := `framework = "pytorch" OR framework = "sklearn"`
		pageSize := int32(10)
		listOptions := models.ExperimentRunListOptions{
			Pagination: models.Pagination{
				PageSize:    &pageSize,
				FilterQuery: &filterQuery,
			},
		}

		result, err := experimentRunRepo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items)) // pytorch and sklearn
	})

	// Test ILIKE operator
	t.Run("ILIKEFilter", func(t *testing.T) {
		filterQuery := `name ILIKE "%experiment%"`
		pageSize := int32(10)
		listOptions := models.ExperimentRunListOptions{
			Pagination: models.Pagination{
				PageSize:    &pageSize,
				FilterQuery: &filterQuery,
			},
		}

		result, err := experimentRunRepo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 3, len(result.Items)) // All experiment runs contain "experiment"
	})

	// Test mixed core and custom property filter
	t.Run("MixedCoreAndCustomFilter", func(t *testing.T) {
		filterQuery := `name ILIKE "%pytorch%" AND learning_rate < 0.01`
		pageSize := int32(10)
		listOptions := models.ExperimentRunListOptions{
			Pagination: models.Pagination{
				PageSize:    &pageSize,
				FilterQuery: &filterQuery,
			},
		}

		result, err := experimentRunRepo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items))
		assert.Equal(t, fmt.Sprintf("%d:pytorch-experiment-run", *experimentID), *result.Items[0].GetAttributes().Name)
	})

	// Test invalid filter query
	t.Run("InvalidFilterQuery", func(t *testing.T) {
		filterQuery := `invalid syntax =`
		pageSize := int32(10)
		listOptions := models.ExperimentRunListOptions{
			Pagination: models.Pagination{
				PageSize:    &pageSize,
				FilterQuery: &filterQuery,
			},
		}

		result, err := experimentRunRepo.List(listOptions)
		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid filter query")
	})

	// Test with parentheses grouping
	t.Run("ParenthesesGrouping", func(t *testing.T) {
		filterQuery := `(framework = "pytorch" OR framework = "tensorflow") AND epochs > 25`
		pageSize := int32(10)
		listOptions := models.ExperimentRunListOptions{
			Pagination: models.Pagination{
				PageSize:    &pageSize,
				FilterQuery: &filterQuery,
			},
		}

		result, err := experimentRunRepo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items)) // pytorch (100 epochs) and tensorflow (50 epochs)
	})
}
