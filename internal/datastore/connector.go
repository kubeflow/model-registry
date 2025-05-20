package datastore

import (
	"errors"
	"fmt"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd"
	"github.com/kubeflow/model-registry/internal/datastore/mlmd"
	"github.com/kubeflow/model-registry/pkg/api"
)

var (
	ErrCreatingDatastore    = errors.New("error creating datastore")
	ErrUnsupportedDatastore = errors.New("unsupported datastore type")
)

type TeardownFunc func() error

type Datastore struct {
	MLMD    mlmd.MLMDConfig
	EmbedMD embedmd.EmbedMDConfig
	Type    string
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

		mlmd := mlmd.NewMLMDService(&ds.MLMD)

		return mlmd, nil
	case "embedmd":
		if err := ds.EmbedMD.Validate(); err != nil {
			return nil, fmt.Errorf("invalid EmbedMD config: %w", err)
		}

		embedmd, err := embedmd.NewEmbedMDService(&ds.EmbedMD)
		if err != nil {
			return nil, fmt.Errorf("error creating EmbedMD service: %w", err)
		}

		return embedmd, nil
	default:
		return nil, fmt.Errorf("%w: %s. Supported types: mlmd, embedmd", ErrUnsupportedDatastore, ds.Type)
	}
}
