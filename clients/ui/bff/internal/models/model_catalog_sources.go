package models

type CatalogSourceProperties struct {
	AccessToken         *string `json:"accessToken,omitempty"`
	AllowedOrganization *string `json:"allowedOrganization,omitempty"`
	YamlCatalogPath     *string `json:"yamlCatalogPath,omitempty"`
}

type CatalogSourceConfig struct {
	Id             string                   `json:"id"`
	Name           string                   `json:"name"`
	Type           string                   `json:"type"`
	Enabled        *bool                    `json:"enabled,omitempty"`
	IncludedModels []string                 `json:"includedModels,omitempty"`
	Labels         []string                 `json:"labels"`
	ExcludedModels []string                 `json:"excludedModels,omitempty"`
	Properties     *CatalogSourceProperties `json:"properties,omitempty"`
}

type CatalogSourceConfigPayload = CatalogSourceConfig

type CatalogSourceConfigList struct {
	Catalogs []CatalogSourceConfig `json:"catalogs,omitempty"`
}
