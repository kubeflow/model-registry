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

func TestExperimentRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Get the actual Experiment type ID from the database
	typeID := getExperimentTypeID(t, db)
	repo := service.NewExperimentRepository(db, typeID)

	t.Run("TestSave", func(t *testing.T) {
		// Test creating a new experiment
		experiment := &models.ExperimentImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ExperimentAttributes{
				Name:       apiutils.Of("test-experiment"),
				ExternalID: apiutils.Of("experiment-ext-123"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Test experiment description"),
				},
				{
					Name:        "owner",
					StringValue: apiutils.Of("test-user"),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "custom-experiment-prop",
					StringValue:      apiutils.Of("custom-experiment-value"),
					IsCustomProperty: true,
				},
			},
		}

		saved, err := repo.Save(experiment)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-experiment", *saved.GetAttributes().Name)
		assert.Equal(t, "experiment-ext-123", *saved.GetAttributes().ExternalID)
		assert.NotNil(t, saved.GetAttributes().CreateTimeSinceEpoch)
		assert.NotNil(t, saved.GetAttributes().LastUpdateTimeSinceEpoch)

		// Test updating the same experiment
		experiment.ID = saved.GetID()
		experiment.GetAttributes().Name = apiutils.Of("updated-experiment")
		experiment.GetAttributes().ExternalID = apiutils.Of("updated-experiment-ext-123")
		// Preserve CreateTimeSinceEpoch from the saved entity (simulating what OpenAPI converter would do)
		experiment.GetAttributes().CreateTimeSinceEpoch = saved.GetAttributes().CreateTimeSinceEpoch

		updated, err := repo.Save(experiment)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, "updated-experiment", *updated.GetAttributes().Name)
		assert.Equal(t, "updated-experiment-ext-123", *updated.GetAttributes().ExternalID)
		assert.Equal(t, *saved.GetAttributes().CreateTimeSinceEpoch, *updated.GetAttributes().CreateTimeSinceEpoch)
		assert.Greater(t, *updated.GetAttributes().LastUpdateTimeSinceEpoch, *saved.GetAttributes().LastUpdateTimeSinceEpoch)
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// First create an experiment to retrieve
		experiment := &models.ExperimentImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ExperimentAttributes{
				Name:       apiutils.Of("get-test-experiment"),
				ExternalID: apiutils.Of("get-experiment-ext-123"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Experiment for get test"),
				},
			},
		}

		saved, err := repo.Save(experiment)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-test-experiment", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-experiment-ext-123", *retrieved.GetAttributes().ExternalID)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "experiment by id not found")
	})

	t.Run("TestList", func(t *testing.T) {
		// Create multiple experiments for listing
		testExperiments := []*models.ExperimentImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ExperimentAttributes{
					Name:       apiutils.Of("list-experiment-1"),
					ExternalID: apiutils.Of("list-experiment-ext-1"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("First experiment"),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ExperimentAttributes{
					Name:       apiutils.Of("list-experiment-2"),
					ExternalID: apiutils.Of("list-experiment-ext-2"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("Second experiment"),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ExperimentAttributes{
					Name:       apiutils.Of("list-experiment-3"),
					ExternalID: apiutils.Of("list-experiment-ext-3"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("Third experiment"),
					},
				},
			},
		}

		for _, experiment := range testExperiments {
			_, err := repo.Save(experiment)
			require.NoError(t, err)
		}

		// Test listing all experiments with basic pagination
		pageSize := int32(10)
		listOptions := models.ExperimentListOptions{}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // At least our 3 test experiments

		// Test listing by name
		listOptions = models.ExperimentListOptions{
			Name: apiutils.Of("list-experiment-1"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-experiment-1", *result.Items[0].GetAttributes().Name)
		}

		// Test listing by external ID
		listOptions = models.ExperimentListOptions{
			ExternalID: apiutils.Of("list-experiment-ext-2"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-experiment-ext-2", *result.Items[0].GetAttributes().ExternalID)
		}

		// Test ordering by ID (deterministic)
		listOptions = models.ExperimentListOptions{
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
		// Create experiments sequentially with time delays to ensure deterministic ordering
		experiment1 := &models.ExperimentImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ExperimentAttributes{
				Name:       apiutils.Of("time-test-experiment-1"),
				ExternalID: apiutils.Of("time-experiment-ext-1"),
			},
		}
		saved1, err := repo.Save(experiment1)
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		experiment2 := &models.ExperimentImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ExperimentAttributes{
				Name:       apiutils.Of("time-test-experiment-2"),
				ExternalID: apiutils.Of("time-experiment-ext-2"),
			},
		}
		saved2, err := repo.Save(experiment2)
		require.NoError(t, err)

		// Test ordering by CREATE_TIME
		pageSize := int32(10)
		listOptions := models.ExperimentListOptions{
			Pagination: models.Pagination{
				OrderBy: apiutils.Of("CREATE_TIME"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test experiments in the results
		var foundExperiment1, foundExperiment2 models.Experiment
		var index1, index2 = -1, -1

		for i, item := range result.Items {
			if *item.GetID() == *saved1.GetID() {
				foundExperiment1 = item
				index1 = i
			}
			if *item.GetID() == *saved2.GetID() {
				foundExperiment2 = item
				index2 = i
			}
		}

		// Verify both experiments were found and experiment1 comes before experiment2 (ascending order)
		require.NotEqual(t, -1, index1, "Experiment 1 should be found in results")
		require.NotEqual(t, -1, index2, "Experiment 2 should be found in results")
		assert.Less(t, index1, index2, "Experiment 1 should come before Experiment 2 when ordered by CREATE_TIME")
		assert.Less(t, *foundExperiment1.GetAttributes().CreateTimeSinceEpoch, *foundExperiment2.GetAttributes().CreateTimeSinceEpoch, "Experiment 1 should have earlier create time")
	})

	t.Run("TestSaveWithProperties", func(t *testing.T) {
		experiment := &models.ExperimentImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ExperimentAttributes{
				Name:       apiutils.Of("props-test-experiment"),
				ExternalID: apiutils.Of("props-experiment-ext-123"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Experiment with properties"),
				},
				{
					Name:        "owner",
					StringValue: apiutils.Of("test-user"),
				},
				{
					Name:        "state",
					StringValue: apiutils.Of("LIVE"),
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

		saved, err := repo.Save(experiment)
		require.NoError(t, err)
		require.NotNil(t, saved)

		// Verify properties were saved
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		assert.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 3) // description, owner, state

		assert.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 2) // team, priority

		// Verify specific properties exist
		foundDescription := false
		foundOwner := false
		foundState := false
		for _, prop := range *retrieved.GetProperties() {
			switch prop.Name {
			case "description":
				foundDescription = true
				assert.Equal(t, "Experiment with properties", *prop.StringValue)
			case "owner":
				foundOwner = true
				assert.Equal(t, "test-user", *prop.StringValue)
			case "state":
				foundState = true
				assert.Equal(t, "LIVE", *prop.StringValue)
			}
		}
		assert.True(t, foundDescription, "description property should exist")
		assert.True(t, foundOwner, "owner property should exist")
		assert.True(t, foundState, "state property should exist")

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

	t.Run("TestUpdateProperties", func(t *testing.T) {
		// Create experiment with initial properties
		experiment := &models.ExperimentImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ExperimentAttributes{
				Name:       apiutils.Of("update-props-experiment"),
				ExternalID: apiutils.Of("update-props-ext-123"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Initial description"),
				},
				{
					Name:        "owner",
					StringValue: apiutils.Of("initial-user"),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "team",
					StringValue:      apiutils.Of("initial-team"),
					IsCustomProperty: true,
				},
			},
		}

		saved, err := repo.Save(experiment)
		require.NoError(t, err)

		// Update with new properties
		experiment.ID = saved.GetID()
		// Preserve CreateTimeSinceEpoch from the saved entity (simulating what OpenAPI converter would do)
		experiment.GetAttributes().CreateTimeSinceEpoch = saved.GetAttributes().CreateTimeSinceEpoch
		experiment.Properties = &[]models.Properties{
			{
				Name:        "description",
				StringValue: apiutils.Of("Updated description"),
			},
			{
				Name:        "owner",
				StringValue: apiutils.Of("updated-user"),
			},
			{
				Name:        "state",
				StringValue: apiutils.Of("LIVE"),
			},
		}
		experiment.CustomProperties = &[]models.Properties{
			{
				Name:             "team",
				StringValue:      apiutils.Of("updated-team"),
				IsCustomProperty: true,
			},
			{
				Name:             "priority",
				IntValue:         apiutils.Of(int32(10)),
				IsCustomProperty: true,
			},
		}

		updated, err := repo.Save(experiment)
		require.NoError(t, err)

		// Verify properties were updated
		retrieved, err := repo.GetByID(*updated.GetID())
		require.NoError(t, err)

		assert.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 3) // description, owner, state

		assert.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 2) // team, priority

		// Verify updated values
		for _, prop := range *retrieved.GetProperties() {
			switch prop.Name {
			case "description":
				assert.Equal(t, "Updated description", *prop.StringValue)
			case "owner":
				assert.Equal(t, "updated-user", *prop.StringValue)
			case "state":
				assert.Equal(t, "LIVE", *prop.StringValue)
			}
		}

		for _, prop := range *retrieved.GetCustomProperties() {
			switch prop.Name {
			case "team":
				assert.Equal(t, "updated-team", *prop.StringValue)
			case "priority":
				assert.Equal(t, int32(10), *prop.IntValue)
			}
		}
	})

	t.Run("TestPagination", func(t *testing.T) {
		// Create multiple experiments for pagination testing
		for i := 0; i < 5; i++ {
			experiment := &models.ExperimentImpl{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ExperimentAttributes{
					Name:       apiutils.Of(fmt.Sprintf("page-experiment-%d", i)),
					ExternalID: apiutils.Of(fmt.Sprintf("page-experiment-ext-%d", i)),
				},
			}
			_, err := repo.Save(experiment)
			require.NoError(t, err)
		}

		// Test pagination with page size 2
		pageSize := int32(2)
		listOptions := models.ExperimentListOptions{
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
		listOptions := models.ExperimentListOptions{
			Name: apiutils.Of("non-existent-experiment"),
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, len(result.Items))
		assert.Equal(t, int32(0), result.Size)
		assert.Empty(t, result.NextPageToken)
	})
}
