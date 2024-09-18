package api

import (
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internals/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetRegisteredModelHandler(t *testing.T) {
	data := mocks.GetRegisteredModelMocks()[0]
	expected := RegisteredModelEnvelope{Data: &data}

	actual, rs, err := setupApiTest[RegisteredModelEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/registered_models/1", nil)
	assert.NoError(t, err)

	//TODO assert the full structure, I couldn't get unmarshalling to work for the full customProperties values
	// this issue is in the test only
	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, expected.Data.Name, actual.Data.Name)
}

func TestGetAllRegisteredModelsHandler(t *testing.T) {
	data := mocks.GetRegisteredModelListMock()
	expected := RegisteredModelListEnvelope{Data: &data}

	actual, rs, err := setupApiTest[RegisteredModelListEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/registered_models", nil)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, expected.Data.Size, actual.Data.Size)
	assert.Equal(t, expected.Data.PageSize, actual.Data.PageSize)
	assert.Equal(t, expected.Data.NextPageToken, actual.Data.NextPageToken)
	assert.Equal(t, len(expected.Data.Items), len(actual.Data.Items))
}

func TestCreateRegisteredModelHandler(t *testing.T) {
	data := mocks.GetRegisteredModelMocks()[0]
	expected := RegisteredModelEnvelope{Data: &data}

	body := RegisteredModelEnvelope{Data: openapi.NewRegisteredModel("Model One")}

	actual, rs, err := setupApiTest[RegisteredModelEnvelope](http.MethodPost, "/api/v1/model_registry/model-registry/registered_models", body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rs.StatusCode)
	assert.Equal(t, expected.Data.Name, actual.Data.Name)
	assert.Equal(t, rs.Header.Get("location"), "/api/v1/model_registry/model-registry/registered_models/1")
}

func TestUpdateRegisteredModelHandler(t *testing.T) {
	data := mocks.GetRegisteredModelMocks()[0]
	expected := RegisteredModelEnvelope{Data: &data}

	body := RegisteredModelEnvelope{Data: openapi.NewRegisteredModel("Model One")}

	actual, rs, err := setupApiTest[RegisteredModelEnvelope](http.MethodPatch, "/api/v1/model_registry/model-registry/registered_models/1", body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, expected.Data.Name, actual.Data.Name)
}

func TestGetAllModelVersionsForRegisteredModelHandler(t *testing.T) {
	data := mocks.GetModelVersionListMock()
	expected := ModelVersionListEnvelope{Data: &data}

	actual, rs, err := setupApiTest[ModelVersionListEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/registered_models/1/versions", nil)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, expected.Data.Size, actual.Data.Size)
	assert.Equal(t, expected.Data.PageSize, actual.Data.PageSize)
	assert.Equal(t, expected.Data.NextPageToken, actual.Data.NextPageToken)
	assert.Equal(t, len(expected.Data.Items), len(actual.Data.Items))
}

func TestCreateModelVersionForRegisteredModelHandler(t *testing.T) {
	data := mocks.GetModelVersionMocks()[0]
	expected := ModelVersionEnvelope{Data: &data}

	body := ModelVersionEnvelope{Data: openapi.NewModelVersion("Version Fifty", "")}
	actual, rs, err := setupApiTest[ModelVersionEnvelope](http.MethodPost, "/api/v1/model_registry/model-registry/registered_models/1/versions", body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rs.StatusCode)
	assert.Equal(t, expected.Data.Name, actual.Data.Name)
	assert.Equal(t, rs.Header.Get("Location"), "/api/v1/model_registry/model-registry/model_versions/1")
}
