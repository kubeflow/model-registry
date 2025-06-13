package mocks

import (
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type MockUserRepository struct{}

func (r *MockUserRepository) GetUser(client kubernetes.KubernetesClientInterface, identity *kubernetes.RequestIdentity) (*models.User, error) {
	return &models.User{
		UserID:       "user",
		ClusterAdmin: true,
	}, nil
}
