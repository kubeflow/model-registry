package catalog

import (
	"context"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/util/yaml"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

type YamlArtifacts struct {
	Protocol string `yaml:"protocol"`
	URI      string `yaml:"uri"`
}

type YamlModel struct {
	model.CatalogModel `yaml:",inline"`
	Artifacts          []YamlArtifacts `yaml:"artifacts"`
}

type YamlCatalog struct {
	Source string      `yaml:"source"`
	Models []YamlModel `yaml:"models"`
}

type yamlCatalogImpl struct {
	contents *YamlCatalog
	source   *CatalogSourceConfig
}

func (y yamlCatalogImpl) GetModel(ctx context.Context, name string) (model.CatalogModel, error) {
	//TODO implement me
	panic("implement me")
}

func (y yamlCatalogImpl) ListModels(ctx context.Context, params ListModelsParams) (model.CatalogModelList, error) {
	//TODO implement me
	panic("implement me")
}

func (y yamlCatalogImpl) GetCatalogSource() (model.CatalogSource, error) {
	return y.source.CatalogSource, nil
}

// TODO start background thread to watch file

var _ ModelProvider = &yamlCatalogImpl{}

const yamlCatalogPath = "yamlCatalogPath"

func NewYamlCatalog(source *CatalogSourceConfig) (ModelProvider, error) {
	var contents YamlCatalog
	properties := source.Properties
	if len(properties) == 0 {
		return nil, fmt.Errorf("missing yaml catalog private properties")
	}
	yamlModelFile, exists := properties[yamlCatalogPath].(string)
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

	// override catalog name from Yaml Catalog File if set
	if len(source.Name) > 0 {
		source.Name = contents.Source
	}

	return &yamlCatalogImpl{source: source, contents: &contents}, nil
}

func init() {
	if err := RegisterCatalogType("yaml", NewYamlCatalog); err != nil {
		panic(err)
	}
}
