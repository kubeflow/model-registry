package db

import (
	"fmt"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"gorm.io/gorm"
)

type DBMigrator interface {
	Migrate() error
	Up(steps *int) error
	Down(steps *int) error
}

func NewDBMigrator(dbType string, db *gorm.DB) (DBMigrator, error) {
	switch dbType {
	case "mysql":
		return mysql.NewMySQLMigrator(db)
	}

	return nil, fmt.Errorf("unsupported database type: %s", dbType)
}
