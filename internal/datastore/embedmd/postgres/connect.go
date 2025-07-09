package postgres

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	// postgresMaxRetries is the maximum number of attempts to retry PostgreSQL connection.
	postgresMaxRetries = 25 // 25 attempts with incremental backoff (1s, 2s, 3s, ..., 25s) it's ~5 minutes
)

type PostgresDBConnector struct {
	DSN          string
	db           *gorm.DB
	connectMutex sync.Mutex
}

func NewPostgresDBConnector(dsn string) *PostgresDBConnector {
	return &PostgresDBConnector{
		DSN: dsn,
	}
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

	glog.V(2).Infof("Attempting to connect with DSN: %q", c.DSN)

	for i := range postgresMaxRetries {
		db, err = gorm.Open(postgres.Open(c.DSN), &gorm.Config{
			Logger:         logger.Default.LogMode(logger.Silent),
			TranslateError: true,
		})
		if err == nil {
			break
		}

		glog.Warningf("Retrying connection to PostgreSQL (attempt %d/%d): %v", i+1, postgresMaxRetries, err)

		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	c.db = db
	glog.Info("Successfully connected to PostgreSQL database")
	return db, nil
}

func (c *PostgresDBConnector) DB() *gorm.DB {
	return c.db
}
