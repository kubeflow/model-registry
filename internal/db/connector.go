package db

import (
	"fmt"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/kubeflow/model-registry/internal/tls"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd/postgres"
	"github.com/kubeflow/model-registry/internal/db/types"
	"gorm.io/gorm"
)

type Connector interface {
	Connect() (*gorm.DB, error)
	DB() *gorm.DB
}

func NewConnector(dbType string, dsn string, tlsConfig *tls.TLSConfig) (Connector, error) {
	switch dbType {
	case "mysql":
		if tlsConfig != nil {
			return mysql.NewMySQLDBConnector(
				dsn,
				tlsConfig,
			), nil
		}

		return mysql.NewMySQLDBConnector(dsn, &tls.TLSConfig{}), nil
	case "postgres":
		return postgres.NewPostgresDBConnector(dsn), nil
	}

	return nil, fmt.Errorf("unsupported database type: %s. Supported types: %s, %s", dbType, types.DatabaseTypeMySQL, types.DatabaseTypePostgres)
}
