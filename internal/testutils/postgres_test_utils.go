package testutils

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd/postgres"
	"github.com/kubeflow/model-registry/internal/tls"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	cont_postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

var (
	sharedPostgresContainer testcontainers.Container
	sharedPostgresDSN       string
	setupPostgresOnce       sync.Once
	cleanupPostgresOnce     sync.Once
	setupPostgresMutex      sync.RWMutex
)

// SetupSharedPostgres initializes a shared PostgreSQL container for all tests
func SetupSharedPostgres() error {
	var setupErr error
	setupPostgresOnce.Do(func() {
		ctx := context.Background()

		// Use PostgreSQL 16 with the specialized module for better reliability
		postgresContainer, err := cont_postgres.Run(ctx, "postgres:16",
			cont_postgres.WithDatabase("test"),
			cont_postgres.WithUsername("postgres"),
			cont_postgres.WithPassword("postgres"),
			testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
		)
		if err != nil {
			setupErr = err
			return
		}

		sharedPostgresContainer = postgresContainer

		// Get connection string using the module's method
		dsn, err := postgresContainer.ConnectionString(ctx)
		if err != nil {
			setupErr = err
			return
		}

		sharedPostgresDSN = dsn
	})
	return setupErr
}

// GetSharedPostgresDB returns a connection to the shared PostgreSQL database
func GetSharedPostgresDB(t *testing.T) (*gorm.DB, func()) {
	setupPostgresMutex.RLock()
	defer setupPostgresMutex.RUnlock()

	require.NotNil(t, sharedPostgresContainer, "Shared PostgreSQL container not initialized. Call SetupSharedPostgres first.")

	dbConnector := postgres.NewPostgresDBConnector(sharedPostgresDSN, &tls.TLSConfig{})
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

// GetSharedPostgresDSN returns the DSN for the shared PostgreSQL database
func GetSharedPostgresDSN(t *testing.T) string {
	setupPostgresMutex.RLock()
	defer setupPostgresMutex.RUnlock()

	require.NotNil(t, sharedPostgresContainer, "Shared PostgreSQL container not initialized. Call SetupSharedPostgres first.")
	return sharedPostgresDSN
}

// CleanupSharedPostgres terminates the shared PostgreSQL container
func CleanupSharedPostgres() {
	cleanupPostgresOnce.Do(func() {
		if sharedPostgresContainer != nil {
			ctx := context.Background()
			err := sharedPostgresContainer.Terminate(ctx)
			if err != nil {
				fmt.Printf("Failed to terminate shared PostgreSQL container: %v\n", err)
			}
		}
	})
}

// SetupPostgresWithMigrations returns a migrated PostgreSQL database connection
func SetupPostgresWithMigrations(t *testing.T, spec *datastore.Spec) (*gorm.DB, func()) {
	db, cleanup := GetSharedPostgresDB(t)

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

// CleanupPostgresTestData cleans up test data from the shared PostgreSQL database
func CleanupPostgresTestData(t *testing.T, db *gorm.DB) {
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

	// Disable triggers and foreign key constraints temporarily (PostgreSQL-specific)
	err := db.Exec("SET session_replication_role = replica").Error
	require.NoError(t, err)

	for _, table := range tables {
		// Use PostgreSQL-specific TRUNCATE with CASCADE for better cleanup
		err := db.Exec("TRUNCATE TABLE IF EXISTS \"" + table + "\" CASCADE").Error
		if err != nil {
			// If truncate fails, try delete (some tables might have foreign key constraints)
			err = db.Exec("DELETE FROM \"" + table + "\"").Error
			if err != nil {
				t.Logf("Failed to clean up table: %s, error: %v", table, err)
			}
		}
	}

	// Re-enable triggers and foreign key constraints (PostgreSQL-specific)
	err = db.Exec("SET session_replication_role = DEFAULT").Error
	require.NoError(t, err)
}

// TestMainPostgresHelper provides a helper function for package-level test setup with PostgreSQL
func TestMainPostgresHelper(m *testing.M) int {
	// Setup shared PostgreSQL container
	err := SetupSharedPostgres()
	if err != nil {
		fmt.Printf("Failed to setup shared PostgreSQL container: %v\n", err)
		return 1
	}

	// Run tests
	code := m.Run()

	// Cleanup shared PostgreSQL container
	CleanupSharedPostgres()

	return code
}
