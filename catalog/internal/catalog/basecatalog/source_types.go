package basecatalog

import (
	apimodels "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

// Source status constants matching the OpenAPI enum values.
const (
	SourceStatusAvailable          = "available"
	SourceStatusPartiallyAvailable = "partially-available"
	SourceStatusError              = "error"
	SourceStatusDisabled           = "disabled"
)

// CommonSourceFields holds the fields shared between ModelSource and MCPSource
// for field-level merge operations.
type CommonSourceFields struct {
	Name       string
	Enabled    *bool
	Labels     []string
	Type       string
	Properties map[string]any
	Origin     string
}

// MergeCommonSourceFields merges two CommonSourceFields structs with field-level override semantics.
// Fields from 'override' take precedence over 'base' when explicitly set (non-zero/non-nil).
func MergeCommonSourceFields(base, override CommonSourceFields) CommonSourceFields {
	result := base
	if override.Name != "" {
		result.Name = override.Name
	}
	if override.Enabled != nil {
		result.Enabled = override.Enabled
	}
	if override.Labels != nil {
		result.Labels = override.Labels
	}
	if override.Type != "" {
		result.Type = override.Type
	}
	if override.Properties != nil {
		result.Properties = override.Properties
	}
	// Origin follows Properties: use override's origin only when Properties are overridden,
	// so relative paths resolve relative to where they were defined.
	if override.Properties != nil && override.Origin != "" {
		result.Origin = override.Origin
	}
	return result
}

// FieldFilter represents a single field filter within a named query
type FieldFilter struct {
	Operator string `json:"operator" yaml:"operator"`
	Value    any    `json:"value" yaml:"value"`
}

// ModelSource is a single entry from the catalog sources YAML file.
type ModelSource struct {
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
func (s ModelSource) GetId() string {
	return s.Id
}

// MCPSource represents a source of MCP servers
type MCPSource struct {
	Name       string         `json:"name" yaml:"name"`
	ID         string         `json:"id" yaml:"id"`
	Type       string         `json:"type" yaml:"type"`
	Enabled    *bool          `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Properties map[string]any `json:"properties" yaml:"properties"`
	Labels     []string       `json:"labels" yaml:"labels"`

	// Origin is the absolute path of the config file this source was loaded from.
	// This is set automatically during loading and used for resolving relative paths.
	// It is not read from YAML; it's set programmatically.
	Origin string `json:"-" yaml:"-"`
}

// GetId returns the ID of the source
func (s MCPSource) GetId() string {
	return s.ID
}

// IsEnabled returns true if the source is enabled.
// If Enabled is nil, the source is considered enabled by default.
func (s *MCPSource) IsEnabled() bool {
	return s.Enabled == nil || *s.Enabled
}
