package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllRegisteredModels(t *testing.T) {
	_ = gofakeit.Seed(0)

	expected := mocks.GenerateMockRegisteredModelList()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	registeredModel := RegisteredModel{}

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On("GET", registeredModelPath).Return(mockData, nil)

	actual, err := registeredModel.GetAllRegisteredModels(mockClient, nil)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.NextPageToken, actual.NextPageToken)
	assert.Equal(t, expected.PageSize, actual.PageSize)
	assert.Equal(t, expected.Size, actual.Size)
	assert.Equal(t, len(expected.Items), len(actual.Items))

	mockClient.AssertExpectations(t)
}

func TestGetAllRegisteredModelsWithPageParams(t *testing.T) {
	gofakeit.Seed(0) //nolint:errcheck

	pageValues := mocks.GenerateMockPageValues()
	expected := mocks.GenerateMockRegisteredModelList()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	reqUrl := fmt.Sprintf("%s?%s", registeredModelPath, pageValues.Encode())

	registeredModel := RegisteredModel{}

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On("GET", reqUrl).Return(mockData, nil)

	actual, err := registeredModel.GetAllRegisteredModels(mockClient, pageValues)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	mockClient.AssertExpectations(t)
}

func TestCreateRegisteredModel(t *testing.T) {
	_ = gofakeit.Seed(0)

	expected := mocks.GenerateMockRegisteredModel()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	registeredModel := RegisteredModel{}

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On("POST", registeredModelPath, mock.Anything).Return(mockData, nil)

	jsonInput, err := json.Marshal(expected)
	assert.NoError(t, err)

	actual, err := registeredModel.CreateRegisteredModel(mockClient, jsonInput)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, *expected.Owner, *actual.Owner)

	mockClient.AssertExpectations(t)
}

func TestGetRegisteredModel(t *testing.T) {
	_ = gofakeit.Seed(0)

	expected := mocks.GenerateMockRegisteredModel()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	registeredModel := RegisteredModel{}

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On("GET", registeredModelPath+"/"+expected.GetId()).Return(mockData, nil)

	actual, err := registeredModel.GetRegisteredModel(mockClient, expected.GetId())
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, *expected.Owner, *actual.Owner)

	mockClient.AssertExpectations(t)
}

func TestUpdateRegisteredModel(t *testing.T) {
	_ = gofakeit.Seed(0)

	expected := mocks.GenerateMockRegisteredModel()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	registeredModel := RegisteredModel{}

	path, err := url.JoinPath(registeredModelPath, expected.GetId())
	assert.NoError(t, err)

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On(http.MethodPatch, path, mock.Anything).Return(mockData, nil)

	jsonInput, err := json.Marshal(expected)
	assert.NoError(t, err)

	actual, err := registeredModel.UpdateRegisteredModel(mockClient, expected.GetId(), jsonInput)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, *expected.Owner, *actual.Owner)

	mockClient.AssertExpectations(t)
}

func TestGetAllModelVersionsByRegisteredModel(t *testing.T) {
	_ = gofakeit.Seed(0)

	expected := mocks.GenerateMockModelVersionList()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	registeredModel := RegisteredModel{}

	mockClient := new(mocks.MockHTTPClient)
	path, err := url.JoinPath(registeredModelPath, "1", versionsPath)
	assert.NoError(t, err)
	mockClient.On("GET", path).Return(mockData, nil)

	actual, err := registeredModel.GetAllModelVersionsForRegisteredModel(mockClient, "1", nil)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.NextPageToken, actual.NextPageToken)
	assert.Equal(t, expected.PageSize, actual.PageSize)
	assert.Equal(t, expected.Size, actual.Size)
	assert.Equal(t, len(expected.Items), len(actual.Items))

	mockClient.AssertExpectations(t)
}

func TestGetAllModelVersionsWithPageParams(t *testing.T) {
	gofakeit.Seed(0) //nolint:errcheck

	pageValues := mocks.GenerateMockPageValues()
	expected := mocks.GenerateMockModelVersionList()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	registeredModel := RegisteredModel{}

	mockClient := new(mocks.MockHTTPClient)
	path, err := url.JoinPath(registeredModelPath, "1", versionsPath)
	assert.NoError(t, err)
	reqUrl := fmt.Sprintf("%s?%s", path, pageValues.Encode())

	mockClient.On("GET", reqUrl).Return(mockData, nil)

	actual, err := registeredModel.GetAllModelVersionsForRegisteredModel(mockClient, "1", pageValues)
	assert.NoError(t, err)
	assert.NotNil(t, actual)

	mockClient.AssertExpectations(t)
}

func TestCreateModelVersionForRegisteredModel(t *testing.T) {
	_ = gofakeit.Seed(0)

	expected := mocks.GenerateMockModelVersion()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	registeredModel := RegisteredModel{}

	path, err := url.JoinPath(registeredModelPath, "1", versionsPath)
	assert.NoError(t, err)
	mockClient := new(mocks.MockHTTPClient)
	mockClient.On(http.MethodPost, path, mock.Anything).Return(mockData, nil)

	jsonInput, err := json.Marshal(expected)
	assert.NoError(t, err)

	actual, err := registeredModel.CreateModelVersionForRegisteredModel(mockClient, "1", jsonInput)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, *expected.Author, *actual.Author)
	assert.Equal(t, expected.RegisteredModelId, actual.RegisteredModelId)

	mockClient.AssertExpectations(t)
}
