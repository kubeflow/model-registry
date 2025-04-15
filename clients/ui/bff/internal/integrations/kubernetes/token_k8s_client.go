package kubernetes

import (
	"context"
	"fmt"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	authnv1 "k8s.io/api/authentication/v1"
	authv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log/slog"
	"strings"
	"time"
)

type TokenKubernetesClient struct {
	SharedClientLogic
}

func (kc *TokenKubernetesClient) IsClusterAdmin(_ *RequestIdentity) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// We cannot list ClusterRoleBindings here because this client is initialized with a user token,
	// which typically does not have permissions to read cluster-scoped RBAC resources.
	// Instead, we use a SelfSubjectAccessReview with wildcard '*' verb and resource,
	// which safely asks the Kubernetes API server: "Can I do everything?"
	// If the review returns allowed=true, it means the user has cluster-admin-equivalent permissions.
	sar := &authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authv1.ResourceAttributes{
				Verb:     "*",
				Resource: "*",
			},
		},
	}

	resp, err := kc.Client.AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
	if err != nil {
		kc.Logger.Error("failed to perform cluster-admin SAR", "error", err)
		return false, fmt.Errorf("failed to verify cluster-admin permissions: %w", err)
	}

	if !resp.Status.Allowed {
		kc.Logger.Info("user is NOT cluster-admin")
		return false, nil
	}

	kc.Logger.Info("user is cluster-admin")
	return true, nil
}

// newTokenKubernetesClient creates a Kubernetes client using a user bearer token.
func newTokenKubernetesClient(token string, logger *slog.Logger) (KubernetesClientInterface, error) {
	baseConfig, err := helper.GetKubeconfig()
	if err != nil {
		logger.Error("failed to get kubeconfig", "error", err)
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	// Start with an anonymous config to avoid preloaded auth
	cfg := rest.AnonymousClientConfig(baseConfig)
	if err != nil {
		logger.Error("failed to create anonymous config", "error", err)
		return nil, fmt.Errorf("failed to create anonymous config: %w", err)
	}
	cfg.BearerToken = token

	// Explicitly clear all other auth mechanisms
	cfg.BearerTokenFile = ""
	cfg.Username = ""
	cfg.Password = ""
	cfg.ExecProvider = nil
	cfg.AuthProvider = nil

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logger.Error("failed to create token-based Kubernetes client", "error", err)
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return &TokenKubernetesClient{
		SharedClientLogic: SharedClientLogic{
			Client: clientset,
			Logger: logger,
			// Token is retained for follow-up calls; do not log it.
			Token: NewBearerToken(token),
		},
	}, nil
}

// RequestIdentity is unused because the token already represents the user identity.
func (kc *TokenKubernetesClient) CanListServicesInNamespace(ctx context.Context, _ *RequestIdentity, namespace string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	for _, verb := range []string{"get", "list"} {
		sar := &authv1.SelfSubjectAccessReview{
			Spec: authv1.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &authv1.ResourceAttributes{
					Verb:      verb,
					Resource:  "services",
					Namespace: namespace,
				},
			},
		}

		resp, err := kc.Client.AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
		if err != nil {
			kc.Logger.Error("self-SAR failed", "namespace", namespace, "verb", verb, "error", err)
			return false, err
		}

		if !resp.Status.Allowed {
			kc.Logger.Error("self-SAR denied", "namespace", namespace, "verb", verb)
			return false, nil
		}
	}

	return true, nil
}

// RequestIdentity is unused because the token already represents the user identity.
func (kc *TokenKubernetesClient) CanAccessServiceInNamespace(ctx context.Context, _ *RequestIdentity, namespace, serviceName string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	sar := &authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authv1.ResourceAttributes{
				Verb:      "get",
				Resource:  "services",
				Namespace: namespace,
				Name:      serviceName,
			},
		},
	}

	resp, err := kc.Client.AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
	if err != nil {
		kc.Logger.Error("self-SAR failed", "service", serviceName, "namespace", namespace, "error", err)
		return false, err
	}
	if !resp.Status.Allowed {
		kc.Logger.Error("self-SAR denied", "service", serviceName, "namespace", namespace)
		return false, nil
	}

	return true, nil
}

// RequestIdentity is unused because the token already represents the user identity.
// This endpoint is used only on dev mode that is why is safe to ignore permissions errors
func (kc *TokenKubernetesClient) GetNamespaces(ctx context.Context, _ *RequestIdentity) ([]corev1.Namespace, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	nsList, err := kc.Client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		kc.Logger.Error("user is not allowed to list namespaces or failed to list namespaces")
		return []corev1.Namespace{}, fmt.Errorf("failed to list namespaces: %w", err)
	}

	return nsList.Items, nil
}

func (kc *TokenKubernetesClient) GetUser(_ *RequestIdentity) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ssr := &authnv1.SelfSubjectReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "SelfSubjectReview",
			APIVersion: "authentication.k8s.io/v1",
		},
	}

	resp, err := kc.Client.AuthenticationV1().SelfSubjectReviews().Create(ctx, ssr, metav1.CreateOptions{})
	if err != nil {
		kc.Logger.Error("failed to get user identity from token", "error", err)
		return "", fmt.Errorf("failed to get user identity: %w", err)
	}

	username := resp.Status.UserInfo.Username
	if username == "" {
		kc.Logger.Error("user identity not found in token")
		return "", fmt.Errorf("no username found in token")
	}

	// If it's a service account, extract the SA name
	const saPrefix = "system:serviceaccount:"
	if strings.HasPrefix(username, saPrefix) {
		parts := strings.SplitN(strings.TrimPrefix(username, saPrefix), ":", 2)
		if len(parts) == 2 {
			return parts[1], nil // Return just the service account name
		}
		kc.Logger.Warn("malformed service account username", "username", username)
	}

	return username, nil
}
