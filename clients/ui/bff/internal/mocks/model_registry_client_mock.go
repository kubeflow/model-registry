package mocks

import (
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/mrserver"
	"log/slog"
	"net/url"

	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/stretchr/testify/mock"
)

type ModelRegistryClientMock struct {
	mock.Mock
}

func NewModelRegistryClient(_ *slog.Logger) (*ModelRegistryClientMock, error) {
	return &ModelRegistryClientMock{}, nil
}

func (m *ModelRegistryClientMock) GetAllRegisteredModels(_ mrserver.HTTPClientInterface, _ url.Values) (*openapi.RegisteredModelList, error) {
	mockData := GetRegisteredModelListMock()
	return &mockData, nil
}

func (m *ModelRegistryClientMock) CreateRegisteredModel(_ mrserver.HTTPClientInterface, _ []byte) (*openapi.RegisteredModel, error) {
	mockData := GetRegisteredModelMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetRegisteredModel(_ mrserver.HTTPClientInterface, id string) (*openapi.RegisteredModel, error) {
	if id == "3" {
		mockData := GetRegisteredModelMocks()[2]
		return &mockData, nil
	}
	mockData := GetRegisteredModelMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) UpdateRegisteredModel(_ mrserver.HTTPClientInterface, _ string, _ []byte) (*openapi.RegisteredModel, error) {
	mockData := GetRegisteredModelMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetAllModelVersions(_ mrserver.HTTPClientInterface) (*openapi.ModelVersionList, error) {
	mockData := GetModelVersionListMock()
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetModelVersion(_ mrserver.HTTPClientInterface, id string) (*openapi.ModelVersion, error) {
	if id == "3" {
		mockData := GetModelVersionMocks()[2]
		return &mockData, nil
	}

	mockData := GetModelVersionMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) CreateModelVersion(_ mrserver.HTTPClientInterface, _ []byte) (*openapi.ModelVersion, error) {
	mockData := GetModelVersionMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) UpdateModelVersion(_ mrserver.HTTPClientInterface, _ string, _ []byte) (*openapi.ModelVersion, error) {
	mockData := GetModelVersionMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetAllModelVersionsForRegisteredModel(_ mrserver.HTTPClientInterface, _ string, _ url.Values) (*openapi.ModelVersionList, error) {
	mockData := GetModelVersionListMock()
	return &mockData, nil
}

func (m *ModelRegistryClientMock) CreateModelVersionForRegisteredModel(_ mrserver.HTTPClientInterface, _ string, _ []byte) (*openapi.ModelVersion, error) {
	mockData := GetModelVersionMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetModelArtifactsByModelVersion(_ mrserver.HTTPClientInterface, _ string, _ url.Values) (*openapi.ModelArtifactList, error) {
	mockData := GetModelArtifactListMock()
	return &mockData, nil
}

func (m *ModelRegistryClientMock) CreateModelArtifactByModelVersion(_ mrserver.HTTPClientInterface, _ string, _ []byte) (*openapi.ModelArtifact, error) {
	mockData := GetModelArtifactMocks()[0]
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetAllArtifacts(_ mrserver.HTTPClientInterface, _ url.Values) (*openapi.ArtifactList, error) {
	mockData := GenerateMockArtifactList()
	return &mockData, nil
}

func (m *ModelRegistryClientMock) GetArtifact(_ mrserver.HTTPClientInterface, _ string) (*openapi.Artifact, error) {
	mockData := GenerateMockArtifact()
	return &mockData, nil
}

func (m *ModelRegistryClientMock) CreateArtifact(_ mrserver.HTTPClientInterface, _ []byte) (*openapi.Artifact, error) {
	mockData := GenerateMockArtifact()
	return &mockData, nil
}

func (m *ModelRegistryClientMock) UpdateModelArtifact(_ mrserver.HTTPClientInterface, _ string, _ []byte) (*openapi.ModelArtifact, error) {
	mockData := GenerateMockModelArtifact()
	return &mockData, nil
}
