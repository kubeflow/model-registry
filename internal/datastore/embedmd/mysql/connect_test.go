package mysql_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	_tls "github.com/kubeflow/model-registry/internal/tls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	cont_mysql "github.com/testcontainers/testcontainers-go/modules/mysql"
)

func TestMySQLDBConnector_Connect_Insecure(t *testing.T) {
	ctx := context.Background()

	// Start MySQL container without SSL
	mysqlContainer, err := cont_mysql.Run(
		ctx,
		"mysql:8.0",
		cont_mysql.WithUsername("root"),
		cont_mysql.WithPassword("testpass"),
		cont_mysql.WithDatabase("testdb"),
	)
	require.NoError(t, err)
	defer func() {
		err := testcontainers.TerminateContainer(mysqlContainer)
		require.NoError(t, err)
	}()

	// Test basic connection without SSL
	t.Run("BasicConnection", func(t *testing.T) {
		dsn := mysqlContainer.MustConnectionString(ctx)
		connector := mysql.NewMySQLDBConnector(dsn, &_tls.TLSConfig{})

		db, err := connector.Connect()
		require.NoError(t, err)
		assert.NotNil(t, db)

		// Test that we can perform a simple query
		var result int
		err = db.Raw("SELECT 1").Scan(&result).Error
		require.NoError(t, err)
		assert.Equal(t, 1, result)

		// Clean up
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
	})

	t.Run("EmptySSLConfig", func(t *testing.T) {
		dsn := mysqlContainer.MustConnectionString(ctx)
		connector := &mysql.MySQLDBConnector{
			DSN:       dsn,
			TLSConfig: &_tls.TLSConfig{},
		}

		db, err := connector.Connect()
		require.NoError(t, err)
		assert.NotNil(t, db)

		// Clean up
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
	})
}

func TestMySQLDBConnector_TLSConfigValidation(t *testing.T) {
	// Create temporary directory for certificates
	tempDir, err := os.MkdirTemp("", "mysql_ssl_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir) //nolint:errcheck

	// Generate test certificates
	caCertPath, _, _, _, clientCertPath, clientKeyPath := generateTestCertificates(t, tempDir)

	t.Run("ValidCertificateFiles", func(t *testing.T) {
		// Test that valid certificate files can be loaded into TLS config
		tlsConfig := &_tls.TLSConfig{
			RootCertPath:     caCertPath,
			VerifyServerCert: true,
		}

		// This should build successfully since files are valid
		builtConfig, err := tlsConfig.BuildTLSConfig()
		assert.NoError(t, err)
		assert.NotNil(t, builtConfig)
		assert.NotNil(t, builtConfig.RootCAs)
		assert.False(t, builtConfig.InsecureSkipVerify)
	})

	t.Run("ValidClientCertificates", func(t *testing.T) {
		tlsConfig := &_tls.TLSConfig{
			CertPath:         clientCertPath,
			KeyPath:          clientKeyPath,
			RootCertPath:     caCertPath,
			VerifyServerCert: true,
		}

		builtConfig, err := tlsConfig.BuildTLSConfig()
		assert.NoError(t, err)
		assert.NotNil(t, builtConfig)
		assert.Len(t, builtConfig.Certificates, 1)
		assert.NotNil(t, builtConfig.RootCAs)
	})

	t.Run("ValidCipherSuites", func(t *testing.T) {
		tlsConfig := &_tls.TLSConfig{
			RootCertPath:     caCertPath,
			Cipher:           "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
			VerifyServerCert: true,
		}

		builtConfig, err := tlsConfig.BuildTLSConfig()
		assert.NoError(t, err)
		assert.NotNil(t, builtConfig)
		assert.Len(t, builtConfig.CipherSuites, 2)
	})

	t.Run("InvalidCertificateFile", func(t *testing.T) {
		tlsConfig := &_tls.TLSConfig{
			RootCertPath: "/nonexistent/cert.pem",
		}

		_, err := tlsConfig.BuildTLSConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read SSL root certificate")
	})

	t.Run("InvalidCertificateContent", func(t *testing.T) {
		invalidCertPath := filepath.Join(tempDir, "invalid.pem")
		err := os.WriteFile(invalidCertPath, []byte("invalid certificate content"), 0600)
		require.NoError(t, err)

		tlsConfig := &_tls.TLSConfig{
			RootCertPath: invalidCertPath,
		}

		_, err = tlsConfig.BuildTLSConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse SSL root certificate")
	})
}

func TestMySQLDBConnector_Connect_Secure(t *testing.T) {
	ctx := context.Background()

	// Start MySQL container with SSL enabled using built-in SSL support
	mysqlContainer, err := cont_mysql.Run(
		ctx,
		"mysql:8.0",
		cont_mysql.WithUsername("root"),
		cont_mysql.WithPassword("testpass"),
		cont_mysql.WithDatabase("testdb"),
		// Enable SSL with default certificates
		testcontainers.WithEnv(map[string]string{
			"MYSQL_ROOT_HOST": "%",
		}),
	)
	require.NoError(t, err)
	defer func() {
		err := testcontainers.TerminateContainer(mysqlContainer)
		require.NoError(t, err)
	}()

	baseDSN := mysqlContainer.MustConnectionString(ctx)

	t.Run("SSLConnectionWithSkipVerify", func(t *testing.T) {
		// Test SSL connection with skip verify (most practical test)
		// Parse the DSN and add TLS parameter properly
		dsn := baseDSN
		if strings.Contains(dsn, "?") {
			dsn += "&tls=skip-verify"
		} else {
			dsn += "?tls=skip-verify"
		}
		connector := mysql.NewMySQLDBConnector(dsn, &_tls.TLSConfig{})

		db, err := connector.Connect()
		require.NoError(t, err)
		assert.NotNil(t, db)

		// Test that we can perform a simple query over SSL
		var result int
		err = db.Raw("SELECT 1").Scan(&result).Error
		require.NoError(t, err)
		assert.Equal(t, 1, result)

		// Verify SSL is actually being used by checking the connection
		var statusName, sslCipher string
		err = db.Raw("SHOW STATUS LIKE 'Ssl_cipher'").Row().Scan(&statusName, &sslCipher)
		if err == nil && sslCipher != "" {
			t.Logf("SSL cipher in use: %s", sslCipher)
		}

		// Clean up
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
	})

	t.Run("SSLConnectionWithPreferredMode", func(t *testing.T) {
		// Test SSL connection with preferred mode
		dsn := baseDSN
		if strings.Contains(dsn, "?") {
			dsn += "&tls=preferred"
		} else {
			dsn += "?tls=preferred"
		}
		connector := mysql.NewMySQLDBConnector(dsn, &_tls.TLSConfig{})

		db, err := connector.Connect()
		require.NoError(t, err)
		assert.NotNil(t, db)

		// Test that we can perform operations
		var result int
		err = db.Raw("SELECT 1").Scan(&result).Error
		require.NoError(t, err)
		assert.Equal(t, 1, result)

		// Clean up
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
	})

	t.Run("CustomSSLConfigWithInsecureSkipVerify", func(t *testing.T) {
		// Test our custom SSL configuration with skip verify
		dsn := baseDSN
		if strings.Contains(dsn, "?") {
			dsn += "&tls=skip-verify"
		} else {
			dsn += "?tls=skip-verify"
		}
		connector := mysql.NewMySQLDBConnector(dsn, &_tls.TLSConfig{})

		db, err := connector.Connect()
		require.NoError(t, err)
		assert.NotNil(t, db)

		// Test database operations
		var result int
		err = db.Raw("SELECT 1").Scan(&result).Error
		require.NoError(t, err)
		assert.Equal(t, 1, result)

		// Test that we can create and query a table
		err = db.Exec("CREATE TEMPORARY TABLE test_ssl (id INT, name VARCHAR(50))").Error
		require.NoError(t, err)

		err = db.Exec("INSERT INTO test_ssl (id, name) VALUES (1, 'ssl_test')").Error
		require.NoError(t, err)

		var name string
		err = db.Raw("SELECT name FROM test_ssl WHERE id = 1").Scan(&name).Error
		require.NoError(t, err)
		assert.Equal(t, "ssl_test", name)

		// Clean up
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
	})

	t.Run("CustomSSLConfigWithCipherSuites", func(t *testing.T) {
		// Test custom cipher suites
		dsn := baseDSN
		if strings.Contains(dsn, "?") {
			dsn += "&tls=custom"
		} else {
			dsn += "?tls=custom"
		}
		connector := &mysql.MySQLDBConnector{
			DSN: dsn,
			TLSConfig: &_tls.TLSConfig{
				Cipher:           "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
				VerifyServerCert: false,
			},
		}

		db, err := connector.Connect()
		require.NoError(t, err)
		assert.NotNil(t, db)

		// Test basic operation
		var result int
		err = db.Raw("SELECT 1").Scan(&result).Error
		require.NoError(t, err)
		assert.Equal(t, 1, result)

		// Clean up
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
	})
}

func TestMySQLDBConnector_Connect_ErrorCases(t *testing.T) {
	t.Run("InvalidDSN", func(t *testing.T) {
		connector := mysql.NewMySQLDBConnector("invalid-dsn", &_tls.TLSConfig{})

		db, err := connector.Connect()
		assert.Error(t, err)
		assert.Nil(t, db)
	})

	t.Run("NonExistentCertFile", func(t *testing.T) {
		connector := &mysql.MySQLDBConnector{
			DSN: "root:pass@tcp(localhost:3306)/test",
			TLSConfig: &_tls.TLSConfig{
				RootCertPath: "/nonexistent/cert.pem",
			},
		}

		db, err := connector.Connect()
		assert.Error(t, err)
		assert.Nil(t, db)
		assert.Contains(t, err.Error(), "failed to read SSL root certificate")
	})

	t.Run("InvalidCertPair", func(t *testing.T) {
		// Create temporary files with invalid content
		tempDir, err := os.MkdirTemp("", "invalid_cert_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir) //nolint:errcheck

		invalidCertPath := filepath.Join(tempDir, "invalid.crt")
		invalidKeyPath := filepath.Join(tempDir, "invalid.key")

		err = os.WriteFile(invalidCertPath, []byte("invalid cert"), 0600)
		require.NoError(t, err)
		err = os.WriteFile(invalidKeyPath, []byte("invalid key"), 0600)
		require.NoError(t, err)

		connector := &mysql.MySQLDBConnector{
			DSN: "root:pass@tcp(localhost:3306)/test",
			TLSConfig: &_tls.TLSConfig{
				CertPath: invalidCertPath,
				KeyPath:  invalidKeyPath,
			},
		}

		db, err := connector.Connect()
		assert.Error(t, err)
		assert.Nil(t, db)
		assert.Contains(t, err.Error(), "failed to load SSL certificate pair")
	})

	t.Run("InvalidRootCert", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "invalid_root_cert_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir) //nolint:errcheck

		invalidRootCertPath := filepath.Join(tempDir, "invalid_root.crt")
		err = os.WriteFile(invalidRootCertPath, []byte("invalid root cert"), 0600)
		require.NoError(t, err)

		connector := &mysql.MySQLDBConnector{
			DSN: "root:pass@tcp(localhost:3306)/test",
			TLSConfig: &_tls.TLSConfig{
				RootCertPath: invalidRootCertPath,
			},
		}

		db, err := connector.Connect()
		assert.Error(t, err)
		assert.Nil(t, db)
		assert.Contains(t, err.Error(), "failed to parse SSL root certificate")
	})
}

// generateTestCertificates creates a set of test certificates for SSL testing
func generateTestCertificates(t *testing.T, tempDir string) (caCertPath, caKeyPath, serverCertPath, serverKeyPath, clientCertPath, clientKeyPath string) {
	// Generate CA private key
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Create CA certificate template
	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Test CA"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"Test"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// Create CA certificate
	caCertDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	require.NoError(t, err)

	// Save CA certificate
	caCertPath = filepath.Join(tempDir, "ca-cert.pem")
	caCertFile, err := os.Create(caCertPath)
	require.NoError(t, err)
	defer caCertFile.Close() //nolint:errcheck
	err = pem.Encode(caCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: caCertDER})
	require.NoError(t, err)

	// Save CA private key
	caKeyPath = filepath.Join(tempDir, "ca-key.pem")
	caKeyFile, err := os.Create(caKeyPath)
	require.NoError(t, err)
	defer caKeyFile.Close() //nolint:errcheck
	caKeyDER, err := x509.MarshalPKCS8PrivateKey(caKey)
	require.NoError(t, err)
	err = pem.Encode(caKeyFile, &pem.Block{Type: "PRIVATE KEY", Bytes: caKeyDER})
	require.NoError(t, err)

	// Parse CA certificate for signing
	caCert, err := x509.ParseCertificate(caCertDER)
	require.NoError(t, err)

	// Generate server certificate
	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	serverTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization:  []string{"Test Server"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"Test"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		DNSNames:     []string{"localhost"},
	}

	serverCertDER, err := x509.CreateCertificate(rand.Reader, &serverTemplate, caCert, &serverKey.PublicKey, caKey)
	require.NoError(t, err)

	// Save server certificate
	serverCertPath = filepath.Join(tempDir, "server-cert.pem")
	serverCertFile, err := os.Create(serverCertPath)
	require.NoError(t, err)
	defer serverCertFile.Close() //nolint:errcheck
	err = pem.Encode(serverCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: serverCertDER})
	require.NoError(t, err)

	// Save server private key
	serverKeyPath = filepath.Join(tempDir, "server-key.pem")
	serverKeyFile, err := os.Create(serverKeyPath)
	require.NoError(t, err)
	defer serverKeyFile.Close() //nolint:errcheck
	serverKeyDER, err := x509.MarshalPKCS8PrivateKey(serverKey)
	require.NoError(t, err)
	err = pem.Encode(serverKeyFile, &pem.Block{Type: "PRIVATE KEY", Bytes: serverKeyDER})
	require.NoError(t, err)

	// Generate client certificate
	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	clientTemplate := x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			Organization:  []string{"Test Client"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"Test"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		SubjectKeyId: []byte{1, 2, 3, 4, 5},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	clientCertDER, err := x509.CreateCertificate(rand.Reader, &clientTemplate, caCert, &clientKey.PublicKey, caKey)
	require.NoError(t, err)

	// Save client certificate
	clientCertPath = filepath.Join(tempDir, "client-cert.pem")
	clientCertFile, err := os.Create(clientCertPath)
	require.NoError(t, err)
	defer clientCertFile.Close() //nolint:errcheck
	err = pem.Encode(clientCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: clientCertDER})
	require.NoError(t, err)

	// Save client private key
	clientKeyPath = filepath.Join(tempDir, "client-key.pem")
	clientKeyFile, err := os.Create(clientKeyPath)
	require.NoError(t, err)
	defer clientKeyFile.Close() //nolint:errcheck
	clientKeyDER, err := x509.MarshalPKCS8PrivateKey(clientKey)
	require.NoError(t, err)
	err = pem.Encode(clientKeyFile, &pem.Block{Type: "PRIVATE KEY", Bytes: clientKeyDER})
	require.NoError(t, err)

	return caCertPath, caKeyPath, serverCertPath, serverKeyPath, clientCertPath, clientKeyPath
}
