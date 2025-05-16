package db

import (
	"fmt"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"gorm.io/gorm"
)

type Connector interface {
	Connect() (*gorm.DB, error)
	DB() *gorm.DB
}

func NewConnector(dbType string, dsn string) (Connector, error) {
	switch dbType {
	case "mysql":
		return mysql.NewMySQLDBConnector(dsn), nil
	}

	return nil, fmt.Errorf("unsupported database type: %s", dbType)
}
