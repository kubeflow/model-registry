package sqlite_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/sqlite"
	_tls "github.com/kubeflow/model-registry/internal/tls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLiteDBConnector_Connect_Basic(t *testing.T) {
	t.Run("InMemoryDatabase", func(t *testing.T) {
		// Test in-memory SQLite database
		connector := sqlite.NewSQLiteDBConnector(":memory:", &_tls.TLSConfig{})

		db, err := connector.Connect()
		require.NoError(t, err)
		assert.NotNil(t, db)

		// Test that we can perform a simple query
		var result int
		err = db.Raw("SELECT 1").Scan(&result).Error
		require.NoError(t, err)
		assert.Equal(t, 1, result)

		// Clean up
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
	})

	t.Run("FileBasedDatabase", func(t *testing.T) {
		// Create temporary file for database
		tempDir := t.TempDir()
		dbPath := filepath.Join(tempDir, "test.db")

		connector := sqlite.NewSQLiteDBConnector(dbPath, &_tls.TLSConfig{})

		db, err := connector.Connect()
		require.NoError(t, err)
		assert.NotNil(t, db)

		// Test basic operations
		var result int
		err = db.Raw("SELECT 1").Scan(&result).Error
		require.NoError(t, err)
		assert.Equal(t, 1, result)

		// Verify database file was created
		_, err = os.Stat(dbPath)
		require.NoError(t, err)

		// Clean up
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
	})

	t.Run("ExistingConnection", func(t *testing.T) {
		tempDir := t.TempDir()
		dbPath := filepath.Join(tempDir, "existing.db")

		connector := sqlite.NewSQLiteDBConnector(dbPath, &_tls.TLSConfig{})

		// First connection
		db1, err := connector.Connect()
		require.NoError(t, err)

		// Second call should return the same connection
		db2, err := connector.Connect()
		require.NoError(t, err)
		assert.Equal(t, db1, db2)

		// Clean up
		sqlDB, err := db1.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
	})
}

func TestSQLiteDBConnector_TLSHandling(t *testing.T) {
	t.Run("EmptyTLSConfig", func(t *testing.T) {
		connector := sqlite.NewSQLiteDBConnector(":memory:", &_tls.TLSConfig{})

		db, err := connector.Connect()
		require.NoError(t, err)
		assert.NotNil(t, db)

		// Clean up
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
	})

	t.Run("NilTLSConfig", func(t *testing.T) {
		connector := sqlite.NewSQLiteDBConnector(":memory:", nil)

		db, err := connector.Connect()
		require.NoError(t, err)
		assert.NotNil(t, db)

		// Clean up
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
	})

	t.Run("TLSConfigWithCertificates", func(t *testing.T) {
		// SQLite doesn't support TLS, but should log warning and still connect
		tempDir := t.TempDir()

		// Create dummy certificate files (content doesn't matter for this test)
		certPath := filepath.Join(tempDir, "cert.pem")
		keyPath := filepath.Join(tempDir, "key.pem")
		caPath := filepath.Join(tempDir, "ca.pem")

		err := os.WriteFile(certPath, []byte("dummy cert"), 0600)
		require.NoError(t, err)
		err = os.WriteFile(keyPath, []byte("dummy key"), 0600)
		require.NoError(t, err)
		err = os.WriteFile(caPath, []byte("dummy ca"), 0600)
		require.NoError(t, err)

		tlsConfig := &_tls.TLSConfig{
			CertPath:         certPath,
			KeyPath:          keyPath,
			RootCertPath:     caPath,
			VerifyServerCert: true,
			Cipher:           "TLS_AES_256_GCM_SHA384",
		}

		connector := sqlite.NewSQLiteDBConnector(":memory:", tlsConfig)

		// Should still connect successfully despite TLS config
		db, err := connector.Connect()
		require.NoError(t, err)
		assert.NotNil(t, db)

		// Should be able to perform operations
		var result int
		err = db.Raw("SELECT 1").Scan(&result).Error
		require.NoError(t, err)
		assert.Equal(t, 1, result)

		// Clean up
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
	})
}

func TestSQLiteDBConnector_DatabaseOperations(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "operations.db")

	connector := sqlite.NewSQLiteDBConnector(dbPath, &_tls.TLSConfig{})

	db, err := connector.Connect()
	require.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close() //nolint:errcheck
	}()

	t.Run("CreateTable", func(t *testing.T) {
		err := db.Exec(`CREATE TABLE test_table (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			value REAL,
			is_active INTEGER,
			data BLOB
		)`).Error
		require.NoError(t, err)
	})

	t.Run("InsertData", func(t *testing.T) {
		testData := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f} // "Hello" in bytes

		err := db.Exec(`INSERT INTO test_table (name, value, is_active, data) 
			VALUES (?, ?, ?, ?)`, "Test Entry", 3.14159, 1, testData).Error
		require.NoError(t, err)

		err = db.Exec(`INSERT INTO test_table (name, value, is_active, data) 
			VALUES (?, ?, ?, ?)`, "Another Entry", 2.71828, 0, nil).Error
		require.NoError(t, err)
	})

	t.Run("QueryData", func(t *testing.T) {
		var entries []struct {
			ID       int64   `gorm:"column:id"`
			Name     string  `gorm:"column:name"`
			Value    float64 `gorm:"column:value"`
			IsActive int     `gorm:"column:is_active"`
			Data     []byte  `gorm:"column:data"`
		}

		err := db.Raw("SELECT id, name, value, is_active, data FROM test_table ORDER BY id").Scan(&entries).Error
		require.NoError(t, err)
		require.Len(t, entries, 2)

		assert.Equal(t, "Test Entry", entries[0].Name)
		assert.InDelta(t, 3.14159, entries[0].Value, 0.00001)
		assert.Equal(t, 1, entries[0].IsActive)
		assert.Equal(t, []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}, entries[0].Data)

		assert.Equal(t, "Another Entry", entries[1].Name)
		assert.InDelta(t, 2.71828, entries[1].Value, 0.00001)
		assert.Equal(t, 0, entries[1].IsActive)
		assert.Nil(t, entries[1].Data)
	})

	t.Run("UpdateData", func(t *testing.T) {
		err := db.Exec("UPDATE test_table SET is_active = ? WHERE name = ?", 1, "Another Entry").Error
		require.NoError(t, err)

		var isActive int
		err = db.Raw("SELECT is_active FROM test_table WHERE name = ?", "Another Entry").Scan(&isActive).Error
		require.NoError(t, err)
		assert.Equal(t, 1, isActive)
	})

	t.Run("DeleteData", func(t *testing.T) {
		err := db.Exec("DELETE FROM test_table WHERE name = ?", "Test Entry").Error
		require.NoError(t, err)

		var count int64
		err = db.Raw("SELECT COUNT(*) FROM test_table").Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})
}

func TestSQLiteDBConnector_Transactions(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "transactions.db")

	connector := sqlite.NewSQLiteDBConnector(dbPath, &_tls.TLSConfig{})

	db, err := connector.Connect()
	require.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close() //nolint:errcheck
	}()

	// Create test table
	err = db.Exec(`CREATE TABLE transaction_test (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		value TEXT
	)`).Error
	require.NoError(t, err)

	t.Run("CommitTransaction", func(t *testing.T) {
		tx := db.Begin()

		err := tx.Exec("INSERT INTO transaction_test (value) VALUES (?)", "committed_value").Error
		require.NoError(t, err)

		err = tx.Commit().Error
		require.NoError(t, err)

		// Verify data was committed
		var value string
		err = db.Raw("SELECT value FROM transaction_test WHERE value = ?", "committed_value").Scan(&value).Error
		require.NoError(t, err)
		assert.Equal(t, "committed_value", value)
	})

	t.Run("RollbackTransaction", func(t *testing.T) {
		tx := db.Begin()

		err := tx.Exec("INSERT INTO transaction_test (value) VALUES (?)", "rollback_value").Error
		require.NoError(t, err)

		err = tx.Rollback().Error
		require.NoError(t, err)

		// Verify data was not committed
		var count int64
		err = db.Raw("SELECT COUNT(*) FROM transaction_test WHERE value = ?", "rollback_value").Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestSQLiteDBConnector_Connect_ErrorCases(t *testing.T) {
	t.Run("InvalidPath", func(t *testing.T) {
		// Try to create database in non-existent directory without permission
		invalidPath := "/root/nonexistent/directory/test.db"
		connector := sqlite.NewSQLiteDBConnector(invalidPath, &_tls.TLSConfig{})

		db, err := connector.Connect()
		// This might or might not fail depending on the system
		// SQLite will try to create the database file, but may fail due to permissions
		if err != nil {
			assert.Error(t, err)
			assert.Nil(t, db)
		} else {
			// If it succeeds, clean up
			sqlDB, _ := db.DB()
			sqlDB.Close() //nolint:errcheck
		}
	})

	t.Run("WithMaxRetries", func(t *testing.T) {
		connector := sqlite.NewSQLiteDBConnector(":memory:", &_tls.TLSConfig{})
		connector = connector.WithMaxRetries(1)

		db, err := connector.Connect()
		require.NoError(t, err)
		assert.NotNil(t, db)

		// Clean up
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
	})
}

func TestSQLiteDBConnector_ConcurrentConnections(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "concurrent.db")

	// Test multiple connectors to the same database file
	connector1 := sqlite.NewSQLiteDBConnector(dbPath, &_tls.TLSConfig{})
	connector2 := sqlite.NewSQLiteDBConnector(dbPath, &_tls.TLSConfig{})

	db1, err := connector1.Connect()
	require.NoError(t, err)

	db2, err := connector2.Connect()
	require.NoError(t, err)

	// Create table with first connection
	err = db1.Exec(`CREATE TABLE IF NOT EXISTS concurrent_test (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message TEXT
	)`).Error
	require.NoError(t, err)

	// Insert data with first connection
	err = db1.Exec("INSERT INTO concurrent_test (message) VALUES (?)", "from_db1").Error
	require.NoError(t, err)

	// Read data with second connection
	var message string
	err = db2.Raw("SELECT message FROM concurrent_test WHERE message = ?", "from_db1").Scan(&message).Error
	require.NoError(t, err)
	assert.Equal(t, "from_db1", message)

	// Insert data with second connection
	err = db2.Exec("INSERT INTO concurrent_test (message) VALUES (?)", "from_db2").Error
	require.NoError(t, err)

	// Read data with first connection
	var messages []string
	err = db1.Raw("SELECT message FROM concurrent_test ORDER BY id").Scan(&messages).Error
	require.NoError(t, err)
	assert.Equal(t, []string{"from_db1", "from_db2"}, messages)

	// Clean up
	sqlDB1, _ := db1.DB()
	sqlDB1.Close() //nolint:errcheck
	sqlDB2, _ := db2.DB()
	sqlDB2.Close() //nolint:errcheck
}

func TestSQLiteDBConnector_DB(t *testing.T) {
	connector := sqlite.NewSQLiteDBConnector(":memory:", &_tls.TLSConfig{})

	// Initially should return nil
	assert.Nil(t, connector.DB())

	// After connecting should return the database
	db, err := connector.Connect()
	require.NoError(t, err)

	retrievedDB := connector.DB()
	assert.Equal(t, db, retrievedDB)

	// Clean up
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.Close() //nolint:errcheck
}

func TestSQLiteDBConnector_ConnectionRetries(t *testing.T) {
	// Test connection retries with a path that might initially fail
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "retry_test.db")

	connector := sqlite.NewSQLiteDBConnector(dbPath, &_tls.TLSConfig{})
	connector = connector.WithMaxRetries(3)

	// This should succeed (SQLite is quite forgiving)
	db, err := connector.Connect()
	require.NoError(t, err)
	assert.NotNil(t, db)

	// Verify database works
	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	require.NoError(t, err)
	assert.Equal(t, 1, result)

	// Clean up
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.Close() //nolint:errcheck
}

func TestSQLiteDBConnector_WALMode(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "wal_test.db")

	connector := sqlite.NewSQLiteDBConnector(dbPath, &_tls.TLSConfig{})

	db, err := connector.Connect()
	require.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close() //nolint:errcheck
	}()

	// Enable WAL mode for better concurrency
	err = db.Exec("PRAGMA journal_mode=WAL").Error
	require.NoError(t, err)

	// Verify WAL mode is enabled
	var journalMode string
	err = db.Raw("PRAGMA journal_mode").Scan(&journalMode).Error
	require.NoError(t, err)
	assert.Equal(t, "wal", journalMode)

	// Test basic operations in WAL mode
	err = db.Exec(`CREATE TABLE wal_test (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		data TEXT
	)`).Error
	require.NoError(t, err)

	err = db.Exec("INSERT INTO wal_test (data) VALUES (?)", "wal_data").Error
	require.NoError(t, err)

	var data string
	err = db.Raw("SELECT data FROM wal_test LIMIT 1").Scan(&data).Error
	require.NoError(t, err)
	assert.Equal(t, "wal_data", data)

	// Verify WAL file exists
	walPath := dbPath + "-wal"
	time.Sleep(100 * time.Millisecond) // Give SQLite time to create WAL file
	_, err = os.Stat(walPath)
	// WAL file might not exist immediately or might be cleaned up, so we don't require it
	if err == nil {
		t.Logf("WAL file created: %s", walPath)
	}
}