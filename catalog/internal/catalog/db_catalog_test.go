package catalog

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/apiutils"
	mr_models "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	os.Exit(testutils.TestMainPostgresHelper(m))
}

func TestDBCatalog(t *testing.T) {
	// Setup test database
	sharedDB, cleanup := testutils.SetupPostgresWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	// Get type IDs
	catalogModelTypeID := getCatalogModelTypeIDForDBTest(t, sharedDB)
	modelArtifactTypeID := getCatalogModelArtifactTypeIDForDBTest(t, sharedDB)
	metricsArtifactTypeID := getCatalogMetricsArtifactTypeIDForDBTest(t, sharedDB)

	// Create repositories
	catalogModelRepo := service.NewCatalogModelRepository(sharedDB, catalogModelTypeID)
	catalogArtifactRepo := service.NewCatalogArtifactRepository(sharedDB, map[string]int32{
		service.CatalogModelArtifactTypeName:   modelArtifactTypeID,
		service.CatalogMetricsArtifactTypeName: metricsArtifactTypeID,
	})
	modelArtifactRepo := service.NewCatalogModelArtifactRepository(sharedDB, modelArtifactTypeID)
	metricsArtifactRepo := service.NewCatalogMetricsArtifactRepository(sharedDB, metricsArtifactTypeID)

	svcs := service.NewServices(
		catalogModelRepo,
		catalogArtifactRepo,
		modelArtifactRepo,
		metricsArtifactRepo,
	)

	// Create DB catalog instance
	dbCatalog := NewDBCatalog(svcs, nil)
	ctx := context.Background()

	t.Run("TestNewDBCatalog", func(t *testing.T) {
		catalog := NewDBCatalog(svcs, nil)
		require.NotNil(t, catalog)

		// Verify it implements the interface
		var _ APIProvider = catalog
	})

	t.Run("TestGetModel_Success", func(t *testing.T) {
		// Create test model
		testModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-get-model"),
				ExternalID: apiutils.Of("test-get-model-ext"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("test-source-id")},
				{Name: "description", StringValue: apiutils.Of("Test model description")},
			},
		}

		savedModel, err := catalogModelRepo.Save(testModel)
		require.NoError(t, err)

		// Test GetModel
		retrievedModel, err := dbCatalog.GetModel(ctx, "test-get-model", "test-source-id")
		require.NoError(t, err)
		require.NotNil(t, retrievedModel)

		assert.Equal(t, "test-get-model", retrievedModel.Name)
		assert.Equal(t, strconv.FormatInt(int64(*savedModel.GetID()), 10), *retrievedModel.Id)
		assert.Equal(t, "test-get-model-ext", *retrievedModel.ExternalId)
		assert.Equal(t, "test-source-id", *retrievedModel.SourceId)
		assert.Equal(t, "Test model description", *retrievedModel.Description)
	})

	t.Run("TestGetModel_NotFound", func(t *testing.T) {
		// Test with non-existent model
		_, err := dbCatalog.GetModel(ctx, "non-existent-model", "test-source-id")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no models found")
		assert.ErrorIs(t, err, api.ErrNotFound)
	})

	t.Run("TestListModels_Success", func(t *testing.T) {
		// Create test models
		sourceIDs := []string{"list-test-source"}

		model1 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("list-test-model-1"),
				ExternalID: apiutils.Of("list-test-1"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("list-test-source")},
				{Name: "description", StringValue: apiutils.Of("First test model")},
			},
		}

		model2 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("list-test-model-2"),
				ExternalID: apiutils.Of("list-test-2"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("list-test-source")},
				{Name: "description", StringValue: apiutils.Of("Second test model")},
			},
		}

		_, err := catalogModelRepo.Save(model1)
		require.NoError(t, err)
		_, err = catalogModelRepo.Save(model2)
		require.NoError(t, err)

		// Test ListModels
		params := ListModelsParams{
			SourceIDs:     sourceIDs,
			PageSize:      10,
			OrderBy:       model.ORDERBYFIELD_CREATE_TIME,
			SortOrder:     model.SORTORDER_ASC,
			NextPageToken: apiutils.Of(""),
		}

		result, err := dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(result.Items), 2, "Should return at least 2 models")
		assert.Equal(t, int32(10), result.PageSize)
		assert.GreaterOrEqual(t, result.Size, int32(2))

		// Verify models are properly mapped
		modelNames := make(map[string]bool)
		for _, model := range result.Items {
			modelNames[model.Name] = true
			// Verify required fields are present
			assert.NotEmpty(t, *model.Id)
			assert.NotEmpty(t, *model.SourceId)
		}

		// Should contain our test models
		foundCount := 0
		if modelNames["list-test-model-1"] {
			foundCount++
		}
		if modelNames["list-test-model-2"] {
			foundCount++
		}
		assert.GreaterOrEqual(t, foundCount, 2, "Should find our test models")
	})

	t.Run("TestListModels_WithPagination", func(t *testing.T) {
		// Test pagination
		sourceIDs := []string{"pagination-test-source"}

		// Create multiple models
		for i := 0; i < 5; i++ {
			model := &models.CatalogModelImpl{
				TypeID: apiutils.Of(int32(catalogModelTypeID)),
				Attributes: &models.CatalogModelAttributes{
					Name:       apiutils.Of(fmt.Sprintf("pagination-test-model-%d", i)),
					ExternalID: apiutils.Of(fmt.Sprintf("pagination-test-%d", i)),
				},
				Properties: &[]mr_models.Properties{
					{Name: "source_id", StringValue: apiutils.Of("pagination-test-source")},
				},
			}
			_, err := catalogModelRepo.Save(model)
			require.NoError(t, err)
		}

		params := ListModelsParams{
			SourceIDs:     sourceIDs,
			PageSize:      3,
			OrderBy:       model.ORDERBYFIELD_CREATE_TIME,
			SortOrder:     model.SORTORDER_ASC,
			NextPageToken: apiutils.Of(""),
		}

		result, err := dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)

		assert.LessOrEqual(t, len(result.Items), 3, "Should respect page size")
		assert.Equal(t, int32(3), result.PageSize)
	})

	t.Run("TestListModels_WithQuery", func(t *testing.T) {
		// Create test models with different properties for query filtering
		sourceIDs := []string{"query-test-source"}

		model1 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("BERT-base-model"),
				ExternalID: apiutils.Of("bert-base-1"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("query-test-source")},
				{Name: "description", StringValue: apiutils.Of("BERT base model for NLP tasks")},
				{Name: "provider", StringValue: apiutils.Of("Hugging Face")},
				{Name: "tasks", StringValue: apiutils.Of(`["text-classification", "question-answering"]`)},
			},
		}

		model2 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("GPT-3.5-turbo"),
				ExternalID: apiutils.Of("gpt-35-turbo-1"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("query-test-source")},
				{Name: "description", StringValue: apiutils.Of("OpenAI GPT model for text generation")},
				{Name: "provider", StringValue: apiutils.Of("OpenAI")},
				{Name: "tasks", StringValue: apiutils.Of(`["text-generation", "conversational"]`)},
			},
		}

		model3 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("ResNet-50-image"),
				ExternalID: apiutils.Of("resnet-50-1"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("query-test-source")},
				{Name: "description", StringValue: apiutils.Of("Deep learning model for image classification")},
				{Name: "provider", StringValue: apiutils.Of("PyTorch")},
				{Name: "tasks", StringValue: apiutils.Of(`["image-classification", "computer-vision"]`)},
			},
		}

		_, err := catalogModelRepo.Save(model1)
		require.NoError(t, err)
		_, err = catalogModelRepo.Save(model2)
		require.NoError(t, err)
		_, err = catalogModelRepo.Save(model3)
		require.NoError(t, err)

		// Test query filtering by name
		params := ListModelsParams{
			Query:         "BERT",
			SourceIDs:     sourceIDs,
			PageSize:      10,
			OrderBy:       model.ORDERBYFIELD_CREATE_TIME,
			SortOrder:     model.SORTORDER_ASC,
			NextPageToken: apiutils.Of(""),
		}

		result, err := dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)

		assert.Equal(t, int32(1), result.Size, "Should return 1 model matching 'BERT'")
		assert.Contains(t, result.Items[0].Name, "BERT", "Should contain BERT model")

		// Test query filtering by description
		params.Query = "NLP"
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)

		assert.Equal(t, int32(1), result.Size, "Should return 1 model with 'NLP' in description")
		assert.Contains(t, result.Items[0].Name, "BERT", "Should contain BERT model")

		// Test query filtering by provider
		params.Query = "OpenAI"
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)

		assert.Equal(t, int32(1), result.Size, "Should return 1 model from 'OpenAI' provider")
		assert.Contains(t, result.Items[0].Name, "GPT", "Should contain GPT model")

		// Test query filtering that should match multiple models
		params.Query = "model"
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, result.Size, int32(3), "Should return at least 3 models matching 'model'")

		// Test query that should return no results
		params.Query = "nonexistent"
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)

		assert.Equal(t, int32(0), result.Size, "Should return 0 models for nonexistent query")

		// Test query filtering by tasks - text-classification
		params.Query = "text-classification"
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)

		assert.Equal(t, int32(1), result.Size, "Should return 1 model with 'text-classification' task")
		assert.Contains(t, result.Items[0].Name, "BERT", "Should contain BERT model")

		// Test query filtering by tasks - image-classification
		params.Query = "image-classification"
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)

		assert.Equal(t, int32(1), result.Size, "Should return 1 model with 'image-classification' task")
		assert.Contains(t, result.Items[0].Name, "ResNet", "Should contain ResNet model")

		// Test query filtering by tasks - conversational
		params.Query = "conversational"
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)

		assert.Equal(t, int32(1), result.Size, "Should return 1 model with 'conversational' task")
		assert.Contains(t, result.Items[0].Name, "GPT", "Should contain GPT model")

		// Test query filtering by tasks - partial match on "classification"
		params.Query = "classification"
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)

		assert.Equal(t, int32(2), result.Size, "Should return 2 models with 'classification' in their tasks")

		// Test query filtering by tasks - computer-vision
		params.Query = "computer-vision"
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)

		assert.Equal(t, int32(1), result.Size, "Should return 1 model with 'computer-vision' task")
		assert.Contains(t, result.Items[0].Name, "ResNet", "Should contain ResNet model")
	})

	t.Run("TestListModels_FilterQuery", func(t *testing.T) {
		// Create test models with diverse properties for filterQuery testing
		sourceIDs := []string{"filterquery-test-source"}

		model1 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("TensorFlow-ResNet50"),
				ExternalID: apiutils.Of("tf-resnet50-1"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("filterquery-test-source")},
				{Name: "description", StringValue: apiutils.Of("Deep learning model for image classification using TensorFlow")},
				{Name: "provider", StringValue: apiutils.Of("Google")},
				{Name: "framework", StringValue: apiutils.Of("TensorFlow")},
				{Name: "tasks", StringValue: apiutils.Of(`["image-classification", "computer-vision"]`)},
				{Name: "accuracy", StringValue: apiutils.Of("0.95")},
			},
		}

		model2 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("PyTorch-BERT"),
				ExternalID: apiutils.Of("pt-bert-1"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("filterquery-test-source")},
				{Name: "description", StringValue: apiutils.Of("BERT model for natural language processing using PyTorch")},
				{Name: "provider", StringValue: apiutils.Of("Hugging Face")},
				{Name: "framework", StringValue: apiutils.Of("PyTorch")},
				{Name: "tasks", StringValue: apiutils.Of(`["text-classification", "question-answering"]`)},
				{Name: "accuracy", StringValue: apiutils.Of("0.92")},
			},
		}

		model3 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("Scikit-learn-LogisticRegression"),
				ExternalID: apiutils.Of("sk-lr-1"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("filterquery-test-source")},
				{Name: "description", StringValue: apiutils.Of("Traditional machine learning model for classification")},
				{Name: "provider", StringValue: apiutils.Of("Scikit-learn")},
				{Name: "framework", StringValue: apiutils.Of("Scikit-learn")},
				{Name: "tasks", StringValue: apiutils.Of(`["classification", "regression"]`)},
				{Name: "accuracy", StringValue: apiutils.Of("0.88")},
			},
		}

		_, err := catalogModelRepo.Save(model1)
		require.NoError(t, err)
		_, err = catalogModelRepo.Save(model2)
		require.NoError(t, err)
		_, err = catalogModelRepo.Save(model3)
		require.NoError(t, err)

		// Test: Basic name filtering with exact match
		params := ListModelsParams{
			FilterQuery:   "name = \"TensorFlow-ResNet50\"",
			SourceIDs:     sourceIDs,
			PageSize:      10,
			OrderBy:       model.ORDERBYFIELD_NAME,
			SortOrder:     model.SORTORDER_ASC,
			NextPageToken: apiutils.Of(""),
		}

		result, err := dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)
		assert.Equal(t, int32(1), result.Size, "Should return 1 model with exact name match")
		assert.Equal(t, "TensorFlow-ResNet50", result.Items[0].Name)

		// Test: LIKE pattern matching
		params.FilterQuery = "name LIKE \"%Tensor%\""
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)
		assert.Equal(t, int32(1), result.Size, "Should return 1 model with LIKE pattern match")
		assert.Contains(t, result.Items[0].Name, "Tensor")

		// Test: LIKE pattern matching with case sensitivity
		params.FilterQuery = "name ILIKE \"%tensor%\""
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)
		assert.Equal(t, int32(1), result.Size, "Should return 1 model with case-insensitive ILIKE match")
		assert.Contains(t, result.Items[0].Name, "Tensor")

		// Test: OR logic
		params.FilterQuery = "name = \"TensorFlow-ResNet50\" OR name = \"PyTorch-BERT\""
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)
		assert.Equal(t, int32(2), result.Size, "Should return 2 models with OR logic")

		// Verify we got the expected models
		modelNames := make(map[string]bool)
		for _, item := range result.Items {
			modelNames[item.Name] = true
		}
		assert.True(t, modelNames["TensorFlow-ResNet50"], "Should contain TensorFlow model")
		assert.True(t, modelNames["PyTorch-BERT"], "Should contain PyTorch model")

		// Test: AND logic
		params.FilterQuery = "name LIKE \"%Tensor%\" AND name LIKE \"%ResNet%\""
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)
		assert.Equal(t, int32(1), result.Size, "Should return 1 model with AND logic")
		assert.Equal(t, "TensorFlow-ResNet50", result.Items[0].Name)

		// Test: Custom property filtering
		params.FilterQuery = "framework.string_value = \"PyTorch\""
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)
		assert.Equal(t, int32(1), result.Size, "Should return 1 model with PyTorch framework")
		assert.Equal(t, "PyTorch-BERT", result.Items[0].Name)

		// Test: Custom property filtering with LIKE
		params.FilterQuery = "provider.string_value LIKE \"%Google%\""
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)
		assert.Equal(t, int32(1), result.Size, "Should return 1 model with Google provider")
		assert.Equal(t, "TensorFlow-ResNet50", result.Items[0].Name)

		// Test: Numeric comparison
		params.FilterQuery = "accuracy.string_value > \"0.90\""
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)
		assert.Equal(t, int32(2), result.Size, "Should return 2 models with accuracy > 0.90")

		// Verify we got the expected models (TensorFlow and PyTorch)
		modelNames = make(map[string]bool)
		for _, item := range result.Items {
			modelNames[item.Name] = true
		}
		assert.True(t, modelNames["TensorFlow-ResNet50"], "Should contain TensorFlow model")
		assert.True(t, modelNames["PyTorch-BERT"], "Should contain PyTorch model")

		// Test: Complex query with multiple conditions
		params.FilterQuery = "(framework.string_value = \"TensorFlow\" OR framework.string_value = \"PyTorch\") AND accuracy.string_value > \"0.90\""
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)
		assert.Equal(t, int32(2), result.Size, "Should return 2 models with complex query")

		// Test: No matches
		params.FilterQuery = "name = \"NonExistentModel\""
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)
		assert.Equal(t, int32(0), result.Size, "Should return 0 models for non-existent name")

		// Test: Empty filterQuery should return all models
		params.FilterQuery = ""
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)
		assert.Equal(t, int32(3), result.Size, "Should return all 3 models with empty filterQuery")

		// Test: Combined with regular query parameter
		params.Query = "BERT"
		params.FilterQuery = "framework.string_value = \"PyTorch\""
		result, err = dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)
		assert.Equal(t, int32(1), result.Size, "Should return 1 model matching both query and filterQuery")
		assert.Equal(t, "PyTorch-BERT", result.Items[0].Name)

		// Test: Invalid filterQuery syntax should return error
		params.Query = ""
		params.FilterQuery = "invalid syntax here"
		_, err = dbCatalog.ListModels(ctx, params)
		require.Error(t, err, "Should return error for invalid filterQuery syntax")
		assert.Contains(t, err.Error(), "invalid filter query", "Error should mention invalid filter query")
	})

	t.Run("TestGetArtifacts_Success", func(t *testing.T) {
		// Create test model
		testModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("artifact-test-model"),
				ExternalID: apiutils.Of("artifact-test-model-ext"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("artifact-test-source")},
			},
		}

		savedModel, err := catalogModelRepo.Save(testModel)
		require.NoError(t, err)

		// Create test artifacts
		modelArtifact := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:         apiutils.Of("test-model-artifact"),
				ExternalID:   apiutils.Of("test-model-artifact-ext"),
				URI:          apiutils.Of("s3://test/model.bin"),
				ArtifactType: apiutils.Of(models.CatalogModelArtifactType),
			},
		}

		metricsArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(metricsArtifactTypeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:         apiutils.Of("test-metrics-artifact"),
				ExternalID:   apiutils.Of("test-metrics-artifact-ext"),
				MetricsType:  models.MetricsTypeAccuracy,
				ArtifactType: apiutils.Of("metrics-artifact"),
			},
		}

		savedModelArt, err := modelArtifactRepo.Save(modelArtifact, savedModel.GetID())
		require.NoError(t, err)
		savedMetricsArt, err := metricsArtifactRepo.Save(metricsArtifact, savedModel.GetID())
		require.NoError(t, err)

		// Test GetArtifacts
		params := ListArtifactsParams{
			PageSize:      10,
			OrderBy:       model.ORDERBYFIELD_CREATE_TIME,
			SortOrder:     model.SORTORDER_ASC,
			NextPageToken: apiutils.Of(""),
		}

		result, err := dbCatalog.GetArtifacts(ctx, "artifact-test-model", "artifact-test-source", params)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(result.Items), 2, "Should return at least 2 artifacts")
		assert.Equal(t, int32(10), result.PageSize)

		// Verify both types of artifacts are returned
		var modelArtifactFound, metricsArtifactFound bool
		artifactIDs := make(map[string]bool)

		for _, artifact := range result.Items {
			if artifact.CatalogModelArtifact != nil {
				modelArtifactFound = true
				artifactIDs[*artifact.CatalogModelArtifact.Id] = true
				assert.Equal(t, "model-artifact", artifact.CatalogModelArtifact.ArtifactType)
			}
			if artifact.CatalogMetricsArtifact != nil {
				metricsArtifactFound = true
				artifactIDs[*artifact.CatalogMetricsArtifact.Id] = true
				assert.Equal(t, "metrics-artifact", artifact.CatalogMetricsArtifact.ArtifactType)
			}
		}

		assert.True(t, modelArtifactFound, "Should find model artifact")
		assert.True(t, metricsArtifactFound, "Should find metrics artifact")

		// Verify our specific artifacts are in the results
		modelArtifactIDStr := strconv.FormatInt(int64(*savedModelArt.GetID()), 10)
		metricsArtifactIDStr := strconv.FormatInt(int64(*savedMetricsArt.GetID()), 10)
		assert.True(t, artifactIDs[modelArtifactIDStr], "Should contain our model artifact")
		assert.True(t, artifactIDs[metricsArtifactIDStr], "Should contain our metrics artifact")
	})

	t.Run("TestGetArtifacts_ModelNotFound", func(t *testing.T) {
		// Test with non-existent model
		params := ListArtifactsParams{
			PageSize:      10,
			OrderBy:       model.ORDERBYFIELD_CREATE_TIME,
			SortOrder:     model.SORTORDER_ASC,
			NextPageToken: apiutils.Of(""),
		}

		_, err := dbCatalog.GetArtifacts(ctx, "non-existent-model", "test-source", params)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid model name")
	})

	t.Run("TestGetArtifacts_WithCustomProperties", func(t *testing.T) {
		// Create model
		testModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("custom-props-model"),
				ExternalID: apiutils.Of("custom-props-model-ext"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("custom-props-source")},
			},
		}

		savedModel, err := catalogModelRepo.Save(testModel)
		require.NoError(t, err)

		// Create artifact with custom properties
		customProps := []mr_models.Properties{
			{Name: "custom_prop_1", StringValue: apiutils.Of("value_1")},
			{Name: "custom_prop_2", StringValue: apiutils.Of("value_2")},
		}

		artifactWithProps := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:         apiutils.Of("artifact-with-props"),
				ExternalID:   apiutils.Of("artifact-with-props-ext"),
				URI:          apiutils.Of("s3://test/props.bin"),
				ArtifactType: apiutils.Of(models.CatalogModelArtifactType),
			},
			CustomProperties: &customProps,
		}

		_, err = modelArtifactRepo.Save(artifactWithProps, savedModel.GetID())
		require.NoError(t, err)

		// Get artifacts and verify custom properties
		params := ListArtifactsParams{
			PageSize:      10,
			OrderBy:       model.ORDERBYFIELD_CREATE_TIME,
			SortOrder:     model.SORTORDER_ASC,
			NextPageToken: apiutils.Of(""),
		}

		result, err := dbCatalog.GetArtifacts(ctx, "custom-props-model", "custom-props-source", params)
		require.NoError(t, err)

		// Find our artifact and check custom properties
		found := false
		for _, artifact := range result.Items {
			if artifact.CatalogModelArtifact != nil &&
				artifact.CatalogModelArtifact.Name != nil &&
				*artifact.CatalogModelArtifact.Name == "artifact-with-props" {

				found = true
				assert.NotNil(t, artifact.CatalogModelArtifact.CustomProperties)

				// Verify custom properties are present and properly converted
				customPropsMap := artifact.CatalogModelArtifact.CustomProperties
				assert.Contains(t, customPropsMap, "custom_prop_1")
				assert.Contains(t, customPropsMap, "custom_prop_2")

				// Verify the values are properly converted to MetadataValue
				prop1 := customPropsMap["custom_prop_1"]
				assert.NotNil(t, prop1.MetadataStringValue)
				assert.Equal(t, "value_1", prop1.MetadataStringValue.StringValue)

				break
			}
		}
		assert.True(t, found, "Should find artifact with custom properties")
	})

	t.Run("TestMappingFunctions", func(t *testing.T) {
		t.Run("TestMapCatalogModelToCatalogModel", func(t *testing.T) {
			// Create a catalog model with various properties
			catalogModel := &models.CatalogModelImpl{
				ID:     apiutils.Of(int32(123)),
				TypeID: apiutils.Of(int32(catalogModelTypeID)),
				Attributes: &models.CatalogModelAttributes{
					Name:                     apiutils.Of("mapping-test-model"),
					ExternalID:               apiutils.Of("mapping-test-ext"),
					CreateTimeSinceEpoch:     apiutils.Of(int64(1234567890)),
					LastUpdateTimeSinceEpoch: apiutils.Of(int64(1234567891)),
				},
				Properties: &[]mr_models.Properties{
					{Name: "source_id", StringValue: apiutils.Of("test-source")},
					{Name: "description", StringValue: apiutils.Of("Test description")},
					{Name: "library_name", StringValue: apiutils.Of("pytorch")},
					{Name: "language", StringValue: apiutils.Of("[\"python\", \"go\"]")},
					{Name: "tasks", StringValue: apiutils.Of("[\"classification\", \"regression\"]")},
				},
			}

			result := mapDBModelToAPIModel(catalogModel)

			assert.Equal(t, "123", *result.Id)
			assert.Equal(t, "mapping-test-model", result.Name)
			assert.Equal(t, "mapping-test-ext", *result.ExternalId)
			assert.Equal(t, "test-source", *result.SourceId)
			assert.Equal(t, "Test description", *result.Description)
			assert.Equal(t, "pytorch", *result.LibraryName)
			assert.Equal(t, "1234567890", *result.CreateTimeSinceEpoch)
			assert.Equal(t, "1234567891", *result.LastUpdateTimeSinceEpoch)

			// Verify JSON arrays are properly parsed
			assert.Equal(t, []string{"python", "go"}, result.Language)
			assert.Equal(t, []string{"classification", "regression"}, result.Tasks)
		})

		t.Run("TestMapCatalogArtifactToCatalogArtifact", func(t *testing.T) {
			// Test model artifact mapping
			var catalogModelArtifact models.CatalogModelArtifact = &models.CatalogModelArtifactImpl{
				ID:     apiutils.Of(int32(456)),
				TypeID: apiutils.Of(int32(modelArtifactTypeID)),
				Attributes: &models.CatalogModelArtifactAttributes{
					Name:       apiutils.Of("test-model-artifact"),
					ExternalID: apiutils.Of("test-model-artifact-ext"),
					URI:        apiutils.Of("s3://test/model.bin"),
				},
			}

			catalogArtifact := models.CatalogArtifact{
				CatalogModelArtifact: catalogModelArtifact,
			}

			result, err := mapDBArtifactToAPIArtifact(catalogArtifact)
			require.NoError(t, err)

			assert.NotNil(t, result.CatalogModelArtifact)
			assert.Nil(t, result.CatalogMetricsArtifact)
			assert.Equal(t, "456", *result.CatalogModelArtifact.Id)
			assert.Equal(t, "test-model-artifact", *result.CatalogModelArtifact.Name)
			assert.Equal(t, "s3://test/model.bin", result.CatalogModelArtifact.Uri)

			// Test metrics artifact mapping
			var catalogMetricsArtifact models.CatalogMetricsArtifact = &models.CatalogMetricsArtifactImpl{
				ID:     apiutils.Of(int32(789)),
				TypeID: apiutils.Of(int32(metricsArtifactTypeID)),
				Attributes: &models.CatalogMetricsArtifactAttributes{
					Name:        apiutils.Of("test-metrics-artifact"),
					ExternalID:  apiutils.Of("test-metrics-artifact-ext"),
					MetricsType: models.MetricsTypePerformance,
				},
			}

			catalogArtifact2 := models.CatalogArtifact{
				CatalogMetricsArtifact: catalogMetricsArtifact,
			}

			result2, err := mapDBArtifactToAPIArtifact(catalogArtifact2)
			require.NoError(t, err)

			assert.Nil(t, result2.CatalogModelArtifact)
			assert.NotNil(t, result2.CatalogMetricsArtifact)
			assert.Equal(t, "789", *result2.CatalogMetricsArtifact.Id)
			assert.Equal(t, "test-metrics-artifact", *result2.CatalogMetricsArtifact.Name)
			assert.Equal(t, "performance-metrics", result2.CatalogMetricsArtifact.MetricsType)
		})

		t.Run("TestMapCatalogArtifact_EmptyArtifact", func(t *testing.T) {
			// Test with empty catalog artifact
			emptyCatalogArtifact := models.CatalogArtifact{}

			_, err := mapDBArtifactToAPIArtifact(emptyCatalogArtifact)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid catalog artifact type")
		})
	})

	t.Run("TestErrorHandling", func(t *testing.T) {
		t.Run("TestGetArtifacts_InvalidModelID", func(t *testing.T) {
			// Create a model with invalid ID format for testing
			// This would be an edge case where the ID isn't a valid integer

			// We can't easily test this directly since IDs are generated as integers
			// But we can test the error case by mocking a scenario

			// For now, let's test a scenario where the model exists but has some issue
			params := ListArtifactsParams{
				PageSize:      10,
				OrderBy:       model.ORDERBYFIELD_CREATE_TIME,
				SortOrder:     model.SORTORDER_ASC,
				NextPageToken: apiutils.Of(""),
			}

			_, err := dbCatalog.GetArtifacts(ctx, "non-existent-model", "test-source", params)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid model name")
		})
	})

	t.Run("TestGetFilterOptions", func(t *testing.T) {
		// Create models with various properties for filter options testing
		model1 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("filter-options-model-1"),
				ExternalID: apiutils.Of("filter-opt-1"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("filter-test-source")},
				{Name: "license", StringValue: apiutils.Of("MIT")},
				{Name: "provider", StringValue: apiutils.Of("HuggingFace")},
				{Name: "maturity", StringValue: apiutils.Of("stable")},
				{Name: "library_name", StringValue: apiutils.Of("transformers")},
				{Name: "language", StringValue: apiutils.Of(`["python", "rust"]`)},
				{Name: "tasks", StringValue: apiutils.Of(`["text-classification", "token-classification"]`)},
			},
		}

		model2 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("filter-options-model-2"),
				ExternalID: apiutils.Of("filter-opt-2"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("filter-test-source")},
				{Name: "license", StringValue: apiutils.Of("Apache-2.0")},
				{Name: "provider", StringValue: apiutils.Of("OpenAI")},
				{Name: "maturity", StringValue: apiutils.Of("experimental")},
				{Name: "library_name", StringValue: apiutils.Of("openai")},
				{Name: "language", StringValue: apiutils.Of(`["python", "javascript"]`)},
				{Name: "tasks", StringValue: apiutils.Of(`["text-generation", "conversational"]`)},
				{Name: "readme", StringValue: apiutils.Of("This is a very long readme that exceeds 100 characters and should be excluded from filter options because it's too verbose for filtering purposes.")},
			},
		}

		model3 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("filter-options-model-3"),
				ExternalID: apiutils.Of("filter-opt-3"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("filter-test-source")},
				{Name: "license", StringValue: apiutils.Of("MIT")},
				{Name: "provider", StringValue: apiutils.Of("PyTorch")},
				{Name: "maturity", StringValue: apiutils.Of("stable")},
				{Name: "language", StringValue: apiutils.Of(`["python"]`)},
				{Name: "tasks", StringValue: apiutils.Of(`["image-classification"]`)},
				{Name: "logo", StringValue: apiutils.Of("https://example.com/logo.png")},
				{Name: "license_link", StringValue: apiutils.Of("https://example.com/license")},
			},
		}

		_, err := catalogModelRepo.Save(model1)
		require.NoError(t, err)
		_, err = catalogModelRepo.Save(model2)
		require.NoError(t, err)
		_, err = catalogModelRepo.Save(model3)
		require.NoError(t, err)

		// Test GetFilterOptions
		filterOptions, err := dbCatalog.GetFilterOptions(ctx)
		require.NoError(t, err)
		require.NotNil(t, filterOptions)
		require.NotNil(t, filterOptions.Filters)

		filters := *filterOptions.Filters

		// Should include short properties
		assert.Contains(t, filters, "license")
		assert.Contains(t, filters, "provider")
		assert.Contains(t, filters, "maturity")
		assert.Contains(t, filters, "library_name")
		assert.Contains(t, filters, "language")
		assert.Contains(t, filters, "tasks")

		// Should exclude internal/verbose fields
		assert.NotContains(t, filters, "source_id", "source_id should be excluded")
		assert.NotContains(t, filters, "logo", "logo should be excluded")
		assert.NotContains(t, filters, "license_link", "license_link should be excluded")
		assert.NotContains(t, filters, "readme", "readme should be excluded (too long)")

		licenseFilter := filters["license"]
		assert.Equal(t, "string", licenseFilter.Type)
		assert.NotNil(t, licenseFilter.Values)
		assert.GreaterOrEqual(t, len(licenseFilter.Values), 2, "Should have at least MIT and Apache-2.0")

		// Convert to string slice for easier checking
		licenseValues := make([]string, 0)
		for _, v := range licenseFilter.Values {
			if strVal, ok := v.(string); ok {
				licenseValues = append(licenseValues, strVal)
			}
		}
		assert.Contains(t, licenseValues, "MIT")
		assert.Contains(t, licenseValues, "Apache-2.0")

		// Verify provider filter options
		providerFilter := filters["provider"]
		assert.Equal(t, "string", providerFilter.Type)
		providerValues := make([]string, 0)
		for _, v := range providerFilter.Values {
			if strVal, ok := v.(string); ok {
				providerValues = append(providerValues, strVal)
			}
		}
		assert.Contains(t, providerValues, "HuggingFace")
		assert.Contains(t, providerValues, "OpenAI")
		assert.Contains(t, providerValues, "PyTorch")

		// Verify JSON array fields are properly parsed and expanded
		languageFilter := filters["language"]
		assert.Equal(t, "string", languageFilter.Type)
		languageValues := make([]string, 0)
		for _, v := range languageFilter.Values {
			if strVal, ok := v.(string); ok {
				languageValues = append(languageValues, strVal)
			}
		}
		// Should contain individual values from JSON arrays
		assert.Contains(t, languageValues, "python")
		assert.Contains(t, languageValues, "rust")
		assert.Contains(t, languageValues, "javascript")

		// Verify tasks are properly expanded
		tasksFilter := filters["tasks"]
		assert.Equal(t, "string", tasksFilter.Type)
		tasksValues := make([]string, 0)
		for _, v := range tasksFilter.Values {
			if strVal, ok := v.(string); ok {
				tasksValues = append(tasksValues, strVal)
			}
		}
		assert.Contains(t, tasksValues, "text-classification")
		assert.Contains(t, tasksValues, "token-classification")
		assert.Contains(t, tasksValues, "text-generation")
		assert.Contains(t, tasksValues, "conversational")
		assert.Contains(t, tasksValues, "image-classification")

		// Verify no duplicates
		pythonCount := 0
		for _, v := range languageValues {
			if v == "python" {
				pythonCount++
			}
		}
		assert.Equal(t, 1, pythonCount, "python should appear only once (deduplicated)")

		// Verify maturity options
		maturityFilter := filters["maturity"]
		maturityValues := make([]string, 0)
		for _, v := range maturityFilter.Values {
			if strVal, ok := v.(string); ok {
				maturityValues = append(maturityValues, strVal)
			}
		}
		assert.Contains(t, maturityValues, "stable")
		assert.Contains(t, maturityValues, "experimental")
	})
}

// Helper functions to get type IDs from database

func getCatalogModelTypeIDForDBTest(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogModelTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogModel type")
	}
	return typeRecord.ID
}

func getCatalogModelArtifactTypeIDForDBTest(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogModelArtifactTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogModelArtifact type")
	}
	return typeRecord.ID
}

func getCatalogMetricsArtifactTypeIDForDBTest(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogMetricsArtifactTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogMetricsArtifact type")
	}
	return typeRecord.ID
}
