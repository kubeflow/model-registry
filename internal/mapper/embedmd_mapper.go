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
	typesMap map[string]int64
}

func NewEmbedMDMapper(typesMap map[string]int64) *EmbedMDMapper {
	return &EmbedMDMapper{
		openAPIConverter:     &generated.OpenAPIToEmbedMDConverterImpl{},
		embedMDConverter:     &generated.EmbedMDToOpenAPIConverterImpl{},
		OpenAPIConverterImpl: &generated.OpenAPIConverterImpl{},
		typesMap:             typesMap,
	}
}

// Utilities for OpenAPI --> EmbedMD mapping, make use of generated Converters

func (e *EmbedMDMapper) MapFromRegisteredModel(registeredModel *openapi.RegisteredModel) (models.RegisteredModel, error) {
	return e.openAPIConverter.ConvertRegisteredModel(&converter.OpenAPIModelWrapper[openapi.RegisteredModel]{
		TypeId: e.typesMap[defaults.RegisteredModelTypeName],
		Model:  registeredModel,
	})
}

func (e *EmbedMDMapper) MapFromModelVersion(modelVersion *openapi.ModelVersion) (models.ModelVersion, error) {
	return e.openAPIConverter.ConvertModelVersion(&converter.OpenAPIModelWrapper[openapi.ModelVersion]{
		TypeId: e.typesMap[defaults.ModelVersionTypeName],
		Model:  modelVersion,
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

func (e *EmbedMDMapper) MapFromModelArtifact(modelArtifact *openapi.ModelArtifact) (models.ModelArtifact, error) {
	return e.openAPIConverter.ConvertModelArtifact(&converter.OpenAPIModelWrapper[openapi.ModelArtifact]{
		TypeId: e.typesMap[defaults.ModelArtifactTypeName],
		Model:  modelArtifact,
	})
}

func (e *EmbedMDMapper) MapFromDocArtifact(docArtifact *openapi.DocArtifact) (models.DocArtifact, error) {
	return e.openAPIConverter.ConvertDocArtifact(&converter.OpenAPIModelWrapper[openapi.DocArtifact]{
		TypeId: e.typesMap[defaults.DocArtifactTypeName],
		Model:  docArtifact,
	})
}

func (e *EmbedMDMapper) MapFromServeModel(serveModel *openapi.ServeModel) (models.ServeModel, error) {
	return e.openAPIConverter.ConvertServeModel(&converter.OpenAPIModelWrapper[openapi.ServeModel]{
		TypeId: e.typesMap[defaults.ServeModelTypeName],
		Model:  serveModel,
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
