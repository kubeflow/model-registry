package proxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	cont_mysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, string, func()) {
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

	dsn := mysqlContainer.MustConnectionString(ctx)
	dbConnector := mysql.NewMySQLDBConnector(dsn, nil)

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

	return db, dsn, cleanup
}

func TestReadinessHandler_NonEmbedMD(t *testing.T) {
	ds := datastore.Datastore{
		Type: "mlmd",
	}
	handler := ReadinessHandler(ds)

	req, err := http.NewRequest("GET", "/readyz/isDirty", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "ok", rr.Body.String())
}

func TestReadinessHandler_EmbedMD_Success(t *testing.T) {
	db, dsn, cleanup := setupTestDB(t)
	defer cleanup()

	// run migrations to create tables
	migrator, err := mysql.NewMySQLMigrator(db)
	require.NoError(t, err)
	err = migrator.Migrate()
	require.NoError(t, err)

	ds := datastore.Datastore{
		Type: "embedmd",
		EmbedMD: embedmd.EmbedMDConfig{
			DatabaseType: "mysql",
			DatabaseDSN:  dsn,
		},
	}

	handler := ReadinessHandler(ds)
	req, err := http.NewRequest("GET", "/readyz/isDirty", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "ok", rr.Body.String())
}

func TestReadinessHandler_EmbedMD_Dirty(t *testing.T) {
	db, dsn, cleanup := setupTestDB(t)
	defer cleanup()

	// run migrations to create tables
	migrator, err := mysql.NewMySQLMigrator(db)
	require.NoError(t, err)
	err = migrator.Migrate()
	require.NoError(t, err)

	// manually set latest migration to dirty
	err = db.Exec("UPDATE schema_migrations SET dirty = 1").Error
	require.NoError(t, err)

	ds := datastore.Datastore{
		Type: "embedmd",
		EmbedMD: embedmd.EmbedMDConfig{
			DatabaseType: "mysql",
			DatabaseDSN:  dsn,
		},
	}

	handler := ReadinessHandler(ds)
	req, err := http.NewRequest("GET", "/readyz/isDirty", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
	assert.Contains(t, rr.Body.String(), "database schema is in dirty state")
}
