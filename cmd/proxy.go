package cmd

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/server/openapi"
	"github.com/spf13/cobra"
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Starts the go OpenAPI proxy server to connect to a metadata store",
	Long: `This command launches the go OpenAPI proxy server.

The server connects to a metadata store, currently only MLMD is supported. It supports options to customize the
hostname and port where it listens.`,
	RunE: runProxyServer,
}

func runProxyServer(cmd *cobra.Command, args []string) error {
	glog.Infof("proxy server started at %s:%v", cfg.Hostname, cfg.Port)

	ds, dsTeardownF, err := datastore.NewDatastore(proxyCfg.DatastoreType, proxyCfg.DatastoreHostname, proxyCfg.DatastorePort)
	if err != nil {
		return fmt.Errorf("error creating datastore: %w", err)
	}

	defer func() {
		if err := dsTeardownF(); err != nil {
			glog.Errorf("error during cleanup: %w", err)
		}
	}()

	ModelRegistryServiceAPIService := openapi.NewModelRegistryServiceAPIService(ds)
	ModelRegistryServiceAPIController := openapi.NewModelRegistryServiceAPIController(ModelRegistryServiceAPIService)

	router := openapi.NewRouter(ModelRegistryServiceAPIController)

	glog.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port), router))
	return nil
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	proxyCmd.Flags().StringVarP(&cfg.Hostname, "hostname", "n", cfg.Hostname, "Proxy server listen hostname")
	proxyCmd.Flags().IntVarP(&cfg.Port, "port", "p", cfg.Port, "Proxy server listen port")

	proxyCmd.Flags().StringVar(&proxyCfg.DatastoreHostname, "mlmd-hostname", proxyCfg.DatastoreHostname, "MLMD hostname")
	if err := proxyCmd.Flags().MarkDeprecated("mlmd-hostname", "please use --datastore-hostname instead"); err != nil {
		glog.Errorf("error marking flag as deprecated: %v", err)
	}

	proxyCmd.Flags().IntVar(&proxyCfg.DatastorePort, "mlmd-port", proxyCfg.DatastorePort, "MLMD port")
	if err := proxyCmd.Flags().MarkDeprecated("mlmd-port", "please use --datastore-port instead"); err != nil {
		glog.Errorf("error marking flag as deprecated: %v", err)
	}

	proxyCmd.Flags().StringVar(&proxyCfg.DatastoreHostname, "datastore-hostname", proxyCfg.DatastoreHostname, "Datastore hostname")
	proxyCmd.Flags().IntVar(&proxyCfg.DatastorePort, "datastore-port", proxyCfg.DatastorePort, "Datastore port")
	proxyCmd.Flags().StringVar(&proxyCfg.DatastoreType, "datastore-type", proxyCfg.DatastoreType, "Datastore type")
}

type ProxyConfig struct {
	DatastoreHostname string
	DatastorePort     int
	DatastoreType     string
}

var proxyCfg = ProxyConfig{
	DatastoreHostname: "localhost",
	DatastorePort:     9090,
	DatastoreType:     "mlmd",
}
