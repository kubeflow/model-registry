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
	connectorOnce      sync.Once
	connectorMutex     sync.RWMutex
)

func NewConnector(dbType string, dsn string, tlsConfig *tls.TLSConfig) (Connector, error) {
	connectorMutex.RLock()
	if _connectorInstance != nil {
		connectorMutex.RUnlock()
		return _connectorInstance, nil
	}
	connectorMutex.RUnlock()

	var err error
	connectorOnce.Do(func() {
		connectorMutex.Lock()
		defer connectorMutex.Unlock()

		switch dbType {
		case "mysql":
			if tlsConfig != nil {
				_connectorInstance = mysql.NewMySQLDBConnector(dsn, tlsConfig)
			} else {
				_connectorInstance = mysql.NewMySQLDBConnector(dsn, &tls.TLSConfig{})
			}
		case "postgres":
			_connectorInstance = postgres.NewPostgresDBConnector(dsn)
		default:
			err = fmt.Errorf("unsupported database type: %s. Supported types: %s, %s", dbType, types.DatabaseTypeMySQL, types.DatabaseTypePostgres)
		}
	})

	if err != nil {
		return nil, err
	}

	return _connectorInstance, nil
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
	connectorOnce = sync.Once{}
}
