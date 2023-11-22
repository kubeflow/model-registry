package mapper

import (
	"fmt"

	"github.com/opendatahub-io/model-registry/internal/converter"
	"github.com/opendatahub-io/model-registry/internal/converter/generated"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/pkg/openapi"
)

type Mapper struct {
	OpenAPIConverter         converter.OpenAPIToMLMDConverter
	MLMDConverter            converter.MLMDToOpenAPIConverter
	RegisteredModelTypeId    int64
	ModelVersionTypeId       int64
	ModelArtifactTypeId      int64
	ServingEnvironmentTypeId int64
	InferenceServiceTypeId   int64
	ServeModelTypeId         int64
}

func NewMapper(registeredModelTypeId int64, modelVersionTypeId int64, modelArtifactTypeId int64, servingEnvironmentTypeId int64, inferenceServiceTypeId int64, serveModelTypeId int64) *Mapper {
	return &Mapper{
		OpenAPIConverter:         &generated.OpenAPIToMLMDConverterImpl{},
		MLMDConverter:            &generated.MLMDToOpenAPIConverterImpl{},
		RegisteredModelTypeId:    registeredModelTypeId,
		ModelVersionTypeId:       modelVersionTypeId,
		ModelArtifactTypeId:      modelArtifactTypeId,
		ServingEnvironmentTypeId: servingEnvironmentTypeId,
		InferenceServiceTypeId:   inferenceServiceTypeId,
		ServeModelTypeId:         serveModelTypeId,
	}
}

// Utilities for OpenAPI --> MLMD mapping, make use of generated Converters

func (m *Mapper) MapFromRegisteredModel(registeredModel *openapi.RegisteredModel) (*proto.Context, error) {
	ctx, err := m.OpenAPIConverter.ConvertRegisteredModel(&converter.OpenAPIModelWrapper[openapi.RegisteredModel]{
		TypeId: m.RegisteredModelTypeId,
		Model:  registeredModel,
	})
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (m *Mapper) MapFromModelVersion(modelVersion *openapi.ModelVersion, registeredModelId string, registeredModelName *string) (*proto.Context, error) {
	ctx, err := m.OpenAPIConverter.ConvertModelVersion(&converter.OpenAPIModelWrapper[openapi.ModelVersion]{
		TypeId:           m.ModelVersionTypeId,
		Model:            modelVersion,
		ParentResourceId: &registeredModelId,
		ModelName:        registeredModelName,
	})
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (m *Mapper) MapFromModelArtifact(modelArtifact *openapi.ModelArtifact, modelVersionId *string) (*proto.Artifact, error) {

	artifact, err := m.OpenAPIConverter.ConvertModelArtifact(&converter.OpenAPIModelWrapper[openapi.ModelArtifact]{
		TypeId:           m.ModelArtifactTypeId,
		Model:            modelArtifact,
		ParentResourceId: modelVersionId,
	})
	if err != nil {
		return nil, err
	}

	return artifact, nil
}

func (m *Mapper) MapFromModelArtifacts(modelArtifacts *[]openapi.ModelArtifact, modelVersionId *string) ([]*proto.Artifact, error) {
	artifacts := []*proto.Artifact{}
	if modelArtifacts == nil {
		return artifacts, nil
	}
	for _, a := range *modelArtifacts {
		mapped, err := m.MapFromModelArtifact(&a, modelVersionId)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, mapped)
	}
	return artifacts, nil
}

func (m *Mapper) MapFromServingEnvironment(servingEnvironment *openapi.ServingEnvironment) (*proto.Context, error) {
	ctx, err := m.OpenAPIConverter.ConvertServingEnvironment(&converter.OpenAPIModelWrapper[openapi.ServingEnvironment]{
		TypeId: m.ServingEnvironmentTypeId,
		Model:  servingEnvironment,
	})
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (m *Mapper) MapFromInferenceService(inferenceService *openapi.InferenceService, servingEnvironmentId string) (*proto.Context, error) {
	ctx, err := m.OpenAPIConverter.ConvertInferenceService(&converter.OpenAPIModelWrapper[openapi.InferenceService]{
		TypeId:           m.InferenceServiceTypeId,
		Model:            inferenceService,
		ParentResourceId: &servingEnvironmentId,
	})
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (m *Mapper) MapFromServeModel(serveModel *openapi.ServeModel, inferenceServiceId string) (*proto.Execution, error) {
	ctx, err := m.OpenAPIConverter.ConvertServeModel(&converter.OpenAPIModelWrapper[openapi.ServeModel]{
		TypeId:           m.ServeModelTypeId,
		Model:            serveModel,
		ParentResourceId: &inferenceServiceId,
	})
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

// Utilities for MLMD --> OpenAPI mapping, make use of generated Converters

func (m *Mapper) MapToRegisteredModel(ctx *proto.Context) (*openapi.RegisteredModel, error) {
	return mapTo(ctx, m.RegisteredModelTypeId, m.MLMDConverter.ConvertRegisteredModel)
}

func (m *Mapper) MapToModelVersion(ctx *proto.Context) (*openapi.ModelVersion, error) {
	return mapTo(ctx, m.ModelVersionTypeId, m.MLMDConverter.ConvertModelVersion)
}

func (m *Mapper) MapToModelArtifact(art *proto.Artifact) (*openapi.ModelArtifact, error) {
	return mapTo(art, m.ModelArtifactTypeId, m.MLMDConverter.ConvertModelArtifact)
}

func (m *Mapper) MapToServingEnvironment(ctx *proto.Context) (*openapi.ServingEnvironment, error) {
	return mapTo(ctx, m.ServingEnvironmentTypeId, m.MLMDConverter.ConvertServingEnvironment)
}

func (m *Mapper) MapToInferenceService(ctx *proto.Context) (*openapi.InferenceService, error) {
	return mapTo(ctx, m.InferenceServiceTypeId, m.MLMDConverter.ConvertInferenceService)
}

func (m *Mapper) MapToServeModel(ex *proto.Execution) (*openapi.ServeModel, error) {
	return mapTo(ex, m.ServeModelTypeId, m.MLMDConverter.ConvertServeModel)
}

type getTypeIder interface {
	GetTypeId() int64
}

func mapTo[S getTypeIder, T any](s S, id int64, convFn func(S) (*T, error)) (*T, error) {
	if s.GetTypeId() != id {
		return nil, fmt.Errorf("invalid TypeId, expected %d but received %d", id, s.GetTypeId())
	}
	return convFn(s)
}
