package converter

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/stretchr/testify/assert"
)

func TestMapEmbedMDCustomProperties(t *testing.T) {
	intValue := int32(1)
	stringValue := "test"
	boolValue := true
	doubleValue := 1.0
	jsonByteValue := []byte("{\"test\": \"test\"}")

	jsonMarshaled, err := json.Marshal(jsonByteValue)
	if err != nil {
		t.Fatalf("failed to marshal json value: %v", err)
	}

	expectedByteValue := base64.StdEncoding.EncodeToString(jsonMarshaled)
	mlmdStructString := "mlmd-struct::CiAKCGxhbmd1YWdlEhQyEgoEGgJlbgoEGgJlcwoEGgJjeg=="
	mlmdStructExpected := "{\"language\":[\"en\",\"es\",\"cz\"]}"

	testCases := []struct {
		name     string
		source   []models.Properties
		expected map[string]openapi.MetadataValue
		wantErr  bool
	}{
		{
			name: "test int value",
			source: []models.Properties{
				{
					Name:     "test",
					IntValue: &intValue,
				},
			},
			expected: map[string]openapi.MetadataValue{
				"test": {
					MetadataIntValue: NewMetadataIntValue("1"),
				},
			},
			wantErr: false,
		},
		{
			name: "test int value with nil",
			source: []models.Properties{
				{
					Name:     "test",
					IntValue: nil,
				},
			},
			expected: map[string]openapi.MetadataValue{},
			wantErr:  false,
		},
		{
			name: "test string value",
			source: []models.Properties{
				{
					Name:        "test",
					StringValue: &stringValue,
				},
			},
			expected: map[string]openapi.MetadataValue{
				"test": {
					MetadataStringValue: NewMetadataStringValue("test"),
				},
			},
			wantErr: false,
		},
		{
			name: "test string value with nil",
			source: []models.Properties{
				{
					Name:        "test",
					StringValue: nil,
				},
			},
			expected: map[string]openapi.MetadataValue{},
			wantErr:  false,
		},
		{
			name: "test string value with mlmd struct prefix",
			source: []models.Properties{
				{
					Name:        "test",
					StringValue: &mlmdStructString,
				},
			},
			expected: map[string]openapi.MetadataValue{
				"test": {
					MetadataStructValue: NewMetadataStructValue(mlmdStructExpected),
				},
			},
			wantErr: false,
		},
		{
			name: "test bool value",
			source: []models.Properties{
				{
					Name:      "test",
					BoolValue: &boolValue,
				},
			},
			expected: map[string]openapi.MetadataValue{
				"test": {
					MetadataBoolValue: NewMetadataBoolValue(true),
				},
			},
			wantErr: false,
		},
		{
			name: "test bool value with nil",
			source: []models.Properties{
				{
					Name:      "test",
					BoolValue: nil,
				},
			},
			expected: map[string]openapi.MetadataValue{},
			wantErr:  false,
		},
		{
			name: "test double value",
			source: []models.Properties{
				{
					Name:        "test",
					DoubleValue: &doubleValue,
				},
			},
			expected: map[string]openapi.MetadataValue{
				"test": {
					MetadataDoubleValue: NewMetadataDoubleValue(1.0),
				},
			},
			wantErr: false,
		},
		{
			name: "test double value with nil",
			source: []models.Properties{
				{
					Name:        "test",
					DoubleValue: nil,
				},
			},
			expected: map[string]openapi.MetadataValue{},
			wantErr:  false,
		},
		{
			name: "test byte value",
			source: []models.Properties{
				{
					Name:      "test",
					ByteValue: &jsonByteValue,
				},
			},
			expected: map[string]openapi.MetadataValue{
				"test": {
					MetadataStructValue: NewMetadataStructValue(expectedByteValue),
				},
			},
			wantErr: false,
		},
		{
			name: "test byte value with nil",
			source: []models.Properties{
				{
					Name:      "test",
					ByteValue: nil,
				},
			},
			expected: map[string]openapi.MetadataValue{},
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapEmbedMDCustomProperties(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDDescription(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test description",
			source: &[]models.Properties{
				{
					Name:        "description",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
		{
			name: "test description with nil",
			source: &[]models.Properties{
				{
					Name:        "description",
					StringValue: nil,
				},
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDDescription(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDOwner(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test owner",
			source: &[]models.Properties{
				{
					Name:        "owner",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
		{
			name: "test owner with nil",
			source: &[]models.Properties{
				{
					Name:        "owner",
					StringValue: nil,
				},
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDOwner(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDAuthor(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test author",
			source: &[]models.Properties{
				{
					Name:        "author",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
		{
			name: "test author with nil",
			source: &[]models.Properties{
				{
					Name:        "author",
					StringValue: nil,
				},
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDAuthor(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyLanguage(t *testing.T) {
	languages := "mlmd-struct::CiAKCGxhbmd1YWdlEhQyEgoEGgJlbgoEGgJlcwoEGgJjeg=="

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected []string
	}{
		{
			name: "test language",
			source: &[]models.Properties{
				{
					Name:        "language",
					StringValue: &languages,
				},
			},
			expected: []string{"en", "es", "cz"},
		},
		{
			name: "test language with nil",
			source: &[]models.Properties{
				{
					Name:        "language",
					StringValue: nil,
				},
			},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyLanguage(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyLibraryName(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test library name",
			source: &[]models.Properties{
				{
					Name:        "library_name",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
		{
			name: "test library name with nil",
			source: &[]models.Properties{
				{
					Name:        "library_name",
					StringValue: nil,
				},
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyLibraryName(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyLicenseLink(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test license link",
			source: &[]models.Properties{
				{
					Name:        "license_link",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
		{
			name: "test license link with nil",
			source: &[]models.Properties{
				{
					Name:        "license_link",
					StringValue: nil,
				},
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyLicenseLink(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyLicense(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test license",
			source: &[]models.Properties{
				{
					Name:        "license",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
		{
			name: "test license with nil",
			source: &[]models.Properties{
				{
					Name:        "license",
					StringValue: nil,
				},
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyLicense(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyLogo(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test logo",
			source: &[]models.Properties{
				{
					Name:        "logo",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
		{
			name: "test logo with nil",
			source: &[]models.Properties{
				{
					Name:        "logo",
					StringValue: nil,
				},
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyLogo(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyMaturity(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test maturity",
			source: &[]models.Properties{
				{
					Name:        "maturity",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
		{
			name: "test maturity with nil",
			source: &[]models.Properties{
				{
					Name:        "maturity",
					StringValue: nil,
				},
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyMaturity(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyReadme(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test readme",
			source: &[]models.Properties{
				{
					Name:        "readme",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
		{
			name: "test readme with nil",
			source: &[]models.Properties{
				{
					Name:        "readme",
					StringValue: nil,
				},
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyReadme(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyTasks(t *testing.T) {
	tasks := "mlmd-struct::Ch4KBXRhc2tzEhUyEwoRGg90ZXh0LWdlbmVyYXRpb24="

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected []string
	}{
		{
			name: "test tasks",
			source: &[]models.Properties{
				{
					Name:        "tasks",
					StringValue: &tasks,
				},
			},
			expected: []string{"text-generation"},
		},
		{
			name: "test tasks with nil",
			source: &[]models.Properties{
				{
					Name:        "tasks",
					StringValue: nil,
				},
			},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyTasks(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDStateRegisteredModel(t *testing.T) {
	stringValue := "test"
	validState := openapi.REGISTEREDMODELSTATE_LIVE.Ptr()

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *openapi.RegisteredModelState
		wantErr  bool
	}{
		{
			name: "test state with invalid value",
			source: &[]models.Properties{
				{
					Name:        "state",
					StringValue: &stringValue,
				},
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "test state with nil",
			source: &[]models.Properties{
				{
					Name:        "state",
					StringValue: nil,
				},
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "test state with valid value",
			source: &[]models.Properties{
				{
					Name:        "state",
					StringValue: (*string)(validState),
				},
			},
			expected: validState,
			wantErr:  false,
		},
		{
			name: "test without state",
			source: &[]models.Properties{
				{
					Name: "test",
				},
			},
			expected: nil,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapEmbedMDStateRegisteredModel(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDStateModelVersion(t *testing.T) {
	stringValue := "test"
	validState := openapi.MODELVERSIONSTATE_LIVE.Ptr()

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *openapi.ModelVersionState
		wantErr  bool
	}{
		{
			name: "test state with invalid value",
			source: &[]models.Properties{
				{
					Name:        "state",
					StringValue: &stringValue,
				},
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "test state with nil",
			source: &[]models.Properties{
				{
					Name:        "state",
					StringValue: nil,
				},
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "test state with valid value",
			source: &[]models.Properties{
				{
					Name:        "state",
					StringValue: (*string)(validState),
				},
			},
			expected: validState,
			wantErr:  false,
		},
		{
			name: "test without state",
			source: &[]models.Properties{
				{
					Name: "test",
				},
			},
			expected: nil,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapEmbedMDStateModelVersion(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDExternalIDRegisteredModel(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.RegisteredModelAttributes
		expected *string
	}{
		{
			name: "test external id",
			source: &models.RegisteredModelAttributes{
				ExternalID: &stringValue,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDExternalIDRegisteredModel(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDNameRegisteredModel(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.RegisteredModelAttributes
		expected string
	}{

		{
			name: "test name",
			source: &models.RegisteredModelAttributes{
				Name: &stringValue,
			},
			expected: stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDNameRegisteredModel(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDCreateTimeSinceEpochRegisteredModel(t *testing.T) {
	now := time.Now().UnixMilli()

	testCases := []struct {
		name     string
		source   *models.RegisteredModelAttributes
		expected *string
	}{
		{
			name: "test create time since epoch",
			source: &models.RegisteredModelAttributes{
				CreateTimeSinceEpoch: &now,
			},
			expected: Int64ToString(&now),
		},
		{
			name: "test create time since epoch with nil",
			source: &models.RegisteredModelAttributes{
				CreateTimeSinceEpoch: nil,
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDCreateTimeSinceEpochRegisteredModel(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDLastUpdateTimeSinceEpochRegisteredModel(t *testing.T) {
	now := time.Now().UnixMilli()

	testCases := []struct {
		name     string
		source   *models.RegisteredModelAttributes
		expected *string
	}{
		{
			name: "test last update time since epoch",
			source: &models.RegisteredModelAttributes{
				LastUpdateTimeSinceEpoch: &now,
			},
			expected: Int64ToString(&now),
		},
		{
			name: "test last update time since epoch with nil",
			source: &models.RegisteredModelAttributes{
				LastUpdateTimeSinceEpoch: nil,
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDLastUpdateTimeSinceEpochRegisteredModel(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDNameModelVersion(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.ModelVersionAttributes
		expected string
	}{
		{
			name: "test name",
			source: &models.ModelVersionAttributes{
				Name: &stringValue,
			},
			expected: stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDNameModelVersion(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDExternalIDModelVersion(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.ModelVersionAttributes
		expected *string
	}{
		{
			name: "test external id",
			source: &models.ModelVersionAttributes{
				ExternalID: &stringValue,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDExternalIDModelVersion(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDCreateTimeSinceEpochModelVersion(t *testing.T) {
	now := time.Now().UnixMilli()

	testCases := []struct {
		name     string
		source   *models.ModelVersionAttributes
		expected *string
	}{
		{
			name: "test create time since epoch",
			source: &models.ModelVersionAttributes{
				CreateTimeSinceEpoch: &now,
			},
			expected: Int64ToString(&now),
		},
		{
			name: "test create time since epoch with nil",
			source: &models.ModelVersionAttributes{
				CreateTimeSinceEpoch: nil,
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDCreateTimeSinceEpochModelVersion(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
func TestMapEmbedMDLastUpdateTimeSinceEpochModelVersion(t *testing.T) {
	now := time.Now().UnixMilli()

	testCases := []struct {
		name     string
		source   *models.ModelVersionAttributes
		expected *string
	}{
		{
			name: "test last update time since epoch",
			source: &models.ModelVersionAttributes{
				LastUpdateTimeSinceEpoch: &now,
			},
			expected: Int64ToString(&now),
		},
		{
			name: "test last update time since epoch with nil",
			source: &models.ModelVersionAttributes{
				LastUpdateTimeSinceEpoch: nil,
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDLastUpdateTimeSinceEpochModelVersion(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDExternalIDServingEnvironment(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.ServingEnvironmentAttributes
		expected *string
	}{
		{
			name: "test external id",
			source: &models.ServingEnvironmentAttributes{
				ExternalID: &stringValue,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDExternalIDServingEnvironment(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDNameServingEnvironment(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.ServingEnvironmentAttributes
		expected string
	}{
		{
			name: "test name",
			source: &models.ServingEnvironmentAttributes{
				Name: &stringValue,
			},
			expected: stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDNameServingEnvironment(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDCreateTimeSinceEpochServingEnvironment(t *testing.T) {
	now := time.Now().UnixMilli()

	testCases := []struct {
		name     string
		source   *models.ServingEnvironmentAttributes
		expected *string
	}{
		{
			name: "test create time since epoch",
			source: &models.ServingEnvironmentAttributes{
				CreateTimeSinceEpoch: &now,
			},
			expected: Int64ToString(&now),
		},
		{
			name: "test create time since epoch with nil",
			source: &models.ServingEnvironmentAttributes{
				CreateTimeSinceEpoch: nil,
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDCreateTimeSinceEpochServingEnvironment(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyRuntime(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test runtime",
			source: &[]models.Properties{
				{
					Name:        "runtime",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyRuntime(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
func TestMapEmbedMDExternalIDInferenceService(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.InferenceServiceAttributes
		expected *string
	}{
		{
			name: "test external id",
			source: &models.InferenceServiceAttributes{
				ExternalID: &stringValue,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDExternalIDInferenceService(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
func TestMapEmbedMDPropertyDesiredStateInferenceService(t *testing.T) {
	stringValue := "test"
	validState := openapi.INFERENCESERVICESTATE_DEPLOYED.Ptr()

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *openapi.InferenceServiceState
		wantErr  bool
	}{
		{
			name: "test desired state with invalid value",
			source: &[]models.Properties{
				{
					Name:        "desired_state",
					StringValue: &stringValue,
				},
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "test desired state with nil",
			source: &[]models.Properties{
				{
					Name:        "desired_state",
					StringValue: nil,
				},
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "test desired state with valid value",
			source: &[]models.Properties{
				{
					Name:        "desired_state",
					StringValue: (*string)(validState),
				},
			},
			expected: validState,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapEmbedMDPropertyDesiredStateInferenceService(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyModelVersionId(t *testing.T) {
	intValue := int32(1)

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test model version id",
			source: &[]models.Properties{
				{
					Name:     "model_version_id",
					IntValue: &intValue,
				},
			},
			expected: Int32ToString(&intValue),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyModelVersionId(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyRegisteredModelId(t *testing.T) {
	intValue := int32(1)

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected string
	}{
		{
			name: "test registered model id",
			source: &[]models.Properties{
				{
					Name:     "registered_model_id",
					IntValue: &intValue,
				},
			},
			expected: *Int32ToString(&intValue),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyRegisteredModelId(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
func TestMapEmbedMDPropertyServingEnvironmentId(t *testing.T) {
	intValue := int32(1)

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected string
	}{
		{
			name: "test serving environment id",
			source: &[]models.Properties{
				{
					Name:     "serving_environment_id",
					IntValue: &intValue,
				},
			},
			expected: *Int32ToString(&intValue),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyServingEnvironmentId(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDNameInferenceService(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.InferenceServiceAttributes
		expected *string
	}{
		{
			name: "test name",
			source: &models.InferenceServiceAttributes{
				Name: &stringValue,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDNameInferenceService(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
func TestMapEmbedMDCreateTimeSinceEpochInferenceService(t *testing.T) {
	now := time.Now().UnixMilli()

	testCases := []struct {
		name     string
		source   *models.InferenceServiceAttributes
		expected *string
	}{
		{
			name: "test create time since epoch",
			source: &models.InferenceServiceAttributes{
				CreateTimeSinceEpoch: &now,
			},
			expected: Int64ToString(&now),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDCreateTimeSinceEpochInferenceService(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDLastUpdateTimeSinceEpochInferenceService(t *testing.T) {
	now := time.Now().UnixMilli()

	testCases := []struct {
		name     string
		source   *models.InferenceServiceAttributes
		expected *string
	}{
		{
			name: "test last update time since epoch",
			source: &models.InferenceServiceAttributes{
				LastUpdateTimeSinceEpoch: &now,
			},
			expected: Int64ToString(&now),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDLastUpdateTimeSinceEpochInferenceService(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
func TestMapEmbedMDNameModelArtifact(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.ModelArtifactAttributes
		expected *string
	}{
		{
			name: "test name",
			source: &models.ModelArtifactAttributes{
				Name: &stringValue,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDNameModelArtifact(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDURIModelArtifact(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.ModelArtifactAttributes
		expected *string
	}{
		{
			name: "test uri",
			source: &models.ModelArtifactAttributes{
				URI: &stringValue,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDURIModelArtifact(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDArtifactTypeModelArtifact(t *testing.T) {
	testString := "test"
	stringValue := "model-artifact"

	testCases := []struct {
		name     string
		source   *models.ModelArtifactAttributes
		expected *string
	}{
		{
			name: "test artifact type",
			source: &models.ModelArtifactAttributes{
				ArtifactType: &testString,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDArtifactTypeModelArtifact(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyModelFormatName(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test model format name",
			source: &[]models.Properties{
				{
					Name:        "model_format_name",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyModelFormatName(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyModelFormatVersion(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test model format version",
			source: &[]models.Properties{
				{
					Name:        "model_format_version",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyModelFormatVersion(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyStorageKey(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test storage key",
			source: &[]models.Properties{
				{
					Name:        "storage_key",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyStorageKey(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyStoragePath(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test storage path",
			source: &[]models.Properties{
				{
					Name:        "storage_path",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyStoragePath(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyServiceAccountName(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test service account name",
			source: &[]models.Properties{
				{
					Name:        "service_account_name",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyServiceAccountName(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyModelSourceKind(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test model source kind",
			source: &[]models.Properties{
				{
					Name:        "model_source_kind",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyModelSourceKind(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyModelSourceClass(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test model source class",
			source: &[]models.Properties{
				{
					Name:        "model_source_class",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyModelSourceClass(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyModelSourceGroup(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test model source group",
			source: &[]models.Properties{
				{
					Name:        "model_source_group",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyModelSourceGroup(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyModelSourceId(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test model source id",
			source: &[]models.Properties{
				{
					Name:        "model_source_id",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyModelSourceId(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyModelSourceName(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected *string
	}{
		{
			name: "test model source name",
			source: &[]models.Properties{
				{
					Name:        "model_source_name",
					StringValue: &stringValue,
				},
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDPropertyModelSourceName(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDExternalIDModelArtifact(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.ModelArtifactAttributes
		expected *string
	}{
		{
			name: "test external id",
			source: &models.ModelArtifactAttributes{
				ExternalID: &stringValue,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDExternalIDModelArtifact(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDCreateTimeSinceEpochModelArtifact(t *testing.T) {
	now := time.Now().UnixMilli()

	testCases := []struct {
		name     string
		source   *models.ModelArtifactAttributes
		expected *string
	}{
		{
			name: "test create time since epoch",
			source: &models.ModelArtifactAttributes{
				CreateTimeSinceEpoch: &now,
			},
			expected: Int64ToString(&now),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDCreateTimeSinceEpochModelArtifact(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDLastUpdateTimeSinceEpochModelArtifact(t *testing.T) {
	now := time.Now().UnixMilli()

	testCases := []struct {
		name     string
		source   *models.ModelArtifactAttributes
		expected *string
	}{
		{
			name: "test last update time since epoch",
			source: &models.ModelArtifactAttributes{
				LastUpdateTimeSinceEpoch: &now,
			},
			expected: Int64ToString(&now),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDLastUpdateTimeSinceEpochModelArtifact(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDStateModelArtifact(t *testing.T) {
	invalidState := "invalid"
	validState := openapi.ARTIFACTSTATE_LIVE

	testCases := []struct {
		name     string
		source   *models.ModelArtifactAttributes
		expected *openapi.ArtifactState
		wantErr  bool
	}{
		{
			name: "test state",
			source: &models.ModelArtifactAttributes{
				State: (*string)(&validState),
			},
			expected: &validState,
			wantErr:  false,
		},
		{
			name: "test invalid state",
			source: &models.ModelArtifactAttributes{
				State: (*string)(&invalidState),
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "test default state",
			source: &models.ModelArtifactAttributes{
				State: nil,
			},
			expected: nil,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapEmbedMDStateModelArtifact(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapEmbedMDURIDocArtifact(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.DocArtifactAttributes
		expected *string
	}{
		{
			name: "test uri",
			source: &models.DocArtifactAttributes{
				URI: &stringValue,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDURIDocArtifact(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDArtifactTypeDocArtifact(t *testing.T) {
	testString := "test"
	stringValue := "doc-artifact"

	testCases := []struct {
		name     string
		source   *models.DocArtifactAttributes
		expected *string
	}{
		{
			name: "test artifact type",
			source: &models.DocArtifactAttributes{
				ArtifactType: &testString,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDArtifactTypeDocArtifact(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDExternalIDDocArtifact(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.DocArtifactAttributes
		expected *string
	}{
		{
			name: "test external id",
			source: &models.DocArtifactAttributes{
				ExternalID: &stringValue,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDExternalIDDocArtifact(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDNameDocArtifact(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.DocArtifactAttributes
		expected *string
	}{
		{
			name: "test name",
			source: &models.DocArtifactAttributes{
				Name: &stringValue,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDNameDocArtifact(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDCreateTimeSinceEpochDocArtifact(t *testing.T) {
	now := time.Now().UnixMilli()

	testCases := []struct {
		name     string
		source   *models.DocArtifactAttributes
		expected *string
	}{
		{
			name: "test create time since epoch",
			source: &models.DocArtifactAttributes{
				CreateTimeSinceEpoch: &now,
			},
			expected: Int64ToString(&now),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDCreateTimeSinceEpochDocArtifact(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDLastUpdateTimeSinceEpochDocArtifact(t *testing.T) {
	now := time.Now().UnixMilli()

	testCases := []struct {
		name     string
		source   *models.DocArtifactAttributes
		expected *string
	}{
		{
			name: "test last update time since epoch",
			source: &models.DocArtifactAttributes{
				LastUpdateTimeSinceEpoch: &now,
			},
			expected: Int64ToString(&now),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDLastUpdateTimeSinceEpochDocArtifact(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDStateDocArtifact(t *testing.T) {
	invalidState := "invalid"
	validState := openapi.ARTIFACTSTATE_LIVE

	testCases := []struct {
		name     string
		source   *models.DocArtifactAttributes
		expected *openapi.ArtifactState
		wantErr  bool
	}{
		{
			name: "test state",
			source: &models.DocArtifactAttributes{
				State: (*string)(&validState),
			},
			expected: &validState,
			wantErr:  false,
		},
		{
			name: "test invalid state",
			source: &models.DocArtifactAttributes{
				State: (*string)(&invalidState),
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "test default state",
			source: &models.DocArtifactAttributes{
				State: nil,
			},
			expected: nil,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapEmbedMDStateDocArtifact(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapEmbedMDExternalIDServeModel(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.ServeModelAttributes
		expected *string
	}{
		{
			name: "test external id",
			source: &models.ServeModelAttributes{
				ExternalID: &stringValue,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDExternalIDServeModel(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDNameServeModel(t *testing.T) {
	stringValue := "test"

	testCases := []struct {
		name     string
		source   *models.ServeModelAttributes
		expected *string
	}{
		{
			name: "test name",
			source: &models.ServeModelAttributes{
				Name: &stringValue,
			},
			expected: &stringValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDNameServeModel(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDLastKnownStateServeModel(t *testing.T) {
	validState := openapi.EXECUTIONSTATE_RUNNING

	testCases := []struct {
		name     string
		source   *models.ServeModelAttributes
		expected *openapi.ExecutionState
		wantErr  bool
	}{
		{
			name: "test last known state",
			source: &models.ServeModelAttributes{
				LastKnownState: (*string)(&validState),
			},
			expected: &validState,
			wantErr:  false,
		},
		{
			name: "test default state",
			source: &models.ServeModelAttributes{
				LastKnownState: nil,
			},
			expected: nil,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapEmbedMDLastKnownStateServeModel(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMapEmbedMDCreateTimeSinceEpochServeModel(t *testing.T) {
	now := time.Now().UnixMilli()

	testCases := []struct {
		name     string
		source   *models.ServeModelAttributes
		expected *string
	}{
		{
			name: "test create time since epoch",
			source: &models.ServeModelAttributes{
				CreateTimeSinceEpoch: &now,
			},
			expected: Int64ToString(&now),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDCreateTimeSinceEpochServeModel(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDLastUpdateTimeSinceEpochServeModel(t *testing.T) {
	now := time.Now().UnixMilli()

	testCases := []struct {
		name     string
		source   *models.ServeModelAttributes
		expected *string
	}{
		{
			name: "test last update time since epoch",
			source: &models.ServeModelAttributes{
				LastUpdateTimeSinceEpoch: &now,
			},
			expected: Int64ToString(&now),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := MapEmbedMDLastUpdateTimeSinceEpochServeModel(tc.source)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMapEmbedMDPropertyModelVersionIdServeModel(t *testing.T) {
	intValue := int32(1)

	testCases := []struct {
		name     string
		source   *[]models.Properties
		expected string
		wantErr  bool
	}{
		{
			name: "test model version id",
			source: &[]models.Properties{
				{
					Name:     "model_version_id",
					IntValue: &intValue,
				},
			},
			expected: *Int32ToString(&intValue),
			wantErr:  false,
		},
		{
			name: "test model version id",
			source: &[]models.Properties{
				{
					Name:     "model_version_id",
					IntValue: nil,
				},
			},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := MapEmbedMDPropertyModelVersionIdServeModel(tc.source)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}
