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

func (r *ModelRegistrySettingsRepository) GetGroups(ctx context.Context, client k8s.KubernetesClientInterface) ([]models.Group, error) {
	groupNames, err := client.GetGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching groups: %w", err)
	}

	var groups []models.Group
	for _, name := range groupNames {
		// Create mock users for each group to make the data more realistic
		var users []string
		switch name {
		case "dora-group-mock":
			users = []string{"dora-user@example.com", "dora-admin@example.com"}
		case "bella-group-mock":
			users = []string{"bella-user@example.com", "bella-maintainer@example.com"}
		default:
			users = []string{fmt.Sprintf("%s-user@example.com", name)}
		}

		groups = append(groups, models.NewGroup(name, users))
	}

	return groups, nil
}
