package converter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// Int32ToString converts int32 to string-based one
func Int32ToString(id *int32) *string {
	if id == nil {
		return nil
	}
	idAsString := strconv.FormatInt(int64(*id), 10)
	return &idAsString
}

// Int64ToInt32 converts int64 to int32 if numeric, otherwise return error
func Int64ToInt32(id *int64) (*int32, error) {
	if id == nil {
		return nil, nil
	}
	if *id > math.MaxInt32 || *id < math.MinInt32 {
		return nil, fmt.Errorf("id is out of range of int32: %d", *id)
	}

	idInt32 := int32(*id)
	return &idInt32, nil
}

// MapOpenAPICustomPropertiesEmbedMD maps OpenAPI custom properties model to embedmd one
func MapOpenAPICustomPropertiesEmbedMD(source *map[string]openapi.MetadataValue) (*[]models.Properties, error) {
	props := make([]models.Properties, 0)

	if source != nil {
		for key, v := range *source {
			value := models.Properties{}

			value.Name = key
			value.IsCustomProperty = true

			switch {
			// bool value
			case v.MetadataBoolValue != nil:
				value.BoolValue = &v.MetadataBoolValue.BoolValue
			// int value
			case v.MetadataIntValue != nil:
				intValue, err := StringToInt32(v.MetadataIntValue.IntValue)
				if err != nil {
					return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, key)
				}
				value.IntValue = &intValue
			// double value
			case v.MetadataDoubleValue != nil:
				value.DoubleValue = &v.MetadataDoubleValue.DoubleValue
			// string value
			case v.MetadataStringValue != nil:
				value.StringValue = &v.MetadataStringValue.StringValue
			// struct value
			case v.MetadataStructValue != nil:
				base64Decoded, err := base64.StdEncoding.DecodeString(v.MetadataStructValue.StructValue)
				if err != nil {
					return nil, fmt.Errorf("%w: unable to decode %w for key %s", api.ErrBadRequest, err, key)
				}

				var structValue structpb.Struct
				err = json.Unmarshal(base64Decoded, &structValue)
				if err != nil {
					return nil, fmt.Errorf("%w: unable to decode %w for key %s", api.ErrBadRequest, err, key)
				}
				encodedStruct, err := encodeStruct(&structValue)
				if err != nil {
					return nil, fmt.Errorf("%w: unable to encode %w for key %s", api.ErrBadRequest, err, key)
				}
				value.StringValue = &encodedStruct
			default:
				return nil, fmt.Errorf("%w: metadataType not found for %s: %v", api.ErrBadRequest, key, v)
			}

			props = append(props, value)
		}
	}

	return &props, nil
}

// MapRegisteredModelTypeIDEmbedMD maps RegisteredModel type id to embedmd one
func MapRegisteredModelTypeIDEmbedMD(source *OpenAPIModelWrapper[openapi.RegisteredModel]) (*int32, error) {
	return Int64ToInt32(&source.TypeId)
}

func encodeStruct(structValue *structpb.Struct) (string, error) {
	binaryData, err := proto.Marshal(structValue)
	if err != nil {
		return "", fmt.Errorf("failed to marshal proto struct: %w", err)
	}
	encodedString := base64.StdEncoding.EncodeToString(binaryData)
	return mlmdStructPrefix + encodedString, nil
}

func convertToStruct(source []string, key string) (*structpb.Struct, error) {
	list := &structpb.ListValue{
		Values: make([]*structpb.Value, len(source)),
	}
	for i, v := range source {
		list.Values[i] = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: v}}
	}
	return &structpb.Struct{Fields: map[string]*structpb.Value{key: {Kind: &structpb.Value_ListValue{ListValue: list}}}}, nil
}

// MapRegisteredModelPropertiesEmbedMD maps RegisteredModel fields to specific embedmd properties
func MapRegisteredModelPropertiesEmbedMD(source *openapi.RegisteredModel) (*[]models.Properties, error) {
	props := make([]models.Properties, 0)
	if source != nil {
		if source.Owner != nil {
			props = append(props, models.Properties{
				Name:             "owner",
				IsCustomProperty: false,
				StringValue:      source.Owner,
			})
		}

		if source.Description != nil {
			props = append(props, models.Properties{
				Name:             "description",
				IsCustomProperty: false,
				StringValue:      source.Description,
			})
		}

		if source.State != nil {
			props = append(props, models.Properties{
				Name:             "state",
				IsCustomProperty: false,
				StringValue:      of(string(*source.State)),
			})
		}

		if source.Language != nil {
			langStruct, err := convertToStruct(source.Language, "language")
			if err != nil {
				return nil, fmt.Errorf("%w: unable to convert to struct %w for key %s", api.ErrBadRequest, err, "language")
			}
			encodedString, err := encodeStruct(langStruct)
			if err != nil {
				return nil, fmt.Errorf("%w: unable to encode struct %w for key %s", api.ErrBadRequest, err, "language")
			}

			props = append(props, models.Properties{
				Name:             "language",
				IsCustomProperty: false,
				StringValue:      &encodedString,
			})
		}

		if source.LibraryName != nil {
			props = append(props, models.Properties{
				Name:             "library_name",
				IsCustomProperty: false,
				StringValue:      source.LibraryName,
			})
		}

		if source.License != nil {
			props = append(props, models.Properties{
				Name:             "license",
				IsCustomProperty: false,
				StringValue:      source.License,
			})
		}

		if source.LicenseLink != nil {
			props = append(props, models.Properties{
				Name:             "license_link",
				IsCustomProperty: false,
				StringValue:      source.LicenseLink,
			})
		}

		if source.Maturity != nil {
			props = append(props, models.Properties{
				Name:             "maturity",
				IsCustomProperty: false,
				StringValue:      source.Maturity,
			})
		}

		if source.Provider != nil {
			props = append(props, models.Properties{
				Name:             "provider",
				IsCustomProperty: false,
				StringValue:      source.Provider,
			})
		}

		if source.Readme != nil {
			props = append(props, models.Properties{
				Name:             "readme",
				IsCustomProperty: false,
				StringValue:      source.Readme,
			})
		}

		if source.Logo != nil {
			props = append(props, models.Properties{
				Name:             "logo",
				IsCustomProperty: false,
				StringValue:      source.Logo,
			})
		}

		if source.Tasks != nil {
			tasksStruct, err := convertToStruct(source.Tasks, "tasks")
			if err != nil {
				return nil, fmt.Errorf("%w: unable to convert to struct %w for key %s", api.ErrBadRequest, err, "tasks")
			}
			encodedString, err := encodeStruct(tasksStruct)
			if err != nil {
				return nil, fmt.Errorf("%w: unable to encode struct %w for key %s", api.ErrBadRequest, err, "tasks")
			}

			props = append(props, models.Properties{
				Name:             "tasks",
				IsCustomProperty: false,
				StringValue:      &encodedString,
			})
		}
	}

	return &props, nil
}

// MapRegisteredModelAttributesEmbedMD maps RegisteredModel attributes to specific embedmd properties
func MapRegisteredModelAttributesEmbedMD(source *openapi.RegisteredModel) (*models.RegisteredModelAttributes, error) {
	attributes := &models.RegisteredModelAttributes{}

	if source != nil {
		attributes.Name = &source.Name
		createdTime, err := StringToInt64(source.CreateTimeSinceEpoch)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, "createTimeSinceEpoch")
		}

		attributes.ExternalID = source.ExternalId

		attributes.CreateTimeSinceEpoch = createdTime

		lastUpdateTime, err := StringToInt64(source.LastUpdateTimeSinceEpoch)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, "lastUpdateTimeSinceEpoch")
		}

		attributes.LastUpdateTimeSinceEpoch = lastUpdateTime
	}

	return attributes, nil
}

// MapModelVersionTypeIDEmbedMD maps ModelVersion type id to embedmd one
func MapModelVersionTypeIDEmbedMD(source *OpenAPIModelWrapper[openapi.ModelVersion]) (*int32, error) {
	return Int64ToInt32(&source.TypeId)
}

// MapModelVersionPropertiesEmbedMD maps ModelVersion fields to specific embedmd properties
func MapModelVersionPropertiesEmbedMD(source *openapi.ModelVersion) (*[]models.Properties, error) {
	props := make([]models.Properties, 0)
	if source != nil {
		if source.Description != nil {
			props = append(props, models.Properties{
				Name:             "description",
				IsCustomProperty: false,
				StringValue:      source.Description,
			})
		}

		if source.State != nil {
			props = append(props, models.Properties{
				Name:             "state",
				IsCustomProperty: false,
				StringValue:      of(string(*source.State)),
			})
		}

		if source.Author != nil {
			props = append(props, models.Properties{
				Name:             "author",
				IsCustomProperty: false,
				StringValue:      source.Author,
			})
		}

		if source.RegisteredModelId != "" {
			registeredModelId, err := StringToInt32(source.RegisteredModelId)
			if err != nil {
				return nil, err
			}
			props = append(props, models.Properties{
				Name:             "registered_model_id",
				IsCustomProperty: false,
				IntValue:         &registeredModelId,
			})
		} else {
			return nil, fmt.Errorf("missing required RegisteredModelId field")
		}
	}

	return &props, nil
}

// MapModelVersionAttributesEmbedMD maps ModelVersion attributes to specific embedmd properties
func MapModelVersionAttributesEmbedMD(source *openapi.ModelVersion) (*models.ModelVersionAttributes, error) {
	attributes := &models.ModelVersionAttributes{}

	if source != nil {
		attributes.Name = &source.Name
		createdTime, err := StringToInt64(source.CreateTimeSinceEpoch)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, "createTimeSinceEpoch")
		}

		attributes.ExternalID = source.ExternalId

		attributes.CreateTimeSinceEpoch = createdTime

		lastUpdateTime, err := StringToInt64(source.LastUpdateTimeSinceEpoch)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, "lastUpdateTimeSinceEpoch")
		}

		attributes.LastUpdateTimeSinceEpoch = lastUpdateTime
	}

	return attributes, nil
}

// MapServingEnvironmentTypeIDEmbedMD maps ServingEnvironment type id to embedmd one
func MapServingEnvironmentTypeIDEmbedMD(source *OpenAPIModelWrapper[openapi.ServingEnvironment]) (*int32, error) {
	return Int64ToInt32(&source.TypeId)
}

// MapServingEnvironmentPropertiesEmbedMD maps ServingEnvironment fields to specific embedmd properties
func MapServingEnvironmentPropertiesEmbedMD(source *openapi.ServingEnvironment) (*[]models.Properties, error) {
	props := make([]models.Properties, 0)
	if source != nil {
		if source.Description != nil {
			props = append(props, models.Properties{
				Name:             "description",
				IsCustomProperty: false,
				StringValue:      source.Description,
			})
		}
	}

	return &props, nil
}

// MapServingEnvironmentAttributesEmbedMD maps ServingEnvironment attributes to specific embedmd properties
func MapServingEnvironmentAttributesEmbedMD(source *openapi.ServingEnvironment) (*models.ServingEnvironmentAttributes, error) {
	attributes := &models.ServingEnvironmentAttributes{}

	if source != nil {
		attributes.Name = &source.Name
		createdTime, err := StringToInt64(source.CreateTimeSinceEpoch)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, "createTimeSinceEpoch")
		}

		attributes.ExternalID = source.ExternalId

		attributes.CreateTimeSinceEpoch = createdTime

		lastUpdateTime, err := StringToInt64(source.LastUpdateTimeSinceEpoch)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, "lastUpdateTimeSinceEpoch")
		}

		attributes.LastUpdateTimeSinceEpoch = lastUpdateTime
	}

	return attributes, nil
}

func MapInferenceServiceTypeIDEmbedMD(source *OpenAPIModelWrapper[openapi.InferenceService]) (*int32, error) {
	return Int64ToInt32(&source.TypeId)
}

func MapInferenceServicePropertiesEmbedMD(source *openapi.InferenceService) (*[]models.Properties, error) {
	props := make([]models.Properties, 0)
	if source != nil {
		if source.Description != nil {
			props = append(props, models.Properties{
				Name:             "description",
				IsCustomProperty: false,
				StringValue:      source.Description,
			})
		}

		if source.Runtime != nil {
			props = append(props, models.Properties{
				Name:             "runtime",
				IsCustomProperty: false,
				StringValue:      source.Runtime,
			})
		}

		if source.DesiredState != nil {
			props = append(props, models.Properties{
				Name:             "desired_state",
				IsCustomProperty: false,
				StringValue:      of(string(*source.DesiredState)),
			})
		}

		if source.RegisteredModelId != "" {
			registeredModelId, err := StringToInt32(source.RegisteredModelId)
			if err != nil {
				return nil, err
			}
			props = append(props, models.Properties{
				Name:             "registered_model_id",
				IsCustomProperty: false,
				IntValue:         &registeredModelId,
			})
		} else {
			return nil, fmt.Errorf("missing required RegisteredModelId field")
		}

		if source.ServingEnvironmentId != "" {
			servingEnvironmentId, err := StringToInt32(source.ServingEnvironmentId)
			if err != nil {
				return nil, err
			}
			props = append(props, models.Properties{
				Name:             "serving_environment_id",
				IsCustomProperty: false,
				IntValue:         &servingEnvironmentId,
			})
		} else {
			return nil, fmt.Errorf("missing required ServingEnvironmentId field")
		}

		if source.ModelVersionId != nil {
			modelVersionId, err := StringToInt32(*source.ModelVersionId)
			if err != nil {
				return nil, err
			}
			props = append(props, models.Properties{
				Name:             "model_version_id",
				IsCustomProperty: false,
				IntValue:         &modelVersionId,
			})
		}
	}

	return &props, nil
}

func MapInferenceServiceAttributesEmbedMD(source *openapi.InferenceService) (*models.InferenceServiceAttributes, error) {
	attributes := &models.InferenceServiceAttributes{}

	if source != nil {
		attributes.Name = source.Name
		createdTime, err := StringToInt64(source.CreateTimeSinceEpoch)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, "createTimeSinceEpoch")
		}

		attributes.ExternalID = source.ExternalId

		attributes.CreateTimeSinceEpoch = createdTime

		lastUpdateTime, err := StringToInt64(source.LastUpdateTimeSinceEpoch)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, "lastUpdateTimeSinceEpoch")
		}

		attributes.LastUpdateTimeSinceEpoch = lastUpdateTime
	}

	return attributes, nil
}

func MapModelArtifactTypeIDEmbedMD(source *OpenAPIModelWrapper[openapi.ModelArtifact]) (*int32, error) {
	return Int64ToInt32(&source.TypeId)
}

func MapModelArtifactPropertiesEmbedMD(source *openapi.ModelArtifact) (*[]models.Properties, error) {
	props := make([]models.Properties, 0)
	if source != nil {
		if source.Description != nil {
			props = append(props, models.Properties{
				Name:             "description",
				IsCustomProperty: false,
				StringValue:      source.Description,
			})
		}
		if source.ModelFormatName != nil {
			props = append(props, models.Properties{
				Name:             "model_format_name",
				IsCustomProperty: false,
				StringValue:      source.ModelFormatName,
			})
		}
		if source.ModelFormatVersion != nil {
			props = append(props, models.Properties{
				Name:             "model_format_version",
				IsCustomProperty: false,
				StringValue:      source.ModelFormatVersion,
			})
		}
		if source.StorageKey != nil {
			props = append(props, models.Properties{
				Name:             "storage_key",
				IsCustomProperty: false,
				StringValue:      source.StorageKey,
			})
		}
		if source.StoragePath != nil {
			props = append(props, models.Properties{
				Name:             "storage_path",
				IsCustomProperty: false,
				StringValue:      source.StoragePath,
			})
		}
		if source.ServiceAccountName != nil {
			props = append(props, models.Properties{
				Name:             "service_account_name",
				IsCustomProperty: false,
				StringValue:      source.ServiceAccountName,
			})
		}
		if source.ModelSourceKind != nil {
			props = append(props, models.Properties{
				Name:             "model_source_kind",
				IsCustomProperty: false,
				StringValue:      source.ModelSourceKind,
			})
		}
		if source.ModelSourceClass != nil {
			props = append(props, models.Properties{
				Name:             "model_source_class",
				IsCustomProperty: false,
				StringValue:      source.ModelSourceClass,
			})
		}
		if source.ModelSourceGroup != nil {
			props = append(props, models.Properties{
				Name:             "model_source_group",
				IsCustomProperty: false,
				StringValue:      source.ModelSourceGroup,
			})
		}
		if source.ModelSourceId != nil {
			props = append(props, models.Properties{
				Name:             "model_source_id",
				IsCustomProperty: false,
				StringValue:      source.ModelSourceId,
			})
		}
		if source.ModelSourceName != nil {
			props = append(props, models.Properties{
				Name:             "model_source_name",
				IsCustomProperty: false,
				StringValue:      source.ModelSourceName,
			})
		}

	}

	return &props, nil
}

func MapModelArtifactAttributesEmbedMD(source *openapi.ModelArtifact) (*models.ModelArtifactAttributes, error) {
	attributes := &models.ModelArtifactAttributes{}

	if source != nil {
		attributes.Name = source.Name

		attributes.URI = source.Uri

		if source.State != nil {
			state, ok := models.Artifact_State_name[models.Artifact_State_value[string(*source.State)]]
			if !ok {
				return nil, fmt.Errorf("invalid state: %s", string(*source.State))
			}
			attributes.State = &state
		}

		createdTime, err := StringToInt64(source.CreateTimeSinceEpoch)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, "createTimeSinceEpoch")
		}

		attributes.ExternalID = source.ExternalId

		attributes.CreateTimeSinceEpoch = createdTime

		lastUpdateTime, err := StringToInt64(source.LastUpdateTimeSinceEpoch)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, "lastUpdateTimeSinceEpoch")
		}

		attributes.LastUpdateTimeSinceEpoch = lastUpdateTime
	}

	return attributes, nil
}

func MapDocArtifactTypeIDEmbedMD(source *OpenAPIModelWrapper[openapi.DocArtifact]) (*int32, error) {
	return Int64ToInt32(&source.TypeId)
}

func MapDocArtifactPropertiesEmbedMD(source *openapi.DocArtifact) (*[]models.Properties, error) {
	props := make([]models.Properties, 0)
	if source != nil {
		if source.Description != nil {
			props = append(props, models.Properties{
				Name:             "description",
				IsCustomProperty: false,
				StringValue:      source.Description,
			})
		}
	}

	return &props, nil
}

func MapDocArtifactAttributesEmbedMD(source *openapi.DocArtifact) (*models.DocArtifactAttributes, error) {
	attributes := &models.DocArtifactAttributes{}

	if source != nil {
		attributes.Name = source.Name

		attributes.URI = source.Uri

		if source.State != nil {
			state, ok := models.Artifact_State_name[models.Artifact_State_value[string(*source.State)]]
			if !ok {
				return nil, fmt.Errorf("invalid state: %s", string(*source.State))
			}
			attributes.State = &state
		}

		createdTime, err := StringToInt64(source.CreateTimeSinceEpoch)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, "createTimeSinceEpoch")
		}

		attributes.ExternalID = source.ExternalId

		attributes.CreateTimeSinceEpoch = createdTime

		lastUpdateTime, err := StringToInt64(source.LastUpdateTimeSinceEpoch)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, "lastUpdateTimeSinceEpoch")
		}

		attributes.LastUpdateTimeSinceEpoch = lastUpdateTime
	}

	return attributes, nil
}

func MapServeModelTypeIDEmbedMD(source *OpenAPIModelWrapper[openapi.ServeModel]) (*int32, error) {
	return Int64ToInt32(&source.TypeId)
}

func MapServeModelPropertiesEmbedMD(source *openapi.ServeModel) (*[]models.Properties, error) {
	props := make([]models.Properties, 0)

	if source != nil {
		if source.Description != nil {
			props = append(props, models.Properties{
				Name:             "description",
				IsCustomProperty: false,
				StringValue:      source.Description,
			})
		}

		if source.ModelVersionId != "" {
			modelVersionId, err := StringToInt32(source.ModelVersionId)
			if err != nil {
				return nil, err
			}
			props = append(props, models.Properties{
				Name:             "model_version_id",
				IsCustomProperty: false,
				IntValue:         &modelVersionId,
			})
		} else {
			return nil, fmt.Errorf("missing required ModelVersionId field")
		}
	}

	return &props, nil
}

func MapServeModelAttributesEmbedMD(source *openapi.ServeModel) (*models.ServeModelAttributes, error) {
	attributes := &models.ServeModelAttributes{}

	if source != nil {
		attributes.Name = source.Name

		if source.LastKnownState != nil {
			lastKnownState, ok := models.Execution_State_name[models.Execution_State_value[string(*source.LastKnownState)]]
			if !ok {
				return nil, fmt.Errorf("invalid last known state: %s", string(*source.LastKnownState))
			}
			attributes.LastKnownState = &lastKnownState
		}

		createdTime, err := StringToInt64(source.CreateTimeSinceEpoch)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, "createTimeSinceEpoch")
		}

		attributes.CreateTimeSinceEpoch = createdTime

		lastUpdateTime, err := StringToInt64(source.LastUpdateTimeSinceEpoch)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to decode as int64 %w for key %s", api.ErrBadRequest, err, "lastUpdateTimeSinceEpoch")
		}

		attributes.LastUpdateTimeSinceEpoch = lastUpdateTime

		attributes.ExternalID = source.ExternalId
	}

	return attributes, nil
}
