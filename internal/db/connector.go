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
	case types.DatabaseTypeMySQL:
		_connectorInstance = mysql.NewMySQLDBConnector(dsn, tlsConfig)
	case types.DatabaseTypePostgres:
		_connectorInstance = postgres.NewPostgresDBConnector(dsn, tlsConfig)
	default:
		return fmt.Errorf("unsupported database type: %s. Supported types: %s, %s", dbType, types.DatabaseTypeMySQL, types.DatabaseTypePostgres)
	}

	return nil
}

func SetDB(connectedDB *gorm.DB) {
	connectorMutex.Lock()
	defer connectorMutex.Unlock()
	_connectorInstance = ConnectedConnector{ConnectedDB: connectedDB}
}

func GetConnector() Connector {
	connectorMutex.RLock()
	defer connectorMutex.RUnlock()

	return _connectorInstance
}

func ClearConnector() {
	connectorMutex.Lock()
	defer connectorMutex.Unlock()

	_connectorInstance = nil
}

// ConnectedConnector satifies the connector interface for an already connected
// gorm.DB instance.
type ConnectedConnector struct {
	ConnectedDB *gorm.DB
}

func (c ConnectedConnector) Connect() (*gorm.DB, error) {
	return c.ConnectedDB, nil
}

func (c ConnectedConnector) DB() *gorm.DB {
	return c.ConnectedDB
}
