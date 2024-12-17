package repositories

import (
	"fmt"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type NamespaceRepository struct{}

func NewNamespaceRepository() *NamespaceRepository {
	return &NamespaceRepository{}
}

func (r *NamespaceRepository) GetNamespaces(client k8s.KubernetesClientInterface, user string, groups []string) ([]models.NamespaceModel, error) {

	namespaces, err := client.GetNamespaces(user, groups)
	if err != nil {
		return nil, fmt.Errorf("error fetching namespaces: %w", err)
	}

	var namespaceModels = []models.NamespaceModel{}
	for _, ns := range namespaces {
		namespaceModels = append(namespaceModels, models.NewNamespaceModelFromNamespace(ns.Name))
	}

	return namespaceModels, nil
}
