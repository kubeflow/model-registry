package postgres

import (
	"fmt"

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
	fmt.Printf("Attempting to connect with DSN: %q\n", c.DSN)
	db, err := gorm.Open(postgres.Open(c.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	c.db = db
	fmt.Printf("Successfully connected to PostgreSQL database\n")
	return db, nil
} 

func (c *PostgresDBConnector) DB() *gorm.DB {
	return c.db
}