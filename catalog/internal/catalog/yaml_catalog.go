package catalog

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"

	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/golang/glog"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

type yamlModel struct {
	model.CatalogModel `yaml:",inline"`
	Artifacts          []*model.CatalogModelArtifact `yaml:"artifacts"`
}

type yamlCatalog struct {
	Source string      `yaml:"source"`
	Models []yamlModel `yaml:"models"`
}

type yamlCatalogImpl struct {
	modelsLock sync.RWMutex
	models     map[string]*yamlModel
}

var _ CatalogSourceProvider = &yamlCatalogImpl{}

func (y *yamlCatalogImpl) GetModel(ctx context.Context, name string) (*model.CatalogModel, error) {
	y.modelsLock.RLock()
	defer y.modelsLock.RUnlock()

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

func (y *yamlCatalogImpl) GetArtifacts(ctx context.Context, name string) (*model.CatalogModelArtifactList, error) {
	y.modelsLock.RLock()
	defer y.modelsLock.RUnlock()

	ym := y.models[name]
	if ym == nil {
		return nil, nil
	}

	count := len(ym.Artifacts)
	if count > math.MaxInt32 {
		count = math.MaxInt32
	}

	list := model.CatalogModelArtifactList{
		Items:    make([]model.CatalogModelArtifact, count),
		PageSize: int32(count),
		Size:     int32(count),
	}
	for i := range list.Items {
		list.Items[i] = *ym.Artifacts[i]
	}
	return &list, nil
}

func (y *yamlCatalogImpl) load(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s file: %v", yamlCatalogPath, err)
	}

	var contents yamlCatalog
	if err = yaml.UnmarshalStrict(bytes, &contents); err != nil {
		return fmt.Errorf("failed to parse %s file: %v", yamlCatalogPath, err)
	}

	models := make(map[string]*yamlModel, len(contents.Models))
	for i := range contents.Models {
		models[contents.Models[i].Name] = &contents.Models[i]
	}

	y.modelsLock.Lock()
	defer y.modelsLock.Unlock()
	y.models = models

	return nil
}

const yamlCatalogPath = "yamlCatalogPath"

func newYamlCatalog(source *CatalogSourceConfig) (CatalogSourceProvider, error) {
	yamlModelFile, exists := source.Properties[yamlCatalogPath].(string)
	if !exists || yamlModelFile == "" {
		return nil, fmt.Errorf("missing %s string property", yamlCatalogPath)
	}

	yamlModelFile, err := filepath.Abs(yamlModelFile)
	if err != nil {
		return nil, fmt.Errorf("abs: %w", err)
	}

	p := &yamlCatalogImpl{}
	err = p.load(yamlModelFile)
	if err != nil {
		return nil, err
	}

	go func() {
		changes, err := getMonitor().Path(yamlModelFile)
		if err != nil {
			glog.Errorf("unable to watch YAML catalog file: %v", err)
			// Not fatal, we just won't get automatic updates.
		}

		for range changes {
			glog.Infof("Reloading YAML catalog %s", yamlModelFile)

			err = p.load(yamlModelFile)
			if err != nil {
				glog.Errorf("unable to load YAML catalog: %v", err)
			}
		}
	}()

	return p, nil
}

func init() {
	if err := RegisterCatalogType("yaml", newYamlCatalog); err != nil {
		panic(err)
	}
}
