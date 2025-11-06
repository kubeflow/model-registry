package service_test

import (
	"fmt"
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/internal/apiutils"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCatalogArtifactRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupPostgresWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	// Get the catalog artifact type IDs
	modelArtifactTypeID := getCatalogModelArtifactTypeID(t, sharedDB)
	metricsArtifactTypeID := getCatalogMetricsArtifactTypeID(t, sharedDB)

	// Create unified artifact repository with both types
	artifactTypeMap := map[string]int32{
		service.CatalogModelArtifactTypeName:   modelArtifactTypeID,
		service.CatalogMetricsArtifactTypeName: metricsArtifactTypeID,
	}
	repo := service.NewCatalogArtifactRepository(sharedDB, artifactTypeMap)

	// Also get CatalogModel type ID for creating parent entities
	catalogModelTypeID := getCatalogModelTypeID(t, sharedDB)
	catalogModelRepo := service.NewCatalogModelRepository(sharedDB, catalogModelTypeID)
	modelArtifactRepo := service.NewCatalogModelArtifactRepository(sharedDB, modelArtifactTypeID)
	metricsArtifactRepo := service.NewCatalogMetricsArtifactRepository(sharedDB, metricsArtifactTypeID)

	// Create shared test data
	catalogModel := &models.CatalogModelImpl{
		TypeID: apiutils.Of(int32(catalogModelTypeID)),
		Attributes: &models.CatalogModelAttributes{
			Name:       apiutils.Of("test-catalog-model-for-artifacts"),
			ExternalID: apiutils.Of("catalog-model-artifacts-ext-123"),
		},
	}
	savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
	require.NoError(t, err)

	t.Run("GetByID_ModelArtifact", func(t *testing.T) {
		// Create a model artifact using the specific repository
		modelArtifact := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:         apiutils.Of("test-model-artifact-getbyid"),
				ExternalID:   apiutils.Of("model-art-getbyid-ext-123"),
				URI:          apiutils.Of("s3://test-bucket/model.bin"),
				ArtifactType: apiutils.Of(models.CatalogModelArtifactType),
			},
		}
		savedModelArtifact, err := modelArtifactRepo.Save(modelArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)

		// Retrieve using unified repository
		retrieved, err := repo.GetByID(*savedModelArtifact.GetID())
		require.NoError(t, err)

		// Verify it's a model artifact
		assert.NotNil(t, retrieved.CatalogModelArtifact)
		assert.Nil(t, retrieved.CatalogMetricsArtifact)
		assert.Equal(t, "test-model-artifact-getbyid", *retrieved.CatalogModelArtifact.GetAttributes().Name)
		assert.Equal(t, "model-art-getbyid-ext-123", *retrieved.CatalogModelArtifact.GetAttributes().ExternalID)
		assert.Equal(t, "s3://test-bucket/model.bin", *retrieved.CatalogModelArtifact.GetAttributes().URI)
	})

	t.Run("GetByID_MetricsArtifact", func(t *testing.T) {
		// Create a metrics artifact using the specific repository
		metricsArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(metricsArtifactTypeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:         apiutils.Of("test-metrics-artifact-getbyid"),
				ExternalID:   apiutils.Of("metrics-art-getbyid-ext-123"),
				MetricsType:  models.MetricsTypeAccuracy,
				ArtifactType: apiutils.Of("metrics-artifact"),
			},
		}
		savedMetricsArtifact, err := metricsArtifactRepo.Save(metricsArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)

		// Retrieve using unified repository
		retrieved, err := repo.GetByID(*savedMetricsArtifact.GetID())
		require.NoError(t, err)

		// Verify it's a metrics artifact
		assert.Nil(t, retrieved.CatalogModelArtifact)
		assert.NotNil(t, retrieved.CatalogMetricsArtifact)
		assert.Equal(t, "test-metrics-artifact-getbyid", *retrieved.CatalogMetricsArtifact.GetAttributes().Name)
		assert.Equal(t, "metrics-art-getbyid-ext-123", *retrieved.CatalogMetricsArtifact.GetAttributes().ExternalID)
		assert.Equal(t, models.MetricsTypeAccuracy, retrieved.CatalogMetricsArtifact.GetAttributes().MetricsType)
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		nonExistentID := int32(99999)
		_, err := repo.GetByID(nonExistentID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "catalog artifact by id not found")
	})

	t.Run("List_AllArtifacts", func(t *testing.T) {
		// Create test artifacts of both types
		modelArtifact1 := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:         apiutils.Of("test-model-artifact-list-1"),
				ExternalID:   apiutils.Of("model-list-1-ext"),
				URI:          apiutils.Of("s3://test/model1.bin"),
				ArtifactType: apiutils.Of(models.CatalogModelArtifactType),
			},
		}

		modelArtifact2 := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:         apiutils.Of("test-model-artifact-list-2"),
				ExternalID:   apiutils.Of("model-list-2-ext"),
				URI:          apiutils.Of("s3://test/model2.bin"),
				ArtifactType: apiutils.Of(models.CatalogModelArtifactType),
			},
		}

		metricsArtifact1 := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(metricsArtifactTypeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:         apiutils.Of("test-metrics-artifact-list-1"),
				ExternalID:   apiutils.Of("metrics-list-1-ext"),
				MetricsType:  models.MetricsTypeAccuracy,
				ArtifactType: apiutils.Of("metrics-artifact"),
			},
		}

		// Save artifacts
		savedModelArt1, err := modelArtifactRepo.Save(modelArtifact1, savedCatalogModel.GetID())
		require.NoError(t, err)
		savedModelArt2, err := modelArtifactRepo.Save(modelArtifact2, savedCatalogModel.GetID())
		require.NoError(t, err)
		savedMetricsArt1, err := metricsArtifactRepo.Save(metricsArtifact1, savedCatalogModel.GetID())
		require.NoError(t, err)

		// List all artifacts for the parent resource
		listOptions := models.CatalogArtifactListOptions{
			ParentResourceID: savedCatalogModel.GetID(),
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Should return all 3 artifacts (2 model + 1 metrics)
		assert.GreaterOrEqual(t, len(result.Items), 3, "Should return at least the 3 artifacts we created")

		// Verify we got both types
		var modelArtifactCount, metricsArtifactCount int
		artifactIDs := make(map[int32]bool)

		for _, artifact := range result.Items {
			if artifact.CatalogModelArtifact != nil {
				modelArtifactCount++
				artifactIDs[*artifact.CatalogModelArtifact.GetID()] = true
			} else if artifact.CatalogMetricsArtifact != nil {
				metricsArtifactCount++
				artifactIDs[*artifact.CatalogMetricsArtifact.GetID()] = true
			}
		}

		assert.GreaterOrEqual(t, modelArtifactCount, 2, "Should have at least 2 model artifacts")
		assert.GreaterOrEqual(t, metricsArtifactCount, 1, "Should have at least 1 metrics artifact")

		// Verify our specific artifacts are in the results
		assert.True(t, artifactIDs[*savedModelArt1.GetID()], "Should contain first model artifact")
		assert.True(t, artifactIDs[*savedModelArt2.GetID()], "Should contain second model artifact")
		assert.True(t, artifactIDs[*savedMetricsArt1.GetID()], "Should contain metrics artifact")
	})

	t.Run("List_FilterByArtifactType_ModelArtifact", func(t *testing.T) {
		// Filter by model artifact type only
		artifactType := "model-artifact"
		listOptions := models.CatalogArtifactListOptions{
			ParentResourceID: savedCatalogModel.GetID(),
			ArtifactType:     &artifactType,
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// All results should be model artifacts
		for _, artifact := range result.Items {
			assert.NotNil(t, artifact.CatalogModelArtifact, "Should only return model artifacts")
			assert.Nil(t, artifact.CatalogMetricsArtifact, "Should not return metrics artifacts")
		}
	})

	t.Run("List_FilterByArtifactType_MetricsArtifact", func(t *testing.T) {
		// Filter by metrics artifact type only
		artifactType := "metrics-artifact"
		listOptions := models.CatalogArtifactListOptions{
			ParentResourceID: savedCatalogModel.GetID(),
			ArtifactType:     &artifactType,
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// All results should be metrics artifacts
		for _, artifact := range result.Items {
			assert.Nil(t, artifact.CatalogModelArtifact, "Should not return model artifacts")
			assert.NotNil(t, artifact.CatalogMetricsArtifact, "Should only return metrics artifacts")
		}
	})

	t.Run("List_FilterByExternalID", func(t *testing.T) {
		// Create artifact with specific external ID for filtering
		testArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(metricsArtifactTypeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:         apiutils.Of("external-id-filter-test"),
				ExternalID:   apiutils.Of("unique-external-id-123"),
				MetricsType:  models.MetricsTypePerformance,
				ArtifactType: apiutils.Of("metrics-artifact"),
			},
		}
		savedArtifact, err := metricsArtifactRepo.Save(testArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)

		// Filter by external ID
		externalID := "unique-external-id-123"
		listOptions := models.CatalogArtifactListOptions{
			ExternalID: &externalID,
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 1, "Should find exactly one artifact with the external ID")

		// Verify it's the correct artifact
		artifact := result.Items[0]
		assert.NotNil(t, artifact.CatalogMetricsArtifact)
		assert.Equal(t, *savedArtifact.GetID(), *artifact.CatalogMetricsArtifact.GetID())
		assert.Equal(t, "unique-external-id-123", *artifact.CatalogMetricsArtifact.GetAttributes().ExternalID)
	})

	t.Run("List_WithPagination", func(t *testing.T) {
		// Create multiple artifacts for pagination testing
		for i := 0; i < 5; i++ {
			artifact := &models.CatalogModelArtifactImpl{
				TypeID: apiutils.Of(int32(modelArtifactTypeID)),
				Attributes: &models.CatalogModelArtifactAttributes{
					Name:         apiutils.Of(fmt.Sprintf("pagination-test-%d", i)),
					ExternalID:   apiutils.Of(fmt.Sprintf("pagination-ext-%d", i)),
					URI:          apiutils.Of(fmt.Sprintf("s3://test/pagination-%d.bin", i)),
					ArtifactType: apiutils.Of(models.CatalogModelArtifactType),
				},
			}
			_, err := modelArtifactRepo.Save(artifact, savedCatalogModel.GetID())
			require.NoError(t, err)
		}

		// Test pagination
		pageSize := int32(3)
		listOptions := models.CatalogArtifactListOptions{
			ParentResourceID: savedCatalogModel.GetID(),
			Pagination: dbmodels.Pagination{
				PageSize: &pageSize,
				OrderBy:  apiutils.Of("ID"),
			},
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.LessOrEqual(t, len(result.Items), 3, "Should respect page size limit")
		assert.GreaterOrEqual(t, len(result.Items), 1, "Should return at least one item")
	})

	t.Run("List_InvalidArtifactType", func(t *testing.T) {
		// Test with invalid artifact type
		invalidType := "invalid-artifact-type"
		listOptions := models.CatalogArtifactListOptions{
			ParentResourceID: savedCatalogModel.GetID(),
			ArtifactType:     &invalidType,
		}

		_, err := repo.List(listOptions)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid catalog artifact type")
		assert.Contains(t, err.Error(), "invalid-artifact-type")
	})

	t.Run("List_WithCustomProperties", func(t *testing.T) {
		// Create artifacts with custom properties
		customProps := []dbmodels.Properties{
			{
				Name:        "custom_prop_1",
				StringValue: apiutils.Of("custom_value_1"),
			},
			{
				Name:        "custom_prop_2",
				StringValue: apiutils.Of("custom_value_2"),
			},
		}

		artifactWithCustomProps := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:         apiutils.Of("artifact-with-custom-props"),
				ExternalID:   apiutils.Of("custom-props-ext"),
				URI:          apiutils.Of("s3://test/custom-props.bin"),
				ArtifactType: apiutils.Of(models.CatalogModelArtifactType),
			},
			CustomProperties: &customProps,
		}

		savedArtifact, err := modelArtifactRepo.Save(artifactWithCustomProps, savedCatalogModel.GetID())
		require.NoError(t, err)

		// Retrieve using unified repository
		retrieved, err := repo.GetByID(*savedArtifact.GetID())
		require.NoError(t, err)

		// Verify custom properties are preserved
		require.NotNil(t, retrieved.CatalogModelArtifact)
		assert.NotNil(t, retrieved.CatalogModelArtifact.GetCustomProperties())

		customPropsMap := make(map[string]string)
		for _, prop := range *retrieved.CatalogModelArtifact.GetCustomProperties() {
			if prop.StringValue != nil {
				customPropsMap[prop.Name] = *prop.StringValue
			}
		}

		assert.Equal(t, "custom_value_1", customPropsMap["custom_prop_1"])
		assert.Equal(t, "custom_value_2", customPropsMap["custom_prop_2"])
	})

	t.Run("MappingErrors", func(t *testing.T) {
		// Test error handling for invalid type mapping
		// This would typically happen if there's data inconsistency in the database

		// We can't easily test this without directly manipulating the database
		// but we can test the GetByID with an artifact that has an unknown type
		// by temporarily modifying the repository's type mapping

		// Create a repository with incomplete type mapping
		incompleteTypeMap := map[string]int32{
			service.CatalogModelArtifactTypeName: modelArtifactTypeID,
			// Missing CatalogMetricsArtifactTypeName intentionally
		}
		incompleteRepo := service.NewCatalogArtifactRepository(sharedDB, incompleteTypeMap)

		// Create a metrics artifact first using the complete repo
		metricsArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(metricsArtifactTypeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:         apiutils.Of("test-mapping-error"),
				ExternalID:   apiutils.Of("mapping-error-ext"),
				MetricsType:  models.MetricsTypeAccuracy,
				ArtifactType: apiutils.Of("metrics-artifact"),
			},
		}
		savedMetricsArtifact, err := metricsArtifactRepo.Save(metricsArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)

		// Try to retrieve using incomplete repo - should get mapping error
		_, err = incompleteRepo.GetByID(*savedMetricsArtifact.GetID())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid catalog artifact type")
	})

	t.Run("TestNameOrdering", func(t *testing.T) {
		// Create a new model for this test
		testModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-model-for-name-ordering"),
				ExternalID: apiutils.Of("test-model-name-ordering-ext"),
			},
		}
		savedTestModel, err := catalogModelRepo.Save(testModel)
		require.NoError(t, err)

		// Create artifacts with various names (including null)
		testArtifacts := []struct {
			name *string
			desc string
		}{
			{apiutils.Of("zebra-artifact"), "zebra"},
			{apiutils.Of("alpha-artifact"), "alpha"},
			{apiutils.Of("beta-artifact"), "beta"},
			{nil, "null-name"}, // Artifact with no name (like real model artifacts)
			{apiutils.Of("gamma-artifact"), "gamma"},
		}

		for _, artifact := range testArtifacts {
			metricsArtifact := &models.CatalogMetricsArtifactImpl{
				TypeID: apiutils.Of(int32(metricsArtifactTypeID)),
				Attributes: &models.CatalogMetricsArtifactAttributes{
					Name:         artifact.name,
					ExternalID:   apiutils.Of(fmt.Sprintf("name-test-%s", artifact.desc)),
					MetricsType:  models.MetricsTypePerformance,
					ArtifactType: apiutils.Of("metrics-artifact"),
				},
			}
			_, err := metricsArtifactRepo.Save(metricsArtifact, savedTestModel.GetID())
			require.NoError(t, err)
		}

		// Test NAME ordering ASC
		listOptions := models.CatalogArtifactListOptions{
			ParentResourceID: savedTestModel.GetID(),
			Pagination: dbmodels.Pagination{
				OrderBy:   apiutils.Of("NAME"),
				SortOrder: apiutils.Of("ASC"),
			},
		}
		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Extract artifact names (including nulls)
		var foundArtifacts []struct {
			name *string
			desc string
		}
		for _, artifact := range result.Items {
			if artifact.CatalogMetricsArtifact != nil {
				name := artifact.CatalogMetricsArtifact.GetAttributes().Name
				foundArtifacts = append(foundArtifacts, struct {
					name *string
					desc string
				}{name, fmt.Sprintf("%v", name)})
			}
		}

		require.GreaterOrEqual(t, len(foundArtifacts), 5, "Should find all test artifacts")

		// Find positions of named artifacts
		var alphaIdx, betaIdx, gammaIdx, zebraIdx, nullIdx int = -1, -1, -1, -1, -1
		for i, artifact := range foundArtifacts {
			if artifact.name != nil {
				switch *artifact.name {
				case "alpha-artifact":
					alphaIdx = i
				case "beta-artifact":
					betaIdx = i
				case "gamma-artifact":
					gammaIdx = i
				case "zebra-artifact":
					zebraIdx = i
				}
			} else {
				nullIdx = i
			}
		}

		// Verify ASC ordering: alpha < beta < gamma < zebra, and null at the end
		require.NotEqual(t, -1, alphaIdx, "alpha-artifact not found")
		require.NotEqual(t, -1, betaIdx, "beta-artifact not found")
		require.NotEqual(t, -1, gammaIdx, "gamma-artifact not found")
		require.NotEqual(t, -1, zebraIdx, "zebra-artifact not found")
		require.NotEqual(t, -1, nullIdx, "null-name artifact not found")

		assert.Less(t, alphaIdx, betaIdx, "alpha should come before beta in ASC")
		assert.Less(t, betaIdx, gammaIdx, "beta should come before gamma in ASC")
		assert.Less(t, gammaIdx, zebraIdx, "gamma should come before zebra in ASC")
		assert.Less(t, zebraIdx, nullIdx, "named artifacts should come before null in ASC")

		// Test NAME ordering DESC
		listOptions = models.CatalogArtifactListOptions{
			ParentResourceID: savedTestModel.GetID(),
			Pagination: dbmodels.Pagination{
				OrderBy:   apiutils.Of("NAME"),
				SortOrder: apiutils.Of("DESC"),
			},
		}
		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Extract artifact names from DESC results
		foundArtifacts = []struct {
			name *string
			desc string
		}{}
		for _, artifact := range result.Items {
			if artifact.CatalogMetricsArtifact != nil {
				name := artifact.CatalogMetricsArtifact.GetAttributes().Name
				foundArtifacts = append(foundArtifacts, struct {
					name *string
					desc string
				}{name, fmt.Sprintf("%v", name)})
			}
		}

		// Find positions in DESC order
		alphaIdx, betaIdx, gammaIdx, zebraIdx, nullIdx = -1, -1, -1, -1, -1
		for i, artifact := range foundArtifacts {
			if artifact.name != nil {
				switch *artifact.name {
				case "alpha-artifact":
					alphaIdx = i
				case "beta-artifact":
					betaIdx = i
				case "gamma-artifact":
					gammaIdx = i
				case "zebra-artifact":
					zebraIdx = i
				}
			} else {
				nullIdx = i
			}
		}

		// Verify DESC ordering: In SQL DESC, NULL comes first, then zebra > gamma > beta > alpha
		assert.Less(t, nullIdx, zebraIdx, "null should come first in DESC (SQL default behavior)")
		assert.Less(t, zebraIdx, gammaIdx, "zebra should come before gamma in DESC")
		assert.Less(t, gammaIdx, betaIdx, "gamma should come before beta in DESC")
		assert.Less(t, betaIdx, alphaIdx, "beta should come before alpha in DESC")
	})

	t.Run("TestNameOrderingPagination", func(t *testing.T) {
		// Create a new model for this test
		testModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-model-for-name-pagination"),
				ExternalID: apiutils.Of("test-model-name-pagination-ext"),
			},
		}
		savedTestModel, err := catalogModelRepo.Save(testModel)
		require.NoError(t, err)

		// Create artifacts with sequential names for pagination testing
		artifactNames := []string{
			"artifact-alpha",
			"artifact-beta",
			"artifact-gamma",
			"artifact-delta",
			"artifact-epsilon",
		}

		for i, name := range artifactNames {
			metricsArtifact := &models.CatalogMetricsArtifactImpl{
				TypeID: apiutils.Of(int32(metricsArtifactTypeID)),
				Attributes: &models.CatalogMetricsArtifactAttributes{
					Name:         apiutils.Of(name),
					ExternalID:   apiutils.Of(fmt.Sprintf("pagination-test-%d", i)),
					MetricsType:  models.MetricsTypePerformance,
					ArtifactType: apiutils.Of("metrics-artifact"),
				},
			}
			_, err := metricsArtifactRepo.Save(metricsArtifact, savedTestModel.GetID())
			require.NoError(t, err)
		}

		// Test pagination with NAME ordering (ASC)
		pageSize := int32(2)
		listOptions := models.CatalogArtifactListOptions{
			ParentResourceID: savedTestModel.GetID(),
			Pagination: dbmodels.Pagination{
				OrderBy:   apiutils.Of("NAME"),
				SortOrder: apiutils.Of("ASC"),
				PageSize:  &pageSize,
			},
		}

		// First page
		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Filter to only our test artifacts
		var page1Artifacts []string
		for _, artifact := range result.Items {
			if artifact.CatalogMetricsArtifact != nil && artifact.CatalogMetricsArtifact.GetAttributes().Name != nil {
				name := *artifact.CatalogMetricsArtifact.GetAttributes().Name
				if name == "artifact-alpha" || name == "artifact-beta" || name == "artifact-gamma" || name == "artifact-delta" || name == "artifact-epsilon" {
					page1Artifacts = append(page1Artifacts, name)
				}
			}
		}

		require.LessOrEqual(t, len(page1Artifacts), 2, "First page should have at most 2 artifacts")
		require.GreaterOrEqual(t, len(page1Artifacts), 1, "First page should have at least 1 artifact")
		assert.NotNil(t, result.NextPageToken, "Should have next page token")

		// Verify first page ordering
		if len(page1Artifacts) >= 2 {
			assert.Less(t, page1Artifacts[0], page1Artifacts[1], "First page should be ordered")
		}

		// Second page
		listOptions.Pagination.NextPageToken = &result.NextPageToken
		result2, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result2)

		var page2Artifacts []string
		for _, artifact := range result2.Items {
			if artifact.CatalogMetricsArtifact != nil && artifact.CatalogMetricsArtifact.GetAttributes().Name != nil {
				name := *artifact.CatalogMetricsArtifact.GetAttributes().Name
				if name == "artifact-alpha" || name == "artifact-beta" || name == "artifact-gamma" || name == "artifact-delta" || name == "artifact-epsilon" {
					page2Artifacts = append(page2Artifacts, name)
				}
			}
		}

		require.GreaterOrEqual(t, len(page2Artifacts), 1, "Second page should have at least 1 artifact")

		// Verify second page ordering
		if len(page2Artifacts) >= 2 {
			assert.Less(t, page2Artifacts[0], page2Artifacts[1], "Second page should be ordered")
		}

		// Verify no overlap between pages
		for _, name1 := range page1Artifacts {
			for _, name2 := range page2Artifacts {
				assert.NotEqual(t, name1, name2, "Pages should not have overlapping artifacts")
			}
		}

		// Verify page 2 comes after page 1
		if len(page1Artifacts) > 0 && len(page2Artifacts) > 0 {
			assert.Less(t, page1Artifacts[len(page1Artifacts)-1], page2Artifacts[0], "Page 2 should continue where page 1 ended")
		}

		// Test DESC pagination
		listOptions = models.CatalogArtifactListOptions{
			ParentResourceID: savedTestModel.GetID(),
			Pagination: dbmodels.Pagination{
				OrderBy:   apiutils.Of("NAME"),
				SortOrder: apiutils.Of("DESC"),
				PageSize:  &pageSize,
			},
		}

		resultDesc, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, resultDesc)

		var pageDescArtifacts []string
		for _, artifact := range resultDesc.Items {
			if artifact.CatalogMetricsArtifact != nil && artifact.CatalogMetricsArtifact.GetAttributes().Name != nil {
				name := *artifact.CatalogMetricsArtifact.GetAttributes().Name
				if name == "artifact-alpha" || name == "artifact-beta" || name == "artifact-gamma" || name == "artifact-delta" || name == "artifact-epsilon" {
					pageDescArtifacts = append(pageDescArtifacts, name)
				}
			}
		}

		require.GreaterOrEqual(t, len(pageDescArtifacts), 1, "DESC first page should have at least 1 artifact")

		// Verify DESC ordering
		if len(pageDescArtifacts) >= 2 {
			assert.Greater(t, pageDescArtifacts[0], pageDescArtifacts[1], "DESC page should be reverse ordered")
		}
	})
}
