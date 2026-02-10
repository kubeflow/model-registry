// Package main provides the unified catalog server entry point.
// This server hosts all registered catalog plugins under a single process.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"gorm.io/gorm"

	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd"
	"github.com/kubeflow/model-registry/internal/db"
	"github.com/kubeflow/model-registry/pkg/catalog/plugin"

	// Import plugins - their init() registers them
	_ "github.com/kubeflow/model-registry/catalog/plugins/model"
	_ "github.com/kubeflow/model-registry/catalog/plugins/mcp"
	// _ "github.com/kubeflow/model-registry/catalog/plugins/datasets"  // future
)

func main() {
	var (
		listenAddr   string
		sourcesPath  string
		databaseType string
		databaseDSN  string
	)

	flag.StringVar(&listenAddr, "listen", ":8080", "Address to listen on")
	flag.StringVar(&sourcesPath, "sources", "/config/sources.yaml", "Path to catalog sources config")
	flag.StringVar(&databaseType, "db-type", "postgres", "Database type (postgres or mysql)")
	flag.StringVar(&databaseDSN, "db-dsn", "", "Database connection string")
	flag.Parse()

	// Initialize glog for backwards compatibility
	_ = flag.Set("logtostderr", "true")

	// Set up structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("starting catalog server",
		"listen", listenAddr,
		"sources", sourcesPath,
		"plugins", plugin.Names(),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		logger.Info("received shutdown signal", "signal", sig)
		cancel()
	}()

	// Load config
	cfg, err := plugin.LoadConfig(sourcesPath)
	if err != nil {
		glog.Fatalf("Failed to load config: %v", err)
	}

	logger.Info("loaded config",
		"apiVersion", cfg.APIVersion,
		"kind", cfg.Kind,
		"catalogs", len(cfg.Catalogs),
	)

	// Setup database
	gormDB, err := setupDatabase(databaseType, databaseDSN)
	if err != nil {
		glog.Fatalf("Failed to connect to database: %v", err)
	}

	// Create and initialize server
	server := plugin.NewServer(cfg, []string{sourcesPath}, gormDB, logger)
	if err := server.Init(ctx); err != nil {
		glog.Fatalf("Failed to initialize plugins: %v", err)
	}

	// Mount routes and start
	router := server.MountRoutes()

	if err := server.Start(ctx); err != nil {
		glog.Fatalf("Failed to start plugins: %v", err)
	}

	logger.Info("catalog server ready",
		"listen", listenAddr,
		"plugins", plugin.Names(),
	)

	// Create HTTP server with graceful shutdown
	httpServer := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}

	// Start HTTP server in goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	logger.Info("shutting down...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", "error", err)
	}

	if err := server.Stop(shutdownCtx); err != nil {
		logger.Error("plugin shutdown error", "error", err)
	}

	logger.Info("catalog server stopped")
}

func setupDatabase(dbType, dsn string) (*gorm.DB, error) {
	if dsn == "" {
		// Try to get from environment
		dsn = os.Getenv("DATABASE_DSN")
		if dsn == "" {
			return nil, fmt.Errorf("database DSN is required (use -db-dsn flag or DATABASE_DSN environment variable)")
		}
	}

	if dbType == "" {
		dbType = os.Getenv("DATABASE_TYPE")
		if dbType == "" {
			dbType = "postgres"
		}
	}

	// Create embedmd connector
	cfg := &embedmd.EmbedMDConfig{
		DatabaseType: dbType,
		DatabaseDSN:  dsn,
	}

	connector, err := datastore.NewConnector("embedmd", cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connector: %w", err)
	}

	// Connect to initialize the database
	// We need a minimal spec just to establish the connection
	_, err = connector.Connect(datastore.NewSpec())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get the GORM DB from the db package
	dbConnector, ok := db.GetConnector()
	if !ok {
		return nil, fmt.Errorf("database connector not available")
	}

	gormDB, err := dbConnector.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to get GORM connection: %w", err)
	}

	return gormDB, nil
}
