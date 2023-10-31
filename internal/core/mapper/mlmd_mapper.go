package mapper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
	"google.golang.org/protobuf/types/known/structpb"
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
		// bool value
		case v.MetadataBoolValue != nil:
			value.Value = &proto.Value_BoolValue{BoolValue: *v.MetadataBoolValue.BoolValue}
		// int value
		case v.MetadataIntValue != nil:
			intValue, err := IdToInt64(*v.MetadataIntValue.IntValue)
			if err != nil {
				return nil, fmt.Errorf("unable to decode as int64 %w for key %s", err, key)
			}
			value.Value = &proto.Value_IntValue{IntValue: *intValue}
		// double value
		case v.MetadataDoubleValue != nil:
			value.Value = &proto.Value_DoubleValue{DoubleValue: *v.MetadataDoubleValue.DoubleValue}
		// string value
		case v.MetadataStringValue != nil:
			value.Value = &proto.Value_StringValue{StringValue: *v.MetadataStringValue.StringValue}
		// struct value
		case v.MetadataStructValue != nil:
			data, err := base64.StdEncoding.DecodeString(*v.MetadataStructValue.StructValue)
			if err != nil {
				return nil, fmt.Errorf("unable to decode %w for key %s", err, key)
			}
			var asMap map[string]interface{}
			err = json.Unmarshal(data, &asMap)
			if err != nil {
				return nil, fmt.Errorf("unable to decode %w for key %s", err, key)
			}
			asStruct, err := structpb.NewStruct(asMap)
			if err != nil {
				return nil, fmt.Errorf("unable to decode %w for key %s", err, key)
			}
			value.Value = &proto.Value_StructValue{
				StructValue: asStruct,
			}
		default:
			return nil, fmt.Errorf("type mapping not found for %s:%v", key, v)
		}

		props[key] = &value
	}

	return props, nil
}

func (m *Mapper) MapToArtifactState(oapiState *openapi.ArtifactState) *proto.Artifact_State {
	if oapiState == nil {
		return nil
	}

	state := (proto.Artifact_State)(proto.Artifact_State_value[string(*oapiState)])
	return &state
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
		Id:               idInt,
		TypeId:           &m.RegisteredModelTypeId,
		Name:             registeredModel.Name,
		ExternalId:       registeredModel.ExternalID,
		CustomProperties: customProps,
	}, nil
}

func (m *Mapper) MapFromModelVersion(modelVersion *openapi.ModelVersion, registeredModelId int64, registeredModelName *string) (*proto.Context, error) {
	fullName := PrefixWhenOwned(&registeredModelId, *modelVersion.Name)
	customProps := make(map[string]*proto.Value)
	if modelVersion.CustomProperties != nil {
		customProps, _ = m.MapToProperties(*modelVersion.CustomProperties)
	}

	var idAsInt *int64
	if modelVersion.Id != nil {
		var err error
		idAsInt, err = IdToInt64(*modelVersion.Id)
		if err != nil {
			return nil, err
		}
	}
	ctx := &proto.Context{
		Id:         idAsInt,
		Name:       &fullName,
		TypeId:     &m.ModelVersionTypeId,
		ExternalId: modelVersion.ExternalID,
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

func (m *Mapper) MapFromModelArtifact(modelArtifact openapi.ModelArtifact, modelVersionId *int64) *proto.Artifact {
	// openapi.Artifact is defined with optional name, so build arbitrary name for this artifact if missing
	var artifactName string
	if modelArtifact.Name != nil {
		artifactName = *modelArtifact.Name
	} else {
		artifactName = uuid.New().String()
	}
	// build fullName for mlmd storage
	fullName := PrefixWhenOwned(modelVersionId, artifactName)

	customProps := make(map[string]*proto.Value)
	if modelArtifact.CustomProperties != nil {
		customProps, _ = m.MapToProperties(*modelArtifact.CustomProperties)
	}

	var idAsInt *int64
	if modelArtifact.Id != nil {
		idAsInt, _ = IdToInt64(*modelArtifact.Id)
	}

	return &proto.Artifact{
		Id:               idAsInt,
		TypeId:           &m.ModelArtifactTypeId,
		Name:             &fullName,
		Uri:              modelArtifact.Uri,
		ExternalId:       modelArtifact.ExternalID,
		State:            m.MapToArtifactState(modelArtifact.State),
		CustomProperties: customProps,
	}
}

func (m *Mapper) MapFromModelArtifacts(modelArtifacts *[]openapi.ModelArtifact, modelVersionId *int64) ([]*proto.Artifact, error) {
	artifacts := []*proto.Artifact{}
	if modelArtifacts == nil {
		return artifacts, nil
	}
	for _, a := range *modelArtifacts {
		artifacts = append(artifacts, m.MapFromModelArtifact(a, modelVersionId))
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
		case *proto.Value_BoolValue:
			customValue.MetadataBoolValue = &openapi.MetadataBoolValue{
				BoolValue: &typedValue.BoolValue,
			}
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
		case *proto.Value_StructValue:
			sv := typedValue.StructValue
			asMap := sv.AsMap()
			asJSON, err := json.Marshal(asMap)
			if err != nil {
				return nil, err
			}
			b64 := base64.StdEncoding.EncodeToString(asJSON)
			customValue.MetadataStructValue = &openapi.MetadataStructValue{
				StructValue: &b64,
			}
		default:
			return nil, fmt.Errorf("type mapping not found for %s:%v", key, v)
		}

		data[key] = customValue
	}

	return data, nil
}

func (m *Mapper) MapFromArtifactState(mlmdState *proto.Artifact_State) *openapi.ArtifactState {
	if mlmdState == nil {
		return nil
	}

	state := mlmdState.String()
	return (*openapi.ArtifactState)(&state)
}

func (m *Mapper) MapToRegisteredModel(ctx *proto.Context) (*openapi.RegisteredModel, error) {
	if ctx.GetTypeId() != m.RegisteredModelTypeId {
		return nil, fmt.Errorf("invalid TypeId, exptected %d but received %d", m.RegisteredModelTypeId, ctx.GetTypeId())
	}

	customProps, err := m.MapFromProperties(ctx.CustomProperties)
	if err != nil {
		return nil, err
	}

	idString := strconv.FormatInt(*ctx.Id, 10)

	model := &openapi.RegisteredModel{
		Id:               &idString,
		Name:             ctx.Name,
		ExternalID:       ctx.ExternalId,
		CustomProperties: &customProps,
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

	name := NameFromOwned(*ctx.Name)
	modelVersion := &openapi.ModelVersion{
		// ModelName: &modelName,
		Id:         &idString,
		Name:       &name,
		ExternalID: ctx.ExternalId,
		// Author:   &author,
		CustomProperties: &metadata,
	}

	return modelVersion, nil
}

func (m *Mapper) MapToModelArtifact(artifact *proto.Artifact) (*openapi.ModelArtifact, error) {
	if artifact.GetTypeId() != m.ModelArtifactTypeId {
		return nil, fmt.Errorf("invalid TypeId, exptected %d but received %d", m.ModelArtifactTypeId, artifact.GetTypeId())
	}

	customProps, err := m.MapFromProperties(artifact.CustomProperties)
	if err != nil {
		return nil, err
	}

	_, err = m.MapFromProperties(artifact.Properties)
	if err != nil {
		return nil, err
	}

	name := NameFromOwned(*artifact.Name)
	modelArtifact := &openapi.ModelArtifact{
		Id:               IdToString(*artifact.Id),
		Uri:              artifact.Uri,
		Name:             &name,
		ExternalID:       artifact.ExternalId,
		State:            m.MapFromArtifactState(artifact.State),
		CustomProperties: &customProps,
	}

	return modelArtifact, nil
}

// For owned entity such as ModelVersion
// for potentially owned entity such as ModelArtifact
// compose the mlmd fullname by using ownerId as prefix
func PrefixWhenOwned(ownerId *int64, entityName string) string {
	if ownerId != nil {
		return fmt.Sprintf("%d:%s", *ownerId, entityName)
	}
	uuidPrefix := uuid.New().String()
	return fmt.Sprintf("%s:%s", uuidPrefix, entityName)
}

// For owned entity such as ModelVersion
// for potentially owned entity such as ModelArtifact
// derive the entity name from the mlmd fullname
func NameFromOwned(fullName string) string {
	return strings.Split(fullName, ":")[1]
}
