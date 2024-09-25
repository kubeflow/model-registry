package data

import (
	"fmt"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations"
)

type ModelRegistryModel struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

func (m ModelRegistryModel) FetchAllModelRegistries(client k8s.KubernetesClientInterface) ([]ModelRegistryModel, error) {

	resources, err := client.GetServiceDetails()
	if err != nil {
		return nil, fmt.Errorf("error fetching model registries: %w", err)
	}

	var registries []ModelRegistryModel = []ModelRegistryModel{}
	for _, item := range resources {
		registry := ModelRegistryModel{
			Name:        item.Name,
			Description: item.Description,
			DisplayName: item.DisplayName,
		}
		registries = append(registries, registry)
	}

	return registries, nil
}
