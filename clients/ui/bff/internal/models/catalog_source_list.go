package models

type CatalogSourceStatus struct {
	State   string  `json:"state,omitempty"`
	Reason  *string `json:"reason,omitempty"`
	Message *string `json:"message,omitempty"`
}

type CatalogSource struct {
	Id      string               `json:"id"`
	Name    string               `json:"name"`
	Enabled *bool                `json:"enabled,omitempty"`
	Labels  []string             `json:"labels"`
	Status  *CatalogSourceStatus `json:"status,omitempty"`
}

type CatalogSourceList struct {
	NextPageToken string          `json:"nextPageToken"`
	PageSize      int32           `json:"pageSize"`
	Size          int32           `json:"size"`
	Items         []CatalogSource `json:"items,omitempty"`
}
