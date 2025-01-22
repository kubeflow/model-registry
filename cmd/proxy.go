package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/internal/mlmdtypes"
	"github.com/kubeflow/model-registry/internal/server/openapi"
	"github.com/kubeflow/model-registry/pkg/core"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	glog.Infof("proxy server started at %s:%v", cfg.Hostname, cfg.Port)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	routers := make([]openapi.Router, 0)

	disableService := proxyCfg.DisableService
	if _, ok := map[string]struct{}{"": {}, CatalogService: {}, RegistryService: {}}[disableService]; !ok {
		return fmt.Errorf("invalid disable-service: %v", disableService)
	}
	if disableService != CatalogService {

		// TODO read yaml catalog file and instantiate ModelCatalogAPI implementations
		ModelCatalogServiceAPIService := openapi.NewModelCatalogServiceAPIService(map[string]openapi.ModelCatalogApi{})
		ModelCatalogServiceAPIController := openapi.NewModelCatalogServiceAPIController(ModelCatalogServiceAPIService)
		routers = append(routers, ModelCatalogServiceAPIController)
		glog.Infof("started catalog service")
	}
	if disableService != RegistryService {

		mlmdAddr := fmt.Sprintf("%s:%d", proxyCfg.MLMDHostname, proxyCfg.MLMDPort)
		glog.Infof("connecting to MLMD server %s..", mlmdAddr)
		conn, err := grpc.DialContext( // nolint:staticcheck
			ctxTimeout,
			mlmdAddr,
			grpc.WithReturnConnectionError(), // nolint:staticcheck
			grpc.WithBlock(),                 // nolint:staticcheck
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return fmt.Errorf("error dialing connection to mlmd server %s: %v", mlmdAddr, err)
		}
		defer conn.Close()
		glog.Infof("connected to MLMD server")

		mlmdTypeNamesConfig := mlmdtypes.NewMLMDTypeNamesConfigFromDefaults()
		_, err = mlmdtypes.CreateMLMDTypes(conn, mlmdTypeNamesConfig)
		if err != nil {
			return fmt.Errorf("error creating MLMD types: %v", err)
		}
		service, err := core.NewModelRegistryService(conn, mlmdTypeNamesConfig)
		if err != nil {
			return fmt.Errorf("error creating core service: %v", err)
		}

		// TODO make registry API optional to support standalone Catalog deployments
		ModelRegistryServiceAPIService := openapi.NewModelRegistryServiceAPIService(service)
		ModelRegistryServiceAPIController := openapi.NewModelRegistryServiceAPIController(ModelRegistryServiceAPIService)
		routers = append(routers, ModelRegistryServiceAPIController)

	}

	router := openapi.NewRouter(routers...)

	glog.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port), router))
	return nil
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	proxyCmd.Flags().StringVarP(&cfg.Hostname, "hostname", "n", cfg.Hostname, "Proxy server listen hostname")
	proxyCmd.Flags().IntVarP(&cfg.Port, "port", "p", cfg.Port, "Proxy server listen port")

	proxyCmd.Flags().StringVar(&proxyCfg.MLMDHostname, "mlmd-hostname", proxyCfg.MLMDHostname, "MLMD hostname")
	proxyCmd.Flags().IntVar(&proxyCfg.MLMDPort, "mlmd-port", proxyCfg.MLMDPort, "MLMD port")

	proxyCmd.Flags().StringVar(&proxyCfg.DisableService, "disable-service", proxyCfg.DisableService, "Optional name of service/endpoint to disable, can be either \"catalog\" or \"registry\"")
}

const (
	CatalogService  string = "catalog"
	RegistryService string = "registry"
)

type ProxyConfig struct {
	MLMDHostname string
	MLMDPort     int

	DisableService string
}

var proxyCfg = ProxyConfig{
	MLMDHostname: "localhost",
	MLMDPort:     9090,
}
