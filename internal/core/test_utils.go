package core

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/internal/tls"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	cont_mysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/gorm"
)

// setupTestDB creates a MySQL testcontainer with migrations for testing
func setupTestDB(t *testing.T) (*gorm.DB, func()) {
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

	dbConnector := mysql.NewMySQLDBConnector(mysqlContainer.MustConnectionString(ctx), &tls.TLSConfig{})
	require.NoError(t, err)

	db, err := dbConnector.Connect()
	require.NoError(t, err)

	// Run migrations
	migrator, err := mysql.NewMySQLMigrator(db)
	require.NoError(t, err)
	err = migrator.Migrate()
	require.NoError(t, err)

	// Return cleanup function
	cleanup := func() {
		sqlDB, err := db.DB()
		require.NoError(t, err)
		sqlDB.Close() //nolint:errcheck
		err = testcontainers.TerminateContainer(mysqlContainer)
		require.NoError(t, err)
	}

	return db, cleanup
}

// getTypeIDs retrieves all type IDs from the database for testing
func getTypeIDs(t *testing.T, db *gorm.DB) map[string]int64 {
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
		err := db.Where("name = ?", typeName).First(&typeRecord).Error
		require.NoError(t, err, "Failed to find type: %s", typeName)
		typesMap[typeName] = int64(typeRecord.ID)
	}

	return typesMap
}

// SetupModelRegistryService creates a complete ModelRegistryService with all repositories for testing
func SetupModelRegistryService(t *testing.T) (*ModelRegistryService, func()) {
	db, cleanup := setupTestDB(t)

	// Get all type IDs from the database
	typesMap := getTypeIDs(t, db)

	// Create all repositories
	artifactRepo := service.NewArtifactRepository(db, typesMap[defaults.ModelArtifactTypeName], typesMap[defaults.DocArtifactTypeName])
	modelArtifactRepo := service.NewModelArtifactRepository(db, typesMap[defaults.ModelArtifactTypeName])
	docArtifactRepo := service.NewDocArtifactRepository(db, typesMap[defaults.DocArtifactTypeName])
	registeredModelRepo := service.NewRegisteredModelRepository(db, typesMap[defaults.RegisteredModelTypeName])
	modelVersionRepo := service.NewModelVersionRepository(db, typesMap[defaults.ModelVersionTypeName])
	servingEnvironmentRepo := service.NewServingEnvironmentRepository(db, typesMap[defaults.ServingEnvironmentTypeName])
	inferenceServiceRepo := service.NewInferenceServiceRepository(db, typesMap[defaults.InferenceServiceTypeName])
	serveModelRepo := service.NewServeModelRepository(db, typesMap[defaults.ServeModelTypeName])

	// Create the core service
	service := NewModelRegistryService(
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

	return service, cleanup
}
