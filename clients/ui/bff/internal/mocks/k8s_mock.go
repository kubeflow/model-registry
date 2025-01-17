package mocks

import (
	"context"
	"fmt"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

const (
	KubeflowUserIDHeaderValue = "user@example.com"
	DoraNonAdminUser          = "doraNonAdmin@example.com"
	BellaNonAdminUser         = "bellaNonAdmin@example.com"
	DoraServiceGroup          = "dora-service-group"
	DoraNamespaceGroup        = "dora-namespace-group"
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

func NewKubernetesClient(logger *slog.Logger, ctx context.Context, cancel context.CancelFunc) (k8s.KubernetesClientInterface, error) {

	projectRoot, err := getProjectRoot()
	if err != nil {
		logger.Error("failed to find project root to locate binaries", slog.String("error", err.Error()))
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
		logger.Error("failed to start test environment", slog.String("error", err.Error()))
		cancel()
		os.Exit(1)
	}

	mockK8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		logger.Error("failed to create Kubernetes client", slog.String("error", err.Error()))
		cancel()
		os.Exit(1)
	}

	nativeK8sClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logger.Error("failed to create native KubernetesNativeClient client", slog.String("error", err.Error()))
		cancel()
		os.Exit(1)
	}

	err = setupMock(mockK8sClient, ctx)
	if err != nil {
		logger.Error("failed on mock setup", slog.String("error", err.Error()))
		cancel()
		os.Exit(1)
	}

	return &KubernetesClientMock{
		KubernetesClient: &k8s.KubernetesClient{
			ControllerRuntimeClient: mockK8sClient,
			KubernetesNativeClient:  nativeK8sClient,
			Logger:                  logger,
			StopFn:                  cancel,
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

	err := createNamespace(mockK8sClient, ctx, "kubeflow")
	if err != nil {
		return err
	}

	err = createNamespace(mockK8sClient, ctx, "dora-namespace")
	if err != nil {
		return err
	}

	err = createNamespace(mockK8sClient, ctx, "bella-namespace")
	if err != nil {
		return err
	}

	err = createService(mockK8sClient, ctx, "model-registry", "kubeflow", "Model Registry", "Model Registry Description", "10.0.0.10", "model-registry")
	if err != nil {
		return err
	}
	err = createService(mockK8sClient, ctx, "model-registry-one", "kubeflow", "Model Registry One", "Model Registry One description", "10.0.0.11", "model-registry")
	if err != nil {
		return err
	}
	err = createService(mockK8sClient, ctx, "model-registry-dora", "dora-namespace", "Model Registry Dora", "Model Registry Dora description", "10.0.0.12", "model-registry")
	if err != nil {
		return err
	}
	err = createService(mockK8sClient, ctx, "model-registry-bella", "bella-namespace", "Model Registry Bella", "Model Registry Bella description", "10.0.0.13", "model-registry")
	if err != nil {
		return err
	}
	err = createService(mockK8sClient, ctx, "non-model-registry", "kubeflow", "Not a Model Registry", "Not a Model Registry Bella description", "10.0.0.14", "")
	if err != nil {
		return err
	}

	err = createClusterAdminRBAC(mockK8sClient, ctx, KubeflowUserIDHeaderValue)
	if err != nil {
		return fmt.Errorf("failed to create cluster admin RBAC: %w", err)
	}

	err = createNamespaceRestrictedRBAC(mockK8sClient, ctx, DoraNonAdminUser, "dora-namespace")
	if err != nil {
		return fmt.Errorf("failed to create namespace-restricted RBAC: %w", err)
	}

	err = createNamespaceRestrictedRBAC(mockK8sClient, ctx, BellaNonAdminUser, "bella-namespace")
	if err != nil {
		return fmt.Errorf("failed to create namespace-restricted RBAC: %w", err)
	}

	err = createGroupAccessRBAC(mockK8sClient, ctx, DoraServiceGroup, "dora-namespace", "model-registry-dora")
	if err != nil {
		return fmt.Errorf("failed to create group-based RBAC: %w", err)
	}

	err = createGroupNamespaceAccessRBAC(mockK8sClient, ctx, DoraNamespaceGroup, "dora-namespace")
	if err != nil {
		return fmt.Errorf("failed to set up group access to namespace: %w", err)
	}

	return nil
}

func (m *KubernetesClientMock) GetServiceDetails(sessionCtx context.Context, namespace string) ([]k8s.ServiceDetails, error) {
	originalServices, err := m.KubernetesClient.GetServiceDetails(sessionCtx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get service details: %w", err)
	}

	for i := range originalServices {
		originalServices[i].ClusterIP = "127.0.0.1"
		originalServices[i].HTTPPort = 8080
	}

	return originalServices, nil
}

func (m *KubernetesClientMock) GetServiceDetailsByName(sessionCtx context.Context, namespace string, serviceName string) (k8s.ServiceDetails, error) {
	originalService, err := m.KubernetesClient.GetServiceDetailsByName(sessionCtx, namespace, serviceName)
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

func createService(k8sClient client.Client, ctx context.Context, name string, namespace string, displayName string, description string, clusterIP string, componentLabel string) error {

	annotations := map[string]string{}

	if displayName != "" {
		annotations["displayName"] = displayName
	}

	if description != "" {
		annotations["description"] = description
	}

	labels := map[string]string{}
	if componentLabel != "" {
		labels["component"] = componentLabel
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
			Labels:      labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"component": k8s.ComponentLabelValue,
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

func createNamespace(k8sClient client.Client, ctx context.Context, namespace string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	err := k8sClient.Create(ctx, ns)
	if err != nil {
		return fmt.Errorf("failed to create namespace %s: %w", namespace, err)
	}

	return nil
}

func createClusterAdminRBAC(k8sClient client.Client, ctx context.Context, username string) error {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("cluster-admin-binding-%s", username),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "User",
				Name: username,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	err := k8sClient.Create(ctx, clusterRoleBinding)
	if err != nil {
		return fmt.Errorf("failed to create cluster admin ClusterRoleBinding: %w", err)
	}

	return nil
}

func createNamespaceRestrictedRBAC(k8sClient client.Client, ctx context.Context, username, namespace string) error {
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "namespace-restricted-role",
			Namespace: namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"services", "namespaces"},
				Verbs:     []string{"get", "list"},
			},
		},
	}

	err := k8sClient.Create(ctx, role)
	if err != nil {
		return fmt.Errorf("failed to create Role: %w", err)
	}

	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "namespace-restricted-binding",
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "User",
				Name: username,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     "namespace-restricted-role",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	err = k8sClient.Create(ctx, roleBinding)
	if err != nil {
		return fmt.Errorf("failed to create RoleBinding: %w", err)
	}

	return nil
}

func createGroupAccessRBAC(k8sClient client.Client, ctx context.Context, groupName, namespace, serviceName string) error {
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "group-model-registry-access",
			Namespace: namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"services"},
				Verbs:     []string{"get", "list"},
				ResourceNames: []string{
					serviceName,
				},
			},
		},
	}

	if err := k8sClient.Create(ctx, role); err != nil {
		return fmt.Errorf("failed to create Role for group: %w", err)
	}

	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "group-access-binding",
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "Group",
				Name: groupName,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     "group-model-registry-access",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	if err := k8sClient.Create(ctx, roleBinding); err != nil {
		return fmt.Errorf("failed to create RoleBinding for group: %w", err)
	}

	return nil
}

func createGroupNamespaceAccessRBAC(k8sClient client.Client, ctx context.Context, groupName, namespace string) error {

	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "group-namespace-access-role",
			Namespace: namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces", "services"},
				Verbs:     []string{"get", "list"},
			},
		},
	}

	if err := k8sClient.Create(ctx, role); err != nil {
		return fmt.Errorf("failed to create Role for group namespace access: %w", err)
	}

	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "group-namespace-access-binding",
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "Group",
				Name: groupName,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     "group-namespace-access-role",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	if err := k8sClient.Create(ctx, roleBinding); err != nil {
		return fmt.Errorf("failed to create RoleBinding for group namespace access: %w", err)
	}

	return nil
}

func strPtr(s string) *string {
	return &s
}
