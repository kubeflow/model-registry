package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseURLTemplate(t *testing.T) {
	expected := "/v1/model_registry/demo-registry/registered_models/111-222-333/versions"
	tmpl := "/v1/model_registry/:model_registry_id/registered_models/:registered_model_id/versions"
	params := map[string]string{"model_registry_id": "demo-registry", "registered_model_id": "111-222-333"}

	actual := ParseURLTemplate(tmpl, params)

	assert.Equal(t, expected, actual)
}

func TestParseURLTemplateWhenEmpty(t *testing.T) {
	actual := ParseURLTemplate("", nil)
	assert.Empty(t, actual)
}

func TestParseOriginListAllowAll(t *testing.T) {
	expected := []string{"*"}

	actual, ok := ParseOriginList("*")

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestParseOriginListEmpty(t *testing.T) {
	actual, ok := ParseOriginList("")

	assert.False(t, ok)
	assert.Empty(t, actual)
}

func TestParseOriginListSingle(t *testing.T) {
	expected := []string{"http://test.com"}

	actual, ok := ParseOriginList("http://test.com")

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestParseOriginListMultiple(t *testing.T) {
	expected := []string{"http://test.com", "http://test2.com"}
	actual, ok := ParseOriginList("http://test.com,http://test2.com")
	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestParseOriginListMultipleAndSpaces(t *testing.T) {
	expected := []string{"http://test.com", "http://test2.com"}
	actual, ok := ParseOriginList("http://test.com, http://test2.com")
	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}
