package models

import "github.com/kubeflow/model-registry/internal/platform/db/entity"

const (
	SortOrderAsc  = entity.SortOrderAsc
	SortOrderDesc = entity.SortOrderDesc
	OrderByID     = entity.OrderByID
)

var (
	DefaultSortOrder = entity.DefaultSortOrder
	DefaultPageSize  = entity.DefaultPageSize
	DefaultOrderBy   = entity.DefaultOrderBy
)

type Pagination = entity.Pagination
