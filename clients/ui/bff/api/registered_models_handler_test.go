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
		"/api/v1/model-registry/model-registry/registered_models/1", nil)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	testApp.GetRegisteredModelHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var registeredModelRes TypedEnvelope[openapi.RegisteredModel]
	err = json.Unmarshal(body, &registeredModelRes)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	var expected = TypedEnvelope[openapi.RegisteredModel]{
		"registered_model": mocks.GetRegisteredModelMocks()[0],
	}

	//TODO assert the full structure, I couldn't get unmarshalling to work for the full customProperties values
	// this issue is in the test only
	assert.Equal(t, expected["registered_model"].Name, registeredModelRes["registered_model"].Name)
}

func TestGetAllRegisteredModelsHandler(t *testing.T) {
	mockMRClient, _ := mocks.NewModelRegistryClient(nil)
	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		modelRegistryClient: mockMRClient,
	}

	req, err := http.NewRequest(http.MethodGet,
		"/api/v1/model-registry/model-registry/registered_models", nil)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	testApp.GetAllRegisteredModelsHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var registeredModelsListRes TypedEnvelope[openapi.RegisteredModelList]
	err = json.Unmarshal(body, &registeredModelsListRes)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	var expected = TypedEnvelope[openapi.RegisteredModelList]{
		"registered_model_list": mocks.GetRegisteredModelListMock(),
	}

	assert.Equal(t, expected["registered_model_list"].Size, registeredModelsListRes["registered_model_list"].Size)
	assert.Equal(t, expected["registered_model_list"].PageSize, registeredModelsListRes["registered_model_list"].PageSize)
	assert.Equal(t, expected["registered_model_list"].NextPageToken, registeredModelsListRes["registered_model_list"].NextPageToken)
	assert.Equal(t, len(expected["registered_model_list"].Items), len(registeredModelsListRes["registered_model_list"].Items))
}

func TestCreateRegisteredModelHandler(t *testing.T) {
	mockMRClient, _ := mocks.NewModelRegistryClient(nil)
	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		modelRegistryClient: mockMRClient,
	}

	newModel := openapi.NewRegisteredModelCreate("Model One")
	newModelJSON, err := newModel.MarshalJSON()
	assert.NoError(t, err)

	reqBody := bytes.NewReader(newModelJSON)

	req, err := http.NewRequest(http.MethodPost,
		"/api/v1/model-registry/model-registry/registered_models", reqBody)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	testApp.CreateRegisteredModelHandler(rr, req, nil)
	rs := rr.Result()

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	assert.NoError(t, err)
	var registeredModelRes openapi.RegisteredModel
	err = json.Unmarshal(body, &registeredModelRes)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var expected = mocks.GetRegisteredModelMocks()[0]

	assert.Equal(t, expected.Name, registeredModelRes.Name)
	assert.NotEmpty(t, rs.Header.Get("location"))
}
