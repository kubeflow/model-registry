package db

import (
	"fmt"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd/postgres"
	"github.com/kubeflow/model-registry/internal/db/types"
	"gorm.io/gorm"
)

type Connector interface {
	Connect() (*gorm.DB, error)
	DB() *gorm.DB
}

func NewConnector(dbType string, dsn string) (Connector, error) {
	switch dbType {
	case types.DatabaseTypeMySQL:
		return mysql.NewMySQLDBConnector(dsn), nil
	case types.DatabaseTypePostgres:
		return postgres.NewPostgresDBConnector(dsn), nil
	}

	return nil, fmt.Errorf("unsupported database type: %s. Supported types: %s, %s", dbType, types.DatabaseTypeMySQL, types.DatabaseTypePostgres)
}
