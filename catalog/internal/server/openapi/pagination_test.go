package openapi

import (
	"strconv"
	"testing"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/stretchr/testify/assert"
)

func createCatalogSource(id int) model.CatalogSource {
	return model.CatalogSource{
		Id:   "source" + strconv.Itoa(id),
		Name: "Source " + strconv.Itoa(id),
	}
}

func createCatalogSources(count int) []model.CatalogSource {
	sources := make([]model.CatalogSource, count)
	for i := 0; i < count; i++ {
		sources[i] = createCatalogSource(i)
	}
	return sources
}

func TestPaginateSources(t *testing.T) {
	allSources := createCatalogSources(25)

	testCases := []struct {
		name               string
		items              []model.CatalogSource
		pageSize           string
		orderBy            model.OrderByField
		nextPageToken      string
		expectedItemsCount int
		expectedNextToken  bool
		expectedFirstID    string
		expectedLastID     string
	}{
		{
			name:               "First page, full page",
			items:              allSources,
			pageSize:           "10",
			orderBy:            "ID",
			nextPageToken:      "",
			expectedItemsCount: 10,
			expectedNextToken:  true,
			expectedFirstID:    "source0",
			expectedLastID:     "source9",
		},
		{
			name:               "Second page, full page",
			items:              allSources,
			pageSize:           "10",
			orderBy:            "ID",
			nextPageToken:      (&stringCursor{Value: "source9", ID: "source9"}).String(),
			expectedItemsCount: 10,
			expectedNextToken:  true,
			expectedFirstID:    "source10",
			expectedLastID:     "source19",
		},
		{
			name:               "Last page, partial page",
			items:              allSources,
			pageSize:           "10",
			orderBy:            "ID",
			nextPageToken:      (&stringCursor{Value: "source19", ID: "source19"}).String(),
			expectedItemsCount: 5,
			expectedNextToken:  false,
			expectedFirstID:    "source20",
			expectedLastID:     "source24",
		},
		{
			name:               "Page size larger than items",
			items:              allSources,
			pageSize:           "30",
			orderBy:            "ID",
			nextPageToken:      "",
			expectedItemsCount: 25,
			expectedNextToken:  false,
			expectedFirstID:    "source0",
			expectedLastID:     "source24",
		},
		{
			name:               "Empty items",
			items:              []model.CatalogSource{},
			pageSize:           "10",
			orderBy:            "ID",
			nextPageToken:      "",
			expectedItemsCount: 0,
			expectedNextToken:  false,
		},
		{
			name:               "Order by Name, first page",
			items:              allSources,
			pageSize:           "5",
			orderBy:            "NAME",
			nextPageToken:      "",
			expectedItemsCount: 5,
			expectedNextToken:  true,
			expectedFirstID:    "source0",
			expectedLastID:     "source4",
		},
		{
			name:               "Order by Name, second page",
			items:              allSources,
			pageSize:           "5",
			orderBy:            "NAME",
			nextPageToken:      (&stringCursor{Value: "Source 4", ID: "source4"}).String(),
			expectedItemsCount: 5,
			expectedNextToken:  true,
			expectedFirstID:    "source5",
			expectedLastID:     "source9",
		},
		{
			name:               "Invalid token",
			items:              allSources,
			pageSize:           "10",
			orderBy:            "ID",
			nextPageToken:      "invalid-token",
			expectedItemsCount: 10,
			expectedNextToken:  true,
			expectedFirstID:    "source0",
			expectedLastID:     "source9",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paginator, err := newPaginator[model.CatalogSource](tc.pageSize, tc.orderBy, "", tc.nextPageToken)
			if !assert.NoError(t, err) {
				return
			}

			pagedItems, newNextPageToken := paginator.Paginate(tc.items)

			assert.Equal(t, tc.expectedItemsCount, len(pagedItems))
			if tc.expectedNextToken {
				assert.NotEmpty(t, newNextPageToken)
			} else {
				assert.Empty(t, newNextPageToken)
			}

			if tc.expectedItemsCount > 0 {
				assert.Equal(t, tc.expectedFirstID, pagedItems[0].Id)
				assert.Equal(t, tc.expectedLastID, pagedItems[len(pagedItems)-1].Id)
			}
		})
	}
}

func TestNewPaginator_InvalidPageSize(t *testing.T) {
	testCases := []struct {
		name        string
		pageSize    string
		expectError bool
		errContains string
	}{
		{
			name:        "pageSize=0 returns error",
			pageSize:    "0",
			expectError: true,
			errContains: "pageSize must be at least 1",
		},
		{
			name:        "negative pageSize returns error",
			pageSize:    "-5",
			expectError: true,
			errContains: "pageSize must be at least 1",
		},
		{
			name:        "valid pageSize=1 works",
			pageSize:    "1",
			expectError: false,
		},
		{
			name:        "empty pageSize uses default",
			pageSize:    "",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paginator, err := newPaginator[model.CatalogSource](tc.pageSize, "ID", "", "")
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)
				assert.Nil(t, paginator)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, paginator)
			}
		})
	}
}

func TestPaginateSources_NoDuplicates(t *testing.T) {
	allSources := createCatalogSources(100)
	pageSize := "10"
	orderBy := "ID"

	seenItems := make(map[string]struct{}, len(allSources))
	totalSeen := 0

	paginator, err := newPaginator[model.CatalogSource](pageSize, model.OrderByField(orderBy), "", "")
	if !assert.NoError(t, err) {
		return
	}

	for paginator != nil {
		var pagedItems []model.CatalogSource
		pagedItems, paginator = paginator.Paginate(allSources)

		for _, item := range pagedItems {
			if _, ok := seenItems[item.Id]; ok {
				t.Errorf("Duplicate item found: %s", item.Id)
			}
			seenItems[item.Id] = struct{}{}
		}

		totalSeen += len(pagedItems)
	}

	assert.Equal(t, len(allSources), totalSeen, "Total number of items seen should match the original slice")
}
