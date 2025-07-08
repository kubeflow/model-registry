package openapi

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindSources(t *testing.T) {
	// Setup test cases
	testCases := []struct {
		name           string
		catalogs       map[string]catalog.CatalogSource
		nameFilter     string
		pageSize       string
		orderBy        model.OrderByField
		sortOrder      model.SortOrder
		nextPageToken  string
		expectedStatus int
		expectedSize   int32
		expectedItems  int
		checkSorting   bool
	}{
		{
			name:           "Empty catalog list",
			catalogs:       map[string]catalog.CatalogSource{},
			nameFilter:     "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedSize:   0,
			expectedItems:  0,
		},
		{
			name: "Single catalog",
			catalogs: map[string]catalog.CatalogSource{
				"catalog1": {
					Metadata: model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1"},
				},
			},
			nameFilter:     "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedSize:   1,
			expectedItems:  1,
		},
		{
			name: "Multiple catalogs with no filter",
			catalogs: map[string]catalog.CatalogSource{
				"catalog1": {
					Metadata: model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1"},
				},
				"catalog2": {
					Metadata: model.CatalogSource{Id: "catalog2", Name: "Test Catalog 2"},
				},
				"catalog3": {
					Metadata: model.CatalogSource{Id: "catalog3", Name: "Another Catalog"},
				},
			},
			nameFilter:     "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedSize:   3,
			expectedItems:  3,
		},
		{
			name: "Filter by name",
			catalogs: map[string]catalog.CatalogSource{
				"catalog1": {
					Metadata: model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1"},
				},
				"catalog2": {
					Metadata: model.CatalogSource{Id: "catalog2", Name: "Test Catalog 2"},
				},
				"catalog3": {
					Metadata: model.CatalogSource{Id: "catalog3", Name: "Another Catalog"},
				},
			},
			nameFilter:     "Test",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedSize:   2,
			expectedItems:  2,
		},
		{
			name: "Filter by name case insensitive",
			catalogs: map[string]catalog.CatalogSource{
				"catalog1": {
					Metadata: model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1"},
				},
				"catalog2": {
					Metadata: model.CatalogSource{Id: "catalog2", Name: "Test Catalog 2"},
				},
				"catalog3": {
					Metadata: model.CatalogSource{Id: "catalog3", Name: "Another Catalog"},
				},
			},
			nameFilter:     "test",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedSize:   2,
			expectedItems:  2,
		},
		{
			name: "Pagination - limit results",
			catalogs: map[string]catalog.CatalogSource{
				"catalog1": {
					Metadata: model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1"},
				},
				"catalog2": {
					Metadata: model.CatalogSource{Id: "catalog2", Name: "Test Catalog 2"},
				},
				"catalog3": {
					Metadata: model.CatalogSource{Id: "catalog3", Name: "Another Catalog"},
				},
			},
			nameFilter:     "",
			pageSize:       "2",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedSize:   3, // Total size should be 3
			expectedItems:  2, // But only 2 items returned due to page size
		},
		{
			name: "Default page size",
			catalogs: map[string]catalog.CatalogSource{
				"catalog1": {
					Metadata: model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1"},
				},
				"catalog2": {
					Metadata: model.CatalogSource{Id: "catalog2", Name: "Test Catalog 2"},
				},
			},
			nameFilter:     "",
			pageSize:       "", // Empty to test default
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedSize:   2,
			expectedItems:  2,
		},
		{
			name: "Invalid page size",
			catalogs: map[string]catalog.CatalogSource{
				"catalog1": {
					Metadata: model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1"},
				},
			},
			nameFilter:     "",
			pageSize:       "invalid",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Sort by ID ascending",
			catalogs: map[string]catalog.CatalogSource{
				"catalog2": {
					Metadata: model.CatalogSource{Id: "catalog2", Name: "B Catalog"},
				},
				"catalog1": {
					Metadata: model.CatalogSource{Id: "catalog1", Name: "A Catalog"},
				},
				"catalog3": {
					Metadata: model.CatalogSource{Id: "catalog3", Name: "C Catalog"},
				},
			},
			nameFilter:     "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedSize:   3,
			expectedItems:  3,
			checkSorting:   true,
		},
		{
			name: "Sort by ID descending",
			catalogs: map[string]catalog.CatalogSource{
				"catalog2": {
					Metadata: model.CatalogSource{Id: "catalog2", Name: "B Catalog"},
				},
				"catalog1": {
					Metadata: model.CatalogSource{Id: "catalog1", Name: "A Catalog"},
				},
				"catalog3": {
					Metadata: model.CatalogSource{Id: "catalog3", Name: "C Catalog"},
				},
			},
			nameFilter:     "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_DESC,
			expectedStatus: http.StatusOK,
			expectedSize:   3,
			expectedItems:  3,
			checkSorting:   true,
		},
		{
			name: "Sort by name ascending",
			catalogs: map[string]catalog.CatalogSource{
				"catalog2": {
					Metadata: model.CatalogSource{Id: "catalog2", Name: "B Catalog"},
				},
				"catalog1": {
					Metadata: model.CatalogSource{Id: "catalog1", Name: "A Catalog"},
				},
				"catalog3": {
					Metadata: model.CatalogSource{Id: "catalog3", Name: "C Catalog"},
				},
			},
			nameFilter:     "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_NAME,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedSize:   3,
			expectedItems:  3,
			checkSorting:   true,
		},
		{
			name: "Sort by name descending",
			catalogs: map[string]catalog.CatalogSource{
				"catalog2": {
					Metadata: model.CatalogSource{Id: "catalog2", Name: "B Catalog"},
				},
				"catalog1": {
					Metadata: model.CatalogSource{Id: "catalog1", Name: "A Catalog"},
				},
				"catalog3": {
					Metadata: model.CatalogSource{Id: "catalog3", Name: "C Catalog"},
				},
			},
			nameFilter:     "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_NAME,
			sortOrder:      model.SORTORDER_DESC,
			expectedStatus: http.StatusOK,
			expectedSize:   3,
			expectedItems:  3,
			checkSorting:   true,
		},
		{
			name: "Invalid sort order",
			catalogs: map[string]catalog.CatalogSource{
				"catalog1": {
					Metadata: model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1"},
				},
			},
			nameFilter:     "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      "INVALID",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid order by field",
			catalogs: map[string]catalog.CatalogSource{
				"catalog1": {
					Metadata: model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1"},
				},
			},
			nameFilter:     "",
			pageSize:       "10",
			orderBy:        "INVALID",
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Default sort order (ID ascending)",
			catalogs: map[string]catalog.CatalogSource{
				"catalog2": {
					Metadata: model.CatalogSource{Id: "catalog2", Name: "B Catalog"},
				},
				"catalog1": {
					Metadata: model.CatalogSource{Id: "catalog1", Name: "A Catalog"},
				},
				"catalog3": {
					Metadata: model.CatalogSource{Id: "catalog3", Name: "C Catalog"},
				},
			},
			nameFilter:     "",
			pageSize:       "10",
			orderBy:        "", // Empty to test default
			sortOrder:      "", // Empty to test default
			expectedStatus: http.StatusOK,
			expectedSize:   3,
			expectedItems:  3,
			checkSorting:   true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create service with test catalogs
			service := NewModelCatalogServiceAPIService(catalog.NewSourceCollection(tc.catalogs))

			// Call FindSources
			resp, err := service.FindSources(
				context.Background(),
				tc.nameFilter,
				tc.pageSize,
				tc.orderBy,
				tc.sortOrder,
				tc.nextPageToken,
			)

			// Check response status
			assert.Equal(t, tc.expectedStatus, resp.Code)

			// If we expect an error, we don't need to check the response body
			if tc.expectedStatus != http.StatusOK {
				assert.NotNil(t, err)
				return
			}

			// For successful responses, check the response body
			require.NotNil(t, resp.Body)

			// Type assertion to access the CatalogSourceList
			sourceList, ok := resp.Body.(model.CatalogSourceList)
			require.True(t, ok, "Response body should be a CatalogSourceList")

			// Check the size matches expected
			assert.Equal(t, tc.expectedSize, sourceList.Size)

			// Check the number of items matches expected
			assert.Equal(t, tc.expectedItems, len(sourceList.Items))

			// Check that page size is set correctly
			if tc.pageSize == "" {
				// Default page size should be 10
				assert.Equal(t, int32(10), sourceList.PageSize)
			} else if pageSizeInt, err := strconv.ParseInt(tc.pageSize, 10, 32); err == nil {
				assert.Equal(t, int32(pageSizeInt), sourceList.PageSize)
			}

			// Check sorting if required
			if tc.checkSorting && len(sourceList.Items) > 1 {
				switch tc.orderBy {
				case model.ORDERBYFIELD_ID, "":
					// Check ID sorting
					for i := 0; i < len(sourceList.Items)-1; i++ {
						if tc.sortOrder == "" || tc.sortOrder == model.SORTORDER_ASC {
							assert.LessOrEqual(t,
								sourceList.Items[i].Id,
								sourceList.Items[i+1].Id,
								"Items should be sorted by ID in ascending order")
						} else {
							assert.GreaterOrEqual(t,
								sourceList.Items[i].Id,
								sourceList.Items[i+1].Id,
								"Items should be sorted by ID in descending order")
						}
					}
				case model.ORDERBYFIELD_NAME:
					// Check name sorting
					for i := 0; i < len(sourceList.Items)-1; i++ {
						if tc.sortOrder == "" || tc.sortOrder == model.SORTORDER_ASC {
							assert.LessOrEqual(t,
								sourceList.Items[i].Name,
								sourceList.Items[i+1].Name,
								"Items should be sorted by name in ascending order")
						} else {
							assert.GreaterOrEqual(t,
								sourceList.Items[i].Name,
								sourceList.Items[i+1].Name,
								"Items should be sorted by name in descending order")
						}
					}
				}
			}
		})
	}
}

// Define a mock model provider
type mockModelProvider struct {
	models    map[string]*model.CatalogModel
	artifacts map[string][]model.CatalogModelArtifact
}

// Implement GetModel method for the mock provider
func (m *mockModelProvider) GetModel(ctx context.Context, name string) (*model.CatalogModel, error) {
	model, exists := m.models[name]
	if !exists {
		return nil, nil
	}
	return model, nil
}

func (m *mockModelProvider) ListModels(ctx context.Context, params catalog.ListModelsParams) (model.CatalogModelList, error) {
	return model.CatalogModelList{}, nil
}

func (m *mockModelProvider) GetArtifacts(ctx context.Context, name string) (*model.CatalogModelArtifactList, error) {
	artifacts, exists := m.artifacts[name]
	if !exists {
		return &model.CatalogModelArtifactList{
			Items:         []model.CatalogModelArtifact{},
			Size:          0,
			PageSize:      0, // Or a default page size if applicable
			NextPageToken: "",
		}, nil
	}
	return &model.CatalogModelArtifactList{
		Items:         artifacts,
		Size:          int32(len(artifacts)),
		PageSize:      int32(len(artifacts)),
		NextPageToken: "",
	}, nil
}

func TestGetModel(t *testing.T) {
	testCases := []struct {
		name           string
		sources        map[string]catalog.CatalogSource
		sourceID       string
		modelName      string
		expectedStatus int
		expectedModel  *model.CatalogModel
	}{
		{
			name: "Existing model in source",
			sources: map[string]catalog.CatalogSource{
				"source1": {
					Metadata: model.CatalogSource{Id: "source1", Name: "Test Source"},
					Provider: &mockModelProvider{
						models: map[string]*model.CatalogModel{
							"test-model": {
								Name: "test-model",
							},
						},
					},
				},
			},
			sourceID:       "source1",
			modelName:      "test-model",
			expectedStatus: http.StatusOK,
			expectedModel: &model.CatalogModel{
				Name: "test-model",
			},
		},
		{
			name: "Non-existing source",
			sources: map[string]catalog.CatalogSource{
				"source1": {
					Metadata: model.CatalogSource{Id: "source1", Name: "Test Source"},
				},
			},
			sourceID:       "source2",
			modelName:      "test-model",
			expectedStatus: http.StatusNotFound,
			expectedModel:  nil,
		},
		{
			name: "Existing source, non-existing model",
			sources: map[string]catalog.CatalogSource{
				"source1": {
					Metadata: model.CatalogSource{Id: "source1", Name: "Test Source"},
					Provider: &mockModelProvider{
						models: map[string]*model.CatalogModel{},
					},
				},
			},
			sourceID:       "source1",
			modelName:      "test-model",
			expectedStatus: http.StatusNotFound,
			expectedModel:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create service with test sources
			service := NewModelCatalogServiceAPIService(catalog.NewSourceCollection(tc.sources))

			// Call GetModel
			resp, _ := service.GetModel(
				context.Background(),
				tc.sourceID,
				tc.modelName,
			)

			// Check response status
			assert.Equal(t, tc.expectedStatus, resp.Code)

			// If we expect an error or not found, we don't need to check the response body
			if tc.expectedStatus != http.StatusOK {
				return
			}

			// For successful responses, check the response body
			require.NotNil(t, resp.Body)

			// Type assertion to access the Model
			model, ok := resp.Body.(*model.CatalogModel)
			require.True(t, ok, "Response body should be a Model")

			// Check the model details
			assert.Equal(t, tc.expectedModel.Name, model.Name)
		})
	}
}

func TestGetAllModelArtifacts(t *testing.T) {
	testCases := []struct {
		name              string
		sources           map[string]catalog.CatalogSource
		sourceID          string
		modelName         string
		expectedStatus    int
		expectedArtifacts []model.CatalogModelArtifact
	}{
		{
			name: "Existing artifacts for model in source",
			sources: map[string]catalog.CatalogSource{
				"source1": {
					Metadata: model.CatalogSource{Id: "source1", Name: "Test Source"},
					Provider: &mockModelProvider{
						artifacts: map[string][]model.CatalogModelArtifact{
							"test-model": {
								{
									Uri: "s3://bucket/artifact1",
								},
								{
									Uri: "s3://bucket/artifact2",
								},
							},
						},
					},
				},
			},
			sourceID:       "source1",
			modelName:      "test-model",
			expectedStatus: http.StatusOK,
			expectedArtifacts: []model.CatalogModelArtifact{
				{
					Uri: "s3://bucket/artifact1",
				},
				{
					Uri: "s3://bucket/artifact2",
				},
			},
		},
		{
			name: "Non-existing source",
			sources: map[string]catalog.CatalogSource{
				"source1": {
					Metadata: model.CatalogSource{Id: "source1", Name: "Test Source"},
				},
			},
			sourceID:          "source2",
			modelName:         "test-model",
			expectedStatus:    http.StatusNotFound,
			expectedArtifacts: nil,
		},
		{
			name: "Existing source, no artifacts for model",
			sources: map[string]catalog.CatalogSource{
				"source1": {
					Metadata: model.CatalogSource{Id: "source1", Name: "Test Source"},
					Provider: &mockModelProvider{
						artifacts: map[string][]model.CatalogModelArtifact{},
					},
				},
			},
			sourceID:          "source1",
			modelName:         "test-model",
			expectedStatus:    http.StatusOK,
			expectedArtifacts: []model.CatalogModelArtifact{}, // Should be an empty slice, not nil
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create service with test sources
			service := NewModelCatalogServiceAPIService(catalog.NewSourceCollection(tc.sources))

			// Call GetAllModelArtifacts
			resp, _ := service.GetAllModelArtifacts(
				context.Background(),
				tc.sourceID,
				tc.modelName,
			)

			// Check response status
			assert.Equal(t, tc.expectedStatus, resp.Code)

			// If we expect an error or not found, we don't need to check the response body
			if tc.expectedStatus != http.StatusOK {
				return
			}

			// For successful responses, check the response body
			require.NotNil(t, resp.Body)

			// Type assertion to access the list of artifacts
			artifactList, ok := resp.Body.(*model.CatalogModelArtifactList)
			require.True(t, ok, "Response body should be a CatalogModelArtifactList")

			// Check the artifacts
			assert.Equal(t, tc.expectedArtifacts, artifactList.Items)
			assert.Equal(t, int32(len(tc.expectedArtifacts)), artifactList.Size)
		})
	}
}
