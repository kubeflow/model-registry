package testutils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(TestMainHelper(m))
}

func TestSharedMySQLUtility(t *testing.T) {
	t.Run("GetSharedMySQLDB", func(t *testing.T) {
		db, cleanup := GetSharedMySQLDB(t)
		defer cleanup()

		// Test basic connectivity
		var result int
		err := db.Raw("SELECT 1").Scan(&result).Error
		require.NoError(t, err)
		assert.Equal(t, 1, result)
	})

	t.Run("GetSharedMySQLDSN", func(t *testing.T) {
		dsn := GetSharedMySQLDSN(t)
		assert.NotEmpty(t, dsn)
		assert.Contains(t, dsn, "root:root@tcp(")
		assert.Contains(t, dsn, ")/test")
	})

	t.Run("SetupMySQLWithMigrations", func(t *testing.T) {
		db, cleanup := SetupMySQLWithMigrations(t)
		defer cleanup()

		// Test that migrations were applied
		var schemaVersion int
		err := db.Raw("SELECT schema_version FROM MLMDEnv LIMIT 1").Scan(&schemaVersion).Error
		require.NoError(t, err)
		assert.Equal(t, 10, schemaVersion)
	})

	t.Run("CleanupTestData", func(t *testing.T) {
		db, cleanup := SetupMySQLWithMigrations(t)
		defer cleanup()

		// Insert some test data
		err := db.Exec("INSERT INTO Artifact (id, type_id, uri, name, external_id, create_time_since_epoch, last_update_time_since_epoch) VALUES (1, 1, 'test://uri', 'test', 'ext-1', 1000, 1000)").Error
		require.NoError(t, err)

		// Verify data exists
		var count int64
		err = db.Raw("SELECT COUNT(*) FROM Artifact WHERE id = 1").Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// Clean up test data
		CleanupTestData(t, db)

		// Verify data is gone
		err = db.Raw("SELECT COUNT(*) FROM Artifact WHERE id = 1").Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}
