package models

import "github.com/kubeflow/hub/internal/platform/db/entity"

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
