package openapi

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// timeToMillisStringPointer converts time.Time to *string representing milliseconds since epoch.
func timeToMillisStringPointer(t time.Time) *string {
	s := strconv.FormatInt(t.UnixMilli(), 10)
	return &s
}

// pointerOrDefault returns the value pointed to by p, or def if p is nil.
func pointerOrDefault(p *string, def string) string {
	if p == nil {
		return def
	}
	return *p
}

func TestFindModels(t *testing.T) {
	// Define common models for testing
	time1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	time3 := time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)
	time4 := time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC)

	// Updated model definitions to match OpenAPI schema (no direct Id or Published, use Name, CreateTime, LastUpdateTime)
	modelA := &model.CatalogModel{Name: "Model A", CreateTimeSinceEpoch: timeToMillisStringPointer(time1), LastUpdateTimeSinceEpoch: timeToMillisStringPointer(time4)}
	modelB := &model.CatalogModel{Name: "Model B", CreateTimeSinceEpoch: timeToMillisStringPointer(time2), LastUpdateTimeSinceEpoch: timeToMillisStringPointer(time3)}
	modelC := &model.CatalogModel{Name: "Another Model C", CreateTimeSinceEpoch: timeToMillisStringPointer(time3), LastUpdateTimeSinceEpoch: timeToMillisStringPointer(time2)}
	modelD := &model.CatalogModel{Name: "My Model D", CreateTimeSinceEpoch: timeToMillisStringPointer(time4), LastUpdateTimeSinceEpoch: timeToMillisStringPointer(time1)}

	testCases := []struct {
		name              string
		sourceID          string
		mockModels        map[string]*model.CatalogModel
		filterQuery       string
		q                 string
		pageSize          string
		orderBy           model.OrderByField
		sortOrder         model.SortOrder
		nextPageToken     string
		expectedStatus    int
		expectedModelList *model.CatalogModelList
	}{
		{
			name:     "Successful query with no filters",
			sourceID: "source1",
			mockModels: map[string]*model.CatalogModel{
				"modelA": modelA, "modelB": modelB, "modelC": modelC, "modelD": modelD,
			},
			q:              "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_NAME,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedModelList: &model.CatalogModelList{
				Items:         []model.CatalogModel{*modelC, *modelA, *modelB, *modelD}, // Sorted by Name ASC: Another Model C, Model A, Model B, My Model D
				Size:          4,
				PageSize:      10, // Default page size
				NextPageToken: "",
			},
		},
		{
			name:     "Filter by query 'Model'",
			sourceID: "source1",
			mockModels: map[string]*model.CatalogModel{
				"modelA": modelA, "modelB": modelB, "modelC": modelC, "modelD": modelD,
			},
			q:              "Model",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_NAME,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedModelList: &model.CatalogModelList{
				Items:         []model.CatalogModel{*modelC, *modelA, *modelB, *modelD}, // Corrected to include modelC and sorted by name ASC
				Size:          4,                                                        // Corrected from 3 to 4
				PageSize:      10,
				NextPageToken: "",
			},
		},
		{
			name:     "Filter by query 'model' (case insensitive)",
			sourceID: "source1",
			mockModels: map[string]*model.CatalogModel{
				"modelA": modelA, "modelB": modelB, "modelC": modelC, "modelD": modelD,
			},
			q:              "model",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_NAME,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedModelList: &model.CatalogModelList{
				Items:         []model.CatalogModel{*modelC, *modelA, *modelB, *modelD}, // Corrected to include modelC and sorted by name ASC
				Size:          4,                                                        // Corrected from 3 to 4
				PageSize:      10,
				NextPageToken: "",
			},
		},
		{
			name:     "Page size limit",
			sourceID: "source1",
			mockModels: map[string]*model.CatalogModel{
				"modelA": modelA, "modelB": modelB, "modelC": modelC, "modelD": modelD,
			},
			q:              "",
			pageSize:       "2",
			orderBy:        model.ORDERBYFIELD_NAME,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedModelList: &model.CatalogModelList{
				Items:         []model.CatalogModel{*modelC, *modelA}, // First 2 after sorting by Name ASC
				Size:          4,                                      // Total size remains 4
				PageSize:      2,
				NextPageToken: (&stringCursor{Value: "Model A", ID: "Model A"}).String(),
			},
		},
		{
			name:     "Sort by ID Descending (mocked as Name Descending)",
			sourceID: "source1",
			mockModels: map[string]*model.CatalogModel{
				"modelA": modelA, "modelB": modelB, "modelC": modelC, "modelD": modelD,
			},
			q:              "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_DESC,
			expectedStatus: http.StatusOK,
			expectedModelList: &model.CatalogModelList{
				Items:         []model.CatalogModel{*modelD, *modelB, *modelA, *modelC}, // Sorted by Name DESC
				Size:          4,
				PageSize:      10,
				NextPageToken: "",
			},
		},
		{
			name:     "Sort by CreateTime Ascending",
			sourceID: "source1",
			mockModels: map[string]*model.CatalogModel{
				"modelA": modelA, "modelB": modelB, "modelC": modelC, "modelD": modelD,
			},
			q:              "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_CREATE_TIME,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedModelList: &model.CatalogModelList{
				Items:         []model.CatalogModel{*modelA, *modelB, *modelC, *modelD}, // Sorted by CreateTime ASC
				Size:          4,
				PageSize:      10,
				NextPageToken: "",
			},
		},
		{
			name:     "Sort by LastUpdateTime Descending",
			sourceID: "source1",
			mockModels: map[string]*model.CatalogModel{
				"modelA": modelA, "modelB": modelB, "modelC": modelC, "modelD": modelD,
			},
			q:              "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_LAST_UPDATE_TIME,
			sortOrder:      model.SORTORDER_DESC,
			expectedStatus: http.StatusOK,
			expectedModelList: &model.CatalogModelList{
				Items:         []model.CatalogModel{*modelA, *modelB, *modelC, *modelD}, // Corrected to be sorted by LastUpdateTime DESC (modelA has latest time4, modelD has earliest time1)
				Size:          4,
				PageSize:      10,
				NextPageToken: "",
			},
		},
		{
			name:           "Invalid source ID",
			sourceID:       "unknown-source",
			mockModels:     map[string]*model.CatalogModel{},
			q:              "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK, // Changed from http.StatusNotFound to http.StatusOK with an empty list -- now the source ID is just a field in the CatalogModel
			expectedModelList: &model.CatalogModelList{
				Items:         []model.CatalogModel{},
				Size:          0,
				PageSize:      10,
				NextPageToken: "",
			},
		},
		{
			name:     "Invalid pageSize string",
			sourceID: "source1",
			mockModels: map[string]*model.CatalogModel{
				"modelA": modelA,
			},
			q:                 "",
			pageSize:          "abc",
			orderBy:           model.ORDERBYFIELD_ID,
			sortOrder:         model.SORTORDER_ASC,
			expectedStatus:    http.StatusBadRequest,
			expectedModelList: nil,
		},
		{
			name:     "Unsupported orderBy field",
			sourceID: "source1",
			mockModels: map[string]*model.CatalogModel{
				"modelA": modelA,
			},
			q:              "",
			pageSize:       "10",
			orderBy:        "UNSUPPORTED_FIELD",
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK, // Changed from http.StatusBadRequest to http.StatusOK -- in model registry we fallback to ID if the order by field is unsupported
			expectedModelList: &model.CatalogModelList{
				Items: []model.CatalogModel{
					*modelA,
				},
				Size:          1,
				PageSize:      10,
				NextPageToken: "",
			},
		},
		{
			name:     "Unsupported sortOrder field",
			sourceID: "source1",
			mockModels: map[string]*model.CatalogModel{
				"modelA": modelA,
			},
			q:              "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      "UNSUPPORTED_ORDER",
			expectedStatus: http.StatusOK, // Changed from http.StatusBadRequest to http.StatusOK -- in model registry we fallback to ASC if the sort order field is unsupported
			expectedModelList: &model.CatalogModelList{
				Items: []model.CatalogModel{
					*modelA,
				},
				Size:          1,
				PageSize:      10,
				NextPageToken: "",
			},
		},
		{
			name:           "Empty models in source",
			sourceID:       "source1",
			mockModels:     map[string]*model.CatalogModel{},
			q:              "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedModelList: &model.CatalogModelList{
				Items:         []model.CatalogModel{},
				Size:          0,
				PageSize:      10,
				NextPageToken: "",
			},
		},
		{
			name:     "Default sort (ID ascending) and default page size",
			sourceID: "source1",
			mockModels: map[string]*model.CatalogModel{
				"modelB": modelB, "modelA": modelA, "modelD": modelD, "modelC": modelC,
			},
			q:              "",
			pageSize:       "", // Default page size
			orderBy:        "", // Default order by ID
			sortOrder:      "", // Default sort order ASC
			expectedStatus: http.StatusOK,
			expectedModelList: &model.CatalogModelList{
				Items:         []model.CatalogModel{*modelC, *modelA, *modelB, *modelD}, // Sorted by Name ASC (as ID is mocked to use Name)
				Size:          4,
				PageSize:      10, // Default page size
				NextPageToken: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock source collection
			sources := catalog.NewSourceCollection()
			sources.Merge("",
				map[string]model.CatalogSource{
					"source1": model.CatalogSource{Id: "source1", Name: "Test Source 1"},
				},
			)

			provider := &mockModelProvider{
				models: tc.mockModels,
			}

			service := NewModelCatalogServiceAPIService(provider, sources)

			resp, err := service.FindModels(
				context.Background(),
				[]string{tc.sourceID},
				tc.q,
				tc.filterQuery,
				tc.pageSize,
				tc.orderBy,
				tc.sortOrder,
				tc.nextPageToken,
			)

			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.expectedStatus != http.StatusOK {
				assert.Error(t, err)
				return
			}

			require.NotNil(t, resp.Body)
			modelList, ok := resp.Body.(model.CatalogModelList)
			require.True(t, ok, "Response body should be a CatalogModelList")

			assert.Equal(t, tc.expectedModelList.Size, modelList.Size)
			assert.Equal(t, tc.expectedModelList.PageSize, modelList.PageSize)
			if !assert.Equal(t, tc.expectedModelList.NextPageToken, modelList.NextPageToken) && tc.expectedModelList.NextPageToken != "" {
				assert.Equal(t, decodeStringCursor(tc.expectedModelList.NextPageToken), decodeStringCursor(modelList.NextPageToken))
			}

			// Deep equality check for items
			assert.Equal(t, tc.expectedModelList.Items, modelList.Items)
		})
	}
}

func TestFindSources(t *testing.T) {
	// Setup test cases
	trueValue := true
	testCases := []struct {
		name           string
		catalogs       map[string]model.CatalogSource
		nameFilter     string
		pageSize       string
		orderBy        model.OrderByField
		sortOrder      model.SortOrder
		nextPageToken  string
		expectedStatus int
		expectedSize   int32
		expectedItems  int
		checkSorting   bool
		expectedLabels int
	}{
		{
			name:           "Empty catalog list",
			catalogs:       map[string]model.CatalogSource{},
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
			catalogs: map[string]model.CatalogSource{
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1", Enabled: &trueValue},
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
			catalogs: map[string]model.CatalogSource{
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1", Enabled: &trueValue},
				"catalog2": model.CatalogSource{Id: "catalog2", Name: "Test Catalog 2", Enabled: &trueValue},
				"catalog3": model.CatalogSource{Id: "catalog3", Name: "Another Catalog", Enabled: &trueValue},
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
			catalogs: map[string]model.CatalogSource{
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1", Enabled: &trueValue},
				"catalog2": model.CatalogSource{Id: "catalog2", Name: "Test Catalog 2", Enabled: &trueValue},
				"catalog3": model.CatalogSource{Id: "catalog3", Name: "Another Catalog", Enabled: &trueValue},
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
			catalogs: map[string]model.CatalogSource{
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1", Enabled: &trueValue},
				"catalog2": model.CatalogSource{Id: "catalog2", Name: "Test Catalog 2", Enabled: &trueValue},
				"catalog3": model.CatalogSource{Id: "catalog3", Name: "Another Catalog", Enabled: &trueValue},
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
			catalogs: map[string]model.CatalogSource{
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1", Enabled: &trueValue},
				"catalog2": model.CatalogSource{Id: "catalog2", Name: "Test Catalog 2", Enabled: &trueValue},
				"catalog3": model.CatalogSource{Id: "catalog3", Name: "Another Catalog", Enabled: &trueValue},
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
			catalogs: map[string]model.CatalogSource{
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1", Enabled: &trueValue},
				"catalog2": model.CatalogSource{Id: "catalog2", Name: "Test Catalog 2", Enabled: &trueValue},
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
			catalogs: map[string]model.CatalogSource{
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1", Enabled: &trueValue},
			},
			nameFilter:     "",
			pageSize:       "invalid",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Sort by ID ascending",
			catalogs: map[string]model.CatalogSource{
				"catalog2": model.CatalogSource{Id: "catalog2", Name: "B Catalog", Enabled: &trueValue},
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "A Catalog", Enabled: &trueValue},
				"catalog3": model.CatalogSource{Id: "catalog3", Name: "C Catalog", Enabled: &trueValue},
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
			catalogs: map[string]model.CatalogSource{
				"catalog2": model.CatalogSource{Id: "catalog2", Name: "B Catalog", Enabled: &trueValue},
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "A Catalog", Enabled: &trueValue},
				"catalog3": model.CatalogSource{Id: "catalog3", Name: "C Catalog", Enabled: &trueValue},
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
			catalogs: map[string]model.CatalogSource{
				"catalog2": model.CatalogSource{Id: "catalog2", Name: "B Catalog", Enabled: &trueValue},
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "A Catalog", Enabled: &trueValue},
				"catalog3": model.CatalogSource{Id: "catalog3", Name: "C Catalog", Enabled: &trueValue},
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
			catalogs: map[string]model.CatalogSource{
				"catalog2": model.CatalogSource{Id: "catalog2", Name: "B Catalog", Enabled: &trueValue},
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "A Catalog", Enabled: &trueValue},
				"catalog3": model.CatalogSource{Id: "catalog3", Name: "C Catalog"},
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
			catalogs: map[string]model.CatalogSource{
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1"},
			},
			nameFilter:     "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      "INVALID",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid order by field",
			catalogs: map[string]model.CatalogSource{
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1"},
			},
			nameFilter:     "",
			pageSize:       "10",
			orderBy:        "INVALID",
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Default sort order (ID ascending)",
			catalogs: map[string]model.CatalogSource{
				"catalog2": model.CatalogSource{Id: "catalog2", Name: "B Catalog"},
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "A Catalog"},
				"catalog3": model.CatalogSource{Id: "catalog3", Name: "C Catalog"},
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
		{
			name: "Labels should be returned if set",
			catalogs: map[string]model.CatalogSource{
				"catalog1": model.CatalogSource{Id: "catalog1", Name: "Test Catalog 1", Labels: []string{"label1", "label2"}},
				"catalog2": model.CatalogSource{Id: "catalog2", Name: "Test Catalog 2", Labels: []string{"label3", "label4"}},
				"catalog3": model.CatalogSource{Id: "catalog3", Name: "Test Catalog 3", Labels: []string{"label5", "label6"}},
			},
			nameFilter:     "",
			pageSize:       "10",
			orderBy:        model.ORDERBYFIELD_ID,
			sortOrder:      model.SORTORDER_ASC,
			expectedStatus: http.StatusOK,
			expectedSize:   3,
			expectedItems:  3,
			checkSorting:   true,
			expectedLabels: 6,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create service with test catalogs
			sources := catalog.NewSourceCollection()
			sources.Merge("", tc.catalogs)
			service := NewModelCatalogServiceAPIService(&mockModelProvider{}, sources)

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

			labels := make([]string, 0)
			for _, item := range sourceList.Items {
				labels = append(labels, item.Labels...)
			}
			assert.Equal(t, tc.expectedLabels, len(labels))
		})
	}
}

// Define a mock model provider
type mockModelProvider struct {
	models    map[string]*model.CatalogModel
	artifacts map[string][]model.CatalogArtifact
}

// Implement GetModel method for the mock provider
func (m *mockModelProvider) GetModel(ctx context.Context, name string, sourceID string) (*model.CatalogModel, error) {
	model, exists := m.models[name]
	if !exists {
		return nil, nil
	}
	return model, nil
}

func (m *mockModelProvider) ListModels(ctx context.Context, params catalog.ListModelsParams) (model.CatalogModelList, error) {
	var filteredModels []*model.CatalogModel
	for _, mdl := range m.models {
		if params.Query == "" || strings.Contains(strings.ToLower(mdl.Name), strings.ToLower(params.Query)) {
			filteredModels = append(filteredModels, mdl)
		}
	}

	// Sort the filtered models
	sort.SliceStable(filteredModels, func(i, j int) bool {
		cmp := 0
		switch params.OrderBy {
		case model.ORDERBYFIELD_CREATE_TIME:
			// Parse CreateTimeSinceEpoch strings to int64 for comparison
			t1, _ := strconv.ParseInt(pointerOrDefault(filteredModels[i].CreateTimeSinceEpoch, "0"), 10, 64)
			t2, _ := strconv.ParseInt(pointerOrDefault(filteredModels[j].CreateTimeSinceEpoch, "0"), 10, 64)
			cmp = int(t1 - t2)
		case model.ORDERBYFIELD_LAST_UPDATE_TIME:
			// Parse LastUpdateTimeSinceEpoch strings to int64 for comparison
			t1, _ := strconv.ParseInt(pointerOrDefault(filteredModels[i].LastUpdateTimeSinceEpoch, "0"), 10, 64)
			t2, _ := strconv.ParseInt(pointerOrDefault(filteredModels[j].LastUpdateTimeSinceEpoch, "0"), 10, 64)
			cmp = int(t1 - t2)
		case model.ORDERBYFIELD_NAME:
			fallthrough
		default:
			cmp = strings.Compare(filteredModels[i].Name, filteredModels[j].Name)
		}

		if params.SortOrder == model.SORTORDER_DESC {
			return cmp > 0
		}
		return cmp < 0
	})

	totalSize := int32(len(filteredModels))
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// Apply pagination - limit items to page size
	endIndex := int(pageSize)
	if endIndex > len(filteredModels) {
		endIndex = len(filteredModels)
	}

	pagedModels := filteredModels[:endIndex]
	items := make([]model.CatalogModel, len(pagedModels))
	for i, mdl := range pagedModels {
		items[i] = *mdl
	}

	nextPageToken := ""
	if len(filteredModels) > int(pageSize) {
		lastItem := pagedModels[len(pagedModels)-1]
		nextPageToken = (&stringCursor{Value: lastItem.Name, ID: lastItem.Name}).String()
	}

	return model.CatalogModelList{
		Items:         items,
		Size:          totalSize,
		PageSize:      pageSize,
		NextPageToken: nextPageToken,
	}, nil
}

func (m *mockModelProvider) GetArtifacts(ctx context.Context, name string, sourceID string, params catalog.ListArtifactsParams) (model.CatalogArtifactList, error) {
	artifacts, exists := m.artifacts[name]
	if !exists {
		return model.CatalogArtifactList{
			Items:         []model.CatalogArtifact{},
			Size:          0,
			PageSize:      0, // Or a default page size if applicable
			NextPageToken: "",
		}, nil
	}
	return model.CatalogArtifactList{
		Items:         artifacts,
		Size:          int32(len(artifacts)),
		PageSize:      int32(len(artifacts)),
		NextPageToken: "",
	}, nil
}

func TestGetModel(t *testing.T) {
	testCases := []struct {
		name           string
		sources        map[string]model.CatalogSource
		sourceID       string
		modelName      string
		expectedStatus int
		expectedModel  *model.CatalogModel
		provider       catalog.APIProvider
	}{
		{
			name: "Existing model in source",
			sources: map[string]model.CatalogSource{
				"source1": model.CatalogSource{Id: "source1", Name: "Test Source"},
			},
			provider: &mockModelProvider{
				models: map[string]*model.CatalogModel{
					"test-model": {
						Name: "test-model",
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
			sources: map[string]model.CatalogSource{
				"source1": model.CatalogSource{Id: "source1", Name: "Test Source"},
			},
			provider: &mockModelProvider{
				models: map[string]*model.CatalogModel{},
			},
			sourceID:       "source2",
			modelName:      "test-model",
			expectedStatus: http.StatusNotFound,
			expectedModel:  nil,
		},
		{
			name: "Existing source, non-existing model",
			sources: map[string]model.CatalogSource{
				"source1": model.CatalogSource{Id: "source1", Name: "Test Source"},
			},
			provider: &mockModelProvider{
				models: map[string]*model.CatalogModel{},
			},
			sourceID:       "source1",
			modelName:      "test-model",
			expectedStatus: http.StatusNotFound,
			expectedModel:  nil,
		},
		{
			name: "Model name with an escaped slash and version",
			sources: map[string]model.CatalogSource{
				"source1": model.CatalogSource{Id: "source1", Name: "Test Source"},
			},
			provider: &mockModelProvider{
				models: map[string]*model.CatalogModel{
					"some/model:v1.0.0": {
						Name: "some/model:v1.0.0",
					},
				},
			},
			sourceID:       "source1",
			modelName:      "some%2Fmodel%3Av1.0.0",
			expectedStatus: http.StatusOK,
			expectedModel: &model.CatalogModel{
				Name: "some/model:v1.0.0",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create service with test sources
			sources := catalog.NewSourceCollection()
			sources.Merge("", tc.sources)
			service := NewModelCatalogServiceAPIService(tc.provider, sources)

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
		sources           map[string]model.CatalogSource
		sourceID          string
		modelName         string
		expectedStatus    int
		expectedArtifacts []model.CatalogArtifact
		provider          catalog.APIProvider
	}{
		{
			name: "Existing artifacts for model in source",
			sources: map[string]model.CatalogSource{
				"source1": model.CatalogSource{Id: "source1", Name: "Test Source"},
			},
			provider: &mockModelProvider{
				artifacts: map[string][]model.CatalogArtifact{
					"test-model": {
						{
							CatalogModelArtifact: &model.CatalogModelArtifact{
								Uri: "s3://bucket/artifact1",
							},
						},
						{
							CatalogModelArtifact: &model.CatalogModelArtifact{
								Uri: "s3://bucket/artifact2",
							},
						},
					},
				},
			},
			sourceID:       "source1",
			modelName:      "test-model",
			expectedStatus: http.StatusOK,
			expectedArtifacts: []model.CatalogArtifact{
				{
					CatalogModelArtifact: &model.CatalogModelArtifact{
						Uri: "s3://bucket/artifact1",
					},
				},
				{
					CatalogModelArtifact: &model.CatalogModelArtifact{
						Uri: "s3://bucket/artifact2",
					},
				},
			},
		},
		{
			name: "Non-existing source",
			sources: map[string]model.CatalogSource{
				"source1": model.CatalogSource{Id: "source1", Name: "Test Source"},
			},
			provider: &mockModelProvider{
				artifacts: map[string][]model.CatalogArtifact{},
			},
			sourceID:          "source2",
			modelName:         "test-model",
			expectedStatus:    http.StatusOK, // Changed from http.StatusNotFound to http.StatusOK -- having the same behavior as the model registry
			expectedArtifacts: []model.CatalogArtifact{},
		},
		{
			name: "Existing source, no artifacts for model",
			sources: map[string]model.CatalogSource{
				"source1": model.CatalogSource{Id: "source1", Name: "Test Source"},
			},
			provider: &mockModelProvider{
				artifacts: map[string][]model.CatalogArtifact{},
			},
			sourceID:          "source1",
			modelName:         "test-model",
			expectedStatus:    http.StatusOK,
			expectedArtifacts: []model.CatalogArtifact{}, // Should be an empty slice, not nil
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create service with test sources
			sources := catalog.NewSourceCollection()
			sources.Merge("", tc.sources)
			service := NewModelCatalogServiceAPIService(tc.provider, sources)

			// Call GetAllModelArtifacts
			resp, _ := service.GetAllModelArtifacts(
				context.Background(),
				tc.sourceID,
				tc.modelName,
				"",
				"10",
				model.ORDERBYFIELD_CREATE_TIME,
				model.SORTORDER_ASC,
				"",
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
			artifactList, ok := resp.Body.(model.CatalogArtifactList)
			require.True(t, ok, "Response body should be a CatalogArtifactList")

			// Check the artifacts
			assert.Equal(t, tc.expectedArtifacts, artifactList.Items)
			assert.Equal(t, int32(len(tc.expectedArtifacts)), artifactList.Size)
		})
	}
}
