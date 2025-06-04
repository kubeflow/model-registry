package db

import (
	"fmt"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd/postgres"
	"gorm.io/gorm"
)

type DBMigrator interface {
	Migrate() error
	Up(steps *int) error
	Down(steps *int) error
}

func NewDBMigrator(dbType string, db *gorm.DB) (DBMigrator, error) {
	switch dbType {
	case embedmd.DatabaseTypeMySQL:
		return mysql.NewMySQLMigrator(db)
	case embedmd.DatabaseTypePostgres:
		return postgres.NewPostgresMigrator(db)
	}

	return nil, fmt.Errorf("unsupported database type: %s. Supported types: %s, %s", dbType, embedmd.DatabaseTypeMySQL, embedmd.DatabaseTypePostgres)
}
