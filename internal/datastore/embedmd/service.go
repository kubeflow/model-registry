package embedmd

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/core"
	"github.com/kubeflow/model-registry/internal/db"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/db/types"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/internal/tls"
	"github.com/kubeflow/model-registry/pkg/api"
)

type EmbedMDConfig struct {
	DatabaseType string
	DatabaseDSN  string
	TLSConfig    *tls.TLSConfig
}

func (c *EmbedMDConfig) Validate() error {
	if c.DatabaseType != types.DatabaseTypeMySQL && c.DatabaseType != types.DatabaseTypePostgres {
		return fmt.Errorf("unsupported database type: %s. Supported types: %s, %s", c.DatabaseType, types.DatabaseTypeMySQL, types.DatabaseTypePostgres)
	}

	return nil
}

type EmbedMDService struct {
	*EmbedMDConfig
	dbConnector db.Connector
}

func NewEmbedMDService(cfg *EmbedMDConfig) (*EmbedMDService, error) {
	err := db.Init(cfg.DatabaseType, cfg.DatabaseDSN, cfg.TLSConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database connector: %w", err)
	}

	dbConnector, ok := db.GetConnector()
	if !ok {
		return nil, fmt.Errorf("database connector not initialized")
	}

	return &EmbedMDService{
		EmbedMDConfig: cfg,
		dbConnector:   dbConnector,
	}, nil
}

func (s *EmbedMDService) Connect() (api.ModelRegistryApi, error) {
	glog.Infof("Connecting to EmbedMD service...")

	connectedDB, err := s.dbConnector.Connect()
	if err != nil {
		return nil, err
	}

	glog.Infof("Connected to EmbedMD service")

	migrator, err := db.NewDBMigrator(s.DatabaseType, connectedDB)
	if err != nil {
		return nil, err
	}

	glog.Infof("Running migrations...")

	err = migrator.Migrate()
	if err != nil {
		return nil, err
	}

	glog.Infof("Migrations completed")

	typeRepository := service.NewTypeRepository(connectedDB)

	glog.Infof("Getting types...")

	types, err := typeRepository.GetAll()
	if err != nil {
		return nil, err
	}

	typesMap := make(map[string]int64)

	for _, t := range types {
		typesMap[*t.GetAttributes().Name] = int64(*t.GetID())
	}

	glog.Infof("Types retrieved")

	// Add debug logging to see what types are actually available
	glog.V(2).Infof("DEBUG: Available types in typesMap:")
	for typeName, typeID := range typesMap {
		glog.V(2).Infof("  %s = %d", typeName, typeID)
	}

	// Validate that all required types are registered
	requiredTypes := []string{
		defaults.ModelArtifactTypeName,
		defaults.DocArtifactTypeName,
		defaults.DataSetTypeName,
		defaults.MetricTypeName,
		defaults.ParameterTypeName,
		defaults.MetricHistoryTypeName,
		defaults.RegisteredModelTypeName,
		defaults.ModelVersionTypeName,
		defaults.ServingEnvironmentTypeName,
		defaults.InferenceServiceTypeName,
		defaults.ServeModelTypeName,
		defaults.ExperimentTypeName,
		defaults.ExperimentRunTypeName,
	}

	for _, requiredType := range requiredTypes {
		if _, exists := typesMap[requiredType]; !exists {
			return nil, fmt.Errorf("required type '%s' not found in database. Please ensure all migrations have been applied", requiredType)
		}
	}

	glog.Infof("All required types validated successfully")

	artifactRepository := service.NewArtifactRepository(connectedDB, typesMap[defaults.ModelArtifactTypeName], typesMap[defaults.DocArtifactTypeName], typesMap[defaults.DataSetTypeName], typesMap[defaults.MetricTypeName], typesMap[defaults.ParameterTypeName], typesMap[defaults.MetricHistoryTypeName])
	modelArtifactRepository := service.NewModelArtifactRepository(connectedDB, typesMap[defaults.ModelArtifactTypeName])
	docArtifactRepository := service.NewDocArtifactRepository(connectedDB, typesMap[defaults.DocArtifactTypeName])
	registeredModelRepository := service.NewRegisteredModelRepository(connectedDB, typesMap[defaults.RegisteredModelTypeName])
	modelVersionRepository := service.NewModelVersionRepository(connectedDB, typesMap[defaults.ModelVersionTypeName])
	servingEnvironmentRepository := service.NewServingEnvironmentRepository(connectedDB, typesMap[defaults.ServingEnvironmentTypeName])
	inferenceServiceRepository := service.NewInferenceServiceRepository(connectedDB, typesMap[defaults.InferenceServiceTypeName])
	serveModelRepository := service.NewServeModelRepository(connectedDB, typesMap[defaults.ServeModelTypeName])
	experimentRepository := service.NewExperimentRepository(connectedDB, typesMap[defaults.ExperimentTypeName])
	experimentRunRepository := service.NewExperimentRunRepository(connectedDB, typesMap[defaults.ExperimentRunTypeName])

	dataSetRepository := service.NewDataSetRepository(connectedDB, typesMap[defaults.DataSetTypeName])
	metricRepository := service.NewMetricRepository(connectedDB, typesMap[defaults.MetricTypeName])
	parameterRepository := service.NewParameterRepository(connectedDB, typesMap[defaults.ParameterTypeName])
	metricHistoryRepository := service.NewMetricHistoryRepository(connectedDB, typesMap[defaults.MetricHistoryTypeName])

	modelRegistryService := core.NewModelRegistryService(
		artifactRepository,
		modelArtifactRepository,
		docArtifactRepository,
		registeredModelRepository,
		modelVersionRepository,
		servingEnvironmentRepository,
		inferenceServiceRepository,
		serveModelRepository,
		experimentRepository,
		experimentRunRepository,
		dataSetRepository,
		metricRepository,
		parameterRepository,
		metricHistoryRepository,
		typesMap,
	)

	glog.Infof("EmbedMD service connected")

	return modelRegistryService, nil
}

func (s *EmbedMDService) Teardown() error {
	return nil
}
