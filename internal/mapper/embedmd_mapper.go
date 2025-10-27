package mapper

import (
	"fmt"

	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/converter/generated"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

type EmbedMDMapper struct {
	openAPIConverter converter.OpenAPIToEmbedMDConverter
	embedMDConverter converter.EmbedMDToOpenAPIConverter
	*generated.OpenAPIConverterImpl
	*generated.OpenAPIReconcilerImpl
	typesMap map[string]int32
}

func NewEmbedMDMapper(typesMap map[string]int32) *EmbedMDMapper {
	return &EmbedMDMapper{
		openAPIConverter:      &generated.OpenAPIToEmbedMDConverterImpl{},
		embedMDConverter:      &generated.EmbedMDToOpenAPIConverterImpl{},
		OpenAPIConverterImpl:  &generated.OpenAPIConverterImpl{},
		OpenAPIReconcilerImpl: &generated.OpenAPIReconcilerImpl{},
		typesMap:              typesMap,
	}
}

// Utilities for OpenAPI --> EmbedMD mapping, make use of generated Converters

func (e *EmbedMDMapper) MapFromRegisteredModel(registeredModel *openapi.RegisteredModel) (models.RegisteredModel, error) {
	return e.openAPIConverter.ConvertRegisteredModel(&converter.OpenAPIModelWrapper[openapi.RegisteredModel]{
		TypeId: e.typesMap[defaults.RegisteredModelTypeName],
		Model:  registeredModel,
	})
}

func (e *EmbedMDMapper) MapFromModelVersion(modelVersion *openapi.ModelVersion, parentResourceId *string) (models.ModelVersion, error) {
	return e.openAPIConverter.ConvertModelVersion(&converter.OpenAPIModelWrapper[openapi.ModelVersion]{
		TypeId:           e.typesMap[defaults.ModelVersionTypeName],
		Model:            modelVersion,
		ParentResourceId: parentResourceId,
	})
}

func (e *EmbedMDMapper) MapFromServingEnvironment(servingEnvironment *openapi.ServingEnvironment) (models.ServingEnvironment, error) {
	return e.openAPIConverter.ConvertServingEnvironment(&converter.OpenAPIModelWrapper[openapi.ServingEnvironment]{
		TypeId: e.typesMap[defaults.ServingEnvironmentTypeName],
		Model:  servingEnvironment,
	})
}

func (e *EmbedMDMapper) MapFromInferenceService(inferenceService *openapi.InferenceService, servingEnvironmentId string) (models.InferenceService, error) {
	return e.openAPIConverter.ConvertInferenceService(&converter.OpenAPIModelWrapper[openapi.InferenceService]{
		TypeId:           e.typesMap[defaults.InferenceServiceTypeName],
		Model:            inferenceService,
		ParentResourceId: &servingEnvironmentId,
	})
}

func (e *EmbedMDMapper) MapFromModelArtifact(modelArtifact *openapi.ModelArtifact, parentResourceId *string) (models.ModelArtifact, error) {
	return e.openAPIConverter.ConvertModelArtifact(&converter.OpenAPIModelWrapper[openapi.ModelArtifact]{
		TypeId:           e.typesMap[defaults.ModelArtifactTypeName],
		Model:            modelArtifact,
		ParentResourceId: parentResourceId,
	})
}

func (e *EmbedMDMapper) MapFromDocArtifact(docArtifact *openapi.DocArtifact, parentResourceId *string) (models.DocArtifact, error) {
	return e.openAPIConverter.ConvertDocArtifact(&converter.OpenAPIModelWrapper[openapi.DocArtifact]{
		TypeId:           e.typesMap[defaults.DocArtifactTypeName],
		Model:            docArtifact,
		ParentResourceId: parentResourceId,
	})
}

func (e *EmbedMDMapper) MapFromServeModel(serveModel *openapi.ServeModel, parentResourceId *string) (models.ServeModel, error) {
	return e.openAPIConverter.ConvertServeModel(&converter.OpenAPIModelWrapper[openapi.ServeModel]{
		TypeId:           e.typesMap[defaults.ServeModelTypeName],
		Model:            serveModel,
		ParentResourceId: parentResourceId,
	})
}

func (e *EmbedMDMapper) MapFromExperiment(experiment *openapi.Experiment) (models.Experiment, error) {
	return e.openAPIConverter.ConvertExperiment(&converter.OpenAPIModelWrapper[openapi.Experiment]{
		TypeId: e.typesMap[defaults.ExperimentTypeName],
		Model:  experiment,
	})
}

func (e *EmbedMDMapper) MapFromExperimentRun(experimentRun *openapi.ExperimentRun, parentResourceId *string) (models.ExperimentRun, error) {
	return e.openAPIConverter.ConvertExperimentRun(&converter.OpenAPIModelWrapper[openapi.ExperimentRun]{
		TypeId:           e.typesMap[defaults.ExperimentRunTypeName],
		Model:            experimentRun,
		ParentResourceId: parentResourceId,
	})
}

func (e *EmbedMDMapper) MapFromMetric(metric *openapi.Metric, parentResourceId *string) (models.Metric, error) {
	return e.openAPIConverter.ConvertMetric(&converter.OpenAPIModelWrapper[openapi.Metric]{
		TypeId:           e.typesMap[defaults.MetricTypeName],
		Model:            metric,
		ParentResourceId: parentResourceId,
	})
}

func (e *EmbedMDMapper) MapFromParameter(parameter *openapi.Parameter, parentResourceId *string) (models.Parameter, error) {
	return e.openAPIConverter.ConvertParameter(&converter.OpenAPIModelWrapper[openapi.Parameter]{
		TypeId:           e.typesMap[defaults.ParameterTypeName],
		Model:            parameter,
		ParentResourceId: parentResourceId,
	})
}

func (e *EmbedMDMapper) MapFromDataSet(dataSet *openapi.DataSet) (models.DataSet, error) {
	return e.openAPIConverter.ConvertDataSet(&converter.OpenAPIModelWrapper[openapi.DataSet]{
		TypeId: e.typesMap[defaults.DataSetTypeName],
		Model:  dataSet,
	})
}

// Utilities for EmbedMD --> OpenAPI mapping, make use of generated Converters

func (e *EmbedMDMapper) MapToRegisteredModel(registeredModel models.RegisteredModel) (*openapi.RegisteredModel, error) {
	if registeredModel == nil {
		return nil, fmt.Errorf("registered model is nil")
	}

	return e.embedMDConverter.ConvertRegisteredModel(&models.RegisteredModelImpl{
		ID:               registeredModel.GetID(),
		TypeID:           registeredModel.GetTypeID(),
		Attributes:       registeredModel.GetAttributes(),
		Properties:       registeredModel.GetProperties(),
		CustomProperties: registeredModel.GetCustomProperties(),
	})
}

func (e *EmbedMDMapper) MapToModelVersion(modelVersion models.ModelVersion) (*openapi.ModelVersion, error) {
	return e.embedMDConverter.ConvertModelVersion(&models.ModelVersionImpl{
		ID:               modelVersion.GetID(),
		TypeID:           modelVersion.GetTypeID(),
		Attributes:       modelVersion.GetAttributes(),
		Properties:       modelVersion.GetProperties(),
		CustomProperties: modelVersion.GetCustomProperties(),
	})
}

func (e *EmbedMDMapper) MapToServingEnvironment(servingEnvironment models.ServingEnvironment) (*openapi.ServingEnvironment, error) {
	return e.embedMDConverter.ConvertServingEnvironment(&models.ServingEnvironmentImpl{
		ID:               servingEnvironment.GetID(),
		TypeID:           servingEnvironment.GetTypeID(),
		Attributes:       servingEnvironment.GetAttributes(),
		Properties:       servingEnvironment.GetProperties(),
		CustomProperties: servingEnvironment.GetCustomProperties(),
	})
}

func (e *EmbedMDMapper) MapToInferenceService(inferenceService models.InferenceService) (*openapi.InferenceService, error) {
	return e.embedMDConverter.ConvertInferenceService(&models.InferenceServiceImpl{
		ID:               inferenceService.GetID(),
		TypeID:           inferenceService.GetTypeID(),
		Attributes:       inferenceService.GetAttributes(),
		Properties:       inferenceService.GetProperties(),
		CustomProperties: inferenceService.GetCustomProperties(),
	})
}

func (e *EmbedMDMapper) MapToModelArtifact(modelArtifact models.ModelArtifact) (*openapi.ModelArtifact, error) {
	return e.embedMDConverter.ConvertModelArtifact(&models.ModelArtifactImpl{
		ID:               modelArtifact.GetID(),
		TypeID:           modelArtifact.GetTypeID(),
		Attributes:       modelArtifact.GetAttributes(),
		Properties:       modelArtifact.GetProperties(),
		CustomProperties: modelArtifact.GetCustomProperties(),
	})
}

func (e *EmbedMDMapper) MapToDocArtifact(docArtifact models.DocArtifact) (*openapi.DocArtifact, error) {
	return e.embedMDConverter.ConvertDocArtifact(&models.DocArtifactImpl{
		ID:               docArtifact.GetID(),
		TypeID:           docArtifact.GetTypeID(),
		Attributes:       docArtifact.GetAttributes(),
		Properties:       docArtifact.GetProperties(),
		CustomProperties: docArtifact.GetCustomProperties(),
	})
}

func (e *EmbedMDMapper) MapToServeModel(serveModel models.ServeModel) (*openapi.ServeModel, error) {
	return e.embedMDConverter.ConvertServeModel(&models.ServeModelImpl{
		ID:               serveModel.GetID(),
		TypeID:           serveModel.GetTypeID(),
		Attributes:       serveModel.GetAttributes(),
		Properties:       serveModel.GetProperties(),
		CustomProperties: serveModel.GetCustomProperties(),
	})
}

func (e *EmbedMDMapper) MapToExperiment(experiment models.Experiment) (*openapi.Experiment, error) {
	if experiment == nil {
		return nil, fmt.Errorf("experiment is nil")
	}

	return e.embedMDConverter.ConvertExperiment(&models.ExperimentImpl{
		ID:               experiment.GetID(),
		TypeID:           experiment.GetTypeID(),
		Attributes:       experiment.GetAttributes(),
		Properties:       experiment.GetProperties(),
		CustomProperties: experiment.GetCustomProperties(),
	})
}

func (e *EmbedMDMapper) MapToExperimentRun(experimentRun models.ExperimentRun) (*openapi.ExperimentRun, error) {
	if experimentRun == nil {
		return nil, fmt.Errorf("experiment run is nil")
	}

	return e.embedMDConverter.ConvertExperimentRun(&models.ExperimentRunImpl{
		ID:               experimentRun.GetID(),
		TypeID:           experimentRun.GetTypeID(),
		Attributes:       experimentRun.GetAttributes(),
		Properties:       experimentRun.GetProperties(),
		CustomProperties: experimentRun.GetCustomProperties(),
	})
}

func (e *EmbedMDMapper) MapToMetric(metricHistory models.MetricHistory) (*openapi.Metric, error) {
	if metricHistory == nil {
		return nil, fmt.Errorf("metric history is nil")
	}

	// Convert MetricHistoryAttributes to MetricAttributes
	var metricAttributes *models.MetricAttributes
	if metricHistory.GetAttributes() != nil {
		metricAttributes = &models.MetricAttributes{
			Name:                     metricHistory.GetAttributes().Name,
			URI:                      metricHistory.GetAttributes().URI,
			State:                    metricHistory.GetAttributes().State,
			ArtifactType:             metricHistory.GetAttributes().ArtifactType,
			ExternalID:               metricHistory.GetAttributes().ExternalID,
			CreateTimeSinceEpoch:     metricHistory.GetAttributes().CreateTimeSinceEpoch,
			LastUpdateTimeSinceEpoch: metricHistory.GetAttributes().LastUpdateTimeSinceEpoch,
		}
	}

	return e.embedMDConverter.ConvertMetric(&models.MetricImpl{
		ID:               metricHistory.GetID(),
		TypeID:           metricHistory.GetTypeID(),
		Attributes:       metricAttributes,
		Properties:       metricHistory.GetProperties(),
		CustomProperties: metricHistory.GetCustomProperties(),
	})
}

func (e *EmbedMDMapper) MapToMetricFromMetric(metric models.Metric) (*openapi.Metric, error) {
	if metric == nil {
		return nil, fmt.Errorf("metric is nil")
	}

	return e.embedMDConverter.ConvertMetric(&models.MetricImpl{
		ID:               metric.GetID(),
		TypeID:           metric.GetTypeID(),
		Attributes:       metric.GetAttributes(),
		Properties:       metric.GetProperties(),
		CustomProperties: metric.GetCustomProperties(),
	})
}

func (e *EmbedMDMapper) MapToDataSet(dataSet models.DataSet) (*openapi.DataSet, error) {
	if dataSet == nil {
		return nil, fmt.Errorf("data set is nil")
	}

	return e.embedMDConverter.ConvertDataSet(&models.DataSetImpl{
		ID:               dataSet.GetID(),
		TypeID:           dataSet.GetTypeID(),
		Attributes:       dataSet.GetAttributes(),
		Properties:       dataSet.GetProperties(),
		CustomProperties: dataSet.GetCustomProperties(),
	})
}

func (e *EmbedMDMapper) MapToParameter(parameter models.Parameter) (*openapi.Parameter, error) {
	if parameter == nil {
		return nil, fmt.Errorf("parameter is nil")
	}

	return e.embedMDConverter.ConvertParameter(&models.ParameterImpl{
		ID:               parameter.GetID(),
		TypeID:           parameter.GetTypeID(),
		Attributes:       parameter.GetAttributes(),
		Properties:       parameter.GetProperties(),
		CustomProperties: parameter.GetCustomProperties(),
	})
}

func (e *EmbedMDMapper) MapToArtifact(artifact models.Artifact) (*openapi.Artifact, error) {
	if artifact.ModelArtifact != nil {
		modelArtifact, err := e.MapToModelArtifact(*artifact.ModelArtifact)
		if err != nil {
			return nil, err
		}
		return &openapi.Artifact{ModelArtifact: modelArtifact}, nil
	} else if artifact.DocArtifact != nil {
		docArtifact, err := e.MapToDocArtifact(*artifact.DocArtifact)
		if err != nil {
			return nil, err
		}
		return &openapi.Artifact{DocArtifact: docArtifact}, nil
	} else if artifact.DataSet != nil {
		dataSet, err := e.MapToDataSet(*artifact.DataSet)
		if err != nil {
			return nil, err
		}
		return &openapi.Artifact{DataSet: dataSet}, nil
	} else if artifact.Metric != nil {
		metric, err := e.MapToMetricFromMetric(*artifact.Metric)
		if err != nil {
			return nil, err
		}
		return &openapi.Artifact{Metric: metric}, nil
	} else if artifact.Parameter != nil {
		parameter, err := e.MapToParameter(*artifact.Parameter)
		if err != nil {
			return nil, err
		}
		return &openapi.Artifact{Parameter: parameter}, nil
	}

	return nil, fmt.Errorf("unknown artifact type")
}
