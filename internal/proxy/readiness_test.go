package proxy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/core"
	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/kubeflow/model-registry/internal/db"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/internal/tls"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	cont_mysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/gorm"
)

// Package-level shared database instance
var (
	sharedDB                   *gorm.DB
	sharedDSN                  string
	mysqlContainer             *cont_mysql.MySQLContainer
	sharedModelRegistryService api.ModelRegistryApi
)

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
		testcontainers.WithEnv(map[string]string{
			"MYSQL_ROOT_HOST": "%",
		}),
	)
	if err != nil {
		panic("Failed to start MySQL container: " + err.Error())
	}
	mysqlContainer = container

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
	sharedDSN = mysqlContainer.MustConnectionString(ctx)
	dbConnector := mysql.NewMySQLDBConnector(sharedDSN, &tls.TLSConfig{})
	sharedDB, err = dbConnector.Connect()
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Initialize the global db connector for health checks
	err = db.Init("mysql", sharedDSN, &tls.TLSConfig{})
	if err != nil {
		panic("Failed to initialize db: " + err.Error())
	}

	// Run migrations
	migrator, err := mysql.NewMySQLMigrator(sharedDB)
	if err != nil {
		panic("Failed to create migrator: " + err.Error())
	}
	err = migrator.Migrate()
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Setup model registry service
	sharedModelRegistryService = setupModelRegistryService()

	// Run all tests
	code := m.Run()

	os.Exit(code)
}

// getTypeIDs retrieves all type IDs from the database for testing
func getTypeIDs() map[string]int64 {
	typesMap := make(map[string]int64)

	typeNames := []string{
		defaults.RegisteredModelTypeName,
		defaults.ModelVersionTypeName,
		defaults.DocArtifactTypeName,
		defaults.ModelArtifactTypeName,
		defaults.ServingEnvironmentTypeName,
		defaults.InferenceServiceTypeName,
		defaults.ServeModelTypeName,
	}

	for _, typeName := range typeNames {
		var typeRecord schema.Type
		err := sharedDB.Where("name = ?", typeName).First(&typeRecord).Error
		if err != nil {
			panic("Failed to find type: " + typeName + ": " + err.Error())
		}
		typesMap[typeName] = int64(typeRecord.ID)
	}

	return typesMap
}

// setupModelRegistryService creates a complete ModelRegistryService with all repositories for testing
func setupModelRegistryService() api.ModelRegistryApi {
	// Get all type IDs from the database
	typesMap := getTypeIDs()

	// Create all repositories
	artifactRepo := service.NewArtifactRepository(sharedDB, typesMap[defaults.ModelArtifactTypeName], typesMap[defaults.DocArtifactTypeName])
	modelArtifactRepo := service.NewModelArtifactRepository(sharedDB, typesMap[defaults.ModelArtifactTypeName])
	docArtifactRepo := service.NewDocArtifactRepository(sharedDB, typesMap[defaults.DocArtifactTypeName])
	registeredModelRepo := service.NewRegisteredModelRepository(sharedDB, typesMap[defaults.RegisteredModelTypeName])
	modelVersionRepo := service.NewModelVersionRepository(sharedDB, typesMap[defaults.ModelVersionTypeName])
	servingEnvironmentRepo := service.NewServingEnvironmentRepository(sharedDB, typesMap[defaults.ServingEnvironmentTypeName])
	inferenceServiceRepo := service.NewInferenceServiceRepository(sharedDB, typesMap[defaults.InferenceServiceTypeName])
	serveModelRepo := service.NewServeModelRepository(sharedDB, typesMap[defaults.ServeModelTypeName])

	// Create the core service
	service := core.NewModelRegistryService(
		artifactRepo,
		modelArtifactRepo,
		docArtifactRepo,
		registeredModelRepo,
		modelVersionRepo,
		servingEnvironmentRepo,
		inferenceServiceRepo,
		serveModelRepo,
		typesMap,
	)

	return service
}

// cleanupSchemaState resets schema_migrations table to clean state
func cleanupSchemaState(t *testing.T) {
	// Reset schema_migrations to clean state
	err := sharedDB.Exec("UPDATE schema_migrations SET dirty = 0").Error
	require.NoError(t, err)
}

// setDirtySchemaState sets schema_migrations to dirty state for testing
func setDirtySchemaState(t *testing.T) {
	err := sharedDB.Exec("UPDATE schema_migrations SET dirty = 1").Error
	require.NoError(t, err)
}

// createTestDatastore creates a datastore config for testing
func createTestDatastore() datastore.Datastore {
	return datastore.Datastore{
		Type: "embedmd",
		EmbedMD: embedmd.EmbedMDConfig{
			DatabaseType: "mysql",
			DatabaseDSN:  sharedDSN,
		},
	}
}

func TestReadinessHandler_NonEmbedMD(t *testing.T) {
	ds := datastore.Datastore{
		Type: "mlmd",
	}
	dbHealthChecker := NewDatabaseHealthChecker(ds)
	handler := GeneralReadinessHandler(ds, dbHealthChecker)

	req, err := http.NewRequest("GET", "/readyz/isDirty", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, responseOK, rr.Body.String())
}

func TestReadinessHandler_EmbedMD_Success(t *testing.T) {
	// Ensure clean state before test
	cleanupSchemaState(t)

	ds := createTestDatastore()

	dbHealthChecker := NewDatabaseHealthChecker(ds)
	handler := GeneralReadinessHandler(ds, dbHealthChecker)
	req, err := http.NewRequest("GET", "/readyz/isDirty", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, responseOK, rr.Body.String())
}

func TestReadinessHandler_EmbedMD_Dirty(t *testing.T) {
	// Set dirty state for this test
	setDirtySchemaState(t)
	defer cleanupSchemaState(t) // Cleanup after test

	ds := createTestDatastore()

	dbHealthChecker := NewDatabaseHealthChecker(ds)
	handler := GeneralReadinessHandler(ds, dbHealthChecker)
	req, err := http.NewRequest("GET", "/readyz/isDirty", nil)
	require.NoError(t, err)

	// Retry logic for CI robustness
	var rr *httptest.ResponseRecorder
	var responseBody string
	maxRetries := 3

	for i := range maxRetries {
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		responseBody = rr.Body.String()

		// If we get the expected dirty state error, test passes
		if rr.Code == http.StatusServiceUnavailable &&
			(strings.Contains(responseBody, "database schema is in dirty state") ||
				strings.Contains(responseBody, "schema_migrations query error")) {
			break
		}

		// If it's a connection error and not the last retry, wait and try again
		if i < maxRetries-1 && strings.Contains(responseBody, "connection refused") {
			time.Sleep(time.Duration(i+1) * 500 * time.Millisecond) // 500ms, 1s, 1.5s
			continue
		}

		break
	}

	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
	assert.Contains(t, rr.Body.String(), "database schema is in dirty state")
}

func TestGeneralReadinessHandler_WithModelRegistry_Success(t *testing.T) {
	// Ensure clean state before test
	cleanupSchemaState(t)

	ds := createTestDatastore()

	// Create both health checkers
	dbHealthChecker := NewDatabaseHealthChecker(ds)
	mrHealthChecker := NewModelRegistryHealthChecker(sharedModelRegistryService)
	handler := GeneralReadinessHandler(ds, dbHealthChecker, mrHealthChecker)

	req, err := http.NewRequest("GET", "/readyz/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, responseOK, rr.Body.String())
}

func TestGeneralReadinessHandler_WithModelRegistry_JSONFormat(t *testing.T) {
	// Ensure clean state before test
	cleanupSchemaState(t)

	ds := createTestDatastore()

	// Create both health checkers
	dbHealthChecker := NewDatabaseHealthChecker(ds)
	mrHealthChecker := NewModelRegistryHealthChecker(sharedModelRegistryService)
	handler := GeneralReadinessHandler(ds, dbHealthChecker, mrHealthChecker)

	req, err := http.NewRequest("GET", "/readyz/health?format=json", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	// Parse and validate JSON response
	var healthStatus HealthStatus
	err = json.Unmarshal(rr.Body.Bytes(), &healthStatus)
	require.NoError(t, err)

	assert.Equal(t, StatusPass, healthStatus.Status)
	assert.Contains(t, healthStatus.Checks, HealthCheckDatabase)
	assert.Contains(t, healthStatus.Checks, HealthCheckModelRegistry)
	assert.Contains(t, healthStatus.Checks, HealthCheckMeta)

	// Check database health details
	dbCheck := healthStatus.Checks[HealthCheckDatabase]
	assert.Equal(t, StatusPass, dbCheck.Status)
	assert.Equal(t, "database is healthy", dbCheck.Message)

	// Check model registry health details
	mrCheck := healthStatus.Checks[HealthCheckModelRegistry]
	assert.Equal(t, StatusPass, mrCheck.Status)
	assert.Equal(t, "model registry service is healthy", mrCheck.Message)
	assert.Equal(t, float64(5), mrCheck.Details[detailTotalResourcesChecked])
	assert.Equal(t, true, mrCheck.Details[detailRegisteredModelsAccessible])
	assert.Equal(t, true, mrCheck.Details[detailArtifactsAccessible])
}

func TestGeneralReadinessHandler_WithModelRegistry_DatabaseFail(t *testing.T) {
	// Set dirty state to make database check fail
	setDirtySchemaState(t)
	defer cleanupSchemaState(t) // Cleanup after test

	ds := createTestDatastore()

	// Create both health checkers
	dbHealthChecker := NewDatabaseHealthChecker(ds)
	mrHealthChecker := NewModelRegistryHealthChecker(sharedModelRegistryService)
	handler := GeneralReadinessHandler(ds, dbHealthChecker, mrHealthChecker)

	req, err := http.NewRequest("GET", "/readyz/health?format=json", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)

	// Parse and validate JSON response
	var healthStatus HealthStatus
	err = json.Unmarshal(rr.Body.Bytes(), &healthStatus)
	require.NoError(t, err)

	assert.Equal(t, StatusFail, healthStatus.Status)

	// Database should fail
	dbCheck := healthStatus.Checks[HealthCheckDatabase]
	assert.Equal(t, StatusFail, dbCheck.Status)
	assert.Contains(t, dbCheck.Message, "database schema is in dirty state")

	// Model registry should still pass (if database connection works for queries)
	mrCheck := healthStatus.Checks[HealthCheckModelRegistry]
	assert.Equal(t, StatusPass, mrCheck.Status)
}

func TestGeneralReadinessHandler_WithModelRegistry_ModelRegistryNil(t *testing.T) {
	// Ensure clean state before test
	cleanupSchemaState(t)

	ds := createTestDatastore()

	// Create health checkers - with nil model registry service
	dbHealthChecker := NewDatabaseHealthChecker(ds)
	mrHealthChecker := NewModelRegistryHealthChecker(nil)
	handler := GeneralReadinessHandler(ds, dbHealthChecker, mrHealthChecker)

	req, err := http.NewRequest("GET", "/readyz/health?format=json", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)

	// Parse and validate JSON response
	var healthStatus HealthStatus
	err = json.Unmarshal(rr.Body.Bytes(), &healthStatus)
	require.NoError(t, err)

	assert.Equal(t, StatusFail, healthStatus.Status)

	// Database should pass
	dbCheck := healthStatus.Checks[HealthCheckDatabase]
	assert.Equal(t, StatusPass, dbCheck.Status)

	// Model registry should fail
	mrCheck := healthStatus.Checks[HealthCheckModelRegistry]
	assert.Equal(t, StatusFail, mrCheck.Status)
	assert.Equal(t, "model registry service not available", mrCheck.Message)
}

func TestGeneralReadinessHandler_SimpleTextResponse_Failure(t *testing.T) {
	// Set dirty state to make database check fail
	setDirtySchemaState(t)
	defer cleanupSchemaState(t) // Cleanup after test

	ds := createTestDatastore()

	// Create both health checkers
	dbHealthChecker := NewDatabaseHealthChecker(ds)
	mrHealthChecker := NewModelRegistryHealthChecker(sharedModelRegistryService)
	handler := GeneralReadinessHandler(ds, dbHealthChecker, mrHealthChecker)

	req, err := http.NewRequest("GET", "/readyz/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
	// Should return the first failed check's error message
	assert.Contains(t, rr.Body.String(), "database schema is in dirty state")
}

func TestDatabaseHealthChecker_EmptyDSN(t *testing.T) {
	ds := datastore.Datastore{
		Type: "embedmd",
		EmbedMD: embedmd.EmbedMDConfig{
			DatabaseType: "mysql",
			DatabaseDSN:  "", // Empty DSN
		},
	}

	checker := NewDatabaseHealthChecker(ds)
	result := checker.Check()

	assert.Equal(t, HealthCheckDatabase, result.Name)
	assert.Equal(t, StatusFail, result.Status)
	assert.Equal(t, "database DSN not configured", result.Message)
}

func TestGeneralReadinessHandler_MultipleFailures(t *testing.T) {
	// Test with both database and model registry failing
	setDirtySchemaState(t)
	defer cleanupSchemaState(t)

	ds := createTestDatastore()

	dbHealthChecker := NewDatabaseHealthChecker(ds)
	mrHealthChecker := NewModelRegistryHealthChecker(nil) // Nil service to make it fail
	handler := GeneralReadinessHandler(ds, dbHealthChecker, mrHealthChecker)

	req, err := http.NewRequest("GET", "/readyz/health?format=json", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)

	var healthStatus HealthStatus
	err = json.Unmarshal(rr.Body.Bytes(), &healthStatus)
	require.NoError(t, err)

	assert.Equal(t, StatusFail, healthStatus.Status)

	// Both checks should fail
	dbCheck := healthStatus.Checks[HealthCheckDatabase]
	assert.Equal(t, StatusFail, dbCheck.Status)

	mrCheck := healthStatus.Checks[HealthCheckModelRegistry]
	assert.Equal(t, StatusFail, mrCheck.Status)
}
