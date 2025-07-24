package service_test

import (
	"os"
	"testing"

	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	os.Exit(testutils.TestMainHelper(m))
}

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	db, dbCleanup := testutils.SetupMySQLWithMigrations(t)

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

// Helper functions to get type IDs from the database
func getRegisteredModelTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.RegisteredModelTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find RegisteredModel type")
	return int64(typeRecord.ID)
}

func getModelVersionTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.ModelVersionTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find ModelVersion type")
	return int64(typeRecord.ID)
}

func getModelArtifactTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.ModelArtifactTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find ModelArtifact type")
	return int64(typeRecord.ID)
}

func getDocArtifactTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.DocArtifactTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find DocArtifact type")
	return int64(typeRecord.ID)
}

func getServingEnvironmentTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.ServingEnvironmentTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find ServingEnvironment type")
	return int64(typeRecord.ID)
}

func getInferenceServiceTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.InferenceServiceTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find InferenceService type")
	return int64(typeRecord.ID)
}

func getServeModelTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.ServeModelTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find ServeModel type")
	return int64(typeRecord.ID)
}

func getExperimentTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.ExperimentTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find Experiment type")
	return int64(typeRecord.ID)
}

func getExperimentRunTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.ExperimentRunTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find ExperimentRun type")
	return int64(typeRecord.ID)
}

func getDataSetTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.DataSetTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find DataSet type")
	return int64(typeRecord.ID)
}

func getMetricTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.MetricTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find Metric type")
	return int64(typeRecord.ID)
}

func getParameterTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.ParameterTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find Parameter type")
	return int64(typeRecord.ID)
}

func getMetricHistoryTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.MetricHistoryTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find MetricHistory type")
	return int64(typeRecord.ID)
}
