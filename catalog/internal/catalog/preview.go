package catalog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// PreviewConfig represents the parsed preview request configuration.
type PreviewConfig struct {
	Type           string         `json:"type" yaml:"type"`
	IncludedModels []string       `json:"includedModels,omitempty" yaml:"includedModels,omitempty"`
	ExcludedModels []string       `json:"excludedModels,omitempty" yaml:"excludedModels,omitempty"`
	Properties     map[string]any `json:"properties,omitempty" yaml:"properties,omitempty"`
}

// ParsePreviewConfig parses the uploaded config bytes into a PreviewConfig.
// Extra fields (like name, id, enabled) are ignored so users can paste
// a full source config entry directly for preview.
func ParsePreviewConfig(configBytes []byte) (*PreviewConfig, error) {
	var config PreviewConfig
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if config.Type == "" {
		return nil, fmt.Errorf("missing required field: type")
	}

	// Validate filter patterns early
	if err := ValidateSourceFilters(config.IncludedModels, config.ExcludedModels); err != nil {
		return nil, err
	}

	return &config, nil
}

// PreviewSourceModels loads models from the source configuration and returns
// preview results showing which models would be included or excluded.
// If catalogDataBytes is provided, it will be used directly instead of reading from yamlCatalogPath.
func PreviewSourceModels(ctx context.Context, config *PreviewConfig, catalogDataBytes []byte) ([]model.ModelPreviewResult, error) {
	// Load all model names from the source (without filtering)
	modelNames, err := loadModelNamesFromSource(ctx, config, catalogDataBytes)
	if err != nil {
		return nil, err
	}

	// Create a ModelFilter from the config
	filter, err := NewModelFilter(config.IncludedModels, config.ExcludedModels)
	if err != nil {
		return nil, fmt.Errorf("invalid filter configuration: %w", err)
	}

	// Create preview results for each model
	results := make([]model.ModelPreviewResult, 0, len(modelNames))
	for _, name := range modelNames {
		included := filter == nil || filter.Allows(name)
		results = append(results, model.ModelPreviewResult{
			Name:     name,
			Included: included,
		})
	}

	return results, nil
}

// loadModelNamesFromSource loads model names from the specified source type.
// If catalogDataBytes is provided, it takes precedence over reading from file paths.
func loadModelNamesFromSource(ctx context.Context, config *PreviewConfig, catalogDataBytes []byte) ([]string, error) {
	switch config.Type {
	case "yaml":
		return loadYamlModelNames(ctx, config, catalogDataBytes)
	case "hf", "huggingface":
		return loadHFModelNames(ctx, config)
	default:
		return nil, fmt.Errorf("unsupported source type: %s", config.Type)
	}
}

// loadHFModelNames fetches model names from the HuggingFace API for preview.
// For HF sources, includedModels specifies which models to fetch from HuggingFace.
// This function calls the HF API to validate models exist and get their actual names.
func loadHFModelNames(ctx context.Context, config *PreviewConfig) ([]string, error) {
	// SECURITY: Override the URL property to prevent SSRF attacks.
	// An attacker could otherwise set a custom URL to leak the HF API key
	// to an attacker-controlled domain.
	// We ensure requests only go to the official HuggingFace API.
	if config.Properties == nil {
		config.Properties = make(map[string]any)
	}

	if customURL, exists := config.Properties["url"]; exists {
		glog.Warningf("HuggingFace preview: custom URL %q was ignored for security reasons (SSRF prevention)", customURL)
		delete(config.Properties, "url")
	}

	// Create HF preview provider (reuses hfModelProvider from hf_catalog.go)
	provider, err := NewHFPreviewProvider(config)
	if err != nil {
		return nil, err
	}

	// Fetch model names from HuggingFace API
	return provider.FetchModelNamesForPreview(ctx, config.IncludedModels)
}

// loadYamlModelNames loads model names from a YAML catalog.
// If catalogDataBytes is provided (stateless mode), it is used directly.
// Otherwise, the catalog is read from yamlCatalogPath (path mode).
func loadYamlModelNames(ctx context.Context, config *PreviewConfig, catalogDataBytes []byte) ([]string, error) {
	var catalogBytes []byte

	if len(catalogDataBytes) > 0 {
		// Stateless mode: use provided catalog data directly
		catalogBytes = catalogDataBytes
	} else {
		// Path mode: read from yamlCatalogPath
		path, ok := config.Properties[yamlCatalogPathKey].(string)
		if !ok || path == "" {
			return nil, fmt.Errorf("missing required property: %s (provide catalogData file or set yamlCatalogPath in config)", yamlCatalogPathKey)
		}

		// Resolve relative paths - for preview, we use the current working directory
		if !filepath.IsAbs(path) {
			cwd, err := os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("failed to get working directory: %w", err)
			}
			path = filepath.Join(cwd, path)
		}

		// Read the catalog file
		var err error
		catalogBytes, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read catalog file %s: %w", path, err)
		}
	}

	var catalog yamlCatalog
	if err := yaml.UnmarshalStrict(catalogBytes, &catalog); err != nil {
		return nil, fmt.Errorf("failed to parse catalog file: %w", err)
	}

	// Extract model names
	names := make([]string, 0, len(catalog.Models))
	for _, m := range catalog.Models {
		names = append(names, m.Name)
	}

	return names, nil
}
