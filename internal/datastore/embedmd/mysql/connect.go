package mysql

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLDBConnector struct {
	DSN                 string
	SSLCertPath         string
	SSLKeyPath          string
	SSLRootCertPath     string
	SSLCAPath           string
	SSLCipher           string
	SSLVerifyServerCert bool
	db                  *gorm.DB
}

func NewMySQLDBConnector(
	dsn,
	sslCertPath,
	sslKeyPath,
	sslRootCertPath,
	sslCAPath,
	sslCipher string,
	sslVerifyServerCert bool,
) *MySQLDBConnector {
	return &MySQLDBConnector{
		DSN:                 dsn,
		SSLCertPath:         sslCertPath,
		SSLKeyPath:          sslKeyPath,
		SSLRootCertPath:     sslRootCertPath,
		SSLCAPath:           sslCAPath,
		SSLCipher:           sslCipher,
		SSLVerifyServerCert: sslVerifyServerCert,
	}
}

func (c *MySQLDBConnector) Connect() (*gorm.DB, error) {
	if err := c.registerTLSConfig(); err != nil {
		return nil, err
	}

	db, err := gorm.Open(gorm_mysql.Open(c.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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

func (c *MySQLDBConnector) registerTLSConfig() error {
	var rootCAs *x509.CertPool
	var err error

	// Skip TLS registration only if no SSL configuration is provided at all
	// If SSLVerifyServerCert is false or SSLCipher is set, we should register a custom TLS config
	if c.SSLCertPath == "" && c.SSLKeyPath == "" && c.SSLRootCertPath == "" && c.SSLCAPath == "" && c.SSLCipher == "" && c.SSLVerifyServerCert {
		return nil
	}

	if c.SSLRootCertPath != "" || c.SSLCAPath != "" {
		rootCAs = x509.NewCertPool()

		if c.SSLRootCertPath != "" {
			rootCert, err := os.ReadFile(c.SSLRootCertPath)
			if err != nil {
				return fmt.Errorf("failed to read SSL root certificate from %s: %w", c.SSLRootCertPath, err)
			}
			if !rootCAs.AppendCertsFromPEM(rootCert) {
				return fmt.Errorf("failed to parse SSL root certificate from %s", c.SSLRootCertPath)
			}
		}

		if c.SSLCAPath != "" && c.SSLCAPath != c.SSLRootCertPath {
			caCert, err := os.ReadFile(c.SSLCAPath)
			if err != nil {
				return fmt.Errorf("failed to read SSL CA certificate from %s: %w", c.SSLCAPath, err)
			}
			if !rootCAs.AppendCertsFromPEM(caCert) {
				return fmt.Errorf("failed to parse SSL CA certificate from %s", c.SSLCAPath)
			}
		}
	} else if c.SSLVerifyServerCert {
		rootCAs, err = x509.SystemCertPool()
		if err != nil {
			return fmt.Errorf("failed to get system cert pool: %w", err)
		}
	}

	tlsConfig := &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: !c.SSLVerifyServerCert,
	}

	if c.SSLCertPath != "" && c.SSLKeyPath != "" {
		cert, err := tls.LoadX509KeyPair(c.SSLCertPath, c.SSLKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load SSL certificate pair (cert: %s, key: %s): %w",
				c.SSLCertPath, c.SSLKeyPath, err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	if c.SSLCipher != "" {
		cipherSuites := parseCipherSuites(c.SSLCipher)
		if len(cipherSuites) > 0 {
			tlsConfig.CipherSuites = cipherSuites
		}
	}

	if err := mysql.RegisterTLSConfig("custom", tlsConfig); err != nil {
		return fmt.Errorf("failed to register TLS config: %w", err)
	}

	return nil
}

// parseCipherSuites parses a comma-separated list of cipher suite names
// and returns the corresponding cipher suite IDs
func parseCipherSuites(cipherStr string) []uint16 {
	if cipherStr == "" {
		return nil
	}

	cipherMap := map[string]uint16{
		"TLS_RSA_WITH_RC4_128_SHA":                      tls.TLS_RSA_WITH_RC4_128_SHA,
		"TLS_RSA_WITH_3DES_EDE_CBC_SHA":                 tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		"TLS_RSA_WITH_AES_128_CBC_SHA":                  tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		"TLS_RSA_WITH_AES_256_CBC_SHA":                  tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		"TLS_RSA_WITH_AES_128_CBC_SHA256":               tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		"TLS_RSA_WITH_AES_128_GCM_SHA256":               tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		"TLS_RSA_WITH_AES_256_GCM_SHA384":               tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":              tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":          tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":          tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_RC4_128_SHA":                tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
		"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":           tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":            tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256":       tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":         tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":         tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256":       tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":         tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384":       tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256":   tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256": tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
		"TLS_AES_128_GCM_SHA256":                        tls.TLS_AES_128_GCM_SHA256,
		"TLS_AES_256_GCM_SHA384":                        tls.TLS_AES_256_GCM_SHA384,
		"TLS_CHACHA20_POLY1305_SHA256":                  tls.TLS_CHACHA20_POLY1305_SHA256,
	}

	var cipherSuites []uint16
	ciphers := strings.Split(strings.TrimSpace(cipherStr), ",")

	for _, cipher := range ciphers {
		cipher = strings.TrimSpace(cipher)
		if cipherID, exists := cipherMap[cipher]; exists {
			cipherSuites = append(cipherSuites, cipherID)
		}
	}

	return cipherSuites
}
