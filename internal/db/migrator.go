package db

import (
	"fmt"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd/postgres"
	"github.com/kubeflow/model-registry/internal/db/types"
	"gorm.io/gorm"
)

type DBMigrator interface {
	Migrate() error
	Up(steps *int) error
	Down(steps *int) error
}

func NewDBMigrator(db *gorm.DB) (DBMigrator, error) {
	switch db.Name() {
	case types.DatabaseTypeMySQL:
		return mysql.NewMySQLMigrator(db)
	case types.DatabaseTypePostgres:
		return postgres.NewPostgresMigrator(db)
	}

	return nil, fmt.Errorf("unsupported database type: %s. Supported types: %s, %s", db.Name(), types.DatabaseTypeMySQL, types.DatabaseTypePostgres)
}
