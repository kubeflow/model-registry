package cmd

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/opendatahub-io/model-registry/internal/server/openapi"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

var (
	// proxyCmd represents the proxy command
	proxyCmd = &cobra.Command{
		Use:   "proxy",
		Short: "Starts the ml-metadata go OpenAPI proxy",
		Long: `This command launches the ml-metadata go OpenAPI proxy server.

The server connects to a mlmd CPP server. It supports options to customize the 
hostname and port where it listens.'`,
		RunE: runProxyServer,
	}
)

func runProxyServer(cmd *cobra.Command, args []string) error {
	glog.Infof("proxy server started at %s:%v", cfg.Hostname, cfg.Port)

	ModelRegistryServiceAPIService := openapi.NewModelRegistryServiceAPIService()
	ModelRegistryServiceAPIController := openapi.NewModelRegistryServiceAPIController(ModelRegistryServiceAPIService)

	router := openapi.NewRouter(ModelRegistryServiceAPIController)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port), router))
	return nil
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// proxyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	proxyCmd.Flags().StringVarP(&cfg.Hostname, "hostname", "n", cfg.Hostname, "Proxy server listen hostname")
	proxyCmd.Flags().IntVarP(&cfg.Port, "port", "p", cfg.Port, "Proxy server listen port")
}
