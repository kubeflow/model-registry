package converter

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func setup(t *testing.T) *assert.Assertions {
	return assert.New(t)
}

func TestStringToInt64(t *testing.T) {
	assertion := setup(t)

	valid := "12345"
	converted, err := StringToInt64(&valid)
	assertion.Nil(err)
	assertion.Equal(int64(12345), *converted)
	assertion.Nil(StringToInt64(nil))
}

func TestStringToInt64InvalidNumber(t *testing.T) {
	assertion := setup(t)

	invalid := "not-a-number"
	converted, err := StringToInt64(&invalid)
	assertion.NotNil(err)
	assertion.Nil(converted)
}

func TestInt64ToString(t *testing.T) {
	assertion := setup(t)

	valid := int64(54321)
	converted := Int64ToString(&valid)
	assertion.Equal("54321", *converted)
	assertion.Nil(Int64ToString(nil))
}

func TestStringToInt32(t *testing.T) {
	assertion := setup(t)

	valid := "12345"
	converted, err := StringToInt32(valid)
	assertion.Nil(err)
	assertion.Equal(int32(12345), converted)
}

func TestStringToInt32InvalidNumber(t *testing.T) {
	assertion := setup(t)

	invalid := "not-a-number"
	_, err := StringToInt32(invalid)
	assertion.NotNil(err)
}

func TestMetadataValueBool(t *testing.T) {
	data := make(map[string]openapi.MetadataValue)
	key := "my bool"
	mdValue := true
	data[key] = openapi.MetadataBoolValueAsMetadataValue(&openapi.MetadataBoolValue{BoolValue: &mdValue})

	roundTripAndAssert(t, data, key)
}

func TestMetadataValueInt(t *testing.T) {
	data := make(map[string]openapi.MetadataValue)
	key := "my int"
	mdValue := "987"
	data[key] = openapi.MetadataIntValueAsMetadataValue(&openapi.MetadataIntValue{IntValue: &mdValue})

	roundTripAndAssert(t, data, key)
}

func TestMetadataValueIntFailure(t *testing.T) {
	data := make(map[string]openapi.MetadataValue)
	key := "my int"
	mdValue := "not a number"
	data[key] = openapi.MetadataIntValueAsMetadataValue(&openapi.MetadataIntValue{IntValue: &mdValue})

	assertion := setup(t)
	asGRPC, err := MapOpenAPICustomProperties(&data)
	if err == nil {
		assertion.Fail("Did not expected a converted value but an error: %v", asGRPC)
	}
}

func TestMetadataValueDouble(t *testing.T) {
	data := make(map[string]openapi.MetadataValue)
	key := "my double"
	mdValue := 3.1415
	data[key] = openapi.MetadataDoubleValueAsMetadataValue(&openapi.MetadataDoubleValue{DoubleValue: &mdValue})

	roundTripAndAssert(t, data, key)
}

func TestMetadataValueString(t *testing.T) {
	data := make(map[string]openapi.MetadataValue)
	key := "my string"
	mdValue := "Hello, World!"
	data[key] = openapi.MetadataStringValueAsMetadataValue(&openapi.MetadataStringValue{StringValue: &mdValue})

	roundTripAndAssert(t, data, key)
}

func TestMetadataValueStruct(t *testing.T) {
	data := make(map[string]openapi.MetadataValue)
	key := "my struct"

	myMap := make(map[string]interface{})
	myMap["name"] = "John Doe"
	myMap["age"] = 47
	asJSON, err := json.Marshal(myMap)
	if err != nil {
		t.Error(err)
	}
	b64 := base64.StdEncoding.EncodeToString(asJSON)
	data[key] = openapi.MetadataStructValueAsMetadataValue(&openapi.MetadataStructValue{StructValue: &b64})

	roundTripAndAssert(t, data, key)
}

func TestMetadataValueProtoUnsupported(t *testing.T) {
	data := make(map[string]openapi.MetadataValue)
	key := "my proto"

	myMap := make(map[string]interface{})
	myMap["name"] = "John Doe"
	myMap["age"] = 47
	asJSON, err := json.Marshal(myMap)
	if err != nil {
		t.Error(err)
	}
	b64 := base64.StdEncoding.EncodeToString(asJSON)
	typeDef := "map[string]openapi.MetadataValue"
	data[key] = openapi.MetadataProtoValueAsMetadataValue(&openapi.MetadataProtoValue{
		Type:       &typeDef,
		ProtoValue: &b64,
	})

	assertion := setup(t)
	asGRPC, err := MapOpenAPICustomProperties(&data)
	if err == nil {
		assertion.Fail("Did not expected a converted value but an error: %v", asGRPC)
	}
}

func roundTripAndAssert(t *testing.T, data map[string]openapi.MetadataValue, key string) {
	assertion := setup(t)

	// first half
	asGRPC, err := MapOpenAPICustomProperties(&data)
	if err != nil {
		t.Error(err)
	}
	assertion.Contains(maps.Keys(asGRPC), key)

	// second half
	unmarshall, err := MapMLMDCustomProperties(asGRPC)
	if err != nil {
		t.Error(err)
	}
	assertion.Equal(data, unmarshall, "result of round-trip shall be equal to original data")
}

func TestPrefixWhenOwned(t *testing.T) {
	assertion := setup(t)

	owner := "owner"
	entity := "name"
	assertion.Equal("owner:name", PrefixWhenOwned(&owner, entity))
}

func TestPrefixWhenOwnedWithoutOwner(t *testing.T) {
	assertion := setup(t)

	entity := "name"
	prefixed := PrefixWhenOwned(nil, entity)
	assertion.Equal(2, len(strings.Split(prefixed, ":")))
	assertion.Equal("name", strings.Split(prefixed, ":")[1])
}

func TestMapRegisteredModelProperties(t *testing.T) {
	assertion := setup(t)

	props, err := MapRegisteredModelProperties(&openapi.RegisteredModel{
		Description: of("super description"),
	})
	assertion.Nil(err)
	assertion.Equal(1, len(props))
	assertion.Equal("super description", props["description"].GetStringValue())

	props, err = MapRegisteredModelProperties(&openapi.RegisteredModel{})
	assertion.Nil(err)
	assertion.Equal(0, len(props))
}

func TestMapRegisteredModelType(t *testing.T) {
	assertion := setup(t)

	typeName := MapRegisteredModelType(&openapi.RegisteredModel{})
	assertion.NotNil(typeName)
	assertion.Equal(RegisteredModelTypeName, *typeName)
}

func TestMapModelVersionProperties(t *testing.T) {
	assertion := setup(t)

	props, err := MapModelVersionProperties(&OpenAPIModelWrapper[openapi.ModelVersion]{
		TypeId:           123,
		ParentResourceId: of("123"),
		ModelName:        of("MyModel"),
		Model: &openapi.ModelVersion{
			Name:        of("v1"),
			Description: of("my model version description"),
			Author:      of("John Doe"),
		},
	})
	assertion.Nil(err)
	assertion.Equal(4, len(props))
	assertion.Equal("my model version description", props["description"].GetStringValue())
	assertion.Equal("v1", props["version"].GetStringValue())
	assertion.Equal("John Doe", props["author"].GetStringValue())
}

func TestMapModelVersionType(t *testing.T) {
	assertion := setup(t)

	typeName := MapModelVersionType(&openapi.ModelVersion{})
	assertion.NotNil(typeName)
	assertion.Equal(ModelVersionTypeName, *typeName)
}

func TestMapModelVersionName(t *testing.T) {
	assertion := setup(t)

	name := MapModelVersionName(&OpenAPIModelWrapper[openapi.ModelVersion]{
		TypeId:           123,
		ParentResourceId: of("123"),
		ModelName:        of("MyModel"),
		Model: &openapi.ModelVersion{
			Name: of("v1"),
		},
	})
	assertion.NotNil(name)
	assertion.Equal("123:v1", *name)
}

func TestMapModelArtifactProperties(t *testing.T) {
	assertion := setup(t)

	props, err := MapModelArtifactProperties(&openapi.ModelArtifact{
		Name:               of("v1"),
		Description:        of("my model art description"),
		ModelFormatName:    of("sklearn"),
		ModelFormatVersion: of("1.0"),
		StorageKey:         of("storage-key"),
		StoragePath:        of("storage-path"),
		ServiceAccountName: of("service-account-name"),
	})
	assertion.Nil(err)
	assertion.Equal(6, len(props))
	assertion.Equal("my model art description", props["description"].GetStringValue())
	assertion.Equal("sklearn", props["model_format_name"].GetStringValue())
	assertion.Equal("1.0", props["model_format_version"].GetStringValue())
	assertion.Equal("storage-key", props["storage_key"].GetStringValue())
	assertion.Equal("storage-path", props["storage_path"].GetStringValue())
	assertion.Equal("service-account-name", props["service_account_name"].GetStringValue())

	props, err = MapModelArtifactProperties(&openapi.ModelArtifact{
		Name: of("v1"),
	})
	assertion.Nil(err)
	assertion.Equal(0, len(props))
}

func TestMapModelArtifactType(t *testing.T) {
	assertion := setup(t)

	typeName := MapModelArtifactType(&openapi.ModelArtifact{})
	assertion.NotNil(typeName)
	assertion.Equal(ModelArtifactTypeName, *typeName)
}

func TestMapModelArtifactName(t *testing.T) {
	assertion := setup(t)

	name := MapModelArtifactName(&OpenAPIModelWrapper[openapi.ModelArtifact]{
		TypeId:           123,
		ParentResourceId: of("parent"),
		Model: &openapi.ModelArtifact{
			Name: of("v1"),
		},
	})
	assertion.NotNil(name)
	assertion.Equal("parent:v1", *name)

	name = MapModelArtifactName(&OpenAPIModelWrapper[openapi.ModelArtifact]{
		TypeId:           123,
		ParentResourceId: of("parent"),
		Model: &openapi.ModelArtifact{
			Name: nil,
		},
	})
	assertion.NotNil(name)
	assertion.Regexp("parent:.*", *name)

	name = MapModelArtifactName(&OpenAPIModelWrapper[openapi.ModelArtifact]{
		TypeId: 123,
		Model: &openapi.ModelArtifact{
			Name: of("v1"),
		},
	})
	assertion.NotNil(name)
	assertion.Regexp(".*:v1", *name)
}

func TestMapOpenAPIModelArtifactState(t *testing.T) {
	assertion := setup(t)

	state, err := MapOpenAPIModelArtifactState(of(openapi.ARTIFACTSTATE_LIVE))
	assertion.Nil(err)
	assertion.NotNil(state)
	assertion.Equal(string(openapi.ARTIFACTSTATE_LIVE), state.String())

	state, err = MapOpenAPIModelArtifactState(nil)
	assertion.Nil(err)
	assertion.Nil(state)
}

func TestMapStringPropertyWithMissingKey(t *testing.T) {
	assertion := setup(t)

	notPresent := MapStringProperty(map[string]*proto.Value{}, "not_present")

	assertion.Nil(notPresent)
}

func TestMapDescription(t *testing.T) {
	assertion := setup(t)

	extracted := MapDescription(map[string]*proto.Value{
		"description": {
			Value: &proto.Value_StringValue{
				StringValue: "my-description",
			},
		},
	})

	assertion.Equal("my-description", *extracted)
}

func TestPropertyRuntime(t *testing.T) {
	assertion := setup(t)

	extracted := MapPropertyRuntime(map[string]*proto.Value{
		"runtime": {
			Value: &proto.Value_StringValue{
				StringValue: "my-runtime",
			},
		},
	})

	assertion.Equal("my-runtime", *extracted)
}

func TestMapModelArtifactFormatName(t *testing.T) {
	assertion := setup(t)

	extracted := MapModelArtifactFormatName(map[string]*proto.Value{
		"model_format_name": {
			Value: &proto.Value_StringValue{
				StringValue: "my-name",
			},
		},
	})

	assertion.Equal("my-name", *extracted)
}

func TestMapModelArtifactFormatVersion(t *testing.T) {
	assertion := setup(t)

	extracted := MapModelArtifactFormatVersion(map[string]*proto.Value{
		"model_format_version": {
			Value: &proto.Value_StringValue{
				StringValue: "my-version",
			},
		},
	})

	assertion.Equal("my-version", *extracted)
}

func TestMapModelArtifactStorageKey(t *testing.T) {
	assertion := setup(t)

	extracted := MapModelArtifactStorageKey(map[string]*proto.Value{
		"storage_key": {
			Value: &proto.Value_StringValue{
				StringValue: "my-key",
			},
		},
	})

	assertion.Equal("my-key", *extracted)
}

func TestMapModelArtifactStoragePath(t *testing.T) {
	assertion := setup(t)

	extracted := MapModelArtifactStoragePath(map[string]*proto.Value{
		"storage_path": {
			Value: &proto.Value_StringValue{
				StringValue: "my-path",
			},
		},
	})

	assertion.Equal("my-path", *extracted)
}

func TestMapModelArtifactServiceAccountName(t *testing.T) {
	assertion := setup(t)

	extracted := MapModelArtifactServiceAccountName(map[string]*proto.Value{
		"service_account_name": {
			Value: &proto.Value_StringValue{
				StringValue: "my-account",
			},
		},
	})

	assertion.Equal("my-account", *extracted)
}

func TestMapPropertyModelVersionId(t *testing.T) {
	assertion := setup(t)

	extracted := MapPropertyModelVersionId(map[string]*proto.Value{
		"model_version_id": {
			Value: &proto.Value_IntValue{
				IntValue: 123,
			},
		},
	})

	assertion.Equal("123", *extracted)
}

func TestMapPropertyModelVersionIdAsValue(t *testing.T) {
	assertion := setup(t)

	extracted := MapPropertyModelVersionIdAsValue(map[string]*proto.Value{
		"model_version_id": {
			Value: &proto.Value_IntValue{
				IntValue: 123,
			},
		},
	})

	assertion.Equal("123", extracted)
}

func TestMapPropertyRegisteredModelId(t *testing.T) {
	assertion := setup(t)

	extracted := MapPropertyRegisteredModelId(map[string]*proto.Value{
		"registered_model_id": {
			Value: &proto.Value_IntValue{
				IntValue: 123,
			},
		},
	})

	assertion.Equal("123", extracted)
}

func TestMapPropertyServingEnvironmentId(t *testing.T) {
	assertion := setup(t)

	extracted := MapPropertyServingEnvironmentId(map[string]*proto.Value{
		"serving_environment_id": {
			Value: &proto.Value_IntValue{
				IntValue: 123,
			},
		},
	})

	assertion.Equal("123", extracted)
}

func TestMapNameFromOwned(t *testing.T) {
	assertion := setup(t)

	name := MapNameFromOwned(of("prefix:name"))
	assertion.Equal("name", *name)

	name = MapNameFromOwned(of("name"))
	assertion.Equal("name", *name)

	name = MapNameFromOwned(of("prefix:name:postfix"))
	assertion.Equal("name", *name)

	name = MapNameFromOwned(nil)
	assertion.Nil(name)
}

func TestMapArtifactType(t *testing.T) {
	assertion := setup(t)

	artifactType, err := MapArtifactType(&proto.Artifact{
		Type: of(ModelArtifactTypeName),
	})
	assertion.Nil(err)
	assertion.Equal("model-artifact", artifactType)

	artifactType, err = MapArtifactType(&proto.Artifact{
		Type: of("Invalid"),
	})
	assertion.NotNil(err)
	assertion.Equal("", artifactType)
}

func TestMapMLMDModelArtifactState(t *testing.T) {
	assertion := setup(t)

	artifactState := MapMLMDModelArtifactState(proto.Artifact_LIVE.Enum())
	assertion.NotNil(artifactState)
	assertion.Equal("LIVE", string(*artifactState))

	artifactState = MapMLMDModelArtifactState(nil)
	assertion.Nil(artifactState)
}

func TestMapRegisteredModelState(t *testing.T) {
	assertion := setup(t)

	state := MapRegisteredModelState(&proto.Context{Properties: map[string]*proto.Value{
		"state": {Value: &proto.Value_StringValue{StringValue: string(openapi.REGISTEREDMODELSTATE_LIVE)}},
	}})
	assertion.NotNil(state)
	assertion.Equal(openapi.REGISTEREDMODELSTATE_LIVE, *state)

	state = MapRegisteredModelState(&proto.Context{Properties: map[string]*proto.Value{}})
	assertion.Nil(state)

	state = MapRegisteredModelState(nil)
	assertion.Nil(state)
}

func TestMapModelVersionState(t *testing.T) {
	assertion := setup(t)

	state := MapModelVersionState(&proto.Context{Properties: map[string]*proto.Value{
		"state": {Value: &proto.Value_StringValue{StringValue: string(openapi.MODELVERSIONSTATE_LIVE)}},
	}})
	assertion.NotNil(state)
	assertion.Equal(openapi.MODELVERSIONSTATE_LIVE, *state)

	state = MapModelVersionState(&proto.Context{Properties: map[string]*proto.Value{}})
	assertion.Nil(state)

	state = MapModelVersionState(nil)
	assertion.Nil(state)
}

func TestMapInferenceServiceState(t *testing.T) {
	assertion := setup(t)

	state := MapInferenceServiceState(&proto.Context{Properties: map[string]*proto.Value{
		"state": {Value: &proto.Value_StringValue{StringValue: string(openapi.INFERENCESERVICESTATE_DEPLOYED)}},
	}})
	assertion.NotNil(state)
	assertion.Equal(openapi.INFERENCESERVICESTATE_DEPLOYED, *state)

	state = MapInferenceServiceState(&proto.Context{Properties: map[string]*proto.Value{}})
	assertion.Nil(state)

	state = MapInferenceServiceState(nil)
	assertion.Nil(state)
}

func TestMapServingEnvironmentType(t *testing.T) {
	assertion := setup(t)

	typeName := MapServingEnvironmentType(&openapi.ServingEnvironment{})
	assertion.NotNil(typeName)
	assertion.Equal(ServingEnvironmentTypeName, *typeName)
}

func TestMapInferenceServiceType(t *testing.T) {
	assertion := setup(t)

	typeName := MapInferenceServiceType(&openapi.InferenceService{})
	assertion.NotNil(typeName)
	assertion.Equal(InferenceServiceTypeName, *typeName)
}

func TestMapInferenceServiceProperties(t *testing.T) {
	assertion := setup(t)

	props, err := MapInferenceServiceProperties(&openapi.InferenceService{
		Description:          of("my custom description"),
		ModelVersionId:       of("1"),
		Runtime:              of("my-runtime"),
		RegisteredModelId:    "2",
		ServingEnvironmentId: "3",
	})
	assertion.Nil(err)
	assertion.Equal(5, len(props))
	assertion.Equal("my custom description", props["description"].GetStringValue())
	assertion.Equal(int64(1), props["model_version_id"].GetIntValue())
	assertion.Equal("my-runtime", props["runtime"].GetStringValue())
	assertion.Equal(int64(2), props["registered_model_id"].GetIntValue())
	assertion.Equal(int64(3), props["serving_environment_id"].GetIntValue())

	// serving and model id must be provided and must be a valid numeric id
	props, err = MapInferenceServiceProperties(&openapi.InferenceService{})
	assertion.NotNil(err)
	assertion.Equal(0, len(props))
}

func TestMapServeModelType(t *testing.T) {
	assertion := setup(t)

	typeName := MapServeModelType(&openapi.ServeModel{})
	assertion.NotNil(typeName)
	assertion.Equal(ServeModelTypeName, *typeName)
}

func TestMapServeModelProperties(t *testing.T) {
	assertion := setup(t)

	props, err := MapServeModelProperties(&openapi.ServeModel{
		Description:    of("my custom description"),
		ModelVersionId: "1",
	})
	assertion.Nil(err)
	assertion.Equal(2, len(props))
	assertion.Equal("my custom description", props["description"].GetStringValue())
	assertion.Equal(int64(1), props["model_version_id"].GetIntValue())

	// serving and model id must be provided and must be a valid numeric id
	props, err = MapServeModelProperties(&openapi.ServeModel{})
	assertion.NotNil(err)
	assertion.Equal(0, len(props))
}

func TestMapServingEnvironmentProperties(t *testing.T) {
	assertion := setup(t)

	props, err := MapServingEnvironmentProperties(&openapi.ServingEnvironment{
		Description: of("my description"),
	})
	assertion.Nil(err)
	assertion.Equal(1, len(props))

	props, err = MapServingEnvironmentProperties(&openapi.ServingEnvironment{})
	assertion.Nil(err)
	assertion.Equal(0, len(props))
}

func TestMapInferenceServiceName(t *testing.T) {
	assertion := setup(t)

	name := MapInferenceServiceName(&OpenAPIModelWrapper[openapi.InferenceService]{
		TypeId:           123,
		ParentResourceId: of("123"),
		ModelName:        of("MyModel"),
		Model: &openapi.InferenceService{
			Name: of("inf-service"),
		},
	})
	assertion.NotNil(name)
	assertion.Equal("123:inf-service", *name)
}

func TestMapServeModelName(t *testing.T) {
	assertion := setup(t)

	name := MapServeModelName(&OpenAPIModelWrapper[openapi.ServeModel]{
		TypeId:           123,
		ParentResourceId: of("parent"),
		Model: &openapi.ServeModel{
			Name: of("v1"),
		},
	})
	assertion.NotNil(name)
	assertion.Equal("parent:v1", *name)

	name = MapServeModelName(&OpenAPIModelWrapper[openapi.ServeModel]{
		TypeId:           123,
		ParentResourceId: of("parent"),
		Model: &openapi.ServeModel{
			Name: nil,
		},
	})
	assertion.NotNil(name)
	assertion.Regexp("parent:.*", *name)

	name = MapServeModelName(&OpenAPIModelWrapper[openapi.ServeModel]{
		TypeId: 123,
		Model: &openapi.ServeModel{
			Name: of("v1"),
		},
	})
	assertion.NotNil(name)
	assertion.Regexp(".*:v1", *name)
}

func TestMapLastKnownState(t *testing.T) {
	assertion := setup(t)

	state, err := MapLastKnownState(of(openapi.EXECUTIONSTATE_RUNNING))
	assertion.Nil(err)
	assertion.Equal("RUNNING", state.String())

	state, err = MapLastKnownState(nil)
	assertion.Nil(err)
	assertion.Nil(state)
}

func TestMapIntProperty(t *testing.T) {
	assertion := setup(t)

	props := map[string]*proto.Value{
		"key": {
			Value: &proto.Value_IntValue{
				IntValue: 10,
			},
		},
	}

	assertion.Equal("10", *MapIntProperty(props, "key"))
	assertion.Nil(MapIntProperty(props, "not-present"))
}

func TestMapIntPropertyAsValue(t *testing.T) {
	assertion := setup(t)

	props := map[string]*proto.Value{
		"key": {
			Value: &proto.Value_IntValue{
				IntValue: 10,
			},
		},
	}

	assertion.Equal("10", MapIntPropertyAsValue(props, "key"))
	assertion.Equal("", MapIntPropertyAsValue(props, "not-present"))
}

func TestMapMLMDServeModelLastKnownState(t *testing.T) {
	assertion := setup(t)

	state := MapMLMDServeModelLastKnownState(of(proto.Execution_COMPLETE))
	assertion.NotNil(state)
	assertion.Equal("COMPLETE", string(*state))

	state = MapMLMDServeModelLastKnownState(nil)
	assertion.Nil(state)
}
