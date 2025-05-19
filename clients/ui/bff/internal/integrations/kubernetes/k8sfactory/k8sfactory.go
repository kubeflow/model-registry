package k8sfactory

import (
	"fmt"
	"log/slog"

	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	redhat "github.com/kubeflow/model-registry/ui/bff/internal/redhat/integrations/kubernetes"
)

func NewKubernetesClientFactory(cfg config.EnvConfig, logger *slog.Logger) (kubernetes.KubernetesClientFactory, error) {
	switch cfg.AuthMethod {

	case config.AuthMethodInternal:
		k8sFactory, err := kubernetes.NewStaticClientFactory(logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create static client factory: %w", err)
		}
		return k8sFactory, nil

	case config.AuthMethodUser:
		k8sFactory := kubernetes.NewTokenClientFactory(logger, cfg)
		return k8sFactory, nil

	//TODO LUCAS red hat only code
	case config.AuthMethodRedHatUser:
		k8sFactory := redhat.NewRHOAIClientFactory(logger, cfg)
		return k8sFactory, nil

	default:
		return nil, fmt.Errorf("invalid auth method: %q", cfg.AuthMethod)
	}
}
