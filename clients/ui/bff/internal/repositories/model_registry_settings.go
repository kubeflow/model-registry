package repositories

import (
	"context"
	"fmt"

	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type ModelRegistrySettingsRepository struct {
}

func NewModelRegistrySettingsRepository() *ModelRegistrySettingsRepository {
	return &ModelRegistrySettingsRepository{}
}

func (r *ModelRegistrySettingsRepository) GetGroups(ctx context.Context, client k8s.KubernetesClientInterface) ([]models.GroupModel, error) {
	groupNames, err := client.GetGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching groups: %w", err)
	}

	var groups []models.GroupModel
	for _, name := range groupNames {
		groups = append(groups, models.NewGroupModel(name))
	}

	return groups, nil
}
