package integrations

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

const ComponentName = "model-registry-server"

type KubernetesClientInterface interface {
	GetServiceNames() ([]string, error)
	GetServiceDetailsByName(serviceName string) (ServiceDetails, error)
	GetServiceDetails() ([]ServiceDetails, error)
	BearerToken() (string, error)
	Shutdown(ctx context.Context, logger *slog.Logger) error
}

type ServiceDetails struct {
	Name        string
	DisplayName string
	Description string
	ClusterIP   string
	HTTPPort    int32
}

type KubernetesClient struct {
	Client     client.Client
	Mgr        ctrl.Manager
	Token      string
	Logger     *slog.Logger
	StopFn     context.CancelFunc // Store a function to cancel the context for graceful shutdown
	mgrStopped chan struct{}
}

func NewKubernetesClient(logger *slog.Logger) (KubernetesClientInterface, error) {
	// Create a context with a cancel function is used for shutdown the kubernetes client
	ctx, cancel := context.WithCancel(ctrl.SetupSignalHandler())

	kubeconfig, err := helper.GetKubeconfig()
	if err != nil {
		logger.Error("failed to get kubeconfig", "error", err)
		os.Exit(1)
	}

	scheme, err := helper.BuildScheme()
	if err != nil {
		logger.Error("failed to build Kubernetes scheme", "error", err)
		os.Exit(1)
	}

	// Create the manager with caching capabilities
	mgr, err := ctrl.NewManager(kubeconfig, ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: "0", // disable metrics serving
		},
		HealthProbeBindAddress: "0", // disable health probe serving
		LeaderElection:         false,
		//Namespace:              "namespace", //TODO (ederign) do we need to specify the namespace to operate in
		//There is also cache filters and Sync periods to assess later.
	})

	if err != nil {
		logger.Error("unable to create manager", "error", err)
		cancel()
		os.Exit(1)
	}

	// Channel to signal when the manager has stopped
	mgrStopped := make(chan struct{})

	// Start the manager in a goroutine
	go func() {
		defer close(mgrStopped) // Signal that the manager has stopped
		if err := mgr.Start(ctx); err != nil {
			logger.Error("problem running manager", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for the cache to sync before using the client
	if !mgr.GetCache().WaitForCacheSync(ctx) {
		cancel()
		return nil, fmt.Errorf("failed to wait for cache to sync")
	}

	kc := &KubernetesClient{
		Client:     mgr.GetClient(),
		Mgr:        mgr,
		Token:      kubeconfig.BearerToken,
		Logger:     logger,
		StopFn:     cancel,
		mgrStopped: mgrStopped, // Store the stop channel

		//Namespace:    namespace, //TODO (ederign) do we need to restrict service list by namespace?
	}
	return kc, nil
}

func (kc *KubernetesClient) Shutdown(ctx context.Context, logger *slog.Logger) error {
	logger.Info("shutting down Kubernetes manager...")

	// Use the saved cancel function to stop the manager
	kc.StopFn()

	// Wait for the manager to stop or for the context to be canceled
	select {
	case <-kc.mgrStopped:
		logger.Info("Kubernetes manager stopped successfully")
		return nil
	case <-ctx.Done():
		logger.Error("context canceled while waiting for Kubernetes manager to stop")
		return ctx.Err()
	case <-time.After(30 * time.Second):
		logger.Error("timeout while waiting for Kubernetes manager to stop")
		return fmt.Errorf("timeout while waiting for Kubernetes manager to stop")
	}
}

func (kc *KubernetesClient) BearerToken() (string, error) {
	return kc.Token, nil
}

func (kc *KubernetesClient) GetServiceNames() ([]string, error) {
	//TODO (ederign) when we develop the front-end, implement subject access review here
	// and check if the username has actually permissions to access that server
	// currently on kf dashboard, the user name comes in kubeflow-userid

	//TODO (ederign) we should consider and rethinking listing all services on cluster
	// what if we have thousand of those?
	// we should consider label filtering for instance

	serviceList := &corev1.ServiceList{}
	//TODO (ederign) review the context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := kc.Client.List(ctx, serviceList, &client.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var serviceNames []string
	for _, service := range serviceList.Items {
		if value, ok := service.Spec.Selector["component"]; ok && value == ComponentName {
			serviceNames = append(serviceNames, service.Name)
		}
	}

	if len(serviceNames) == 0 {
		return nil, fmt.Errorf("no services found with component: %s", ComponentName)
	}

	return serviceNames, nil
}

func (kc *KubernetesClient) GetServiceDetails() ([]ServiceDetails, error) {
	//TODO (ederign) review the context timeout

	//TODO (ederign) when we develop the front-end, implement subject access review here
	// and check if the username has actually permissions to access that server
	// currently on kf dashboard, the user name comes in kubeflow-userid

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel() // Ensure the context is canceled to free up resources

	serviceList := &corev1.ServiceList{}

	err := kc.Client.List(ctx, serviceList, &client.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var services []ServiceDetails

	for _, service := range serviceList.Items {
		if svcComponent, exists := service.Spec.Selector["component"]; exists && svcComponent == ComponentName {
			var httpPort int32
			hasHTTPPort := false
			for _, port := range service.Spec.Ports {
				if port.Name == "http-api" {
					httpPort = port.Port
					hasHTTPPort = true
					break
				}
			}
			if !hasHTTPPort {
				kc.Logger.Error("service missing HTTP port", "serviceName", service.Name)
				continue
			}

			if service.Spec.ClusterIP == "" {
				kc.Logger.Error("service missing valid ClusterIP", "serviceName", service.Name)
				continue
			}

			displayName := ""
			description := ""

			if service.Annotations != nil {
				displayName = service.Annotations["displayName"]
				description = service.Annotations["description"]
			}

			if displayName == "" {
				kc.Logger.Warn("service missing displayName annotation", "serviceName", service.Name)
			}

			if description == "" {
				kc.Logger.Warn("service missing description annotation", "serviceName", service.Name)
			}

			serviceDetails := ServiceDetails{
				Name:        service.Name,
				DisplayName: displayName,
				Description: description,
				ClusterIP:   service.Spec.ClusterIP,
				HTTPPort:    httpPort,
			}

			services = append(services, serviceDetails)
		}
	}

	return services, nil
}

func (kc *KubernetesClient) GetServiceDetailsByName(serviceName string) (ServiceDetails, error) {
	//TODO (ederign) when we develop the front-end, implement subject access review here
	// and check if the username has actually permissions to access that server
	// currently on kf dashboard, the user name comes in kubeflow-userid

	services, err := kc.GetServiceDetails()
	if err != nil {
		return ServiceDetails{}, fmt.Errorf("failed to get service details: %w", err)
	}

	for _, service := range services {
		if service.Name == serviceName {
			return service, nil
		}
	}

	return ServiceDetails{}, fmt.Errorf("service %s not found", serviceName)
}
