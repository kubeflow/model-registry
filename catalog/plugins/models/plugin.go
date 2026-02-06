// Package models provides the model catalog plugin for the unified catalog server.
// This plugin wraps the existing catalog internals and exposes them via the plugin interface.
package models

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"reflect"
	"sync/atomic"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/catalog/internal/server/openapi"
	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd"
	"github.com/kubeflow/model-registry/pkg/catalog/plugin"
)

const (
	// PluginName is the identifier for this plugin.
	PluginName = "models"

	// PluginVersion is the API version.
	PluginVersion = "v1alpha1"
)

// ModelCatalogPlugin implements the CatalogPlugin interface for model catalogs.
type ModelCatalogPlugin struct {
	cfg       plugin.Config
	logger    *slog.Logger
	loader    *catalog.Loader
	dbCatalog catalog.APIProvider
	services  service.Services
	sources   *catalog.SourceCollection
	labels    *catalog.LabelCollection
	healthy   atomic.Bool
	started   atomic.Bool
}

// Name returns the plugin name.
func (p *ModelCatalogPlugin) Name() string {
	return PluginName
}

// Version returns the plugin API version.
func (p *ModelCatalogPlugin) Version() string {
	return PluginVersion
}

// Description returns a human-readable description.
func (p *ModelCatalogPlugin) Description() string {
	return "Model catalog for ML models"
}

// BasePath returns the API base path for this plugin.
func (p *ModelCatalogPlugin) BasePath() string {
	return "/api/model_catalog/v1alpha1"
}

// Init initializes the plugin with configuration.
func (p *ModelCatalogPlugin) Init(ctx context.Context, cfg plugin.Config) error {
	p.cfg = cfg
	p.logger = cfg.Logger
	if p.logger == nil {
		p.logger = slog.Default()
	}

	p.logger.Info("initializing model catalog plugin")

	// Build paths from config sources
	// The paths are the config file origins that contain the sources
	paths := make([]string, 0)
	originPaths := make(map[string]bool)

	for _, src := range cfg.Section.Sources {
		if src.Origin != "" {
			originPath := src.Origin
			if !originPaths[originPath] {
				paths = append(paths, originPath)
				originPaths[originPath] = true
			}
		}
	}

	// If no origins found in sources, use ConfigPaths
	if len(paths) == 0 {
		paths = cfg.ConfigPaths
	}

	// Convert paths to absolute
	absPaths := make([]string, 0, len(paths))
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			absPath = path
		}
		absPaths = append(absPaths, absPath)
	}

	// Initialize the services from the database connection
	services, err := p.initServices(cfg.DB)
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}
	p.services = services

	// Create the loader with existing catalog code
	p.loader = catalog.NewLoader(services, absPaths)
	p.sources = p.loader.Sources
	p.labels = p.loader.Labels

	// Create the DB catalog provider
	p.dbCatalog = catalog.NewDBCatalog(services, p.sources)

	p.logger.Info("model catalog plugin initialized", "paths", absPaths)

	return nil
}

// initServices creates the service layer from the database connection.
func (p *ModelCatalogPlugin) initServices(db *gorm.DB) (service.Services, error) {
	// Get the datastore spec for the catalog
	spec := service.DatastoreSpec()

	// We need to create the RepoSet from the existing database
	// This requires the types to already be registered in the database
	repoSet, err := p.createRepoSet(db, spec)
	if err != nil {
		return service.Services{}, fmt.Errorf("failed to create repo set: %w", err)
	}

	// Extract repositories from the RepoSet
	catalogModelRepo, err := getRepository[models.CatalogModelRepository](repoSet)
	if err != nil {
		return service.Services{}, fmt.Errorf("failed to get catalog model repository: %w", err)
	}

	catalogArtifactRepo, err := getRepository[models.CatalogArtifactRepository](repoSet)
	if err != nil {
		return service.Services{}, fmt.Errorf("failed to get catalog artifact repository: %w", err)
	}

	catalogModelArtifactRepo, err := getRepository[models.CatalogModelArtifactRepository](repoSet)
	if err != nil {
		return service.Services{}, fmt.Errorf("failed to get catalog model artifact repository: %w", err)
	}

	catalogMetricsArtifactRepo, err := getRepository[models.CatalogMetricsArtifactRepository](repoSet)
	if err != nil {
		return service.Services{}, fmt.Errorf("failed to get catalog metrics artifact repository: %w", err)
	}

	catalogSourceRepo, err := getRepository[models.CatalogSourceRepository](repoSet)
	if err != nil {
		return service.Services{}, fmt.Errorf("failed to get catalog source repository: %w", err)
	}

	propertyOptionsRepo, err := getRepository[models.PropertyOptionsRepository](repoSet)
	if err != nil {
		return service.Services{}, fmt.Errorf("failed to get property options repository: %w", err)
	}

	return service.NewServices(
		catalogModelRepo,
		catalogArtifactRepo,
		catalogModelArtifactRepo,
		catalogMetricsArtifactRepo,
		catalogSourceRepo,
		propertyOptionsRepo,
	), nil
}

// createRepoSet creates a RepoSet from the database using the spec.
// This uses the embedmd connector logic to initialize repositories.
func (p *ModelCatalogPlugin) createRepoSet(db *gorm.DB, spec *datastore.Spec) (datastore.RepoSet, error) {
	// Create a connector that uses the existing database
	connector, err := datastore.NewConnector("embedmd", &embedmd.EmbedMDConfig{DB: db})
	if err != nil {
		return nil, fmt.Errorf("failed to create connector: %w", err)
	}

	return connector.Connect(spec)
}

// Start begins background operations (hot-reload, watchers).
func (p *ModelCatalogPlugin) Start(ctx context.Context) error {
	p.logger.Info("starting model catalog plugin")

	if err := p.loader.Start(ctx); err != nil {
		return fmt.Errorf("failed to start loader: %w", err)
	}

	p.started.Store(true)
	p.healthy.Store(true)

	p.logger.Info("model catalog plugin started")
	return nil
}

// Stop gracefully shuts down the plugin.
func (p *ModelCatalogPlugin) Stop(ctx context.Context) error {
	p.logger.Info("stopping model catalog plugin")
	p.started.Store(false)
	p.healthy.Store(false)
	return nil
}

// Healthy returns true if the plugin is functioning correctly.
func (p *ModelCatalogPlugin) Healthy() bool {
	return p.healthy.Load()
}

// RegisterRoutes mounts the plugin's HTTP routes on the provided router.
func (p *ModelCatalogPlugin) RegisterRoutes(router chi.Router) error {
	p.logger.Info("registering model catalog routes")

	// Create the OpenAPI service using existing handlers
	apiService := openapi.NewModelCatalogServiceAPIService(
		p.dbCatalog,
		p.sources,
		p.labels,
		p.services.CatalogSourceRepository,
	)

	// Create the controller
	apiController := openapi.NewModelCatalogServiceAPIController(apiService)

	// Mount routes - remove the base path prefix since chi.Router already handles that
	for _, route := range apiController.OrderedRoutes() {
		// Remove the /api/model_catalog/v1alpha1 prefix from the pattern
		pattern := route.Pattern
		basePath := "/api/model_catalog/v1alpha1"
		if len(pattern) > len(basePath) && pattern[:len(basePath)] == basePath {
			pattern = pattern[len(basePath):]
		}
		if pattern == "" {
			pattern = "/"
		}

		router.Method(route.Method, pattern, route.HandlerFunc)
		p.logger.Debug("registered route", "method", route.Method, "pattern", pattern)
	}

	return nil
}

// Migrations returns database migrations for this plugin.
func (p *ModelCatalogPlugin) Migrations() []plugin.Migration {
	// The model catalog uses the existing database schema from embedmd
	// No additional migrations are needed as the schema is managed by the datastore layer
	return nil
}

// getRepository extracts a repository of type T from the RepoSet.
func getRepository[T any](rs datastore.RepoSet) (T, error) {
	var zero T
	t := reflect.TypeFor[T]()

	repo, err := rs.Repository(t)
	if err != nil {
		return zero, err
	}

	result, ok := repo.(T)
	if !ok {
		return zero, fmt.Errorf("repository type mismatch: expected %T, got %T", zero, repo)
	}

	return result, nil
}
