package cmd

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/catalog"
	openapi "github.com/kubeflow/model-registry/internal/server/openapi/catalog"
	"github.com/spf13/cobra"
)

var catalogCfg = struct {
	ListenAddress string
	ConfigPath    string
}{
	ListenAddress: "0.0.0.0:8080",
	ConfigPath:    "sources.yaml",
}

var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "Catalog API server",
	Long:  `Launch the API server for the model catalog`,
	RunE:  runCatalogServer,
}

func init() {
	rootCmd.AddCommand(catalogCmd)

	catalogCmd.Flags().StringVarP(&catalogCfg.ListenAddress, "listen", "l", catalogCfg.ListenAddress, "Address to listen on")
	catalogCmd.Flags().StringVar(&catalogCfg.ConfigPath, "catalogs-path", catalogCfg.ConfigPath, "Path to catalog source configuration file")
}

func runCatalogServer(cmd *cobra.Command, args []string) error {
	sources, err := catalog.LoadCatalogSources(catalogCfg.ConfigPath)
	if err != nil {
		return fmt.Errorf("error loading catalog sources: %v", err)
	}

	svc := openapi.NewModelCatalogServiceAPIService(sources)
	ctrl := openapi.NewModelCatalogServiceAPIController(svc)

	glog.Infof("Catalog API server listening on %s", catalogCfg.ListenAddress)
	return http.ListenAndServe(catalogCfg.ListenAddress, openapi.NewRouter(ctrl))
}
