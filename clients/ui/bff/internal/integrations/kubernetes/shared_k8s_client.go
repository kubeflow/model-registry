package kubernetes

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

var modelRegistryGVR = schema.GroupVersionResource{
	Group:    "modelregistry.opendatahub.io",
	Version:  "v1beta1",
	Resource: "modelregistries",
}

type SharedClientLogic struct {
	Client        kubernetes.Interface
	DynamicClient dynamic.Interface
	Logger        *slog.Logger
	Token         BearerToken
}

func (kc *SharedClientLogic) GetServiceNames(sessionCtx context.Context, namespace string) ([]string, error) {
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

func (kc *SharedClientLogic) GetServiceDetails(sessionCtx context.Context, namespace string) ([]ServiceDetails, error) {

	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sessionLogger := sessionCtx.Value(constants.TraceLoggerKey).(*slog.Logger)

	labelSelector := fmt.Sprintf("component=%s", ComponentLabelValue)

	serviceList, err := kc.Client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var services []ServiceDetails

	for _, service := range serviceList.Items {
		serviceDetails, err := buildServiceDetails(&service, sessionLogger)
		if err != nil {
			sessionLogger.Warn("skipping service", "error", err)
			continue
		}
		services = append(services, *serviceDetails)

	}

	return services, nil
}

func buildServiceDetails(service *corev1.Service, logger *slog.Logger) (*ServiceDetails, error) {
	if service == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}

	var httpPort int32
	var isHTTPS bool
	hasHTTPPort := false
	for _, port := range service.Spec.Ports {
		if port.Name == "http-api" || port.Name == "https-api" {
			httpPort = port.Port
			isHTTPS = port.Name == "https-api"
			hasHTTPPort = true
			break
		}
	}
	if !hasHTTPPort {
		logger.Error("service missing HTTP/HTTPS port", "serviceName", service.Name)
		return nil, fmt.Errorf("service %q missing required 'http-api' or 'https-api' port", service.Name)
	}

	if service.Spec.ClusterIP == "" {
		logger.Error("service missing valid ClusterIP", "serviceName", service.Name)
		return nil, fmt.Errorf("service %q missing ClusterIP", service.Name)
	}

	displayName := ""
	description := ""
	externalAddressRest := ""

	// Check for annotations including external-address-rest
	if service.Annotations != nil {
		displayName = service.Annotations["displayName"]
		description = service.Annotations["description"]

		// Look for external-address-rest annotation with any prefix
		for key, value := range service.Annotations {
			if strings.HasSuffix(key, "/external-address-rest") {
				externalAddressRest = value
				break
			}
		}
	}

	return &ServiceDetails{
		Name:                service.Name,
		DisplayName:         displayName,
		Description:         description,
		ClusterIP:           service.Spec.ClusterIP,
		HTTPPort:            httpPort,
		IsHTTPS:             isHTTPS,
		ExternalAddressRest: externalAddressRest,
	}, nil
}

func (kc *SharedClientLogic) GetServiceDetailsByName(sessionCtx context.Context, namespace string, serviceName string, serviceType string) (ServiceDetails, error) {
	if namespace == "" || serviceName == "" {
		return ServiceDetails{}, fmt.Errorf("namespace and serviceName cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sessionLogger := sessionCtx.Value(constants.TraceLoggerKey).(*slog.Logger)

	service, err := kc.Client.CoreV1().Services(namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		return ServiceDetails{}, fmt.Errorf("failed to get service %q in namespace %q: %w", serviceName, namespace, err)
	}
	if serviceType != "" && service.Labels["component"] != serviceType {
		return ServiceDetails{}, fmt.Errorf("service %q in namespace %q is not a %s", serviceName, namespace, serviceType)
	}

	details, err := buildServiceDetails(service, sessionLogger)
	if err != nil {
		return ServiceDetails{}, err
	}
	return *details, nil
}

func (kc *SharedClientLogic) BearerToken() (string, error) {
	// Token is retained for follow-up calls; do not log it.
	return kc.Token.Raw(), nil
}

func (kc *SharedClientLogic) GetGroups(ctx context.Context) ([]string, error) {
	kc.Logger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development until we have a definition on how to create model registries")
	return []string{}, nil
}

// CreateModelRegistry creates a ModelRegistry custom resource in the specified namespace
func (kc *SharedClientLogic) CreateModelRegistry(ctx context.Context, namespace string, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}
	if obj == nil {
		return nil, fmt.Errorf("model registry object cannot be nil")
	}

	// Match existing pattern - use Background for timeout consistency
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	kc.Logger.Info("creating ModelRegistry CR", "namespace", namespace, "name", obj.GetName())

	result, err := kc.DynamicClient.Resource(modelRegistryGVR).Namespace(namespace).Create(
		ctxWithTimeout,
		obj,
		metav1.CreateOptions{},
	)
	if err != nil {
		kc.Logger.Error("failed to create ModelRegistry CR", "error", err, "namespace", namespace, "name", obj.GetName())
		return nil, fmt.Errorf("failed to create ModelRegistry CR: %w", err)
	}

	kc.Logger.Info("successfully created ModelRegistry CR", "namespace", namespace, "name", result.GetName())
	return result, nil
}

// CreateSecret creates a Kubernetes Secret in the specified namespace
func (kc *SharedClientLogic) CreateSecret(ctx context.Context, namespace string, secret *corev1.Secret) (*corev1.Secret, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}
	if secret == nil {
		return nil, fmt.Errorf("secret object cannot be nil")
	}

	// Match existing pattern - use Background for timeout consistency
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	kc.Logger.Info("creating Secret", "namespace", namespace, "name", secret.Name)

	result, err := kc.Client.CoreV1().Secrets(namespace).Create(
		ctxWithTimeout,
		secret,
		metav1.CreateOptions{},
	)
	if err != nil {
		kc.Logger.Error("failed to create Secret", "error", err, "namespace", namespace, "name", secret.Name)
		return nil, fmt.Errorf("failed to create Secret: %w", err)
	}

	kc.Logger.Info("successfully created Secret", "namespace", namespace, "name", result.Name)
	return result, nil
}
