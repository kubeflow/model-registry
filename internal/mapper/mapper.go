package mapper

import (
	"fmt"

	"github.com/opendatahub-io/model-registry/internal/constants"
	"github.com/opendatahub-io/model-registry/internal/converter"
	"github.com/opendatahub-io/model-registry/internal/converter/generated"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/pkg/openapi"
)

type Mapper struct {
	OpenAPIConverter converter.OpenAPIToMLMDConverter
	MLMDConverter    converter.MLMDToOpenAPIConverter
	MLMDTypes        map[string]int64
}

func NewMapper(mlmdTypes map[string]int64) *Mapper {
	return &Mapper{
		OpenAPIConverter: &generated.OpenAPIToMLMDConverterImpl{},
		MLMDConverter:    &generated.MLMDToOpenAPIConverterImpl{},
		MLMDTypes:        mlmdTypes,
	}
}

// Utilities for OpenAPI --> MLMD mapping, make use of generated Converters

func (m *Mapper) MapFromRegisteredModel(registeredModel *openapi.RegisteredModel) (*proto.Context, error) {
	ctx, err := m.OpenAPIConverter.ConvertRegisteredModel(&converter.OpenAPIModelWrapper[openapi.RegisteredModel]{
		TypeId: m.MLMDTypes[constants.RegisteredModelTypeName],
		Model:  registeredModel,
	})
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (m *Mapper) MapFromModelVersion(modelVersion *openapi.ModelVersion, registeredModelId string, registeredModelName *string) (*proto.Context, error) {
	ctx, err := m.OpenAPIConverter.ConvertModelVersion(&converter.OpenAPIModelWrapper[openapi.ModelVersion]{
		TypeId:           m.MLMDTypes[constants.ModelVersionTypeName],
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
		TypeId:           m.MLMDTypes[constants.ModelArtifactTypeName],
		Model:            modelArtifact,
		ParentResourceId: modelVersionId,
	})
	if err != nil {
		return nil, err
	}

	return artifact, nil
}

func (m *Mapper) MapFromModelArtifacts(modelArtifacts []openapi.ModelArtifact, modelVersionId *string) ([]*proto.Artifact, error) {
	artifacts := []*proto.Artifact{}
	if modelArtifacts == nil {
		return artifacts, nil
	}
	for _, a := range modelArtifacts {
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
		TypeId: m.MLMDTypes[constants.ServingEnvironmentTypeName],
		Model:  servingEnvironment,
	})
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (m *Mapper) MapFromInferenceService(inferenceService *openapi.InferenceService, servingEnvironmentId string) (*proto.Context, error) {
	ctx, err := m.OpenAPIConverter.ConvertInferenceService(&converter.OpenAPIModelWrapper[openapi.InferenceService]{
		TypeId:           m.MLMDTypes[constants.InferenceServiceTypeName],
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
		TypeId:           m.MLMDTypes[constants.ServeModelTypeName],
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
	return mapTo(ctx, m.MLMDTypes, constants.RegisteredModelTypeName, m.MLMDConverter.ConvertRegisteredModel)
}

func (m *Mapper) MapToModelVersion(ctx *proto.Context) (*openapi.ModelVersion, error) {
	return mapTo(ctx, m.MLMDTypes, constants.ModelVersionTypeName, m.MLMDConverter.ConvertModelVersion)
}

func (m *Mapper) MapToModelArtifact(art *proto.Artifact) (*openapi.ModelArtifact, error) {
	return mapTo(art, m.MLMDTypes, constants.ModelArtifactTypeName, m.MLMDConverter.ConvertModelArtifact)
}

func (m *Mapper) MapToServingEnvironment(ctx *proto.Context) (*openapi.ServingEnvironment, error) {
	return mapTo(ctx, m.MLMDTypes, constants.ServingEnvironmentTypeName, m.MLMDConverter.ConvertServingEnvironment)
}

func (m *Mapper) MapToInferenceService(ctx *proto.Context) (*openapi.InferenceService, error) {
	return mapTo(ctx, m.MLMDTypes, constants.InferenceServiceTypeName, m.MLMDConverter.ConvertInferenceService)
}

func (m *Mapper) MapToServeModel(ex *proto.Execution) (*openapi.ServeModel, error) {
	return mapTo(ex, m.MLMDTypes, constants.ServeModelTypeName, m.MLMDConverter.ConvertServeModel)
}

type getTypeIder interface {
	GetTypeId() int64
	GetType() string
}

func mapTo[S getTypeIder, T any](s S, typesMap map[string]int64, typeName string, convFn func(S) (*T, error)) (*T, error) {
	id, ok := typesMap[typeName]
	if !ok {
		return nil, fmt.Errorf("unknown type name provided: %s", typeName)
	}

	if s.GetTypeId() != id {
		return nil, fmt.Errorf("invalid entity: expected %s but received %s, please check the provided id", typeName, s.GetType())
	}
	return convFn(s)
}
