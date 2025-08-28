package models

type CatalogSourceProperties struct {
	YamlCatalogPath *string  `json:"yamlCatalogPath"`
	Models          []string `json:"models"`
	ExcludedModels  []string `json:"exclidedModels"`
	APIKey          *string  `json:"apiKey"`
	URL             *string  `json:"url"`
	ModelLimit      *int     `json:"modelLimit"`
}

type CatalogSource struct {
	Name       string                   `json:"name"`
	Id         string                   `json:"id"`
	Type       string                   `json:"type"`
	Enabled    *bool                    `json:"enabled"`
	Properties *CatalogSourceProperties `json:"properties"`
}
