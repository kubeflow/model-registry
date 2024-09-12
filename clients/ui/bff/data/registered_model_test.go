package data

import (
	"encoding/json"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/kubeflow/model-registry/ui/bff/internals/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"testing"
)

func TestGetAllRegisteredModels(t *testing.T) {
	gofakeit.Seed(0)

	expected := mocks.GenerateMockRegisteredModelList()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	mrClient := ModelRegistryClient{}

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On("GET", registeredModelPath).Return(mockData, nil)

	actual, err := mrClient.GetAllRegisteredModels(mockClient)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.NextPageToken, actual.NextPageToken)
	assert.Equal(t, expected.PageSize, actual.PageSize)
	assert.Equal(t, expected.Size, actual.Size)
	assert.Equal(t, len(expected.Items), len(actual.Items))

	mockClient.AssertExpectations(t)
}

func TestCreateRegisteredModel(t *testing.T) {
	gofakeit.Seed(0)

	expected := mocks.GenerateRegisteredModel()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	mrClient := ModelRegistryClient{}

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On("POST", registeredModelPath, mock.Anything).Return(mockData, nil)

	jsonInput, err := json.Marshal(expected)
	assert.NoError(t, err)

	actual, err := mrClient.CreateRegisteredModel(mockClient, jsonInput)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, *expected.Owner, *actual.Owner)

	mockClient.AssertExpectations(t)
}

func TestGetRegisteredModel(t *testing.T) {
	gofakeit.Seed(0)

	expected := mocks.GenerateRegisteredModel()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	mrClient := ModelRegistryClient{}

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On("GET", registeredModelPath+"/"+expected.GetId()).Return(mockData, nil)

	actual, err := mrClient.GetRegisteredModel(mockClient, expected.GetId())
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, *expected.Owner, *actual.Owner)

	mockClient.AssertExpectations(t)
}

func TestUpdateRegisteredModel(t *testing.T) {
	gofakeit.Seed(0)

	expected := mocks.GenerateRegisteredModel()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	mrClient := ModelRegistryClient{}

	path, err := url.JoinPath(registeredModelPath, expected.GetId())
	assert.NoError(t, err)

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On(http.MethodPatch, path, mock.Anything).Return(mockData, nil)

	jsonInput, err := json.Marshal(expected)
	assert.NoError(t, err)

	actual, err := mrClient.UpdateRegisteredModel(mockClient, expected.GetId(), jsonInput)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, *expected.Owner, *actual.Owner)

	mockClient.AssertExpectations(t)
}
