package core

import (
	"fmt"

	"github.com/opendatahub-io/model-registry/internal/converter"
	"github.com/opendatahub-io/model-registry/internal/converter/generated"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
)

type Mapper struct {
	OpenAPIConverter      converter.OpenAPIToMLMDConverter
	MLMDConverter         converter.MLMDToOpenAPIConverter
	RegisteredModelTypeId int64
	ModelVersionTypeId    int64
	ModelArtifactTypeId   int64
}

func NewMapper(registeredModelTypeId int64, modelVersionTypeId int64, modelArtifactTypeId int64) *Mapper {
	return &Mapper{
		OpenAPIConverter:      &generated.OpenAPIToMLMDConverterImpl{},
		MLMDConverter:         &generated.MLMDToOpenAPIConverterImpl{},
		RegisteredModelTypeId: registeredModelTypeId,
		ModelVersionTypeId:    modelVersionTypeId,
		ModelArtifactTypeId:   modelArtifactTypeId,
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

// Utilities for MLMD --> OpenAPI mapping, make use of generated Converters

func (m *Mapper) MapToRegisteredModel(ctx *proto.Context) (*openapi.RegisteredModel, error) {
	if ctx.GetTypeId() != m.RegisteredModelTypeId {
		return nil, fmt.Errorf("invalid TypeId, expected %d but received %d", m.RegisteredModelTypeId, ctx.GetTypeId())
	}

	return m.MLMDConverter.ConvertRegisteredModel(ctx)
}

func (m *Mapper) MapToModelVersion(ctx *proto.Context) (*openapi.ModelVersion, error) {
	if ctx.GetTypeId() != m.ModelVersionTypeId {
		return nil, fmt.Errorf("invalid TypeId, expected %d but received %d", m.ModelVersionTypeId, ctx.GetTypeId())
	}

	return m.MLMDConverter.ConvertModelVersion(ctx)
}

func (m *Mapper) MapToModelArtifact(artifact *proto.Artifact) (*openapi.ModelArtifact, error) {
	if artifact.GetTypeId() != m.ModelArtifactTypeId {
		return nil, fmt.Errorf("invalid TypeId, expected %d but received %d", m.ModelArtifactTypeId, artifact.GetTypeId())
	}

	return m.MLMDConverter.ConvertModelArtifact(artifact)
}
