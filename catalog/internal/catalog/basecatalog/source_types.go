package basecatalog

import (
	apimodels "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

// FieldFilter represents a single field filter within a named query
type FieldFilter struct {
	Operator string `json:"operator" yaml:"operator"`
	Value    any    `json:"value" yaml:"value"`
}

// Source is a single entry from the catalog sources YAML file.
type Source struct {
	apimodels.CatalogSource `json:",inline"`

	// Catalog type to use, must match one of the registered types
	Type string `json:"type"`

	// Properties used for configuring the catalog connection based on catalog implementation
	Properties map[string]any `json:"properties,omitempty"`

	// Origin is the absolute path of the config file this source was loaded from.
	// This is set automatically during loading and used for resolving relative paths.
	// It is not read from YAML; it's set programmatically.
	Origin string `json:"-" yaml:"-"`
}

// GetId returns the ID of the source
func (s Source) GetId() string {
	return s.Id
}
