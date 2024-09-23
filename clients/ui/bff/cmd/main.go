package main

import (
	"flag"
	"fmt"
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
	flag.IntVar(&cfg.Port, "port", getEnvAsInt("PORT", 4000), "API server port")
	flag.BoolVar(&cfg.MockK8Client, "mock-k8s-client", false, "Use mock Kubernetes client")
	flag.BoolVar(&cfg.MockMRClient, "mock-mr-client", false, "Use mock Model Registry client")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

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

	logger.Info("starting server", "addr", srv.Addr)

	err = srv.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
	}
	os.Exit(1)
}

func getEnvAsInt(name string, defaultVal int) int {
	if value, exists := os.LookupEnv(name); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultVal
}
