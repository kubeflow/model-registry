package db

import (
	"fmt"
	"slices"
	"sync"

	"github.com/kubeflow/model-registry/internal/platform/tls"
	"gorm.io/gorm"
)

type Connector interface {
	Connect() (*gorm.DB, error)
	DB() *gorm.DB
}

type ConnectorFactory func(dsn string, tlsConfig *tls.TLSConfig) Connector

var (
	_connectorInstance Connector
	connectorMutex     sync.RWMutex

	connectorFactories   = make(map[string]ConnectorFactory)
	connectorFactoriesMu sync.RWMutex
)

func RegisterConnectorFactory(dbType string, factory ConnectorFactory) {
	connectorFactoriesMu.Lock()
	defer connectorFactoriesMu.Unlock()

	if _, exists := connectorFactories[dbType]; exists {
		panic(fmt.Sprintf("duplicate connector factory for database type %q", dbType))
	}

	connectorFactories[dbType] = factory
}

func Init(dbType string, dsn string, tlsConfig *tls.TLSConfig) error {
	connectorMutex.Lock()
	defer connectorMutex.Unlock()

	if tlsConfig == nil {
		tlsConfig = &tls.TLSConfig{}
	}

	connectorFactoriesMu.RLock()
	factory, ok := connectorFactories[dbType]
	connectorFactoriesMu.RUnlock()
	if !ok {
		return fmt.Errorf("unsupported database type: %s. Registered types: %v", dbType, registeredConnectorTypes())
	}

	_connectorInstance = factory(dsn, tlsConfig)
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

type ConnectedConnector struct {
	ConnectedDB *gorm.DB
}

func (c ConnectedConnector) Connect() (*gorm.DB, error) {
	return c.ConnectedDB, nil
}

func (c ConnectedConnector) DB() *gorm.DB {
	return c.ConnectedDB
}

func registeredConnectorTypes() []string {
	connectorFactoriesMu.RLock()
	defer connectorFactoriesMu.RUnlock()

	types := make([]string, 0, len(connectorFactories))
	for dbType := range connectorFactories {
		types = append(types, dbType)
	}
	slices.Sort(types)
	return types
}
