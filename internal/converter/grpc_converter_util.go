package converter

import (
	"fmt"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/db"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

func ConvertTypeIDToName(id int64) (*string, error) {
	if id == 0 {
		return nil, nil
	}
	if globalDB == nil {
		return nil, fmt.Errorf("converter DB connection has not been initialized")
	}
	var name string
	if err := globalDB.Model(db.Type{}).Select("name").Where("id = ?", id).First(&name).Error; err != nil {
		return nil, fmt.Errorf("error getting type name for type id %d: %w", id, err)
	}
	return &name, nil
}

func ConvertToProtoArtifactProperties(source []db.ArtifactProperty) (map[string]*proto.Value, error) {
	result := make(map[string]*proto.Value)
	for _, prop := range source {
		err := convertToProtoMetadataProperty(&prop, result, false)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func ConvertToProtoArtifactCustomProperties(source []db.ArtifactProperty) (map[string]*proto.Value, error) {
	result := make(map[string]*proto.Value)
	for _, prop := range source {
		err := convertToProtoMetadataProperty(&prop, result, true)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func ConvertToProtoContextProperties(source []db.ContextProperty) (map[string]*proto.Value, error) {
	result := make(map[string]*proto.Value)
	for _, prop := range source {
		err := convertToProtoMetadataProperty(&prop, result, false)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func ConvertToProtoContextCustomProperties(source []db.ContextProperty) (map[string]*proto.Value, error) {
	result := make(map[string]*proto.Value)
	for _, prop := range source {
		err := convertToProtoMetadataProperty(&prop, result, true)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func ConvertToProtoExecutionProperties(source []db.ExecutionProperty) (map[string]*proto.Value, error) {
	result := make(map[string]*proto.Value)
	for _, prop := range source {
		err := convertToProtoMetadataProperty(&prop, result, false)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func ConvertToProtoExecutionCustomProperties(source []db.ExecutionProperty) (map[string]*proto.Value, error) {
	result := make(map[string]*proto.Value)
	for _, prop := range source {
		err := convertToProtoMetadataProperty(&prop, result, true)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func convertToProtoMetadataProperty(p db.MetadataProperty, result map[string]*proto.Value, customProperty bool) error {
	if p.GetIsCustomProperty() == customProperty {
		value, err := convertToProtoValue(p)
		if err != nil {
			return err
		}
		result[p.GetName()] = value
	}
	return nil
}

func convertToProtoValue(prop db.MetadataProperty) (*proto.Value, error) {
	result := proto.Value{}
	if prop.GetIntValue() != nil {
		result.Value = &proto.Value_IntValue{IntValue: *prop.GetIntValue()}
	}
	if prop.GetDoubleValue() != nil {
		result.Value = &proto.Value_DoubleValue{DoubleValue: *prop.GetDoubleValue()}
	}
	if prop.GetStringValue() != nil {
		result.Value = &proto.Value_StringValue{StringValue: *prop.GetStringValue()}
	}
	if prop.GetByteValue() != nil {
		value := structpb.Struct{}
		err := value.UnmarshalJSON(*prop.GetByteValue())
		if err != nil {
			return nil, err
		}
		result.Value = &proto.Value_StructValue{StructValue: &value}
	}
	if prop.GetProtoValue() != nil {
		value := anypb.Any{
			TypeUrl: *prop.GetTypeURL(),
			Value:   *prop.GetProtoValue(),
		}
		result.Value = &proto.Value_ProtoValue{ProtoValue: &value}
	}
	if prop.GetBoolValue() != nil {
		result.Value = &proto.Value_BoolValue{BoolValue: *prop.GetBoolValue()}
	}
	return &result, nil
}

func ConvertProtoArtifactProperties(source *proto.Artifact) ([]db.ArtifactProperty, error) {
	var result []db.ArtifactProperty
	var id int64
	if source.Id != nil {
		id = *source.Id
	}
	err := convertArtifactProtoProperties(source.Properties, id, false, &result)
	if err != nil {
		return nil, err
	}
	err = convertArtifactProtoProperties(source.CustomProperties, id, true, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ConvertProtoContextProperties(source *proto.Context) ([]db.ContextProperty, error) {
	var result []db.ContextProperty
	var id int64
	if source.Id != nil {
		id = *source.Id
	}
	err := convertContextProtoProperties(source.Properties, id, false, &result)
	if err != nil {
		return nil, err
	}
	err = convertContextProtoProperties(source.CustomProperties, id, true, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ConvertProtoExecutionProperties(source *proto.Execution) ([]db.ExecutionProperty, error) {
	var result []db.ExecutionProperty
	var id int64
	if source.Id != nil {
		id = *source.Id
	}
	err := convertExecutionProtoProperties(source.Properties, id, false, &result)
	if err != nil {
		return nil, err
	}
	err = convertExecutionProtoProperties(source.CustomProperties, id, true, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func convertArtifactProtoProperties(source map[string]*proto.Value, id int64, customProperty bool, result *[]db.ArtifactProperty) error {
	for name, prop := range source {
		if prop != nil {
			ap := db.ArtifactProperty{Name: name, ArtifactID: id, IsCustomProperty: customProperty}
			err := convertProtoValue(prop, &ap)
			if err != nil {
				return err
			}
			*result = append(*result, ap)
		}
	}
	return nil
}

func convertContextProtoProperties(source map[string]*proto.Value, id int64, customProperty bool, result *[]db.ContextProperty) error {
	for name, prop := range source {
		if prop != nil {
			ap := db.ContextProperty{Name: name, ContextID: id, IsCustomProperty: customProperty}
			err := convertProtoValue(prop, &ap)
			if err != nil {
				return err
			}
			*result = append(*result, ap)
		}
	}
	return nil
}

func convertExecutionProtoProperties(source map[string]*proto.Value, id int64, customProperty bool, result *[]db.ExecutionProperty) error {
	for name, prop := range source {
		if prop != nil {
			ap := db.ExecutionProperty{Name: name, ExecutionID: id, IsCustomProperty: customProperty}
			err := convertProtoValue(prop, &ap)
			if err != nil {
				return err
			}
			*result = append(*result, ap)
		}
	}
	return nil
}

func convertProtoValue(source *proto.Value, target db.MetadataProperty) error {
	switch p := source.Value.(type) {
	case *proto.Value_IntValue:
		target.SetIntValue(&p.IntValue)
	case *proto.Value_DoubleValue:
		target.SetDoubleValue(&p.DoubleValue)
	case *proto.Value_StringValue:
		target.SetStringValue(&p.StringValue)
	case *proto.Value_StructValue:
		bytes, err := p.StructValue.MarshalJSON()
		if err != nil {
			return fmt.Errorf("error converting struct value: %w", err)
		}
		target.SetByteValue(&bytes)
	case *proto.Value_ProtoValue:
		target.SetTypeURL(&p.ProtoValue.TypeUrl)
		target.SetProtoValue(&p.ProtoValue.Value)
	case *proto.Value_BoolValue:
		target.SetBoolValue(&p.BoolValue)
	}
	return nil
}

func ConvertArtifact_State(source *proto.Artifact_State) *int8 {
	if source == nil {
		return nil
	}
	state := int8((*source).Number())
	return &state
}

func ConvertToArtifact_State(source *int8) (*proto.Artifact_State, error) {
	if source == nil {
		return nil, nil
	}
	if _, ok := proto.Artifact_State_name[int32(*source)]; !ok {
		return nil, fmt.Errorf("invalid artifact state %d", *source)
	}
	state := proto.Artifact_State(*source)
	return &state, nil
}

func ConvertExecution_State(source *proto.Execution_State) *int8 {
	if source == nil {
		return nil
	}
	state := int8((*source).Number())
	return &state
}

func ConvertToExecution_State(source *int8) (*proto.Execution_State, error) {
	if source == nil {
		return nil, nil
	}
	if _, ok := proto.Execution_State_name[int32(*source)]; !ok {
		return nil, fmt.Errorf("invalid execution state %d", *source)
	}
	state := proto.Execution_State(*source)
	return &state, nil
}

func ConvertProtoEventType(source *proto.Event_Type) int8 {
	if source == nil {
		return int8(proto.Event_UNKNOWN)
	}
	return int8((*source).Number())
}

func ConvertToProtoEventType(source int8) (*proto.Event_Type, error) {
	if _, ok := proto.Event_Type_name[int32(source)]; !ok {
		return nil, fmt.Errorf("invalid event type %d", source)
	}
	state := proto.Event_Type(source)
	return &state, nil
}

func ConvertProtoEventPath(source *proto.Event_Path) ([]db.EventPath, error) {
	var result []db.EventPath
	if source != nil {
		for _, step := range source.Steps {
			switch s := step.Value.(type) {
			case *proto.Event_Path_Step_Index:
				i := int(s.Index)
				result = append(result, db.EventPath{
					IsIndexStep: true,
					StepIndex:   &i,
				})
			case *proto.Event_Path_Step_Key:
				result = append(result, db.EventPath{
					IsIndexStep: false,
					StepKey:     &s.Key,
				})
			}
		}
	}
	return result, nil
}

func ConvertToProtoEventPath(source []db.EventPath) (*proto.Event_Path, error) {
	if len(source) == 0 {
		return nil, nil
	}
	var result proto.Event_Path
	for _, step := range source {
		if step.IsIndexStep {
			result.Steps = append(result.Steps, &proto.Event_Path_Step{
				Value: &proto.Event_Path_Step_Index{Index: int64(*step.StepIndex)},
			})
		} else {
			result.Steps = append(result.Steps, &proto.Event_Path_Step{
				Value: &proto.Event_Path_Step_Key{Key: *step.StepKey},
			})
		}
	}
	return &result, nil
}
