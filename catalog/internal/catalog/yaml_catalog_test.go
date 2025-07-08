package catalog

import (
	"context"
	"testing"

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

func testYAMLProvider(t *testing.T, path string) CatalogSourceProvider {
	provider, err := newYamlCatalog(&CatalogSourceConfig{
		Properties: map[string]any{
			yamlCatalogPath: path,
		},
	})
	if err != nil {
		t.Fatalf("newYamlCatalog(%s) failed: %v", path, err)
	}
	return provider
}
