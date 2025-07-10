package openapi

import (
	"strconv"
	"testing"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/db/models"
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
		pageSize           int32
		orderBy            string
		nextPageToken      string
		expectedItemsCount int
		expectedNextToken  bool
		expectedFirstID    string
		expectedLastID     string
	}{
		{
			name:               "First page, full page",
			items:              allSources,
			pageSize:           10,
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
			pageSize:           10,
			orderBy:            "ID",
			nextPageToken:      encodeStringCursor(stringCursor{Value: "source9", ID: "source9"}),
			expectedItemsCount: 10,
			expectedNextToken:  true,
			expectedFirstID:    "source10",
			expectedLastID:     "source19",
		},
		{
			name:               "Last page, partial page",
			items:              allSources,
			pageSize:           10,
			orderBy:            "ID",
			nextPageToken:      encodeStringCursor(stringCursor{Value: "source19", ID: "source19"}),
			expectedItemsCount: 5,
			expectedNextToken:  false,
			expectedFirstID:    "source20",
			expectedLastID:     "source24",
		},
		{
			name:               "Page size larger than items",
			items:              allSources,
			pageSize:           30,
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
			pageSize:           10,
			orderBy:            "ID",
			nextPageToken:      "",
			expectedItemsCount: 0,
			expectedNextToken:  false,
		},
		{
			name:               "Order by Name, first page",
			items:              allSources,
			pageSize:           5,
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
			pageSize:           5,
			orderBy:            "NAME",
			nextPageToken:      encodeStringCursor(stringCursor{Value: "Source 4", ID: "source4"}),
			expectedItemsCount: 5,
			expectedNextToken:  true,
			expectedFirstID:    "source5",
			expectedLastID:     "source9",
		},
		{
			name:               "Invalid token",
			items:              allSources,
			pageSize:           10,
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
			pagination := &models.Pagination{
				PageSize:      &tc.pageSize,
				OrderBy:       &tc.orderBy,
				NextPageToken: &tc.nextPageToken,
			}

			pagedItems, newNextPageToken := paginateSources(tc.items, pagination)

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

func TestPaginateSources_NoDuplicates(t *testing.T) {
	allSources := createCatalogSources(100)
	pageSize := int32(10)
	orderBy := "ID"

	seenItems := make(map[string]bool)
	var nextPageToken string
	totalSeen := 0

	for {
		pagination := &models.Pagination{
			PageSize:      &pageSize,
			OrderBy:       &orderBy,
			NextPageToken: &nextPageToken,
		}

		pagedItems, newNextPageToken := paginateSources(allSources, pagination)

		for _, item := range pagedItems {
			if _, ok := seenItems[item.Id]; ok {
				t.Errorf("Duplicate item found: %s", item.Id)
			}
			seenItems[item.Id] = true
		}

		totalSeen += len(pagedItems)

		if newNextPageToken == "" {
			break
		}
		nextPageToken = newNextPageToken
	}

	assert.Equal(t, len(allSources), totalSeen, "Total number of items seen should match the original slice")
}
