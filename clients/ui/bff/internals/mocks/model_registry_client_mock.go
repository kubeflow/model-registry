package mocks

import (
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/integrations"
	"github.com/stretchr/testify/mock"
	"log/slog"
)

type ModelRegistryClientMock struct {
	mock.Mock
}

func NewModelRegistryClient(logger *slog.Logger) (*ModelRegistryClientMock, error) {
	return &ModelRegistryClientMock{}, nil
}

func (m *ModelRegistryClientMock) GetAllRegisteredModels(client integrations.HTTPClientInterface) (*openapi.RegisteredModelList, error) {
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
