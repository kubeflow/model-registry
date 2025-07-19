package models

const (
	SortOrderAsc  = "ASC"
	SortOrderDesc = "DESC"
	OrderByID     = "id"
)

var (
	DefaultSortOrder = SortOrderAsc
	DefaultPageSize  = int32(0)
	DefaultOrderBy   = OrderByID
)

type Pagination struct {
	PageSize      *int32  `json:"pageSize,omitempty"`
	OrderBy       *string `json:"orderBy,omitempty"`
	SortOrder     *string `json:"sortOrder,omitempty"`
	NextPageToken *string `json:"nextPageToken,omitempty"`
	FilterQuery   *string `json:"filterQuery,omitempty"`
}

func (p *Pagination) GetNextPageToken() string {
	if p.NextPageToken == nil {
		return ""
	}

	return *p.NextPageToken
}

func (p *Pagination) GetOrderBy() string {
	if p.OrderBy == nil {
		return DefaultOrderBy
	}

	return *p.OrderBy
}

func (p *Pagination) GetSortOrder() string {
	if p.SortOrder == nil {
		return DefaultSortOrder
	}

	return *p.SortOrder
}

func (p *Pagination) GetPageSize() int32 {
	if p.PageSize == nil {
		return DefaultPageSize
	}

	return *p.PageSize
}

func (p *Pagination) GetFilterQuery() string {
	if p.FilterQuery == nil {
		return ""
	}

	return *p.FilterQuery
}

func (p *Pagination) SetNextPageToken(token *string) {
	p.NextPageToken = token
}
