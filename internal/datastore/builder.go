package datastore

import (
	"errors"
	"fmt"

	"github.com/kubeflow/model-registry/pkg/api"
)

var (
	ErrCreatingDatastore    = errors.New("error creating datastore")
	ErrUnsupportedDatastore = errors.New("unsupported datastore type")
)

type Builder interface {
	Build() (api.ModelRegistryApi, error)
	Teardown() error
}

func NewDatastore(dsType string, dsHostname string, dsPort int) (api.ModelRegistryApi, func() error, error) {
	switch dsType {
	case "mlmd":
		mlmd := NewMLMDService(dsHostname, dsPort)

		svc, err := mlmd.Build()
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %w", ErrCreatingDatastore, err)
		}

		return svc, mlmd.Teardown, nil
	case "inmemory":
		inmemory := NewInMemoryService()

		svc, err := inmemory.Build()
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %w", ErrCreatingDatastore, err)
		}

		return svc, inmemory.Teardown, nil
	default:
		return nil, nil, fmt.Errorf("%w: %s", ErrUnsupportedDatastore, dsType)
	}
}
