package basecatalog

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// ReadSourceConfig reads, parses, and validates a sources configuration file.
func ReadSourceConfig(path string) (*SourceConfig, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &SourceConfig{}
	if err = yaml.UnmarshalStrict(bytes, config); err != nil {
		return nil, err
	}

	if config.HasDeprecatedCatalogs() {
		glog.Warningf("Configuration file %s uses deprecated 'catalogs' field. Please rename to 'model_catalogs'.", path)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration in %s: %w", path, err)
	}

	return config, nil
}

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
	ModelCatalogs []ModelSource `yaml:"model_catalogs,omitempty" json:"model_catalogs,omitempty"`

	// MCPCatalogs contains MCP catalog source definitions
	MCPCatalogs []MCPSource `yaml:"mcp_catalogs,omitempty" json:"mcp_catalogs,omitempty"`

	// Labels contains label definitions for the catalogs
	Labels []map[string]any `yaml:"labels,omitempty" json:"labels,omitempty"`

	// NamedQueries contains predefined query filters
	NamedQueries map[string]map[string]FieldFilter `yaml:"namedQueries,omitempty" json:"namedQueries,omitempty"`

	// DEPRECATED: Use ModelCatalogs instead
	// This field is maintained for backwards compatibility
	Catalogs []ModelSource `yaml:"catalogs,omitempty" json:"catalogs,omitempty"`
}

// GetModelCatalogs returns the merged list of model catalogs, combining the
// new model_catalogs field with the deprecated catalogs field.
// If there are ID conflicts, model_catalogs takes precedence.
func (c *SourceConfig) GetModelCatalogs() []ModelSource {
	if len(c.Catalogs) == 0 {
		return c.ModelCatalogs
	}

	if len(c.ModelCatalogs) == 0 {
		return c.Catalogs
	}

	// Both fields have values. Concatenate the two lists (with
	// ModelCatalogs coming before Catalogs), and remove duplicate entries
	// from Catalogs.

	merged := make([]ModelSource, len(c.ModelCatalogs), len(c.ModelCatalogs)+len(c.Catalogs))
	copy(merged, c.ModelCatalogs)

	mcIDs := make(map[string]struct{}, len(merged))
	for _, catalog := range merged {
		mcIDs[catalog.GetId()] = struct{}{}
	}

	for _, catalog := range c.Catalogs {
		if _, exists := mcIDs[catalog.GetId()]; !exists {
			merged = append(merged, catalog)
		}
	}

	return merged
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

	// Check for duplicate IDs within MCP catalogs, and cross-type collisions with model catalogs
	mcpSeen := make(map[string]bool)
	for _, source := range c.MCPCatalogs {
		id := source.GetId()
		if id == "" {
			return fmt.Errorf("MCP catalog source missing id")
		}
		if mcpSeen[id] {
			return fmt.Errorf("duplicate MCP catalog id: %s", id)
		}
		if seen[id] {
			return fmt.Errorf("id %q used in both model_catalogs and mcp_catalogs", id)
		}
		mcpSeen[id] = true
	}

	return nil
}
