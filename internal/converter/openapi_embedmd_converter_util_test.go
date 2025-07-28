package converter

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/stretchr/testify/assert"
)

func TestMapOpenAPICustomPropertiesEmbedMD(t *testing.T) {
	boolValue := true
	intValue := int32(1)
	doubleValue := 1.0
	stringValue := "test"
	structValue := map[string]interface{}{
		"language": []string{"en", "es", "cz"},
	}
	structValueBytes, err := json.Marshal(structValue)
	if err != nil {
		t.Fatalf("failed to marshal struct value: %v", err)
	}

	expectedStructValue := "mlmd-struct::CiAKCGxhbmd1YWdlEhQyEgoEGgJlbgoEGgJlcwoEGgJjeg=="

	invalidIntValue := "invalid"
	invalidStructValue := "invalid"

	testCases := []struct {
		name     string
		source   *map[string]openapi.MetadataValue
		expected *[]models.Properties
		wantErr  bool
	}{
		{
			name: "test custom properties",
			source: &map[string]openapi.MetadataValue{
				"bool": {
					MetadataBoolValue: &openapi.MetadataBoolValue{
						BoolValue: boolValue,
					},
				},
				"int": {
					MetadataIntValue: &openapi.MetadataIntValue{
						IntValue: strconv.Itoa(int(intValue)),
					},
				},
				"double": {
					MetadataDoubleValue: &openapi.MetadataDoubleValue{
						DoubleValue: doubleValue,
					},
				},
				"string": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue: stringValue,
					},
				},
				"struct": {
					MetadataStructValue: &openapi.MetadataStructValue{
						StructValue: base64.StdEncoding.EncodeToString(structValueBytes),
					},
				},
			},
			expected: &[]models.Properties{
				{
					Name:             "bool",
					BoolValue:        &boolValue,
					IsCustomProperty: true,
				}, {
					Name:             "int",
					IntValue:         &intValue,
					IsCustomProperty: true,
				},
				{
					Name:             "double",
					DoubleValue:      &doubleValue,
					IsCustomProperty: true,
				},
				{
					Name:             "string",
					StringValue:      &stringValue,
					IsCustomProperty: true,
				},
				{
					Name:             "struct",
					StringValue:      &expectedStructValue,
					IsCustomProperty: true,
				},
			},
			wantErr: false,
		},
		{
			name: "test invalid int value",
			source: &map[string]openapi.MetadataValue{
				"int": {
					MetadataIntValue: &openapi.MetadataIntValue{
						IntValue: invalidIntValue,
					},
				},
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "test invalid struct value",
			source: &map[string]openapi.MetadataValue{
				"struct": {
					MetadataStructValue: &openapi.MetadataStructValue{
						StructValue: invalidStructValue,
					},
				},
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapOpenAPICustomPropertiesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Since map iteration order is not guaranteed, we need to compare elements regardless of order
				if tc.expected != nil && actual != nil {
					assert.Equal(t, len(*tc.expected), len(*actual), "Expected and actual slices should have the same length")

					// Create maps for easier comparison
					expectedMap := make(map[string]models.Properties)
					for _, prop := range *tc.expected {
						expectedMap[prop.Name] = prop
					}

					actualMap := make(map[string]models.Properties)
					for _, prop := range *actual {
						actualMap[prop.Name] = prop
					}

					// Compare each property by name
					for name, expectedProp := range expectedMap {
						actualProp, exists := actualMap[name]
						assert.True(t, exists, "Property %s should exist in actual result", name)
						assert.Equal(t, expectedProp, actualProp, "Property %s should match", name)
					}
				} else {
					assert.Equal(t, tc.expected, actual)
				}
			}
		})
	}
}

func TestMapRegisteredModelTypeIDEmbedMD(t *testing.T) {
	testId := int64(1)
	testId32 := int32(testId)

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.RegisteredModel]
		expected *int32
		wantErr  bool
	}{
		{
			name: "test registered model type id",
			source: &OpenAPIModelWrapper[openapi.RegisteredModel]{
				TypeId: testId,
			},
			expected: &testId32,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapRegisteredModelTypeIDEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapRegisteredModelPropertiesEmbedMD(t *testing.T) {
	owner := "test"
	description := "test"
	state := openapi.REGISTEREDMODELSTATE_LIVE
	language := []string{"en", "es", "cz"}
	tasks := []string{"text-generation"}
	provider := "test"
	logo := "test"
	license := "test"
	licenseLink := "test"
	libraryName := "test"
	maturity := "test"
	readme := "test"

	expectedLanguage := "mlmd-struct::CiAKCGxhbmd1YWdlEhQyEgoEGgJlbgoEGgJlcwoEGgJjeg=="
	expectedTasks := "mlmd-struct::Ch4KBXRhc2tzEhUyEwoRGg90ZXh0LWdlbmVyYXRpb24="

	testCases := []struct {
		name     string
		source   *openapi.RegisteredModel
		expected *[]models.Properties
		wantErr  bool
	}{
		{
			name: "test registered model properties",
			source: &openapi.RegisteredModel{
				Owner:       &owner,
				Description: &description,
				State:       &state,
				Language:    language,
				Tasks:       tasks,
				Provider:    &provider,
				Logo:        &logo,
				License:     &license,
				LicenseLink: &licenseLink,
				LibraryName: &libraryName,
				Maturity:    &maturity,
				Readme:      &readme,
			},
			expected: &[]models.Properties{
				{
					Name:             "owner",
					StringValue:      &owner,
					IsCustomProperty: false,
				},
				{
					Name:             "description",
					StringValue:      &description,
					IsCustomProperty: false,
				},
				{
					Name:             "state",
					StringValue:      apiutils.Of(string(state)),
					IsCustomProperty: false,
				},
				{
					Name:             "language",
					StringValue:      &expectedLanguage,
					IsCustomProperty: false,
				},
				{
					Name:             "library_name",
					StringValue:      &libraryName,
					IsCustomProperty: false,
				},
				{
					Name:             "license",
					StringValue:      &license,
					IsCustomProperty: false,
				},
				{
					Name:             "license_link",
					StringValue:      &licenseLink,
					IsCustomProperty: false,
				},
				{
					Name:             "maturity",
					StringValue:      &maturity,
					IsCustomProperty: false,
				},
				{
					Name:             "provider",
					StringValue:      &provider,
					IsCustomProperty: false,
				},
				{
					Name:             "readme",
					StringValue:      &readme,
					IsCustomProperty: false,
				},
				{
					Name:             "logo",
					StringValue:      &logo,
					IsCustomProperty: false,
				},
				{
					Name:             "tasks",
					StringValue:      &expectedTasks,
					IsCustomProperty: false,
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapRegisteredModelPropertiesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapRegisteredModelAttributesEmbedMD(t *testing.T) {
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	name := "test"
	externalId := "test"

	testCases := []struct {
		name     string
		source   *openapi.RegisteredModel
		expected *models.RegisteredModelAttributes
		wantErr  bool
	}{
		{
			name: "test registered model attributes",
			source: &openapi.RegisteredModel{
				Name:                     name,
				CreateTimeSinceEpoch:     &nowStr,
				LastUpdateTimeSinceEpoch: &nowStr,
				ExternalId:               &externalId,
			},
			expected: &models.RegisteredModelAttributes{
				Name:                     &name,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapRegisteredModelAttributesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapModelVersionTypeIDEmbedMD(t *testing.T) {
	testId := int64(1)
	testId32 := int32(testId)

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.ModelVersion]
		expected *int32
		wantErr  bool
	}{
		{
			name: "test model version type id",
			source: &OpenAPIModelWrapper[openapi.ModelVersion]{
				TypeId: testId,
			},
			expected: &testId32,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapModelVersionTypeIDEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapModelVersionPropertiesEmbedMD(t *testing.T) {
	description := "test"
	state := openapi.MODELVERSIONSTATE_LIVE
	author := "test"
	registeredModelId := int32(1)

	testCases := []struct {
		name     string
		source   *openapi.ModelVersion
		expected *[]models.Properties
		wantErr  bool
	}{
		{
			name: "test model version properties",
			source: &openapi.ModelVersion{
				Description:       &description,
				State:             &state,
				Author:            &author,
				RegisteredModelId: strconv.Itoa(int(registeredModelId)),
			},
			expected: &[]models.Properties{
				{
					Name:             "description",
					StringValue:      &description,
					IsCustomProperty: false,
				},
				{
					Name:             "state",
					StringValue:      apiutils.Of(string(state)),
					IsCustomProperty: false,
				},
				{
					Name:             "author",
					StringValue:      &author,
					IsCustomProperty: false,
				},
				{
					Name:             "registered_model_id",
					IntValue:         &registeredModelId,
					IsCustomProperty: false,
				},
			},
			wantErr: false,
		},
		{
			name: "test model version properties with missing registered model id",
			source: &openapi.ModelVersion{
				Description: &description,
				State:       &state,
				Author:      &author,
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapModelVersionPropertiesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapModelVersionAttributesEmbedMD(t *testing.T) {
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	name := "test"
	externalId := "test"
	parentResourceId := "test-parent-id"
	expectedPrefixedName := "test-parent-id:test"

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.ModelVersion]
		expected *models.ModelVersionAttributes
		wantErr  bool
	}{
		{
			name: "test model version attributes",
			source: &OpenAPIModelWrapper[openapi.ModelVersion]{
				Model: &openapi.ModelVersion{
					Name:                     name,
					CreateTimeSinceEpoch:     &nowStr,
					LastUpdateTimeSinceEpoch: &nowStr,
					ExternalId:               &externalId,
				},
				ParentResourceId: &parentResourceId,
			},
			expected: &models.ModelVersionAttributes{
				Name:                     &expectedPrefixedName,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapModelVersionAttributesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapServingEnvironmentTypeIDEmbedMD(t *testing.T) {
	testId := int64(1)
	testId32 := int32(testId)

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.ServingEnvironment]
		expected *int32
		wantErr  bool
	}{
		{
			name: "test serving environment type id",
			source: &OpenAPIModelWrapper[openapi.ServingEnvironment]{
				TypeId: testId,
			},
			expected: &testId32,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapServingEnvironmentTypeIDEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapServingEnvironmentPropertiesEmbedMD(t *testing.T) {
	description := "test"

	testCases := []struct {
		name     string
		source   *openapi.ServingEnvironment
		expected *[]models.Properties
		wantErr  bool
	}{
		{
			name: "test serving environment properties",
			source: &openapi.ServingEnvironment{
				Description: &description,
			},
			expected: &[]models.Properties{
				{
					Name:             "description",
					StringValue:      &description,
					IsCustomProperty: false,
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapServingEnvironmentPropertiesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapServingEnvironmentAttributesEmbedMD(t *testing.T) {
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	name := "test"
	externalId := "test"

	testCases := []struct {
		name     string
		source   *openapi.ServingEnvironment
		expected *models.ServingEnvironmentAttributes
		wantErr  bool
	}{
		{
			name: "test serving environment attributes",
			source: &openapi.ServingEnvironment{
				Name:                     name,
				CreateTimeSinceEpoch:     &nowStr,
				LastUpdateTimeSinceEpoch: &nowStr,
				ExternalId:               &externalId,
			},
			expected: &models.ServingEnvironmentAttributes{
				Name:                     &name,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapServingEnvironmentAttributesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapInferenceServiceTypeIDEmbedMD(t *testing.T) {
	testId := int64(1)
	testId32 := int32(testId)

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.InferenceService]
		expected *int32
		wantErr  bool
	}{
		{
			name: "test inference service type id",
			source: &OpenAPIModelWrapper[openapi.InferenceService]{
				TypeId: testId,
			},
			expected: &testId32,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapInferenceServiceTypeIDEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapInferenceServicePropertiesEmbedMD(t *testing.T) {
	description := "test"
	runtime := "test"
	desiredState := openapi.INFERENCESERVICESTATE_DEPLOYED
	registeredModelId := int32(1)
	servingEnvironmentId := int32(1)

	testCases := []struct {
		name     string
		source   *openapi.InferenceService
		expected *[]models.Properties
		wantErr  bool
	}{
		{
			name: "test inference service properties",
			source: &openapi.InferenceService{
				Description:          &description,
				Runtime:              &runtime,
				DesiredState:         &desiredState,
				RegisteredModelId:    strconv.Itoa(int(registeredModelId)),
				ServingEnvironmentId: strconv.Itoa(int(servingEnvironmentId)),
			},
			expected: &[]models.Properties{
				{
					Name:             "description",
					StringValue:      &description,
					IsCustomProperty: false,
				},
				{
					Name:             "runtime",
					StringValue:      &runtime,
					IsCustomProperty: false,
				},
				{
					Name:             "desired_state",
					StringValue:      apiutils.Of(string(desiredState)),
					IsCustomProperty: false,
				},
				{
					Name:             "registered_model_id",
					IntValue:         &registeredModelId,
					IsCustomProperty: false,
				},
				{
					Name:             "serving_environment_id",
					IntValue:         &servingEnvironmentId,
					IsCustomProperty: false,
				},
			},
			wantErr: false,
		},
		{
			name: "test inference service properties with missing registered model id",
			source: &openapi.InferenceService{
				Description:  &description,
				Runtime:      &runtime,
				DesiredState: &desiredState,
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "test inference service properties with missing serving environment id",
			source: &openapi.InferenceService{
				Description:       &description,
				Runtime:           &runtime,
				DesiredState:      &desiredState,
				RegisteredModelId: strconv.Itoa(int(registeredModelId)),
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapInferenceServicePropertiesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapInferenceServiceAttributesEmbedMD(t *testing.T) {
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	name := "test"
	parentResourceId := "test-parent-id"
	expectedPrefixedName := "test-parent-id:test"
	externalId := "test"

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.InferenceService]
		expected *models.InferenceServiceAttributes
		wantErr  bool
	}{
		{
			name: "test inference service attributes",
			source: &OpenAPIModelWrapper[openapi.InferenceService]{
				Model: &openapi.InferenceService{
					Name:                     &name,
					CreateTimeSinceEpoch:     &nowStr,
					LastUpdateTimeSinceEpoch: &nowStr,
					ExternalId:               &externalId,
				},
				ParentResourceId: &parentResourceId,
			},
			expected: &models.InferenceServiceAttributes{
				Name:                     &expectedPrefixedName,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapInferenceServiceAttributesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapModelArtifactTypeIDEmbedMD(t *testing.T) {
	testId := int64(1)
	testId32 := int32(testId)

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.ModelArtifact]
		expected *int32
		wantErr  bool
	}{
		{
			name: "test model artifact type id",
			source: &OpenAPIModelWrapper[openapi.ModelArtifact]{
				TypeId: testId,
			},
			expected: &testId32,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapModelArtifactTypeIDEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapModelArtifactPropertiesEmbedMD(t *testing.T) {
	description := "test"
	modelFormatName := "test"
	modelFormatVersion := "test"
	storageKey := "test"
	storagePath := "test"
	serviceAccountName := "test"
	modelSourceKind := "test"
	modelSourceClass := "test"
	modelSourceGroup := "test"
	modelSourceId := "test"
	modelSourceName := "test"

	testCases := []struct {
		name     string
		source   *openapi.ModelArtifact
		expected *[]models.Properties
		wantErr  bool
	}{
		{
			name: "test model artifact properties",
			source: &openapi.ModelArtifact{
				Description:        &description,
				ModelFormatName:    &modelFormatName,
				ModelFormatVersion: &modelFormatVersion,
				StorageKey:         &storageKey,
				StoragePath:        &storagePath,
				ServiceAccountName: &serviceAccountName,
				ModelSourceKind:    &modelSourceKind,
				ModelSourceClass:   &modelSourceClass,
				ModelSourceGroup:   &modelSourceGroup,
				ModelSourceId:      &modelSourceId,
				ModelSourceName:    &modelSourceName,
			},
			expected: &[]models.Properties{
				{
					Name:             "description",
					StringValue:      &description,
					IsCustomProperty: false,
				},
				{
					Name:             "model_format_name",
					StringValue:      &modelFormatName,
					IsCustomProperty: false,
				},
				{
					Name:             "model_format_version",
					StringValue:      &modelFormatVersion,
					IsCustomProperty: false,
				},
				{
					Name:             "storage_key",
					StringValue:      &storageKey,
					IsCustomProperty: false,
				},
				{
					Name:             "storage_path",
					StringValue:      &storagePath,
					IsCustomProperty: false,
				},
				{
					Name:             "service_account_name",
					StringValue:      &serviceAccountName,
					IsCustomProperty: false,
				},
				{
					Name:             "model_source_kind",
					StringValue:      &modelSourceKind,
					IsCustomProperty: false,
				},
				{
					Name:             "model_source_class",
					StringValue:      &modelSourceClass,
					IsCustomProperty: false,
				},
				{
					Name:             "model_source_group",
					StringValue:      &modelSourceGroup,
					IsCustomProperty: false,
				},
				{
					Name:             "model_source_id",
					StringValue:      &modelSourceId,
					IsCustomProperty: false,
				},
				{
					Name:             "model_source_name",
					StringValue:      &modelSourceName,
					IsCustomProperty: false,
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapModelArtifactPropertiesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapModelArtifactAttributesEmbedMD(t *testing.T) {
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	name := "test"
	externalId := "test"
	state := openapi.ARTIFACTSTATE_LIVE
	uri := "test"
	parentResourceId := "test-parent-id"
	expectedPrefixedName := "test-parent-id:test"

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.ModelArtifact]
		expected *models.ModelArtifactAttributes
		wantErr  bool
	}{
		{
			name: "test model artifact attributes",
			source: &OpenAPIModelWrapper[openapi.ModelArtifact]{
				Model: &openapi.ModelArtifact{
					Name:                     &name,
					CreateTimeSinceEpoch:     &nowStr,
					LastUpdateTimeSinceEpoch: &nowStr,
					ExternalId:               &externalId,
					State:                    &state,
					Uri:                      &uri,
				},
				ParentResourceId: &parentResourceId,
			},
			expected: &models.ModelArtifactAttributes{
				Name:                     &expectedPrefixedName,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
				State:                    apiutils.Of(string(state)),
				URI:                      &uri,
			},
			wantErr: false,
		},
		{
			name: "test model artifact attributes with nil state",
			source: &OpenAPIModelWrapper[openapi.ModelArtifact]{
				Model: &openapi.ModelArtifact{
					Name:                     &name,
					CreateTimeSinceEpoch:     &nowStr,
					LastUpdateTimeSinceEpoch: &nowStr,
					ExternalId:               &externalId,
					State:                    nil, // nil state should remain nil
					Uri:                      &uri,
				},
				ParentResourceId: &parentResourceId,
			},
			expected: &models.ModelArtifactAttributes{
				Name:                     &expectedPrefixedName,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
				State:                    nil,
				URI:                      &uri,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapModelArtifactAttributesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapDocArtifactTypeIDEmbedMD(t *testing.T) {
	testId := int64(1)
	testId32 := int32(testId)

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.DocArtifact]
		expected *int32
		wantErr  bool
	}{
		{
			name: "test doc artifact type id",
			source: &OpenAPIModelWrapper[openapi.DocArtifact]{
				TypeId: testId,
			},
			expected: &testId32,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapDocArtifactTypeIDEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapDocArtifactPropertiesEmbedMD(t *testing.T) {
	description := "test"

	testCases := []struct {
		name     string
		source   *openapi.DocArtifact
		expected *[]models.Properties
		wantErr  bool
	}{
		{
			name: "test doc artifact properties",
			source: &openapi.DocArtifact{
				Description: &description,
			},
			expected: &[]models.Properties{
				{
					Name:             "description",
					StringValue:      &description,
					IsCustomProperty: false,
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapDocArtifactPropertiesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapDocArtifactAttributesEmbedMD(t *testing.T) {
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	name := "test"
	externalId := "test"
	state := openapi.ARTIFACTSTATE_LIVE
	uri := "test"
	parentResourceId := "test-parent-id"
	expectedPrefixedName := "test-parent-id:test"

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.DocArtifact]
		expected *models.DocArtifactAttributes
		wantErr  bool
	}{
		{
			name: "test doc artifact attributes",
			source: &OpenAPIModelWrapper[openapi.DocArtifact]{
				Model: &openapi.DocArtifact{
					Name:                     &name,
					CreateTimeSinceEpoch:     &nowStr,
					LastUpdateTimeSinceEpoch: &nowStr,
					ExternalId:               &externalId,
					State:                    &state,
					Uri:                      &uri,
				},
				ParentResourceId: &parentResourceId,
			},
			expected: &models.DocArtifactAttributes{
				Name:                     &expectedPrefixedName,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
				State:                    apiutils.Of(string(state)),
				URI:                      &uri,
			},
			wantErr: false,
		},
		{
			name: "test doc artifact attributes with nil state",
			source: &OpenAPIModelWrapper[openapi.DocArtifact]{
				Model: &openapi.DocArtifact{
					Name:                     &name,
					CreateTimeSinceEpoch:     &nowStr,
					LastUpdateTimeSinceEpoch: &nowStr,
					ExternalId:               &externalId,
					State:                    nil, // nil state should remain nil
					Uri:                      &uri,
				},
				ParentResourceId: &parentResourceId,
			},
			expected: &models.DocArtifactAttributes{
				Name:                     &expectedPrefixedName,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
				State:                    nil,
				URI:                      &uri,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapDocArtifactAttributesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapServeModelTypeIDEmbedMD(t *testing.T) {
	testId := int64(1)
	testId32 := int32(testId)

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.ServeModel]
		expected *int32
		wantErr  bool
	}{
		{
			name: "test serve model type id",
			source: &OpenAPIModelWrapper[openapi.ServeModel]{
				TypeId: testId,
			},
			expected: &testId32,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapServeModelTypeIDEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapServeModelPropertiesEmbedMD(t *testing.T) {
	description := "test"
	modelVersionId := int32(1)

	testCases := []struct {
		name     string
		source   *openapi.ServeModel
		expected *[]models.Properties
		wantErr  bool
	}{
		{
			name: "test serve model properties",
			source: &openapi.ServeModel{
				Description:    &description,
				ModelVersionId: strconv.Itoa(int(modelVersionId)),
			},
			expected: &[]models.Properties{
				{
					Name:             "description",
					StringValue:      &description,
					IsCustomProperty: false,
				},
				{
					Name:             "model_version_id",
					IntValue:         &modelVersionId,
					IsCustomProperty: false,
				},
			},
			wantErr: false,
		},
		{
			name: "test serve model properties with missing model version id",
			source: &openapi.ServeModel{
				Description: &description,
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapServeModelPropertiesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapServeModelAttributesEmbedMD(t *testing.T) {
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	name := "test"
	externalId := "test"
	lastKnownState := openapi.EXECUTIONSTATE_RUNNING
	parentResourceId := "test-parent-id"
	expectedPrefixedName := "test-parent-id:test"

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.ServeModel]
		expected *models.ServeModelAttributes
		wantErr  bool
	}{
		{
			name: "test serve model attributes",
			source: &OpenAPIModelWrapper[openapi.ServeModel]{
				Model: &openapi.ServeModel{
					Name:                     &name,
					CreateTimeSinceEpoch:     &nowStr,
					LastUpdateTimeSinceEpoch: &nowStr,
					ExternalId:               &externalId,
					LastKnownState:           &lastKnownState,
				},
				ParentResourceId: &parentResourceId,
			},
			expected: &models.ServeModelAttributes{
				Name:                     &expectedPrefixedName,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
				LastKnownState:           apiutils.Of(string(lastKnownState)),
			},
			wantErr: false,
		},
		{
			name: "test serve model attributes with nil last known state",
			source: &OpenAPIModelWrapper[openapi.ServeModel]{
				Model: &openapi.ServeModel{
					Name:                     &name,
					CreateTimeSinceEpoch:     &nowStr,
					LastUpdateTimeSinceEpoch: &nowStr,
					ExternalId:               &externalId,
					LastKnownState:           nil, // nil state should remain nil
				},
				ParentResourceId: &parentResourceId,
			},
			expected: &models.ServeModelAttributes{
				Name:                     &expectedPrefixedName,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
				LastKnownState:           nil,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapServeModelAttributesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapExperimentRunAttributesEmbedMD(t *testing.T) {
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	name := "test"
	externalId := "test"
	parentResourceId := "test-parent-id"
	expectedPrefixedName := "test-parent-id:test"

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.ExperimentRun]
		expected *models.ExperimentRunAttributes
		wantErr  bool
	}{
		{
			name: "test experiment run attributes",
			source: &OpenAPIModelWrapper[openapi.ExperimentRun]{
				Model: &openapi.ExperimentRun{
					Name:                     &name,
					CreateTimeSinceEpoch:     &nowStr,
					LastUpdateTimeSinceEpoch: &nowStr,
					ExternalId:               &externalId,
				},
				ParentResourceId: &parentResourceId,
			},
			expected: &models.ExperimentRunAttributes{
				Name:                     &expectedPrefixedName,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapExperimentRunAttributesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapMetricAttributesEmbedMD(t *testing.T) {
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	name := "test"
	externalId := "test"
	state := openapi.ARTIFACTSTATE_LIVE
	parentResourceId := "test-parent-id"
	expectedPrefixedName := "test-parent-id:test"

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.Metric]
		expected *models.MetricAttributes
		wantErr  bool
	}{
		{
			name: "test metric attributes",
			source: &OpenAPIModelWrapper[openapi.Metric]{
				Model: &openapi.Metric{
					Name:                     &name,
					CreateTimeSinceEpoch:     &nowStr,
					LastUpdateTimeSinceEpoch: &nowStr,
					ExternalId:               &externalId,
					State:                    &state,
				},
				ParentResourceId: &parentResourceId,
			},
			expected: &models.MetricAttributes{
				Name:                     &expectedPrefixedName,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
				State:                    apiutils.Of(string(state)),
				URI:                      nil, // Metric artifacts don't have URI
			},
			wantErr: false,
		},
		{
			name: "test metric attributes with nil state",
			source: &OpenAPIModelWrapper[openapi.Metric]{
				Model: &openapi.Metric{
					Name:                     &name,
					CreateTimeSinceEpoch:     &nowStr,
					LastUpdateTimeSinceEpoch: &nowStr,
					ExternalId:               &externalId,
					State:                    nil,
				},
				ParentResourceId: &parentResourceId,
			},
			expected: &models.MetricAttributes{
				Name:                     &expectedPrefixedName,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
				State:                    nil,
				URI:                      nil,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapMetricAttributesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapParameterAttributesEmbedMD(t *testing.T) {
	now := time.Now().Unix()
	nowStr := strconv.FormatInt(now, 10)
	name := "test"
	externalId := "test"
	state := openapi.ARTIFACTSTATE_LIVE
	parentResourceId := "test-parent-id"
	expectedPrefixedName := "test-parent-id:test"

	testCases := []struct {
		name     string
		source   *OpenAPIModelWrapper[openapi.Parameter]
		expected *models.ParameterAttributes
		wantErr  bool
	}{
		{
			name: "test parameter attributes",
			source: &OpenAPIModelWrapper[openapi.Parameter]{
				Model: &openapi.Parameter{
					Name:                     &name,
					CreateTimeSinceEpoch:     &nowStr,
					LastUpdateTimeSinceEpoch: &nowStr,
					ExternalId:               &externalId,
					State:                    &state,
				},
				ParentResourceId: &parentResourceId,
			},
			expected: &models.ParameterAttributes{
				Name:                     &expectedPrefixedName,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
				State:                    apiutils.Of(string(state)),
				URI:                      nil, // Parameter artifacts don't have URI
			},
			wantErr: false,
		},
		{
			name: "test parameter attributes with nil state",
			source: &OpenAPIModelWrapper[openapi.Parameter]{
				Model: &openapi.Parameter{
					Name:                     &name,
					CreateTimeSinceEpoch:     &nowStr,
					LastUpdateTimeSinceEpoch: &nowStr,
					ExternalId:               &externalId,
					State:                    nil,
				},
				ParentResourceId: &parentResourceId,
			},
			expected: &models.ParameterAttributes{
				Name:                     &expectedPrefixedName,
				CreateTimeSinceEpoch:     &now,
				LastUpdateTimeSinceEpoch: &now,
				ExternalID:               &externalId,
				State:                    nil,
				URI:                      nil,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapParameterAttributesEmbedMD(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}
