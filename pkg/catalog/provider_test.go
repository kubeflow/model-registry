package catalog

import (
	"context"
	"testing"
)

func TestProviderRegistry(t *testing.T) {
	type Entity struct {
		Name string
	}
	type Artifact struct {
		URI string
	}

	registry := NewProviderRegistry[*Entity, *Artifact]()

	// Test Register
	providerFunc := func(ctx context.Context, source *Source, reldir string) (<-chan Record[*Entity, *Artifact], error) {
		ch := make(chan Record[*Entity, *Artifact])
		close(ch)
		return ch, nil
	}

	err := registry.Register("test", providerFunc)
	if err != nil {
		t.Errorf("Register failed: %v", err)
	}

	// Test duplicate registration
	err = registry.Register("test", providerFunc)
	if err == nil {
		t.Error("Expected error for duplicate registration, got nil")
	}

	// Test Get
	fn, ok := registry.Get("test")
	if !ok {
		t.Error("Expected to find registered provider")
	}
	if fn == nil {
		t.Error("Expected non-nil provider function")
	}

	// Test Get non-existent
	_, ok = registry.Get("nonexistent")
	if ok {
		t.Error("Expected not to find non-existent provider")
	}

	// Test Has
	if !registry.Has("test") {
		t.Error("Expected Has to return true for registered provider")
	}
	if registry.Has("nonexistent") {
		t.Error("Expected Has to return false for non-existent provider")
	}

	// Test Names
	names := registry.Names()
	if len(names) != 1 || names[0] != "test" {
		t.Errorf("Expected names to be [test], got %v", names)
	}
}

func TestMustRegister(t *testing.T) {
	type Entity struct{}
	type Artifact struct{}

	registry := NewProviderRegistry[*Entity, *Artifact]()

	// Should not panic
	registry.MustRegister("test", func(ctx context.Context, source *Source, reldir string) (<-chan Record[*Entity, *Artifact], error) {
		return nil, nil
	})

	// Should panic on duplicate
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic on duplicate MustRegister")
		}
	}()

	registry.MustRegister("test", func(ctx context.Context, source *Source, reldir string) (<-chan Record[*Entity, *Artifact], error) {
		return nil, nil
	})
}
