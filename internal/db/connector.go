package db

import (
	"fmt"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"gorm.io/gorm"
)

type Connector interface {
	Connect() (*gorm.DB, error)
	DB() *gorm.DB
}

type SSLConfig struct {
	SSLCert             string
	SSLKey              string
	SSLRootCert         string
	SSLCA               string
	SSLCipher           string
	SSLVerifyServerCert bool
}

func NewConnector(dbType string, dsn string, sslConfig *SSLConfig) (Connector, error) {
	switch dbType {
	case "mysql":
		if sslConfig != nil {
			return mysql.NewMySQLDBConnector(
				dsn,
				sslConfig.SSLCert,
				sslConfig.SSLKey,
				sslConfig.SSLRootCert,
				sslConfig.SSLCA,
				sslConfig.SSLCipher,
				sslConfig.SSLVerifyServerCert,
			), nil
		}

		return mysql.NewMySQLDBConnector(dsn, "", "", "", "", "", false), nil
	}

	return nil, fmt.Errorf("unsupported database type: %s", dbType)
}
