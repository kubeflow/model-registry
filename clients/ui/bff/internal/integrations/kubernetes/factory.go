package kubernetes

import (
	"context"
	"fmt"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"log/slog"
)

// ─── STATIC FACTORY (INTERNAL) ──────────────────────────────────────────
// uses the credentials of the running backend to create a single instance of the client
// If running inside the cluster, it uses the pod's service account.
// If running locally (e.g. for development), it uses the current user's kubeconfig context.
type KubernetesClientFactory interface {
	GetClient(ctx context.Context) (KubernetesClientInterface, error)
}

type StaticClientFactory struct {
	Logger *slog.Logger
	Client KubernetesClientInterface
}

func NewStaticClientFactory(logger *slog.Logger) (KubernetesClientFactory, error) {
	client, err := newInternalKubernetesClient(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create service account client: %w", err)
	}
	return &StaticClientFactory{
		Client: client,
		Logger: logger,
	}, nil
}

func (f *StaticClientFactory) GetClient(_ context.Context) (KubernetesClientInterface, error) {
	return f.Client, nil
}

//
// ─── TOKEN FACTORY (USER TOKEN) ────────────────────────────────────────────────
// uses a user-provided Bearer token for client creation.
// each user has a separate client instance.
//

type TokenClientFactory struct {
	logger *slog.Logger
}

func NewTokenClientFactory(logger *slog.Logger) KubernetesClientFactory {
	return &TokenClientFactory{logger: logger}
}

func (f *TokenClientFactory) GetClient(ctx context.Context) (KubernetesClientInterface, error) {
	identityVal := ctx.Value(constants.RequestIdentityKey)
	if identityVal == nil {
		return nil, fmt.Errorf("missing RequestIdentity in context")
	}

	identity, ok := identityVal.(*RequestIdentity)
	if !ok || identity.Token == "" {
		return nil, fmt.Errorf("invalid or missing identity token")
	}

	return newTokenKubernetesClient(identity.Token, f.logger)
}
