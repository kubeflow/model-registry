package integrations

import (
	"context"
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log/slog"
)

type KubernetesClientInterface interface {
	GetServiceNames() ([]string, error)
	GetServiceDetailsByName(serviceName string) (ServiceDetails, error)
	GetServiceDetails() ([]ServiceDetails, error)
	BearerToken() (string, error)
}

type ServiceDetails struct {
	Name        string
	DisplayName string
	Description string
	ClusterIP   string
	HTTPPort    int32
}

type KubernetesClient struct {
	ClientSet *kubernetes.Clientset
	Namespace string
	Token     string
	//TODO (ederign) How and on which frequency should we update this cache?
	//dont forget about mutexes
	ServiceCache map[string]ServiceDetails
}

func (kc *KubernetesClient) BearerToken() (string, error) {

	return kc.Token, nil
}

func NewKubernetesClient(logger *slog.Logger) (KubernetesClientInterface, error) {
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientcmd.NewDefaultClientConfigLoadingRules(), &clientcmd.ConfigOverrides{})
	restConfig, err := kubeConfig.ClientConfig()

	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes restConfig: %w", err)
	}

	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes namespace: %w", err)
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes clientset: %w", err)
	}
	//fetching services
	services, err := clientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list model-registry-server services: %w", err)
	}

	//building serviceCache
	serviceCache, err := buildModelRegistryServiceCache(logger, *services)
	if err != nil {
		return nil, err
	}

	kc := &KubernetesClient{
		ClientSet:    clientSet,
		Namespace:    namespace,
		Token:        restConfig.BearerToken,
		ServiceCache: serviceCache,
	}

	return kc, nil
}

func buildModelRegistryServiceCache(logger *slog.Logger, services v1.ServiceList) (map[string]ServiceDetails, error) {
	serviceCache := make(map[string]ServiceDetails)
	for _, service := range services.Items {
		if svcComponent, exists := service.Spec.Selector["component"]; exists && svcComponent == "model-registry-server" {
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
				logger.Error("service missing HTTP port", "serviceName", service.Name)
				continue
			}
			if service.Spec.ClusterIP == "" {
				logger.Error("service missing valid ClusterIP", "serviceName", service.Name)
				continue
			}

			//TODO (acreasy) DisplayName and Description need to be included and not given a zero value once we
			// know how this will be implemented.
			serviceCache[service.Name] = ServiceDetails{
				Name:      service.Name,
				ClusterIP: service.Spec.ClusterIP,
				HTTPPort:  httpPort,
			}
		}
	}
	return serviceCache, nil
}

func (kc *KubernetesClient) GetServiceNames() ([]string, error) {
	//TODO (ederign) when we develop the front-end, implement subject access review here
	// and check if the username has actually permissions to access that server
	// currently on kf dashboard, the user name comes in kubeflow-userid

	var serviceNames []string

	for _, service := range kc.ServiceCache {
		if service.Name != "" {
			serviceNames = append(serviceNames, service.Name)
		}
	}
	return serviceNames, nil
}

func (kc *KubernetesClient) GetServiceDetails() ([]ServiceDetails, error) {
	var services []ServiceDetails

	for _, service := range kc.ServiceCache {
		if service.Name != "" {
			services = append(services, ServiceDetails{
				Name:        service.Name,
				DisplayName: service.DisplayName,
				Description: service.Description,
			})
		}
	}
	return services, nil
}

func (kc *KubernetesClient) GetServiceDetailsByName(serviceName string) (ServiceDetails, error) {

	service, exists := kc.ServiceCache[serviceName]
	if !exists {
		return ServiceDetails{}, fmt.Errorf("service %s not found in cache", serviceName)
	}

	return service, nil
}
