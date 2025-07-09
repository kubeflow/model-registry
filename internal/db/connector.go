package db

import (
	"fmt"
	"sync"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd/postgres"
	"github.com/kubeflow/model-registry/internal/db/types"
	"github.com/kubeflow/model-registry/internal/tls"
	"gorm.io/gorm"
)

type Connector interface {
	Connect() (*gorm.DB, error)
	DB() *gorm.DB
}

var (
	_connectorInstance Connector
	connectorMutex     sync.RWMutex
)

func Init(dbType string, dsn string, tlsConfig *tls.TLSConfig) error {
	connectorMutex.Lock()
	defer connectorMutex.Unlock()

	if tlsConfig == nil {
		tlsConfig = &tls.TLSConfig{}
	}

	switch dbType {
	case "mysql":
		_connectorInstance = mysql.NewMySQLDBConnector(dsn, tlsConfig)
	case "postgres":
		_connectorInstance = postgres.NewPostgresDBConnector(dsn, tlsConfig)
	default:
		return fmt.Errorf("unsupported database type: %s. Supported types: %s, %s", dbType, types.DatabaseTypeMySQL, types.DatabaseTypePostgres)
	}

	return nil
}

func GetConnector() (Connector, bool) {
	connectorMutex.RLock()
	defer connectorMutex.RUnlock()

	return _connectorInstance, _connectorInstance != nil
}

func ClearConnector() {
	connectorMutex.Lock()
	defer connectorMutex.Unlock()

	_connectorInstance = nil
}

