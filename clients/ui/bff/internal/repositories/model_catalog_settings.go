package repositories

import (
	"context"
	"fmt"
	"log/slog"

	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type ModelCatalogSettingsRepositoryInterface interface {
	GetAllCatalogSourceConfigs(ctx context.Context, k8sClient k8s.KubernetesClientInterface, namespace string) (*models.CatalogSourceConfigList, error)
	GetCatalogSourceConfig(ctx context.Context, k8sClient k8s.KubernetesClientInterface, namespace string, catalogSourceId string) (*models.CatalogSourceConfig, error)
	CreateCatalogSourceConfig(ctx context.Context, k8sClient k8s.KubernetesClientInterface, namespace string, payload models.CatalogSourceConfigPayload) (*models.CatalogSourceConfig, error)
	UpdateCatalogSourceConfig(ctx context.Context, k8sClient k8s.KubernetesClientInterface, namespace string, payload models.CatalogSourceConfigPayload) (*models.CatalogSourceConfig, error)
	DeleteCatalogSourceConfig(ctx context.Context, k8sClient k8s.KubernetesClientInterface, namespace string, catalogSourceId string) (*models.CatalogSourceConfig, error)
}

type ModelCatalogSettingsRepository struct {
}

func NewModelCatalogSettingsRepository(logger *slog.Logger) (*ModelCatalogSettingsRepository, error) {
	return &ModelCatalogSettingsRepository{}, nil
}

func (r *ModelCatalogSettingsRepository) GetAllCatalogSourceConfigs(_ context.Context, _ k8s.KubernetesClientInterface, namespace string) (*models.CatalogSourceConfigList, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *ModelCatalogSettingsRepository) GetCatalogSourceConfig(ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	catalogSourceId string,
) (*models.CatalogSourceConfig, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *ModelCatalogSettingsRepository) CreateCatalogSourceConfig(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	payload models.CatalogSourceConfigPayload,
) (*models.CatalogSourceConfig, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *ModelCatalogSettingsRepository) UpdateCatalogSourceConfig(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	payload models.CatalogSourceConfigPayload,
) (*models.CatalogSourceConfig, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *ModelCatalogSettingsRepository) DeleteCatalogSourceConfig(
	ctx context.Context,
	client k8s.KubernetesClientInterface,
	namespace string,
	catalogSourceId string,
) (*models.CatalogSourceConfig, error) {
	return nil, fmt.Errorf("not implemented yet")
}
