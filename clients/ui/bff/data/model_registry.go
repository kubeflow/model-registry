package data

import (
	"fmt"

	k8s "github.com/kubeflow/model-registry/ui/bff/integrations"
)

type ModelRegistryModel struct {
	Name string `json:"name"`
}

func (m ModelRegistryModel) FetchAllModelRegistry(client k8s.KubernetesClientInterface) ([]ModelRegistryModel, error) {

	resources, err := client.FetchServiceNamesByComponent(k8s.ModelRegistryServiceComponentSelector)
	if err != nil {
		return nil, fmt.Errorf("error fetching model registries: %w", err)
	}

	var registries []ModelRegistryModel = []ModelRegistryModel{}
	for _, item := range resources {
		registry := ModelRegistryModel{
			Name: item,
		}
		registries = append(registries, registry)
	}

	return registries, nil
}
