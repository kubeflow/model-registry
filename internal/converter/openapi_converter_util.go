package converter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

type OpenAPIModel interface {
	openapi.Artifact |
		openapi.RegisteredModel |
		openapi.ModelVersion |
		openapi.ModelArtifact |
		openapi.DocArtifact |
		openapi.DataSet |
		openapi.Metric |
		openapi.Parameter |
		openapi.ServingEnvironment |
		openapi.InferenceService |
		openapi.ServeModel |
		openapi.Experiment |
		openapi.ExperimentRun
}

type OpenapiUpdateWrapper[
	M OpenAPIModel,
] struct {
	Existing *M
	Update   *M
}

func NewOpenapiUpdateWrapper[
	M OpenAPIModel,
](existing *M, update *M) OpenapiUpdateWrapper[M] {
	return OpenapiUpdateWrapper[M]{
		Existing: existing,
		Update:   update,
	}
}

func InitWithExisting[M OpenAPIModel](source OpenapiUpdateWrapper[M]) M {
	return *source.Existing
}

func InitWithUpdate[M OpenAPIModel](source OpenapiUpdateWrapper[M]) M {
	if source.Update != nil {
		return *source.Update
	}
	var m M
	return m
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

	if entityName == "" {
		return prefix
	}

	return fmt.Sprintf("%s:%s", prefix, entityName)
}

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

// Int64ToString converts numeric id to string-based one
func Int64ToString(id *int64) *string {
	if id == nil {
		return nil
	}

	idAsString := strconv.FormatInt(*id, 10)
	return &idAsString
}

// StringToInt64 converts string-based id to int64 if numeric, otherwise return error
func StringToInt64(id *string) (*int64, error) {
	if id == nil {
		return nil, nil
	}

	idAsInt64, err := strconv.ParseInt(*id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid numeric string: %v", err)
	}

	return &idAsInt64, nil
}

// StringToInt32 converts string-based numeric value (a OpenAPI string literal consisting only of digits) to int32 if numeric, otherwise return error
func StringToInt32(idString string) (int32, error) {
	idAsInt32, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric string: %v", err)
	}

	return int32(idAsInt32), nil
}

// ValidateStepIds validates and parses a comma-separated string of step IDs
// Returns error if any step ID is not a valid integer
func ValidateStepIds(stepIds string) error {
	if stepIds == "" {
		return nil
	}

	parts := strings.Split(stepIds, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue // skip empty parts
		}

		_, err := StringToInt32(part)
		if err != nil {
			return fmt.Errorf("invalid step ID '%s': must be a valid integer", part)
		}
	}

	return nil
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
