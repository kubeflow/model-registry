package postgres

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	_tls "github.com/kubeflow/model-registry/internal/tls"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	// postgresMaxRetriesDefault is the maximum number of attempts to retry PostgreSQL connection.
	postgresMaxRetriesDefault = 25 // 25 attempts with incremental backoff (1s, 2s, 3s, ..., 25s) it's ~5 minutes
)

type PostgresDBConnector struct {
	DSN          string
	TLSConfig    *_tls.TLSConfig
	db           *gorm.DB
	connectMutex sync.Mutex
	maxRetries   int
}

func NewPostgresDBConnector(
	dsn string,
	tlsConfig *_tls.TLSConfig,
) *PostgresDBConnector {
	return &PostgresDBConnector{
		DSN:        dsn,
		TLSConfig:  tlsConfig,
		maxRetries: postgresMaxRetriesDefault,
	}
}

func (c *PostgresDBConnector) WithMaxRetries(maxRetries int) *PostgresDBConnector {
	c.maxRetries = maxRetries

	return c
}

func (c *PostgresDBConnector) Connect() (*gorm.DB, error) {
	// Use mutex to ensure only one connection attempt at a time
	c.connectMutex.Lock()
	defer c.connectMutex.Unlock()

	// If we already have a working connection, return it
	if c.db != nil {
		return c.db, nil
	}

	var db *gorm.DB
	var err error

	dsn := c.DSN
	if c.needsTLSConfig() {
		dsn, err = c.BuildDSNWithTLS()
		if err != nil {
			return nil, fmt.Errorf("failed to build DSN with TLS: %w", err)
		}
	}

	for i := range c.maxRetries {
		glog.V(2).Infof("Attempting to connect with DSN: %q (attempt %d/%d)", dsn, i+1, c.maxRetries)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:         logger.Default.LogMode(logger.Silent),
			TranslateError: true,
		})
		if err == nil {
			break
		}

		glog.Warningf("Retrying connection to PostgreSQL (attempt %d/%d): %v", i+1, c.maxRetries, err)

		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	glog.Info("Successfully connected to PostgreSQL database")

	c.db = db

	return db, nil
}

func (c *PostgresDBConnector) DB() *gorm.DB {
	return c.db
}

func (c *PostgresDBConnector) needsTLSConfig() bool {
	if c.TLSConfig == nil {
		return false
	}

	// Log warning if cipher configuration is specified (not supported by PostgreSQL)
	if c.TLSConfig.Cipher != "" {
		glog.Warningf("SSL cipher configuration is not supported for PostgreSQL connections, ignoring cipher: %s", c.TLSConfig.Cipher)
	}

	return c.TLSConfig.CertPath != "" || c.TLSConfig.KeyPath != "" || c.TLSConfig.RootCertPath != "" || c.TLSConfig.CAPath != "" || c.TLSConfig.VerifyServerCert
}

// BuildDSNWithTLS builds a PostgreSQL DSN with SSL/TLS parameters based on the TLS configuration
func (c *PostgresDBConnector) BuildDSNWithTLS() (string, error) {
	if c.TLSConfig == nil {
		return c.DSN, nil
	}

	// Parse the existing DSN to determine format (URL or key=value pairs)
	if strings.HasPrefix(c.DSN, "postgres://") || strings.HasPrefix(c.DSN, "postgresql://") {
		return c.buildURLDSNWithTLS()
	}

	return c.buildKeyValueDSNWithTLS()
}

// buildURLDSNWithTLS handles URL-format DSNs (postgres://...)
func (c *PostgresDBConnector) buildURLDSNWithTLS() (string, error) {
	parsedURL, err := url.Parse(c.DSN)
	if err != nil {
		return "", fmt.Errorf("failed to parse PostgreSQL URL DSN: %w", err)
	}

	query := parsedURL.Query()

	// Set SSL mode based on TLS configuration
	if c.TLSConfig.VerifyServerCert {
		query.Set("sslmode", "verify-full")
	} else if c.TLSConfig.CertPath != "" || c.TLSConfig.KeyPath != "" || c.TLSConfig.RootCertPath != "" || c.TLSConfig.CAPath != "" {
		query.Set("sslmode", "require")
	}

	// Add certificate paths
	if c.TLSConfig.CertPath != "" {
		query.Set("sslcert", c.TLSConfig.CertPath)
	}
	if c.TLSConfig.KeyPath != "" {
		query.Set("sslkey", c.TLSConfig.KeyPath)
	}
	if c.TLSConfig.RootCertPath != "" {
		query.Set("sslrootcert", c.TLSConfig.RootCertPath)
	} else if c.TLSConfig.CAPath != "" {
		query.Set("sslrootcert", c.TLSConfig.CAPath)
	}

	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), nil
}

// buildKeyValueDSNWithTLS handles key=value format DSNs
func (c *PostgresDBConnector) buildKeyValueDSNWithTLS() (string, error) {
	dsn := c.DSN

	// Set SSL mode based on TLS configuration
	if c.TLSConfig.VerifyServerCert {
		dsn = c.addOrUpdateDSNParam(dsn, "sslmode", "verify-full")
	} else if c.TLSConfig.CertPath != "" || c.TLSConfig.KeyPath != "" || c.TLSConfig.RootCertPath != "" || c.TLSConfig.CAPath != "" {
		dsn = c.addOrUpdateDSNParam(dsn, "sslmode", "require")
	}

	// Add certificate paths
	if c.TLSConfig.CertPath != "" {
		dsn = c.addOrUpdateDSNParam(dsn, "sslcert", c.TLSConfig.CertPath)
	}
	if c.TLSConfig.KeyPath != "" {
		dsn = c.addOrUpdateDSNParam(dsn, "sslkey", c.TLSConfig.KeyPath)
	}
	if c.TLSConfig.RootCertPath != "" {
		dsn = c.addOrUpdateDSNParam(dsn, "sslrootcert", c.TLSConfig.RootCertPath)
	} else if c.TLSConfig.CAPath != "" {
		dsn = c.addOrUpdateDSNParam(dsn, "sslrootcert", c.TLSConfig.CAPath)
	}

	return dsn, nil
}

// addOrUpdateDSNParam adds or updates a parameter in a key=value DSN string
func (c *PostgresDBConnector) addOrUpdateDSNParam(dsn, key, value string) string {
	// Split DSN into individual parameters
	parts := strings.Fields(dsn)
	keyPrefix := key + "="
	updated := false

	// Update existing parameter or collect all parts
	for i, part := range parts {
		if strings.HasPrefix(part, keyPrefix) {
			parts[i] = keyPrefix + value
			updated = true
			break
		}
	}

	// Add new parameter if it didn't exist
	if !updated {
		parts = append(parts, keyPrefix+value)
	}

	return strings.Join(parts, " ")
}
