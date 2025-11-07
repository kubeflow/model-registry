package catalog

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apimodels "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

func TestNewHFModelProvider_MissingAPIKey(t *testing.T) {
	source := &Source{
		CatalogSource: apimodels.CatalogSource{
			Id:   "test_hf",
			Name: "Test HF",
		},
		Type: "hf",
		Properties: map[string]any{
			"includedModels": []string{"test/model"},
		},
	}

	ctx := context.Background()
	_, err := newHFModelProvider(ctx, source, "")
	if err == nil {
		t.Fatal("Expected error for missing API key, got nil")
	}
	if !strings.Contains(err.Error(), "missing or invalid 'apiKey'") {
		t.Fatalf("Expected API key error, got: %s", err.Error())
	}
}

func TestNewHFModelProvider_MissingSourceID(t *testing.T) {
	source := &Source{
		CatalogSource: apimodels.CatalogSource{
			Id:   "", // Empty source ID
			Name: "Test HF",
		},
		Type: "hf",
		Properties: map[string]any{
			"apiKey":         "test-key",
			"includedModels": []string{"test/model"},
		},
	}

	ctx := context.Background()
	_, err := newHFModelProvider(ctx, source, "")
	if err == nil {
		t.Fatal("Expected error for missing source ID, got nil")
	}
	if !strings.Contains(err.Error(), "missing source ID") {
		t.Fatalf("Expected source ID error, got: %s", err.Error())
	}
}

func TestNewHFModelProvider_MissingIncludedModels(t *testing.T) {
	// Create mock server for credential validation
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/whoami-v2" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name": "test-user"}`))
		}
	}))
	defer server.Close()

	source := &Source{
		CatalogSource: apimodels.CatalogSource{
			Id:   "test_hf",
			Name: "Test HF",
		},
		Type: "hf",
		Properties: map[string]any{
			"apiKey": "test-key",
			"url":    server.URL,
		},
	}

	ctx := context.Background()
	_, err := newHFModelProvider(ctx, source, "")
	if err == nil {
		t.Fatal("Expected error for missing includedModels, got nil")
	}
	if !strings.Contains(err.Error(), "includedModels") {
		t.Fatalf("Expected includedModels error, got: %s", err.Error())
	}
}

func TestNewHFModelProvider_InvalidCredentials(t *testing.T) {
	// Create mock server that returns 401
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	source := &Source{
		CatalogSource: apimodels.CatalogSource{
			Id:   "test_hf",
			Name: "Test HF",
		},
		Type: "hf",
		Properties: map[string]any{
			"apiKey":         "invalid-key",
			"url":            server.URL,
			"includedModels": []string{"test/model"},
		},
	}

	ctx := context.Background()
	_, err := newHFModelProvider(ctx, source, "")
	if err == nil {
		t.Fatal("Expected error for invalid credentials, got nil")
	}
	if !strings.Contains(err.Error(), "invalid HuggingFace API credentials") {
		t.Fatalf("Expected credential validation error, got: %s", err.Error())
	}
}

func TestNewHFModelProvider_ExcludedModels(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/whoami-v2" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name": "test-user"}`))
			return
		}

		if strings.HasPrefix(r.URL.Path, "/api/models/") {
			modelName := strings.TrimPrefix(r.URL.Path, "/api/models/")
			modelInfo := map[string]interface{}{
				"id":     modelName,
				"author": "test-author",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(modelInfo)
		}
	}))
	defer server.Close()

	source := &Source{
		CatalogSource: apimodels.CatalogSource{
			Id:   "test_hf",
			Name: "Test HF",
		},
		Type: "hf",
		Properties: map[string]any{
			"apiKey":         "test-api-key",
			"url":            server.URL,
			"includedModels": []any{"test/model1", "test/model2", "test/model3"},
			"excludedModels": []any{"test/model2"},
		},
	}

	ctx := context.Background()
	recordsChan, err := newHFModelProvider(ctx, source, "")
	if err != nil {
		t.Fatalf("Failed to create HF model provider: %v", err)
	}

	// Read records from channel
	var records []ModelProviderRecord
	for record := range recordsChan {
		records = append(records, record)
	}

	// Should have 2 models (model1 and model3, model2 is excluded)
	if len(records) != 2 {
		t.Fatalf("Expected 2 model records (model2 excluded), got %d", len(records))
	}

	// Verify model2 is not in the results
	for _, record := range records {
		if record.Model.GetAttributes() != nil && record.Model.GetAttributes().Name != nil {
			if *record.Model.GetAttributes().Name == "test/model2" {
				t.Fatal("Expected model2 to be excluded, but it was included")
			}
		}
	}
}

func TestHFModel_PopulateFromHFInfo(t *testing.T) {
	hfm := &hfModel{}
	hfInfo := &hfModelInfo{
		ID:          "test/model",
		Author:      "test-author",
		LibraryName: "transformers",
		Task:        "text-generation",
		PipelineTag: "text-generation",
		Tags:        []string{"pytorch", "nlp"},
		Downloads:   1000,
		Sha:         "abc123",
		CreatedAt:   "2024-01-01T00:00:00Z",
		UpdatedAt:   "2024-01-02T00:00:00Z",
		CardData: &hfCard{
			Data: map[string]interface{}{
				"description":  "Test model description",
				"readme":       "# Test Model\nThis is a test.",
				"license":      "apache-2.0",
				"license_link": "https://www.apache.org/licenses/LICENSE-2.0",
			},
		},
		Config: &hfConfig{
			Architectures: []string{"BertForMaskedLM"},
			ModelType:     "bert",
		},
	}

	hfm.populateFromHFInfo(hfInfo, "test-source-id", "test/model")

	// Verify CatalogModel fields are populated
	if hfm.Name != "test/model" {
		t.Fatalf("Expected Name 'test/model', got '%s'", hfm.Name)
	}

	if hfm.ExternalId == nil || *hfm.ExternalId != "test/model" {
		t.Fatalf("Expected ExternalId 'test/model', got '%v'", hfm.ExternalId)
	}

	if hfm.SourceId == nil || *hfm.SourceId != "test-source-id" {
		t.Fatalf("Expected SourceId 'test-source-id', got '%v'", hfm.SourceId)
	}

	if hfm.Description == nil || *hfm.Description != "Test model description" {
		t.Fatalf("Expected Description 'Test model description', got '%v'", hfm.Description)
	}

	if hfm.Readme == nil || !strings.Contains(*hfm.Readme, "Test Model") {
		t.Fatalf("Expected Readme to contain 'Test Model', got '%v'", hfm.Readme)
	}

	if hfm.License == nil || *hfm.License != "apache-2.0" {
		t.Fatalf("Expected License 'apache-2.0', got '%v'", hfm.License)
	}

	if hfm.Provider == nil || *hfm.Provider != "test-author" {
		t.Fatalf("Expected Provider 'test-author', got '%v'", hfm.Provider)
	}

	if hfm.LibraryName == nil || *hfm.LibraryName != "transformers" {
		t.Fatalf("Expected LibraryName 'transformers', got '%v'", hfm.LibraryName)
	}

	if len(hfm.Tasks) != 1 || hfm.Tasks[0] != "text-generation" {
		t.Fatalf("Expected Tasks ['text-generation'], got '%v'", hfm.Tasks)
	}

	if hfm.CustomProperties == nil {
		t.Fatal("Expected CustomProperties to be set, got nil")
	}

	customProps := *hfm.CustomProperties
	if _, ok := customProps["hf_tags"]; !ok {
		t.Fatal("Expected hf_tags in custom properties")
	}
	if _, ok := customProps["hf_downloads"]; !ok {
		t.Fatal("Expected hf_downloads in custom properties")
	}
	if _, ok := customProps["hf_sha"]; !ok {
		t.Fatal("Expected hf_sha in custom properties")
	}
}

func TestNewHFModelProvider_EmptyIncludedModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/whoami-v2" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name": "test-user"}`))
		}
	}))
	defer server.Close()

	source := &Source{
		CatalogSource: apimodels.CatalogSource{
			Id:   "test_hf",
			Name: "Test HF",
		},
		Type: "hf",
		Properties: map[string]any{
			"apiKey":         "test-api-key",
			"url":            server.URL,
			"includedModels": []string{}, // Empty list
		},
	}

	ctx := context.Background()
	_, err := newHFModelProvider(ctx, source, "")
	if err == nil {
		t.Fatal("Expected error for empty includedModels, got nil")
	}
	if !strings.Contains(err.Error(), "property should be a list") {
		t.Fatalf("Expected empty list error, got: %s", err.Error())
	}
}
