package embedmd

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/db"
	"github.com/kubeflow/model-registry/internal/db/types"
	"github.com/kubeflow/model-registry/internal/tls"
)

const connectorType = "embedmd"

func init() {
	datastore.Register(connectorType, func(cfg any) (datastore.Connector, error) {
		emdbCfg, ok := cfg.(*EmbedMDConfig)
		if !ok {
			return nil, fmt.Errorf("invalid EmbedMD config type (%T)", cfg)
		}

		if err := emdbCfg.Validate(); err != nil {
			return nil, fmt.Errorf("invalid EmbedMD config: %w", err)
		}

		return NewEmbedMDService(cfg.(*EmbedMDConfig))
	})
}

type EmbedMDConfig struct {
	DatabaseType string
	DatabaseDSN  string
	TLSConfig    *tls.TLSConfig
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
	err := db.Init(cfg.DatabaseType, cfg.DatabaseDSN, cfg.TLSConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database connector: %w", err)
	}

	dbConnector, ok := db.GetConnector()
	if !ok {
		return nil, fmt.Errorf("database connector not initialized")
	}

	return &EmbedMDService{
		EmbedMDConfig: cfg,
		dbConnector:   dbConnector,
	}, nil
}

func (s *EmbedMDService) Connect(spec datastore.RepoSetSpec) (datastore.RepoSet, error) {
	glog.Infof("Connecting to EmbedMD service...")

	connectedDB, err := s.dbConnector.Connect()
	if err != nil {
		return nil, err
	}

	glog.Infof("Connected to EmbedMD service")

	migrator, err := db.NewDBMigrator(s.DatabaseType, connectedDB)
	if err != nil {
		return nil, err
	}

	glog.Infof("Running migrations...")

	err = migrator.Migrate()
	if err != nil {
		return nil, err
	}

	glog.Infof("Migrations completed")

	return newRepoSet(connectedDB, spec)
}

func (s EmbedMDService) Type() string {
	return connectorType
}
