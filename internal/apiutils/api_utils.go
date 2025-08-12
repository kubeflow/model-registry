package apiutils

import (
	"fmt"
	"strconv"

	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/api"
	model "github.com/kubeflow/model-registry/pkg/openapi"
)

func BuildListOperationOptions(listOptions api.ListOptions) (*proto.ListOperationOptions, error) {
	result := proto.ListOperationOptions{}
	if listOptions.PageSize != nil {
		result.MaxResultSize = listOptions.PageSize
	}
	if listOptions.NextPageToken != nil {
		result.NextPageToken = listOptions.NextPageToken
	}
	if listOptions.FilterQuery != nil {
		result.FilterQuery = listOptions.FilterQuery
	}
	if listOptions.OrderBy != nil {
		so := listOptions.SortOrder

		// default is DESC
		isAsc := false
		if so != nil && *so == "ASC" {
			isAsc = true
		}

		var orderByField proto.ListOperationOptions_OrderByField_Field
		if val, ok := proto.ListOperationOptions_OrderByField_Field_value[*listOptions.OrderBy]; ok {
			orderByField = proto.ListOperationOptions_OrderByField_Field(val)
		}

		result.OrderByField = &proto.ListOperationOptions_OrderByField{
			Field: &orderByField,
			IsAsc: &isAsc,
		}
	}
	return &result, nil
}

// ZeroIfNil return the zeroed value if input is a nil pointer
func ZeroIfNil[T any](input *T) T {
	if input != nil {
		return *input
	}
	return *new(T)
}

// of returns a pointer to the provided literal/const input
func Of[E any](e E) *E {
	return &e
}

func StrPtr(notEmpty string) *string {
	if notEmpty == "" {
		return nil
	}
	return &notEmpty
}

func BuildListOption(filterQuery string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (api.ListOptions, error) {
	var filterQueryPtr *string
	if filterQuery != "" {
		filterQueryPtr = &filterQuery
	}
	var pageSizeInt32 *int32
	if pageSize != "" {
		conv, err := converter.StringToInt32(pageSize)
		if err != nil {
			return api.ListOptions{}, fmt.Errorf("%v: %w", err, api.ErrBadRequest)
		}
		pageSizeInt32 = &conv
	}
	var orderByString *string
	if orderBy != "" {
		orderByString = (*string)(&orderBy)
	}
	var sortOrderString *string
	if sortOrder != "" {
		sortOrderString = (*string)(&sortOrder)
	}
	var nextPageTokenParam *string
	if nextPageToken != "" {
		nextPageTokenParam = &nextPageToken
	}
	return api.ListOptions{
		FilterQuery:   filterQueryPtr,
		PageSize:      pageSizeInt32,
		OrderBy:       orderByString,
		SortOrder:     sortOrderString,
		NextPageToken: nextPageTokenParam,
	}, nil
}

// BuildListOptionWithFilterTranslation builds list options and translates filter queries from REST API to MLMD format
func BuildListOptionWithFilterTranslation(filterQuery string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string, entityType EntityType) (api.ListOptions, error) {
	// First build the basic list options
	listOptions, err := BuildListOption(filterQuery, pageSize, orderBy, sortOrder, nextPageToken)
	if err != nil {
		return api.ListOptions{}, err
	}

	// Translate the filter query if present
	if listOptions.FilterQuery != nil && *listOptions.FilterQuery != "" {
		translatedQuery, err := TranslateFilterQuery(*listOptions.FilterQuery, entityType)
		if err != nil {
			return api.ListOptions{}, fmt.Errorf("filter query translation failed: %v: %w", err, api.ErrBadRequest)
		}
		listOptions.FilterQuery = &translatedQuery
	}

	return listOptions, nil
}

// ValidateIDAsInt32 validates and converts a string ID to int32
// Returns an error with api.ErrBadRequest if the ID is invalid
func ValidateIDAsInt32(id string, entityName string) (int32, error) {
	convertedId, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s ID: %v: %w", entityName, err, api.ErrBadRequest)
	}
	return int32(convertedId), nil
}

// ValidateIDAsInt32Ptr validates and converts a string pointer ID to int32 pointer
// Returns nil if the input is nil, otherwise validates and converts
func ValidateIDAsInt32Ptr(id *string, entityName string) (*int32, error) {
	if id == nil {
		return nil, nil
	}
	convertedId, err := strconv.ParseInt(*id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid %s ID: %v: %w", entityName, err, api.ErrBadRequest)
	}
	result := int32(convertedId)
	return &result, nil
}

// ValidateIDAsInt64 validates and converts a string ID to int64
// Returns an error with api.ErrBadRequest if the ID is invalid
func ValidateIDAsInt64(id string, entityName string) (int64, error) {
	convertedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s ID: %v: %w", entityName, err, api.ErrBadRequest)
	}
	return convertedId, nil
}

// ValidateIDAsInt64Ptr validates and converts a string pointer ID to int64 pointer
// Returns nil if the input is nil, otherwise validates and converts
func ValidateIDAsInt64Ptr(id *string, entityName string) (*int64, error) {
	if id == nil {
		return nil, nil
	}
	convertedId, err := strconv.ParseInt(*id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid %s ID: %v: %w", entityName, err, api.ErrBadRequest)
	}
	return &convertedId, nil
}
