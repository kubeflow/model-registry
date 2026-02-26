package models

// CatalogRestEntityType represents catalog-specific REST API entity types
type CatalogRestEntityType string

const (
	RestEntityCatalogModel    CatalogRestEntityType = "CatalogModel"
	RestEntityCatalogArtifact CatalogRestEntityType = "CatalogArtifact"
	RestEntityMCPServer       CatalogRestEntityType = "MCPServer"
	RestEntityMCPServerTool   CatalogRestEntityType = "MCPServerTool"
)
