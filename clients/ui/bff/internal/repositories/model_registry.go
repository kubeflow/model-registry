package repositories

import (
	"fmt"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type ModelRegistryRepository struct {
}

func NewModelRegistryRepository() *ModelRegistryRepository {
	return &ModelRegistryRepository{}
}

func (m *ModelRegistryRepository) GetAllModelRegistries(client k8s.KubernetesClientInterface, namespace string) ([]models.ModelRegistryModel, error) {

	resources, err := client.GetServiceDetails(namespace)
	if err != nil {
		return nil, fmt.Errorf("error fetching model registries: %w", err)
	}

	var registries = []models.ModelRegistryModel{}
	for _, item := range resources {
		registry := models.ModelRegistryModel{
			Name:        item.Name,
			Description: item.Description,
			DisplayName: item.DisplayName,
		}
		registries = append(registries, registry)
	}

	return registries, nil
}
