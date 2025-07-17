package mysql_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	_tls "github.com/kubeflow/model-registry/internal/tls"
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

// Package-level shared database instance
var (
	sharedDB       *gorm.DB
	mysqlContainer *cont_mysql.MySQLContainer
)

// TestMain sets up the shared database instance and tears it down after all tests
func TestMain(m *testing.M) {
	ctx := context.Background()

	// Create MySQL container once for all tests
	container, err := cont_mysql.Run(
		ctx,
		"mysql:8",
		cont_mysql.WithUsername("root"),
		cont_mysql.WithPassword("root"),
		cont_mysql.WithDatabase("test"),
		cont_mysql.WithConfigFile(filepath.Join("testdata", "testdb.cnf")),
		// Enable SSL with default certificates
		testcontainers.WithEnv(map[string]string{
			"MYSQL_ROOT_HOST": "%",
		}),
	)
	if err != nil {
		panic("Failed to start MySQL container: " + err.Error())
	}
	mysqlContainer = container

	// Cleanup: close DB connection and terminate container
	defer func() {
		if sharedDB != nil {
			if sqlDB, err := sharedDB.DB(); err == nil {
				sqlDB.Close() //nolint:errcheck
			}
		}

		if mysqlContainer != nil {
			testcontainers.TerminateContainer(mysqlContainer) //nolint:errcheck
		}
	}()

	// Connect to the database
	dbConnector := mysql.NewMySQLDBConnector(mysqlContainer.MustConnectionString(ctx), &_tls.TLSConfig{})
	sharedDB, err = dbConnector.Connect()
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Run all tests
	code := m.Run()

	os.Exit(code)
}

// cleanupTestData truncates all tables to provide clean state between tests
func cleanupTestData(t *testing.T, db *gorm.DB) {
	// List of tables to clean up (in order to respect foreign key constraints)
	tables := []string{
		"TypeProperty",
		"Type",
		"MLMDEnv",
		"schema_migrations",
		// Add other tables as needed
	}

	// Disable foreign key checks temporarily
	err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error
	require.NoError(t, err)

	// Truncate all tables
	for _, table := range tables {
		err := db.Exec("TRUNCATE TABLE " + table).Error
		if err != nil {
			// Table might not exist, which is okay
			t.Logf("Could not truncate table %s: %v", table, err)
		}
	}

	// Re-enable foreign key checks
	err = db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
	require.NoError(t, err)
}

func TestMigrations(t *testing.T) {
	// Clean up any existing data
	cleanupTestData(t, sharedDB)

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
	cleanupTestData(t, sharedDB)

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
