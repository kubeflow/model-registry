package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDBConnector struct {
	DSN string
	db  *gorm.DB
}

func NewMySQLDBConnector(dsn string) *MySQLDBConnector {
	return &MySQLDBConnector{
		DSN: dsn,
	}
}

func (c *MySQLDBConnector) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	c.db = db

	return db, nil
}

func (c *MySQLDBConnector) DB() *gorm.DB {
	return c.db
}
