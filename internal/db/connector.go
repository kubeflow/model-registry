package db

import (
	"fmt"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd/postgres"
	"gorm.io/gorm"
)

type Connector interface {
	Connect() (*gorm.DB, error)
	DB() *gorm.DB
}

func NewConnector(dbType string, dsn string) (Connector, error) {
	switch dbType {
	case embedmd.DatabaseTypeMySQL:
		return mysql.NewMySQLDBConnector(dsn), nil
	case embedmd.DatabaseTypePostgres:
		return postgres.NewPostgresDBConnector(dsn), nil
	}

	return nil, fmt.Errorf("unsupported database type: %s. Supported types: %s, %s", dbType, embedmd.DatabaseTypeMySQL, embedmd.DatabaseTypePostgres)
}
