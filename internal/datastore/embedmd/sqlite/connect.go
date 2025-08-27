package sqlite

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	_tls "github.com/kubeflow/model-registry/internal/tls"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	// sqliteMaxRetriesDefault is the maximum number of attempts to retry SQLite connection.
	sqliteMaxRetriesDefault = 5 // SQLite is file-based, so fewer retries needed
)

type SQLiteDBConnector struct {
	DSN          string
	TLSConfig    *_tls.TLSConfig // Not used for SQLite but kept for interface consistency
	db           *gorm.DB
	connectMutex sync.Mutex
	maxRetries   int
}

func NewSQLiteDBConnector(
	dsn string,
	tlsConfig *_tls.TLSConfig,
) *SQLiteDBConnector {
	return &SQLiteDBConnector{
		DSN:        dsn,
		TLSConfig:  tlsConfig,
		maxRetries: sqliteMaxRetriesDefault,
	}
}

func (c *SQLiteDBConnector) WithMaxRetries(maxRetries int) *SQLiteDBConnector {
	c.maxRetries = maxRetries

	return c
}

func (c *SQLiteDBConnector) Connect() (*gorm.DB, error) {
	// Use mutex to ensure only one connection attempt at a time
	c.connectMutex.Lock()
	defer c.connectMutex.Unlock()

	// If we already have a working connection, return it
	if c.db != nil {
		return c.db, nil
	}

	var db *gorm.DB
	var err error

	// Log warning if TLS configuration is specified (not supported by SQLite)
	if c.TLSConfig != nil && c.needsTLSConfig() {
		glog.Warningf("TLS configuration is not supported for SQLite connections, ignoring TLS settings")
	}

	for i := range c.maxRetries {
		glog.V(2).Infof("Attempting to connect to SQLite database: %q (attempt %d/%d)", c.DSN, i+1, c.maxRetries)
		db, err = gorm.Open(sqlite.Open(c.DSN), &gorm.Config{
			Logger:         logger.Default.LogMode(logger.Silent),
			TranslateError: true,
		})
		if err == nil {
			break
		}

		glog.Warningf("Retrying connection to SQLite (attempt %d/%d): %v", i+1, c.maxRetries, err)

		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
	}

	glog.Info("Successfully connected to SQLite database")

	c.db = db

	return db, nil
}

func (c *SQLiteDBConnector) DB() *gorm.DB {
	return c.db
}

func (c *SQLiteDBConnector) needsTLSConfig() bool {
	if c.TLSConfig == nil {
		return false
	}

	return c.TLSConfig.CertPath != "" || c.TLSConfig.KeyPath != "" || c.TLSConfig.RootCertPath != "" || c.TLSConfig.CAPath != "" || c.TLSConfig.Cipher != "" || c.TLSConfig.VerifyServerCert
}