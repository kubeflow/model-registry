package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"log/slog"
	"net/http"
	"strings"
)

func NewKubernetesClientFactory(cfg config.EnvConfig, logger *slog.Logger) (KubernetesClientFactory, error) {
	switch cfg.AuthMethod {

	case config.AuthMethodInternal:
		k8sFactory, err := NewStaticClientFactory(logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create static client factory: %w", err)
		}
		return k8sFactory, nil

	case config.AuthMethodUser:
		k8sFactory := NewTokenClientFactory(logger, cfg)
		return k8sFactory, nil

	default:
		return nil, fmt.Errorf("invalid auth method: %q", cfg.AuthMethod)
	}
}

// ─── STATIC FACTORY (INTERNAL) ──────────────────────────────────────────
// uses the credentials of the running backend to create a single instance of the client
// If running inside the cluster, it uses the pod's service account.
// If running locally (e.g. for development), it uses the current user's kubeconfig context.
type KubernetesClientFactory interface {
	GetClient(ctx context.Context) (KubernetesClientInterface, error)
	ExtractRequestIdentity(httpHeader http.Header) (*RequestIdentity, error)
	ValidateRequestIdentity(identity *RequestIdentity) error
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

func (f *StaticClientFactory) ExtractRequestIdentity(httpHeader http.Header) (*RequestIdentity, error) {

	userID := httpHeader.Get(constants.KubeflowUserIDHeader)
	//`kubeflow-userid`: Contains the user's email address.
	if userID == "" {
		return nil, errors.New("missing required kubeflow-userid header")
	}

	userGroupsHeader := httpHeader.Get(constants.KubeflowUserGroupsIdHeader)
	// Note: The functionality for `kubeflow-groups` is not fully operational at Kubeflow platform at this time
	// but it's supported on Model Registry BFF
	//`kubeflow-groups`: Holds a comma-separated list of user groups.
	groups := []string{}
	if userGroupsHeader != "" {
		for _, g := range strings.Split(userGroupsHeader, ",") {
			groups = append(groups, strings.TrimSpace(g))
		}
	}
	identity := &RequestIdentity{
		UserID: userID,
		Groups: groups,
	}
	return identity, nil
}

func (f *StaticClientFactory) ValidateRequestIdentity(identity *RequestIdentity) error {
	if identity == nil {
		return errors.New("missing identity")
	}
	if identity.UserID == "" {
		return errors.New("user ID (kubeflow-userid) required for internal authentication")
	}
	return nil
}

//
// ─── TOKEN FACTORY (USER TOKEN) ────────────────────────────────────────────────
// uses a user-provided Bearer token for client creation.
// each user has a separate client instance.
//

type TokenClientFactory struct {
	Logger *slog.Logger
	Header string
	Prefix string
}

func NewTokenClientFactory(logger *slog.Logger, cfg config.EnvConfig) KubernetesClientFactory {
	return &TokenClientFactory{
		Logger: logger,
		Header: cfg.AuthTokenHeader,
		Prefix: cfg.AuthTokenPrefix,
	}
}

func (f *TokenClientFactory) ExtractRequestIdentity(httpHeader http.Header) (*RequestIdentity, error) {
	raw := httpHeader.Get(f.Header)
	if raw == "" {
		return nil, fmt.Errorf("missing required Header: %s", f.Header)
	}

	token := raw
	if f.Prefix != "" {
		if !strings.HasPrefix(raw, f.Prefix) {
			return nil, fmt.Errorf("expected token Header %s to start with Prefix %q", f.Header, f.Prefix)
		}
		token = strings.TrimPrefix(raw, f.Prefix)
	}

	return &RequestIdentity{
		Token: strings.TrimSpace(token),
	}, nil
}

func (f *TokenClientFactory) ValidateRequestIdentity(identity *RequestIdentity) error {

	if identity == nil {
		return errors.New("missing identity")
	}

	if identity.Token == "" {
		return errors.New("token is required for token-based authentication")
	}

	return nil
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

	return newTokenKubernetesClient(identity.Token, f.Logger)
}
