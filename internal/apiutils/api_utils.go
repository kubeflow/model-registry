package apiutils

import (
	"fmt"

	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/pkg/api"
	model "github.com/kubeflow/model-registry/pkg/openapi"
)

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

func BuildListOption(pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (api.ListOptions, error) {
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
		PageSize:      pageSizeInt32,
		OrderBy:       orderByString,
		SortOrder:     sortOrderString,
		NextPageToken: nextPageTokenParam,
	}, nil
}
