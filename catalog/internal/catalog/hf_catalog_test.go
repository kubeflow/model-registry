package catalog

import (
	"context"
	"net/http"
	"testing"

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
			expectedLicense:   stringPtr("mit"),
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
			expectedLicense:   stringPtr("apache-2.0"),
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
					t.Errorf("License = %v, want %v", hfm.License, tt.expectedLicense)
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

func TestIsModelExcluded(t *testing.T) {
	tests := []struct {
		name      string
		modelName string
		patterns  []string
		want      bool
	}{
		{
			name:      "exact match",
			modelName: "test/model",
			patterns:  []string{"test/model"},
			want:      true,
		},
		{
			name:      "no match",
			modelName: "test/model",
			patterns:  []string{"other/model"},
			want:      false,
		},
		{
			name:      "prefix match with wildcard",
			modelName: "test/model-v1",
			patterns:  []string{"test/model*"},
			want:      true,
		},
		{
			name:      "prefix match without wildcard",
			modelName: "test/model-v1",
			patterns:  []string{"test/model"},
			want:      false,
		},
		{
			name:      "multiple patterns, one matches",
			modelName: "excluded/model",
			patterns:  []string{"other/model", "excluded/model"},
			want:      true,
		},
		{
			name:      "wildcard at start",
			modelName: "test/model",
			patterns:  []string{"*test*"},
			want:      false, // wildcard only works at end
		},
		{
			name:      "empty patterns",
			modelName: "test/model",
			patterns:  []string{},
			want:      false,
		},
		{
			name:      "wildcard matches exact",
			modelName: "test/model",
			patterns:  []string{"test/model*"},
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isModelExcluded(tt.modelName, tt.patterns); got != tt.want {
				t.Errorf("isModelExcluded() = %v, want %v", got, tt.want)
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
				if len(record.Artifacts) != 0 {
					t.Errorf("Artifacts should be empty, got %d", len(record.Artifacts))
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
				License:     stringPtr("mit"),
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

	// Verify hf_downloads
	if downloadsVal, ok := customProps["hf_downloads"]; !ok {
		t.Error("Expected hf_downloads in custom properties")
	} else if downloadsVal.MetadataIntValue == nil {
		t.Error("hf_downloads should be an int value")
	} else if downloadsVal.MetadataIntValue.IntValue != "5000" {
		t.Errorf("hf_downloads = %v, want '5000'", downloadsVal.MetadataIntValue.IntValue)
	}

	// Verify hf_sha
	if shaVal, ok := customProps["hf_sha"]; !ok {
		t.Error("Expected hf_sha in custom properties")
	} else if shaVal.MetadataStringValue == nil || shaVal.MetadataStringValue.StringValue != "sha123" {
		t.Errorf("hf_sha = %v, want 'sha123'", shaVal.MetadataStringValue)
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
