package sqlite_test

import (
	"path/filepath"
	"testing"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlitedriver "gorm.io/driver/sqlite"
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

// setupTestDB creates a temporary SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	// Create temporary file for SQLite database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	// Create GORM database connection
	db, err := gorm.Open(sqlitedriver.Open(dbPath), &gorm.Config{
		TranslateError: true,
	})
	require.NoError(t, err)

	// Ensure cleanup happens properly
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	})

	return db
}

func TestMigrations(t *testing.T) {
	db := setupTestDB(t)

	// Create migrator
	migrator, err := sqlite.NewSQLiteMigrator(db)
	require.NoError(t, err)

	// Run migrations
	err = migrator.Migrate()
	require.NoError(t, err)

	// Verify MLMDEnv table exists and has expected data
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

	// Verify at least one type exists by checking if any record exists
	var typeExists bool
	err = db.Raw("SELECT EXISTS(SELECT 1 FROM Type LIMIT 1)").Scan(&typeExists).Error
	require.NoError(t, err)
	assert.True(t, typeExists)
}

func TestDownMigrations(t *testing.T) {
	db := setupTestDB(t)

	migrator, err := sqlite.NewSQLiteMigrator(db)
	require.NoError(t, err)

	// Run migrations first
	err = migrator.Migrate()
	require.NoError(t, err)

	// Verify tables exist before down migration
	var tableCount int64
	err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name != 'schema_migrations'").Scan(&tableCount).Error
	require.NoError(t, err)
	assert.Greater(t, tableCount, int64(0))

	// Run down migrations
	err = migrator.Down(nil)
	require.NoError(t, err)

	// Verify most tables are dropped (should be significantly fewer than before)
	var afterTableCount int64
	err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name != 'schema_migrations'").Scan(&afterTableCount).Error
	require.NoError(t, err)
	// Should have fewer tables than before (complete rollback may not remove all tables)
	assert.LessOrEqual(t, afterTableCount, tableCount/2)
}

func TestStepMigrations(t *testing.T) {
	db := setupTestDB(t)

	migrator, err := sqlite.NewSQLiteMigrator(db)
	require.NoError(t, err)

	// Test migrating up in steps
	steps := 5
	err = migrator.Up(&steps)
	require.NoError(t, err)

	// Verify some tables exist but not all
	var tableCount int64
	err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name != 'schema_migrations'").Scan(&tableCount).Error
	require.NoError(t, err)
	assert.Greater(t, tableCount, int64(0))
	assert.Less(t, tableCount, int64(15)) // Should be less than total tables

	// Test migrating down in steps
	downSteps := -2
	err = migrator.Down(&downSteps)
	require.NoError(t, err)

	// Verify some tables were removed
	var newTableCount int64
	err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name != 'schema_migrations'").Scan(&newTableCount).Error
	require.NoError(t, err)
	assert.Less(t, newTableCount, tableCount)
}

func TestMigrationValidation(t *testing.T) {
	db := setupTestDB(t)

	migrator, err := sqlite.NewSQLiteMigrator(db)
	require.NoError(t, err)

	// Test invalid step parameters
	t.Run("InvalidUpSteps", func(t *testing.T) {
		negativeSteps := -5
		err := migrator.Up(&negativeSteps)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "steps cannot be negative")
	})

	t.Run("InvalidDownSteps", func(t *testing.T) {
		positiveSteps := 5
		err := migrator.Down(&positiveSteps)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "steps cannot be positive")
	})
}

func TestSQLiteSpecificFeatures(t *testing.T) {
	db := setupTestDB(t)

	migrator, err := sqlite.NewSQLiteMigrator(db)
	require.NoError(t, err)

	// Run migrations
	err = migrator.Migrate()
	require.NoError(t, err)

	t.Run("TestSQLiteIntegerPrimaryKey", func(t *testing.T) {
		// Verify that our AUTOINCREMENT primary keys work correctly
		// Test by inserting and checking auto-increment behavior
		err = db.Exec(`INSERT INTO Artifact (type_id, uri) VALUES (1, 'test://uri')`).Error
		require.NoError(t, err)

		var maxID int64
		err = db.Raw("SELECT id FROM Artifact ORDER BY id DESC LIMIT 1").Scan(&maxID).Error
		require.NoError(t, err)
		assert.Greater(t, maxID, int64(0))
	})

	t.Run("TestSQLiteBooleanHandling", func(t *testing.T) {
		// Test boolean values (stored as INTEGER in SQLite)
		// Insert a test record with boolean values
		err = db.Exec(`INSERT INTO ArtifactProperty 
			(artifact_id, name, is_custom_property, bool_value) 
			VALUES (1, 'test_bool', 1, 0)`).Error
		require.NoError(t, err)

		// Query back the boolean value
		var boolValue int
		err = db.Raw("SELECT bool_value FROM ArtifactProperty WHERE name = 'test_bool'").Scan(&boolValue).Error
		require.NoError(t, err)
		assert.Equal(t, 0, boolValue) // SQLite stores boolean as 0/1
	})

	t.Run("TestSQLiteTextHandling", func(t *testing.T) {
		// Test TEXT field handling (SQLite's dynamic typing)
		err = db.Exec(`INSERT INTO ArtifactProperty 
			(artifact_id, name, is_custom_property, string_value) 
			VALUES (2, 'test_string', 0, 'This is a long text string that would be MEDIUMTEXT in MySQL')`).Error
		require.NoError(t, err)

		var stringValue string
		err = db.Raw("SELECT string_value FROM ArtifactProperty WHERE name = 'test_string'").Scan(&stringValue).Error
		require.NoError(t, err)
		assert.Contains(t, stringValue, "long text string")
	})

	t.Run("TestSQLiteBlobHandling", func(t *testing.T) {
		// Test BLOB field handling
		testData := []byte{0x01, 0x02, 0x03, 0x04, 0xFF}
		err = db.Exec(`INSERT INTO ArtifactProperty 
			(artifact_id, name, is_custom_property, byte_value) 
			VALUES (3, 'test_blob', 0, ?)`, testData).Error
		require.NoError(t, err)

		var blobValue []byte
		err = db.Raw("SELECT byte_value FROM ArtifactProperty WHERE name = 'test_blob'").Row().Scan(&blobValue)
		require.NoError(t, err)
		assert.Equal(t, testData, blobValue)
	})
}

func TestConcurrentAccess(t *testing.T) {
	// Create temporary file for SQLite database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "concurrent_test.db")

	// Test that SQLite can handle concurrent connections
	db1, err := gorm.Open(sqlitedriver.Open(dbPath), &gorm.Config{
		TranslateError: true,
	})
	require.NoError(t, err)

	db2, err := gorm.Open(sqlitedriver.Open(dbPath), &gorm.Config{
		TranslateError: true,
	})
	require.NoError(t, err)

	// Run migrations on first connection
	migrator1, err := sqlite.NewSQLiteMigrator(db1)
	require.NoError(t, err)
	err = migrator1.Migrate()
	require.NoError(t, err)

	// Verify second connection can read the schema
	var tableCount int64
	err = db2.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table'").Scan(&tableCount).Error
	require.NoError(t, err)
	assert.Greater(t, tableCount, int64(0))
}