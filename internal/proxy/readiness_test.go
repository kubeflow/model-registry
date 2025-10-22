package proxy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/core"
	"github.com/kubeflow/model-registry/internal/db"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/kubeflow/model-registry/internal/tls"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	os.Exit(testutils.TestMainHelper(m))
}

func setupTestDB(t *testing.T) (*gorm.DB, string, api.ModelRegistryApi, func()) {
	sharedDB, cleanup := testutils.SetupMySQLWithMigrations(t, service.DatastoreSpec())
	dsn := testutils.GetSharedMySQLDSN(t)
	svc := setupModelRegistryService(sharedDB)

	// Initialize global db connector for health checks
	err := db.Init("mysql", dsn, &tls.TLSConfig{})
	require.NoError(t, err)

	return sharedDB, dsn, svc, cleanup
}

// getTypeIDs retrieves all type IDs from the database for testing
func getTypeIDs(sharedDB *gorm.DB) map[string]int32 {
	typesMap := map[string]int32{}

	typeNames := []string{
		defaults.RegisteredModelTypeName,
		defaults.ModelVersionTypeName,
		defaults.DocArtifactTypeName,
		defaults.ModelArtifactTypeName,
		defaults.ServingEnvironmentTypeName,
		defaults.InferenceServiceTypeName,
		defaults.ServeModelTypeName,
		defaults.ExperimentTypeName,
		defaults.ExperimentRunTypeName,
		defaults.DataSetTypeName,
		defaults.MetricTypeName,
		defaults.ParameterTypeName,
		defaults.MetricHistoryTypeName,
	}

	for _, typeName := range typeNames {
		var typeRecord schema.Type
		err := sharedDB.Where("name = ?", typeName).First(&typeRecord).Error
		if err != nil {
			panic("Failed to find type: " + typeName + ": " + err.Error())
		}
		typesMap[typeName] = typeRecord.ID
	}

	return typesMap
}

// setupModelRegistryService creates a complete ModelRegistryService with all repositories for testing
func setupModelRegistryService(sharedDB *gorm.DB) api.ModelRegistryApi {
	// Get all type IDs from the database
	typesMap := getTypeIDs(sharedDB)

	// Create all repositories
	artifactRepo := service.NewArtifactRepository(sharedDB, map[string]int32{
		defaults.ModelArtifactTypeName: typesMap[defaults.ModelArtifactTypeName],
		defaults.DocArtifactTypeName:   typesMap[defaults.DocArtifactTypeName],
		defaults.DataSetTypeName:       typesMap[defaults.DataSetTypeName],
		defaults.MetricTypeName:        typesMap[defaults.MetricTypeName],
		defaults.ParameterTypeName:     typesMap[defaults.ParameterTypeName],
		defaults.MetricHistoryTypeName: typesMap[defaults.MetricHistoryTypeName],
	})
	modelArtifactRepo := service.NewModelArtifactRepository(sharedDB, typesMap[defaults.ModelArtifactTypeName])
	docArtifactRepo := service.NewDocArtifactRepository(sharedDB, typesMap[defaults.DocArtifactTypeName])
	registeredModelRepo := service.NewRegisteredModelRepository(sharedDB, typesMap[defaults.RegisteredModelTypeName])
	modelVersionRepo := service.NewModelVersionRepository(sharedDB, typesMap[defaults.ModelVersionTypeName])
	servingEnvironmentRepo := service.NewServingEnvironmentRepository(sharedDB, typesMap[defaults.ServingEnvironmentTypeName])
	inferenceServiceRepo := service.NewInferenceServiceRepository(sharedDB, typesMap[defaults.InferenceServiceTypeName])
	serveModelRepo := service.NewServeModelRepository(sharedDB, typesMap[defaults.ServeModelTypeName])
	experimentRepo := service.NewExperimentRepository(sharedDB, typesMap[defaults.ExperimentTypeName])
	experimentRunRepo := service.NewExperimentRunRepository(sharedDB, typesMap[defaults.ExperimentRunTypeName])
	dataSetRepo := service.NewDataSetRepository(sharedDB, typesMap[defaults.DataSetTypeName])
	metricRepo := service.NewMetricRepository(sharedDB, typesMap[defaults.MetricTypeName])
	parameterRepo := service.NewParameterRepository(sharedDB, typesMap[defaults.ParameterTypeName])
	metricHistoryRepo := service.NewMetricHistoryRepository(sharedDB, typesMap[defaults.MetricHistoryTypeName])

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
		experimentRepo,
		experimentRunRepo,
		dataSetRepo,
		metricRepo,
		parameterRepo,
		metricHistoryRepo,
		typesMap,
	)

	return service
}

// cleanupSchemaState resets schema_migrations table to clean state
func cleanupSchemaState(t *testing.T, sharedDB *gorm.DB) {
	// Reset schema_migrations to clean state
	err := sharedDB.Exec("UPDATE schema_migrations SET dirty = 0").Error

	require.NoError(t, err)
}

// setDirtySchemaState sets schema_migrations to dirty state for testing
func setDirtySchemaState(t *testing.T, sharedDB *gorm.DB) {
	err := sharedDB.Exec("UPDATE schema_migrations SET dirty = 1").Error
	require.NoError(t, err)
}

func TestReadinessHandler_EmbedMD_Success(t *testing.T) {
	// Ensure clean state before test
	sharedDB, _, _, cleanup := setupTestDB(t)
	defer cleanup()

	cleanupSchemaState(t, sharedDB)

	dbHealthChecker := NewDatabaseHealthChecker()
	handler := GeneralReadinessHandler(dbHealthChecker)
	req, err := http.NewRequest("GET", "/readyz/isDirty", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, responseOK, rr.Body.String())
}

func TestReadinessHandler_EmbedMD_Dirty(t *testing.T) {
	// Set dirty state for this test
	sharedDB, _, _, cleanup := setupTestDB(t)
	defer cleanup()

	setDirtySchemaState(t, sharedDB)
	defer cleanupSchemaState(t, sharedDB)

	dbHealthChecker := NewDatabaseHealthChecker()
	handler := GeneralReadinessHandler(dbHealthChecker)
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
	sharedDB, _, sharedModelRegistryService, cleanup := setupTestDB(t)
	defer cleanup()

	cleanupSchemaState(t, sharedDB)

	// Create both health checkers
	dbHealthChecker := NewDatabaseHealthChecker()
	mrHealthChecker := NewModelRegistryHealthChecker(sharedModelRegistryService)
	handler := GeneralReadinessHandler(dbHealthChecker, mrHealthChecker)

	req, err := http.NewRequest("GET", "/readyz/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, responseOK, rr.Body.String())
}

func TestGeneralReadinessHandler_WithModelRegistry_JSONFormat(t *testing.T) {
	// Ensure clean state before test
	sharedDB, _, sharedModelRegistryService, cleanup := setupTestDB(t)
	defer cleanup()

	cleanupSchemaState(t, sharedDB)

	// Create both health checkers
	dbHealthChecker := NewDatabaseHealthChecker()
	mrHealthChecker := NewModelRegistryHealthChecker(sharedModelRegistryService)
	handler := GeneralReadinessHandler(dbHealthChecker, mrHealthChecker)

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
	sharedDB, _, sharedModelRegistryService, cleanup := setupTestDB(t)
	defer cleanup()

	setDirtySchemaState(t, sharedDB)
	defer cleanupSchemaState(t, sharedDB)

	// Create both health checkers
	dbHealthChecker := NewDatabaseHealthChecker()
	mrHealthChecker := NewModelRegistryHealthChecker(sharedModelRegistryService)
	handler := GeneralReadinessHandler(dbHealthChecker, mrHealthChecker)

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
	sharedDB, _, _, cleanup := setupTestDB(t)
	defer cleanup()

	cleanupSchemaState(t, sharedDB)

	// Create health checkers - with nil model registry service
	dbHealthChecker := NewDatabaseHealthChecker()
	mrHealthChecker := NewModelRegistryHealthChecker(nil)
	handler := GeneralReadinessHandler(dbHealthChecker, mrHealthChecker)

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
	sharedDB, _, sharedModelRegistryService, cleanup := setupTestDB(t)
	defer cleanup()

	setDirtySchemaState(t, sharedDB)
	defer cleanupSchemaState(t, sharedDB)

	// Create both health checkers
	dbHealthChecker := NewDatabaseHealthChecker()
	mrHealthChecker := NewModelRegistryHealthChecker(sharedModelRegistryService)
	handler := GeneralReadinessHandler(dbHealthChecker, mrHealthChecker)

	req, err := http.NewRequest("GET", "/readyz/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
	// Should return the first failed check's error message
	assert.Contains(t, rr.Body.String(), "database schema is in dirty state")
}

func TestGeneralReadinessHandler_MultipleFailures(t *testing.T) {
	// Test with both database and model registry failing
	sharedDB, _, _, cleanup := setupTestDB(t)
	defer cleanup()

	setDirtySchemaState(t, sharedDB)
	defer cleanupSchemaState(t, sharedDB)

	dbHealthChecker := NewDatabaseHealthChecker()
	mrHealthChecker := NewModelRegistryHealthChecker(nil) // Nil service to make it fail
	handler := GeneralReadinessHandler(dbHealthChecker, mrHealthChecker)

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
