package catalog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

func TestNewHfCatalog_MissingAPIKey(t *testing.T) {
	source := &CatalogSourceConfig{
		CatalogSource: openapi.CatalogSource{
			Id:   "test_hf",
			Name: "Test HF",
		},
		Type: "hf",
		Properties: map[string]any{
			"url": "https://huggingface.co",
		},
	}

	_, err := newHfCatalog(source)
	if err == nil {
		t.Fatal("Expected error for missing API key, got nil")
	}
	if err.Error() != "missing or invalid 'apiKey' property for HuggingFace catalog" {
		t.Fatalf("Expected specific error message, got: %s", err.Error())
	}
}

func TestNewHfCatalog_WithValidCredentials(t *testing.T) {
	// Create mock server that returns valid response for credential validation
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-api-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch r.URL.Path {
		case "/api/whoami-v2":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name": "test-user", "type": "user"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	source := &CatalogSourceConfig{
		CatalogSource: openapi.CatalogSource{
			Id:   "test_hf",
			Name: "Test HF",
		},
		Type: "hf",
		Properties: map[string]any{
			"apiKey":     "test-api-key",
			"url":        server.URL,
			"modelLimit": 10,
		},
	}

	catalog, err := newHfCatalog(source)
	if err != nil {
		t.Fatalf("Failed to create HF catalog: %v", err)
	}

	hfCatalog := catalog.(*hfCatalogImpl)

	// Test that methods return appropriate responses for stub implementation
	ctx := context.Background()

	// Test GetModel - should return not implemented error
	model, err := hfCatalog.GetModel(ctx, "test-model")
	if err == nil {
		t.Fatal("Expected not implemented error, got nil")
	}
	if model != nil {
		t.Fatal("Expected nil model, got non-nil")
	}

	// Test ListModels - should return empty list
	listParams := ListModelsParams{
		Query:     "",
		OrderBy:   openapi.ORDERBYFIELD_NAME,
		SortOrder: openapi.SORTORDER_ASC,
	}
	modelList, err := hfCatalog.ListModels(ctx, listParams)
	if err != nil {
		t.Fatalf("Failed to list models: %v", err)
	}
	if len(modelList.Items) != 0 {
		t.Fatalf("Expected 0 models, got %d", len(modelList.Items))
	}

	// Test GetArtifacts - should return empty list
	artifacts, err := hfCatalog.GetArtifacts(ctx, "test-model")
	if err != nil {
		t.Fatalf("Failed to get artifacts: %v", err)
	}
	if artifacts == nil {
		t.Fatal("Expected artifacts list, got nil")
	}
	if len(artifacts.Items) != 0 {
		t.Fatalf("Expected 0 artifacts, got %d", len(artifacts.Items))
	}
}

func TestNewHfCatalog_InvalidCredentials(t *testing.T) {
	// Create mock server that returns 401
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	source := &CatalogSourceConfig{
		CatalogSource: openapi.CatalogSource{
			Id:   "test_hf",
			Name: "Test HF",
		},
		Type: "hf",
		Properties: map[string]any{
			"apiKey": "invalid-key",
			"url":    server.URL,
		},
	}

	_, err := newHfCatalog(source)
	if err == nil {
		t.Fatal("Expected error for invalid credentials, got nil")
	}
	if !strings.Contains(err.Error(), "invalid HuggingFace API credentials") {
		t.Fatalf("Expected credential validation error, got: %s", err.Error())
	}
}

func TestNewHfCatalog_DefaultConfiguration(t *testing.T) {
	// Create mock server for default HuggingFace URL
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "test-user"}`))
	}))
	defer server.Close()

	source := &CatalogSourceConfig{
		CatalogSource: openapi.CatalogSource{
			Id:   "test_hf",
			Name: "Test HF",
		},
		Type: "hf",
		Properties: map[string]any{
			"apiKey": "test-key",
			"url":    server.URL, // Override default for testing
		},
	}

	catalog, err := newHfCatalog(source)
	if err != nil {
		t.Fatalf("Failed to create HF catalog with defaults: %v", err)
	}

	hfCatalog := catalog.(*hfCatalogImpl)
	if hfCatalog.apiKey != "test-key" {
		t.Fatalf("Expected apiKey 'test-key', got '%s'", hfCatalog.apiKey)
	}
	if hfCatalog.baseURL != server.URL {
		t.Fatalf("Expected baseURL '%s', got '%s'", server.URL, hfCatalog.baseURL)
	}
}
