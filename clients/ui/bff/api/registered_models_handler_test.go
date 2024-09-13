package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internals/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRegisteredModelHandler(t *testing.T) {
	mockMRClient, _ := mocks.NewModelRegistryClient(nil)
	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		modelRegistryClient: mockMRClient,
	}

	req, err := http.NewRequest(http.MethodGet,
		"/api/v1/model_registry/model-registry/registered_models/1", nil)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	testApp.GetRegisteredModelHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var registeredModelRes RegisteredModelEnvelope
	err = json.Unmarshal(body, &registeredModelRes)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	mockModel := mocks.GetRegisteredModelMocks()[0]

	var expected = RegisteredModelEnvelope{
		Data: &mockModel,
	}

	//TODO assert the full structure, I couldn't get unmarshalling to work for the full customProperties values
	// this issue is in the test only
	assert.Equal(t, expected.Data.Name, registeredModelRes.Data.Name)
}

func TestGetAllRegisteredModelsHandler(t *testing.T) {
	mockMRClient, _ := mocks.NewModelRegistryClient(nil)
	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		modelRegistryClient: mockMRClient,
	}

	req, err := http.NewRequest(http.MethodGet,
		"/api/v1/model_registry/model-registry/registered_models", nil)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	testApp.GetAllRegisteredModelsHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var registeredModelsListRes RegisteredModelListEnvelope
	err = json.Unmarshal(body, &registeredModelsListRes)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	modelList := mocks.GetRegisteredModelListMock()

	var expected = RegisteredModelListEnvelope{
		Data: &modelList,
	}

	assert.Equal(t, expected.Data.Size, registeredModelsListRes.Data.Size)
	assert.Equal(t, expected.Data.PageSize, registeredModelsListRes.Data.PageSize)
	assert.Equal(t, expected.Data.NextPageToken, registeredModelsListRes.Data.NextPageToken)
	assert.Equal(t, len(expected.Data.Items), len(registeredModelsListRes.Data.Items))
}

func TestCreateRegisteredModelHandler(t *testing.T) {
	mockMRClient, _ := mocks.NewModelRegistryClient(nil)
	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		modelRegistryClient: mockMRClient,
	}

	newModel := openapi.NewRegisteredModel("Model One")
	newEnvelope := RegisteredModelEnvelope{Data: newModel}

	newModelJSON, err := json.Marshal(newEnvelope)
	assert.NoError(t, err)

	reqBody := bytes.NewReader(newModelJSON)

	req, err := http.NewRequest(http.MethodPost,
		"/api/v1/model_registry/model-registry/registered_models", reqBody)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	testApp.CreateRegisteredModelHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var actual RegisteredModelEnvelope
	err = json.Unmarshal(body, &actual)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var expected = mocks.GetRegisteredModelMocks()[0]

	assert.Equal(t, expected.Name, actual.Data.Name)
	assert.NotEmpty(t, rs.Header.Get("location"))
}
