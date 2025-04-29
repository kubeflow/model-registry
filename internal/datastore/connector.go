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

type TeardownFunc func() error

type Datastore struct {
	MLMD MLMDConfig
	Type string
}

type Connector interface {
	Connect() (api.ModelRegistryApi, error)
	Teardown() error
}

func NewConnector(ds Datastore) (Connector, error) {
	switch ds.Type {
	case "mlmd":
		if err := ds.MLMD.Validate(); err != nil {
			return nil, fmt.Errorf("invalid MLMD config: %w", err)
		}

		mlmd := NewMLMDService(&ds.MLMD)

		return mlmd, nil
	default:
		return nil, fmt.Errorf("%w: %s. Supported types: mlmd", ErrUnsupportedDatastore, ds.Type)
	}
}
