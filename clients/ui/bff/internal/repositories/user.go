package repositories

import (
	"fmt"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetUser(client k8s.KubernetesClientInterface, identity *k8s.RequestIdentity) (*models.User, error) {

	isAdmin, err := client.IsClusterAdmin(identity)
	if err != nil {
		return nil, fmt.Errorf("error getting user info: %w", err)
	}

	user, err := client.GetUser(identity)
	if err != nil {
		return nil, fmt.Errorf("error getting user info: %w", err)
	}

	var res = models.User{
		UserID:       user,
		ClusterAdmin: isAdmin,
	}

	return &res, nil
}
