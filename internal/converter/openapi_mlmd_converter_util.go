package converter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/openapi"
	"google.golang.org/protobuf/types/known/structpb"
)

// StringToInt64 converts string-based id to int64 if numeric, otherwise return error
func StringToInt64(id *string) (*int64, error) {
	if id == nil {
		return nil, nil
	}

	idAsInt, err := strconv.Atoi(*id)
	if err != nil {
		return nil, err
	}

	idInt64 := int64(idAsInt)
	return &idInt64, nil
}

// Int64ToString converts numeric id to string-based one
func Int64ToString(id *int64) *string {
	if id == nil {
		return nil
	}

	idAsString := strconv.FormatInt(*id, 10)
	return &idAsString
}

// StringToInt32 converts string-based numeric value (a OpenAPI string literal consisting only of digits) to int32 if numeric, otherwise return error
func StringToInt32(idString string) (int32, error) {
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return 0, err
	}

	idInt32 := int32(idInt)
	return idInt32, nil
}

// MapOpenAPICustomProperties maps OpenAPI custom properties model to MLMD one
func MapOpenAPICustomProperties(source *map[string]openapi.MetadataValue) (map[string]*proto.Value, error) {
	props := make(map[string]*proto.Value)

	if source != nil {
		for key, v := range *source {
			value := proto.Value{}

			switch {
			// bool value
			case v.MetadataBoolValue != nil:
				value.Value = &proto.Value_BoolValue{BoolValue: *v.MetadataBoolValue.BoolValue}
			// int value
			case v.MetadataIntValue != nil:
				intValue, err := StringToInt64(v.MetadataIntValue.IntValue)
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
	}

	return props, nil
}

// PrefixWhenOwned compose the mlmd fullname by using ownerId as prefix
// For owned entity such as ModelVersion
// for potentially owned entity such as ModelArtifact
func PrefixWhenOwned(ownerId *string, entityName string) string {
	var prefix string
	if ownerId != nil {
		prefix = *ownerId
	} else {
		prefix = uuid.New().String()
	}
	prefixedName := fmt.Sprintf("%s:%s", prefix, entityName)
	return prefixedName
}

// REGISTERED MODEL

// MapRegisteredModelProperties maps RegisteredModel fields to specific MLMD properties
func MapRegisteredModelProperties(source *openapi.RegisteredModel) (map[string]*proto.Value, error) {
	props := make(map[string]*proto.Value)
	if source != nil {
		if source.Description != nil {
			props["description"] = &proto.Value{
				Value: &proto.Value_StringValue{
					StringValue: *source.Description,
				},
			}
		}
	}
	return props, nil
}

// MapRegisteredModelType return RegisteredModel corresponding MLMD context type
func MapRegisteredModelType(_ *openapi.RegisteredModel) *string {
	return of(RegisteredModelTypeName)
}

// MODEL VERSION

// MapModelVersionProperties maps ModelVersion fields to specific MLMD properties
func MapModelVersionProperties(source *OpenAPIModelWrapper[openapi.ModelVersion]) (map[string]*proto.Value, error) {
	props := make(map[string]*proto.Value)
	if source != nil {
		if (*source.Model).Description != nil {
			props["description"] = &proto.Value{
				Value: &proto.Value_StringValue{
					StringValue: *(*source.Model).Description,
				},
			}
		}
		if (*source).ModelName != nil {
			props["model_name"] = &proto.Value{
				Value: &proto.Value_StringValue{
					StringValue: *(*source).ModelName,
				},
			}
		}
		props["version"] = &proto.Value{
			Value: &proto.Value_StringValue{
				StringValue: *(*source.Model).Name,
			},
		}
		// TODO: not available for now
		props["author"] = &proto.Value{
			Value: &proto.Value_StringValue{
				StringValue: "",
			},
		}
	}
	return props, nil
}

// MapModelVersionType return ModelVersion corresponding MLMD context type
func MapModelVersionType(_ *openapi.ModelVersion) *string {
	return of(ModelVersionTypeName)
}

// MapModelVersionName maps the user-provided name into MLMD one, i.e., prefixing it with
// either the parent resource id or a generated uuid
func MapModelVersionName(source *OpenAPIModelWrapper[openapi.ModelVersion]) *string {
	return of(PrefixWhenOwned(source.ParentResourceId, *(*source).Model.Name))
}

// MODEL ARTIFACT

// MapModelArtifactProperties maps ModelArtifact fields to specific MLMD properties
func MapModelArtifactProperties(source *openapi.ModelArtifact) (map[string]*proto.Value, error) {
	props := make(map[string]*proto.Value)
	if source != nil {
		if source.Description != nil {
			props["description"] = &proto.Value{
				Value: &proto.Value_StringValue{
					StringValue: *source.Description,
				},
			}
		}
		if source.Runtime != nil {
			props["runtime"] = &proto.Value{
				Value: &proto.Value_StringValue{
					StringValue: *source.Runtime,
				},
			}
		}
		if source.ModelFormatName != nil {
			props["model_format_name"] = &proto.Value{
				Value: &proto.Value_StringValue{
					StringValue: *source.ModelFormatName,
				},
			}
		}
		if source.ModelFormatVersion != nil {
			props["model_format_version"] = &proto.Value{
				Value: &proto.Value_StringValue{
					StringValue: *source.ModelFormatVersion,
				},
			}
		}
		if source.StorageKey != nil {
			props["storage_key"] = &proto.Value{
				Value: &proto.Value_StringValue{
					StringValue: *source.StorageKey,
				},
			}
		}
		if source.StoragePath != nil {
			props["storage_path"] = &proto.Value{
				Value: &proto.Value_StringValue{
					StringValue: *source.StoragePath,
				},
			}
		}
		if source.ServiceAccountName != nil {
			props["service_account_name"] = &proto.Value{
				Value: &proto.Value_StringValue{
					StringValue: *source.ServiceAccountName,
				},
			}
		}
	}
	return props, nil
}

// MapModelArtifactType return ModelArtifact corresponding MLMD context type
func MapModelArtifactType(_ *openapi.ModelArtifact) *string {
	return of(ModelArtifactTypeName)
}

// MapModelArtifactName maps the user-provided name into MLMD one, i.e., prefixing it with
// either the parent resource id or a generated uuid. If not provided, autogenerate the name
// itself
func MapModelArtifactName(source *OpenAPIModelWrapper[openapi.ModelArtifact]) *string {
	// openapi.Artifact is defined with optional name, so build arbitrary name for this artifact if missing
	var artifactName string
	if (*source).Model.Name != nil {
		artifactName = *(*source).Model.Name
	} else {
		artifactName = uuid.New().String()
	}
	return of(PrefixWhenOwned(source.ParentResourceId, artifactName))
}

func MapOpenAPIModelArtifactState(source *openapi.ArtifactState) (*proto.Artifact_State, error) {
	if source == nil {
		return nil, nil
	}

	val, ok := proto.Artifact_State_value[string(*source)]
	if !ok {
		return nil, fmt.Errorf("invalid artifact state: %s", string(*source))
	}

	return (*proto.Artifact_State)(&val), nil
}

// SERVING ENVIRONMENT

// MapServingEnvironmentType return ServingEnvironment corresponding MLMD context type
func MapServingEnvironmentType(_ *openapi.ServingEnvironment) *string {
	return of(ServingEnvironmentTypeName)
}

// MapServingEnvironmentProperties maps ServingEnvironment fields to specific MLMD properties
func MapServingEnvironmentProperties(source *openapi.ServingEnvironment) (map[string]*proto.Value, error) {
	props := make(map[string]*proto.Value)
	if source != nil {
		if source.Description != nil {
			props["description"] = &proto.Value{
				Value: &proto.Value_StringValue{
					StringValue: *source.Description,
				},
			}
		}
	}
	return props, nil
}

// INFERENCE SERVICE

// MapInferenceServiceType return InferenceService corresponding MLMD context type
func MapInferenceServiceType(_ *openapi.InferenceService) *string {
	return of(InferenceServiceTypeName)
}

// MapInferenceServiceProperties maps InferenceService fields to specific MLMD properties
func MapInferenceServiceProperties(source *openapi.InferenceService) (map[string]*proto.Value, error) {
	props := make(map[string]*proto.Value)
	if source != nil {
		if source.Description != nil {
			props["description"] = &proto.Value{
				Value: &proto.Value_StringValue{
					StringValue: *source.Description,
				},
			}
		}

		registeredModelId, err := StringToInt64(&source.RegisteredModelId)
		if err != nil {
			return nil, err
		}
		props["registered_model_id"] = &proto.Value{
			Value: &proto.Value_IntValue{
				IntValue: *registeredModelId,
			},
		}

		servingEnvironmentId, err := StringToInt64(&source.ServingEnvironmentId)
		if err != nil {
			return nil, err
		}
		props["serving_environment_id"] = &proto.Value{
			Value: &proto.Value_IntValue{
				IntValue: *servingEnvironmentId,
			},
		}

		if source.ModelVersionId != nil {
			modelVersionId, err := StringToInt64(source.ModelVersionId)
			if err != nil {
				return nil, err
			}
			props["model_version_id"] = &proto.Value{
				Value: &proto.Value_IntValue{
					IntValue: *modelVersionId,
				},
			}
		}

	}
	return props, nil
}

// MapInferenceServiceName maps the user-provided name into MLMD one, i.e., prefixing it with
// either the parent resource id or a generated uuid
// ref: > InferenceService context is actually a child of ServingEnvironment parent context
func MapInferenceServiceName(source *OpenAPIModelWrapper[openapi.InferenceService]) *string {
	return of(PrefixWhenOwned(source.ParentResourceId, *(*source).Model.Name))
}

// SERVE MODEL

// MapServeModelType return ServeModel corresponding MLMD context type
func MapServeModelType(_ *openapi.ServeModel) *string {
	return of(ServeModelTypeName)
}

// MapServeModelProperties maps ServeModel fields to specific MLMD properties
func MapServeModelProperties(source *openapi.ServeModel) (map[string]*proto.Value, error) {
	props := make(map[string]*proto.Value)
	if source != nil {
		if source.Description != nil {
			props["description"] = &proto.Value{
				Value: &proto.Value_StringValue{
					StringValue: *source.Description,
				},
			}
		}

		modelVersionId, err := StringToInt64(&source.ModelVersionId)
		if err != nil {
			return nil, err
		}
		props["model_version_id"] = &proto.Value{
			Value: &proto.Value_IntValue{
				IntValue: *modelVersionId,
			},
		}

	}
	return props, nil
}

// MapServeModelName maps the user-provided name into MLMD one, i.e., prefixing it with
// either the parent resource id or a generated uuid. If not provided, autogenerate the name
// itself
func MapServeModelName(source *OpenAPIModelWrapper[openapi.ServeModel]) *string {
	// openapi.ServeModel is defined with optional name, so build arbitrary name for this artifact if missing
	var serveModelName string
	if (*source).Model.Name != nil {
		serveModelName = *(*source).Model.Name
	} else {
		serveModelName = uuid.New().String()
	}
	return of(PrefixWhenOwned(source.ParentResourceId, serveModelName))
}

// MapLastKnownState maps LastKnownState field from ServeModel to Execution
func MapLastKnownState(source *openapi.ExecutionState) (*proto.Execution_State, error) {
	if source == nil {
		return nil, nil
	}

	val, ok := proto.Execution_State_value[string(*source)]
	if !ok {
		return nil, fmt.Errorf("invalid execution state: %s", string(*source))
	}

	return (*proto.Execution_State)(&val), nil
}

// of returns a pointer to the provided literal/const input
func of[E any](e E) *E {
	return &e
}
