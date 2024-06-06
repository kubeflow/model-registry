package integrations

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesClientInterface interface {
	FetchServiceNamesByComponent(componentValue string) ([]string, error)
}

type KubernetesClient struct {
	ClientSet *kubernetes.Clientset
	Namespace string
}

func NewKubernetesClient() (KubernetesClientInterface, error) {
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

	return &KubernetesClient{ClientSet: clientSet, Namespace: namespace}, nil
}

func (kc *KubernetesClient) FetchServiceNamesByComponent(componentValue string) ([]string, error) {

	services, err := kc.ClientSet.CoreV1().Services(kc.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var serviceNames []string
	for _, service := range services.Items {
		if svcComponent, exists := service.Spec.Selector["component"]; exists && svcComponent == componentValue {
			serviceNames = append(serviceNames, service.Name)
		}
	}
	return serviceNames, nil
}
