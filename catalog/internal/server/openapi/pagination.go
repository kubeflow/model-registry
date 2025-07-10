package openapi

import (
	"encoding/base64"
	"fmt"
	"strings"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/db/models"
)

func buildPagination(pageSize string, orderBy string, sortOrder string, nextPageToken string) (*models.Pagination, error) {
	var pageSizeInt32 *int32
	if pageSize != "" {
		conv, err := converter.StringToInt32(pageSize)
		if err != nil {
			return nil, fmt.Errorf("error converting page size to int32: %w", err)
		}
		pageSizeInt32 = &conv
	} else {
		defaultPageSize := int32(10)
		pageSizeInt32 = &defaultPageSize
	}

	var orderByString *string
	if orderBy != "" {
		orderByString = &orderBy
	}

	var sortOrderString *string
	if sortOrder != "" {
		sortOrderString = &sortOrder
	}

	var nextPageTokenParam *string
	if nextPageToken != "" {
		nextPageTokenParam = &nextPageToken
	}

	return &models.Pagination{
		PageSize:      pageSizeInt32,
		OrderBy:       orderByString,
		SortOrder:     sortOrderString,
		NextPageToken: nextPageTokenParam,
	}, nil
}

type stringCursor struct {
	Value string
	ID    string
}

func encodeStringCursor(c stringCursor) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.Value, c.ID)))
}

func decodeStringCursor(encoded string) (stringCursor, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return stringCursor{}, err
	}
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return stringCursor{}, fmt.Errorf("invalid cursor format")
	}
	return stringCursor{
		Value: parts[0],
		ID:    parts[1],
	}, nil
}

func paginateSources(items []model.CatalogSource, pagination *models.Pagination) ([]model.CatalogSource, string) {
	startIndex := 0
	if pagination.GetNextPageToken() != "" {
		cursor, err := decodeStringCursor(pagination.GetNextPageToken())
		if err == nil {
			for i, item := range items {
				var itemValue string
				switch model.OrderByField(strings.ToUpper(pagination.GetOrderBy())) {
				case model.ORDERBYFIELD_ID, "":
					if item.Id != "" {
						itemValue = item.Id
					}
				case model.ORDERBYFIELD_NAME:
					if item.Name != "" {
						itemValue = item.Name
					}
				}
				if item.Id != "" && item.Id == cursor.ID && itemValue == cursor.Value {
					startIndex = i + 1
					break
				}
			}
		}
	}

	var pagedItems []model.CatalogSource
	var newNextPageToken string

	if startIndex < len(items) {
		limit := int(pagination.GetPageSize())
		endIndex := startIndex + limit
		if endIndex > len(items) {
			endIndex = len(items)
		}
		pagedItems = items[startIndex:endIndex]

		if endIndex < len(items) {
			lastItem := pagedItems[len(pagedItems)-1]
			var lastItemValue string
			switch model.OrderByField(strings.ToUpper(pagination.GetOrderBy())) {
			case model.ORDERBYFIELD_ID, "":
				if lastItem.Id != "" {
					lastItemValue = lastItem.Id
				}
			case model.ORDERBYFIELD_NAME:
				if lastItem.Name != "" {
					lastItemValue = lastItem.Name
				}
			}
			if lastItem.Id != "" {
				newNextPageToken = encodeStringCursor(stringCursor{
					Value: lastItemValue,
					ID:    lastItem.Id,
				})
			}
		}
	} else {
		pagedItems = []model.CatalogSource{}
	}

	return pagedItems, newNextPageToken
}
