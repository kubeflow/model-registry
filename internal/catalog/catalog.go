package catalog

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"k8s.io/apimachinery/pkg/util/yaml"
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
	openapi.CatalogSource `json:",inline"`

	// Catalog type to use, must match one of the registered types
	Type string `json:"type"`

	// private properties used for configuring the catalog connection based on catalog implementation
	PrivateProperties map[string]interface{} `json:"privateProperties,omitempty"`
}

type CatalogsConfig struct {
	Catalogs []CatalogSource `json:"catalogs"`
}

type CatalogTypeRegisterFunc func (source *CatalogSource) (ModelCatalogApi, error)

var registeredCatalogTypes = make(map[string]CatalogTypeRegisterFunc, 0)

func RegisterCatalogType(catalogType string, callback CatalogTypeRegisterFunc) error {
	if _, exists := registeredCatalogTypes[catalogType]; exists {
		return fmt.Errorf("catalog type %s already exists", catalogType)
	}
	registeredCatalogTypes[catalogType] = callback
	return nil
}

func LoadCatalogSources(catalogsPath string) (map[string]ModelCatalogApi, error) {
	config := CatalogsConfig{}
	bytes, err := os.ReadFile(catalogsPath)
	if err != nil {
		return nil, err
	}

	if err = yaml.UnmarshalStrict(bytes, &config); err != nil {
		return nil, err
	}

	catalogs := make(map[string]ModelCatalogApi)
	for _, catalog := range config.Catalogs {
		catalogType := catalog.Type
		glog.Infof("reading config type %s...", catalogType)
		registerFunc, ok := registeredCatalogTypes[catalogType]
		if !ok {
			return nil, fmt.Errorf("catalog type %s not registered", catalogType)
		}
		id := catalog.GetId()
		if len(id) == 0 {
			return nil, fmt.Errorf("invalid catalog id %s", id)
		}
		if _, exists := catalogs[id]; exists {
			return nil, fmt.Errorf("duplicate catalog id %s", id)
		}
		api, err := registerFunc(&catalog)
		if err != nil {
			return nil, fmt.Errorf("error reading catalog type %s with id %s: %v", catalogType, id, err)
		}
		catalogs[id] = api
		glog.Infof("loaded config %s of type %s", id, catalogType)
	}
	return catalogs, nil
}
