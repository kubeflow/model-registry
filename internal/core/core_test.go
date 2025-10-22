package core_test

import (
	"os"
	"testing"

	"github.com/kubeflow/model-registry/internal/core"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	os.Exit(testutils.TestMainHelper(m))
}

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	db, dbCleanup := testutils.SetupMySQLWithMigrations(t, service.DatastoreSpec())

	// Clean up test data before each test
	testutils.CleanupTestData(t, db)

	// Return combined cleanup function
	cleanup := func() {
		// Clean up test data after each test
		testutils.CleanupTestData(t, db)
		dbCleanup()
	}

	return db, cleanup
}

// getTypeIDs retrieves all type IDs from the database for testing
func getTypeIDs(t *testing.T, db *gorm.DB) map[string]int32 {
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
		defaults.MetricHistoryTypeName,
		defaults.ParameterTypeName,
	}

	for _, typeName := range typeNames {
		var typeRecord schema.Type
		err := db.Where("name = ?", typeName).First(&typeRecord).Error
		require.NoError(t, err, "Failed to find type: %s", typeName)
		typesMap[typeName] = typeRecord.ID
	}

	return typesMap
}

// createModelRegistryService creates a ModelRegistryService from a database instance
func createModelRegistryService(t *testing.T, db *gorm.DB) *core.ModelRegistryService {
	// Get all type IDs from the database
	typesMap := getTypeIDs(t, db)

	// Create all repositories
	artifactRepo := service.NewArtifactRepository(db, map[string]int32{
		defaults.ModelArtifactTypeName: typesMap[defaults.ModelArtifactTypeName],
		defaults.DocArtifactTypeName:   typesMap[defaults.DocArtifactTypeName],
		defaults.DataSetTypeName:       typesMap[defaults.DataSetTypeName],
		defaults.MetricTypeName:        typesMap[defaults.MetricTypeName],
		defaults.ParameterTypeName:     typesMap[defaults.ParameterTypeName],
		defaults.MetricHistoryTypeName: typesMap[defaults.MetricHistoryTypeName],
	})
	modelArtifactRepo := service.NewModelArtifactRepository(db, typesMap[defaults.ModelArtifactTypeName])
	docArtifactRepo := service.NewDocArtifactRepository(db, typesMap[defaults.DocArtifactTypeName])
	registeredModelRepo := service.NewRegisteredModelRepository(db, typesMap[defaults.RegisteredModelTypeName])
	modelVersionRepo := service.NewModelVersionRepository(db, typesMap[defaults.ModelVersionTypeName])
	servingEnvironmentRepo := service.NewServingEnvironmentRepository(db, typesMap[defaults.ServingEnvironmentTypeName])
	inferenceServiceRepo := service.NewInferenceServiceRepository(db, typesMap[defaults.InferenceServiceTypeName])
	serveModelRepo := service.NewServeModelRepository(db, typesMap[defaults.ServeModelTypeName])
	experimentRepo := service.NewExperimentRepository(db, typesMap[defaults.ExperimentTypeName])
	experimentRunRepo := service.NewExperimentRunRepository(db, typesMap[defaults.ExperimentRunTypeName])
	dataSetRepo := service.NewDataSetRepository(db, typesMap[defaults.DataSetTypeName])
	metricRepo := service.NewMetricRepository(db, typesMap[defaults.MetricTypeName])
	parameterRepo := service.NewParameterRepository(db, typesMap[defaults.ParameterTypeName])
	metricHistoryRepo := service.NewMetricHistoryRepository(db, typesMap[defaults.MetricHistoryTypeName])

	// Create the core service
	return core.NewModelRegistryService(
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
}

// SetupModelRegistryService creates a complete ModelRegistryService with all repositories for testing
// This now uses the shared database infrastructure from testutils
func SetupModelRegistryService(t *testing.T) (*core.ModelRegistryService, func()) {
	db, cleanup := setupTestDB(t)

	// Create the core service
	service := createModelRegistryService(t, db)

	return service, cleanup
}
