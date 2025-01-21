package datastore

import (
	"os"

	"github.com/golang/glog"
	services "github.com/kubeflow/model-registry/internal/services/inmemory"
	"github.com/kubeflow/model-registry/pkg/api"
)

type InMemoryService struct{}

func NewInMemoryService() *InMemoryService {
	return &InMemoryService{}
}

func (s *InMemoryService) Build() (api.ModelRegistryApi, error) {
	inMemory := services.NewInMemory()

	path := os.Getenv("SEED_DATA_PATH")

	if path != "" {
		glog.Infof("Seeding in-memory service with data from %s", path)

		err := inMemory.Seed(path)
		if err != nil {
			return nil, err
		}

		glog.Info("Successfully seeded in-memory service")
	}

	return inMemory, nil
}

func (s *InMemoryService) Teardown() error {
	return nil
}
