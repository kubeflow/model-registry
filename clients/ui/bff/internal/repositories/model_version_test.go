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

func TestGetModelVersion(t *testing.T) {
	_ = gofakeit.Seed(0)

	expected := mocks.GenerateMockModelVersion()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	modelVersion := ModelVersion{}

	path, err := url.JoinPath(modelVersionPath, expected.GetId())
	assert.NoError(t, err)

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On("GET", path).Return(mockData, nil)

	actual, err := modelVersion.GetModelVersion(mockClient, expected.GetId())
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, *expected.Author, *actual.Author)

	mockClient.AssertExpectations(t)
}

func TestGetAllModelVersions(t *testing.T) {
	_ = gofakeit.Seed(0)

	expected := mocks.GenerateMockModelVersionList()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	modelVersion := ModelVersion{}

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On("GET", modelVersionPath).Return(mockData, nil)

	actual, err := modelVersion.GetAllModelVersions(mockClient)
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

func TestCreateModelVersion(t *testing.T) {
	_ = gofakeit.Seed(0)

	expected := mocks.GenerateMockModelVersion()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	modelVersion := ModelVersion{}

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On("POST", modelVersionPath, mock.Anything).Return(mockData, nil)

	jsonInput, err := json.Marshal(expected)
	assert.NoError(t, err)

	actual, err := modelVersion.CreateModelVersion(mockClient, jsonInput)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, *expected.Author, *actual.Author)

	mockClient.AssertExpectations(t)
}

func TestUpdateModelVersion(t *testing.T) {
	_ = gofakeit.Seed(0)

	expected := mocks.GenerateMockModelVersion()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	modelVersion := ModelVersion{}

	path, err := url.JoinPath(modelVersionPath, expected.GetId())
	assert.NoError(t, err)

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On(http.MethodPatch, path, mock.Anything).Return(mockData, nil)

	jsonInput, err := json.Marshal(expected)
	assert.NoError(t, err)

	actual, err := modelVersion.UpdateModelVersion(mockClient, expected.GetId(), jsonInput)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, *expected.Author, *actual.Author)

	mockClient.AssertExpectations(t)
}

func TestGetModelArtifactsByModelVersion(t *testing.T) {
	_ = gofakeit.Seed(0)

	expected := mocks.GenerateMockModelArtifactList()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	modelVersion := ModelVersion{}

	path, err := url.JoinPath(modelVersionPath, "1", artifactsByModelVersionPath)
	assert.NoError(t, err)

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On(http.MethodGet, path, mock.Anything).Return(mockData, nil)

	actual, err := modelVersion.GetModelArtifactsByModelVersion(mockClient, "1", nil)
	assert.NoError(t, err)

	assert.NotNil(t, actual)
	assert.Equal(t, expected.Size, actual.Size)
	assert.Equal(t, expected.NextPageToken, actual.NextPageToken)
	assert.Equal(t, expected.PageSize, actual.PageSize)
	assert.Equal(t, len(expected.Items), len(actual.Items))
}

func TestGetModelArtifactsByModelVersionWithPageParams(t *testing.T) {
	gofakeit.Seed(0) //nolint:errcheck

	pageValues := mocks.GenerateMockPageValues()
	expected := mocks.GenerateMockModelArtifactList()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	modelVersion := ModelVersion{}

	path, err := url.JoinPath(modelVersionPath, "1", artifactsByModelVersionPath)
	assert.NoError(t, err)
	reqUrl := fmt.Sprintf("%s?%s", path, pageValues.Encode())

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On(http.MethodGet, reqUrl, mock.Anything).Return(mockData, nil)

	actual, err := modelVersion.GetModelArtifactsByModelVersion(mockClient, "1", pageValues)
	assert.NoError(t, err)

	assert.NotNil(t, actual)
	mockClient.AssertExpectations(t)
}

func TestCreateModelArtifactByModelVersion(t *testing.T) {
	_ = gofakeit.Seed(0)

	expected := mocks.GenerateMockModelArtifact()

	mockData, err := json.Marshal(expected)
	assert.NoError(t, err)

	modelVersion := ModelVersion{}

	path, err := url.JoinPath(modelVersionPath, "1", artifactsByModelVersionPath)
	assert.NoError(t, err)

	mockClient := new(mocks.MockHTTPClient)
	mockClient.On(http.MethodPost, path, mock.Anything).Return(mockData, nil)

	jsonInnput, err := json.Marshal(expected)
	assert.NoError(t, err)

	actual, err := modelVersion.CreateModelArtifactByModelVersion(mockClient, "1", jsonInnput)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.ArtifactType, actual.ArtifactType)
}
