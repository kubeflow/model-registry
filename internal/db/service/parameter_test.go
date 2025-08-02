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

func TestParameterRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Get the actual Parameter type ID from the database
	typeID := getParameterTypeID(t, db)
	repo := service.NewParameterRepository(db, typeID)

	// Also get model version type for creating related entities
	registeredModelTypeID := getRegisteredModelTypeID(t, db)
	registeredModelRepo := service.NewRegisteredModelRepository(db, registeredModelTypeID)

	modelVersionTypeID := getModelVersionTypeID(t, db)
	modelVersionRepo := service.NewModelVersionRepository(db, modelVersionTypeID)

	t.Run("TestSave", func(t *testing.T) {
		// Test creating a new parameter
		parameter := &models.ParameterImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ParameterAttributes{
				Name:         apiutils.Of("test-parameter"),
				ExternalID:   apiutils.Of("parameter-ext-123"),
				URI:          apiutils.Of("s3://bucket/parameter.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("parameter"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Test parameter description"),
				},
				{
					Name:        "value",
					StringValue: apiutils.Of("0.001"),
				},
				{
					Name:        "parameter_type",
					StringValue: apiutils.Of("number"),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "custom-parameter-prop",
					StringValue:      apiutils.Of("custom-parameter-value"),
					IsCustomProperty: true,
				},
			},
		}

		saved, err := repo.Save(parameter, nil)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-parameter", *saved.GetAttributes().Name)
		assert.Equal(t, "parameter-ext-123", *saved.GetAttributes().ExternalID)
		assert.Equal(t, "s3://bucket/parameter.json", *saved.GetAttributes().URI)
		assert.Equal(t, "LIVE", *saved.GetAttributes().State)
		assert.Equal(t, "parameter", *saved.GetAttributes().ArtifactType)
		assert.NotNil(t, saved.GetAttributes().CreateTimeSinceEpoch)
		assert.NotNil(t, saved.GetAttributes().LastUpdateTimeSinceEpoch)

		// Verify properties were saved
		assert.NotNil(t, saved.GetProperties())
		assert.Len(t, *saved.GetProperties(), 3) // description, value, parameter_type

		// Verify specific properties
		var foundDescription, foundValue, foundType bool
		for _, prop := range *saved.GetProperties() {
			switch prop.Name {
			case "description":
				foundDescription = true
				assert.Equal(t, "Test parameter description", *prop.StringValue)
			case "value":
				foundValue = true
				assert.Equal(t, "0.001", *prop.StringValue)
			case "parameter_type":
				foundType = true
				assert.Equal(t, "number", *prop.StringValue)
			}
		}
		assert.True(t, foundDescription, "description property should exist")
		assert.True(t, foundValue, "value property should exist")
		assert.True(t, foundType, "parameter_type property should exist")

		// Verify custom properties were saved
		assert.NotNil(t, saved.GetCustomProperties())
		assert.Len(t, *saved.GetCustomProperties(), 1)

		var foundCustomProp bool
		for _, prop := range *saved.GetCustomProperties() {
			if prop.Name == "custom-parameter-prop" {
				foundCustomProp = true
				assert.Equal(t, "custom-parameter-value", *prop.StringValue)
				assert.True(t, prop.IsCustomProperty)
			}
		}
		assert.True(t, foundCustomProp, "custom-parameter-prop should exist")

		// Test updating the same parameter
		parameter.ID = saved.GetID()
		parameter.GetAttributes().Name = apiutils.Of("updated-parameter")
		parameter.GetAttributes().State = apiutils.Of("PENDING")
		// Preserve CreateTimeSinceEpoch from the saved entity (simulating what OpenAPI converter would do)
		parameter.GetAttributes().CreateTimeSinceEpoch = saved.GetAttributes().CreateTimeSinceEpoch

		updated, err := repo.Save(parameter, nil)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, "updated-parameter", *updated.GetAttributes().Name)
		assert.Equal(t, "PENDING", *updated.GetAttributes().State)
		assert.Equal(t, *saved.GetAttributes().CreateTimeSinceEpoch, *updated.GetAttributes().CreateTimeSinceEpoch)
		assert.Greater(t, *updated.GetAttributes().LastUpdateTimeSinceEpoch, *saved.GetAttributes().LastUpdateTimeSinceEpoch)
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// First create a parameter to retrieve
		parameter := &models.ParameterImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ParameterAttributes{
				Name:         apiutils.Of("get-test-parameter"),
				ExternalID:   apiutils.Of("get-parameter-ext-123"),
				URI:          apiutils.Of("s3://bucket/get-parameter.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("parameter"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Parameter for get test"),
				},
				{
					Name:        "value",
					StringValue: apiutils.Of("0.005"),
				},
				{
					Name:        "parameter_type",
					StringValue: apiutils.Of("string"),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "experiment-phase",
					StringValue:      apiutils.Of("validation"),
					IsCustomProperty: true,
				},
			},
		}

		saved, err := repo.Save(parameter, nil)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-test-parameter", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-parameter-ext-123", *retrieved.GetAttributes().ExternalID)
		assert.Equal(t, "s3://bucket/get-parameter.json", *retrieved.GetAttributes().URI)
		assert.Equal(t, "LIVE", *retrieved.GetAttributes().State)

		// Verify type-specific properties were retrieved
		assert.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 3)

		var foundDescription, foundValue, foundType bool
		for _, prop := range *retrieved.GetProperties() {
			switch prop.Name {
			case "description":
				foundDescription = true
				assert.Equal(t, "Parameter for get test", *prop.StringValue)
			case "value":
				foundValue = true
				assert.Equal(t, "0.005", *prop.StringValue)
			case "parameter_type":
				foundType = true
				assert.Equal(t, "string", *prop.StringValue)
			}
		}
		assert.True(t, foundDescription, "description should be retrieved")
		assert.True(t, foundValue, "value should be retrieved")
		assert.True(t, foundType, "parameter_type should be retrieved")

		// Verify custom properties were retrieved
		assert.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 1)

		customProp := (*retrieved.GetCustomProperties())[0]
		assert.Equal(t, "experiment-phase", customProp.Name)
		assert.Equal(t, "validation", *customProp.StringValue)
		assert.True(t, customProp.IsCustomProperty)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parameter by id not found")
	})

	t.Run("TestList", func(t *testing.T) {
		// Create multiple parameters for listing
		testParameters := []*models.ParameterImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ParameterAttributes{
					Name:         apiutils.Of("list-parameter-1"),
					ExternalID:   apiutils.Of("list-parameter-ext-1"),
					URI:          apiutils.Of("s3://bucket/list-parameter-1.json"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("parameter"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("First parameter"),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ParameterAttributes{
					Name:         apiutils.Of("list-parameter-2"),
					ExternalID:   apiutils.Of("list-parameter-ext-2"),
					URI:          apiutils.Of("s3://bucket/list-parameter-2.json"),
					State:        apiutils.Of("PENDING"),
					ArtifactType: apiutils.Of("parameter"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("Second parameter"),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ParameterAttributes{
					Name:         apiutils.Of("list-parameter-3"),
					ExternalID:   apiutils.Of("list-parameter-ext-3"),
					URI:          apiutils.Of("s3://bucket/list-parameter-3.json"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("parameter"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("Third parameter"),
					},
				},
			},
		}

		for _, parameter := range testParameters {
			_, err := repo.Save(parameter, nil)
			require.NoError(t, err)
		}

		// Test listing all parameters with basic pagination
		pageSize := int32(10)
		listOptions := models.ParameterListOptions{}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // At least our 3 test parameters

		// Test listing by name
		listOptions = models.ParameterListOptions{
			Name: apiutils.Of("list-parameter-1"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-parameter-1", *result.Items[0].GetAttributes().Name)
		}

		// Test listing by external ID
		listOptions = models.ParameterListOptions{
			ExternalID: apiutils.Of("list-parameter-ext-2"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-parameter-ext-2", *result.Items[0].GetAttributes().ExternalID)
		}

		// Test ordering by ID (deterministic)
		listOptions = models.ParameterListOptions{
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
		// Create parameters sequentially with time delays to ensure deterministic ordering
		parameter1 := &models.ParameterImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ParameterAttributes{
				Name:         apiutils.Of("time-test-parameter-1"),
				URI:          apiutils.Of("s3://bucket/time-parameter-1.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("parameter"),
			},
		}
		saved1, err := repo.Save(parameter1, nil)
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		parameter2 := &models.ParameterImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ParameterAttributes{
				Name:         apiutils.Of("time-test-parameter-2"),
				URI:          apiutils.Of("s3://bucket/time-parameter-2.json"),
				State:        apiutils.Of("PENDING"),
				ArtifactType: apiutils.Of("parameter"),
			},
		}
		saved2, err := repo.Save(parameter2, nil)
		require.NoError(t, err)

		// Test ordering by CREATE_TIME
		pageSize := int32(10)
		listOptions := models.ParameterListOptions{
			Pagination: models.Pagination{
				OrderBy: apiutils.Of("CREATE_TIME"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test parameters in the results
		var foundParameter1, foundParameter2 models.Parameter
		var index1, index2 = -1, -1

		for i, item := range result.Items {
			if *item.GetID() == *saved1.GetID() {
				foundParameter1 = item
				index1 = i
			}
			if *item.GetID() == *saved2.GetID() {
				foundParameter2 = item
				index2 = i
			}
		}

		// Verify both parameters were found and parameter1 comes before parameter2 (ascending order)
		require.NotEqual(t, -1, index1, "Parameter 1 should be found in results")
		require.NotEqual(t, -1, index2, "Parameter 2 should be found in results")
		assert.Less(t, index1, index2, "Parameter 1 should come before Parameter 2 when ordered by CREATE_TIME")
		assert.Less(t, *foundParameter1.GetAttributes().CreateTimeSinceEpoch, *foundParameter2.GetAttributes().CreateTimeSinceEpoch, "Parameter 1 should have earlier create time")
	})

	t.Run("TestSaveWithModelVersion", func(t *testing.T) {
		// First create a registered model and model version
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model-for-parameter"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		modelVersion := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(modelVersionTypeID)),
			Attributes: &models.ModelVersionAttributes{
				Name: apiutils.Of("test-model-version-for-parameter"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedModelVersion, err := modelVersionRepo.Save(modelVersion)
		require.NoError(t, err)

		// Test creating a parameter with model version attribution
		parameter := &models.ParameterImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ParameterAttributes{
				Name:         apiutils.Of(fmt.Sprintf("%d:model-version-parameter", *savedModelVersion.GetID())),
				URI:          apiutils.Of("s3://bucket/model-version-parameter.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("parameter"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Parameter associated with model version"),
				},
			},
		}

		saved, err := repo.Save(parameter, savedModelVersion.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, fmt.Sprintf("%d:model-version-parameter", *savedModelVersion.GetID()), *saved.GetAttributes().Name)
		assert.Equal(t, "s3://bucket/model-version-parameter.json", *saved.GetAttributes().URI)

		// Test listing by model version ID
		listOptions := models.ParameterListOptions{
			ParentResourceID: savedModelVersion.GetID(),
		}
		listOptions.PageSize = apiutils.Of(int32(10))

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 1) // Should find our parameter

		// Verify the found parameter
		found := false
		for _, item := range result.Items {
			if *item.GetID() == *saved.GetID() {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find the parameter associated with the model version")
	})

	t.Run("TestSaveWithProperties", func(t *testing.T) {
		parameter := &models.ParameterImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ParameterAttributes{
				Name:         apiutils.Of("props-test-parameter"),
				ExternalID:   apiutils.Of("props-parameter-ext-123"),
				URI:          apiutils.Of("s3://bucket/props-parameter.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("parameter"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Parameter with properties"),
				},
				{
					Name:        "value",
					StringValue: apiutils.Of("learning_rate=0.001"),
				},
				{
					Name:        "type",
					StringValue: apiutils.Of("hyperparameter"),
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

		saved, err := repo.Save(parameter, nil)
		require.NoError(t, err)
		require.NotNil(t, saved)

		// Verify properties were saved
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		assert.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 3) // description, value, type

		assert.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 2) // team, priority

		// Verify specific properties exist
		foundDescription := false
		foundValue := false
		foundType := false
		for _, prop := range *retrieved.GetProperties() {
			switch prop.Name {
			case "description":
				foundDescription = true
				assert.Equal(t, "Parameter with properties", *prop.StringValue)
			case "value":
				foundValue = true
				assert.Equal(t, "learning_rate=0.001", *prop.StringValue)
			case "type":
				foundType = true
				assert.Equal(t, "hyperparameter", *prop.StringValue)
			}
		}
		assert.True(t, foundDescription, "description property should exist")
		assert.True(t, foundValue, "value property should exist")
		assert.True(t, foundType, "type property should exist")

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

	t.Run("TestPagination", func(t *testing.T) {
		// Create multiple parameters for pagination testing
		for i := 0; i < 5; i++ {
			parameter := &models.ParameterImpl{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ParameterAttributes{
					Name:         apiutils.Of(fmt.Sprintf("page-parameter-%d", i)),
					URI:          apiutils.Of(fmt.Sprintf("s3://bucket/page-parameter-%d.json", i)),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("parameter"),
				},
			}
			_, err := repo.Save(parameter, nil)
			require.NoError(t, err)
		}

		// Test pagination with page size 2
		pageSize := int32(2)
		listOptions := models.ParameterListOptions{
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
		listOptions := models.ParameterListOptions{
			Name: apiutils.Of("non-existent-parameter"),
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, len(result.Items))
		assert.Equal(t, int32(0), result.Size)
		assert.Empty(t, result.NextPageToken)
	})

	t.Run("TestListByState", func(t *testing.T) {
		// Create parameters with different states for testing
		parameters := []*models.ParameterImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ParameterAttributes{
					Name:         apiutils.Of("state-parameter-1"),
					URI:          apiutils.Of("s3://bucket/state-parameter-1.json"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("parameter"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("Live parameter"),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ParameterAttributes{
					Name:         apiutils.Of("state-parameter-2"),
					URI:          apiutils.Of("s3://bucket/state-parameter-2.json"),
					State:        apiutils.Of("PENDING"),
					ArtifactType: apiutils.Of("parameter"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("Pending parameter"),
					},
				},
			},
		}

		for _, parameter := range parameters {
			_, err := repo.Save(parameter, nil)
			require.NoError(t, err)
		}

		// Test listing all parameters (should find at least our 2 test parameters)
		listOptions := models.ParameterListOptions{}
		listOptions.PageSize = apiutils.Of(int32(10))

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Should find at least one parameter with each state
		foundLive := false
		foundPending := false
		for _, item := range result.Items {
			if item.GetAttributes().State != nil {
				if *item.GetAttributes().State == "LIVE" {
					foundLive = true
				}
				if *item.GetAttributes().State == "PENDING" {
					foundPending = true
				}
			}
		}
		assert.True(t, foundLive, "Should find parameter with LIVE state")
		assert.True(t, foundPending, "Should find parameter with PENDING state")
	})
}
