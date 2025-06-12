package mysql

import (
	"github.com/golang/glog"
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
	glog.V(2).Infof("Attempting to connect with DSN: %q", c.DSN)
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	glog.Info("Successfully connected to MySQL database")

	c.db = db

	return db, nil
}

func (c *MySQLDBConnector) DB() *gorm.DB {
	return c.db
}
