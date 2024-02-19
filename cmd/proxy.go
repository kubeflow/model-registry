package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/opendatahub-io/model-registry/internal/mlmdtypes"
	"github.com/opendatahub-io/model-registry/internal/server/openapi"
	"github.com/opendatahub-io/model-registry/pkg/core"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mlmdAddr := fmt.Sprintf("%s:%d", proxyCfg.MLMDHostname, proxyCfg.MLMDPort)
	glog.Infof("connecting to MLMD server %s..", mlmdAddr)
	conn, err := grpc.DialContext(
		ctxTimeout,
		mlmdAddr,
		grpc.WithReturnConnectionError(),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("error dialing connection to mlmd server %s: %v", mlmdAddr, err)
	}
	defer conn.Close()
	glog.Infof("connected to MLMD server")

	_, err = mlmdtypes.CreateMLMDTypes(conn)
	if err != nil {
		return fmt.Errorf("error creating MLMD types: %v", err)
	}
	service, err := core.NewModelRegistryService(conn)
	if err != nil {
		return fmt.Errorf("error creating core service: %v", err)
	}

	ModelRegistryServiceAPIService := openapi.NewModelRegistryServiceAPIService(service)
	ModelRegistryServiceAPIController := openapi.NewModelRegistryServiceAPIController(ModelRegistryServiceAPIService)

	router := openapi.NewRouter(ModelRegistryServiceAPIController)

	glog.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port), router))
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
