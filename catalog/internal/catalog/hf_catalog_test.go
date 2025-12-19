package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apimodels "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/db/models"
)

func TestPopulateFromHFInfo(t *testing.T) {
	tests := []struct {
		name               string
		hfInfo             *hfModelInfo
		sourceId           string
		originalModelName  string
		expectedName       string
		expectedExternalID *string
		expectedSourceID   *string
		expectedProvider   *string
		expectedLicense    *string
		expectedLibrary    *string
		hasReadme          bool
		hasDescription     bool
		hasTasks           bool
		hasCustomProps     bool
	}{
		{
			name: "complete model info",
			hfInfo: &hfModelInfo{
				ID:          "test-org/test-model",
				Author:      "test-author",
				Sha:         "abc123",
				CreatedAt:   "2023-01-01T00:00:00Z",
				UpdatedAt:   "2023-01-02T00:00:00Z",
				Downloads:   1000,
				Tags:        []string{"license:mit", "transformers", "pytorch"},
				PipelineTag: "text-generation",
				Task:        "text-generation",
				LibraryName: "transformers",
				Config: &hfConfig{
					Architectures: []string{"GPT2LMHeadModel"},
					ModelType:     "gpt2",
				},
				CardData: &hfCard{
					Data: map[string]interface{}{
						"description": "A test model description",
					},
				},
			},
			sourceId:          "test-source-id",
			originalModelName: "test-org/test-model",
			expectedName:      "test-org/test-model",
			expectedProvider:  stringPtr("test-author"),
			expectedLicense:   stringPtr("MIT"),
			expectedLibrary:   stringPtr("transformers"),
			hasTasks:          true,
			hasCustomProps:    true,
			hasReadme:         false, // No README fetching in unit tests
		},
		{
			name: "model with ModelID fallback",
			hfInfo: &hfModelInfo{
				ModelID: "fallback-model-id",
				Author:  "another-author",
			},
			sourceId:          "source-2",
			originalModelName: "original-name",
			expectedName:      "fallback-model-id",
			expectedProvider:  stringPtr("another-author"),
		},
		{
			name: "model with original name fallback",
			hfInfo: &hfModelInfo{
				Author: "author-3",
			},
			sourceId:          "source-3",
			originalModelName: "fallback-original-name",
			expectedName:      "fallback-original-name",
			expectedProvider:  stringPtr("author-3"),
		},
		{
			name: "model with license in tags",
			hfInfo: &hfModelInfo{
				ID:   "test/licensed-model",
				Tags: []string{"license:apache-2.0", "other-tag"},
			},
			sourceId:          "source-4",
			originalModelName: "test/licensed-model",
			expectedName:      "test/licensed-model",
			expectedLicense:   stringPtr("Apache 2.0"),
			hasCustomProps:    true,
		},
		{
			name: "model with tasks",
			hfInfo: &hfModelInfo{
				ID:          "test/task-model",
				Task:        "text-classification",
				PipelineTag: "sentiment-analysis",
			},
			sourceId:          "source-5",
			originalModelName: "test/task-model",
			expectedName:      "test/task-model",
			hasTasks:          true,
		},
		{
			name: "model with description in cardData",
			hfInfo: &hfModelInfo{
				ID: "test/desc-model",
				CardData: &hfCard{
					Data: map[string]interface{}{
						"description": "This is a test description",
					},
				},
			},
			sourceId:          "source-6",
			originalModelName: "test/desc-model",
			expectedName:      "test/desc-model",
			hasDescription:    true,
		},
		{
			name: "minimal model info",
			hfInfo: &hfModelInfo{
				ID: "minimal/model",
			},
			sourceId:          "source-7",
			originalModelName: "minimal/model",
			expectedName:      "minimal/model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock provider with HTTP client to avoid nil pointer
			// Note: README fetching will fail, but that's expected in unit tests
			provider := &hfModelProvider{
				sourceId: tt.sourceId,
				client:   &http.Client{},
			}

			// Create hfModel and populate it
			hfm := &hfModel{}
			ctx := context.Background()
			hfm.populateFromHFInfo(ctx, provider, tt.hfInfo, tt.sourceId, tt.originalModelName)

			// Verify name
			if hfm.Name != tt.expectedName {
				t.Errorf("Name = %v, want %v", hfm.Name, tt.expectedName)
			}

			// Verify ExternalID
			if tt.expectedExternalID != nil {
				if hfm.ExternalId == nil || *hfm.ExternalId != *tt.expectedExternalID {
					t.Errorf("ExternalId = %v, want %v", hfm.ExternalId, tt.expectedExternalID)
				}
			} else if tt.hfInfo.ID != "" {
				// If hfInfo has ID, ExternalId should be set
				if hfm.ExternalId == nil || *hfm.ExternalId != tt.hfInfo.ID {
					t.Errorf("ExternalId = %v, want %v", hfm.ExternalId, tt.hfInfo.ID)
				}
			}

			// Verify SourceID
			if tt.expectedSourceID != nil {
				if hfm.SourceId == nil || *hfm.SourceId != *tt.expectedSourceID {
					t.Errorf("SourceId = %v, want %v", hfm.SourceId, tt.expectedSourceID)
				}
			} else if tt.sourceId != "" {
				if hfm.SourceId == nil || *hfm.SourceId != tt.sourceId {
					t.Errorf("SourceId = %v, want %v", hfm.SourceId, tt.sourceId)
				}
			}

			// Verify Provider
			if tt.expectedProvider != nil {
				if hfm.Provider == nil || *hfm.Provider != *tt.expectedProvider {
					t.Errorf("Provider = %v, want %v", hfm.Provider, tt.expectedProvider)
				}
			}

			// Verify License
			if tt.expectedLicense != nil {
				if hfm.License == nil || *hfm.License != *tt.expectedLicense {
					t.Errorf("License = %v, want %v", *hfm.License, *tt.expectedLicense)
				}
			}

			// Verify LibraryName
			if tt.expectedLibrary != nil {
				if hfm.LibraryName == nil || *hfm.LibraryName != *tt.expectedLibrary {
					t.Errorf("LibraryName = %v, want %v", hfm.LibraryName, tt.expectedLibrary)
				}
			}

			// Verify Tasks
			if tt.hasTasks {
				if len(hfm.Tasks) == 0 {
					t.Error("Expected tasks to be set, but got empty slice")
				}
			}

			// Verify Description
			if tt.hasDescription {
				if hfm.Description == nil {
					t.Error("Expected description to be set, but got nil")
				}
			}

			// Verify CustomProperties
			if tt.hasCustomProps {
				if hfm.GetCustomProperties() == nil || len(hfm.GetCustomProperties()) == 0 {
					t.Error("Expected custom properties to be set, but got nil or empty")
				}
			}

			// Verify timestamps if present
			if tt.hfInfo.CreatedAt != "" {
				if hfm.CreateTimeSinceEpoch == nil {
					t.Error("Expected CreateTimeSinceEpoch to be set")
				}
			}
			if tt.hfInfo.UpdatedAt != "" {
				if hfm.LastUpdateTimeSinceEpoch == nil {
					t.Error("Expected LastUpdateTimeSinceEpoch to be set")
				}
			}
		})
	}
}

func TestConvertHFModelToRecord(t *testing.T) {
	tests := []struct {
		name              string
		hfInfo            *hfModelInfo
		originalModelName string
		sourceId          string
		verifyFunc        func(t *testing.T, record ModelProviderRecord)
	}{
		{
			name: "complete model conversion",
			hfInfo: &hfModelInfo{
				ID:          "test-org/complete-model",
				Author:      "test-author",
				CreatedAt:   "2023-01-01T00:00:00Z",
				UpdatedAt:   "2023-01-02T00:00:00Z",
				Tags:        []string{"license:mit"},
				LibraryName: "transformers",
				Task:        "text-generation",
				CardData: &hfCard{
					Data: map[string]interface{}{
						"description": "A complete test model",
					},
				},
			},
			originalModelName: "test-org/complete-model",
			sourceId:          "test-source",
			verifyFunc: func(t *testing.T, record ModelProviderRecord) {
				if record.Model == nil {
					t.Fatal("Model should not be nil")
				}
				attrs := record.Model.GetAttributes()
				if attrs == nil {
					t.Fatal("Attributes should not be nil")
				}
				if attrs.Name == nil || *attrs.Name != "test-org/complete-model" {
					t.Errorf("Name = %v, want 'test-org/complete-model'", attrs.Name)
				}
				if attrs.ExternalID == nil || *attrs.ExternalID != "test-org/complete-model" {
					t.Errorf("ExternalID = %v, want 'test-org/complete-model'", attrs.ExternalID)
				}
				if attrs.CreateTimeSinceEpoch == nil {
					t.Error("CreateTimeSinceEpoch should be set")
				}
				if attrs.LastUpdateTimeSinceEpoch == nil {
					t.Error("LastUpdateTimeSinceEpoch should be set")
				}
				if record.Model.GetProperties() == nil || len(*record.Model.GetProperties()) == 0 {
					t.Error("Properties should be set")
				}
				// Should have one artifact with hf:// URI
				if len(record.Artifacts) != 1 {
					t.Errorf("Expected 1 artifact with hf:// URI, got %d", len(record.Artifacts))
				}
				if len(record.Artifacts) == 1 {
					artifact := record.Artifacts[0]
					if artifact.CatalogModelArtifact == nil {
						t.Error("CatalogModelArtifact should not be nil")
					}
					if artifact.CatalogModelArtifact.GetAttributes().URI == nil {
						t.Error("Artifact URI should not be nil")
					} else if !strings.HasPrefix(*artifact.CatalogModelArtifact.GetAttributes().URI, "hf://") {
						t.Errorf("Artifact URI should start with hf://, got %s", *artifact.CatalogModelArtifact.GetAttributes().URI)
					}
				}
			},
		},
		{
			name: "minimal model conversion",
			hfInfo: &hfModelInfo{
				ID: "minimal/model",
			},
			originalModelName: "minimal/model",
			sourceId:          "source-1",
			verifyFunc: func(t *testing.T, record ModelProviderRecord) {
				if record.Model == nil {
					t.Fatal("Model should not be nil")
				}
				attrs := record.Model.GetAttributes()
				if attrs == nil || attrs.Name == nil {
					t.Fatal("Attributes and Name should not be nil")
				}
				if *attrs.Name != "minimal/model" {
					t.Errorf("Name = %v, want 'minimal/model'", attrs.Name)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &hfModelProvider{
				sourceId: tt.sourceId,
			}
			ctx := context.Background()
			record := provider.convertHFModelToRecord(ctx, tt.hfInfo, tt.originalModelName)
			tt.verifyFunc(t, record)
		})
	}
}

func TestConvertHFModelProperties(t *testing.T) {
	tests := []struct {
		name         string
		catalogModel *apimodels.CatalogModel
		wantProps    bool
		wantCustom   bool
		verifyFunc   func(t *testing.T, props []models.Properties, customProps []models.Properties)
	}{
		{
			name: "model with all properties",
			catalogModel: &apimodels.CatalogModel{
				Name:        "test-model",
				Description: stringPtr("Test description"),
				Readme:      stringPtr("# Test README"),
				Provider:    stringPtr("test-provider"),
				License:     stringPtr("MIT License"),
				LibraryName: stringPtr("transformers"),
				SourceId:    stringPtr("source-1"),
				Tasks:       []string{"text-generation"},
			},
			wantProps:  true,
			wantCustom: false,
			verifyFunc: func(t *testing.T, props []models.Properties, customProps []models.Properties) {
				if len(props) == 0 {
					t.Error("Expected properties to be set")
				}
			},
		},
		{
			name: "model with custom properties",
			catalogModel: func() *apimodels.CatalogModel {
				model := &apimodels.CatalogModel{
					Name: "test-model",
				}
				customProps := map[string]apimodels.MetadataValue{
					"hf_tags": {
						MetadataStringValue: &apimodels.MetadataStringValue{
							StringValue: `["tag1","tag2"]`,
						},
					},
				}
				model.SetCustomProperties(customProps)
				return model
			}(),
			wantProps:  false,
			wantCustom: true,
			verifyFunc: func(t *testing.T, props []models.Properties, customProps []models.Properties) {
				if len(customProps) == 0 {
					t.Error("Expected custom properties to be set")
				}
			},
		},
		{
			name: "model with minimal properties",
			catalogModel: &apimodels.CatalogModel{
				Name: "minimal-model",
			},
			wantProps:  false,
			wantCustom: false,
			verifyFunc: func(t *testing.T, props []models.Properties, customProps []models.Properties) {
				if len(props) != 0 {
					t.Errorf("Expected no properties, got %d", len(props))
				}
				if len(customProps) != 0 {
					t.Errorf("Expected no custom properties, got %d", len(customProps))
				}
			},
		},
		{
			name: "model with tasks",
			catalogModel: &apimodels.CatalogModel{
				Name:  "task-model",
				Tasks: []string{"classification", "generation"},
			},
			wantProps: true,
			verifyFunc: func(t *testing.T, props []models.Properties, customProps []models.Properties) {
				// Should have tasks property
				foundTasks := false
				for _, prop := range props {
					if prop.Name == "tasks" {
						foundTasks = true
						break
					}
				}
				if !foundTasks {
					t.Error("Expected tasks property to be present")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			props, customProps := convertHFModelProperties(tt.catalogModel)
			if (len(props) > 0) != tt.wantProps {
				t.Errorf("Properties presence = %v, want %v", len(props) > 0, tt.wantProps)
			}
			if (len(customProps) > 0) != tt.wantCustom {
				t.Errorf("Custom properties presence = %v, want %v", len(customProps) > 0, tt.wantCustom)
			}
			if tt.verifyFunc != nil {
				tt.verifyFunc(t, props, customProps)
			}
		})
	}
}

func TestHFModelProviderWithModelFilter(t *testing.T) {
	tests := []struct {
		name           string
		includedModels []string
		excludedModels []string
		modelName      string
		wantAllowed    bool
		description    string
	}{
		{
			name:           "model matches included pattern",
			includedModels: []string{"ibm-granite/*"},
			excludedModels: nil,
			modelName:      "ibm-granite/granite-4.0-h-small",
			wantAllowed:    true,
			description:    "Model matching included pattern should be allowed",
		},
		{
			name:           "model does not match included pattern",
			includedModels: []string{"ibm-granite/*"},
			excludedModels: nil,
			modelName:      "meta-llama/Llama-3.2-1B",
			wantAllowed:    false,
			description:    "Model not matching included pattern should be excluded",
		},
		{
			name:           "model matches excluded pattern",
			includedModels: []string{"ibm-granite/*"},
			excludedModels: []string{"*-beta"},
			modelName:      "ibm-granite/granite-4.0-h-beta",
			wantAllowed:    false,
			description:    "Model matching excluded pattern should be excluded even if it matches included",
		},
		{
			name:           "model matches included but not excluded",
			includedModels: []string{"ibm-granite/*"},
			excludedModels: []string{"*-beta"},
			modelName:      "ibm-granite/granite-4.0-h-small",
			wantAllowed:    true,
			description:    "Model matching included but not excluded should be allowed",
		},
		{
			name:           "case insensitive matching",
			includedModels: []string{"IBM-Granite/*"},
			excludedModels: nil,
			modelName:      "ibm-granite/granite-4.0-h-small",
			wantAllowed:    true,
			description:    "Filtering should be case-insensitive",
		},
		{
			name:           "no included patterns allows all",
			includedModels: nil,
			excludedModels: []string{"*-beta"},
			modelName:      "test/model",
			wantAllowed:    true,
			description:    "No included patterns means all models are allowed (unless excluded)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a filter with the test patterns
			filter, err := NewModelFilter(tt.includedModels, tt.excludedModels)
			if err != nil {
				t.Fatalf("NewModelFilter() error = %v, want nil", err)
			}

			// Create a provider with the filter
			provider := &hfModelProvider{
				filter: filter,
			}

			// Test that the filter works correctly
			got := provider.filter.Allows(tt.modelName)
			if got != tt.wantAllowed {
				t.Errorf("ModelFilter.Allows(%q) = %v, want %v. %s", tt.modelName, got, tt.wantAllowed, tt.description)
			}
		})
	}
}

func TestPopulateFromHFInfoWithCustomProperties(t *testing.T) {
	hfInfo := &hfModelInfo{
		ID:        "test/custom-props-model",
		Sha:       "sha123",
		Downloads: 5000,
		Tags:      []string{"tag1", "tag2", "license:apache-2.0"},
		Config: &hfConfig{
			Architectures: []string{"BertModel", "BertForSequenceClassification"},
			ModelType:     "bert",
		},
	}

	provider := &hfModelProvider{
		sourceId: "test-source",
	}

	hfm := &hfModel{}
	ctx := context.Background()
	hfm.populateFromHFInfo(ctx, provider, hfInfo, "test-source", "test/custom-props-model")

	customProps := hfm.GetCustomProperties()
	if customProps == nil {
		t.Fatal("Custom properties should not be nil")
	}

	// Verify hf_tags
	if tagsVal, ok := customProps["hf_tags"]; !ok {
		t.Error("Expected hf_tags in custom properties")
	} else if tagsVal.MetadataStringValue == nil {
		t.Error("hf_tags should be a string value")
	}

	// Verify hf_architectures
	if archVal, ok := customProps["hf_architectures"]; !ok {
		t.Error("Expected hf_architectures in custom properties")
	} else if archVal.MetadataStringValue == nil {
		t.Error("hf_architectures should be a string value")
	}

	// Verify hf_model_type
	if modelTypeVal, ok := customProps["hf_model_type"]; !ok {
		t.Error("Expected hf_model_type in custom properties")
	} else if modelTypeVal.MetadataStringValue == nil || modelTypeVal.MetadataStringValue.StringValue != "bert" {
		t.Errorf("hf_model_type = %v, want 'bert'", modelTypeVal.MetadataStringValue)
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

func TestParseModelPattern(t *testing.T) {
	tests := []struct {
		pattern      string
		expectedType PatternType
		expectedOrg  string
		expectedPfx  string
	}{
		// Exact patterns
		{"meta-llama/Llama-2-7b-chat", PatternExact, "", ""},
		{"gpt2", PatternExact, "", ""},
		{"openai-community/gpt2", PatternExact, "", ""},

		// Org/* patterns
		{"ibm-granite/*", PatternOrgAll, "ibm-granite", ""},
		{"meta-llama/*", PatternOrgAll, "meta-llama", ""},
		{"openai/*", PatternOrgAll, "openai", ""},

		// Org/prefix* patterns
		{"meta-llama/Llama-2-*", PatternOrgPrefix, "meta-llama", "Llama-2-"},
		{"ibm-granite/granite-3*", PatternOrgPrefix, "ibm-granite", "granite-3"},
		{"mistralai/Mistral-*", PatternOrgPrefix, "mistralai", "Mistral-"},

		// Invalid patterns - would try to list all HuggingFace models
		{"*", PatternInvalid, "", ""},
		{"*/*", PatternInvalid, "", ""},
		{"*/something", PatternInvalid, "", ""},
		{"*/prefix*", PatternInvalid, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			pType, org, prefix := parseModelPattern(tt.pattern)
			assert.Equal(t, tt.expectedType, pType, "pattern type mismatch")
			assert.Equal(t, tt.expectedOrg, org, "org mismatch")
			assert.Equal(t, tt.expectedPfx, prefix, "prefix mismatch")
		})
	}
}

func TestParseNextCursor(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected string
	}{
		{
			name:     "empty header",
			header:   "",
			expected: "",
		},
		{
			name:     "valid next link",
			header:   `<https://huggingface.co/api/models?author=ibm-granite&limit=100&cursor=abc123>; rel="next"`,
			expected: "abc123",
		},
		{
			name:     "next link with other params",
			header:   `<https://huggingface.co/api/models?author=ibm-granite&cursor=xyz789&limit=100>; rel="next"`,
			expected: "xyz789",
		},
		{
			name:     "multiple links",
			header:   `<https://example.com/first>; rel="first", <https://huggingface.co/api/models?cursor=page2>; rel="next"`,
			expected: "page2",
		},
		{
			name:     "no next link",
			header:   `<https://example.com/first>; rel="first"`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseNextCursor(tt.header)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestListModelsByAuthor(t *testing.T) {
	// Setup mock HF server
	mux := http.NewServeMux()
	callCount := 0

	// Mock /api/whoami-v2 for credential validation
	mux.HandleFunc("/api/whoami-v2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"name": "test-user"})
	})

	// Mock /api/models list endpoint with pagination
	mux.HandleFunc("/api/models", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		author := r.URL.Query().Get("author")
		cursor := r.URL.Query().Get("cursor")
		search := r.URL.Query().Get("search")

		if author == "test-org" {
			switch cursor {
			case "":
				// First page - return 100 items to simulate full page (triggers pagination)
				models := make([]map[string]interface{}, 100)
				for i := range 100 {
					models[i] = map[string]interface{}{"id": fmt.Sprintf("test-org/model-%d", i+1)}
				}
				// Add Link header for next page
				w.Header().Set("Link", `<https://huggingface.co/api/models?author=test-org&cursor=page2>; rel="next"`)
				_ = json.NewEncoder(w).Encode(models)
			case "page2":
				// Second page (last) - return fewer than 100 to indicate end
				models := []map[string]interface{}{
					{"id": "test-org/model-101"},
					{"id": "test-org/model-102"},
				}
				_ = json.NewEncoder(w).Encode(models)
			}
		} else if author == "search-org" && search != "" {
			// Search results
			models := []map[string]interface{}{
				{"id": "search-org/" + search + "-match1"},
				{"id": "search-org/" + search + "-match2"},
				{"id": "search-org/other-model"}, // Should be filtered out
			}
			_ = json.NewEncoder(w).Encode(models)
		} else {
			_ = json.NewEncoder(w).Encode([]map[string]interface{}{})
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	os.Setenv("HF_API_KEY", "test-api-key")
	defer os.Unsetenv("HF_API_KEY")

	t.Run("lists all models from org with pagination", func(t *testing.T) {
		callCount = 0
		config := &PreviewConfig{
			Type: "hf",
			Properties: map[string]any{
				"url": server.URL,
			},
		}

		provider, err := NewHFPreviewProvider(config)
		require.NoError(t, err)

		models, err := provider.listModelsByAuthor(context.Background(), "test-org", "")
		require.NoError(t, err)

		// Should have 100 from first page + 2 from second page = 102
		assert.Len(t, models, 102)
		assert.Contains(t, models, "test-org/model-1")
		assert.Contains(t, models, "test-org/model-100")
		assert.Contains(t, models, "test-org/model-101")
		assert.Contains(t, models, "test-org/model-102")

		// Should have made 2 API calls (2 pages)
		assert.Equal(t, 2, callCount)
	})

	t.Run("filters by search prefix", func(t *testing.T) {
		config := &PreviewConfig{
			Type: "hf",
			Properties: map[string]any{
				"url": server.URL,
			},
		}

		provider, err := NewHFPreviewProvider(config)
		require.NoError(t, err)

		models, err := provider.listModelsByAuthor(context.Background(), "search-org", "prefix")
		require.NoError(t, err)

		// Should only include models starting with "prefix"
		assert.Len(t, models, 2)
		assert.Contains(t, models, "search-org/prefix-match1")
		assert.Contains(t, models, "search-org/prefix-match2")
		// "other-model" should be filtered out
		assert.NotContains(t, models, "search-org/other-model")
	})

	t.Run("respects maxModels limit", func(t *testing.T) {
		callCount = 0
		config := &PreviewConfig{
			Type: "hf",
			Properties: map[string]any{
				"url":       server.URL,
				"maxModels": 50, // Limit to 50 models
			},
		}

		provider, err := NewHFPreviewProvider(config)
		require.NoError(t, err)
		assert.Equal(t, 50, provider.maxModels)

		models, err := provider.listModelsByAuthor(context.Background(), "test-org", "")
		require.NoError(t, err)

		// Should stop at 50 models (first page has 100, but we limit to 50)
		assert.Len(t, models, 50)

		// Should have only made 1 API call (stopped before second page)
		assert.Equal(t, 1, callCount)
	})

	t.Run("uses default maxModels when not specified", func(t *testing.T) {
		config := &PreviewConfig{
			Type: "hf",
			Properties: map[string]any{
				"url": server.URL,
			},
		}

		provider, err := NewHFPreviewProvider(config)
		require.NoError(t, err)

		// Should use default (500)
		assert.Equal(t, 500, provider.maxModels)
	})

	t.Run("maxModels 0 means no limit", func(t *testing.T) {
		callCount = 0
		config := &PreviewConfig{
			Type: "hf",
			Properties: map[string]any{
				"url":       server.URL,
				"maxModels": 0, // No limit
			},
		}

		provider, err := NewHFPreviewProvider(config)
		require.NoError(t, err)
		assert.Equal(t, 0, provider.maxModels)

		models, err := provider.listModelsByAuthor(context.Background(), "test-org", "")
		require.NoError(t, err)

		// Should get all 102 models (100 from page 1 + 2 from page 2)
		assert.Len(t, models, 102)

		// Should have made 2 API calls
		assert.Equal(t, 2, callCount)
	})
}

func TestFetchModelNamesForPreviewWithPatterns(t *testing.T) {
	// Setup mock HF server
	mux := http.NewServeMux()

	mux.HandleFunc("/api/whoami-v2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"name": "test-user"})
	})

	// Mock list API
	mux.HandleFunc("/api/models", func(w http.ResponseWriter, r *http.Request) {
		author := r.URL.Query().Get("author")
		if author == "test-org" {
			models := []map[string]interface{}{
				{"id": "test-org/model-a"},
				{"id": "test-org/model-b"},
			}
			_ = json.NewEncoder(w).Encode(models)
		} else {
			_ = json.NewEncoder(w).Encode([]map[string]interface{}{})
		}
	})

	// Mock individual model endpoints
	mux.HandleFunc("/api/models/exact-org/exact-model", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"id": "exact-org/exact-model"})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	os.Setenv("HF_API_KEY", "test-api-key")
	defer os.Unsetenv("HF_API_KEY")

	t.Run("mixed patterns: org/* and exact", func(t *testing.T) {
		config := &PreviewConfig{
			Type: "hf",
			IncludedModels: []string{
				"test-org/*",            // Should use list API
				"exact-org/exact-model", // Should use direct fetch
			},
			Properties: map[string]any{
				"url": server.URL,
			},
		}

		provider, err := NewHFPreviewProvider(config)
		require.NoError(t, err)

		names, err := provider.FetchModelNamesForPreview(context.Background(), config.IncludedModels)
		require.NoError(t, err)

		assert.Len(t, names, 3)
		assert.Contains(t, names, "test-org/model-a")
		assert.Contains(t, names, "test-org/model-b")
		assert.Contains(t, names, "exact-org/exact-model")
	})

	t.Run("rejects * wildcard pattern", func(t *testing.T) {
		config := &PreviewConfig{
			Type: "hf",
			IncludedModels: []string{
				"*",
			},
			Properties: map[string]any{
				"url": server.URL,
			},
		}

		provider, err := NewHFPreviewProvider(config)
		require.NoError(t, err)

		_, err = provider.FetchModelNamesForPreview(context.Background(), config.IncludedModels)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "wildcard pattern")
		assert.Contains(t, err.Error(), "not supported")
	})

	t.Run("rejects */* wildcard pattern", func(t *testing.T) {
		config := &PreviewConfig{
			Type: "hf",
			IncludedModels: []string{
				"*/*",
			},
			Properties: map[string]any{
				"url": server.URL,
			},
		}

		provider, err := NewHFPreviewProvider(config)
		require.NoError(t, err)

		_, err = provider.FetchModelNamesForPreview(context.Background(), config.IncludedModels)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "wildcard pattern")
		assert.Contains(t, err.Error(), "not supported")
	})

	t.Run("rejects */prefix pattern", func(t *testing.T) {
		config := &PreviewConfig{
			Type: "hf",
			IncludedModels: []string{
				"*/Llama-*",
			},
			Properties: map[string]any{
				"url": server.URL,
			},
		}

		provider, err := NewHFPreviewProvider(config)
		require.NoError(t, err)

		_, err = provider.FetchModelNamesForPreview(context.Background(), config.IncludedModels)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "wildcard pattern")
		assert.Contains(t, err.Error(), "not supported")
	})
}

func TestPreviewSourceModelsWithHFPatterns(t *testing.T) {
	// Setup mock HF server
	mux := http.NewServeMux()

	mux.HandleFunc("/api/whoami-v2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"name": "test-user"})
	})

	mux.HandleFunc("/api/models", func(w http.ResponseWriter, r *http.Request) {
		author := r.URL.Query().Get("author")
		if author == "test-org" {
			models := []map[string]interface{}{
				{"id": "test-org/model-stable"},
				{"id": "test-org/model-experimental"},
				{"id": "test-org/model-draft"},
			}
			_ = json.NewEncoder(w).Encode(models)
		} else {
			_ = json.NewEncoder(w).Encode([]map[string]interface{}{})
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	os.Setenv("HF_API_KEY", "test-api-key")
	defer os.Unsetenv("HF_API_KEY")

	t.Run("org/* pattern with excludedModels filter", func(t *testing.T) {
		// Note: We test the filtering logic by calling NewHFPreviewProvider directly
		// rather than PreviewSourceModels, because PreviewSourceModels has SSRF
		// protection that removes custom URLs (which breaks mock server testing).
		// The filtering behavior is the same - we're testing the filter logic here.
		includedModels := []string{"test-org/*"}
		excludedModels := []string{"*-experimental", "*-draft"}

		config := &PreviewConfig{
			Type:           "hf",
			IncludedModels: includedModels,
			ExcludedModels: excludedModels,
			Properties: map[string]any{
				"url": server.URL,
			},
		}

		// Create provider and fetch model names (bypassing SSRF protection for testing)
		provider, err := NewHFPreviewProvider(config)
		require.NoError(t, err)

		modelNames, err := provider.FetchModelNamesForPreview(context.Background(), includedModels)
		require.NoError(t, err)
		require.Len(t, modelNames, 3)

		// Create filter and apply it (same logic as PreviewSourceModels)
		filter, err := NewModelFilter(includedModels, excludedModels)
		require.NoError(t, err)

		var included, excluded []string
		for _, name := range modelNames {
			if filter.Allows(name) {
				included = append(included, name)
			} else {
				excluded = append(excluded, name)
			}
		}

		assert.Len(t, included, 1)
		assert.Contains(t, included, "test-org/model-stable")

		assert.Len(t, excluded, 2)
		assert.Contains(t, excluded, "test-org/model-experimental")
		assert.Contains(t, excluded, "test-org/model-draft")
	})
}

func TestConvertHFModelToRecord_CreatesArtifactWithHFProtocol(t *testing.T) {
	tests := []struct {
		name               string
		hfInfo             *hfModelInfo
		expectedArtifacts  int
		expectedURI        string
		expectedExternalID string
	}{
		{
			name: "creates artifact with hf:// URI for valid model",
			hfInfo: &hfModelInfo{
				ID:        "meta-llama/Llama-2-7b",
				Author:    "meta",
				ModelID:   "meta-llama/Llama-2-7b",
				CreatedAt: "2023-01-01T00:00:00Z",
			},
			expectedArtifacts:  1,
			expectedURI:        "hf://meta-llama/Llama-2-7b",
			expectedExternalID: "meta-llama/Llama-2-7b",
		},
		{
			name: "creates artifact for gated model",
			hfInfo: &hfModelInfo{
				ID:      "ibm-granite/granite-3.0-8b-instruct",
				Author:  "ibm-granite",
				ModelID: "ibm-granite/granite-3.0-8b-instruct",
				Gated:   gatedString("auto"),
			},
			expectedArtifacts:  1,
			expectedURI:        "hf://ibm-granite/granite-3.0-8b-instruct",
			expectedExternalID: "ibm-granite/granite-3.0-8b-instruct",
		},
		{
			name: "no artifact when ExternalId is nil",
			hfInfo: &hfModelInfo{
				ModelID: "test-model",
				Author:  "test-author",
				// ID is empty, so ExternalId will be nil
			},
			expectedArtifacts: 0,
		},
		{
			name: "no artifact when ExternalId is empty string",
			hfInfo: &hfModelInfo{
				ID:      "",
				ModelID: "test-model",
				Author:  "test-author",
			},
			expectedArtifacts: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &hfModelProvider{
				client:   &http.Client{},
				sourceId: "test-source",
			}

			ctx := context.Background()
			record := provider.convertHFModelToRecord(ctx, tt.hfInfo, "original-model-name")

			// Check artifact count
			assert.Len(t, record.Artifacts, tt.expectedArtifacts)

			if tt.expectedArtifacts > 0 {
				// Verify artifact structure
				artifact := record.Artifacts[0]
				require.NotNil(t, artifact.CatalogModelArtifact)
				require.NotNil(t, artifact.CatalogModelArtifact.GetAttributes())

				attrs := artifact.CatalogModelArtifact.GetAttributes()

				// Check URI has hf:// prefix
				require.NotNil(t, attrs.URI)
				assert.Equal(t, tt.expectedURI, *attrs.URI)

				// Check artifact type
				require.NotNil(t, attrs.ArtifactType)
				assert.Equal(t, "model-artifact", *attrs.ArtifactType)

				// Check external ID matches model ID
				require.NotNil(t, attrs.ExternalID)
				assert.Equal(t, tt.expectedExternalID, *attrs.ExternalID)

				// Check artifact name is set
				require.NotNil(t, attrs.Name)
				assert.Contains(t, *attrs.Name, "-hf-artifact")
			}
		})
	}
}

func TestConvertHFModelToRecord_ArtifactTimestamps(t *testing.T) {
	hfInfo := &hfModelInfo{
		ID:        "test-org/test-model",
		Author:    "test-author",
		CreatedAt: "2023-01-01T00:00:00Z",
		UpdatedAt: "2023-06-01T12:00:00Z",
	}

	provider := &hfModelProvider{
		client:   &http.Client{},
		sourceId: "test-source",
	}

	ctx := context.Background()
	record := provider.convertHFModelToRecord(ctx, hfInfo, "test-org/test-model")

	require.Len(t, record.Artifacts, 1)
	artifact := record.Artifacts[0]
	attrs := artifact.CatalogModelArtifact.GetAttributes()

	// Verify timestamps are copied from model
	require.NotNil(t, attrs.CreateTimeSinceEpoch)
	require.NotNil(t, attrs.LastUpdateTimeSinceEpoch)

	// Timestamps should match the model's timestamps
	require.NotNil(t, record.Model.GetAttributes().CreateTimeSinceEpoch)
	require.NotNil(t, record.Model.GetAttributes().LastUpdateTimeSinceEpoch)

	assert.Equal(t, *record.Model.GetAttributes().CreateTimeSinceEpoch, *attrs.CreateTimeSinceEpoch)
	assert.Equal(t, *record.Model.GetAttributes().LastUpdateTimeSinceEpoch, *attrs.LastUpdateTimeSinceEpoch)
}
