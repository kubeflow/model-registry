package api

import (
	"io"
	"log/slog"

	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
)

// NewTestApp exposes a minimal constructor that allows tests and downstream
// extensions to configure specific App dependencies without invoking the
// production bootstrap logic.
func NewTestApp(cfg config.EnvConfig, logger *slog.Logger, factory k8s.KubernetesClientFactory, repos *repositories.Repositories) *App {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	return &App{
		config:                  cfg,
		logger:                  logger,
		kubernetesClientFactory: factory,
		repositories:            repos,
	}
}
