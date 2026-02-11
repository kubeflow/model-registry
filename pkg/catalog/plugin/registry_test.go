package plugin

import (
	"context"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// mockPlugin is a minimal CatalogPlugin implementation for testing.
type mockPlugin struct {
	name        string
	version     string
	description string
	healthy     bool
}

func (p *mockPlugin) Name() string                               { return p.name }
func (p *mockPlugin) Version() string                            { return p.version }
func (p *mockPlugin) Description() string                        { return p.description }
func (p *mockPlugin) Init(ctx context.Context, cfg Config) error { return nil }
func (p *mockPlugin) Start(ctx context.Context) error            { return nil }
func (p *mockPlugin) Stop(ctx context.Context) error             { return nil }
func (p *mockPlugin) Healthy() bool                              { return p.healthy }
func (p *mockPlugin) RegisterRoutes(router chi.Router) error     { return nil }
func (p *mockPlugin) Migrations() []Migration { return nil }

func TestRegister(t *testing.T) {
	Reset() // Clear any existing plugins

	plugin := &mockPlugin{
		name:        "test-plugin",
		version:     "v1",
		description: "Test plugin",
		healthy:     true,
	}

	Register(plugin)

	// Verify plugin was registered
	names := Names()
	assert.Equal(t, 1, len(names))
	assert.Equal(t, "test-plugin", names[0])

	// Verify Get works
	p, ok := Get("test-plugin")
	assert.True(t, ok)
	assert.Equal(t, "test-plugin", p.Name())
	assert.Equal(t, "v1", p.Version())

	// Verify non-existent plugin
	_, ok = Get("non-existent")
	assert.False(t, ok)

	Reset()
}

func TestRegisterDuplicate(t *testing.T) {
	Reset()

	plugin1 := &mockPlugin{name: "duplicate"}
	plugin2 := &mockPlugin{name: "duplicate"}

	Register(plugin1)

	// Should panic on duplicate registration
	assert.Panics(t, func() {
		Register(plugin2)
	})

	Reset()
}

func TestAll(t *testing.T) {
	Reset()

	plugin1 := &mockPlugin{name: "plugin-a"}
	plugin2 := &mockPlugin{name: "plugin-b"}
	plugin3 := &mockPlugin{name: "plugin-c"}

	Register(plugin1)
	Register(plugin2)
	Register(plugin3)

	all := All()
	assert.Equal(t, 3, len(all))

	// Verify registration order is preserved
	assert.Equal(t, "plugin-a", all[0].Name())
	assert.Equal(t, "plugin-b", all[1].Name())
	assert.Equal(t, "plugin-c", all[2].Name())

	Reset()
}

func TestCount(t *testing.T) {
	Reset()

	assert.Equal(t, 0, Count())

	Register(&mockPlugin{name: "p1"})
	assert.Equal(t, 1, Count())

	Register(&mockPlugin{name: "p2"})
	assert.Equal(t, 2, Count())

	Reset()
	assert.Equal(t, 0, Count())
}
