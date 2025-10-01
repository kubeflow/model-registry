package testutils

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"testing"

	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/kubeflow/model-registry/internal/tls"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	cont_mysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/gorm"
)

var (
	sharedMySQLContainer testcontainers.Container
	sharedMySQLDSN       string
	setupOnce            sync.Once
	cleanupOnce          sync.Once
	setupMutex           sync.RWMutex
)

// SetupSharedMySQL initializes a shared MySQL container for all tests
func SetupSharedMySQL() error {
	var setupErr error
	setupOnce.Do(func() {
		ctx := context.Background()

		// Use MySQL 8.0 with the specialized module for better reliability
		mysqlContainer, err := cont_mysql.Run(ctx, "mysql:8.3",
			cont_mysql.WithDatabase("test"),
			cont_mysql.WithUsername("root"),
			cont_mysql.WithPassword("root"),
			cont_mysql.WithConfigFile(filepath.Join("testdata", "testdb.cnf")),
		)
		if err != nil {
			setupErr = err
			return
		}

		sharedMySQLContainer = mysqlContainer

		// Get connection string using the module's method
		dsn, err := mysqlContainer.ConnectionString(ctx)
		if err != nil {
			setupErr = err
			return
		}

		sharedMySQLDSN = dsn
	})
	return setupErr
}

// GetSharedMySQLDB returns a connection to the shared MySQL database
func GetSharedMySQLDB(t *testing.T) (*gorm.DB, func()) {
	setupMutex.RLock()
	defer setupMutex.RUnlock()

	require.NotNil(t, sharedMySQLContainer, "Shared MySQL container not initialized. Call SetupSharedMySQL first.")

	dbConnector := mysql.NewMySQLDBConnector(sharedMySQLDSN, &tls.TLSConfig{})
	db, err := dbConnector.Connect()
	require.NoError(t, err)

	// Return cleanup function that only closes the DB connection, not the container
	cleanup := func() {
		sqlDB, err := db.DB()
		if err != nil {
			t.Logf("Failed to get underlying sql.DB: %v", err)
			return
		}
		sqlDB.Close() //nolint:errcheck
	}

	return db, cleanup
}

// GetSharedMySQLDSN returns the DSN for the shared MySQL database
func GetSharedMySQLDSN(t *testing.T) string {
	setupMutex.RLock()
	defer setupMutex.RUnlock()

	require.NotNil(t, sharedMySQLContainer, "Shared MySQL container not initialized. Call SetupSharedMySQL first.")
	return sharedMySQLDSN
}

// CleanupSharedMySQL terminates the shared MySQL container
func CleanupSharedMySQL() {
	cleanupOnce.Do(func() {
		if sharedMySQLContainer != nil {
			ctx := context.Background()
			err := sharedMySQLContainer.Terminate(ctx)
			if err != nil {
				fmt.Printf("Failed to terminate shared MySQL container: %v\n", err)
			}
		}
	})
}

// SetupMySQLWithMigrations returns a migrated MySQL database connection
func SetupMySQLWithMigrations(t *testing.T, spec *datastore.Spec) (*gorm.DB, func()) {
	db, cleanup := GetSharedMySQLDB(t)

	ds, err := datastore.NewConnector("embedmd", &embedmd.EmbedMDConfig{DB: db})
	if err != nil {
		t.Fatalf("unable get datastore connector: %v", err)
	}

	_, err = ds.Connect(spec)
	if err != nil {
		t.Fatalf("unable to connect to datastore: %v", err)
	}

	return db, cleanup
}

// CleanupTestData cleans up test data from the shared database
func CleanupTestData(t *testing.T, db *gorm.DB) {
	// List of tables to clean up (in reverse dependency order)
	// Note: Type and TypeProperty tables are excluded because they contain
	// essential system data that should not be cleaned up between tests
	tables := []string{
		"ArtifactProperty",
		"ContextProperty",
		"ExecutionProperty",
		"ParentContext",
		"Attribution",
		"Association",
		"Event",
		"Artifact",
		"Execution",
		"Context",
		// "Type", // DO NOT clean up - contains essential system types
		// "TypeProperty", // DO NOT clean up - contains essential system properties
	}

	for _, table := range tables {
		// Use raw SQL to truncate tables quickly
		err := db.Exec("TRUNCATE TABLE " + table).Error
		if err != nil {
			// If truncate fails, try delete (some tables might have foreign key constraints)
			err = db.Exec("DELETE FROM " + table).Error
			require.NoError(t, err, "Failed to clean up table: "+table)
		}
	}
}

// TestMainHelper provides a helper function for package-level test setup
func TestMainHelper(m *testing.M) int {
	// Setup shared MySQL container
	err := SetupSharedMySQL()
	if err != nil {
		fmt.Printf("Failed to setup shared MySQL container: %v\n", err)
		return 1
	}

	// Run tests
	code := m.Run()

	// Cleanup shared MySQL container
	CleanupSharedMySQL()

	return code
}
