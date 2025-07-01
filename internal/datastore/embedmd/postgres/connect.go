package postgres

import (
	"github.com/golang/glog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDBConnector struct {
	DSN string
	db  *gorm.DB
}

func NewPostgresDBConnector(dsn string) *PostgresDBConnector {
	return &PostgresDBConnector{
		DSN: dsn,
	}
}

func (c *PostgresDBConnector) Connect() (*gorm.DB, error) {
	glog.V(2).Infof("Attempting to connect with DSN: %q", c.DSN)
	db, err := gorm.Open(postgres.Open(c.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	c.db = db
	glog.Info("Successfully connected to PostgreSQL database")
	return db, nil
} 

func (c *PostgresDBConnector) DB() *gorm.DB {
	return c.db
}