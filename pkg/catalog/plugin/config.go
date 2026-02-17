package plugin

import (
	"fmt"
	"maps"
	"os"

	"k8s.io/apimachinery/pkg/util/yaml"
)

// CatalogSourcesConfig is the root configuration structure for multi-catalog sources.yaml.
type CatalogSourcesConfig struct {
	// APIVersion identifies the config format version (e.g., "catalog/v1alpha1").
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`

	// Kind identifies the config type (e.g., "CatalogSources").
	Kind string `json:"kind" yaml:"kind"`

	// Catalogs maps plugin names to their configurations.
	// The key is the plugin name (e.g., "models", "datasets").
	Catalogs map[string]CatalogSection `json:"catalogs" yaml:"catalogs"`
}

// CatalogSection contains configuration for a single catalog plugin.
type CatalogSection struct {
	// Sources is the list of data sources for this catalog.
	Sources []SourceConfig `json:"sources" yaml:"sources"`

	// Labels defines custom labels available in this catalog.
	Labels []map[string]any `json:"labels,omitempty" yaml:"labels,omitempty"`

	// NamedQueries defines preset filter queries.
	NamedQueries map[string]map[string]FieldFilter `json:"namedQueries,omitempty" yaml:"namedQueries,omitempty"`
}

// SourceConfig represents a single data source configuration.
// This is a unified structure that works across all catalog types.
type SourceConfig struct {
	// ID is the unique identifier for this source.
	ID string `json:"id" yaml:"id"`

	// Name is the human-readable display name.
	Name string `json:"name" yaml:"name"`

	// Type identifies the provider type (e.g., "yaml", "http", "hf").
	Type string `json:"type" yaml:"type"`

	// Enabled indicates whether this source should be loaded.
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`

	// Labels are tags for filtering and categorization.
	Labels []string `json:"labels,omitempty" yaml:"labels,omitempty"`

	// Properties contains provider-specific configuration.
	Properties map[string]any `json:"properties,omitempty" yaml:"properties,omitempty"`

	// IncludedItems are glob patterns for items to include.
	IncludedItems []string `json:"includedItems,omitempty" yaml:"includedItems,omitempty"`

	// ExcludedItems are glob patterns for items to exclude.
	ExcludedItems []string `json:"excludedItems,omitempty" yaml:"excludedItems,omitempty"`

	// Origin is set programmatically to the config file path.
	Origin string `json:"-" yaml:"-"`
}

// FieldFilter represents a filter condition for named queries.
type FieldFilter struct {
	Operator string `json:"operator" yaml:"operator"`
	Value    any    `json:"value" yaml:"value"`
}

// IsEnabled returns true if this source is enabled (defaults to true if nil).
func (s SourceConfig) IsEnabled() bool {
	return s.Enabled == nil || *s.Enabled
}

// LoadConfig loads a CatalogSourcesConfig from a YAML file.
func LoadConfig(path string) (*CatalogSourcesConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	return ParseConfig(data, path)
}

// ParseConfig parses a CatalogSourcesConfig from YAML bytes.
// The origin parameter is used to set the Origin field on sources.
func ParseConfig(data []byte, origin string) (*CatalogSourcesConfig, error) {
	cfg := &CatalogSourcesConfig{}
	if err := yaml.UnmarshalStrict(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set origin on all sources
	for catalogName := range cfg.Catalogs {
		section := cfg.Catalogs[catalogName]
		for i := range section.Sources {
			section.Sources[i].Origin = origin
		}
		cfg.Catalogs[catalogName] = section
	}

	return cfg, nil
}

// LoadConfigs loads and merges multiple config files.
// Later files take precedence over earlier ones for overlapping sources.
func LoadConfigs(paths []string) (*CatalogSourcesConfig, error) {
	if len(paths) == 0 {
		return &CatalogSourcesConfig{
			Catalogs: make(map[string]CatalogSection),
		}, nil
	}

	// Load first config as base
	result, err := LoadConfig(paths[0])
	if err != nil {
		return nil, err
	}

	// Merge subsequent configs
	for _, path := range paths[1:] {
		cfg, err := LoadConfig(path)
		if err != nil {
			return nil, err
		}

		result = MergeConfigs(result, cfg)
	}

	return result, nil
}

// MergeConfigs merges two configs, with override taking precedence.
func MergeConfigs(base, override *CatalogSourcesConfig) *CatalogSourcesConfig {
	result := &CatalogSourcesConfig{
		APIVersion: override.APIVersion,
		Kind:       override.Kind,
		Catalogs:   make(map[string]CatalogSection),
	}

	if result.APIVersion == "" {
		result.APIVersion = base.APIVersion
	}
	if result.Kind == "" {
		result.Kind = base.Kind
	}

	// Copy base catalogs
	maps.Copy(result.Catalogs, base.Catalogs)

	// Merge override catalogs
	for name, overrideSection := range override.Catalogs {
		if baseSection, exists := result.Catalogs[name]; exists {
			result.Catalogs[name] = mergeCatalogSections(baseSection, overrideSection)
		} else {
			result.Catalogs[name] = overrideSection
		}
	}

	return result
}

// mergeCatalogSections merges two CatalogSections.
func mergeCatalogSections(base, override CatalogSection) CatalogSection {
	result := CatalogSection{
		Sources:      make([]SourceConfig, 0),
		Labels:       override.Labels,
		NamedQueries: make(map[string]map[string]FieldFilter),
	}

	if result.Labels == nil {
		result.Labels = base.Labels
	}

	// Build source map from base
	sourceMap := make(map[string]SourceConfig)
	for _, s := range base.Sources {
		sourceMap[s.ID] = s
	}

	// Merge override sources
	for _, s := range override.Sources {
		if existing, ok := sourceMap[s.ID]; ok {
			sourceMap[s.ID] = mergeSourceConfigs(existing, s)
		} else {
			sourceMap[s.ID] = s
		}
	}

	// Convert back to slice
	for _, s := range sourceMap {
		result.Sources = append(result.Sources, s)
	}

	// Merge named queries
	maps.Copy(result.NamedQueries, base.NamedQueries)
	for name, filters := range override.NamedQueries {
		if existing, ok := result.NamedQueries[name]; ok {
			// Merge field filters
			maps.Copy(existing, filters)
			result.NamedQueries[name] = existing
		} else {
			result.NamedQueries[name] = filters
		}
	}

	return result
}

// mergeSourceConfigs merges two SourceConfigs with field-level merging.
func mergeSourceConfigs(base, override SourceConfig) SourceConfig {
	result := base

	result.ID = override.ID

	if override.Name != "" {
		result.Name = override.Name
	}

	if override.Type != "" {
		result.Type = override.Type
	}

	if override.Enabled != nil {
		result.Enabled = override.Enabled
	}

	if override.Labels != nil {
		result.Labels = override.Labels
	}

	if override.Properties != nil {
		result.Properties = override.Properties
		result.Origin = override.Origin
	}

	if override.IncludedItems != nil {
		result.IncludedItems = override.IncludedItems
	}

	if override.ExcludedItems != nil {
		result.ExcludedItems = override.ExcludedItems
	}

	return result
}
