package converter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/pkg/openapi"
)

// MapMLMDCustomProperties maps MLMD custom properties model to OpenAPI one
func MapMLMDCustomProperties(source map[string]*proto.Value) (map[string]openapi.MetadataValue, error) {
	data := make(map[string]openapi.MetadataValue)

	for key, v := range source {
		// data[key] = v.Value
		customValue := openapi.MetadataValue{}

		switch typedValue := v.Value.(type) {
		case *proto.Value_BoolValue:
			customValue.MetadataBoolValue = &openapi.MetadataBoolValue{
				BoolValue: &typedValue.BoolValue,
			}
		case *proto.Value_IntValue:
			customValue.MetadataIntValue = &openapi.MetadataIntValue{
				IntValue: Int64ToString(&typedValue.IntValue),
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

// MapNameFromOwned derive the entity name from the mlmd fullname
// for owned entity such as ModelVersion
// for potentially owned entity such as ModelArtifact
func MapNameFromOwned(source *string) *string {
	if source == nil {
		return nil
	}

	exploded := strings.Split(*source, ":")
	if len(exploded) == 1 {
		return source
	}
	return &exploded[1]
}

// REGISTERED MODEL

// MODEL VERSION

// MODEL ARTIFACT

func MapArtifactType(source *proto.Artifact) (string, error) {
	if source.Type != nil && *source.Type == ModelArtifactTypeName {
		return "model-artifact", nil
	}
	return "", fmt.Errorf("invalid artifact type found: %v", source.Type)
}

func MapRegisteredModelState(source *proto.Context) *openapi.RegisteredModelState {
	if source == nil || source.GetProperties() == nil {
		return nil
	}

	state, ok := source.GetProperties()["state"]
	if !ok {
		return nil
	}
	str := state.GetStringValue()
	return (*openapi.RegisteredModelState)(&str)
}

func MapModelVersionState(source *proto.Context) *openapi.ModelVersionState {
	if source == nil || source.GetProperties() == nil {
		return nil
	}

	state, ok := source.GetProperties()["state"]
	if !ok {
		return nil
	}
	str := state.GetStringValue()
	return (*openapi.ModelVersionState)(&str)
}

func MapInferenceServiceState(source *proto.Context) *openapi.InferenceServiceState {
	if source == nil || source.GetProperties() == nil {
		return nil
	}

	state, ok := source.GetProperties()["state"]
	if !ok {
		return nil
	}
	str := state.GetStringValue()
	return (*openapi.InferenceServiceState)(&str)
}

func MapMLMDModelArtifactState(source *proto.Artifact_State) *openapi.ArtifactState {
	if source == nil {
		return nil
	}

	state := source.String()
	return (*openapi.ArtifactState)(&state)
}

// MapStringProperty maps string proto.Value property to specific string field
func MapStringProperty(properties map[string]*proto.Value, key string) *string {
	val, ok := properties[key]
	if ok {
		res := val.GetStringValue()
		if res != "" {
			return &res
		}
	}

	return nil
}

// MapIntProperty maps int proto.Value property to specific string field
func MapIntProperty(properties map[string]*proto.Value, key string) *string {
	val, ok := properties[key]
	if ok {
		res := val.GetIntValue()
		return Int64ToString(&res)
	}

	return nil
}

// MapIntPropertyAsValue maps int proto.Value property to specific string field
func MapIntPropertyAsValue(properties map[string]*proto.Value, key string) string {
	val := MapIntProperty(properties, key)
	if val != nil {
		return *val
	}
	return ""
}

func MapDescription(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "description")
}

func MapModelArtifactFormatName(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "model_format_name")
}

func MapModelArtifactFormatVersion(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "model_format_version")
}

func MapModelArtifactStorageKey(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "storage_key")
}

func MapModelArtifactStoragePath(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "storage_path")
}

func MapModelArtifactServiceAccountName(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "service_account_name")
}

func MapPropertyModelVersionId(properties map[string]*proto.Value) *string {
	return MapIntProperty(properties, "model_version_id")
}

func MapPropertyModelVersionIdAsValue(properties map[string]*proto.Value) string {
	return MapIntPropertyAsValue(properties, "model_version_id")
}

func MapPropertyRegisteredModelId(properties map[string]*proto.Value) string {
	return MapIntPropertyAsValue(properties, "registered_model_id")
}

func MapPropertyServingEnvironmentId(properties map[string]*proto.Value) string {
	return MapIntPropertyAsValue(properties, "serving_environment_id")
}

// INFERENCE SERVICE

func MapPropertyRuntime(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "runtime")
}

// SERVE MODEL

func MapMLMDServeModelLastKnownState(source *proto.Execution_State) *openapi.ExecutionState {
	if source == nil {
		return nil
	}

	state := source.String()
	return (*openapi.ExecutionState)(&state)
}
