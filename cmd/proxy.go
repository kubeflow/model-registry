package cmd

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/golang/glog"

	"github.com/kubeflow/model-registry/internal/datastore"
	"github.com/kubeflow/model-registry/internal/proxy"
	"github.com/kubeflow/model-registry/internal/server/openapi"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/spf13/cobra"
)

type ProxyConfig struct {
	DatastoreHostname string
	DatastorePort     int
	DatastoreType     string
}

const (
	// datastoreUnavailableMessage is the message returned when the datastore service is down or unavailable.
	datastoreUnavailableMessage = "Datastore service is down or unavailable. Please check that the database is reachable and try again later."
)

var proxyCfg = ProxyConfig{
	DatastoreHostname: "localhost",
	DatastorePort:     9090,
	DatastoreType:     "mlmd",
}

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
	var (
		dsTeardownF datastore.TeardownFunc
		wg          sync.WaitGroup
	)

	router := proxy.NewDynamicRouter()

	router.SetRouter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, datastoreUnavailableMessage, http.StatusServiceUnavailable)
	}))

	errChan := make(chan error, 1)

	wg.Add(2)

	go func() {
		defer close(errChan)
		wg.Wait()
	}()

	// Start the connection to the Datastore server in a separate goroutine, so that
	// we can start the proxy server and start serving requests while we wait
	// for the connection to be established.
	go func() {
		var (
			ds  api.ModelRegistryApi
			err error
		)

		defer wg.Done()

		ds, dsTeardownF, err = datastore.NewDatastore(proxyCfg.DatastoreType, proxyCfg.DatastoreHostname, proxyCfg.DatastorePort)
		if err != nil {
			errChan <- fmt.Errorf("error creating datastore: %w", err)
			return
		}

		ModelRegistryServiceAPIService := openapi.NewModelRegistryServiceAPIService(ds)
		ModelRegistryServiceAPIController := openapi.NewModelRegistryServiceAPIController(ModelRegistryServiceAPIService)

		router.SetRouter(openapi.NewRouter(ModelRegistryServiceAPIController))

		glog.Infof("connected to Datastore server")
	}()

	// Start the proxy server in a separate goroutine so that we can handle
	// errors from both the proxy server and the connection to the Datastore server.
	go func() {
		defer wg.Done()

		glog.Infof("proxy server started at %s:%v", cfg.Hostname, cfg.Port)

		err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port), router)
		if err != nil {
			errChan <- fmt.Errorf("error starting proxy server: %w", err)
		}
	}()

	defer func() {
		if dsTeardownF != nil {
			glog.Info("closing connection to datastore service")

			//nolint:errcheck
			dsTeardownF()
		}
	}()

	// Wait for either the Datastore server connection or the proxy server to return an error
	// or for both to finish successfully.
	return <-errChan
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
