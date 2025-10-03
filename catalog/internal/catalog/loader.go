package catalog

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/golang/glog"
	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	apimodels "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/apiutils"
	mrmodels "github.com/kubeflow/model-registry/internal/db/models"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// ModelProviderRecord contains one model and its associated artifacts.
type ModelProviderRecord struct {
	Model     dbmodels.CatalogModel
	Artifacts []dbmodels.CatalogArtifact
}

// ModelProviderFunc emits models and related data in the channel it returns. It is
// expected to spawn a goroutine and return immediately. The returned channel must
// close when the goroutine ends. The goroutine should end when the context is
// canceled, but may end sooner.
type ModelProviderFunc func(ctx context.Context, source *Source, reldir string) (<-chan ModelProviderRecord, error)

var registeredModelProviders = map[string]ModelProviderFunc{}

func RegisterModelProvider(name string, callback ModelProviderFunc) error {
	if _, exists := registeredModelProviders[name]; exists {
		return fmt.Errorf("provider type %s already exists", name)
	}
	registeredModelProviders[name] = callback
	return nil
}

// LoadCatalogSources processes sources YAML files, returning the source data,
// and also loads models from those sources into the database.
func LoadCatalogSources(ctx context.Context, services service.Services, paths []string) (*SourceCollection, error) {
	l := newLoader(services)

	for _, path := range paths {
		err := l.Load(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}

		go func(path string) {
			changes, err := getMonitor().Path(ctx, path)
			if err != nil {
				glog.Errorf("unable to watch sources file (%s): %v", path, err)
				// Not fatal, we just won't get automatic updates.
			}

			for range changes {
				glog.Infof("Reloading sources %s", path)

				err = l.Load(ctx, path)
				if err != nil {
					glog.Errorf("unable to load sources: %v", err)
				}
			}
		}(path)
	}

	return l.Sources, nil
}

// sourceConfig is the structure for the catalog sources YAML file.
type sourceConfig struct {
	Catalogs []Source `json:"catalogs"`
}

// Source is a single entry from the catalog sources YAML file.
type Source struct {
	apimodels.CatalogSource `json:",inline"`

	// Catalog type to use, must match one of the registered types
	Type string `json:"type"`

	// Properties used for configuring the catalog connection based on catalog implementation
	Properties map[string]any `json:"properties,omitempty"`
}

type loader struct {
	Sources   *SourceCollection
	services  service.Services
	closersMu sync.Mutex
	closers   map[string]func()
}

func newLoader(services service.Services) *loader {
	return &loader{
		Sources:  NewSourceCollection(),
		services: services,
		closers:  map[string]func(){},
	}
}

// Load processes (or re-processes) a sources config file.
func (l *loader) Load(ctx context.Context, path string) error {
	// Get absolute path of the catalog config file
	path, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %v", path, err)
	}

	config, err := l.read(path)
	if err != nil {
		return err
	}

	err = l.updateSources(path, config)
	if err != nil {
		return err
	}

	return l.updateDatabase(ctx, path, config)
}

func (l *loader) read(path string) (*sourceConfig, error) {
	config := &sourceConfig{}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = yaml.UnmarshalStrict(bytes, &config); err != nil {
		return nil, err
	}

	enabledSources := make([]Source, 0, len(config.Catalogs))

	// Remove disabled sources and explicitly set enabled on the others.
	for _, source := range config.Catalogs {
		// If enabled is explicitly set to false, skip
		if source.HasEnabled() && *source.Enabled == false {
			continue
		}
		// If not explicitly set, default to enabled
		source.CatalogSource.Enabled = apiutils.Of(true)
		enabledSources = append(enabledSources, source)
	}
	config.Catalogs = enabledSources

	return config, nil
}

func (l *loader) updateSources(path string, config *sourceConfig) error {
	sources := make(map[string]apimodels.CatalogSource, len(config.Catalogs))

	for _, source := range config.Catalogs {
		glog.Infof("reading config type %s...", source.Type)
		id := source.GetId()
		if len(id) == 0 {
			return fmt.Errorf("invalid catalog id %s", id)
		}
		if _, exists := sources[id]; exists {
			return fmt.Errorf("duplicate catalog id %s", id)
		}

		sources[id] = source.CatalogSource

		glog.Infof("loaded source %s of type %s", id, source.Type)
	}

	return l.Sources.Merge(path, sources)
}

func (l *loader) updateDatabase(ctx context.Context, path string, config *sourceConfig) error {
	ctx, cancel := context.WithCancel(ctx)

	l.closersMu.Lock()
	if l.closers[path] != nil {
		l.closers[path]()
	}
	l.closers[path] = cancel
	l.closersMu.Unlock()

	records := l.readProviderRecords(ctx, path, config)

	go func() {
		for record := range records {
			attr := record.Model.GetAttributes()
			if attr == nil || attr.Name == nil {
				continue
			}

			glog.Infof("Loading model %s with %d artifact(s)", *attr.Name, len(record.Artifacts))

			model, err := l.services.CatalogModelRepository.Save(record.Model)
			if err != nil {
				glog.Errorf("%s: unable to save: %v", *attr.Name, err)
				continue
			}

			modelID := model.GetID()
			if modelID == nil {
				glog.Errorf("%s: model has no ID after save")
				continue
			}

			// Remove any catalog model artifacts that existed
			// before. Any other artifact types will be added to
			// what's there.
			err = l.services.CatalogArtifactRepository.DeleteByParentID(service.CatalogModelArtifactTypeName, *modelID)
			if err != nil {
				glog.Errorf("%s: unable to remove old catalog model artifacts: %v", err)
			}

			for i, artifact := range record.Artifacts {
				switch {
				case artifact.CatalogModelArtifact != nil:
					_, err = l.services.CatalogModelArtifactRepository.Save(artifact.CatalogModelArtifact, modelID)
				case artifact.CatalogMetricsArtifact != nil:
					_, err = l.services.CatalogMetricsArtifactRepository.Save(artifact.CatalogMetricsArtifact, modelID)
				default:
					err = errors.New("unknown artifact type")
				}

				if err != nil {
					glog.Errorf("%s, artifact %d: %v", *attr.Name, i, err)
				}
			}
		}
	}()

	return nil
}

// readProviderRecords calls the provider for every configured source and
// merges the returned channels together. The returned channel is closed when
// the last provider channel is closed.
func (l *loader) readProviderRecords(ctx context.Context, path string, config *sourceConfig) <-chan ModelProviderRecord {
	configDir := filepath.Dir(path)

	ch := make(chan ModelProviderRecord)
	var wg sync.WaitGroup

	for _, source := range config.Catalogs {
		glog.Infof("Reading models from %s source %s", source.Type, source.Id)

		registerFunc, ok := registeredModelProviders[source.Type]
		if !ok {
			glog.Errorf("catalog type %s not registered", source.Type)
			continue
		}

		records, err := registerFunc(ctx, &source, configDir)
		if err != nil {
			glog.Errorf("error reading catalog type %s with id %s: %v", source.Type, source.Id, err)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range records {
				// Set source_id on every returned model.
				l.setModelSourceID(r.Model, source.Id)

				ch <- r
			}
		}()
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	return ch
}

func (l *loader) setModelSourceID(model dbmodels.CatalogModel, sourceID string) {
	if model == nil {
		return
	}

	// Add a source_id property to the model's properties list.. the hard
	// way, because we use pointers to slices for some reason.

	props := model.GetProperties()
	if props == nil {
		if modelImpl, ok := model.(*dbmodels.CatalogModelImpl); ok {
			newProps := make([]mrmodels.Properties, 0, 1)
			modelImpl.Properties = &newProps
			props = &newProps
		} else {
			// Can't do anything with this.
			return
		}
	}

	for _, property := range *props {
		if property.Name == "source_id" {
			// Already has a source_id, just update it
			property.StringValue = &sourceID
			return
		}
	}

	*props = append(*props, mrmodels.NewStringProperty("source_id", sourceID, false))
}
