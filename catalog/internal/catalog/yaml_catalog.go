package catalog

import (
	"context"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/util/yaml"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

type yamlArtifacts struct {
	Protocol string `yaml:"protocol"`
	URI      string `yaml:"uri"`
}

type yamlModel struct {
	model.CatalogModel `yaml:",inline"`
	Artifacts          []yamlArtifacts `yaml:"artifacts"`
}

type yamlCatalog struct {
	Source string      `yaml:"source"`
	Models []yamlModel `yaml:"models"`
}

type yamlCatalogImpl struct {
	models map[string]*yamlModel
	source *CatalogSourceConfig
}

var _ CatalogSourceProvider = &yamlCatalogImpl{}

func (y *yamlCatalogImpl) GetModel(ctx context.Context, name string) (*model.CatalogModel, error) {
	ym := y.models[name]
	if ym == nil {
		return nil, nil
	}
	cp := ym.CatalogModel
	return &cp, nil
}

func (y *yamlCatalogImpl) ListModels(ctx context.Context, params ListModelsParams) (model.CatalogModelList, error) {
	//TODO implement me
	panic("implement me")
}

// TODO start background thread to watch file

const yamlCatalogPath = "yamlCatalogPath"

func newYamlCatalog(source *CatalogSourceConfig) (CatalogSourceProvider, error) {
	yamlModelFile, exists := source.Properties[yamlCatalogPath].(string)
	if !exists || yamlModelFile == "" {
		return nil, fmt.Errorf("missing %s string property", yamlCatalogPath)
	}
	bytes, err := os.ReadFile(yamlModelFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s file: %v", yamlCatalogPath, err)
	}

	var contents yamlCatalog
	if err = yaml.UnmarshalStrict(bytes, &contents); err != nil {
		return nil, fmt.Errorf("failed to parse %s file: %v", yamlCatalogPath, err)
	}

	// override catalog name from Yaml Catalog File if set
	if source.Name != "" {
		source.Name = contents.Source
	}

	models := make(map[string]*yamlModel, len(contents.Models))
	for i := range contents.Models {
		models[contents.Models[i].Name] = &contents.Models[i]
	}

	return &yamlCatalogImpl{
		models: models,
		source: source,
	}, nil
}

func init() {
	if err := RegisterCatalogType("yaml", newYamlCatalog); err != nil {
		panic(err)
	}
}
