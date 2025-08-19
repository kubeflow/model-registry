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
	"k8s.io/client-go/kubernetes"
)

type SharedClientLogic struct {
	Client kubernetes.Interface
	Logger *slog.Logger
	Token  BearerToken
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

func (kc *SharedClientLogic) GetServiceDetailsByName(sessionCtx context.Context, namespace string, serviceName string) (ServiceDetails, error) {
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

func (kc *SharedClientLogic) GetDatabaseSecretValue(ctx context.Context, namespace, secretName, key string) (string, error) {
	secret, err := kc.Client.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		kc.Logger.Error("failed to get secret", "namespace", namespace, "secret", secretName, "error", err)
		return "", fmt.Errorf("failed to get secret %q: %w", secretName, err)
	}

	val, ok := secret.Data[key]
	if !ok {
		kc.Logger.Warn("key not found in secret", "secret", secretName, "key", key)
		return "", fmt.Errorf("key %q not found in secret %q", key, secretName)
	}

	return string(val), nil
}

func (kc *SharedClientLogic) CreateDatabaseSecret(ctx context.Context, name string, namespace string, database string, databaseUsername string, databasePassword string, dryRun bool) (*corev1.Secret, error) {
	kc.Logger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development until we have a definition on how to create model registries")
	return &corev1.Secret{}, nil
}

func (kc *SharedClientLogic) GetModelRegistrySettings(ctx context.Context, namespace string, labelSelector string) ([]unstructured.Unstructured, error) {
	kc.Logger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development until we have a definition on how to create model registries")
	return []unstructured.Unstructured{
		newUnstructuredModelRegistry("model-registry", namespace),
		newUnstructuredModelRegistry("model-registry-dora", namespace),
		newUnstructuredModelRegistry("model-registry-bella", namespace),
	}, nil
}

func (kc *SharedClientLogic) GetModelRegistrySettingsByName(ctx context.Context, namespace string, name string) (unstructured.Unstructured, error) {
	kc.Logger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development until we have a definition on how to create model registries")
	return newUnstructuredModelRegistry("model-registry", namespace), nil
}

func (kc *SharedClientLogic) CreateModelRegistryKind(ctx context.Context, namespace string, modelRegistryKind unstructured.Unstructured, dryRun bool) (unstructured.Unstructured, error) {
	kc.Logger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development until we have a definition on how to create model registries")

	// No actual logic; just simulate success
	return newUnstructuredModelRegistry("model-registry", namespace), nil
}

// This function is a temporary function to create a sample model registry kind until we have a real implementation
func newUnstructuredModelRegistry(name string, namespace string) unstructured.Unstructured {
	creationTime, _ := time.Parse(time.RFC3339, "2024-03-14T08:01:42Z")
	lastTransitionTime, _ := time.Parse(time.RFC3339, "2024-03-22T09:30:02Z")

	return unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "modelregistry.io/v1alpha1",
			"kind":       "ModelRegistry",
			"metadata": map[string]interface{}{
				"name":              name,
				"namespace":         namespace,
				"creationTimestamp": creationTime,
				"annotations":       map[string]interface{}{},
			},
			"spec": map[string]interface{}{
				"grpc": map[string]interface{}{},
				"rest": map[string]interface{}{},
				"istio": map[string]interface{}{
					"gateway": map[string]interface{}{
						"grpc": map[string]interface{}{
							"tls": map[string]interface{}{},
						},
						"rest": map[string]interface{}{
							"tls": map[string]interface{}{},
						},
					},
				},
				"databaseConfig": map[string]interface{}{
					"databaseType":                "mysql",
					"database":                    "model-registry",
					"host":                        "model-registry-db",
					"port":                        5432,
					"skipDBCreation":              false,
					"username":                    "mlmduser",
					"sslRootCertificateConfigMap": "ssl-config-map",
					"sslRootCertificateSecret":    "ssl-secret",
					// PasswordSecret intentionally omitted
				},
			},
			"status": map[string]interface{}{
				"conditions": []interface{}{
					map[string]interface{}{
						"type":               "Progressing",
						"status":             "True",
						"reason":             "CreatedDeployment",
						"message":            "Deployment for custom resource " + name + " was successfully created",
						"lastTransitionTime": lastTransitionTime,
					},
				},
			},
		},
	}
}
