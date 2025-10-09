package cmd

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/catalog/internal/server/openapi"
	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd"
	"github.com/spf13/cobra"
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
	ds, err := datastore.NewConnector("embedmd", &embedmd.EmbedMDConfig{
		DatabaseType: "postgres", // We only support postgres right now
		DatabaseDSN:  "",         // Empty DSN, see https://www.postgresql.org/docs/current/libpq-envars.html
	})
	if err != nil {
		return fmt.Errorf("error creating datastore: %w", err)
	}

	repoSet, err := ds.Connect(service.DatastoreSpec())
	if err != nil {
		return fmt.Errorf("error initializing datastore: %v", err)
	}

	services := service.NewServices(
		getRepo[models.CatalogModelRepository](repoSet),
		getRepo[models.CatalogArtifactRepository](repoSet),
		getRepo[models.CatalogModelArtifactRepository](repoSet),
		getRepo[models.CatalogMetricsArtifactRepository](repoSet),
	)

	sources, err := catalog.LoadCatalogSources(context.Background(), services, catalogCfg.ConfigPath)
	if err != nil {
		return fmt.Errorf("error loading catalog sources: %v", err)
	}

	modelRepo := getRepo[models.CatalogModelRepository](repoSet)
	artifactRepo := getRepo[models.CatalogArtifactRepository](repoSet)
	metricsArtifactRepo := getRepo[models.CatalogMetricsArtifactRepository](repoSet)

	svc := openapi.NewModelCatalogServiceAPIService(catalog.NewDBCatalog(
		modelRepo,
		artifactRepo,
	), sources)
	ctrl := openapi.NewModelCatalogServiceAPIController(svc)

	err = catalog.LoadPerformanceMetricsData(catalogCfg.PerformanceMetricsPath, modelRepo, metricsArtifactRepo, repoSet.TypeMap())
	if err != nil {
		return fmt.Errorf("error loading performance metrics data: %v", err)
	}

	glog.Infof("Catalog API server listening on %s", catalogCfg.ListenAddress)
	return http.ListenAndServe(catalogCfg.ListenAddress, openapi.NewRouter(ctrl))
}

func getRepo[T any](repoSet datastore.RepoSet) T {
	repo, err := repoSet.Repository(reflect.TypeFor[T]())
	if err != nil {
		panic(fmt.Sprintf("unable to get repository: %v", err))
	}

	return repo.(T)
}
