package mysql

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	_tls "github.com/kubeflow/model-registry/internal/tls"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	// mysqlMaxRetriesDefault is the maximum number of attempts to retry MySQL connection.
	mysqlMaxRetriesDefault = 25 // 25 attempts with incremental backoff (1s, 2s, 3s, ..., 25s) it's ~5 minutes
)

type MySQLDBConnector struct {
	DSN          string
	TLSConfig    *_tls.TLSConfig
	db           *gorm.DB
	connectMutex sync.Mutex
	maxRetries   int
}

func NewMySQLDBConnector(
	dsn string,
	tlsConfig *_tls.TLSConfig,
) *MySQLDBConnector {
	return &MySQLDBConnector{
		DSN:        dsn,
		TLSConfig:  tlsConfig,
		maxRetries: mysqlMaxRetriesDefault,
	}
}

func (c *MySQLDBConnector) WithMaxRetries(maxRetries int) *MySQLDBConnector {
	c.maxRetries = maxRetries

	return c
}

func (c *MySQLDBConnector) Connect() (*gorm.DB, error) {
	// Use mutex to ensure only one connection attempt at a time
	c.connectMutex.Lock()
	defer c.connectMutex.Unlock()

	// If we already have a working connection, return it
	if c.db != nil {
		return c.db, nil
	}

	var db *gorm.DB
	var err error

	if c.needsTLSConfig() {
		if err := c.registerTLSConfig(); err != nil {
			return nil, err
		}

		cfg, err := mysql.ParseDSN(c.DSN)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DSN: %w", err)
		}

		cfg.TLSConfig = "custom"

		c.DSN = cfg.FormatDSN()
	}

	for i := range c.maxRetries {
		db, err = gorm.Open(gorm_mysql.Open(c.DSN), &gorm.Config{
			Logger:         logger.Default.LogMode(logger.Silent),
			TranslateError: true,
		})
		if err == nil {
			break
		}

		glog.Warningf("Retrying connection to MySQL (attempt %d/%d): %v", i+1, c.maxRetries, err)

		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	glog.Info("Successfully connected to MySQL database")

	c.db = db

	return db, nil
}

func (c *MySQLDBConnector) DB() *gorm.DB {
	return c.db
}

func (c *MySQLDBConnector) needsTLSConfig() bool {
	return c.TLSConfig != nil && (c.TLSConfig.CertPath != "" || c.TLSConfig.KeyPath != "" || c.TLSConfig.RootCertPath != "" || c.TLSConfig.CAPath != "" || c.TLSConfig.Cipher != "" || c.TLSConfig.VerifyServerCert)
}

func (c *MySQLDBConnector) registerTLSConfig() error {
	tlsConfig, err := c.TLSConfig.BuildTLSConfig()
	if err != nil {
		return err
	}

	if err := mysql.RegisterTLSConfig("custom", tlsConfig); err != nil {
		return fmt.Errorf("failed to register TLS config: %w", err)
	}

	return nil
}
