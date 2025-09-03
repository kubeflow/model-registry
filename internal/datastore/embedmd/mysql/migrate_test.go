package mysql_test

import (
	"os"
	"testing"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestMain(m *testing.M) {
	os.Exit(testutils.TestMainHelper(m))
}

func TestMigrations(t *testing.T) {
	sharedDB, cleanup := testutils.GetSharedMySQLDB(t)
	defer cleanup()

	// Create migrator
	migrator, err := mysql.NewMySQLMigrator(sharedDB)
	require.NoError(t, err)

	// Run migrations
	err = migrator.Migrate()
	require.NoError(t, err)

	// Verify MLMDEnv table
	var schemaVersion int
	err = sharedDB.Raw("SELECT schema_version FROM MLMDEnv LIMIT 1").Scan(&schemaVersion).Error
	require.NoError(t, err)
	assert.Equal(t, 10, schemaVersion)

	// Verify Type table has expected entries
	var count int64
	err = sharedDB.Model(&Type{}).Count(&count).Error
	require.NoError(t, err)
	assert.Greater(t, count, int64(0))

	// Verify TypeProperty table has expected entries
	err = sharedDB.Model(&TypeProperty{}).Count(&count).Error
	require.NoError(t, err)
	assert.Greater(t, count, int64(0))
}

func TestDownMigrations(t *testing.T) {
	// Clean up any existing data
	sharedDB, cleanup := testutils.GetSharedMySQLDB(t)
	defer cleanup()

	migrator, err := mysql.NewMySQLMigrator(sharedDB)
	require.NoError(t, err)

	// Run migrations first
	err = migrator.Migrate()
	require.NoError(t, err)

	// Down migrations
	err = migrator.Down(nil)
	require.NoError(t, err)

	// Verify tables don't exist (except schema_migrations)
	var count int64
	err = sharedDB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name != 'schema_migrations'").Scan(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
