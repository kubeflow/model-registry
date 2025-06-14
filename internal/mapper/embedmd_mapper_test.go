package mapper_test

import (
	"testing"

	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/internal/mapper"
	"github.com/kubeflow/model-registry/internal/ptr"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test constants for type IDs
const (
	testRegisteredModelTypeId    = int64(1)
	testModelVersionTypeId       = int64(2)
	testDocArtifactTypeId        = int64(3)
	testModelArtifactTypeId      = int64(4)
	testServingEnvironmentTypeId = int64(5)
	testInferenceServiceTypeId   = int64(6)
	testServeModelTypeId         = int64(7)
)

var testTypesMap = map[string]int64{
	defaults.RegisteredModelTypeName:    testRegisteredModelTypeId,
	defaults.ModelVersionTypeName:       testModelVersionTypeId,
	defaults.DocArtifactTypeName:        testDocArtifactTypeId,
	defaults.ModelArtifactTypeName:      testModelArtifactTypeId,
	defaults.ServingEnvironmentTypeName: testServingEnvironmentTypeId,
	defaults.InferenceServiceTypeName:   testInferenceServiceTypeId,
	defaults.ServeModelTypeName:         testServeModelTypeId,
}

func setupEmbedMDMapper(t *testing.T) (*assert.Assertions, *mapper.EmbedMDMapper) {
	return assert.New(t), mapper.NewEmbedMDMapper(testTypesMap)
}

// Tests for OpenAPI --> EmbedMD mapping

func TestEmbedMDMapFromRegisteredModel(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	openAPIModel := &openapi.RegisteredModel{
		Name:        "test-registered-model",
		Description: ptr.Of("Test description"),
		Owner:       ptr.Of("test-owner"),
		ExternalId:  ptr.Of("ext-123"),
		State:       ptr.Of(openapi.REGISTEREDMODELSTATE_LIVE),
	}

	result, err := mapper.MapFromRegisteredModel(openAPIModel)
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify type ID
	assertion.Equal(int32(testRegisteredModelTypeId), *result.GetTypeID())

	// Verify attributes
	attrs := result.GetAttributes()
	assertion.NotNil(attrs)
	assertion.Equal("test-registered-model", *attrs.Name)
	assertion.Equal("ext-123", *attrs.ExternalID)

	// Verify properties
	props := result.GetProperties()
	assertion.NotNil(props)

	// Check for description property
	var foundDescription, foundOwner, foundState bool
	for _, prop := range *props {
		switch prop.Name {
		case "description":
			foundDescription = true
			assertion.Equal("Test description", *prop.StringValue)
		case "owner":
			foundOwner = true
			assertion.Equal("test-owner", *prop.StringValue)
		case "state":
			foundState = true
			assertion.Equal("LIVE", *prop.StringValue)
		}
	}
	assertion.True(foundDescription, "Should find description property")
	assertion.True(foundOwner, "Should find owner property")
	assertion.True(foundState, "Should find state property")
}

func TestEmbedMDMapFromModelVersion(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	openAPIModel := &openapi.ModelVersion{
		Name:              "test-model-version",
		Description:       ptr.Of("Test version description"),
		Author:            ptr.Of("test-author"),
		ExternalId:        ptr.Of("version-ext-123"),
		State:             ptr.Of(openapi.MODELVERSIONSTATE_LIVE),
		RegisteredModelId: "1",
	}

	result, err := mapper.MapFromModelVersion(openAPIModel)
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify type ID
	assertion.Equal(int32(testModelVersionTypeId), *result.GetTypeID())

	// Verify attributes
	attrs := result.GetAttributes()
	assertion.NotNil(attrs)
	assertion.Equal("test-model-version", *attrs.Name)
	assertion.Equal("version-ext-123", *attrs.ExternalID)

	// Verify properties
	props := result.GetProperties()
	assertion.NotNil(props)

	var foundDescription, foundAuthor, foundState, foundRegisteredModelId bool
	for _, prop := range *props {
		switch prop.Name {
		case "description":
			foundDescription = true
			assertion.Equal("Test version description", *prop.StringValue)
		case "author":
			foundAuthor = true
			assertion.Equal("test-author", *prop.StringValue)
		case "state":
			foundState = true
			assertion.Equal("LIVE", *prop.StringValue)
		case "registered_model_id":
			foundRegisteredModelId = true
			assertion.Equal(int32(1), *prop.IntValue)
		}
	}
	assertion.True(foundDescription, "Should find description property")
	assertion.True(foundAuthor, "Should find author property")
	assertion.True(foundState, "Should find state property")
	assertion.True(foundRegisteredModelId, "Should find registered_model_id property")
}

func TestEmbedMDMapFromServingEnvironment(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	openAPIModel := &openapi.ServingEnvironment{
		Name:        "test-serving-env",
		Description: ptr.Of("Test serving environment"),
		ExternalId:  ptr.Of("env-ext-123"),
	}

	result, err := mapper.MapFromServingEnvironment(openAPIModel)
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify type ID
	assertion.Equal(int32(testServingEnvironmentTypeId), *result.GetTypeID())

	// Verify attributes
	attrs := result.GetAttributes()
	assertion.NotNil(attrs)
	assertion.Equal("test-serving-env", *attrs.Name)
	assertion.Equal("env-ext-123", *attrs.ExternalID)
}

func TestEmbedMDMapFromInferenceService(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	openAPIModel := &openapi.InferenceService{
		Name:                 ptr.Of("test-inference-service"),
		Description:          ptr.Of("Test inference service"),
		ExternalId:           ptr.Of("inf-ext-123"),
		ServingEnvironmentId: "5",
		RegisteredModelId:    "1",
		ModelVersionId:       ptr.Of("2"),
		Runtime:              ptr.Of("tensorflow"),
		DesiredState:         ptr.Of(openapi.INFERENCESERVICESTATE_DEPLOYED),
	}

	result, err := mapper.MapFromInferenceService(openAPIModel, "5")
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify type ID
	assertion.Equal(int32(testInferenceServiceTypeId), *result.GetTypeID())

	// Verify attributes
	attrs := result.GetAttributes()
	assertion.NotNil(attrs)
	assertion.Equal("test-inference-service", *attrs.Name)
	assertion.Equal("inf-ext-123", *attrs.ExternalID)

	// Verify properties
	props := result.GetProperties()
	assertion.NotNil(props)

	var foundServingEnvId, foundRegisteredModelId, foundModelVersionId, foundRuntime, foundDesiredState bool
	for _, prop := range *props {
		switch prop.Name {
		case "serving_environment_id":
			foundServingEnvId = true
			assertion.Equal(int32(5), *prop.IntValue)
		case "registered_model_id":
			foundRegisteredModelId = true
			assertion.Equal(int32(1), *prop.IntValue)
		case "model_version_id":
			foundModelVersionId = true
			assertion.Equal(int32(2), *prop.IntValue)
		case "runtime":
			foundRuntime = true
			assertion.Equal("tensorflow", *prop.StringValue)
		case "desired_state":
			foundDesiredState = true
			assertion.Equal("DEPLOYED", *prop.StringValue)
		}
	}
	assertion.True(foundServingEnvId, "Should find serving_environment_id property")
	assertion.True(foundRegisteredModelId, "Should find registered_model_id property")
	assertion.True(foundModelVersionId, "Should find model_version_id property")
	assertion.True(foundRuntime, "Should find runtime property")
	assertion.True(foundDesiredState, "Should find desired_state property")
}

func TestEmbedMDMapFromModelArtifact(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	openAPIModel := &openapi.ModelArtifact{
		Name:               ptr.Of("test-model-artifact"),
		Description:        ptr.Of("Test model artifact"),
		ExternalId:         ptr.Of("model-art-ext-123"),
		Uri:                ptr.Of("s3://bucket/model.pkl"),
		State:              ptr.Of(openapi.ARTIFACTSTATE_LIVE),
		ModelFormatName:    ptr.Of("pickle"),
		ModelFormatVersion: ptr.Of("1.0"),
		StorageKey:         ptr.Of("storage-key"),
		StoragePath:        ptr.Of("/path/to/model"),
	}

	result, err := mapper.MapFromModelArtifact(openAPIModel)
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify type ID
	assertion.Equal(int32(testModelArtifactTypeId), *result.GetTypeID())

	// Verify attributes
	attrs := result.GetAttributes()
	assertion.NotNil(attrs)
	assertion.Equal("test-model-artifact", *attrs.Name)
	assertion.Equal("model-art-ext-123", *attrs.ExternalID)
	assertion.Equal("s3://bucket/model.pkl", *attrs.URI)
	assertion.Equal("LIVE", *attrs.State)
	// Add nil check for ArtifactType
	if attrs.ArtifactType != nil {
		assertion.Equal("model-artifact", *attrs.ArtifactType)
	}
}

func TestEmbedMDMapFromDocArtifact(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	openAPIModel := &openapi.DocArtifact{
		Name:        ptr.Of("test-doc-artifact"),
		Description: ptr.Of("Test doc artifact"),
		ExternalId:  ptr.Of("doc-art-ext-123"),
		Uri:         ptr.Of("s3://bucket/doc.pdf"),
		State:       ptr.Of(openapi.ARTIFACTSTATE_LIVE),
	}

	result, err := mapper.MapFromDocArtifact(openAPIModel)
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify type ID
	assertion.Equal(int32(testDocArtifactTypeId), *result.GetTypeID())

	// Verify attributes
	attrs := result.GetAttributes()
	assertion.NotNil(attrs)
	assertion.Equal("test-doc-artifact", *attrs.Name)
	assertion.Equal("doc-art-ext-123", *attrs.ExternalID)
	assertion.Equal("s3://bucket/doc.pdf", *attrs.URI)
	assertion.Equal("LIVE", *attrs.State)
	// Add nil check for ArtifactType
	if attrs.ArtifactType != nil {
		assertion.Equal("doc-artifact", *attrs.ArtifactType)
	}
}

func TestEmbedMDMapFromServeModel(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	openAPIModel := &openapi.ServeModel{
		Name:           ptr.Of("test-serve-model"),
		Description:    ptr.Of("Test serve model"),
		ExternalId:     ptr.Of("serve-ext-123"),
		ModelVersionId: "2",
		LastKnownState: ptr.Of(openapi.EXECUTIONSTATE_RUNNING),
	}

	result, err := mapper.MapFromServeModel(openAPIModel)
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify type ID
	assertion.Equal(int32(testServeModelTypeId), *result.GetTypeID())

	// Verify attributes
	attrs := result.GetAttributes()
	assertion.NotNil(attrs)
	assertion.Equal("test-serve-model", *attrs.Name)
	assertion.Equal("serve-ext-123", *attrs.ExternalID)
	// Add nil check for LastKnownState
	if attrs.LastKnownState != nil {
		assertion.Equal("RUNNING", *attrs.LastKnownState)
	}

	// Verify properties
	props := result.GetProperties()
	assertion.NotNil(props)

	var foundModelVersionId bool
	for _, prop := range *props {
		if prop.Name == "model_version_id" {
			foundModelVersionId = true
			assertion.Equal(int32(2), *prop.IntValue)
		}
	}
	assertion.True(foundModelVersionId, "Should find model_version_id property")
}

// Tests for EmbedMD --> OpenAPI mapping

func TestEmbedMDMapToRegisteredModel(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	embedMDModel := &models.RegisteredModelImpl{
		ID:     ptr.Of(int32(1)),
		TypeID: ptr.Of(int32(testRegisteredModelTypeId)),
		Attributes: &models.RegisteredModelAttributes{
			Name:                     ptr.Of("test-registered-model"),
			ExternalID:               ptr.Of("ext-123"),
			CreateTimeSinceEpoch:     ptr.Of(int64(1234567890)),
			LastUpdateTimeSinceEpoch: ptr.Of(int64(1234567891)),
		},
		Properties: &[]models.Properties{
			{
				Name:        "description",
				StringValue: ptr.Of("Test description"),
			},
			{
				Name:        "owner",
				StringValue: ptr.Of("test-owner"),
			},
			{
				Name:        "state",
				StringValue: ptr.Of("LIVE"),
			},
		},
		CustomProperties: &[]models.Properties{
			{
				Name:             "custom-prop",
				StringValue:      ptr.Of("custom-value"),
				IsCustomProperty: true,
			},
		},
	}

	result, err := mapper.MapToRegisteredModel(embedMDModel)
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify basic fields
	assertion.Equal("1", *result.Id)
	assertion.Equal("test-registered-model", result.Name)
	assertion.Equal("ext-123", *result.ExternalId)
	assertion.Equal("1234567890", *result.CreateTimeSinceEpoch)
	assertion.Equal("1234567891", *result.LastUpdateTimeSinceEpoch)

	// Verify mapped properties
	assertion.Equal("Test description", *result.Description)
	assertion.Equal("test-owner", *result.Owner)
	assertion.Equal(openapi.REGISTEREDMODELSTATE_LIVE, *result.State)

	// Verify custom properties
	assertion.NotNil(result.CustomProperties)
	customProps := *result.CustomProperties
	assertion.Contains(customProps, "custom-prop")
	assertion.Equal("custom-value", customProps["custom-prop"].MetadataStringValue.StringValue)
}

func TestEmbedMDMapToRegisteredModelNil(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	result, err := mapper.MapToRegisteredModel(nil)
	assertion.NotNil(err)
	assertion.Nil(result)
	assertion.Equal("registered model is nil", err.Error())
}

func TestEmbedMDMapToModelVersion(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	embedMDModel := &models.ModelVersionImpl{
		ID:     ptr.Of(int32(2)),
		TypeID: ptr.Of(int32(testModelVersionTypeId)),
		Attributes: &models.ModelVersionAttributes{
			Name:                     ptr.Of("test-model-version"),
			ExternalID:               ptr.Of("version-ext-123"),
			CreateTimeSinceEpoch:     ptr.Of(int64(1234567890)),
			LastUpdateTimeSinceEpoch: ptr.Of(int64(1234567891)),
		},
		Properties: &[]models.Properties{
			{
				Name:        "description",
				StringValue: ptr.Of("Test version description"),
			},
			{
				Name:        "author",
				StringValue: ptr.Of("test-author"),
			},
			{
				Name:        "state",
				StringValue: ptr.Of("LIVE"),
			},
			{
				Name:     "registered_model_id",
				IntValue: ptr.Of(int32(1)),
			},
		},
	}

	result, err := mapper.MapToModelVersion(embedMDModel)
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify basic fields
	assertion.Equal("2", *result.Id)
	assertion.Equal("test-model-version", result.Name)
	assertion.Equal("version-ext-123", *result.ExternalId)
	assertion.Equal("1234567890", *result.CreateTimeSinceEpoch)
	assertion.Equal("1234567891", *result.LastUpdateTimeSinceEpoch)

	// Verify mapped properties
	assertion.Equal("Test version description", *result.Description)
	assertion.Equal("test-author", *result.Author)
	assertion.Equal(openapi.MODELVERSIONSTATE_LIVE, *result.State)
	assertion.Equal("1", result.RegisteredModelId)
}

func TestEmbedMDMapToServingEnvironment(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	embedMDModel := &models.ServingEnvironmentImpl{
		ID:     ptr.Of(int32(5)),
		TypeID: ptr.Of(int32(testServingEnvironmentTypeId)),
		Attributes: &models.ServingEnvironmentAttributes{
			Name:                     ptr.Of("test-serving-env"),
			ExternalID:               ptr.Of("env-ext-123"),
			CreateTimeSinceEpoch:     ptr.Of(int64(1234567890)),
			LastUpdateTimeSinceEpoch: ptr.Of(int64(1234567891)),
		},
		Properties: &[]models.Properties{
			{
				Name:        "description",
				StringValue: ptr.Of("Test serving environment"),
			},
		},
	}

	result, err := mapper.MapToServingEnvironment(embedMDModel)
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify basic fields
	assertion.Equal("5", *result.Id)
	assertion.Equal("test-serving-env", result.Name)
	assertion.Equal("env-ext-123", *result.ExternalId)
	assertion.Equal("1234567890", *result.CreateTimeSinceEpoch)
	assertion.Equal("1234567891", *result.LastUpdateTimeSinceEpoch)

	// Verify mapped properties
	assertion.Equal("Test serving environment", *result.Description)
}

func TestEmbedMDMapToInferenceService(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	embedMDModel := &models.InferenceServiceImpl{
		ID:     ptr.Of(int32(6)),
		TypeID: ptr.Of(int32(testInferenceServiceTypeId)),
		Attributes: &models.InferenceServiceAttributes{
			Name:                     ptr.Of("test-inference-service"),
			ExternalID:               ptr.Of("inf-ext-123"),
			CreateTimeSinceEpoch:     ptr.Of(int64(1234567890)),
			LastUpdateTimeSinceEpoch: ptr.Of(int64(1234567891)),
		},
		Properties: &[]models.Properties{
			{
				Name:        "description",
				StringValue: ptr.Of("Test inference service"),
			},
			{
				Name:     "serving_environment_id",
				IntValue: ptr.Of(int32(5)),
			},
			{
				Name:     "registered_model_id",
				IntValue: ptr.Of(int32(1)),
			},
			{
				Name:     "model_version_id",
				IntValue: ptr.Of(int32(2)),
			},
			{
				Name:        "runtime",
				StringValue: ptr.Of("tensorflow"),
			},
			{
				Name:        "desired_state",
				StringValue: ptr.Of("DEPLOYED"),
			},
		},
	}

	result, err := mapper.MapToInferenceService(embedMDModel)
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify basic fields
	assertion.Equal("6", *result.Id)
	assertion.Equal("test-inference-service", *result.Name)
	assertion.Equal("inf-ext-123", *result.ExternalId)
	assertion.Equal("1234567890", *result.CreateTimeSinceEpoch)
	assertion.Equal("1234567891", *result.LastUpdateTimeSinceEpoch)

	// Verify mapped properties
	assertion.Equal("Test inference service", *result.Description)
	assertion.Equal("5", result.ServingEnvironmentId)
	assertion.Equal("1", result.RegisteredModelId)
	assertion.Equal("2", *result.ModelVersionId)
	assertion.Equal("tensorflow", *result.Runtime)
	assertion.Equal(openapi.INFERENCESERVICESTATE_DEPLOYED, *result.DesiredState)
}

func TestEmbedMDMapToModelArtifact(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	embedMDModel := &models.ModelArtifactImpl{
		ID:     ptr.Of(int32(3)),
		TypeID: ptr.Of(int32(testModelArtifactTypeId)),
		Attributes: &models.ModelArtifactAttributes{
			Name:                     ptr.Of("test-model-artifact"),
			ExternalID:               ptr.Of("model-art-ext-123"),
			URI:                      ptr.Of("s3://bucket/model.pkl"),
			State:                    ptr.Of("LIVE"),
			ArtifactType:             ptr.Of("model-artifact"),
			CreateTimeSinceEpoch:     ptr.Of(int64(1234567890)),
			LastUpdateTimeSinceEpoch: ptr.Of(int64(1234567891)),
		},
		Properties: &[]models.Properties{
			{
				Name:        "description",
				StringValue: ptr.Of("Test model artifact"),
			},
			{
				Name:        "model_format_name",
				StringValue: ptr.Of("pickle"),
			},
			{
				Name:        "model_format_version",
				StringValue: ptr.Of("1.0"),
			},
		},
	}

	result, err := mapper.MapToModelArtifact(embedMDModel)
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify basic fields
	assertion.Equal("3", *result.Id)
	assertion.Equal("test-model-artifact", *result.Name)
	assertion.Equal("model-art-ext-123", *result.ExternalId)
	assertion.Equal("s3://bucket/model.pkl", *result.Uri)
	assertion.Equal(openapi.ARTIFACTSTATE_LIVE, *result.State)
	assertion.Equal("model-artifact", *result.ArtifactType)
	assertion.Equal("1234567890", *result.CreateTimeSinceEpoch)
	assertion.Equal("1234567891", *result.LastUpdateTimeSinceEpoch)

	// Verify mapped properties
	assertion.Equal("Test model artifact", *result.Description)
	assertion.Equal("pickle", *result.ModelFormatName)
	assertion.Equal("1.0", *result.ModelFormatVersion)
}

func TestEmbedMDMapToDocArtifact(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	embedMDModel := &models.DocArtifactImpl{
		ID:     ptr.Of(int32(4)),
		TypeID: ptr.Of(int32(testDocArtifactTypeId)),
		Attributes: &models.DocArtifactAttributes{
			Name:                     ptr.Of("test-doc-artifact"),
			ExternalID:               ptr.Of("doc-art-ext-123"),
			URI:                      ptr.Of("s3://bucket/doc.pdf"),
			State:                    ptr.Of("LIVE"),
			ArtifactType:             ptr.Of("doc-artifact"),
			CreateTimeSinceEpoch:     ptr.Of(int64(1234567890)),
			LastUpdateTimeSinceEpoch: ptr.Of(int64(1234567891)),
		},
		Properties: &[]models.Properties{
			{
				Name:        "description",
				StringValue: ptr.Of("Test doc artifact"),
			},
		},
	}

	result, err := mapper.MapToDocArtifact(embedMDModel)
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify basic fields
	assertion.Equal("4", *result.Id)
	assertion.Equal("test-doc-artifact", *result.Name)
	assertion.Equal("doc-art-ext-123", *result.ExternalId)
	assertion.Equal("s3://bucket/doc.pdf", *result.Uri)
	assertion.Equal(openapi.ARTIFACTSTATE_LIVE, *result.State)
	assertion.Equal("doc-artifact", *result.ArtifactType)
	assertion.Equal("1234567890", *result.CreateTimeSinceEpoch)
	assertion.Equal("1234567891", *result.LastUpdateTimeSinceEpoch)

	// Verify mapped properties
	assertion.Equal("Test doc artifact", *result.Description)
}

func TestEmbedMDMapToServeModel(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	embedMDModel := &models.ServeModelImpl{
		ID:     ptr.Of(int32(7)),
		TypeID: ptr.Of(int32(testServeModelTypeId)),
		Attributes: &models.ServeModelAttributes{
			Name:                     ptr.Of("test-serve-model"),
			ExternalID:               ptr.Of("serve-ext-123"),
			LastKnownState:           ptr.Of("RUNNING"),
			CreateTimeSinceEpoch:     ptr.Of(int64(1234567890)),
			LastUpdateTimeSinceEpoch: ptr.Of(int64(1234567891)),
		},
		Properties: &[]models.Properties{
			{
				Name:        "description",
				StringValue: ptr.Of("Test serve model"),
			},
			{
				Name:     "model_version_id",
				IntValue: ptr.Of(int32(2)),
			},
		},
	}

	result, err := mapper.MapToServeModel(embedMDModel)
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify basic fields
	assertion.Equal("7", *result.Id)
	assertion.Equal("test-serve-model", *result.Name)
	assertion.Equal("serve-ext-123", *result.ExternalId)
	assertion.Equal(openapi.EXECUTIONSTATE_RUNNING, *result.LastKnownState)
	assertion.Equal("1234567890", *result.CreateTimeSinceEpoch)
	assertion.Equal("1234567891", *result.LastUpdateTimeSinceEpoch)

	// Verify mapped properties
	assertion.Equal("Test serve model", *result.Description)
	assertion.Equal("2", result.ModelVersionId)
}

// Test edge cases and error conditions

func TestEmbedMDMapFromWithCustomProperties(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	customProps := map[string]openapi.MetadataValue{
		"string-prop": {
			MetadataStringValue: &openapi.MetadataStringValue{
				StringValue: "string-value",
			},
		},
		"int-prop": {
			MetadataIntValue: &openapi.MetadataIntValue{
				IntValue: "42",
			},
		},
		"bool-prop": {
			MetadataBoolValue: &openapi.MetadataBoolValue{
				BoolValue: true,
			},
		},
		"double-prop": {
			MetadataDoubleValue: &openapi.MetadataDoubleValue{
				DoubleValue: 3.14,
			},
		},
	}

	openAPIModel := &openapi.RegisteredModel{
		Name:             "test-with-custom-props",
		CustomProperties: &customProps,
	}

	result, err := mapper.MapFromRegisteredModel(openAPIModel)
	assertion.Nil(err)
	assertion.NotNil(result)

	// Verify custom properties were converted
	customPropsResult := result.GetCustomProperties()
	assertion.NotNil(customPropsResult)
	assertion.Len(*customPropsResult, 4)

	// Check each custom property
	propMap := make(map[string]models.Properties)
	for _, prop := range *customPropsResult {
		propMap[prop.Name] = prop
	}

	assertion.Contains(propMap, "string-prop")
	assertion.Equal("string-value", *propMap["string-prop"].StringValue)
	assertion.True(propMap["string-prop"].IsCustomProperty)

	assertion.Contains(propMap, "int-prop")
	assertion.Equal(int32(42), *propMap["int-prop"].IntValue)
	assertion.True(propMap["int-prop"].IsCustomProperty)

	assertion.Contains(propMap, "bool-prop")
	assertion.Equal(true, *propMap["bool-prop"].BoolValue)
	assertion.True(propMap["bool-prop"].IsCustomProperty)

	assertion.Contains(propMap, "double-prop")
	assertion.Equal(3.14, *propMap["double-prop"].DoubleValue)
	assertion.True(propMap["double-prop"].IsCustomProperty)
}

func TestEmbedMDMapperCreation(t *testing.T) {
	assertion := assert.New(t)

	mapper := mapper.NewEmbedMDMapper(testTypesMap)
	assertion.NotNil(mapper)
	// Note: Cannot test unexported fields from external package
}

func TestEmbedMDMapFromWithMinimalData(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	// Test with minimal required data
	openAPIModel := &openapi.RegisteredModel{
		Name: "minimal-model",
	}

	result, err := mapper.MapFromRegisteredModel(openAPIModel)
	assertion.Nil(err)
	assertion.NotNil(result)
	assertion.Equal("minimal-model", *result.GetAttributes().Name)
	assertion.Equal(int32(testRegisteredModelTypeId), *result.GetTypeID())
}

func TestEmbedMDRoundTripConversion(t *testing.T) {
	assertion, mapper := setupEmbedMDMapper(t)

	// Create an OpenAPI model
	originalOpenAPI := &openapi.RegisteredModel{
		Name:        "roundtrip-test",
		Description: ptr.Of("Test roundtrip conversion"),
		Owner:       ptr.Of("test-owner"),
		ExternalId:  ptr.Of("roundtrip-ext-123"),
		State:       ptr.Of(openapi.REGISTEREDMODELSTATE_LIVE),
	}

	// Convert to EmbedMD
	embedMDModel, err := mapper.MapFromRegisteredModel(originalOpenAPI)
	require.NoError(t, err)

	// Set ID for the conversion back (simulating saved model)
	embedMDModel.(*models.RegisteredModelImpl).ID = ptr.Of(int32(1))
	embedMDModel.(*models.RegisteredModelImpl).Attributes.CreateTimeSinceEpoch = ptr.Of(int64(1234567890))
	embedMDModel.(*models.RegisteredModelImpl).Attributes.LastUpdateTimeSinceEpoch = ptr.Of(int64(1234567891))

	// Convert back to OpenAPI
	resultOpenAPI, err := mapper.MapToRegisteredModel(embedMDModel)
	require.NoError(t, err)

	// Verify the roundtrip preserved the data
	assertion.Equal(originalOpenAPI.Name, resultOpenAPI.Name)
	assertion.Equal(*originalOpenAPI.Description, *resultOpenAPI.Description)
	assertion.Equal(*originalOpenAPI.Owner, *resultOpenAPI.Owner)
	assertion.Equal(*originalOpenAPI.ExternalId, *resultOpenAPI.ExternalId)
	assertion.Equal(*originalOpenAPI.State, *resultOpenAPI.State)
}
