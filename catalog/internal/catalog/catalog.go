package catalog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/util/yaml"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

type SortDirection int

const (
	SortDirectionAscending SortDirection = iota
	SortDirectionDescending
)

type SortField int

const (
	SortByUnspecified SortField = iota
	SortByName
	SortByPublished
)

type ListModelsParams struct {
	Query         string
	SortBy        SortField
	SortDirection SortDirection
}

// CatalogSourceProvider is implemented by catalog source types, e.g. YamlCatalog
type CatalogSourceProvider interface {
	// GetModel returns model metadata for a single model by its name. If
	// nothing is found with the name provided it returns nil, without an
	// error.
	GetModel(ctx context.Context, name string) (*model.CatalogModel, error)
	ListModels(ctx context.Context, params ListModelsParams) (model.CatalogModelList, error)
}

// CatalogSourceConfig is a single entry from the catalog sources YAML file.
type CatalogSourceConfig struct {
	model.CatalogSource `json:",inline"`

	// Catalog type to use, must match one of the registered types
	Type string `json:"type"`

	// Properties used for configuring the catalog connection based on catalog implementation
	Properties map[string]any `json:"properties,omitempty"`
}

// sourceConfig is the structure for the catalog sources YAML file.
type sourceConfig struct {
	Catalogs []CatalogSourceConfig `json:"catalogs"`
}

type CatalogTypeRegisterFunc func(source *CatalogSourceConfig) (CatalogSourceProvider, error)

var registeredCatalogTypes = make(map[string]CatalogTypeRegisterFunc, 0)

func RegisterCatalogType(catalogType string, callback CatalogTypeRegisterFunc) error {
	if _, exists := registeredCatalogTypes[catalogType]; exists {
		return fmt.Errorf("catalog type %s already exists", catalogType)
	}
	registeredCatalogTypes[catalogType] = callback
	return nil
}

type CatalogSource struct {
	Provider CatalogSourceProvider
	Metadata model.CatalogSource
}

func LoadCatalogSources(catalogsPath string) (map[string]CatalogSource, error) {
	// Get absolute path of the catalog config file
	absConfigPath, err := filepath.Abs(catalogsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for %s: %v", catalogsPath, err)
	}

	// Get the directory of the config file to resolve relative paths
	configDir := filepath.Dir(absConfigPath)

	// Save current working directory
	originalWd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %v", err)
	}

	// Change to the config directory to make relative paths work
	if err := os.Chdir(configDir); err != nil {
		return nil, fmt.Errorf("failed to change to config directory %s: %v", configDir, err)
	}

	// Ensure we restore the original working directory when we're done
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			glog.Errorf("failed to restore original working directory %s: %v", originalWd, err)
		}
	}()

	config := sourceConfig{}
	bytes, err := os.ReadFile(absConfigPath)
	if err != nil {
		return nil, err
	}

	if err = yaml.UnmarshalStrict(bytes, &config); err != nil {
		return nil, err
	}

	catalogs := make(map[string]CatalogSource, len(config.Catalogs))
	for _, catalogConfig := range config.Catalogs {
		catalogType := catalogConfig.Type
		glog.Infof("reading config type %s...", catalogType)
		registerFunc, ok := registeredCatalogTypes[catalogType]
		if !ok {
			return nil, fmt.Errorf("catalog type %s not registered", catalogType)
		}
		id := catalogConfig.GetId()
		if len(id) == 0 {
			return nil, fmt.Errorf("invalid catalog id %s", id)
		}
		if _, exists := catalogs[id]; exists {
			return nil, fmt.Errorf("duplicate catalog id %s", id)
		}
		provider, err := registerFunc(&catalogConfig)
		if err != nil {
			return nil, fmt.Errorf("error reading catalog type %s with id %s: %v", catalogType, id, err)
		}

		catalogs[id] = CatalogSource{
			Provider: provider,
			Metadata: catalogConfig.CatalogSource,
		}

		glog.Infof("loaded config %s of type %s", id, catalogType)
	}

	return catalogs, nil
}
