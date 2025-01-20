package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kubeflow/model-registry/ui/bff/internal/api"
	"github.com/kubeflow/model-registry/ui/bff/internal/config"

	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	var cfg config.EnvConfig
	flag.IntVar(&cfg.Port, "port", getEnvAsInt("PORT", 8080), "API server port")
	flag.BoolVar(&cfg.MockK8Client, "mock-k8s-client", false, "Use mock Kubernetes client")
	flag.BoolVar(&cfg.MockMRClient, "mock-mr-client", false, "Use mock Model Registry client")
	flag.BoolVar(&cfg.DevMode, "dev-mode", false, "Use development mode for access to local K8s cluster")
	flag.IntVar(&cfg.DevModePort, "dev-mode-port", getEnvAsInt("DEV_MODE_PORT", 8080), "Use port when in development mode")
	flag.BoolVar(&cfg.StandaloneMode, "standalone-mode", false, "Use standalone mode for enabling endpoints in standalone mode")
	flag.StringVar(&cfg.StaticAssetsDir, "static-assets-dir", "./static", "Configure frontend static assets root directory")
	flag.StringVar(&cfg.LogLevel, "log-level", getEnvAsString("LOG_LEVEL", "info"), "Sets server log level, possible values: debug, info, warn, error, fatal")
	flag.StringVar(&cfg.AllowedOrigins, "allowed-origins", getEnvAsString("ALLOWED_ORIGINS", ""), "Sets allowed origins for CORS purposes, accepts a comma separated list of origins or * to allow all, default none")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLogLevelFromString(cfg.LogLevel),
	}))

	app, err := api.NewApp(cfg, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Start the server in a goroutine
	go func() {
		logger.Info("starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server ListenAndServe", "error", err)
		}
	}()

	// Graceful shutdown setup
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Wait for shutdown signal
	<-shutdownCh
	logger.Info("shutting down gracefully...")

	// Create a context with timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the HTTP server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
	}

	// Shutdown the Kubernetes manager gracefully
	if err := app.Shutdown(ctx, logger); err != nil {
		logger.Error("failed to shutdown Kubernetes manager", "error", err)
	}

	logger.Info("server stopped")
	os.Exit(0)

}

func getEnvAsInt(name string, defaultVal int) int {
	if value, exists := os.LookupEnv(name); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultVal
}

func getEnvAsString(name string, defaultVal string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	return defaultVal
}

func getLogLevelFromString(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "fatal":
		return slog.LevelError

	}
	return slog.LevelInfo
}
