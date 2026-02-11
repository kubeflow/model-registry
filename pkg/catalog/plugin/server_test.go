package plugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testPlugin is a CatalogPlugin implementation for testing.
type testPlugin struct {
	name        string
	version     string
	description string
	healthy     bool
	initCalled  bool
	startCalled bool
	stopCalled  bool
}

func (p *testPlugin) Name() string        { return p.name }
func (p *testPlugin) Version() string     { return p.version }
func (p *testPlugin) Description() string { return p.description }

func (p *testPlugin) Init(ctx context.Context, cfg Config) error {
	p.initCalled = true
	return nil
}

func (p *testPlugin) Start(ctx context.Context) error {
	p.startCalled = true
	return nil
}

func (p *testPlugin) Stop(ctx context.Context) error {
	p.stopCalled = true
	return nil
}

func (p *testPlugin) Healthy() bool { return p.healthy }

func (p *testPlugin) RegisterRoutes(router chi.Router) error {
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	return nil
}

func (p *testPlugin) Migrations() []Migration { return nil }

func TestServerInit(t *testing.T) {
	Reset()

	plugin := &testPlugin{
		name:    "test",
		version: "v1",
		healthy: true,
	}
	Register(plugin)

	cfg := &CatalogSourcesConfig{
		Catalogs: map[string]CatalogSection{
			"test": {
				Sources: []SourceConfig{},
			},
		},
	}

	server := NewServer(cfg, []string{}, nil, nil)
	err := server.Init(context.Background())
	require.NoError(t, err)

	assert.True(t, plugin.initCalled)
	assert.Equal(t, 1, len(server.Plugins()))

	Reset()
}

func TestServerInitializesUnconfiguredPlugins(t *testing.T) {
	Reset()

	plugin := &testPlugin{
		name:    "test",
		version: "v1",
		healthy: true,
	}
	Register(plugin)

	// Empty config - no catalogs configured
	cfg := &CatalogSourcesConfig{
		Catalogs: map[string]CatalogSection{},
	}

	server := NewServer(cfg, []string{}, nil, nil)
	err := server.Init(context.Background())
	require.NoError(t, err)

	// Plugin should still be initialized even without config
	assert.True(t, plugin.initCalled)
	assert.Equal(t, 1, len(server.Plugins()))

	Reset()
}

func TestServerHealthEndpoint(t *testing.T) {
	Reset()

	cfg := &CatalogSourcesConfig{
		Catalogs: map[string]CatalogSection{},
	}

	server := NewServer(cfg, []string{}, nil, nil)
	err := server.Init(context.Background())
	require.NoError(t, err)

	router := server.MountRoutes()

	// Test /healthz
	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "ok")

	Reset()
}

func TestServerReadyEndpoint(t *testing.T) {
	Reset()

	healthyPlugin := &testPlugin{
		name:    "healthy",
		version: "v1",
		healthy: true,
	}
	Register(healthyPlugin)

	cfg := &CatalogSourcesConfig{
		Catalogs: map[string]CatalogSection{
			"healthy": {},
		},
	}

	server := NewServer(cfg, []string{}, nil, nil)
	err := server.Init(context.Background())
	require.NoError(t, err)

	router := server.MountRoutes()

	// Test /readyz when all plugins are healthy
	req := httptest.NewRequest("GET", "/readyz", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "ready")

	Reset()
}

func TestServerPluginsEndpoint(t *testing.T) {
	Reset()

	plugin := &testPlugin{
		name:        "test",
		version:     "v1",
		description: "Test plugin",
		healthy:     true,
	}
	Register(plugin)

	cfg := &CatalogSourcesConfig{
		Catalogs: map[string]CatalogSection{
			"test": {},
		},
	}

	server := NewServer(cfg, []string{}, nil, nil)
	err := server.Init(context.Background())
	require.NoError(t, err)

	router := server.MountRoutes()

	// Test /api/plugins
	req := httptest.NewRequest("GET", "/api/plugins", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "test")
	assert.Contains(t, body, "v1")

	Reset()
}

func TestServerStartStop(t *testing.T) {
	Reset()

	plugin := &testPlugin{
		name:    "test",
		version: "v1",
		healthy: true,
	}
	Register(plugin)

	cfg := &CatalogSourcesConfig{
		Catalogs: map[string]CatalogSection{
			"test": {},
		},
	}

	server := NewServer(cfg, []string{}, nil, nil)
	err := server.Init(context.Background())
	require.NoError(t, err)

	// Start
	err = server.Start(context.Background())
	require.NoError(t, err)
	assert.True(t, plugin.startCalled)

	// Stop
	err = server.Stop(context.Background())
	require.NoError(t, err)
	assert.True(t, plugin.stopCalled)

	Reset()
}
