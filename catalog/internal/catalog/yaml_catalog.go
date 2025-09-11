package catalog

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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
	y.modelsLock.RLock()
	defer y.modelsLock.RUnlock()

	var filteredModels []*model.CatalogModel
	for _, ym := range y.models {
		cm := ym.CatalogModel
		if params.Query != "" {
			query := strings.ToLower(params.Query)
			// Check if query matches name, description, tasks, provider, or libraryName
			if !strings.Contains(strings.ToLower(cm.Name), query) &&
				!strings.Contains(strings.ToLower(cm.GetDescription()), query) &&
				!strings.Contains(strings.ToLower(cm.GetProvider()), query) &&
				!strings.Contains(strings.ToLower(cm.GetLibraryName()), query) {

				// Check tasks
				foundInTasks := false
				for _, task := range cm.GetTasks() { // Use GetTasks() for nil safety
					if strings.Contains(strings.ToLower(task), query) {
						foundInTasks = true
						break
					}
				}
				if !foundInTasks {
					continue // Skip if no match in any searchable field
				}
			}
		}
		filteredModels = append(filteredModels, &cm)
	}

	// Sort the filtered models
	sort.Slice(filteredModels, func(i, j int) bool {
		a := filteredModels[i]
		b := filteredModels[j]

		var less bool
		switch params.OrderBy {
		case model.ORDERBYFIELD_CREATE_TIME:
			// Convert CreateTimeSinceEpoch (string) to int64 for comparison
			// Handle potential nil or conversion errors by treating as 0
			aTime, _ := strconv.ParseInt(a.GetCreateTimeSinceEpoch(), 10, 64)
			bTime, _ := strconv.ParseInt(b.GetCreateTimeSinceEpoch(), 10, 64)
			less = aTime < bTime
		case model.ORDERBYFIELD_LAST_UPDATE_TIME:
			// Convert LastUpdateTimeSinceEpoch (string) to int64 for comparison
			// Handle potential nil or conversion errors by treating as 0
			aTime, _ := strconv.ParseInt(a.GetLastUpdateTimeSinceEpoch(), 10, 64)
			bTime, _ := strconv.ParseInt(b.GetLastUpdateTimeSinceEpoch(), 10, 64)
			less = aTime < bTime
		case model.ORDERBYFIELD_NAME:
			fallthrough
		default:
			// Fallback to name sort if an unknown sort field is provided
			less = strings.Compare(a.Name, b.Name) < 0
		}

		if params.SortOrder == model.SORTORDER_DESC {
			return !less
		}
		return less
	})

	count := len(filteredModels)
	if count > math.MaxInt32 {
		count = math.MaxInt32
	}

	list := model.CatalogModelList{
		Items:    make([]model.CatalogModel, count),
		PageSize: int32(count),
		Size:     int32(count),
	}
	for i := range list.Items {
		list.Items[i] = *filteredModels[i]
	}
	return list, nil // Return the struct value directly
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

func isModelExcluded(modelName string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.HasSuffix(pattern, "*") {
			if strings.HasPrefix(modelName, strings.TrimSuffix(pattern, "*")) {
				return true
			}
		} else if modelName == pattern {
			return true
		}
	}
	return false
}

func (y *yamlCatalogImpl) load(path string, excludedModelsList []string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s file: %v", yamlCatalogPath, err)
	}

	var contents yamlCatalog
	if err = yaml.UnmarshalStrict(bytes, &contents); err != nil {
		return fmt.Errorf("failed to parse %s file: %v", yamlCatalogPath, err)
	}

	models := make(map[string]*yamlModel)
	for i := range contents.Models {
		modelName := contents.Models[i].Name
		if isModelExcluded(modelName, excludedModelsList) {
			continue
		}
		models[modelName] = &contents.Models[i]
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

	// Excluded models is an optional source property.
	var excludedModels []string
	if excludedModelsData, ok := source.Properties["excludedModels"]; ok {
		excludedModelsList, ok := excludedModelsData.([]any)
		if !ok {
			return nil, fmt.Errorf("'excludedModels' property should be a list")
		}
		excludedModels = make([]string, len(excludedModelsList))
		for i, v := range excludedModelsList {
			excludedModels[i], ok = v.(string)
			if !ok {
				return nil, fmt.Errorf("invalid entry in 'excludedModels' list, expected a string")
			}
		}
	}

	p := &yamlCatalogImpl{
		models: make(map[string]*yamlModel),
	}
	err = p.load(yamlModelFile, excludedModels)
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

			err = p.load(yamlModelFile, excludedModels)
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
