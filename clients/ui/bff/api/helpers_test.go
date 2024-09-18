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
