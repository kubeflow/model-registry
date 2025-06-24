package embedmd

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/core"
	"github.com/kubeflow/model-registry/internal/db"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/internal/tls"
	"github.com/kubeflow/model-registry/pkg/api"
)

const (
	DatabaseTypeMySQL = "mysql"
)

type EmbedMDConfig struct {
	DatabaseType string
	DatabaseDSN  string
	TLSConfig    *tls.TLSConfig
}

func (c *EmbedMDConfig) Validate() error {
	if c.DatabaseType != DatabaseTypeMySQL {
		return fmt.Errorf("unsupported database type: %s", c.DatabaseType)
	}

	return nil
}

type EmbedMDService struct {
	*EmbedMDConfig
	dbConnector db.Connector
}

func NewEmbedMDService(cfg *EmbedMDConfig) (*EmbedMDService, error) {
	dbConnector, err := db.NewConnector(cfg.DatabaseType, cfg.DatabaseDSN, cfg.TLSConfig)
	if err != nil {
		return nil, err
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

	artifactRepository := service.NewArtifactRepository(connectedDB, typesMap[defaults.ModelArtifactTypeName], typesMap[defaults.DocArtifactTypeName])
	modelArtifactRepository := service.NewModelArtifactRepository(connectedDB, typesMap[defaults.ModelArtifactTypeName])
	docArtifactRepository := service.NewDocArtifactRepository(connectedDB, typesMap[defaults.DocArtifactTypeName])
	registeredModelRepository := service.NewRegisteredModelRepository(connectedDB, typesMap[defaults.RegisteredModelTypeName])
	modelVersionRepository := service.NewModelVersionRepository(connectedDB, typesMap[defaults.ModelVersionTypeName])
	servingEnvironmentRepository := service.NewServingEnvironmentRepository(connectedDB, typesMap[defaults.ServingEnvironmentTypeName])
	inferenceServiceRepository := service.NewInferenceServiceRepository(connectedDB, typesMap[defaults.InferenceServiceTypeName])
	serveModelRepository := service.NewServeModelRepository(connectedDB, typesMap[defaults.ServeModelTypeName])

	modelRegistryService := core.NewModelRegistryService(
		artifactRepository,
		modelArtifactRepository,
		docArtifactRepository,
		registeredModelRepository,
		modelVersionRepository,
		servingEnvironmentRepository,
		inferenceServiceRepository,
		serveModelRepository,
		typesMap,
	)

	glog.Infof("EmbedMD service connected")

	return modelRegistryService, nil
}

func (s *EmbedMDService) Teardown() error {
	return nil
}
