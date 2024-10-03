package mocks

import (
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/url"
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

func (m *ModelRegistryClientMock) CreateRegisteredModel(client integrations.HTTPClientInterface, jsonData []byte) (*openapi.RegisteredModel, error) {
	mockData := GetRegisteredModelMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetRegisteredModel(client integrations.HTTPClientInterface, id string) (*openapi.RegisteredModel, error) {
	mockData := GetRegisteredModelMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) UpdateRegisteredModel(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.RegisteredModel, error) {
	mockData := GetRegisteredModelMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetModelVersion(client integrations.HTTPClientInterface, id string) (*openapi.ModelVersion, error) {
	mockData := GetModelVersionMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) CreateModelVersion(client integrations.HTTPClientInterface, jsonData []byte) (*openapi.ModelVersion, error) {
	mockData := GetModelVersionMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) UpdateModelVersion(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.ModelVersion, error) {
	mockData := GetModelVersionMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetAllModelVersions(_ integrations.HTTPClientInterface, _ string, _ url.Values) (*openapi.ModelVersionList, error) {
	mockData := GetModelVersionListMock()
	return &mockData, nil
}

func (m *ModelRegistryClientMock) CreateModelVersionForRegisteredModel(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.ModelVersion, error) {
	mockData := GetModelVersionMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetModelArtifactsByModelVersion(_ integrations.HTTPClientInterface, _ string, _ url.Values) (*openapi.ModelArtifactList, error) {
	mockData := GetModelArtifactListMock()
	return &mockData, nil
}

func (m *ModelRegistryClientMock) CreateModelArtifactByModelVersion(client integrations.HTTPClientInterface, id string, jsonData []byte) (*openapi.ModelArtifact, error) {
	mockData := GetModelArtifactMocks()[0]
	return &mockData, nil
}
