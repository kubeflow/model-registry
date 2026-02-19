// Package plugin provides a plugin-based architecture for catalog services.
// Catalog types (models, datasets, etc.) register as plugins via init() and
// are mounted under a unified HTTP server.
package plugin

import (
	"context"
	"flag"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// CatalogPlugin defines the interface that all catalog plugins must implement.
// Plugins register themselves via init() using the Register function.
type CatalogPlugin interface {
	// Identity returns the plugin name (e.g., "models", "datasets").
	// This name is used for routing and configuration lookup.
	Name() string

	// Version returns the API version (e.g., "v1alpha1").
	Version() string

	// Description returns a human-readable description of the plugin.
	Description() string

	// Init initializes the plugin with its configuration.
	// Called once during server startup before Start.
	Init(ctx context.Context, cfg Config) error

	// Start begins background operations (hot-reload, watchers, etc.).
	// Called after Init and after database migrations.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the plugin.
	// Called during server shutdown.
	Stop(ctx context.Context) error

	// Healthy returns true if the plugin is functioning correctly.
	// Used for health check endpoints.
	Healthy() bool

	// RegisterRoutes mounts the plugin's HTTP routes on the provided router.
	// The router is already scoped to the plugin's base path.
	RegisterRoutes(router chi.Router) error

	// Migrations returns database migrations for this plugin.
	// Migrations are applied in order during server initialization.
	Migrations() []Migration
}

// BasePathProvider is an optional interface that plugins can implement
// to specify their own API base path. If not implemented, the server
// computes it as /api/{name}_catalog/{version}.
type BasePathProvider interface {
	BasePath() string
}

// SourceKeyProvider is an optional interface that plugins can implement
// to specify which key in the sources.yaml "catalogs" map they respond to.
// If not implemented, the plugin name is used as the config key.
// This allows the plugin name and config key to differ (e.g., plugin "model"
// can read from the "models" config section).
type SourceKeyProvider interface {
	SourceKey() string
}

// CatalogLoader defines the interface for data loading strategies.
// The core Loader[E, A] implements this by default. Plugins can
// register multiple loaders (e.g., core + custom).
type CatalogLoader interface {
	// Start begins loading data and sets up any watchers/background operations.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the loader.
	Stop(ctx context.Context) error
}

// FlagProvider is an optional interface that plugins can implement
// to register custom CLI flags before flag parsing.
type FlagProvider interface {
	// RegisterFlags registers custom CLI flags for this plugin.
	// Called before flag.Parse() during server startup.
	RegisterFlags(fs *flag.FlagSet)
}

// Migration represents a database migration for a plugin.
type Migration struct {
	// Version is a unique identifier for this migration (e.g., "001", "20240101_initial").
	Version string

	// Description provides a human-readable description of what this migration does.
	Description string

	// Up applies the migration.
	Up func(db *gorm.DB) error

	// Down reverts the migration.
	Down func(db *gorm.DB) error
}

// Config is passed to each plugin during Init.
type Config struct {
	// Section contains the plugin-specific configuration from sources.yaml.
	Section CatalogSection

	// DB is the shared database connection.
	DB *gorm.DB

	// Logger is a namespaced logger for this plugin.
	Logger *slog.Logger

	// BasePath is the API base path for this plugin (e.g., "/api/models_catalog/v1alpha1").
	BasePath string

	// ConfigPaths are the paths to all sources.yaml files being used.
	ConfigPaths []string
}
