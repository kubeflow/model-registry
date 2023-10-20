package mapper

import (
	"fmt"
	"log"
	"strconv"

	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
)

type Mapper struct {
	RegisteredModelTypeId int64
	ModelVersionTypeId    int64
	ModelArtifactTypeId   int64
}

func NewMapper(registeredModelTypeId int64, modelVersionTypeId int64, modelArtifactTypeId int64) *Mapper {
	return &Mapper{
		RegisteredModelTypeId: registeredModelTypeId,
		ModelVersionTypeId:    modelVersionTypeId,
		ModelArtifactTypeId:   modelArtifactTypeId,
	}
}

func IdToInt64(idString string) (*int64, error) {
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	idInt64 := int64(idInt)

	return &idInt64, nil
}

func IdToString(idInt int64) *string {
	idString := strconv.FormatInt(idInt, 10)

	return &idString
}

// Internal Model --> MLMD

// Map generic map into MLMD [custom] properties object
func (m *Mapper) MapToProperties(data map[string]openapi.MetadataValue) (map[string]*proto.Value, error) {
	props := make(map[string]*proto.Value)

	for key, v := range data {
		value := proto.Value{}

		switch {
		// int value
		case v.MetadataIntValue != nil:
			intValue, err := IdToInt64(*v.MetadataIntValue.IntValue)
			if err != nil {
				log.Printf("Skipping mapping for %s:%v", key, v)
				continue
			}
			value.Value = &proto.Value_IntValue{IntValue: *intValue}
		// double value
		case v.MetadataDoubleValue != nil:
			value.Value = &proto.Value_DoubleValue{DoubleValue: *v.MetadataDoubleValue.DoubleValue}
		// string value
		case v.MetadataStringValue != nil:
			value.Value = &proto.Value_StringValue{StringValue: *v.MetadataStringValue.StringValue}
		default:
			log.Printf("Type mapping not found for %s:%v", key, v)
			continue
		}

		props[key] = &value
	}

	return props, nil
}

func (m *Mapper) MapFromRegisteredModel(registeredModel *openapi.RegisteredModel) (*proto.Context, error) {

	var idInt *int64
	if registeredModel.Id != nil {
		var err error
		idInt, err = IdToInt64(*registeredModel.Id)
		if err != nil {
			return nil, err
		}
	}

	customProps := make(map[string]*proto.Value)
	if registeredModel.CustomProperties != nil {
		customProps, _ = m.MapToProperties(*registeredModel.CustomProperties)
	}

	return &proto.Context{
		Name:             registeredModel.Name,
		TypeId:           &m.RegisteredModelTypeId,
		ExternalId:       registeredModel.ExternalID,
		Id:               idInt,
		CustomProperties: customProps,
	}, nil
}

func (m *Mapper) MapFromModelVersion(modelVersion *openapi.ModelVersion, registeredModelId int64, registeredModelName *string) (*proto.Context, error) {
	fullName := fmt.Sprintf("%d:%s", registeredModelId, *modelVersion.Name)
	customProps := make(map[string]*proto.Value)
	if modelVersion.CustomProperties != nil {
		customProps, _ = m.MapToProperties(*modelVersion.CustomProperties)
	}
	ctx := &proto.Context{
		Name:   &fullName,
		TypeId: &m.ModelVersionTypeId,
		Properties: map[string]*proto.Value{
			"model_name": {
				Value: &proto.Value_StringValue{
					StringValue: *registeredModelName,
				},
			},
		},
		CustomProperties: customProps,
	}
	if modelVersion.Name != nil {
		ctx.Properties["version"] = &proto.Value{
			Value: &proto.Value_StringValue{
				StringValue: *modelVersion.Name,
			},
		}
	}
	// TODO: missing explicit property in openapi
	// if modelVersion.Author != nil {
	// 	ctx.Properties["author"] = &proto.Value{
	// 		Value: &proto.Value_StringValue{
	// 			StringValue: *modelVersion.Author,
	// 		},
	// 	}
	// }

	return ctx, nil
}

func (m *Mapper) MapFromModelArtifact(modelArtifact openapi.ModelArtifact) *proto.Artifact {
	return &proto.Artifact{
		TypeId: &m.ModelArtifactTypeId,
		// TODO: we should use concatenation between uuid + name
		Name: modelArtifact.Name,
		Uri:  modelArtifact.Uri,
	}
}

func (m *Mapper) MapFromModelArtifacts(modelArtifacts *[]openapi.ModelArtifact) ([]*proto.Artifact, error) {
	artifacts := []*proto.Artifact{}
	if modelArtifacts == nil {
		return artifacts, nil
	}
	for _, a := range *modelArtifacts {
		artifacts = append(artifacts, m.MapFromModelArtifact(a))
	}
	return artifacts, nil
}

//  MLMD --> Internal Model

// Maps MLMD properties into a generic <string, any> map
func (m *Mapper) MapFromProperties(props map[string]*proto.Value) (map[string]openapi.MetadataValue, error) {
	data := make(map[string]openapi.MetadataValue)

	for key, v := range props {
		// data[key] = v.Value
		customValue := openapi.MetadataValue{}

		switch typedValue := v.Value.(type) {
		case *proto.Value_IntValue:
			customValue.MetadataIntValue = &openapi.MetadataIntValue{
				IntValue: IdToString(typedValue.IntValue),
			}
		case *proto.Value_DoubleValue:
			customValue.MetadataDoubleValue = &openapi.MetadataDoubleValue{
				DoubleValue: &typedValue.DoubleValue,
			}
		case *proto.Value_StringValue:
			customValue.MetadataStringValue = &openapi.MetadataStringValue{
				StringValue: &typedValue.StringValue,
			}
		default:
			log.Printf("Type mapping not found for %s:%v", key, v)
			continue
		}

		data[key] = customValue
	}

	return data, nil
}

func (m *Mapper) MapToRegisteredModel(ctx *proto.Context) (*openapi.RegisteredModel, error) {
	if ctx.GetTypeId() != m.RegisteredModelTypeId {
		return nil, fmt.Errorf("invalid TypeId, exptected %d but received %d", m.RegisteredModelTypeId, ctx.GetTypeId())
	}

	_, err := m.MapFromProperties(ctx.CustomProperties)
	if err != nil {
		return nil, err
	}

	idString := strconv.FormatInt(*ctx.Id, 10)

	model := &openapi.RegisteredModel{
		Id:         &idString,
		Name:       ctx.Name,
		ExternalID: ctx.ExternalId,
	}

	return model, nil
}

func (m *Mapper) MapToModelVersion(ctx *proto.Context) (*openapi.ModelVersion, error) {
	if ctx.GetTypeId() != m.ModelVersionTypeId {
		return nil, fmt.Errorf("invalid TypeId, exptected %d but received %d", m.ModelVersionTypeId, ctx.GetTypeId())
	}

	metadata, err := m.MapFromProperties(ctx.CustomProperties)
	if err != nil {
		return nil, err
	}

	// modelName := ctx.GetProperties()["model_name"].GetStringValue()
	// version := ctx.GetProperties()["version"].GetStringValue()
	// author := ctx.GetProperties()["author"].GetStringValue()

	idString := strconv.FormatInt(*ctx.Id, 10)

	modelVersion := &openapi.ModelVersion{
		// ModelName: &modelName,
		Id:   &idString,
		Name: ctx.Name,
		// Author:   &author,
		CustomProperties: &metadata,
	}

	return modelVersion, nil
}

func (m *Mapper) MapToModelArtifact(artifact *proto.Artifact) (*openapi.ModelArtifact, error) {
	if artifact.GetTypeId() != m.ModelArtifactTypeId {
		return nil, fmt.Errorf("invalid TypeId, exptected %d but received %d", m.ModelArtifactTypeId, artifact.GetTypeId())
	}

	_, err := m.MapFromProperties(artifact.CustomProperties)
	if err != nil {
		return nil, err
	}

	_, err = m.MapFromProperties(artifact.Properties)
	if err != nil {
		return nil, err
	}

	modelArtifact := &openapi.ModelArtifact{
		Uri:  artifact.Uri,
		Name: artifact.Name,
	}

	return modelArtifact, nil
}
