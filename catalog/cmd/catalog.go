package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/catalog/internal/leader"
	"github.com/kubeflow/model-registry/catalog/internal/server/openapi"
	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd"
	"github.com/kubeflow/model-registry/internal/db"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var catalogCfg = struct {
	ListenAddress          string
	ConfigPath             []string
	PerformanceMetricsPath []string
}{
	ListenAddress:          "0.0.0.0:8080",
	ConfigPath:             []string{"sources.yaml"},
	PerformanceMetricsPath: []string{},
}

const (
	leaderLockName = "catalog-leader"

	defaultLeaderLockDuration = 60 * time.Second
	defaultLeaderHeartbeat    = 15 * time.Second

	envLeaderLockDuration = "CATALOG_LEADER_LOCK_DURATION"
	envLeaderHeartbeat    = "CATALOG_LEADER_HEARTBEAT"
)

// parseDurationEnv parses a duration from an environment variable,
// falling back to a default value if unset or invalid.
func parseDurationEnv(envName string, defaultVal time.Duration) time.Duration {
	if envVal := os.Getenv(envName); envVal != "" {
		if parsed, err := time.ParseDuration(envVal); err == nil {
			glog.Infof("Using %s: %v", envName, parsed)
			return parsed
		}
		glog.Warningf("Invalid %s value %q, using default %v", envName, envVal, defaultVal)
	}
	return defaultVal
}

// getLeaderElectionConfig reads leader election configuration from environment
// variables, falling back to defaults when unset or invalid.
func getLeaderElectionConfig() (lockDuration, heartbeat time.Duration) {
	lockDuration = parseDurationEnv(envLeaderLockDuration, defaultLeaderLockDuration)
	heartbeat = parseDurationEnv(envLeaderHeartbeat, defaultLeaderHeartbeat)

	// Validate pglock requirement: heartbeat <= lockDuration/2
	if heartbeat > lockDuration/2 {
		glog.Warningf("Heartbeat (%v) exceeds half of lock duration (%v), required by pglock. Using defaults.", heartbeat, lockDuration)
		return defaultLeaderLockDuration, defaultLeaderHeartbeat
	}

	return lockDuration, heartbeat
}

var CatalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "Catalog API server",
	Long: `Launch the API server for the model catalog. Use PostgreSQL's
	environment variables
	(https://www.postgresql.org/docs/current/libpq-envars.html) to
	configure the database connection.`,
	RunE: runCatalogServer,
}

func init() {
	fs := CatalogCmd.Flags()
	fs.StringVarP(&catalogCfg.ListenAddress, "listen", "l", catalogCfg.ListenAddress, "Address to listen on")
	fs.StringSliceVar(&catalogCfg.ConfigPath, "catalogs-path", catalogCfg.ConfigPath, "Path to catalog source configuration file")
	fs.StringSliceVar(&catalogCfg.PerformanceMetricsPath, "performance-metrics", catalogCfg.PerformanceMetricsPath, "Path to performance metrics data directory")
}

func runCatalogServer(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := db.Init(
		"postgres", // We only support postgres right now
		"",         // Empty DSN, see https://www.postgresql.org/docs/current/libpq-envars.html
		nil,        // Default TLS config
	)
	if err != nil {
		return fmt.Errorf("error creating datastore: %w", err)
	}
	gormDB, err := db.GetConnector().Connect()
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	ds, err := datastore.NewConnector("embedmd", &embedmd.EmbedMDConfig{
		DB:                gormDB,
		WaitForMigrations: true,
	})
	if err != nil {
		return fmt.Errorf("error creating datastore: %w", err)
	}

	lockDuration, heartbeat := getLeaderElectionConfig()
	glog.Infof("Leader election configured: lock duration=%v, heartbeat=%v", lockDuration, heartbeat)

	elector, err := leader.NewLeaderElector(gormDB, ctx, leaderLockName, lockDuration, heartbeat)
	if err != nil {
		return fmt.Errorf("error creating leader elector: %w", err)
	}
	elector.OnBecomeLeader(func(leaderCtx context.Context) {
		err := ds.RunMigrations(service.DatastoreSpec())
		if err != nil {
			glog.Errorf("unable to run migrations: %v", err)
		}
	})

	repoSet, err := ds.Connect(service.DatastoreSpec())
	if err != nil {
		return fmt.Errorf("error initializing datastore: %v", err)
	}

	services := service.NewServices(
		getRepo[models.CatalogModelRepository](repoSet),
		getRepo[models.CatalogArtifactRepository](repoSet),
		getRepo[models.CatalogModelArtifactRepository](repoSet),
		getRepo[models.CatalogMetricsArtifactRepository](repoSet),
		getRepo[models.CatalogSourceRepository](repoSet),
		getRepo[models.PropertyOptionsRepository](repoSet),
	)

	loader := catalog.NewLoader(services, catalogCfg.ConfigPath)

	perfLoader, err := catalog.NewPerformanceMetricsLoader(catalogCfg.PerformanceMetricsPath, services.CatalogModelRepository, services.CatalogMetricsArtifactRepository, repoSet.TypeMap())
	if err != nil {
		return fmt.Errorf("error initializing performance metrics: %v", err)
	}
	loader.RegisterEventHandler(perfLoader.Load)

	elector.OnBecomeLeader(func(leaderCtx context.Context) {
		poRefresher := models.NewPropertyOptionsRefresher(leaderCtx, services.PropertyOptionsRepository, time.Second)
		loader.RegisterEventHandler(func(ctx context.Context, record catalog.ModelProviderRecord) error {
			poRefresher.Trigger()
			return nil
		})
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		glog.Infof("Received signal %v, initiating graceful shutdown", sig)
		cancel()
	}()

	glog.Info("Starting loader in read-only mode (standby)")
	if err := loader.StartReadOnly(ctx); err != nil {
		return fmt.Errorf("error starting loader in read-only mode: %v", err)
	}

	// Set up HTTP server (runs continuously regardless of leadership)
	svc := openapi.NewModelCatalogServiceAPIService(
		catalog.NewDBCatalog(services, loader.Sources),
		loader.Sources,
		loader.Labels,
		services.CatalogSourceRepository,
	)
	ctrl := openapi.NewModelCatalogServiceAPIController(svc)

	server := &http.Server{
		Addr:    catalogCfg.ListenAddress,
		Handler: openapi.NewRouter(ctrl),
	}

	g, gctx := errgroup.WithContext(ctx)

	// HTTP server goroutine
	g.Go(func() error {
		glog.Infof("Catalog API server listening on %s", catalogCfg.ListenAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("HTTP server failed: %w", err)
		}
		return nil
	})

	// HTTP server shutdown goroutine
	g.Go(func() error {
		<-gctx.Done()
		glog.Info("Shutting down HTTP server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			glog.Errorf("HTTP server shutdown error: %v", err)
		}
		return nil
	})

	elector.OnBecomeLeader(func(leaderCtx context.Context) {
		glog.Info("Became leader - starting leader-only operations")

		// StartLeader blocks until leaderCtx is cancelled (leadership lost)
		if err := loader.StartLeader(leaderCtx); err != nil && !errors.Is(err, context.Canceled) {
			glog.Errorf("StartLeader exited with error: %v", err)
		}

		glog.Info("Leader callback complete")
	})

	g.Go(func() error {
		if err := elector.Wait(); err != nil {
			return fmt.Errorf("leader elector failed: %w", err)
		}
		return nil
	})

	// Wait for all goroutines and collect errors
	errs := []error{}
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		errs = append(errs, err)
	}

	if err := loader.Shutdown(); err != nil {
		errs = append(errs, fmt.Errorf("loader shutdown error: %w", err))
	}

	return errors.Join(errs...)
}

func getRepo[T any](repoSet datastore.RepoSet) T {
	repo, err := repoSet.Repository(reflect.TypeFor[T]())
	if err != nil {
		panic(fmt.Sprintf("unable to get repository: %v", err))
	}

	return repo.(T)
}
