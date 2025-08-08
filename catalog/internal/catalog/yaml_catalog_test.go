package catalog

import (
	"context"
	"testing"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/stretchr/testify/assert"
)

func TestYAMLCatalogGetModel(t *testing.T) {
	assert := assert.New(t)
	provider := testYAMLProvider(t, "testdata/test-yaml-catalog.yaml")

	model, err := provider.GetModel(context.Background(), "rhelai1/granite-8b-code-base")
	if assert.NoError(err) {
		assert.Equal("rhelai1/granite-8b-code-base", model.Name)

		newLogo := "foobar"
		model.Logo = &newLogo

		model2, err := provider.GetModel(context.Background(), "rhelai1/granite-8b-code-base")
		if assert.NoError(err) {
			assert.NotEqual(model2.Logo, model.Logo, "changes to one returned object should not affect other return values")
		}
	}

	notFound, err := provider.GetModel(context.Background(), "foo")
	assert.NoError(err)
	assert.Nil(notFound)
}

func TestYAMLCatalogGetArtifacts(t *testing.T) {
	assert := assert.New(t)
	provider := testYAMLProvider(t, "testdata/test-yaml-catalog.yaml")

	// Test case 1: Model with artifacts
	artifacts, err := provider.GetArtifacts(context.Background(), "rhelai1/granite-8b-code-base")
	if assert.NoError(err) {
		assert.NotNil(artifacts)
		assert.Equal(int32(1), artifacts.Size)
		assert.Equal(int32(1), artifacts.PageSize)
		assert.Len(artifacts.Items, 1)
		assert.Equal("oci://registry.redhat.io/rhelai1/granite-8b-code-base:1.3-1732870892", artifacts.Items[0].Uri)
	}

	// Test case 2: Model with no artifacts
	noArtifactsModel, err := provider.GetArtifacts(context.Background(), "model-with-no-artifacts")
	if assert.NoError(err) {
		assert.NotNil(noArtifactsModel)
		assert.Equal(int32(0), noArtifactsModel.Size)
		assert.Equal(int32(0), noArtifactsModel.PageSize)
		assert.Len(noArtifactsModel.Items, 0)
	}

	// Test case 3: Model not found
	notFoundArtifacts, err := provider.GetArtifacts(context.Background(), "non-existent-model")
	assert.NoError(err)
	assert.Nil(notFoundArtifacts)
}

func TestYAMLCatalogListModels(t *testing.T) {
	assert := assert.New(t)
	provider := testYAMLProvider(t, "testdata/test-list-models-catalog.yaml")
	ctx := context.Background()

	// Test case 1: List all models, default sort (by name ascending)
	models, err := provider.ListModels(ctx, ListModelsParams{})
	if assert.NoError(err) {
		assert.NotNil(models)
		assert.Equal(int32(6), models.Size)
		assert.Equal(int32(6), models.PageSize)
		assert.Len(models.Items, 6)
		assert.Equal("Z-model", models.Items[0].Name) // Z-model should be first due to string comparison for alphabetical sort
		assert.Equal("another-model-alpha", models.Items[1].Name)
		assert.Equal("model-alpha", models.Items[2].Name)
		assert.Equal("model-beta", models.Items[3].Name)
		assert.Equal("model-gamma", models.Items[4].Name)
		assert.Equal("model-with-no-tasks", models.Items[5].Name)
	}

	// Test case 2: List all models, sort by name ascending
	models, err = provider.ListModels(ctx, ListModelsParams{OrderBy: model.ORDERBYFIELD_NAME, SortOrder: model.SORTORDER_ASC})
	if assert.NoError(err) {
		assert.Equal(int32(6), models.Size)
		assert.Equal("Z-model", models.Items[0].Name)
		assert.Equal("another-model-alpha", models.Items[1].Name)
	}

	// Test case 3: List all models, sort by name descending
	models, err = provider.ListModels(ctx, ListModelsParams{OrderBy: model.ORDERBYFIELD_NAME, SortOrder: model.SORTORDER_DESC})
	if assert.NoError(err) {
		assert.Equal(int32(6), models.Size)
		assert.Equal("model-with-no-tasks", models.Items[0].Name)
		assert.Equal("model-gamma", models.Items[1].Name)
	}

	// Test case 4: List all models, sort by created (CreateTimeSinceEpoch) ascending
	models, err = provider.ListModels(ctx, ListModelsParams{OrderBy: model.ORDERBYFIELD_CREATE_TIME, SortOrder: model.SORTORDER_ASC})
	if assert.NoError(err) {
		assert.Equal(int32(6), models.Size)
		assert.Equal("model-with-no-tasks", models.Items[0].Name) // Jan 1, 2023
		assert.Equal("model-gamma", models.Items[1].Name)         // Feb 1, 2023
	}

	// Test case 5: List all models, sort by published (CreateTimeSinceEpoch) descending
	models, err = provider.ListModels(ctx, ListModelsParams{OrderBy: model.ORDERBYFIELD_CREATE_TIME, SortOrder: model.SORTORDER_DESC})
	if assert.NoError(err) {
		assert.Equal(int32(6), models.Size)
		assert.Equal("Z-model", models.Items[0].Name)             // Aug 2, 2023
		assert.Equal("another-model-alpha", models.Items[1].Name) // May 16, 2023
	}

	// Test case 6: Filter by query "model" (should match all 6 models)
	models, err = provider.ListModels(ctx, ListModelsParams{Query: "model"})
	if assert.NoError(err) {
		assert.Equal(int32(6), models.Size)
		assert.Equal("Z-model", models.Items[0].Name)
		assert.Equal("another-model-alpha", models.Items[1].Name)
		assert.Equal("model-alpha", models.Items[2].Name)
		assert.Equal("model-beta", models.Items[3].Name)
		assert.Equal("model-gamma", models.Items[4].Name)
		assert.Equal("model-with-no-tasks", models.Items[5].Name)
	}

	// Test case 7: Filter by query "text" (should match model-alpha, another-model-alpha)
	models, err = provider.ListModels(ctx, ListModelsParams{Query: "text"})
	if assert.NoError(err) {
		assert.Equal(int32(2), models.Size)
		assert.Equal("another-model-alpha", models.Items[0].Name) // Alphabetical order
		assert.Equal("model-alpha", models.Items[1].Name)
	}

	// Test case 8: Filter by query "nlp" (should match model-alpha, model-gamma, another-model-alpha)
	models, err = provider.ListModels(ctx, ListModelsParams{Query: "nlp"})
	if assert.NoError(err) {
		assert.Equal(int32(3), models.Size)
		assert.Equal("another-model-alpha", models.Items[0].Name)
		assert.Equal("model-alpha", models.Items[1].Name)
		assert.Equal("model-gamma", models.Items[2].Name)
	}

	// Test case 9: Filter by query "IBM" (should match model-alpha, model-gamma)
	models, err = provider.ListModels(ctx, ListModelsParams{Query: "IBM"})
	if assert.NoError(err) {
		assert.Equal(int32(2), models.Size)
		assert.Equal("model-alpha", models.Items[0].Name)
		assert.Equal("model-gamma", models.Items[1].Name)
	}

	// Test case 10: Filter by query "transformers" (should match model-alpha)
	models, err = provider.ListModels(ctx, ListModelsParams{Query: "transformers"})
	if assert.NoError(err) {
		assert.Equal(int32(1), models.Size)
		assert.Equal("model-alpha", models.Items[0].Name)
	}

	// Test case 11: Filter by query "nonexistent" (should return empty list)
	models, err = provider.ListModels(ctx, ListModelsParams{Query: "nonexistent"})
	assert.NoError(err)
	assert.NotNil(models)
	assert.Equal(int32(0), models.Size)
	assert.Equal(int32(0), models.PageSize)
	assert.Len(models.Items, 0)

	// Test case 12: Empty catalog
	emptyProvider := testYAMLProvider(t, "testdata/empty-catalog.yaml") // Assuming an empty-catalog.yaml exists or will be created
	emptyModels, err := emptyProvider.ListModels(ctx, ListModelsParams{})
	assert.NoError(err)
	assert.NotNil(emptyModels)
	assert.Equal(int32(0), emptyModels.Size)
	assert.Equal(int32(0), emptyModels.PageSize)
	assert.Len(emptyModels.Items, 0)

	// Test case 13: Test with excluded models
	excludedProvider := testYAMLProviderWithExclusions(t, "testdata/test-list-models-catalog.yaml", []any{
		"model-alpha",
	})
	excludedModels, err := excludedProvider.ListModels(ctx, ListModelsParams{})
	if assert.NoError(err) {
		assert.NotNil(excludedModels)
		assert.Equal(int32(5), excludedModels.Size)
		for _, m := range excludedModels.Items {
			assert.NotEqual("model-alpha", m.Name)
		}
	}
}

func testYAMLProvider(t *testing.T, path string) CatalogSourceProvider {
	return testYAMLProviderWithExclusions(t, path, nil)
}

func testYAMLProviderWithExclusions(t *testing.T, path string, excludedModels []any) CatalogSourceProvider {
	properties := map[string]any{
		yamlCatalogPath: path,
	}
	if excludedModels != nil {
		properties["excludedModels"] = excludedModels
	}
	provider, err := newYamlCatalog(&CatalogSourceConfig{
		Properties: properties,
	})
	if err != nil {
		t.Fatalf("newYamlCatalog(%s) with exclusions failed: %v", path, err)
	}
	return provider
}
