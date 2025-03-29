package kubernetes

import (
	"context"
	"fmt"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	authv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log/slog"
	"time"
)

type TokenKubernetesClient struct {
	SharedClientLogic
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
	// Reuse CA settings from base config to validate the API server's TLS certificate.
	// This ensures secure communication and prevents x509 trust errors with self-signed or cluster-issued certs.
	cfg.TLSClientConfig = rest.TLSClientConfig{
		CAFile: baseConfig.TLSClientConfig.CAFile,
		CAData: baseConfig.TLSClientConfig.CAData,
	}
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
			kc.Logger.Warn("self-SAR denied", "namespace", namespace, "verb", verb)
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
		kc.Logger.Warn("self-SAR denied", "service", serviceName, "namespace", namespace)
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
		kc.Logger.Warn("user is not allowed to list namespaces or failed to list namespaces")
		return []corev1.Namespace{}, nil
	}

	return nsList.Items, nil
}
