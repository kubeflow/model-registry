package core

import (
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/mapper"
	"github.com/kubeflow/model-registry/pkg/api"
)

// Compile-time assertion to ensure ModelRegistryService implements ModelRegistryApi
var _ api.ModelRegistryApi = (*ModelRegistryService)(nil)

type ModelRegistryService struct {
	artifactRepository           models.ArtifactRepository
	modelArtifactRepository      models.ModelArtifactRepository
	docArtifactRepository        models.DocArtifactRepository
	registeredModelRepository    models.RegisteredModelRepository
	modelVersionRepository       models.ModelVersionRepository
	servingEnvironmentRepository models.ServingEnvironmentRepository
	inferenceServiceRepository   models.InferenceServiceRepository
	serveModelRepository         models.ServeModelRepository
	experimentRepository         models.ExperimentRepository
	experimentRunRepository      models.ExperimentRunRepository
	dataSetRepository            models.DataSetRepository
	metricRepository             models.MetricRepository
	parameterRepository          models.ParameterRepository
	metricHistoryRepository      models.MetricHistoryRepository
	mapper                       mapper.EmbedMDMapper
	typesMap                     map[string]int64
}

func NewModelRegistryService(
	artifactRepository models.ArtifactRepository,
	modelArtifactRepository models.ModelArtifactRepository,
	docArtifactRepository models.DocArtifactRepository,
	registeredModelRepository models.RegisteredModelRepository,
	modelVersionRepository models.ModelVersionRepository,
	servingEnvironmentRepository models.ServingEnvironmentRepository,
	inferenceServiceRepository models.InferenceServiceRepository,
	serveModelRepository models.ServeModelRepository,
	experimentRepository models.ExperimentRepository,
	experimentRunRepository models.ExperimentRunRepository,
	dataSetRepository models.DataSetRepository,
	metricRepository models.MetricRepository,
	parameterRepository models.ParameterRepository,
	metricHistoryRepository models.MetricHistoryRepository,
	typesMap map[string]int64) *ModelRegistryService {
	return &ModelRegistryService{
		artifactRepository:           artifactRepository,
		modelArtifactRepository:      modelArtifactRepository,
		docArtifactRepository:        docArtifactRepository,
		registeredModelRepository:    registeredModelRepository,
		modelVersionRepository:       modelVersionRepository,
		servingEnvironmentRepository: servingEnvironmentRepository,
		inferenceServiceRepository:   inferenceServiceRepository,
		serveModelRepository:         serveModelRepository,
		experimentRepository:         experimentRepository,
		experimentRunRepository:      experimentRunRepository,
		dataSetRepository:            dataSetRepository,
		metricRepository:             metricRepository,
		parameterRepository:          parameterRepository,
		metricHistoryRepository:      metricHistoryRepository,
		mapper:                       *mapper.NewEmbedMDMapper(typesMap),
		typesMap:                     typesMap,
	}
}
