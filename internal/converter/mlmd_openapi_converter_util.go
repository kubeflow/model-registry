package converter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

func NewMetadataStringValue(value string) *openapi.MetadataStringValue {
	result := openapi.NewMetadataStringValueWithDefaults()
	result.StringValue = value
	return result
}

func NewMetadataBoolValue(value bool) *openapi.MetadataBoolValue {
	result := openapi.NewMetadataBoolValueWithDefaults()
	result.BoolValue = value
	return result
}

func NewMetadataDoubleValue(value float64) *openapi.MetadataDoubleValue {
	result := openapi.NewMetadataDoubleValueWithDefaults()
	result.DoubleValue = value
	return result
}

func NewMetadataIntValue(value string) *openapi.MetadataIntValue {
	result := openapi.NewMetadataIntValueWithDefaults()
	result.IntValue = value
	return result
}

func NewMetadataStructValue(value string) *openapi.MetadataStructValue {
	result := openapi.NewMetadataStructValueWithDefaults()
	result.StructValue = value
	return result
}

func NewMetadataProtoValue(typeDef string, value string) *openapi.MetadataProtoValue {
	result := openapi.NewMetadataProtoValueWithDefaults()
	result.Type = typeDef
	result.ProtoValue = value
	return result
}

// MapMLMDCustomProperties maps MLMD custom properties model to OpenAPI one
func MapMLMDCustomProperties(source map[string]*proto.Value) (map[string]openapi.MetadataValue, error) {
	data := make(map[string]openapi.MetadataValue)

	for key, v := range source {
		// data[key] = v.Value
		customValue := openapi.MetadataValue{}

		switch typedValue := v.Value.(type) {
		case *proto.Value_BoolValue:
			customValue.MetadataBoolValue = NewMetadataBoolValue(typedValue.BoolValue)
		case *proto.Value_IntValue:
			customValue.MetadataIntValue = NewMetadataIntValue(strconv.FormatInt(typedValue.IntValue, 10))
		case *proto.Value_DoubleValue:
			customValue.MetadataDoubleValue = NewMetadataDoubleValue(typedValue.DoubleValue)
		case *proto.Value_StringValue:
			customValue.MetadataStringValue = NewMetadataStringValue(typedValue.StringValue)
		case *proto.Value_StructValue:
			sv := typedValue.StructValue
			asMap := sv.AsMap()
			asJSON, err := json.Marshal(asMap)
			if err != nil {
				return nil, err
			}
			b64 := base64.StdEncoding.EncodeToString(asJSON)
			customValue.MetadataStructValue = NewMetadataStructValue(b64)
		default:
			return nil, fmt.Errorf("%w: metadataType not found for %s: %v", api.ErrBadRequest, key, v)
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
	// cat remaining parts of the exploded string
	joined := strings.Join(exploded[1:], ":")
	return &joined
}

// MapName derive the entity name from the mlmd fullname
// for owned entity such as ModelVersion
func MapName(source *string) string {
	if source == nil {
		return ""
	}

	return *MapNameFromOwned(source)
}

// REGISTERED MODEL

func MapRegisteredModelReadme(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "readme")
}

func MapRegisteredModelMaturity(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "maturity")
}

func MapRegisteredModelLanguage(properties map[string]*proto.Value) []string {
	return MapStringSliceProperty(properties, "language")
}

func MapRegisteredModelTasks(properties map[string]*proto.Value) []string {
	return MapStringSliceProperty(properties, "tasks")
}

func MapRegisteredModelProvider(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "provider")
}

func MapRegisteredModelLogo(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "logo")
}

func MapRegisteredModelLicense(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "license")
}

func MapRegisteredModelLicenseLink(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "license_link")
}

func MapRegisteredModelLibraryName(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "library_name")
}

// MODEL VERSION

func MapPropertyAuthor(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "author")
}

func MapRegisteredModelIdFromOwned(source *string) (string, error) {
	if source == nil {
		return "", nil
	}

	exploded := strings.Split(*source, ":")
	if len(exploded) < 2 {
		return "", fmt.Errorf("wrong owned format")
	}
	return exploded[0], nil
}

// ARTIFACT

func MapArtifactType(source *proto.Artifact) (*string, error) {
	if source.Type == nil {
		return nil, fmt.Errorf("artifact type is nil")
	}
	switch *source.Type {
	case defaults.ModelArtifactTypeName:
		return of("model-artifact"), nil
	case defaults.DocArtifactTypeName:
		return of("doc-artifact"), nil
	default:
		return nil, fmt.Errorf("invalid artifact type found: %v", source.Type)
	}
}

// MODEL ARTIFACT

func MapRegisteredModelState(properties map[string]*proto.Value) *openapi.RegisteredModelState {
	state, ok := properties["state"]
	if !ok {
		return nil
	}
	str := state.GetStringValue()
	return (*openapi.RegisteredModelState)(&str)
}

func MapModelVersionState(properties map[string]*proto.Value) *openapi.ModelVersionState {
	state, ok := properties["state"]
	if !ok {
		return nil
	}
	str := state.GetStringValue()
	return (*openapi.ModelVersionState)(&str)
}

func MapInferenceServiceDesiredState(properties map[string]*proto.Value) *openapi.InferenceServiceState {
	state, ok := properties["desired_state"]
	if !ok {
		return nil
	}
	str := state.GetStringValue()
	return (*openapi.InferenceServiceState)(&str)
}

func MapMLMDArtifactState(source *proto.Artifact_State) (st *openapi.ArtifactState) {
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

func MapStringSliceProperty(properties map[string]*proto.Value, key string) []string {
	structValue := properties[key].GetStructValue()
	if structValue == nil {
		return nil
	}

	listValue := structValue.Fields[key].GetListValue()
	if listValue == nil {
		return nil
	}

	anySlice := listValue.AsSlice()
	strSlice := make([]string, 0, len(anySlice))
	for _, v := range anySlice {
		if s, ok := v.(string); ok {
			strSlice = append(strSlice, s)
		}
	}

	return strSlice
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

func MapOwner(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "owner")
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

func MapModelArtifactModelSourceKind(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "model_source_kind")
}

func MapModelArtifactModelSourceClass(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "model_source_class")
}

func MapModelArtifactModelSourceGroup(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "model_source_group")
}

func MapModelArtifactModelSourceId(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "model_source_id")
}

func MapModelArtifactModelSourceName(properties map[string]*proto.Value) *string {
	return MapStringProperty(properties, "model_source_name")
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
