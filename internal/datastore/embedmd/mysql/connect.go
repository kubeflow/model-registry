package mysql

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	_tls "github.com/kubeflow/model-registry/internal/tls"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLDBConnector struct {
	DSN       string
	TLSConfig *_tls.TLSConfig
	db        *gorm.DB
}

func NewMySQLDBConnector(
	dsn string,
	tlsConfig *_tls.TLSConfig,
) *MySQLDBConnector {
	return &MySQLDBConnector{
		DSN:       dsn,
		TLSConfig: tlsConfig,
	}
}

func (c *MySQLDBConnector) Connect() (*gorm.DB, error) {
	if c.needsTLSConfig() {
		if err := c.registerTLSConfig(); err != nil {
			return nil, err
		}

		c.DSN += "&tls=custom"
	}

	db, err := gorm.Open(gorm_mysql.Open(c.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

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
