package cmd

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	"github.com/kubeflow/model-registry/catalog/internal/server/openapi"
	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/datastore/embedmd"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/spf13/cobra"
)

var catalogCfg = struct {
	ListenAddress string
	ConfigPath    []string
}{
	ListenAddress: "0.0.0.0:8080",
	ConfigPath:    []string{"sources.yaml"},
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
}

func runCatalogServer(cmd *cobra.Command, args []string) error {
	ds, err := datastore.NewConnector("embedmd", &embedmd.EmbedMDConfig{
		DatabaseType: "postgres", // We only support postgres right now
		DatabaseDSN:  "",         // Empty DSN, see https://www.postgresql.org/docs/current/libpq-envars.html
	})
	if err != nil {
		return fmt.Errorf("error creating datastore: %w", err)
	}

	_, err = ds.Connect(datastore.RepoSetSpec{
		ArtifactTypes: map[string]any{
			defaults.ModelArtifactTypeName: service.NewModelArtifactRepository,
			defaults.DocArtifactTypeName:   service.NewDocArtifactRepository,
		},
		ContextTypes:   map[string]any{},
		ExecutionTypes: map[string]any{},
	})
	if err != nil {
		return fmt.Errorf("error initializing datastore: %v", err)
	}

	sources, err := catalog.LoadCatalogSources(catalogCfg.ConfigPath)
	if err != nil {
		return fmt.Errorf("error loading catalog sources: %v", err)
	}

	svc := openapi.NewModelCatalogServiceAPIService(sources)
	ctrl := openapi.NewModelCatalogServiceAPIController(svc)

	glog.Infof("Catalog API server listening on %s", catalogCfg.ListenAddress)
	return http.ListenAndServe(catalogCfg.ListenAddress, openapi.NewRouter(ctrl))
}
