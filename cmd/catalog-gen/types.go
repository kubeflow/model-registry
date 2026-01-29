package main

// CatalogConfig is the configuration structure for a catalog.
type CatalogConfig struct {
	APIVersion string          `yaml:"apiVersion"`
	Kind       string          `yaml:"kind"`
	Metadata   CatalogMetadata `yaml:"metadata"`
	Spec       CatalogSpec     `yaml:"spec"`
}

// CatalogMetadata contains catalog metadata.
type CatalogMetadata struct {
	Name string `yaml:"name"`
}

// CatalogSpec contains the catalog specification.
type CatalogSpec struct {
	Package   string           `yaml:"package"`
	Entity    EntityConfig     `yaml:"entity"`
	Artifacts []ArtifactConfig `yaml:"artifacts,omitempty"`
	Providers []ProviderConfig `yaml:"providers,omitempty"`
	API       APIConfig        `yaml:"api"`
}

// EntityConfig defines the main entity type.
type EntityConfig struct {
	Name       string           `yaml:"name"`
	Properties []PropertyConfig `yaml:"properties,omitempty"`
}

// ArtifactConfig defines an artifact type linked to the entity.
type ArtifactConfig struct {
	Name       string           `yaml:"name"`
	Properties []PropertyConfig `yaml:"properties,omitempty"`
}

// PropertyConfig defines a property on an entity or artifact.
type PropertyConfig struct {
	Name     string          `yaml:"name"`
	Type     string          `yaml:"type"`
	Required bool            `yaml:"required,omitempty"`
	Items    *PropertyConfig `yaml:"items,omitempty"` // For array types
}

// ProviderConfig defines a data provider.
type ProviderConfig struct {
	Type string `yaml:"type"`
}

// APIConfig defines API settings.
type APIConfig struct {
	BasePath string `yaml:"basePath"`
	Port     int    `yaml:"port"`
}
