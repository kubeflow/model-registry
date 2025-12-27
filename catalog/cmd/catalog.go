package cmd

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/catalog/internal/mcp"
	"github.com/kubeflow/model-registry/catalog/internal/server/openapi"
	openapimodel "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd"
	"github.com/spf13/cobra"
)

var catalogCfg = struct {
	ListenAddress          string
	ConfigPath             []string
	McpCatalogPath         []string
	PerformanceMetricsPath []string
}{
	ListenAddress:          "0.0.0.0:8080",
	ConfigPath:             []string{"sources.yaml"},
	McpCatalogPath:         []string{},
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
	fs.StringSliceVar(&catalogCfg.McpCatalogPath, "mcp-catalogs-path", catalogCfg.McpCatalogPath, "Path to MCP catalog source configuration file")
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
		getRepo[models.CatalogSourceRepository](repoSet),
		getRepo[models.PropertyOptionsRepository](repoSet),
		getRepo[models.McpServerRepository](repoSet),
	)

	loader := catalog.NewLoader(services, catalogCfg.ConfigPath)

	perfLoader, err := catalog.NewPerformanceMetricsLoader(catalogCfg.PerformanceMetricsPath, services.CatalogModelRepository, services.CatalogMetricsArtifactRepository, repoSet.TypeMap())
	if err != nil {
		return fmt.Errorf("error initializing performance metrics: %v", err)
	}
	loader.RegisterEventHandler(perfLoader.Load)

	poRefresher := models.NewPropertyOptionsRefresher(context.Background(), services.PropertyOptionsRepository, time.Second)
	loader.RegisterEventHandler(func(ctx context.Context, record catalog.ModelProviderRecord) error {
		poRefresher.Trigger()
		return nil
	})

	err = loader.Start(context.Background())
	if err != nil {
		return fmt.Errorf("error loading catalog sources: %v", err)
	}

	svc := openapi.NewModelCatalogServiceAPIService(
		catalog.NewDBCatalog(services, loader.Sources),
		loader.Sources,
		loader.Labels,
		services.CatalogSourceRepository,
	)
	ctrl := openapi.NewModelCatalogServiceAPIController(svc)

	// Initialize MCP Catalog service
	// Always use database-backed provider - if no sources configured, returns empty list
	var mcpLoader *catalog.McpLoader
	if len(catalogCfg.McpCatalogPath) > 0 {
		// Load MCP servers from YAML sources into database
		// Pass the shared SourceCollection so MCP sources appear in unified /sources API
		mcpLoader = catalog.NewMcpLoader(services, catalogCfg.McpCatalogPath, loader.Sources)
		err = mcpLoader.Start(context.Background())
		if err != nil {
			return fmt.Errorf("error loading MCP catalog sources: %v", err)
		}
		glog.Infof("MCP catalog loaded from %d source(s)", len(catalogCfg.McpCatalogPath))
	} else {
		glog.Infof("No MCP catalog sources configured (use --mcp-catalogs-path to specify)")
	}
	mcpProvider := mcp.NewDbMcpCatalogProvider(services.McpServerRepository)
	// Set named query resolver from MCP loader (or shared SourceCollection)
	if mcpLoader != nil {
		mcpProvider.SetNamedQueryResolver(func() map[string]map[string]openapimodel.FieldFilter {
			namedQueries := mcpLoader.GetNamedQueries()
			// Convert catalog.FieldFilter to openapimodel.FieldFilter
			result := make(map[string]map[string]openapimodel.FieldFilter, len(namedQueries))
			for queryName, fieldFilters := range namedQueries {
				result[queryName] = make(map[string]openapimodel.FieldFilter, len(fieldFilters))
				for fieldName, filter := range fieldFilters {
					result[queryName][fieldName] = openapimodel.FieldFilter{
						Operator: filter.Operator,
						Value:    filter.Value,
					}
				}
			}
			return result
		})
	}
	mcpSvc := openapi.NewMcpCatalogServiceAPIService(mcpProvider)
	mcpCtrl := openapi.NewMcpCatalogServiceAPIController(mcpSvc)

	glog.Infof("Catalog API server listening on %s", catalogCfg.ListenAddress)
	return http.ListenAndServe(catalogCfg.ListenAddress, openapi.NewRouter(ctrl, mcpCtrl))
}

func getRepo[T any](repoSet datastore.RepoSet) T {
	repo, err := repoSet.Repository(reflect.TypeFor[T]())
	if err != nil {
		panic(fmt.Sprintf("unable to get repository: %v", err))
	}

	return repo.(T)
}
