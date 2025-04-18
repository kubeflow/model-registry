package k8mocks

import (
	"context"
	"fmt"
	kubernetes2 "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var DefaultTestUsers = []TestUser{
	{
		UserName: "user@example.com",
		Token:    "FAKE_CLUSTER_ADMIN_TOKEN",
		Groups:   []string{},
	},
	{
		UserName: "doraNonAdmin@example.com",
		Token:    "FAKE_DORA_TOKEN",
		Groups:   []string{"dora-namespace-group", "dora-service-group"},
	},
	{
		UserName: "bellaNonAdmin@example.com",
		Token:    "FAKE_BELLA_TOKEN",
		Groups:   []string{},
	},
}

type TestUser struct {
	UserName string
	Token    string
	Groups   []string
}

type TestEnvInput struct {
	Users  []TestUser
	Logger *slog.Logger
	Ctx    context.Context
	Cancel context.CancelFunc
}

func SetupEnvTest(input TestEnvInput) (*envtest.Environment, kubernetes.Interface, error) {
	projectRoot, err := getProjectRoot()
	if err != nil {
		input.Logger.Error("failed to find project root", slog.String("error", err.Error()))
		input.Cancel()
		os.Exit(1)
	}

	testEnv := &envtest.Environment{
		BinaryAssetsDirectory: filepath.Join(projectRoot, "bin", "k8s", fmt.Sprintf("1.29.0-%s-%s", runtime.GOOS, runtime.GOARCH)),
	}

	cfg, err := testEnv.Start()
	if err != nil {
		input.Logger.Error("failed to start envtest", slog.String("error", err.Error()))
		input.Cancel()
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		input.Logger.Error("failed to create clientset", slog.String("error", err.Error()))
		input.Cancel()
		os.Exit(1)
	}

	// bootstrap resources
	err = setupMock(clientset, input.Ctx)
	if err != nil {
		input.Logger.Error("failed to setup mock data", slog.String("error", err.Error()))
		input.Cancel()
		os.Exit(1)
	}

	return testEnv, clientset, nil
}

func setupMock(mockK8sClient kubernetes.Interface, ctx context.Context) error {

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

	err = createNamespace(mockK8sClient, ctx, "bento-namespace")
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

	err = createClusterAdminRBAC(mockK8sClient, ctx, DefaultTestUsers[0].UserName)
	if err != nil {
		return fmt.Errorf("failed to create cluster admin RBAC: %w", err)
	}

	err = createNamespaceRestrictedRBAC(mockK8sClient, ctx, DefaultTestUsers[1].UserName, "dora-namespace")
	if err != nil {
		return fmt.Errorf("failed to create namespace-restricted RBAC: %w", err)
	}

	err = createNamespaceRestrictedRBAC(mockK8sClient, ctx, DefaultTestUsers[2].UserName, "bella-namespace")
	if err != nil {
		return fmt.Errorf("failed to create namespace-restricted RBAC: %w", err)
	}

	err = createGroupAccessRBAC(mockK8sClient, ctx, DefaultTestUsers[1].Groups[1], "dora-namespace", "model-registry-dora")
	if err != nil {
		return fmt.Errorf("failed to create group-based RBAC: %w", err)
	}

	err = createGroupNamespaceAccessRBAC(mockK8sClient, ctx, DefaultTestUsers[1].Groups[0], "dora-namespace")
	if err != nil {
		return fmt.Errorf("failed to set up group access to namespace: %w", err)
	}

	return nil
}

func createNamespace(k8sClient kubernetes.Interface, ctx context.Context, namespace string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	_, err := k8sClient.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create namespace %s: %w", namespace, err)
	}

	return nil
}

func createClusterAdminRBAC(k8sClient kubernetes.Interface, ctx context.Context, username string) error {
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

	_, err := k8sClient.RbacV1().ClusterRoleBindings().Create(ctx, clusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create cluster admin ClusterRoleBinding: %w", err)
	}

	return nil
}

func createNamespaceRestrictedRBAC(k8sClient kubernetes.Interface, ctx context.Context, username, namespace string) error {
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

	_, err := k8sClient.RbacV1().Roles(namespace).Create(ctx, role, metav1.CreateOptions{})
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

	_, err = k8sClient.RbacV1().RoleBindings(namespace).Create(ctx, roleBinding, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create RoleBinding: %w", err)
	}

	return nil
}

func createGroupAccessRBAC(k8sClient kubernetes.Interface, ctx context.Context, groupName, namespace, serviceName string) error {
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

	_, err := k8sClient.RbacV1().Roles(namespace).Create(ctx, role, metav1.CreateOptions{})
	if err != nil {
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

	_, err = k8sClient.RbacV1().RoleBindings(namespace).Create(ctx, roleBinding, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create RoleBinding for group: %w", err)
	}
	return nil
}

func createGroupNamespaceAccessRBAC(k8sClient kubernetes.Interface, ctx context.Context, groupName, namespace string) error {

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

	_, err := k8sClient.RbacV1().Roles(namespace).Create(ctx, role, metav1.CreateOptions{})
	if err != nil {
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

	_, err = k8sClient.RbacV1().RoleBindings(namespace).Create(ctx, roleBinding, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create RoleBinding for group namespace access: %w", err)
	}

	return nil
}

func createService(k8sClient kubernetes.Interface, ctx context.Context, name string, namespace string, displayName string, description string, clusterIP string, componentLabel string) error {

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
				"component": kubernetes2.ComponentLabelValue,
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

	// Create the service using kubernetes.Interface
	_, err := k8sClient.CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	return nil
}

func strPtr(s string) *string {
	return &s
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
