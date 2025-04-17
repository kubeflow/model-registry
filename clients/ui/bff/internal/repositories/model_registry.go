package repositories

import (
	"context"
	"fmt"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"

	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type ModelRegistryRepository struct {
}

func NewModelRegistryRepository() *ModelRegistryRepository {
	return &ModelRegistryRepository{}
}

func (m *ModelRegistryRepository) GetAllModelRegistries(sessionCtx context.Context, client k8s.KubernetesClientInterface, namespace string) ([]models.ModelRegistryModel, error) {

	resources, err := client.GetServiceDetails(sessionCtx, namespace)
	if err != nil {
		return nil, fmt.Errorf("error fetching model registries: %w", err)
	}

	var registries = []models.ModelRegistryModel{}
	for _, s := range resources {
		serverAddress := m.ResolveServerAddress(s.ClusterIP, s.HTTPPort)
		registry := models.ModelRegistryModel{
			Name:          s.Name,
			Description:   s.Description,
			DisplayName:   s.DisplayName,
			ServerAddress: serverAddress,
		}
		registries = append(registries, registry)
	}

	return registries, nil
}

func (m *ModelRegistryRepository) GetModelRegistry(sessionCtx context.Context, client k8s.KubernetesClientInterface, namespace string, modelRegistryID string) (models.ModelRegistryModel, error) {

	s, err := client.GetServiceDetailsByName(sessionCtx, namespace, modelRegistryID)
	if err != nil {
		return models.ModelRegistryModel{}, fmt.Errorf("error fetching model registry: %w", err)
	}

	modelRegistry := models.ModelRegistryModel{
		Name:          s.Name,
		Description:   s.Description,
		DisplayName:   s.DisplayName,
		ServerAddress: m.ResolveServerAddress(s.ClusterIP, s.HTTPPort),
	}

	return modelRegistry, nil
}

func (m *ModelRegistryRepository) ResolveServerAddress(clusterIP string, httpPort int32) string {
	url := fmt.Sprintf("http://%s:%d/api/model_registry/v1alpha3", clusterIP, httpPort)
	return url
}
