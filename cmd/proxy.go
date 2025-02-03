package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/mlmdtypes"
	"github.com/kubeflow/model-registry/internal/proxy"
	"github.com/kubeflow/model-registry/internal/server/openapi"
	"github.com/kubeflow/model-registry/pkg/core"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	// mlmdUnavailableMessage is the message returned when the MLMD server is down or unavailable.
	mlmdUnavailableMessage = "MLMD server is down or unavailable. Please check that the database is reachable and try again later."
	// maxGRPCRetryAttempts is the maximum number of attempts to retry GRPC requests to the MLMD server.
	maxGRPCRetryAttempts = 25 // 25 attempts with incremental backoff (1s, 2s, 3s, ..., 25s) it's ~5 minutes
)

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

	errMLMDChan := make(chan error, 1)
	errProxyChan := make(chan error, 1)

	router := proxy.NewDynamicRouter()

	router.SetRouter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, mlmdUnavailableMessage, http.StatusServiceUnavailable)
	}))

	// Start the connection to the MLMD server in a separate goroutine, so that
	// we can start the proxy server and start serving requests while we wait
	// for the connection to be established.
	go func() {
		defer close(errMLMDChan)

		mlmdAddr := fmt.Sprintf("%s:%d", proxyCfg.MLMDHostname, proxyCfg.MLMDPort)
		glog.Infof("connecting to MLMD server %s..", mlmdAddr)
		conn, err = grpc.NewClient(mlmdAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			errMLMDChan <- fmt.Errorf("error dialing connection to mlmd server %s: %w", mlmdAddr, err)

			return
		}

		mlmdTypeNamesConfig := mlmdtypes.NewMLMDTypeNamesConfigFromDefaults()

		// Backoff and retry GRPC requests to the MLMD server, until the server
		// becomes available or the maximum number of attempts is reached.
		for i := 0; i < maxGRPCRetryAttempts; i++ {
			_, err := mlmdtypes.CreateMLMDTypes(conn, mlmdTypeNamesConfig)
			if err == nil {
				break
			}

			st, ok := status.FromError(err)
			if !ok || st.Code() != codes.Unavailable {
				errMLMDChan <- fmt.Errorf("error creating MLMD types: %w", err)

				return
			}

			time.Sleep(time.Duration(i+1) * time.Second)
		}

		service, err := core.NewModelRegistryService(conn, mlmdTypeNamesConfig)
		if err != nil {
			errMLMDChan <- fmt.Errorf("error creating core service: %w", err)

			return
		}

		ModelRegistryServiceAPIService := openapi.NewModelRegistryServiceAPIService(service)
		ModelRegistryServiceAPIController := openapi.NewModelRegistryServiceAPIController(ModelRegistryServiceAPIService)

		router.SetRouter(openapi.NewRouter(ModelRegistryServiceAPIController))

		glog.Infof("connected to MLMD server")
	}()

	// Start the proxy server in a separate goroutine so that we can handle
	// errors from both the proxy server and the connection to the MLMD server.
	go func() {
		defer close(errProxyChan)

		glog.Infof("proxy server started at %s:%v", cfg.Hostname, cfg.Port)

		err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port), router)
		if err != nil {
			errProxyChan <- fmt.Errorf("error starting proxy server: %w", err)
		}
	}()

	defer func() {
		if conn != nil {
			glog.Info("closing connection to MLMD server")

			conn.Close()
		}
	}()

	// Wait for either the MLMD server connection or the proxy server to return an error
	// or for both to finish successfully.
	for {
		select {
		case err := <-errMLMDChan:
			if err != nil {
				return err
			}

		case err := <-errProxyChan:
			if err != nil {
				return err
			}
		}

		if errMLMDChan == nil && errProxyChan == nil {
			return nil
		}
	}
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
