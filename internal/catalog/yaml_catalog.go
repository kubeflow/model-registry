package catalog

import (
	"context"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

type YamlCatalog struct {
	Models []struct {
		Name                     string   `yaml:"name"`
		Provider                 string   `yaml:"provider"`
		Description              string   `yaml:"description"`
		ReadmeLink               string   `yaml:"readmeLink"`
		Language                 []string `yaml:"language"`
		License                  string   `yaml:"license"`
		LicenseLink              string   `yaml:"licenseLink"`
		LibraryName              string   `yaml:"libraryName"`
		Tags                     []string `yaml:"tags"`
		Tasks                    []string `yaml:"tasks"`
		CreateTimeSinceEpoch     int64    `yaml:"createTimeSinceEpoch"`
		LastUpdateTimeSinceEpoch int64    `yaml:"lastUpdateTimeSinceEpoch"`
		BaseModel                []struct {
			Catalog    string `yaml:"catalog"`
			Repository string `yaml:"repository"`
			Model      string `yaml:"model"`
		} `yaml:"baseModel,omitempty"`
	} `yaml:"models"`
}

type yamlCatalogImpl struct {
	contents *YamlCatalog
	source   *CatalogSource
}

func (y yamlCatalogImpl) GetCatalogModel(ctx context.Context, modelId string) (openapi.CatalogModel, error) {
	//TODO implement me
	panic("implement me")
}

func (y yamlCatalogImpl) GetCatalogModelVersion(ctx context.Context, modelId string, versionId string) (openapi.CatalogModelVersion, error) {
	//TODO implement me
	panic("implement me")
}

func (y yamlCatalogImpl) GetCatalogModelVersions(ctx context.Context, modelId string, nameParam string, externalIdParam string, pageSizeParam string, orderByParam openapi.OrderByField, sortOrderParam openapi.SortOrder, offsetParam string) (openapi.CatalogModelVersionList, error) {
	//TODO implement me
	panic("implement me")
}

func (y yamlCatalogImpl) GetCatalogModels(ctx context.Context, nameParam string, externalIdParam string, pageSizeParam string, orderByParam openapi.OrderByField, sortOrderParam openapi.SortOrder, offsetParam string) (openapi.CatalogModelList, error) {
	//TODO implement me
	panic("implement me")
}

func (y yamlCatalogImpl) GetCatalogSource() (openapi.CatalogSource, error) {
	return y.source.CatalogSource, nil
}

// TODO start background thread to watch file

var _ ModelCatalogApi = &yamlCatalogImpl{}

func NewYamlCatalog(source *CatalogSource) ModelCatalogApi {
	// TODO read file contents from config
	var contents YamlCatalog
	return &yamlCatalogImpl{source: source, contents: &contents}
}

func init() {
	RegisterCatalogType("yaml", NewYamlCatalog)
}
