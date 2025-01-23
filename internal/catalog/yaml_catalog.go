package catalog

import (
	"context"
	"fmt"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
)

type YamlBaseModel struct {
	Catalog    string `yaml:"catalog"`
	Repository string `yaml:"repository"`
	Model      string `yaml:"model"`
}

type YamlModel struct {
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
	BaseModel                []YamlBaseModel `yaml:"baseModel,omitempty"`
}

type YamlCatalog struct {
	Models []YamlModel `yaml:"models"`
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

const yamlCatalogPath = "yamlCatalogPath"

func NewYamlCatalog(source *CatalogSource) (ModelCatalogApi, error) {
	var contents YamlCatalog
	privateProps := source.PrivateProperties
	if len(privateProps) == 0 {
		return nil, fmt.Errorf("missing yaml catalog private properties")
	}
	yamlModelFile, exists := privateProps[yamlCatalogPath].(string)
	if !exists || len(yamlModelFile) == 0 {
		return nil, fmt.Errorf("missing %s string property", yamlCatalogPath)
	}
	bytes, err := os.ReadFile(yamlModelFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s file: %v", yamlCatalogPath, err)
	}
	if err = yaml.UnmarshalStrict(bytes, &contents); err != nil {
		return nil, fmt.Errorf("failed to parse %s file: %v", yamlCatalogPath, err)
	}
	return &yamlCatalogImpl{source: source, contents: &contents}, nil
}

func init() {
	if err := RegisterCatalogType("yaml", NewYamlCatalog); err != nil {
		panic(err)
	}
}
