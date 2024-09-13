package api

import (
	"context"
	"encoding/json"
	"github.com/kubeflow/model-registry/ui/bff/internals/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetModelVersionHandler(t *testing.T) {
	data := mocks.GetModelVersionMocks()[0]
	expected := ModelVersionEnvelope{Data: &data}

	actual, rs, err := setupApiTest[ModelVersionEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/model_versions/1", nil)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, expected.Data.Name, actual.Data.Name)
}
