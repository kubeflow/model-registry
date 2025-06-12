package converter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// MapEmbedMDCustomProperties maps EmbedMD custom properties model to OpenAPI one
func MapEmbedMDCustomProperties(source []models.Properties) (map[string]openapi.MetadataValue, error) {
	data := make(map[string]openapi.MetadataValue)

	for _, v := range source {
		customValue := openapi.MetadataValue{}

		if v.IntValue != nil {
			customValue.MetadataIntValue = NewMetadataIntValue(strconv.FormatInt(int64(*v.IntValue), 10))
		} else if v.StringValue != nil {
			customValue.MetadataStringValue = NewMetadataStringValue(*v.StringValue)
		} else if v.BoolValue != nil {
			customValue.MetadataBoolValue = NewMetadataBoolValue(*v.BoolValue)
		} else if v.DoubleValue != nil {
			customValue.MetadataDoubleValue = NewMetadataDoubleValue(*v.DoubleValue)
		} else if v.ByteValue != nil {
			asJSON, err := json.Marshal(v.ByteValue)
			if err != nil {
				return nil, err
			}
			b64 := base64.StdEncoding.EncodeToString(asJSON)
			customValue.MetadataStructValue = NewMetadataStructValue(b64)
		} else {
			return nil, fmt.Errorf("%w: metadataType not found for %s: %v", api.ErrBadRequest, v.Name, v)
		}

		data[v.Name] = customValue
	}

	return data, nil
}

func MapEmbedMDDescription(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "description" {
			return v.StringValue
		}
	}

	return nil
}

func MapEmbedMDOwner(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "owner" {
			return v.StringValue
		}
	}

	return nil
}

func MapEmbedMDAuthor(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "author" {
			return v.StringValue
		}
	}

	return nil
}

func MapEmbedMDStateRegisteredModel(source *[]models.Properties) (*openapi.RegisteredModelState, error) {
	for _, v := range *source {
		if v.Name == "state" {
			if v.StringValue == nil {
				return nil, fmt.Errorf("%w: state is required", api.ErrBadRequest)
			}

			registeredModelState, err := openapi.NewRegisteredModelStateFromValue(*v.StringValue)
			if err != nil {
				return nil, err
			}

			return registeredModelState, nil
		}
	}

	return nil, nil
}

func MapEmbedMDStateModelVersion(source *[]models.Properties) (*openapi.ModelVersionState, error) {
	for _, v := range *source {
		if v.Name == "state" {
			if v.StringValue == nil {
				return nil, fmt.Errorf("%w: state is required", api.ErrBadRequest)
			}

			modelVersionState, err := openapi.NewModelVersionStateFromValue(*v.StringValue)
			if err != nil {
				return nil, err
			}

			return modelVersionState, nil
		}
	}

	return nil, nil
}

func MapEmbedMDExternalIDRegisteredModel(source *models.RegisteredModelAttributes) *string {
	return source.ExternalID
}

func MapEmbedMDNameRegisteredModel(source *models.RegisteredModelAttributes) string {
	if source.Name == nil {
		return ""
	}
	return *source.Name
}

func MapEmbedMDCreateTimeSinceEpochRegisteredModel(source *models.RegisteredModelAttributes) *string {
	return Int64ToString(source.CreateTimeSinceEpoch)
}

func MapEmbedMDLastUpdateTimeSinceEpochRegisteredModel(source *models.RegisteredModelAttributes) *string {
	return Int64ToString(source.LastUpdateTimeSinceEpoch)
}

func MapEmbedMDNameModelVersion(source *models.ModelVersionAttributes) string {
	return *MapNameFromOwned(source.Name)
}

func MapEmbedMDExternalIDModelVersion(source *models.ModelVersionAttributes) *string {
	return source.ExternalID
}

func MapEmbedMDCreateTimeSinceEpochModelVersion(source *models.ModelVersionAttributes) *string {
	return Int64ToString(source.CreateTimeSinceEpoch)
}

func MapEmbedMDLastUpdateTimeSinceEpochModelVersion(source *models.ModelVersionAttributes) *string {
	return Int64ToString(source.LastUpdateTimeSinceEpoch)
}

func MapEmbedMDExternalIDServingEnvironment(source *models.ServingEnvironmentAttributes) *string {
	return source.ExternalID
}

func MapEmbedMDNameServingEnvironment(source *models.ServingEnvironmentAttributes) string {
	if source.Name == nil {
		return ""
	}
	return *source.Name
}

func MapEmbedMDCreateTimeSinceEpochServingEnvironment(source *models.ServingEnvironmentAttributes) *string {
	return Int64ToString(source.CreateTimeSinceEpoch)
}

func MapEmbedMDLastUpdateTimeSinceEpochServingEnvironment(source *models.ServingEnvironmentAttributes) *string {
	return Int64ToString(source.LastUpdateTimeSinceEpoch)
}

func MapEmbedMDPropertyRuntime(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "runtime" {
			return v.StringValue
		}
	}

	return nil
}

func MapEmbedMDExternalIDInferenceService(source *models.InferenceServiceAttributes) *string {
	return source.ExternalID
}

func MapEmbedMDPropertyDesiredStateInferenceService(source *[]models.Properties) (*openapi.InferenceServiceState, error) {
	for _, v := range *source {
		if v.Name == "desired_state" {
			if v.StringValue == nil {
				return nil, fmt.Errorf("%w: desired_state is required", api.ErrBadRequest)
			}

			inferenceServiceState, err := openapi.NewInferenceServiceStateFromValue(*v.StringValue)
			if err != nil {
				return nil, err
			}

			return inferenceServiceState, nil
		}
	}

	return nil, nil
}

func MapEmbedMDPropertyModelVersionId(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "model_version_id" {
			return Int32ToString(v.IntValue)
		}
	}

	return nil
}

func MapEmbedMDPropertyRegisteredModelId(source *[]models.Properties) string {
	for _, v := range *source {
		if v.Name == "registered_model_id" {
			result := Int32ToString(v.IntValue)
			if result == nil {
				return ""
			}
			return *result
		}
	}

	return ""
}

func MapEmbedMDPropertyServingEnvironmentId(source *[]models.Properties) string {
	for _, v := range *source {
		if v.Name == "serving_environment_id" {
			result := Int32ToString(v.IntValue)
			if result == nil {
				return ""
			}
			return *result
		}
	}

	return ""
}

func MapEmbedMDNameInferenceService(source *models.InferenceServiceAttributes) *string {
	return MapNameFromOwned(source.Name)
}

func MapEmbedMDCreateTimeSinceEpochInferenceService(source *models.InferenceServiceAttributes) *string {
	return Int64ToString(source.CreateTimeSinceEpoch)
}

func MapEmbedMDLastUpdateTimeSinceEpochInferenceService(source *models.InferenceServiceAttributes) *string {
	return Int64ToString(source.LastUpdateTimeSinceEpoch)
}

func MapEmbedMDNameModelArtifact(source *models.ModelArtifactAttributes) *string {
	return MapNameFromOwned(source.Name)
}

func MapEmbedMDURIModelArtifact(source *models.ModelArtifactAttributes) *string {
	return source.URI
}

func MapEmbedMDArtifactTypeModelArtifact(source *models.ModelArtifactAttributes) *string {
	return of("model-artifact")
}

func MapEmbedMDPropertyModelFormatName(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "model_format_name" {
			return v.StringValue
		}
	}

	return nil
}

func MapEmbedMDPropertyModelFormatVersion(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "model_format_version" {
			return v.StringValue
		}
	}

	return nil
}

func MapEmbedMDPropertyStorageKey(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "storage_key" {
			return v.StringValue
		}
	}

	return nil
}

func MapEmbedMDPropertyStoragePath(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "storage_path" {
			return v.StringValue
		}
	}

	return nil
}

func MapEmbedMDPropertyServiceAccountName(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "service_account_name" {
			return v.StringValue
		}
	}

	return nil
}

func MapEmbedMDPropertyModelSourceKind(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "model_source_kind" {
			return v.StringValue
		}
	}

	return nil
}

func MapEmbedMDPropertyModelSourceClass(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "model_source_class" {
			return v.StringValue
		}
	}

	return nil
}

func MapEmbedMDPropertyModelSourceGroup(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "model_source_group" {
			return v.StringValue
		}
	}

	return nil
}

func MapEmbedMDPropertyModelSourceId(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "model_source_id" {
			return v.StringValue
		}
	}

	return nil
}

func MapEmbedMDPropertyModelSourceName(source *[]models.Properties) *string {
	for _, v := range *source {
		if v.Name == "model_source_name" {
			return v.StringValue
		}
	}

	return nil
}

func MapEmbedMDExternalIDModelArtifact(source *models.ModelArtifactAttributes) *string {
	return source.ExternalID
}

func MapEmbedMDCreateTimeSinceEpochModelArtifact(source *models.ModelArtifactAttributes) *string {
	return Int64ToString(source.CreateTimeSinceEpoch)
}

func MapEmbedMDLastUpdateTimeSinceEpochModelArtifact(source *models.ModelArtifactAttributes) *string {
	return Int64ToString(source.LastUpdateTimeSinceEpoch)
}

func MapEmbedMDStateModelArtifact(source *models.ModelArtifactAttributes) (*openapi.ArtifactState, error) {
	defaultState := openapi.ARTIFACTSTATE_UNKNOWN

	if source.State == nil {
		return &defaultState, nil
	}

	return openapi.NewArtifactStateFromValue(*source.State)
}

func MapEmbedMDURIDocArtifact(source *models.DocArtifactAttributes) *string {
	return source.URI
}

func MapEmbedMDArtifactTypeDocArtifact(source *models.DocArtifactAttributes) *string {
	return of("doc-artifact")
}

func MapEmbedMDExternalIDDocArtifact(source *models.DocArtifactAttributes) *string {
	return source.ExternalID
}

func MapEmbedMDNameDocArtifact(source *models.DocArtifactAttributes) *string {
	return MapNameFromOwned(source.Name)
}

func MapEmbedMDCreateTimeSinceEpochDocArtifact(source *models.DocArtifactAttributes) *string {
	return Int64ToString(source.CreateTimeSinceEpoch)
}

func MapEmbedMDLastUpdateTimeSinceEpochDocArtifact(source *models.DocArtifactAttributes) *string {
	return Int64ToString(source.LastUpdateTimeSinceEpoch)
}

func MapEmbedMDStateDocArtifact(source *models.DocArtifactAttributes) (*openapi.ArtifactState, error) {
	defaultState := openapi.ARTIFACTSTATE_UNKNOWN

	if source.State == nil {
		return &defaultState, nil
	}

	return openapi.NewArtifactStateFromValue(*source.State)
}

func MapEmbedMDExternalIDServeModel(source *models.ServeModelAttributes) *string {
	return source.ExternalID
}

func MapEmbedMDNameServeModel(source *models.ServeModelAttributes) *string {
	return MapNameFromOwned(source.Name)
}

func MapEmbedMDLastKnownStateServeModel(source *models.ServeModelAttributes) (*openapi.ExecutionState, error) {
	defaultState := openapi.EXECUTIONSTATE_UNKNOWN

	if source.LastKnownState == nil {
		return &defaultState, nil
	}

	return openapi.NewExecutionStateFromValue(*source.LastKnownState)
}

func MapEmbedMDCreateTimeSinceEpochServeModel(source *models.ServeModelAttributes) *string {
	return Int64ToString(source.CreateTimeSinceEpoch)
}

func MapEmbedMDLastUpdateTimeSinceEpochServeModel(source *models.ServeModelAttributes) *string {
	return Int64ToString(source.LastUpdateTimeSinceEpoch)
}

func MapEmbedMDPropertyModelVersionIdServeModel(source *[]models.Properties) (string, error) {
	modelVersionId := MapEmbedMDPropertyModelVersionId(source)

	if modelVersionId == nil {
		return "", fmt.Errorf("model version id is required")
	}

	return *modelVersionId, nil
}
