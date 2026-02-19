package basecatalog

import (
	"fmt"
)

// SourceConfig represents the configuration format for model catalogs.
//
// Example:
//
//	model_catalogs:
//	  - name: "Organization AI Models"
//	    id: organization_ai_models
//	    type: yaml
//	    enabled: true
//	    properties:
//	      yamlCatalogPath: dev-organization-models.yaml
//	    labels:
//	      - Organization AI
//
//	# DEPRECATED: Use model_catalogs instead
//	# catalogs: []
type SourceConfig struct {
	// ModelCatalogs contains model catalog source definitions
	ModelCatalogs []Source `yaml:"model_catalogs,omitempty" json:"model_catalogs,omitempty"`

	// Labels contains label definitions for the catalogs
	Labels []map[string]any `yaml:"labels,omitempty" json:"labels,omitempty"`

	// NamedQueries contains predefined query filters
	NamedQueries map[string]map[string]FieldFilter `yaml:"namedQueries,omitempty" json:"namedQueries,omitempty"`

	// DEPRECATED: Use ModelCatalogs instead
	// This field is maintained for backwards compatibility
	Catalogs []Source `yaml:"catalogs,omitempty" json:"catalogs,omitempty"`
}

// GetModelCatalogs returns the merged list of model catalogs, combining the
// new model_catalogs field with the deprecated catalogs field.
// If there are ID conflicts, model_catalogs takes precedence.
func (c *SourceConfig) GetModelCatalogs() []Source {
	if len(c.Catalogs) == 0 {
		return c.ModelCatalogs
	}

	if len(c.ModelCatalogs) == 0 {
		return c.Catalogs
	}

	// Both fields have values - merge them with ModelCatalogs taking precedence
	merged := make(map[string]Source)

	// Add deprecated catalogs first (lower priority)
	for _, source := range c.Catalogs {
		id := source.GetId()
		if id != "" {
			merged[id] = source
		}
	}

	// Add model_catalogs (higher priority, will override)
	for _, source := range c.ModelCatalogs {
		id := source.GetId()
		if id != "" {
			merged[id] = source
		}
	}

	// Convert back to slice
	result := make([]Source, 0, len(merged))
	for _, source := range merged {
		result = append(result, source)
	}

	return result
}

// HasDeprecatedCatalogs returns true if the deprecated "catalogs" field is being used
func (c *SourceConfig) HasDeprecatedCatalogs() bool {
	return len(c.Catalogs) > 0
}

// Validate checks the configuration for common errors
func (c *SourceConfig) Validate() error {
	// Check for duplicate IDs within model catalogs
	seen := make(map[string]bool)

	for _, source := range c.GetModelCatalogs() {
		id := source.GetId()
		if id == "" {
			return fmt.Errorf("model catalog source missing id: %+v", source)
		}
		if seen[id] {
			return fmt.Errorf("duplicate model catalog id: %s", id)
		}
		seen[id] = true
	}

	return nil
}
