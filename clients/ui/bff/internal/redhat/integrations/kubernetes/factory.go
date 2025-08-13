package kubernetes

import (
	"log/slog"

	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
)

// NewRHOAIClientFactory creates a KubernetesClientFactory used in RHOAI
// when the authentication method is "user_token_red_hat".
// It injects a custom Kubernetes client that wraps and extends the default behavior.
func NewRHOAIClientFactory(logger *slog.Logger, cfg config.EnvConfig) kubernetes.KubernetesClientFactory {
	return kubernetes.NewTokenClientFactoryWithCustomNewTokenKubernetesClientFn(logger, cfg, NewRHOAIKubernetesClient)
}
