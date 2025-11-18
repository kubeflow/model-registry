package models

type CatalogSourceProperties struct {
	ApiKey              *string  `json:"apiKey,omitempty"`
	AllowedOrganization *string  `json:"allowedOrganization,omitempty"`
	YamlCatalogPath     *string  `json:"yamlCatalogPath,omitempty"`
	IncludedModels      []string `json:"includedModels,omitempty"`
	ExcludedModels      []string `json:"excludedModels,omitempty"`
}

type CatalogSourceConfig struct {
	Id         string                   `json:"id"`
	Name       string                   `json:"name"`
	Type       string                   `json:"type"`
	Enabled    *bool                    `json:"enabled,omitempty"`
	Labels     []string                 `json:"labels"`
	Properties *CatalogSourceProperties `json:"properties,omitempty"`
	IsDefault  *bool                    `json:"isDefault,omitempty"`
}

type CatalogSourceConfigPayload = CatalogSourceConfig

type CatalogSourceConfigList struct {
	Catalogs []CatalogSourceConfig `json:"catalogs,omitempty"`
}
