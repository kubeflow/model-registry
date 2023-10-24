package mapper_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/opendatahub-io/model-registry/internal/core/mapper"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

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

	mapper, assert := setup(t)
	asGRPC, err := mapper.MapToProperties(data)
	if err == nil {
		assert.Fail("Did not expected a converted value but an error: %v", asGRPC)
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

	mapper, assert := setup(t)
	asGRPC, err := mapper.MapToProperties(data)
	if err == nil {
		assert.Fail("Did not expected a converted value but an error: %v", asGRPC)
	}
}

func roundTripAndAssert(t *testing.T, data map[string]openapi.MetadataValue, key string) {
	mapper, assert := setup(t)

	// first half
	asGRPC, err := mapper.MapToProperties(data)
	if err != nil {
		t.Error(err)
	}
	assert.Contains(maps.Keys(asGRPC), key)

	// second half
	unmarshall, err := mapper.MapFromProperties(asGRPC)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(data, unmarshall, "result of round-trip shall be equal to original data")
}

func setup(t *testing.T) (*mapper.Mapper, *assert.Assertions) {
	return mapper.NewMapper(1, 2, 3), assert.New(t)
}
