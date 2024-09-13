package data

import (
	"encoding/json"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/kubeflow/model-registry/ui/bff/internals/mocks"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestGetModelVersion(t *testing.T) {
	gofakeit.Seed(0)

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
