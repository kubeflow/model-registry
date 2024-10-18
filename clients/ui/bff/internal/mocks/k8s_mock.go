package mocks

import (
	"context"
	"fmt"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

type KubernetesClientMock struct {
	*k8s.KubernetesClient
	testEnv *envtest.Environment
}

func (m *KubernetesClientMock) Shutdown(ctx context.Context, logger *slog.Logger) error {
	logger.Info("Shutdown was called in mock")
	m.StopFn()
	err := m.testEnv.Stop()
	if err != nil {
		logger.Error("timeout while waiting for Kubernetes manager to stop")
		return fmt.Errorf("timeout while waiting for Kubernetes manager to stop")
	}
	logger.Info("Shutdown ended successfully")
	return nil
}

func NewKubernetesClient(logger *slog.Logger) (k8s.KubernetesClientInterface, error) {
	ctx, cancel := context.WithCancel(context.Background())
	projectRoot, err := getProjectRoot()
	if err != nil {
		logger.Error("failed to find project root to locate binaries", err)
		cancel()
		os.Exit(1)
	}

	testEnv := &envtest.Environment{
		// The BinaryAssetsDirectory is only required if you want to run the tests directly without call the makefile target test.
		// If not informed it will look for the default path defined in bff which is /usr/local/kubebuilder/.
		// Note that you must have the required binaries setup under the bin directory to perform the tests directly.
		// When we run make test it will be setup and used automatically.
		BinaryAssetsDirectory: filepath.Join(projectRoot, "bin", "k8s", fmt.Sprintf("1.29.0-%s-%s", runtime.GOOS, runtime.GOARCH)),
	}
	cfg, err := testEnv.Start()
	if err != nil {
		logger.Error("failed to start test environment", err)
		cancel()
		os.Exit(1)
	}

	mockK8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		logger.Error("failed to create Kubernetes client", err)
		cancel()
		os.Exit(1)
	}

	err = setupMock(mockK8sClient, ctx)
	if err != nil {
		logger.Error("failed on mock setup", err)
		cancel()
		os.Exit(1)
	}

	return &KubernetesClientMock{
		KubernetesClient: &k8s.KubernetesClient{
			Client: mockK8sClient,
			Logger: logger,
			StopFn: cancel,
		},
		testEnv: testEnv,
	}, nil
}

func getProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			// Found the project root where go.mod is located
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// We reached the root directory and did not find the go.mod
			return "", fmt.Errorf("could not find project root")
		}

		currentDir = parentDir
	}
}

func setupMock(mockK8sClient client.Client, ctx context.Context) error {
	err := createService(mockK8sClient, ctx, "model-registry", "default", "Model Registry", "Model Registry Description", "10.0.0.10")
	if err != nil {
		return err
	}
	err = createService(mockK8sClient, ctx, "model-registry-dora", "default", "Model Registry Dora", "Model Registry Dora description", "10.0.0.11")
	if err != nil {
		return err
	}
	err = createService(mockK8sClient, ctx, "model-registry-bella", "default", "Model Registry Bella", "Model Registry Bella description", "10.0.0.12")
	if err != nil {
		return err
	}
	return nil
}

func (m *KubernetesClientMock) GetServiceDetails() ([]k8s.ServiceDetails, error) {
	originalServices, err := m.KubernetesClient.GetServiceDetails()
	if err != nil {
		return nil, fmt.Errorf("failed to get service details: %w", err)
	}

	for i := range originalServices {
		originalServices[i].ClusterIP = "127.0.0.1"
		originalServices[i].HTTPPort = 8080
	}

	return originalServices, nil
}

func (m *KubernetesClientMock) GetServiceDetailsByName(serviceName string) (k8s.ServiceDetails, error) {
	originalService, err := m.KubernetesClient.GetServiceDetailsByName(serviceName)
	if err != nil {
		return k8s.ServiceDetails{}, fmt.Errorf("failed to get service details: %w", err)
	}
	//changing from cluster service ip to localhost
	originalService.ClusterIP = "127.0.0.1"
	originalService.HTTPPort = 8080

	return originalService, nil
}

func (m *KubernetesClientMock) BearerToken() (string, error) {
	return "FAKE BEARER TOKEN", nil
}

func createService(k8sClient client.Client, ctx context.Context, name string, namespace string, displayName string, description string, clusterIP string) error {

	annotations := map[string]string{}

	if displayName != "" {
		annotations["displayName"] = displayName
	}

	if description != "" {
		annotations["description"] = description
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"component": k8s.ComponentName,
			},
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: clusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:        "http-api",
					Port:        8080,
					Protocol:    corev1.ProtocolTCP,
					AppProtocol: strPtr("http"),
				},
				{
					Name:        "grpc-api",
					Port:        9090,
					Protocol:    corev1.ProtocolTCP,
					AppProtocol: strPtr("grpc"),
				},
			},
		},
	}

	err := k8sClient.Create(ctx, service)
	if err != nil {
		return fmt.Errorf("failed to create services: %w", err)
	}

	serviceList := &corev1.ServiceList{}

	err = k8sClient.List(ctx, serviceList, &client.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list services: %w", err)
	}

	if err != nil {
		return err
	}
	return nil
}

func strPtr(s string) *string {
	return &s
}
