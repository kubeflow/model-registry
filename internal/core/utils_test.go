package core_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/kubeflow/model-registry/internal/core"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd/mysql"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/defaults"
	_tls "github.com/kubeflow/model-registry/internal/tls"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	cont_mysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/gorm"
)

// Package-level shared database instance
var (
	sharedDB       *gorm.DB
	mysqlContainer *cont_mysql.MySQLContainer
	_service       *core.ModelRegistryService
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
		// Enable SSL with default certificates
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
	dbConnector := mysql.NewMySQLDBConnector(mysqlContainer.MustConnectionString(ctx), &_tls.TLSConfig{})
	sharedDB, err = dbConnector.Connect()
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
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

	_service = SetupModelRegistryService()

	// Run all tests
	code := m.Run()

	os.Exit(code)
}

func cleanupTestData(t *testing.T, db *gorm.DB) {
	tables := []string{
		"Context",
		"ContextProperty",
		"Association",
		"Attribution",
		"Event",
		"EventPath",
		"Artifact",
		"ArtifactProperty",
		"Execution",
		"ExecutionProperty",
		"ParentContext",
		"ParentType",
	}

	// Disable foreign key checks temporarily
	err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error
	require.NoError(t, err)

	// Truncate all tables
	for _, table := range tables {
		err := db.Exec("TRUNCATE TABLE " + table).Error
		if err != nil {
			// Table might not exist, which is okay
			t.Logf("Could not truncate table %s: %v", table, err)
		}
	}

	// Re-enable foreign key checks
	err = db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
	require.NoError(t, err)
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

// SetupModelRegistryService creates a complete ModelRegistryService with all repositories for testing
func SetupModelRegistryService() *core.ModelRegistryService {
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
