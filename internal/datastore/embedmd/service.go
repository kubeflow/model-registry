package embedmd

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/db"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/db/types"
	"github.com/kubeflow/model-registry/internal/tls"
	"gorm.io/gorm"
)

const connectorType = "embedmd"

func init() {
	datastore.Register(connectorType, func(cfg any) (datastore.Connector, error) {
		emdbCfg, ok := cfg.(*EmbedMDConfig)
		if !ok {
			return nil, fmt.Errorf("invalid EmbedMD config type (%T)", cfg)
		}

		if err := emdbCfg.Validate(); err != nil {
			return nil, fmt.Errorf("invalid EmbedMD config: %w", err)
		}

		return NewEmbedMDService(cfg.(*EmbedMDConfig))
	})
}

type EmbedMDConfig struct {
	DatabaseType string
	DatabaseDSN  string
	TLSConfig    *tls.TLSConfig

	// DB is an already connected database instance that, if provided, will
	// be used instead of making a new connection.
	DB *gorm.DB
}

func (c *EmbedMDConfig) Validate() error {
	if c.DB == nil {
		if c.DatabaseType != types.DatabaseTypeMySQL && c.DatabaseType != types.DatabaseTypePostgres {
			return fmt.Errorf("unsupported database type: %s. Supported types: %s, %s", c.DatabaseType, types.DatabaseTypeMySQL, types.DatabaseTypePostgres)
		}
	}

	return nil
}

type EmbedMDService struct {
	dbConnector db.Connector
}

func NewEmbedMDService(cfg *EmbedMDConfig) (*EmbedMDService, error) {
	if cfg.DB != nil {
		db.SetDB(cfg.DB)
	} else {
		err := db.Init(cfg.DatabaseType, cfg.DatabaseDSN, cfg.TLSConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize database connector: %w", err)
		}
	}

	dbConnector, ok := db.GetConnector()
	if !ok {
		return nil, fmt.Errorf("database connector not initialized")
	}

	return &EmbedMDService{
		dbConnector: dbConnector,
	}, nil
}

func (s *EmbedMDService) Connect(spec *datastore.Spec) (datastore.RepoSet, error) {
	glog.Infof("Connecting to EmbedMD service...")

	connectedDB, err := s.dbConnector.Connect()
	if err != nil {
		return nil, err
	}

	glog.Infof("Connected to EmbedMD service")

	migrator, err := db.NewDBMigrator(connectedDB)
	if err != nil {
		return nil, err
	}

	glog.Infof("Running migrations...")

	err = migrator.Migrate()
	if err != nil {
		return nil, err
	}

	glog.Infof("Migrations completed")

	glog.Infof("Syncing types...")
	err = s.syncTypes(connectedDB, spec)
	if err != nil {
		return nil, err
	}
	glog.Infof("Syncing types completed")

	return newRepoSet(connectedDB, spec)
}

func (s EmbedMDService) Type() string {
	return connectorType
}

const (
	executionTypeKind int32 = iota
	artifactTypeKind
	contextTypeKind
)

func (s *EmbedMDService) syncTypes(conn *gorm.DB, spec *datastore.Spec) error {
	idMap := make(map[string]int32, len(spec.ExecutionTypes)+len(spec.ArtifactTypes)+len(spec.ContextTypes))
	var errs []error

	typeRepository := service.NewTypeRepository(conn)
	errs = append(errs, s.createTypes(typeRepository, spec.ExecutionTypes, executionTypeKind, idMap))
	errs = append(errs, s.createTypes(typeRepository, spec.ArtifactTypes, artifactTypeKind, idMap))
	errs = append(errs, s.createTypes(typeRepository, spec.ContextTypes, contextTypeKind, idMap))

	typePropertyRepository := service.NewTypePropertyRepository(conn)
	errs = append(errs, s.createTypeProperties(typePropertyRepository, spec.ExecutionTypes, idMap))
	errs = append(errs, s.createTypeProperties(typePropertyRepository, spec.ArtifactTypes, idMap))
	errs = append(errs, s.createTypeProperties(typePropertyRepository, spec.ContextTypes, idMap))

	return errors.Join(errs...)
}

func (s *EmbedMDService) createTypes(repo models.TypeRepository, types map[string]*datastore.SpecType, kind int32, idMap map[string]int32) error {
	var errs []error
	for name := range types {
		t, err := repo.Save(&models.TypeImpl{
			Attributes: &models.TypeAttributes{
				Name:     &name,
				TypeKind: &kind,
			},
		})
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: unable to create type: %w", name, err))
			continue
		}

		id := t.GetID()
		if id == nil {
			errs = append(errs, fmt.Errorf("%s: unable to determine type ID", name))
			continue
		}
		idMap[name] = *id
	}

	return errors.Join(errs...)
}

func (s *EmbedMDService) createTypeProperties(repo models.TypePropertyRepository, types map[string]*datastore.SpecType, idMap map[string]int32) error {
	var errs []error
	for typeName, typeSpec := range types {
		typeID := idMap[typeName]
		if typeID == 0 {
			errs = append(errs, fmt.Errorf("%s: unknown type", typeName))
			continue
		}

		for name, dataType := range typeSpec.Properties {
			_, err := repo.Save(&models.TypePropertyImpl{
				TypeID:   typeID,
				Name:     name,
				DataType: apiutils.Of(int32(dataType)),
			})
			if err != nil {
				errs = append(errs, fmt.Errorf("%s-%s: %w", typeName, name, err))
			}
		}
	}

	return errors.Join(errs...)
}
