package k8mocks

import (
	"context"
	"fmt"
	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log/slog"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sync"
)

type MockedKubernetesClientFactory interface {
	k8s.KubernetesClientFactory
}

func NewMockedKubernetesClientFactory(clientset kubernetes.Interface, testEnv *envtest.Environment, cfg config.EnvConfig, logger *slog.Logger) (k8s.KubernetesClientFactory, error) {
	switch cfg.AuthMethod {
	case config.AuthMethodInternal:
		k8sFactory, err := NewStaticClientFactory(clientset, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create static client factory: %w", err)
		}
		return k8sFactory, nil

	case config.AuthMethodUser:
		k8sFactory, err := NewTokenClientFactory(clientset, testEnv.Config, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create static client factory: %w", err)
		}
		return k8sFactory, nil

	default:
		return nil, fmt.Errorf("invalid auth method: %q", cfg.AuthMethod)
	}
}

// ─── MOCKED STATIC FACTORY (envtest + "INTERNAL ACCOUNT") ──────────────────────────────────────────
type MockedStaticClientFactory struct {
	logger                       *slog.Logger
	serviceAccountMockedK8client k8s.KubernetesClientInterface
	clientset                    kubernetes.Interface
	initErr                      error
	initLock                     sync.Mutex
	realFactoryWithoutClient     k8s.StaticClientFactory
}

func NewStaticClientFactory(clientset kubernetes.Interface, logger *slog.Logger) (k8s.KubernetesClientFactory, error) {
	realFactoryWithoutClient := k8s.StaticClientFactory{
		Logger: logger,
	}
	return &MockedStaticClientFactory{
		logger:                   logger,
		clientset:                clientset,
		realFactoryWithoutClient: realFactoryWithoutClient,
	}, nil
}

func (f *MockedStaticClientFactory) GetClient(_ context.Context) (k8s.KubernetesClientInterface, error) {
	f.initLock.Lock()
	defer f.initLock.Unlock()

	if f.serviceAccountMockedK8client != nil {
		return f.serviceAccountMockedK8client, nil
	}

	f.logger.Info("Initializing mocked service account client")
	client := newMockedInternalKubernetesClientFromClientset(f.clientset, f.logger)
	if client == nil {
		f.initErr = fmt.Errorf("failed to create mocked service account client")
		return nil, f.initErr
	}

	f.serviceAccountMockedK8client = client
	return f.serviceAccountMockedK8client, nil
}

func (f *MockedStaticClientFactory) ExtractRequestIdentity(httpHeader http.Header) (*k8s.RequestIdentity, error) {
	return f.realFactoryWithoutClient.ExtractRequestIdentity(httpHeader)
}
func (f *MockedStaticClientFactory) ValidateRequestIdentity(identity *k8s.RequestIdentity) error {
	return f.realFactoryWithoutClient.ValidateRequestIdentity(identity)
}

// ─── MOCKED TOKEN FACTORY (envtest + "USER TOKEN") ──────────────────────────────
//
// MockedTokenClientFactory simulates token-based client creation in tests.
// It maps fake tokens (like "FAKE_BELLA_TOKEN") to a TestUser (username + groups),
// and creates a Kubernetes client that impersonates that user.
// This is critical for triggering proper RBAC evaluation (e.g., SelfSubjectAccessReview) inside envtest,
// which does not perform real token authentication.
type MockedTokenClientFactory struct {
	logger     *slog.Logger
	clientset  kubernetes.Interface
	restConfig *rest.Config

	clients        map[string]k8s.KubernetesClientInterface
	initLock       sync.Mutex
	realK8sFactory k8s.KubernetesClientFactory
}

// NewTokenClientFactory initializes a factory using a known envtest clientset + config.
func NewTokenClientFactory(clientset kubernetes.Interface, restConfig *rest.Config, logger *slog.Logger) (k8s.KubernetesClientFactory, error) {
	cfg := config.EnvConfig{
		AuthMethod:      config.AuthMethodUser,
		AuthTokenHeader: config.DefaultAuthTokenHeader,
		AuthTokenPrefix: config.DefaultAuthTokenPrefix,
	}
	realFactory := k8s.NewTokenClientFactory(logger, cfg)

	return &MockedTokenClientFactory{
		logger:         logger,
		clientset:      clientset,
		restConfig:     restConfig,
		realK8sFactory: realFactory,
		clients:        make(map[string]k8s.KubernetesClientInterface),
	}, nil
}

func (f *MockedTokenClientFactory) ExtractRequestIdentity(httpHeader http.Header) (*k8s.RequestIdentity, error) {
	return f.realK8sFactory.ExtractRequestIdentity(httpHeader)
}

func (f *MockedTokenClientFactory) ValidateRequestIdentity(identity *k8s.RequestIdentity) error {
	return f.realK8sFactory.ValidateRequestIdentity(identity)
}

// GetClient returns a Kubernetes client for the identity in context,
// impersonating the associated user to allow SelfSubjectAccessReview (SSAR) and RBAC testing.
func (f *MockedTokenClientFactory) GetClient(ctx context.Context) (k8s.KubernetesClientInterface, error) {
	val := ctx.Value(constants.RequestIdentityKey)
	if val == nil {
		return nil, fmt.Errorf("missing RequestIdentity in context")
	}

	identity, ok := val.(*k8s.RequestIdentity)
	if !ok || identity.Token == "" {
		return nil, fmt.Errorf("invalid or missing identity token")
	}

	f.initLock.Lock()
	defer f.initLock.Unlock()

	if client, exists := f.clients[identity.Token]; exists {
		return client, nil
	}

	// Map token to test user identity
	user := findTestUserByToken(identity.Token)
	if user == nil {
		return nil, fmt.Errorf("unknown test token: %s", identity.Token)
	}

	// Create a new rest.Config that impersonates the user.
	// This bypasses the lack of real authentication in envtest and allows RBAC to work properly.
	impersonatedCfg := rest.CopyConfig(f.restConfig)
	impersonatedCfg.Impersonate = rest.ImpersonationConfig{
		UserName: user.UserName,
		Groups:   user.Groups,
	}

	clientset, err := kubernetes.NewForConfig(impersonatedCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create impersonated client: %w", err)
	}

	client := newMockedTokenKubernetesClientFromClientset(clientset, f.logger)
	f.clients[identity.Token] = client
	return client, nil
}

func findTestUserByToken(token string) *TestUser {
	for _, u := range DefaultTestUsers {
		if u.Token == token {
			return &u
		}
	}
	return nil
}
