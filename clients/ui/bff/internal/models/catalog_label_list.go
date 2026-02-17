package models

// CatalogLabel represents a label used to categorize catalog sources.
// Labels can have optional display names and descriptions for UI presentation.
type CatalogLabel struct {
	Name        *string `json:"name"`
	DisplayName *string `json:"displayName,omitempty"`
	Description *string `json:"description,omitempty"`
}

// CatalogLabelList represents a paginated list of catalog labels.
type CatalogLabelList struct {
	NextPageToken string         `json:"nextPageToken"`
	PageSize      int32          `json:"pageSize"`
	Size          int32          `json:"size"`
	Items         []CatalogLabel `json:"items"`
}
