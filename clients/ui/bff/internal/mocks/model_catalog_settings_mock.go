package mocks

import (
	"context"
	"fmt"
	"log/slog"

	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type ModelCatalogSettingsRepositoryMock struct {
}

func NewModelCatalogSettingsRepository(logger *slog.Logger) (*ModelCatalogSettingsRepositoryMock, error) {
	return &ModelCatalogSettingsRepositoryMock{}, nil
}

func (m *ModelCatalogSettingsRepositoryMock) GetAllCatalogSourceConfigs(_ context.Context, _ k8s.KubernetesClientInterface, _ string) (*models.CatalogSourceConfigList, error) {
	allCatalogSourceConfigs := GetCatalogSourceConfigListMock()

	return &allCatalogSourceConfigs, nil
}

func (m *ModelCatalogSettingsRepositoryMock) GetCatalogSourceConfig(_ context.Context, _ k8s.KubernetesClientInterface, _ string, catalogSourceId string) (*models.CatalogSourceConfig, error) {
	catalogSourceConfig := CreateSampleCatalogSource(catalogSourceId, "catalog-source-1", "yaml")

	return &catalogSourceConfig, nil
}

func (m *ModelCatalogSettingsRepositoryMock) CreateCatalogSourceConfig(_ context.Context, _ k8s.KubernetesClientInterface, _ string, sourceConfigPayload models.CatalogSourceConfigPayload) (*models.CatalogSourceConfig, error) {
	var sourceName = sourceConfigPayload.Name
	var sourceId = sourceConfigPayload.Id
	var sourceType = sourceConfigPayload.Type

	if sourceName == "" {
		return nil, fmt.Errorf("source name is required")
	}
	if sourceId == "" {
		return nil, fmt.Errorf("source ID is required")
	}
	if sourceType == "" {
		return nil, fmt.Errorf("source type is required")
	}

	newCatalogSource := CreateSampleCatalogSource(sourceId, sourceName, sourceType)

	return &newCatalogSource, nil
}

func (m *ModelCatalogSettingsRepositoryMock) UpdateCatalogSourceConfig(_ context.Context, _ k8s.KubernetesClientInterface, _ string, sourceConfigPayload models.CatalogSourceConfigPayload) (*models.CatalogSourceConfig, error) {
	catalogSourceId := sourceConfigPayload.Id

	updatedCatalogSource := CreateSampleCatalogSource(catalogSourceId, "Updated Catalog", "yaml")

	return &updatedCatalogSource, nil
}

func (m *ModelCatalogSettingsRepositoryMock) DeleteCatalogSourceConfig(_ context.Context, _ k8s.KubernetesClientInterface, _ string, catalogSourceId string) (*models.CatalogSourceConfig, error) {
	deletedCatalogSource := CreateSampleCatalogSource(catalogSourceId, "Updated Catalog", "yaml")

	return &deletedCatalogSource, nil
}
