package models

type ConfigMapKind struct {
	APIVersion string        `json:"apiVersion"`
	Kind       string        `json:"kind"`
	Metadata   Metadata      `json:"metadata"`
	Data       ConfigMapData `json:"data"`
}

type ConfigMapData struct {
	SamplaCatalogYaml *CatalogContent `json:"sampleCataloYaml"`
	SourcesYaml       *SourcesContent `json:"sourcesYaml"`
}

type CatalogContent struct {
	Source string      `json:"source"`
	Models []BaseModel `json:"models"`
}

type SourcesContent struct {
	Catalogs []CatalogSource `json:"catalogs"`
}

type BaseModel struct {
	Name string `json:"name"`
}
