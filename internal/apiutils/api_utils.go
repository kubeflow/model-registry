package apiutils

import (
	"github.com/opendatahub-io/model-registry/internal/converter"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/pkg/api"
	model "github.com/opendatahub-io/model-registry/pkg/openapi"
)

func BuildListOperationOptions(listOptions api.ListOptions) (*proto.ListOperationOptions, error) {

	result := proto.ListOperationOptions{}
	if listOptions.PageSize != nil {
		result.MaxResultSize = listOptions.PageSize
	}
	if listOptions.NextPageToken != nil {
		result.NextPageToken = listOptions.NextPageToken
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

func BuildListOption(pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (api.ListOptions, error) {
	var pageSizeInt32 *int32
	if pageSize != "" {
		conv, err := converter.StringToInt32(pageSize)
		if err != nil {
			return api.ListOptions{}, err
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
		PageSize:      pageSizeInt32,
		OrderBy:       orderByString,
		SortOrder:     sortOrderString,
		NextPageToken: nextPageTokenParam,
	}, nil
}
