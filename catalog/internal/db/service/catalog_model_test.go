package service_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/internal/apiutils"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCatalogModelRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupPostgresWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	// Create or get the CatalogModel type ID
	typeID := getCatalogModelTypeID(t, sharedDB)
	repo := service.NewCatalogModelRepository(sharedDB, typeID)

	t.Run("TestSave", func(t *testing.T) {
		// Test creating a new catalog model
		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model"),
				ExternalID: apiutils.Of("catalog-ext-123"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Test catalog model description"),
				},
			},
			CustomProperties: &[]dbmodels.Properties{
				{
					Name:        "custom-prop",
					StringValue: apiutils.Of("custom-value"),
				},
			},
		}

		saved, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-catalog-model", *saved.GetAttributes().Name)
		assert.Equal(t, "catalog-ext-123", *saved.GetAttributes().ExternalID)

		// Test updating the same model
		catalogModel.ID = saved.GetID()
		catalogModel.GetAttributes().Name = apiutils.Of("updated-catalog-model")
		// Preserve CreateTimeSinceEpoch from the saved entity
		catalogModel.GetAttributes().CreateTimeSinceEpoch = saved.GetAttributes().CreateTimeSinceEpoch

		updated, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, "updated-catalog-model", *updated.GetAttributes().Name)
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// First create a model to retrieve
		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("get-test-catalog-model"),
				ExternalID: apiutils.Of("get-catalog-ext-123"),
			},
		}

		saved, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-test-catalog-model", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-catalog-ext-123", *retrieved.GetAttributes().ExternalID)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.ErrorIs(t, err, service.ErrCatalogModelNotFound)
	})

	t.Run("TestList", func(t *testing.T) {
		// Create multiple models for listing
		testModels := []*models.CatalogModelImpl{
			{
				Attributes: &models.CatalogModelAttributes{
					Name:       apiutils.Of("list-catalog-model-1"),
					ExternalID: apiutils.Of("list-catalog-ext-1"),
				},
			},
			{
				Attributes: &models.CatalogModelAttributes{
					Name:       apiutils.Of("list-catalog-model-2"),
					ExternalID: apiutils.Of("list-catalog-ext-2"),
				},
			},
		}

		// Save all test models
		var savedModels []models.CatalogModel
		for _, model := range testModels {
			saved, err := repo.Save(model)
			require.NoError(t, err)
			savedModels = append(savedModels, saved)
		}

		// Test listing all models
		listOptions := models.CatalogModelListOptions{}
		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 2) // At least our 2 test models

		// Test filtering by name
		nameFilter := "list-catalog-model-1"
		listOptions = models.CatalogModelListOptions{
			Name: &nameFilter,
		}
		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items))
		assert.Equal(t, "list-catalog-model-1", *result.Items[0].GetAttributes().Name)

		// Test filtering by external ID
		externalIDFilter := "list-catalog-ext-2"
		listOptions = models.CatalogModelListOptions{
			ExternalID: &externalIDFilter,
		}
		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items))
		assert.Equal(t, "list-catalog-ext-2", *result.Items[0].GetAttributes().ExternalID)
	})

	t.Run("TestGetByName", func(t *testing.T) {
		// First create a model to retrieve by name
		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("get-by-name-test-model"),
				ExternalID: apiutils.Of("get-by-name-ext-123"),
			},
		}

		saved, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by name
		retrieved, err := repo.GetByName("get-by-name-test-model")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-by-name-test-model", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-by-name-ext-123", *retrieved.GetAttributes().ExternalID)

		// Test retrieving non-existent name
		_, err = repo.GetByName("non-existent-model")
		assert.ErrorIs(t, err, service.ErrCatalogModelNotFound)
	})

	t.Run("TestUpdateWithID", func(t *testing.T) {
		// First create a model
		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("update-test-model"),
				ExternalID: apiutils.Of("update-ext-123"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "version",
					StringValue: apiutils.Of("1.0.0"),
				},
			},
			CustomProperties: &[]dbmodels.Properties{
				{
					Name:        "environment",
					StringValue: apiutils.Of("dev"),
				},
			},
		}

		saved, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Update the model with ID specified
		updateModel := &models.CatalogModelImpl{
			ID: saved.GetID(), // Specify the ID for update
			Attributes: &models.CatalogModelAttributes{
				Name:                 apiutils.Of("updated-test-model"),
				ExternalID:           apiutils.Of("updated-ext-456"),
				CreateTimeSinceEpoch: saved.GetAttributes().CreateTimeSinceEpoch, // Preserve create time
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "version",
					StringValue: apiutils.Of("2.0.0"), // Updated version
				},
				{
					Name:        "description",
					StringValue: apiutils.Of("Updated description"), // New property
				},
			},
			CustomProperties: &[]dbmodels.Properties{
				{
					Name:        "environment",
					StringValue: apiutils.Of("prod"), // Updated environment
				},
			},
		}

		updated, err := repo.Save(updateModel)
		require.NoError(t, err)
		require.NotNil(t, updated)

		// Verify the update
		assert.Equal(t, *saved.GetID(), *updated.GetID()) // Same ID
		assert.Equal(t, "updated-test-model", *updated.GetAttributes().Name)
		assert.Equal(t, "updated-ext-456", *updated.GetAttributes().ExternalID)

		// Verify properties were updated
		require.NotNil(t, updated.GetProperties())
		assert.Len(t, *updated.GetProperties(), 2)

		// Verify custom properties were updated
		require.NotNil(t, updated.GetCustomProperties())
		assert.Len(t, *updated.GetCustomProperties(), 1)
	})

	t.Run("TestUpdateWithName", func(t *testing.T) {
		// First create a model
		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("update-by-name-model"),
				ExternalID: apiutils.Of("update-by-name-ext-123"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "status",
					StringValue: apiutils.Of("draft"),
				},
			},
		}

		saved, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Update the model without specifying ID (should lookup by name)
		updateModel := &models.CatalogModelImpl{
			// No ID specified - should trigger name lookup in Save method
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("update-by-name-model"), // Same name to trigger lookup
				ExternalID: apiutils.Of("updated-by-name-ext-789"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "status",
					StringValue: apiutils.Of("published"), // Updated status
				},
				{
					Name:        "category",
					StringValue: apiutils.Of("ml-model"), // New property
				},
			},
		}

		updated, err := repo.Save(updateModel)
		require.NoError(t, err)
		require.NotNil(t, updated)

		// Verify the update happened (same ID, updated fields)
		assert.Equal(t, *saved.GetID(), *updated.GetID()) // Should have same ID from lookup
		assert.Equal(t, "update-by-name-model", *updated.GetAttributes().Name)
		assert.Equal(t, "updated-by-name-ext-789", *updated.GetAttributes().ExternalID)

		// Verify properties were updated
		require.NotNil(t, updated.GetProperties())
		assert.Len(t, *updated.GetProperties(), 2)
	})

	t.Run("TestListWithPropertiesAndCustomProperties", func(t *testing.T) {
		// Create a model with both properties and custom properties
		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("props-test-catalog-model"),
				ExternalID: apiutils.Of("props-catalog-ext-123"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "version",
					StringValue: apiutils.Of("1.0.0"),
				},
				{
					Name:     "priority",
					IntValue: apiutils.Of(int32(5)),
				},
			},
			CustomProperties: &[]dbmodels.Properties{
				{
					Name:        "team",
					StringValue: apiutils.Of("ml-team"),
				},
				{
					Name:      "active",
					BoolValue: apiutils.Of(true),
				},
			},
		}

		saved, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, saved)

		// Retrieve and verify properties
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		// Check regular properties
		require.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 2)

		// Check custom properties
		require.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 2)
	})

	t.Run("TestAccuracySorting", func(t *testing.T) {
		// Get the CatalogMetricsArtifact type ID for creating accuracy metrics
		metricsTypeID := getCatalogMetricsArtifactTypeID(t, sharedDB)
		metricsRepo := service.NewCatalogMetricsArtifactRepository(sharedDB, metricsTypeID)

		// Create test models with different accuracy scores
		testModels := []struct {
			name     string
			accuracy *float64 // nil means no accuracy score
		}{
			{"high-accuracy-model", apiutils.Of(95.5)},
			{"medium-accuracy-model", apiutils.Of(75.0)},
			{"low-accuracy-model", apiutils.Of(45.2)},
			{"no-accuracy-model", nil},
			{"zero-accuracy-model", apiutils.Of(0.0)},
		}

		var savedModels []models.CatalogModel
		for _, testModel := range testModels {
			// Create the model
			catalogModel := &models.CatalogModelImpl{
				Attributes: &models.CatalogModelAttributes{
					Name:       apiutils.Of(testModel.name),
					ExternalID: apiutils.Of(testModel.name + "-ext"),
				},
			}

			savedModel, err := repo.Save(catalogModel)
			require.NoError(t, err)
			savedModels = append(savedModels, savedModel)

			// Create accuracy metrics artifact if accuracy score is provided
			if testModel.accuracy != nil {
				metricsArtifact := &models.CatalogMetricsArtifactImpl{
					Attributes: &models.CatalogMetricsArtifactAttributes{
						Name:        apiutils.Of(fmt.Sprintf("accuracy-metrics-%s", testModel.name)),
						ExternalID:  apiutils.Of(fmt.Sprintf("accuracy-metrics-%s", testModel.name)),
						MetricsType: models.MetricsTypeAccuracy,
					},
					CustomProperties: &[]dbmodels.Properties{
						{
							Name:        "overall_average",
							DoubleValue: testModel.accuracy,
						},
						{
							Name:        "benchmark1",
							DoubleValue: apiutils.Of(*testModel.accuracy + 1.0), // Individual benchmark score
						},
						{
							Name:        "benchmark2",
							DoubleValue: apiutils.Of(*testModel.accuracy - 1.0), // Individual benchmark score
						},
					},
				}

				_, err := metricsRepo.Save(metricsArtifact, savedModel.GetID())
				require.NoError(t, err)
			}
		}

		// Test ACCURACY sorting DESC (default)
		listOptions := models.CatalogModelListOptions{
			Pagination: dbmodels.Pagination{
				OrderBy:   apiutils.Of("artifacts.overall_average.double_value"),
				SortOrder: apiutils.Of("DESC"),
			},
		}
		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Verify DESC order: high accuracy first, then medium, low, zero, then models without accuracy scores
		var accuracyModelsFound []string
		for _, model := range result.Items {
			name := *model.GetAttributes().Name
			if name == "high-accuracy-model" || name == "medium-accuracy-model" ||
				name == "low-accuracy-model" || name == "zero-accuracy-model" || name == "no-accuracy-model" {
				accuracyModelsFound = append(accuracyModelsFound, name)
			}
		}

		// We should have found all our test models
		require.GreaterOrEqual(t, len(accuracyModelsFound), 5)

		// Check that high accuracy comes before medium accuracy
		highIdx := findIndex(accuracyModelsFound, "high-accuracy-model")
		mediumIdx := findIndex(accuracyModelsFound, "medium-accuracy-model")
		lowIdx := findIndex(accuracyModelsFound, "low-accuracy-model")
		zeroIdx := findIndex(accuracyModelsFound, "zero-accuracy-model")
		noAccIdx := findIndex(accuracyModelsFound, "no-accuracy-model")

		require.NotEqual(t, -1, highIdx, "high-accuracy-model not found in results")
		require.NotEqual(t, -1, mediumIdx, "medium-accuracy-model not found in results")
		require.NotEqual(t, -1, lowIdx, "low-accuracy-model not found in results")
		require.NotEqual(t, -1, zeroIdx, "zero-accuracy-model not found in results")
		require.NotEqual(t, -1, noAccIdx, "no-accuracy-model not found in results")

		// Verify DESC ordering: high > medium > low > zero > no-accuracy
		assert.Less(t, highIdx, mediumIdx, "high accuracy model should come before medium accuracy")
		assert.Less(t, mediumIdx, lowIdx, "medium accuracy model should come before low accuracy")
		assert.Less(t, lowIdx, zeroIdx, "low accuracy model should come before zero accuracy")
		assert.Less(t, zeroIdx, noAccIdx, "zero accuracy model should come before no accuracy")

		// Test ACCURACY sorting ASC
		listOptions = models.CatalogModelListOptions{
			Pagination: dbmodels.Pagination{
				OrderBy:   apiutils.Of("artifacts.overall_average.double_value"),
				SortOrder: apiutils.Of("ASC"),
			},
		}
		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test models in ASC results
		accuracyModelsFound = []string{}
		for _, model := range result.Items {
			name := *model.GetAttributes().Name
			if name == "high-accuracy-model" || name == "medium-accuracy-model" ||
				name == "low-accuracy-model" || name == "zero-accuracy-model" || name == "no-accuracy-model" {
				accuracyModelsFound = append(accuracyModelsFound, name)
			}
		}

		// Get indices for ASC order
		highIdxAsc := findIndex(accuracyModelsFound, "high-accuracy-model")
		mediumIdxAsc := findIndex(accuracyModelsFound, "medium-accuracy-model")
		lowIdxAsc := findIndex(accuracyModelsFound, "low-accuracy-model")
		zeroIdxAsc := findIndex(accuracyModelsFound, "zero-accuracy-model")
		noAccIdxAsc := findIndex(accuracyModelsFound, "no-accuracy-model")

		// Verify ASC ordering: zero < low < medium < high, with no-accuracy still last
		assert.Less(t, zeroIdxAsc, lowIdxAsc, "zero accuracy model should come before low accuracy in ASC")
		assert.Less(t, lowIdxAsc, mediumIdxAsc, "low accuracy model should come before medium accuracy in ASC")
		assert.Less(t, mediumIdxAsc, highIdxAsc, "medium accuracy model should come before high accuracy in ASC")
		assert.Less(t, highIdxAsc, noAccIdxAsc, "models with accuracy should come before models without accuracy in ASC")

		// Test fallback to standard sorting for non-ACCURACY orderBy
		listOptions = models.CatalogModelListOptions{
			Pagination: dbmodels.Pagination{
				OrderBy:   apiutils.Of("ID"),
				SortOrder: apiutils.Of("ASC"),
			},
		}
		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		// Should not error and should return results (detailed verification not needed since we're testing fallback)
		assert.Greater(t, len(result.Items), 0)
	})

	t.Run("TestAccuracySortingPagination", func(t *testing.T) {
		// Get the CatalogMetricsArtifact type ID for creating accuracy metrics
		metricsTypeID := getCatalogMetricsArtifactTypeID(t, sharedDB)
		metricsRepo := service.NewCatalogMetricsArtifactRepository(sharedDB, metricsTypeID)

		// Create 5 test models with accuracy scores for pagination testing
		// Use unique names to avoid interference with other tests
		testModels := []struct {
			name     string
			accuracy float64
		}{
			{"pagination-test-model-a", 95.0}, // Should be first in DESC order
			{"pagination-test-model-b", 85.0},
			{"pagination-test-model-c", 75.0},
			{"pagination-test-model-d", 65.0},
			{"pagination-test-model-e", 55.0}, // Should be last in DESC order
		}

		var savedModels []models.CatalogModel
		for _, testModel := range testModels {
			// Create the model
			catalogModel := &models.CatalogModelImpl{
				Attributes: &models.CatalogModelAttributes{
					Name:       apiutils.Of(testModel.name),
					ExternalID: apiutils.Of(testModel.name + "-ext"),
				},
			}

			savedModel, err := repo.Save(catalogModel)
			require.NoError(t, err)
			savedModels = append(savedModels, savedModel)

			// Create accuracy metrics artifact
			metricsArtifact := &models.CatalogMetricsArtifactImpl{
				Attributes: &models.CatalogMetricsArtifactAttributes{
					Name:        apiutils.Of(fmt.Sprintf("accuracy-metrics-%s", testModel.name)),
					ExternalID:  apiutils.Of(fmt.Sprintf("accuracy-metrics-%s", testModel.name)),
					MetricsType: models.MetricsTypeAccuracy,
				},
				CustomProperties: &[]dbmodels.Properties{
					{
						Name:        "overall_average",
						DoubleValue: &testModel.accuracy,
					},
				},
			}

			_, err = metricsRepo.Save(metricsArtifact, savedModel.GetID())
			require.NoError(t, err)
		}

		// Test pagination by collecting all pages
		// This approach is more robust and less sensitive to test interference
		listOptions := models.CatalogModelListOptions{
			Pagination: dbmodels.Pagination{
				OrderBy:   apiutils.Of("artifacts.overall_average.double_value"),
				SortOrder: apiutils.Of("DESC"),
				PageSize:  apiutils.Of(int32(2)),
			},
		}

		// Collect all our test models across pages
		var allPaginatedModels []models.CatalogModel
		var pageCount int
		currentToken := (*string)(nil)

		for {
			pageCount++
			if currentToken != nil {
				listOptions.Pagination.NextPageToken = currentToken
			}

			page, err := repo.List(listOptions)
			require.NoError(t, err)
			require.NotNil(t, page)
			assert.Equal(t, int32(2), page.PageSize)

			// Filter to only include our test models
			for _, model := range page.Items {
				if strings.HasPrefix(*model.GetAttributes().Name, "pagination-test-model-") {
					allPaginatedModels = append(allPaginatedModels, model)
				}
			}

			// Stop if no more pages or we've collected all our test models
			if page.NextPageToken == "" || len(allPaginatedModels) >= 5 {
				if page.NextPageToken == "" {
					t.Logf("Pagination completed in %d pages", pageCount)
				}
				break
			}
			currentToken = &page.NextPageToken

			// Safety check to prevent infinite loop
			if pageCount > 10 {
				t.Fatal("Too many pages, might be an infinite loop")
			}
		}

		// Verify we collected all our test models
		assert.GreaterOrEqual(t, len(allPaginatedModels), 5, "Should have found all pagination test models")

		// Extract names and verify ordering (DESC by accuracy)
		var modelNames []string
		for _, model := range allPaginatedModels {
			if strings.HasPrefix(*model.GetAttributes().Name, "pagination-test-model-") {
				modelNames = append(modelNames, *model.GetAttributes().Name)
			}
		}

		// Check that pagination preserved the correct ordering
		// In DESC order: a(95.0) -> b(85.0) -> c(75.0) -> d(65.0) -> e(55.0)
		expectedOrder := []string{
			"pagination-test-model-a", // 95.0 (highest)
			"pagination-test-model-b", // 85.0
			"pagination-test-model-c", // 75.0
			"pagination-test-model-d", // 65.0
			"pagination-test-model-e", // 55.0 (lowest)
		}

		// Verify our test models appear in correct order (allowing for other models in between)
		lastIndex := -1
		for _, expectedModel := range expectedOrder {
			foundIndex := -1
			for i, actualModel := range modelNames {
				if actualModel == expectedModel {
					foundIndex = i
					break
				}
			}
			assert.NotEqual(t, -1, foundIndex, "Should find model %s", expectedModel)
			if foundIndex != -1 {
				assert.Greater(t, foundIndex, lastIndex, "Model %s should appear after previous models in DESC order", expectedModel)
				lastIndex = foundIndex
			}
		}

		// Test ASC pagination briefly to verify token generation works in both directions
		listOptions = models.CatalogModelListOptions{
			Pagination: dbmodels.Pagination{
				OrderBy:   apiutils.Of("artifacts.overall_average.double_value"),
				SortOrder: apiutils.Of("ASC"),
				PageSize:  apiutils.Of(int32(3)),
			},
		}

		pageAsc, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, pageAsc)

		// Just verify that ASC pagination works and generates tokens when there are more results
		if len(pageAsc.Items) == 3 {
			assert.NotEmpty(t, pageAsc.NextPageToken, "Should have next page token in ASC order when page is full")
		}
	})

	t.Run("TestNameOrdering", func(t *testing.T) {
		// Create test models with specific names for ordering
		testModels := []string{
			"zebra-model",
			"alpha-model",
			"beta-model",
			"gamma-model",
			"delta-model",
		}

		var savedModels []models.CatalogModel
		for _, name := range testModels {
			catalogModel := &models.CatalogModelImpl{
				Attributes: &models.CatalogModelAttributes{
					Name:       apiutils.Of(name),
					ExternalID: apiutils.Of(name + "-ext"),
				},
			}

			savedModel, err := repo.Save(catalogModel)
			require.NoError(t, err)
			savedModels = append(savedModels, savedModel)
		}

		// Test NAME ordering ASC
		listOptions := models.CatalogModelListOptions{
			Pagination: dbmodels.Pagination{
				OrderBy:   apiutils.Of("NAME"),
				SortOrder: apiutils.Of("ASC"),
			},
		}
		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Extract our test model names from results
		var foundNames []string
		for _, model := range result.Items {
			name := *model.GetAttributes().Name
			if name == "zebra-model" || name == "alpha-model" || name == "beta-model" ||
				name == "gamma-model" || name == "delta-model" {
				foundNames = append(foundNames, name)
			}
		}

		// Verify we found all our test models
		require.GreaterOrEqual(t, len(foundNames), 5, "Should find all test models")

		// Verify ASC ordering: alpha < beta < delta < gamma < zebra
		alphaIdx := findIndex(foundNames, "alpha-model")
		betaIdx := findIndex(foundNames, "beta-model")
		deltaIdx := findIndex(foundNames, "delta-model")
		gammaIdx := findIndex(foundNames, "gamma-model")
		zebraIdx := findIndex(foundNames, "zebra-model")

		require.NotEqual(t, -1, alphaIdx, "alpha-model not found")
		require.NotEqual(t, -1, betaIdx, "beta-model not found")
		require.NotEqual(t, -1, deltaIdx, "delta-model not found")
		require.NotEqual(t, -1, gammaIdx, "gamma-model not found")
		require.NotEqual(t, -1, zebraIdx, "zebra-model not found")

		assert.Less(t, alphaIdx, betaIdx, "alpha should come before beta in ASC")
		assert.Less(t, betaIdx, deltaIdx, "beta should come before delta in ASC")
		assert.Less(t, deltaIdx, gammaIdx, "delta should come before gamma in ASC")
		assert.Less(t, gammaIdx, zebraIdx, "gamma should come before zebra in ASC")

		// Test NAME ordering DESC
		listOptions = models.CatalogModelListOptions{
			Pagination: dbmodels.Pagination{
				OrderBy:   apiutils.Of("NAME"),
				SortOrder: apiutils.Of("DESC"),
			},
		}
		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Extract our test model names from DESC results
		foundNames = []string{}
		for _, model := range result.Items {
			name := *model.GetAttributes().Name
			if name == "zebra-model" || name == "alpha-model" || name == "beta-model" ||
				name == "gamma-model" || name == "delta-model" {
				foundNames = append(foundNames, name)
			}
		}

		// Verify DESC ordering: zebra > gamma > delta > beta > alpha
		alphaIdxDesc := findIndex(foundNames, "alpha-model")
		betaIdxDesc := findIndex(foundNames, "beta-model")
		deltaIdxDesc := findIndex(foundNames, "delta-model")
		gammaIdxDesc := findIndex(foundNames, "gamma-model")
		zebraIdxDesc := findIndex(foundNames, "zebra-model")

		assert.Less(t, zebraIdxDesc, gammaIdxDesc, "zebra should come before gamma in DESC")
		assert.Less(t, gammaIdxDesc, deltaIdxDesc, "gamma should come before delta in DESC")
		assert.Less(t, deltaIdxDesc, betaIdxDesc, "delta should come before beta in DESC")
		assert.Less(t, betaIdxDesc, alphaIdxDesc, "beta should come before alpha in DESC")
	})

	t.Run("TestNameOrderingPagination", func(t *testing.T) {
		// Create models with sequential names for pagination testing
		testModels := []string{
			"page-test-model-01",
			"page-test-model-02",
			"page-test-model-03",
			"page-test-model-04",
			"page-test-model-05",
		}

		for _, name := range testModels {
			catalogModel := &models.CatalogModelImpl{
				Attributes: &models.CatalogModelAttributes{
					Name:       apiutils.Of(name),
					ExternalID: apiutils.Of(name + "-ext"),
				},
			}

			_, err := repo.Save(catalogModel)
			require.NoError(t, err)
		}

		// Test pagination with NAME ordering
		listOptions := models.CatalogModelListOptions{
			Pagination: dbmodels.Pagination{
				OrderBy:   apiutils.Of("NAME"),
				SortOrder: apiutils.Of("ASC"),
				PageSize:  apiutils.Of(int32(2)),
			},
		}

		// Collect all our test models across pages
		var allPaginatedModels []string
		var pageCount int
		currentToken := (*string)(nil)

		for {
			pageCount++
			if currentToken != nil {
				listOptions.Pagination.NextPageToken = currentToken
			}

			page, err := repo.List(listOptions)
			require.NoError(t, err)
			require.NotNil(t, page)
			assert.Equal(t, int32(2), page.PageSize)

			// Filter to only include our test models
			for _, model := range page.Items {
				name := *model.GetAttributes().Name
				if strings.HasPrefix(name, "page-test-model-") {
					allPaginatedModels = append(allPaginatedModels, name)
				}
			}

			// Stop if no more pages or we've collected all our test models
			if page.NextPageToken == "" || len(allPaginatedModels) >= 5 {
				if page.NextPageToken == "" {
					t.Logf("NAME pagination completed in %d pages", pageCount)
				}
				break
			}
			currentToken = &page.NextPageToken

			// Safety check to prevent infinite loop
			if pageCount > 10 {
				t.Fatal("Too many pages, might be an infinite loop")
			}
		}

		// Verify we collected all our test models
		assert.GreaterOrEqual(t, len(allPaginatedModels), 5, "Should have found all page-test models")

		// Verify ordering is maintained across pages
		expectedOrder := []string{
			"page-test-model-01",
			"page-test-model-02",
			"page-test-model-03",
			"page-test-model-04",
			"page-test-model-05",
		}

		// Verify our test models appear in correct order
		lastIndex := -1
		for _, expectedModel := range expectedOrder {
			foundIndex := findIndex(allPaginatedModels, expectedModel)
			assert.NotEqual(t, -1, foundIndex, "Should find model %s", expectedModel)
			if foundIndex != -1 {
				assert.Greater(t, foundIndex, lastIndex, "Model %s should appear after previous models", expectedModel)
				lastIndex = foundIndex
			}
		}
	})

	t.Run("TestDeleteBySource", func(t *testing.T) {
		// Setup: Create models with different source IDs
		sourceID1 := "test_source_1"
		sourceID2 := "test_source_2"

		model1 := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name: apiutils.Of("model-source-1"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: &sourceID1,
				},
			},
		}

		model2 := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name: apiutils.Of("model-source-1-second"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: &sourceID1,
				},
			},
		}

		model3 := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name: apiutils.Of("model-source-2"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "source_id",
					StringValue: &sourceID2,
				},
			},
		}

		// Save all models
		saved1, err := repo.Save(model1)
		require.NoError(t, err)
		saved2, err := repo.Save(model2)
		require.NoError(t, err)
		saved3, err := repo.Save(model3)
		require.NoError(t, err)

		// Delete by source_id
		err = repo.DeleteBySource(sourceID1)
		require.NoError(t, err)

		// Verify models from source1 are deleted
		_, err = repo.GetByID(*saved1.GetID())
		assert.Error(t, err)
		assert.True(t, errors.Is(err, service.ErrCatalogModelNotFound))

		_, err = repo.GetByID(*saved2.GetID())
		assert.Error(t, err)
		assert.True(t, errors.Is(err, service.ErrCatalogModelNotFound))

		// Verify model from source2 still exists
		retrieved, err := repo.GetByID(*saved3.GetID())
		require.NoError(t, err)
		assert.Equal(t, "model-source-2", *retrieved.GetAttributes().Name)
	})

	t.Run("TestDeleteByID", func(t *testing.T) {
		// Setup: Create a model
		model := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name: apiutils.Of("model-to-delete"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Model that will be deleted"),
				},
			},
		}

		saved, err := repo.Save(model)
		require.NoError(t, err)

		// Delete by ID
		err = repo.DeleteByID(*saved.GetID())
		require.NoError(t, err)

		// Verify model is deleted
		_, err = repo.GetByID(*saved.GetID())
		assert.Error(t, err)
		assert.True(t, errors.Is(err, service.ErrCatalogModelNotFound))
	})

	t.Run("TestDeleteBySourceNonExistent", func(t *testing.T) {
		// Test deleting by non-existent source - should not error
		err := repo.DeleteBySource("non-existent-source")
		require.NoError(t, err)
	})

	t.Run("TestDeleteByIDNonExistent", func(t *testing.T) {
		// Test deleting non-existent ID - should return NotFoundError
		err := repo.DeleteByID(999999)
		require.Error(t, err)
		assert.True(t, errors.Is(err, service.ErrCatalogModelNotFound))
	})

	t.Run("TestGetDistinctSourceIDs", func(t *testing.T) {
		// Get initial count of source IDs
		initialSourceIDs, err := repo.GetDistinctSourceIDs()
		assert.NoError(t, err)
		initialCount := len(initialSourceIDs)

		// Create test data with different source_ids (use unique prefixes to avoid collision with other tests)
		testSourceID1 := "test-distinct-source-1"
		testSourceID2 := "test-distinct-source-2"

		model1 := createTestCatalogModelWithSourceID(t, testSourceID1)
		model2 := createTestCatalogModelWithSourceID(t, testSourceID2)
		model3 := createTestCatalogModelWithSourceID(t, testSourceID1) // duplicate

		_, err = repo.Save(model1)
		assert.NoError(t, err)
		_, err = repo.Save(model2)
		assert.NoError(t, err)
		_, err = repo.Save(model3)
		assert.NoError(t, err)

		// Test distinct source_ids - should have 2 new source IDs added
		sourceIDs, err := repo.GetDistinctSourceIDs()
		assert.NoError(t, err)
		assert.Len(t, sourceIDs, initialCount+2, "Should have exactly 2 new distinct source IDs")
		assert.Contains(t, sourceIDs, testSourceID1)
		assert.Contains(t, sourceIDs, testSourceID2)
	})
}

func createTestCatalogModelWithSourceID(t *testing.T, sourceID string) models.CatalogModel {
	model := &models.CatalogModelImpl{
		Attributes: &models.CatalogModelAttributes{
			Name: apiutils.Of(fmt.Sprintf("test-model-%s", sourceID)),
		},
	}

	// Add source_id as a property
	properties := []dbmodels.Properties{
		{
			Name:        "source_id",
			StringValue: &sourceID,
		},
	}
	model.Properties = &properties

	return model
}

// Helper function to get or create CatalogModel type ID
func getCatalogModelTypeID(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogModelTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogModel type")
	}

	return typeRecord.ID
}

// Helper function to find index of string in slice
func findIndex(slice []string, target string) int {
	for i, item := range slice {
		if item == target {
			return i
		}
	}
	return -1
}

func TestCatalogModelRepository_TimestampPreservation(t *testing.T) {
	sharedDB, cleanup := testutils.SetupPostgresWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	typeID := getCatalogModelTypeID(t, sharedDB)
	repo := service.NewCatalogModelRepository(sharedDB, typeID)

	t.Run("Preserve historical timestamps from YAML catalog", func(t *testing.T) {
		// Simulate loading a model from YAML with historical timestamps
		historicalCreateTime := int64(1739776988000) // From YAML example
		historicalUpdateTime := int64(1746720264000) // From YAML example

		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:                     apiutils.Of("yaml-loaded-model"),
				ExternalID:               apiutils.Of("yaml-model-123"),
				CreateTimeSinceEpoch:     &historicalCreateTime,
				LastUpdateTimeSinceEpoch: &historicalUpdateTime,
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Model loaded from YAML"),
				},
			},
		}

		// Save the model - timestamps should be preserved
		saved, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())

		// Verify historical timestamps were preserved
		savedAttrs := saved.GetAttributes()
		require.NotNil(t, savedAttrs.CreateTimeSinceEpoch)
		require.NotNil(t, savedAttrs.LastUpdateTimeSinceEpoch)
		assert.Equal(t, historicalCreateTime, *savedAttrs.CreateTimeSinceEpoch,
			"CreateTimeSinceEpoch should be preserved from YAML")
		assert.Equal(t, historicalUpdateTime, *savedAttrs.LastUpdateTimeSinceEpoch,
			"LastUpdateTimeSinceEpoch should be preserved from YAML")

		// Reload from database to verify persistence
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		retrievedAttrs := retrieved.GetAttributes()
		assert.Equal(t, historicalCreateTime, *retrievedAttrs.CreateTimeSinceEpoch,
			"CreateTimeSinceEpoch should persist in database")
		assert.Equal(t, historicalUpdateTime, *retrievedAttrs.LastUpdateTimeSinceEpoch,
			"LastUpdateTimeSinceEpoch should persist in database")
	})

	t.Run("Auto-generate timestamps for API-created models", func(t *testing.T) {
		// Model created via API without timestamps
		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("api-created-model"),
				ExternalID: apiutils.Of("api-model-456"),
				// No timestamps set - should be auto-generated
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Model created via API"),
				},
			},
		}

		// Save the model - timestamps should be auto-generated
		saved, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, saved)

		// Verify timestamps were auto-generated (non-zero)
		savedAttrs := saved.GetAttributes()
		require.NotNil(t, savedAttrs.CreateTimeSinceEpoch)
		require.NotNil(t, savedAttrs.LastUpdateTimeSinceEpoch)
		assert.Greater(t, *savedAttrs.CreateTimeSinceEpoch, int64(0),
			"CreateTimeSinceEpoch should be auto-generated")
		assert.Greater(t, *savedAttrs.LastUpdateTimeSinceEpoch, int64(0),
			"LastUpdateTimeSinceEpoch should be auto-generated")
	})

	t.Run("Update existing model from YAML preserves CreateTime", func(t *testing.T) {
		// First save: Create model with historical timestamps
		historicalCreateTime := int64(1739776988000)
		firstUpdateTime := int64(1746720264000)

		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:                     apiutils.Of("yaml-updated-model"),
				ExternalID:               apiutils.Of("yaml-updated-123"),
				CreateTimeSinceEpoch:     &historicalCreateTime,
				LastUpdateTimeSinceEpoch: &firstUpdateTime,
			},
		}

		saved, err := repo.Save(catalogModel)
		require.NoError(t, err)
		savedID := saved.GetID()

		// Second save: Update the model with new LastUpdateTime (simulating catalog reload)
		newerUpdateTime := int64(1750000000000)
		catalogModel.ID = savedID
		catalogModel.GetAttributes().LastUpdateTimeSinceEpoch = &newerUpdateTime

		updated, err := repo.Save(catalogModel)
		require.NoError(t, err)

		// Verify CreateTime is preserved but LastUpdateTime is updated
		updatedAttrs := updated.GetAttributes()
		assert.Equal(t, historicalCreateTime, *updatedAttrs.CreateTimeSinceEpoch,
			"CreateTimeSinceEpoch should be preserved on update")
		assert.Equal(t, newerUpdateTime, *updatedAttrs.LastUpdateTimeSinceEpoch,
			"LastUpdateTimeSinceEpoch should be updated")
	})
}
