package cmd

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/golang/glog"
	mrGrpc "github.com/kubeflow/model-registry/internal/grpc"
	"github.com/kubeflow/model-registry/internal/mlmdtypes"
	"github.com/kubeflow/model-registry/internal/proxy"
	"github.com/kubeflow/model-registry/internal/server/openapi"
	"github.com/kubeflow/model-registry/pkg/core"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const mlmdUnavailableMessage = "MLMD server is down or unavailable. Please check that the database is reachable and try again later."

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Starts the ml-metadata go OpenAPI proxy",
	Long: `This command launches the ml-metadata go OpenAPI proxy server.

The server connects to a mlmd CPP server. It supports options to customize the
hostname and port where it listens.'`,
	RunE: runProxyServer,
}

func runProxyServer(cmd *cobra.Command, args []string) error {
	var conn *grpc.ClientConn
	var err error
	var wg sync.WaitGroup

	errChan := make(chan error, 1)

	router := proxy.NewDynamicRouter()

	router.SetRouter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, mlmdUnavailableMessage, http.StatusServiceUnavailable)
	}))

	// Start the connection to the MLMD server in a separate goroutine, so that
	// we can start the proxy server and start serving requests while we wait
	// for the connection to be established.
	go func() {
		defer close(errChan)

		mlmdAddr := fmt.Sprintf("%s:%d", proxyCfg.MLMDHostname, proxyCfg.MLMDPort)
		glog.Infof("connecting to MLMD server %s..", mlmdAddr)
		conn, err = grpc.NewClient(mlmdAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			errChan <- fmt.Errorf("error dialing connection to mlmd server %s: %w", mlmdAddr, err)

			return
		}

		mlmdTypeNamesConfig := mlmdtypes.NewMLMDTypeNamesConfigFromDefaults()

		_, err = mrGrpc.RetryOnGRPCError[map[string]int64](mlmdtypes.CreateMLMDTypes, conn, mlmdTypeNamesConfig)
		if err != nil {
			errChan <- fmt.Errorf("error creating MLMD types: %w", err)

			return
		}
		service, err := core.NewModelRegistryService(conn, mlmdTypeNamesConfig)
		if err != nil {
			errChan <- fmt.Errorf("error creating core service: %w", err)

			return
		}

		ModelRegistryServiceAPIService := openapi.NewModelRegistryServiceAPIService(service)
		ModelRegistryServiceAPIController := openapi.NewModelRegistryServiceAPIController(ModelRegistryServiceAPIService)

		router.SetRouter(openapi.NewRouter(ModelRegistryServiceAPIController))

		glog.Infof("connected to MLMD server")
	}()

	wg.Add(1)

	// Start the proxy server in a separate goroutine so that we can handle
	// errors from both the proxy server and the connection to the MLMD server.
	go func() {
		defer wg.Done()

		glog.Infof("proxy server started at %s:%v", cfg.Hostname, cfg.Port)

		err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port), router)
		if err != nil {
			errChan <- err
		}
	}()

	defer func() {
		if conn != nil {
			glog.Info("closing connection to MLMD server")

			conn.Close()
		}
	}()

	err = <-errChan
	if err != nil {
		return fmt.Errorf("error starting proxy server: %w", err)
	}

	// Wait for the proxy server to finish serving requests.
	wg.Wait()

	return nil
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	proxyCmd.Flags().StringVarP(&cfg.Hostname, "hostname", "n", cfg.Hostname, "Proxy server listen hostname")
	proxyCmd.Flags().IntVarP(&cfg.Port, "port", "p", cfg.Port, "Proxy server listen port")

	proxyCmd.Flags().StringVar(&proxyCfg.MLMDHostname, "mlmd-hostname", proxyCfg.MLMDHostname, "MLMD hostname")
	proxyCmd.Flags().IntVar(&proxyCfg.MLMDPort, "mlmd-port", proxyCfg.MLMDPort, "MLMD port")
}

type ProxyConfig struct {
	MLMDHostname string
	MLMDPort     int
}

var proxyCfg = ProxyConfig{
	MLMDHostname: "localhost",
	MLMDPort:     9090,
}
