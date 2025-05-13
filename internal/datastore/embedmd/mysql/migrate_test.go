package mysql_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	cont_mysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/gorm"
)

// Type represents the Type table structure
type Type struct {
	ID          int64  `gorm:"primaryKey"`
	Name        string `gorm:"column:name"`
	Version     string `gorm:"column:version"`
	ExternalID  string `gorm:"column:external_id"`
	Description string `gorm:"column:description"`
}

func (Type) TableName() string {
	return "Type"
}

// TypeProperty represents the TypeProperty table structure
type TypeProperty struct {
	ID          int64  `gorm:"primaryKey"`
	TypeID      int64  `gorm:"column:type_id"`
	Name        string `gorm:"column:name"`
	DataType    string `gorm:"column:data_type"`
	Description string `gorm:"column:description"`
}

func (TypeProperty) TableName() string {
	return "TypeProperty"
}

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	mysqlContainer, err := cont_mysql.Run(
		ctx,
		"mysql:5.7",
		cont_mysql.WithUsername("root"),
		cont_mysql.WithPassword("root"),
		cont_mysql.WithDatabase("test"),
		cont_mysql.WithConfigFile(filepath.Join("testdata", "testdb.cnf")),
	)
	require.NoError(t, err)

	dbConnector := mysql.NewMySQLDBConnector(mysqlContainer.MustConnectionString(ctx))
	require.NoError(t, err)

	db, err := dbConnector.Connect()
	require.NoError(t, err)

	// Return cleanup function
	cleanup := func() {
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
		err = testcontainers.TerminateContainer(
			mysqlContainer,
		)
		require.NoError(t, err)
	}

	return db, cleanup
}

func TestMigrations(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create migrator
	migrator, err := mysql.NewMySQLMigrator(db)
	require.NoError(t, err)

	// Run migrations
	err = migrator.Migrate()
	require.NoError(t, err)

	// Verify MLMDEnv table
	var schemaVersion int
	err = db.Raw("SELECT schema_version FROM MLMDEnv LIMIT 1").Scan(&schemaVersion).Error
	require.NoError(t, err)
	assert.Equal(t, 10, schemaVersion)

	// Verify Type table has expected entries
	var count int64
	err = db.Model(&Type{}).Count(&count).Error
	require.NoError(t, err)
	assert.Greater(t, count, int64(0))

	// Verify TypeProperty table has expected entries
	err = db.Model(&TypeProperty{}).Count(&count).Error
	require.NoError(t, err)
	assert.Greater(t, count, int64(0))
}

func TestDownMigrations(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	migrator, err := mysql.NewMySQLMigrator(db)
	require.NoError(t, err)

	// Run migrations
	err = migrator.Migrate()
	require.NoError(t, err)

	// Down migrations
	err = migrator.Down(nil)
	require.NoError(t, err)

	// Verify tables don't exist (except schema_migrations)
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name != 'schema_migrations'").Scan(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
