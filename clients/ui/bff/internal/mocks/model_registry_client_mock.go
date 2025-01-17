package mocks

import (
	"log/slog"
	"net/url"

	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations"
	"github.com/stretchr/testify/mock"
)

type ModelRegistryClientMock struct {
	mock.Mock
}

func NewModelRegistryClient(_ *slog.Logger) (*ModelRegistryClientMock, error) {
	return &ModelRegistryClientMock{}, nil
}

func (m *ModelRegistryClientMock) GetAllRegisteredModels(_ integrations.HTTPClientInterface, _ url.Values) (*openapi.RegisteredModelList, error) {
	mockData := GetRegisteredModelListMock()
	return &mockData, nil
}

func (m *ModelRegistryClientMock) CreateRegisteredModel(_ integrations.HTTPClientInterface, _ []byte) (*openapi.RegisteredModel, error) {
	mockData := GetRegisteredModelMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetRegisteredModel(_ integrations.HTTPClientInterface, id string) (*openapi.RegisteredModel, error) {
	if id == "3" {
		mockData := GetRegisteredModelMocks()[2]
		return &mockData, nil
	}
	mockData := GetRegisteredModelMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) UpdateRegisteredModel(_ integrations.HTTPClientInterface, _ string, _ []byte) (*openapi.RegisteredModel, error) {
	mockData := GetRegisteredModelMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetAllModelVersions(_ integrations.HTTPClientInterface) (*openapi.ModelVersionList, error) {
	mockData := GetModelVersionListMock()
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetModelVersion(_ integrations.HTTPClientInterface, id string) (*openapi.ModelVersion, error) {
	if id == "3" {
		mockData := GetModelVersionMocks()[2]
		return &mockData, nil
	}

	mockData := GetModelVersionMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) CreateModelVersion(_ integrations.HTTPClientInterface, _ []byte) (*openapi.ModelVersion, error) {
	mockData := GetModelVersionMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) UpdateModelVersion(_ integrations.HTTPClientInterface, _ string, _ []byte) (*openapi.ModelVersion, error) {
	mockData := GetModelVersionMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetAllModelVersionsForRegisteredModel(_ integrations.HTTPClientInterface, _ string, _ url.Values) (*openapi.ModelVersionList, error) {
	mockData := GetModelVersionListMock()
	return &mockData, nil
}

func (m *ModelRegistryClientMock) CreateModelVersionForRegisteredModel(_ integrations.HTTPClientInterface, _ string, _ []byte) (*openapi.ModelVersion, error) {
	mockData := GetModelVersionMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetModelArtifactsByModelVersion(_ integrations.HTTPClientInterface, _ string, _ url.Values) (*openapi.ModelArtifactList, error) {
	mockData := GetModelArtifactListMock()
	return &mockData, nil
}

func (m *ModelRegistryClientMock) CreateModelArtifactByModelVersion(_ integrations.HTTPClientInterface, _ string, _ []byte) (*openapi.ModelArtifact, error) {
	mockData := GetModelArtifactMocks()[0]
	return &mockData, nil
}
