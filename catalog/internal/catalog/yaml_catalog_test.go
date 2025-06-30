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
