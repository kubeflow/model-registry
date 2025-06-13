package embedmd

import (
	"fmt"

	"github.com/kubeflow/model-registry/internal/db"
	"github.com/kubeflow/model-registry/internal/db/types"
	"github.com/kubeflow/model-registry/pkg/api"
)

type EmbedMDConfig struct {
	DatabaseType string
	DatabaseDSN  string
}

func (c *EmbedMDConfig) Validate() error {
	if c.DatabaseType != types.DatabaseTypeMySQL && c.DatabaseType != types.DatabaseTypePostgres {
		return fmt.Errorf("unsupported database type: %s. Supported types: %s, %s", c.DatabaseType, types.DatabaseTypeMySQL, types.DatabaseTypePostgres)
	}

	return nil
}

type EmbedMDService struct {
	*EmbedMDConfig
	dbConnector db.Connector
}

func NewEmbedMDService(cfg *EmbedMDConfig) (*EmbedMDService, error) {
	dbConnector, err := db.NewConnector(cfg.DatabaseType, cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	return &EmbedMDService{
		EmbedMDConfig: cfg,
		dbConnector:   dbConnector,
	}, nil
}

func (s *EmbedMDService) Connect() (api.ModelRegistryApi, error) {
	connectedDB, err := s.dbConnector.Connect()
	if err != nil {
		return nil, err
	}

	migrator, err := db.NewDBMigrator(s.DatabaseType, connectedDB)
	if err != nil {
		return nil, err
	}

	err = migrator.Migrate()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *EmbedMDService) Teardown() error {
	return nil
}
