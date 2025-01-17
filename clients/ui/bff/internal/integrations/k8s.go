package integrations

import (
	"context"
	"fmt"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	authv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log/slog"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"time"
)

const ComponentLabelValue = "model-registry"

type KubernetesClientInterface interface {
	GetServiceNames(sessionCtx context.Context, namespace string) ([]string, error)
	GetServiceDetailsByName(sessionCtx context.Context, namespace string, serviceName string) (ServiceDetails, error)
	GetServiceDetails(sessionCtx context.Context, namespace string) ([]ServiceDetails, error)
	BearerToken() (string, error)
	Shutdown(ctx context.Context, logger *slog.Logger) error
	IsInCluster() bool
	PerformSARonGetListServicesByNamespace(user string, groups []string, namespace string) (bool, error)
	PerformSARonSpecificService(user string, groups []string, namespace string, serviceName string) (bool, error)
	IsClusterAdmin(user string) (bool, error)
	GetNamespaces(user string, groups []string) ([]corev1.Namespace, error)
}

type ServiceDetails struct {
	Name        string
	DisplayName string
	Description string
	ClusterIP   string
	HTTPPort    int32
}

type KubernetesClient struct {
	ControllerRuntimeClient client.Client        //Controller-runtime client: used for high-level operations with caching.
	KubernetesNativeClient  kubernetes.Interface //Native KubernetesNativeClient client: only for specific non-cached subresources like SAR.
	Mgr                     ctrl.Manager
	Token                   string
	Logger                  *slog.Logger
	StopFn                  context.CancelFunc // Store a function to cancel the context for graceful shutdown
	mgrStopped              chan struct{}
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

	//Native KubernetesNativeClient client: only for specific non-cached subresources like SAR.
	k8sClient, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		logger.Error("failed to create native KubernetesNativeClient client", "error", err)
		cancel()
		return nil, fmt.Errorf("failed to create KubernetesNativeClient client: %w", err)
	}

	kc := &KubernetesClient{
		ControllerRuntimeClient: mgr.GetClient(),
		KubernetesNativeClient:  k8sClient,
		Mgr:                     mgr,
		Token:                   kubeconfig.BearerToken,
		Logger:                  logger,
		StopFn:                  cancel,
		mgrStopped:              mgrStopped,
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

func (kc *KubernetesClient) IsInCluster() bool {
	_, err := rest.InClusterConfig()
	return err == nil
}

func (kc *KubernetesClient) BearerToken() (string, error) {
	return kc.Token, nil
}

func (kc *KubernetesClient) GetServiceNames(sessionCtx context.Context, namespace string) ([]string, error) {
	services, err := kc.GetServiceDetails(sessionCtx, namespace)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(services))
	for _, svc := range services {
		names = append(names, svc.Name)
	}

	return names, nil
}

func (kc *KubernetesClient) GetServiceDetails(sessionCtx context.Context, namespace string) ([]ServiceDetails, error) {

	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sessionLogger := sessionCtx.Value(constants.TraceLoggerKey).(*slog.Logger)

	serviceList := &corev1.ServiceList{}

	labelSelector := labels.SelectorFromSet(labels.Set{
		"component": ComponentLabelValue,
	})

	err := kc.ControllerRuntimeClient.List(ctx, serviceList, &client.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     namespace,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var services []ServiceDetails

	for _, service := range serviceList.Items {
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
			sessionLogger.Error("service missing HTTP port", "serviceName", service.Name)
			continue
		}

		if service.Spec.ClusterIP == "" {
			sessionLogger.Error("service missing valid ClusterIP", "serviceName", service.Name)
			continue
		}

		displayName := ""
		description := ""

		if service.Annotations != nil {
			displayName = service.Annotations["displayName"]
			description = service.Annotations["description"]
		}

		if displayName == "" {
			sessionLogger.Warn("service missing displayName annotation", "serviceName", service.Name)
		}

		if description == "" {
			sessionLogger.Warn("service missing description annotation", "serviceName", service.Name)
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

	return services, nil
}

func (kc *KubernetesClient) GetServiceDetailsByName(sessionCtx context.Context, namespace string, serviceName string) (ServiceDetails, error) {
	services, err := kc.GetServiceDetails(sessionCtx, namespace)
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

func (kc *KubernetesClient) PerformSARonGetListServicesByNamespace(user string, groups []string, namespace string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	verbs := []string{"get", "list"}
	resource := "services"

	for _, verb := range verbs {
		sar := &authv1.SubjectAccessReview{
			Spec: authv1.SubjectAccessReviewSpec{
				User:   user,
				Groups: groups,
				ResourceAttributes: &authv1.ResourceAttributes{
					Verb:      verb,
					Resource:  resource,
					Namespace: namespace,
				},
			},
		}

		// Perform the SAR using the native KubernetesNativeClient client
		response, err := kc.KubernetesNativeClient.AuthorizationV1().SubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
		if err != nil {
			return false, fmt.Errorf("failed to create SubjectAccessReview for verb %q on resource %q: %w", verb, resource, err)
		}

		if !response.Status.Allowed {
			kc.Logger.Warn("access denied", "user", user, "verb", verb, "resource", resource)
			return false, nil
		}
	}

	return true, nil
}

func (kc *KubernetesClient) IsClusterAdmin(user string) (bool, error) {
	//using a context here, because checking ClusterRoleBindings could be expensive in large clusters
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clusterRoleBindings := &rbacv1.ClusterRoleBindingList{}
	err := kc.ControllerRuntimeClient.List(ctx, clusterRoleBindings)
	if err != nil {
		return false, fmt.Errorf("failed to list ClusterRoleBindings: %w", err)
	}

	for _, crb := range clusterRoleBindings.Items {
		if crb.RoleRef.Name != "cluster-admin" {
			continue
		}
		for _, subject := range crb.Subjects {

			if subject.Kind == "User" && subject.Name == user {
				return true, nil
			}
		}
	}

	return false, nil
}

func (kc *KubernetesClient) GetNamespaces(user string, groups []string) ([]corev1.Namespace, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//list all namespaces
	namespaceList := &corev1.NamespaceList{}
	err := kc.ControllerRuntimeClient.List(ctx, namespaceList)
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	//check user access with SAR for each namespace
	var namespaces []corev1.Namespace
	for _, ns := range namespaceList.Items {
		sar := &authv1.SubjectAccessReview{
			Spec: authv1.SubjectAccessReviewSpec{
				User:   user,
				Groups: groups,
				ResourceAttributes: &authv1.ResourceAttributes{
					Namespace: ns.Name,
					Verb:      "get",
					Resource:  "namespaces",
				},
			},
		}

		response, err := kc.KubernetesNativeClient.AuthorizationV1().SubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
		if err != nil {
			kc.Logger.Error("failed to perform SubjectAccessReview", "namespace", ns.Name, "error", err)
			continue
		}

		if response.Status.Allowed {
			namespaces = append(namespaces, ns)
		}
	}

	return namespaces, nil

}

func (kc *KubernetesClient) PerformSARonSpecificService(user string, groups []string, namespace string, serviceName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resource := "services"
	verb := "get"

	sar := &authv1.SubjectAccessReview{
		Spec: authv1.SubjectAccessReviewSpec{
			User:   user,
			Groups: groups,
			ResourceAttributes: &authv1.ResourceAttributes{
				Verb:      verb,
				Resource:  resource,
				Namespace: namespace,
				Name:      serviceName,
			},
		},
	}

	// Perform the SAR using the native KubernetesNativeClient client
	response, err := kc.KubernetesNativeClient.AuthorizationV1().SubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
	if err != nil {
		return false, fmt.Errorf(
			"failed to create SubjectAccessReview for verb %q on resource %q (service: %q) in namespace %q: %w",
			verb, resource, serviceName, namespace, err,
		)
	}

	if !response.Status.Allowed {
		kc.Logger.Warn(
			"access denied", "user", user, "verb", verb, "resource", resource,
			"namespace", namespace, "service", serviceName,
		)
		return false, nil
	}

	return true, nil
}
