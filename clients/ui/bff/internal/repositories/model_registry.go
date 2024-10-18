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

func (m *ModelRegistryRepository) FetchAllModelRegistries(client k8s.KubernetesClientInterface) ([]models.ModelRegistryModel, error) {

	resources, err := client.GetServiceDetails()
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
