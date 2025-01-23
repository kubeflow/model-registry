package catalog

import (
	"context"
	"fmt"
	"github.com/kubeflow/model-registry/pkg/openapi"
	yaml3 "gopkg.in/yaml.v3"
	"log"
	"os"
)

// ModelCatalogApi is implemented by catalog types, e.g. YamlCatalog
type ModelCatalogApi interface {
	GetCatalogModel(ctx context.Context, modelId string) (openapi.CatalogModel, error)
	GetCatalogModelVersion(ctx context.Context, modelId string, versionId string) (openapi.CatalogModelVersion, error)
	GetCatalogModelVersions(ctx context.Context, modelId string, nameParam string, externalIdParam string, pageSizeParam string, orderByParam openapi.OrderByField, sortOrderParam openapi.SortOrder, offsetParam string) (openapi.CatalogModelVersionList, error)
	GetCatalogModels(ctx context.Context, nameParam string, externalIdParam string, pageSizeParam string, orderByParam openapi.OrderByField, sortOrderParam openapi.SortOrder, offsetParam string) (openapi.CatalogModelList, error)
	GetCatalogSource() (openapi.CatalogSource, error)
}

type CatalogSource struct {
	openapi.CatalogSource

	// Catalog type to use, must match one of the registered types
	Type string `json:"type"`

	// private properties used for configuring the catalog connection based on catalog implementation
	PrivateProperties *map[string]openapi.MetadataValue `json:"privateProperties,omitempty"`
}

type CatalogsConfig struct {
	Catalogs []CatalogSource `json:"catalogs"`
}

type CatalogTypeRegisterFunc func (source *CatalogSource) ModelCatalogApi

var catalogTypes = make(map[string]CatalogTypeRegisterFunc, 0)

func RegisterCatalogType(catalogType string, callback CatalogTypeRegisterFunc) {
	catalogTypes[catalogType] = callback
}

func LoadCatalogSources(catalogsPath string) (map[string]ModelCatalogApi, error) {
	config := CatalogsConfig{}
	f, err := os.Open(catalogsPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml3.NewDecoder(f)
	if err := decoder.Decode(&config); err != nil {
		log.Fatal(err)
	}

	catalogs := make(map[string]ModelCatalogApi)
	for _, catalog := range config.Catalogs {
		catalogType := catalog.Type
		registerFunc, ok := catalogTypes[catalogType]
		if !ok {
			return nil, fmt.Errorf("catalog type %s not registered", catalogType)
		}
		id := catalog.GetId()
		catalogs[id] = registerFunc(&catalog)
	}
	return catalogs, nil
}
