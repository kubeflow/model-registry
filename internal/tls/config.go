package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
)

type TLSConfig struct {
	CertPath         string
	KeyPath          string
	RootCertPath     string
	CAPath           string
	Cipher           string
	VerifyServerCert bool
}

func NewTLSConfig(certPath, keyPath, rootCertPath, caPath, cipher string, verifyServerCert bool) *TLSConfig {
	return &TLSConfig{
		CertPath:         certPath,
		KeyPath:          keyPath,
		RootCertPath:     rootCertPath,
		CAPath:           caPath,
		Cipher:           cipher,
		VerifyServerCert: verifyServerCert,
	}
}

func (c *TLSConfig) BuildTLSConfig() (*tls.Config, error) {
	var rootCAs *x509.CertPool
	var err error

	if c.RootCertPath != "" || c.CAPath != "" {
		rootCAs = x509.NewCertPool()

		if c.RootCertPath != "" {
			rootCert, err := os.ReadFile(c.RootCertPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read SSL root certificate from %s: %w", c.RootCertPath, err)
			}
			if !rootCAs.AppendCertsFromPEM(rootCert) {
				return nil, fmt.Errorf("failed to parse SSL root certificate from %s", c.RootCertPath)
			}
		}

		if c.CAPath != "" && c.CAPath != c.RootCertPath {
			caCert, err := os.ReadFile(c.CAPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read SSL CA certificate from %s: %w", c.CAPath, err)
			}
			if !rootCAs.AppendCertsFromPEM(caCert) {
				return nil, fmt.Errorf("failed to parse SSL CA certificate from %s", c.CAPath)
			}
		}
	} else if c.VerifyServerCert {
		rootCAs, err = x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("failed to get system cert pool: %w", err)
		}
	}

	tlsConfig := &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: !c.VerifyServerCert,
	}

	if c.CertPath != "" && c.KeyPath != "" {
		cert, err := tls.LoadX509KeyPair(c.CertPath, c.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load SSL certificate pair (cert: %s, key: %s): %w",
				c.CertPath, c.KeyPath, err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	if c.Cipher != "" {
		cipherSuites, err := parseCipherSuites(c.Cipher)
		if err != nil {
			return nil, err
		}

		if len(cipherSuites) > 0 {
			tlsConfig.CipherSuites = cipherSuites
		}
	}

	return tlsConfig, nil
}

// parseCipherSuites parses a colon-separated list of cipher suite names
// and returns the corresponding cipher suite IDs
func parseCipherSuites(cipherStr string) ([]uint16, error) {
	if cipherStr == "" {
		return nil, nil
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
	ciphers := strings.Split(strings.TrimSpace(cipherStr), ":")

	for _, cipher := range ciphers {
		cipher = strings.TrimSpace(cipher)
		if cipherID, exists := cipherMap[cipher]; exists {
			cipherSuites = append(cipherSuites, cipherID)
		} else {
			return nil, fmt.Errorf("invalid cipher suite: %s", cipher)
		}
	}

	return cipherSuites, nil
}
