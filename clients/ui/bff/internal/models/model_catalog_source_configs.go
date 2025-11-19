package models

type CatalogSourceConfig struct {
	Id                  string   `json:"id"`
	Name                string   `json:"name"`
	Type                string   `json:"type"`
	Enabled             *bool    `json:"enabled,omitempty"`
	Labels              []string `json:"labels"`
	ApiKey              *string  `json:"apiKey,omitempty"`
	AllowedOrganization *string  `json:"allowedOrganization,omitempty"`
	IncludedModels      []string `json:"includedModels,omitempty"`
	ExcludedModels      []string `json:"excludedModels,omitempty"`
	IsDefault           *bool    `json:"isDefault,omitempty"`
	Yaml                *string  `json:"yaml,omitempty"`
}

type CatalogSourceConfigPayload = CatalogSourceConfig

type CatalogSourceConfigList struct {
	Catalogs []CatalogSourceConfig `json:"catalogs,omitempty"`
}
