package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/postgres"
	_tls "github.com/kubeflow/model-registry/internal/tls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	cont_postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
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

// Package-level shared database instance
var (
	sharedDB          *gorm.DB
	postgresContainer *cont_postgres.PostgresContainer
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Create Postgres container once for all tests
	container, err := cont_postgres.Run(
		ctx,
		"postgres:16",
		cont_postgres.WithUsername("postgres"),
		cont_postgres.WithPassword("testpass"),
		cont_postgres.WithDatabase("testdb"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
	)
	if err != nil {
		panic("Failed to start Postgres container: " + err.Error())
	}
	postgresContainer = container

	defer func() {
		if sharedDB != nil {
			if sqlDB, err := sharedDB.DB(); err == nil {
				sqlDB.Close() //nolint:errcheck
			}
		}

		if postgresContainer != nil {
			testcontainers.TerminateContainer(postgresContainer) //nolint:errcheck
		}
	}()

	// Connect to the database
	dbConnector := postgres.NewPostgresDBConnector(postgresContainer.MustConnectionString(ctx), &_tls.TLSConfig{})
	sharedDB, err = dbConnector.Connect()
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Run all tests
	code := m.Run()

	os.Exit(code)
}

// cleanupTestData truncates tables and drops indexes to provide clean state between tests
func cleanupTestData(t *testing.T, db *gorm.DB) {
	// First, drop any indexes that might conflict with migrations
	dropIndexes := []string{
		"idx_artifact_uri",
		"idx_artifact_create_time_since_epoch",
		"idx_artifact_last_update_time_since_epoch",
		"idx_artifact_external_id",
		"idx_context_name",
		"idx_context_create_time_since_epoch",
		"idx_context_last_update_time_since_epoch",
		"idx_context_external_id",
		"idx_execution_name",
		"idx_execution_create_time_since_epoch",
		"idx_execution_last_update_time_since_epoch",
		"idx_execution_external_id",
		"idx_event_artifact_id",
		"idx_event_execution_id",
		"idx_event_milliseconds_since_epoch",
		"idx_parentcontext_parent_id",
		"idx_parentcontext_child_id",
		"idx_association_context_id",
		"idx_association_execution_id",
		"idx_attribution_context_id",
		"idx_attribution_artifact_id",
		"idx_artifactproperty_artifact_id",
		"idx_artifactproperty_name",
		"idx_contextproperty_context_id",
		"idx_contextproperty_name",
		"idx_executionproperty_execution_id",
		"idx_executionproperty_name",
		"idx_typeproperty_type_id",
		"idx_typeproperty_name",
	}

	for _, index := range dropIndexes {
		err := db.Exec("DROP INDEX IF EXISTS " + index).Error
		if err != nil {
			// Log but don't fail - index might not exist
			t.Logf("Could not drop index %s: %v", index, err)
		}
	}

	// List of tables to clean up (in dependency order - dependent tables first)
	tables := []string{
		"ArtifactProperty",
		"ContextProperty",
		"ExecutionProperty",
		"TypeProperty",
		"ParentContext",
		"Association",
		"Attribution",
		"Event",
		"Artifact",
		"Execution",
		"Context",
		"Type",
		"MLMDEnv",
		"schema_migrations",
	}

	// Disable triggers and foreign key constraints temporarily (PostgreSQL-specific)
	err := db.Exec("SET session_replication_role = replica").Error
	require.NoError(t, err)

	// Truncate all tables
	for _, table := range tables {
		err := db.Exec("TRUNCATE TABLE IF EXISTS \"" + table + "\" CASCADE").Error
		if err != nil {
			// Log but don't fail - table might not exist
			t.Logf("Could not truncate table %s: %v", table, err)
		}
	}

	// Re-enable triggers and foreign key constraints (PostgreSQL-specific)
	err = db.Exec("SET session_replication_role = DEFAULT").Error
	require.NoError(t, err)
}

func TestMigrations(t *testing.T) {
	cleanupTestData(t, sharedDB)

	// Create migrator
	migrator, err := postgres.NewPostgresMigrator(sharedDB)
	require.NoError(t, err)

	// Run migrations
	err = migrator.Migrate()
	require.NoError(t, err)

	// Verify MLMDEnv table
	var schemaVersion int
	err = sharedDB.Raw("SELECT schema_version FROM \"MLMDEnv\" LIMIT 1").Scan(&schemaVersion).Error
	require.NoError(t, err)
	assert.Equal(t, 10, schemaVersion)

	// Verify Type table has expected entries
	var count int64
	err = sharedDB.Model(&Type{}).Count(&count).Error
	require.NoError(t, err)
	assert.Greater(t, count, int64(0))
}

func TestDownMigrations(t *testing.T) {
	cleanupTestData(t, sharedDB)

	migrator, err := postgres.NewPostgresMigrator(sharedDB)
	require.NoError(t, err)

	// Run migrations
	err = migrator.Migrate()
	require.NoError(t, err)

	// Down migrations
	err = migrator.Down(nil)
	require.NoError(t, err)

	// Verify tables don't exist (except schema_migrations)
	var count int64
	err = sharedDB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name != 'schema_migrations'").Scan(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// TestCatalogModelNamespacedNamesMigration verifies that migration 000026 updates existing
// catalog model names from "model_name" to "source_id:model_name".
func TestCatalogModelNamespacedNamesMigration(t *testing.T) {
	cleanupTestData(t, sharedDB)

	migrator, err := postgres.NewPostgresMigrator(sharedDB)
	require.NoError(t, err)

	err = migrator.Migrate()
	require.NoError(t, err)

	// Roll back the last migration (000026) so we can insert old-format data
	downOne := -1
	err = migrator.Down(&downOne)
	require.NoError(t, err)

	// Insert catalog model type (normally created at runtime by syncTypes)
	var typeID int32
	err = sharedDB.Raw(
		`INSERT INTO "Type" (name, type_kind) VALUES ('kf.CatalogModel', 2) RETURNING id`,
	).Scan(&typeID).Error
	require.NoError(t, err)

	// Insert a context (catalog model) with legacy name format (no source_id prefix)
	var contextID int32
	err = sharedDB.Raw(
		`INSERT INTO "Context" (type_id, name, create_time_since_epoch, last_update_time_since_epoch) VALUES (?, 'legacy-model', 0, 0) RETURNING id`,
		typeID,
	).Scan(&contextID).Error
	require.NoError(t, err)

	// Set source_id property so the migration can build source_id:model_name
	err = sharedDB.Exec(
		`INSERT INTO "ContextProperty" (context_id, name, is_custom_property, string_value) VALUES (?, 'source_id', false, 'test-source')`,
		contextID,
	).Error
	require.NoError(t, err)

	// Run migration 000026 up
	upOne := 1
	err = migrator.Up(&upOne)
	require.NoError(t, err)

	// Verify the context name was updated to namespaced format
	var name string
	err = sharedDB.Raw(`SELECT name FROM "Context" WHERE id = ?`, contextID).Scan(&name).Error
	require.NoError(t, err)
	assert.Equal(t, "test-source:legacy-model", name, "catalog model name should be namespaced as source_id:model_name")
}
