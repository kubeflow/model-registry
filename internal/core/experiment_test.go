package core_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpsertExperiment(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful create", func(t *testing.T) {
		experiment := &openapi.Experiment{
			Name:        "test-experiment",
			Description: apiutils.Of("Test experiment description"),
			Owner:       apiutils.Of("test-owner"),
			ExternalId:  apiutils.Of("exp-ext-123"),
			CustomProperties: &map[string]openapi.MetadataValue{
				"project": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "ml-project",
						MetadataType: "MetadataStringValue",
					},
				},
			},
		}

		result, err := service.UpsertExperiment(experiment)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Id)
		assert.Equal(t, "test-experiment", result.Name)
		assert.Equal(t, "exp-ext-123", *result.ExternalId)
		assert.Equal(t, "Test experiment description", *result.Description)
		assert.Equal(t, "test-owner", *result.Owner)
		assert.NotNil(t, result.CreateTimeSinceEpoch)
		assert.NotNil(t, result.LastUpdateTimeSinceEpoch)
		assert.NotNil(t, result.CustomProperties)
		assert.Contains(t, *result.CustomProperties, "project")
	})

	t.Run("successful update", func(t *testing.T) {
		// Create initial experiment
		experiment := &openapi.Experiment{
			Name:        "update-test-experiment",
			Description: apiutils.Of("Original description"),
			Owner:       apiutils.Of("original-owner"),
		}

		created, err := service.UpsertExperiment(experiment)
		require.NoError(t, err)

		// Update the experiment
		created.Description = apiutils.Of("Updated description")
		created.Owner = apiutils.Of("updated-owner")

		updated, err := service.UpsertExperiment(created)

		require.NoError(t, err)
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, "updated-owner", *updated.Owner)
	})

	t.Run("error on nil experiment", func(t *testing.T) {
		_, err := service.UpsertExperiment(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid experiment pointer")
	})

	t.Run("error on duplicate name", func(t *testing.T) {
		experiment := &openapi.Experiment{
			Name: "duplicate-name-test",
		}

		_, err := service.UpsertExperiment(experiment)
		require.NoError(t, err)

		// Try to create another experiment with the same name
		duplicate := &openapi.Experiment{
			Name: "duplicate-name-test",
		}

		_, err = service.UpsertExperiment(duplicate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestGetExperimentById(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get", func(t *testing.T) {
		// Create an experiment
		experiment := &openapi.Experiment{
			Name:        "get-test-experiment",
			Description: apiutils.Of("Get test description"),
			Owner:       apiutils.Of("get-test-owner"),
		}

		created, err := service.UpsertExperiment(experiment)
		require.NoError(t, err)

		// Get the experiment by ID
		retrieved, err := service.GetExperimentById(*created.Id)

		require.NoError(t, err)
		assert.Equal(t, *created.Id, *retrieved.Id)
		assert.Equal(t, "get-test-experiment", retrieved.Name)
		assert.Equal(t, "Get test description", *retrieved.Description)
		assert.Equal(t, "get-test-owner", *retrieved.Owner)
	})

	t.Run("error on non-existent id", func(t *testing.T) {
		_, err := service.GetExperimentById("999999")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("error on invalid id", func(t *testing.T) {
		_, err := service.GetExperimentById("invalid-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
	})
}

func TestGetExperimentByParams(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create test experiments
	experiment1 := &openapi.Experiment{
		Name:       "params-test-experiment-1",
		ExternalId: apiutils.Of("params-exp-ext-1"),
	}
	created1, err := service.UpsertExperiment(experiment1)
	require.NoError(t, err)

	experiment2 := &openapi.Experiment{
		Name:       "params-test-experiment-2",
		ExternalId: apiutils.Of("params-exp-ext-2"),
	}
	created2, err := service.UpsertExperiment(experiment2)
	require.NoError(t, err)

	t.Run("get by name", func(t *testing.T) {
		result, err := service.GetExperimentByParams(apiutils.Of("params-test-experiment-1"), nil)
		require.NoError(t, err)
		assert.Equal(t, *created1.Id, *result.Id)
		assert.Equal(t, "params-test-experiment-1", result.Name)
	})

	t.Run("get by external id", func(t *testing.T) {
		result, err := service.GetExperimentByParams(nil, apiutils.Of("params-exp-ext-2"))
		require.NoError(t, err)
		assert.Equal(t, *created2.Id, *result.Id)
		assert.Equal(t, "params-exp-ext-2", *result.ExternalId)
	})

	t.Run("error on no params", func(t *testing.T) {
		_, err := service.GetExperimentByParams(nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "supply either name or externalId")
	})

	t.Run("error on not found", func(t *testing.T) {
		_, err := service.GetExperimentByParams(apiutils.Of("non-existent-experiment"), nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestGetExperiments(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create multiple test experiments
	for i := 0; i < 5; i++ {
		experiment := &openapi.Experiment{
			Name:        fmt.Sprintf("list-test-experiment-%d", i),
			Description: apiutils.Of(fmt.Sprintf("List test description %d", i)),
			Owner:       apiutils.Of("list-test-owner"),
		}
		_, err := service.UpsertExperiment(experiment)
		require.NoError(t, err)
	}

	t.Run("get all experiments", func(t *testing.T) {
		result, err := service.GetExperiments(api.ListOptions{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, int(result.Size), 5)
		assert.GreaterOrEqual(t, len(result.Items), 5)
	})

	t.Run("get experiments with pagination", func(t *testing.T) {
		pageSize := int32(3)
		result, err := service.GetExperiments(api.ListOptions{
			PageSize: &pageSize,
		})
		require.NoError(t, err)
		assert.Equal(t, pageSize, result.PageSize)
		assert.Equal(t, pageSize, result.Size)
		assert.Equal(t, int(pageSize), len(result.Items))
	})

	t.Run("get experiments with ordering", func(t *testing.T) {
		orderBy := "CREATE_TIME"
		sortOrder := "DESC"
		result, err := service.GetExperiments(api.ListOptions{
			OrderBy:   &orderBy,
			SortOrder: &sortOrder,
		})
		require.NoError(t, err)
		assert.Greater(t, result.Size, int32(0))

		// Verify ordering (newest first) - now with proper constants
		if len(result.Items) > 1 {
			for i := 0; i < len(result.Items)-1; i++ {
				assert.GreaterOrEqual(t, *result.Items[i].CreateTimeSinceEpoch, *result.Items[i+1].CreateTimeSinceEpoch,
					"Items should be in descending order by create time, but %s is not >= %s",
					*result.Items[i].CreateTimeSinceEpoch, *result.Items[i+1].CreateTimeSinceEpoch)
			}
		}
	})
}

func TestExperimentNonEditableFieldsProtection(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("non-editable fields are protected during update", func(t *testing.T) {
		// Create initial experiment
		experiment := &openapi.Experiment{
			Name:        "test-non-editable",
			Description: apiutils.Of("Original description"),
			Owner:       apiutils.Of("original-owner"),
			ExternalId:  apiutils.Of("original-ext-id"),
			CustomProperties: &map[string]openapi.MetadataValue{
				"original_prop": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "original_value",
						MetadataType: "MetadataStringValue",
					},
				},
			},
		}

		created, err := service.UpsertExperiment(experiment)
		require.NoError(t, err)

		// Store original values
		originalId := *created.Id
		originalName := created.Name
		originalCreateTime := *created.CreateTimeSinceEpoch
		originalUpdateTime := *created.LastUpdateTimeSinceEpoch

		// Wait a moment to ensure timestamp difference
		time.Sleep(10 * time.Millisecond)

		// Attempt to update non-editable fields along with editable fields
		updateRequest := &openapi.Experiment{
			Id:                       created.Id,                         // This should be preserved
			Name:                     "HACKED_NAME",                      // This should be ignored (non-editable)
			CreateTimeSinceEpoch:     apiutils.Of("9999999999"),          // This should be ignored (non-editable)
			LastUpdateTimeSinceEpoch: apiutils.Of("8888888888"),          // This should be ignored (non-editable)
			Description:              apiutils.Of("Updated description"), // This should be updated (editable)
			Owner:                    apiutils.Of("updated-owner"),       // This should be updated (editable)
			ExternalId:               apiutils.Of("updated-ext-id"),      // This should be updated (editable)
			CustomProperties: &map[string]openapi.MetadataValue{
				"updated_prop": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "updated_value",
						MetadataType: "MetadataStringValue",
					},
				},
			},
		}

		updated, err := service.UpsertExperiment(updateRequest)
		require.NoError(t, err)

		// Verify non-editable fields are preserved
		assert.Equal(t, originalId, *updated.Id, "ID should not be changeable")
		assert.Equal(t, originalName, updated.Name, "Name should not be changeable")
		assert.Equal(t, originalCreateTime, *updated.CreateTimeSinceEpoch, "CreateTimeSinceEpoch should not be changeable")

		// LastUpdateTimeSinceEpoch should be updated by the system, not by user input
		assert.NotEqual(t, originalUpdateTime, *updated.LastUpdateTimeSinceEpoch, "LastUpdateTimeSinceEpoch should be updated by system")
		assert.NotEqual(t, "8888888888", *updated.LastUpdateTimeSinceEpoch, "LastUpdateTimeSinceEpoch should not use user-provided value")

		// Verify editable fields are updated
		assert.Equal(t, "Updated description", *updated.Description, "Description should be updatable")
		assert.Equal(t, "updated-owner", *updated.Owner, "Owner should be updatable")
		assert.Equal(t, "updated-ext-id", *updated.ExternalId, "ExternalId should be updatable")
		assert.Contains(t, *updated.CustomProperties, "updated_prop", "CustomProperties should be updatable")
	})

	t.Run("partial update preserves existing editable fields", func(t *testing.T) {
		// Create initial experiment with multiple editable fields
		experiment := &openapi.Experiment{
			Name:        "test-partial-update",
			Description: apiutils.Of("Original description"),
			Owner:       apiutils.Of("original-owner"),
			ExternalId:  apiutils.Of("original-ext-id"),
			CustomProperties: &map[string]openapi.MetadataValue{
				"keep_this": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "keep_value",
						MetadataType: "MetadataStringValue",
					},
				},
			},
		}

		created, err := service.UpsertExperiment(experiment)
		require.NoError(t, err)

		// Partial update - only update description, should preserve other editable fields
		partialUpdate := &openapi.Experiment{
			Id:          created.Id,
			Description: apiutils.Of("Updated description only"),
		}

		updated, err := service.UpsertExperiment(partialUpdate)
		require.NoError(t, err)

		// Verify partial update worked correctly
		assert.Equal(t, "Updated description only", *updated.Description, "Description should be updated")
		assert.Equal(t, "original-owner", *updated.Owner, "Owner should be preserved from existing")
		assert.Equal(t, "original-ext-id", *updated.ExternalId, "ExternalId should be preserved from existing")
		assert.Contains(t, *updated.CustomProperties, "keep_this", "CustomProperties should be preserved from existing")
	})
}

func TestGetExperimentsWithFilterQuery(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create test experiments with various properties for filtering
	testExperiments := []struct {
		experiment *openapi.Experiment
	}{
		{
			experiment: &openapi.Experiment{
				Name:        "nlp-experiment-1",
				Description: apiutils.Of("Natural Language Processing experiment"),
				ExternalId:  apiutils.Of("ext-nlp-001"),
				Owner:       apiutils.Of("alice"),
				CustomProperties: &map[string]openapi.MetadataValue{
					"project": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "nlp",
							MetadataType: "MetadataStringValue",
						},
					},
					"team": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "research",
							MetadataType: "MetadataStringValue",
						},
					},
					"budget": {
						MetadataDoubleValue: &openapi.MetadataDoubleValue{
							DoubleValue:  10000.0,
							MetadataType: "MetadataDoubleValue",
						},
					},
					"priority": {
						MetadataIntValue: &openapi.MetadataIntValue{
							IntValue:     "1",
							MetadataType: "MetadataIntValue",
						},
					},
				},
			},
		},
		{
			experiment: &openapi.Experiment{
				Name:        "cv-experiment-2",
				Description: apiutils.Of("Computer Vision experiment with object detection"),
				ExternalId:  apiutils.Of("ext-cv-002"),
				Owner:       apiutils.Of("bob"),
				CustomProperties: &map[string]openapi.MetadataValue{
					"project": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "vision",
							MetadataType: "MetadataStringValue",
						},
					},
					"team": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "engineering",
							MetadataType: "MetadataStringValue",
						},
					},
					"budget": {
						MetadataDoubleValue: &openapi.MetadataDoubleValue{
							DoubleValue:  25000.0,
							MetadataType: "MetadataDoubleValue",
						},
					},
					"priority": {
						MetadataIntValue: &openapi.MetadataIntValue{
							IntValue:     "2",
							MetadataType: "MetadataIntValue",
						},
					},
					"active": {
						MetadataBoolValue: &openapi.MetadataBoolValue{
							BoolValue:    true,
							MetadataType: "MetadataBoolValue",
						},
					},
				},
			},
		},
		{
			experiment: &openapi.Experiment{
				Name:        "nlp-experiment-2",
				Description: apiutils.Of("NLP sentiment analysis experiment"),
				ExternalId:  apiutils.Of("ext-nlp-003"),
				Owner:       apiutils.Of("alice"),
				CustomProperties: &map[string]openapi.MetadataValue{
					"project": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "nlp",
							MetadataType: "MetadataStringValue",
						},
					},
					"team": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "research",
							MetadataType: "MetadataStringValue",
						},
					},
					"budget": {
						MetadataDoubleValue: &openapi.MetadataDoubleValue{
							DoubleValue:  15000.0,
							MetadataType: "MetadataDoubleValue",
						},
					},
					"priority": {
						MetadataIntValue: &openapi.MetadataIntValue{
							IntValue:     "3",
							MetadataType: "MetadataIntValue",
						},
					},
				},
			},
		},
		{
			experiment: &openapi.Experiment{
				Name:        "rl-experiment",
				Description: apiutils.Of("Reinforcement Learning experiment"),
				ExternalId:  apiutils.Of("ext-rl-004"),
				Owner:       apiutils.Of("charlie"),
				CustomProperties: &map[string]openapi.MetadataValue{
					"project": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "reinforcement-learning",
							MetadataType: "MetadataStringValue",
						},
					},
					"team": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "research",
							MetadataType: "MetadataStringValue",
						},
					},
					"budget": {
						MetadataDoubleValue: &openapi.MetadataDoubleValue{
							DoubleValue:  8000.0,
							MetadataType: "MetadataDoubleValue",
						},
					},
				},
			},
		},
	}

	// Create all test experiments
	for _, te := range testExperiments {
		_, err := service.UpsertExperiment(te.experiment)
		require.NoError(t, err)
	}

	testCases := []struct {
		name          string
		filterQuery   string
		expectedCount int
		expectedNames []string
	}{
		{
			name:          "Filter by exact name",
			filterQuery:   "name = 'nlp-experiment-1'",
			expectedCount: 1,
			expectedNames: []string{"nlp-experiment-1"},
		},
		{
			name:          "Filter by name pattern",
			filterQuery:   "name LIKE 'nlp-%'",
			expectedCount: 2,
			expectedNames: []string{"nlp-experiment-1", "nlp-experiment-2"},
		},
		{
			name:          "Filter by description",
			filterQuery:   "description LIKE '%Computer Vision%'",
			expectedCount: 1,
			expectedNames: []string{"cv-experiment-2"},
		},
		{
			name:          "Filter by owner",
			filterQuery:   "owner = 'alice'",
			expectedCount: 2,
			expectedNames: []string{"nlp-experiment-1", "nlp-experiment-2"},
		},
		{
			name:          "Filter by external ID",
			filterQuery:   "externalId = 'ext-cv-002'",
			expectedCount: 1,
			expectedNames: []string{"cv-experiment-2"},
		},
		{
			name:          "Filter by custom property - string",
			filterQuery:   "project = 'nlp'",
			expectedCount: 2,
			expectedNames: []string{"nlp-experiment-1", "nlp-experiment-2"},
		},
		{
			name:          "Filter by custom property - team",
			filterQuery:   "team = 'engineering'",
			expectedCount: 1,
			expectedNames: []string{"cv-experiment-2"},
		},
		{
			name:          "Filter by custom property - numeric comparison",
			filterQuery:   "budget.double_value > 12000",
			expectedCount: 2,
			expectedNames: []string{"cv-experiment-2", "nlp-experiment-2"},
		},
		{
			name:          "Filter by custom property - integer",
			filterQuery:   "priority <= 2",
			expectedCount: 2,
			expectedNames: []string{"nlp-experiment-1", "cv-experiment-2"},
		},
		{
			name:          "Filter by custom property - boolean",
			filterQuery:   "active = true",
			expectedCount: 1,
			expectedNames: []string{"cv-experiment-2"},
		},
		{
			name:          "Complex filter with AND",
			filterQuery:   "project = 'nlp' AND budget > 12000.0",
			expectedCount: 1,
			expectedNames: []string{"nlp-experiment-2"},
		},
		{
			name:          "Complex filter with OR",
			filterQuery:   "owner = 'alice' OR owner = 'charlie'",
			expectedCount: 3,
			expectedNames: []string{"nlp-experiment-1", "nlp-experiment-2", "rl-experiment"},
		},
		{
			name:          "Complex filter with parentheses",
			filterQuery:   "(project = 'nlp' OR project = 'vision') AND budget.double_value < 20000",
			expectedCount: 2,
			expectedNames: []string{"nlp-experiment-1", "nlp-experiment-2"},
		},
		{
			name:          "Case insensitive pattern matching",
			filterQuery:   "name ILIKE '%EXPERIMENT%'",
			expectedCount: 4,
			expectedNames: []string{"nlp-experiment-1", "cv-experiment-2", "nlp-experiment-2", "rl-experiment"},
		},
		{
			name:          "Filter with NOT condition",
			filterQuery:   "team != 'engineering'",
			expectedCount: 3,
			expectedNames: []string{"nlp-experiment-1", "nlp-experiment-2", "rl-experiment"},
		},

		{
			name:          "Filter by name pattern with suffix",
			filterQuery:   "name LIKE '%-2'",
			expectedCount: 2,
			expectedNames: []string{"cv-experiment-2", "nlp-experiment-2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pageSize := int32(10)
			listOptions := api.ListOptions{
				PageSize:    &pageSize,
				FilterQuery: &tc.filterQuery,
			}

			result, err := service.GetExperiments(listOptions)

			require.NoError(t, err)
			require.NotNil(t, result)

			// Extract names from results
			var actualNames []string
			for _, item := range result.Items {
				for _, expectedName := range tc.expectedNames {
					if item.Name == expectedName {
						actualNames = append(actualNames, item.Name)
						break
					}
				}
			}

			assert.Equal(t, tc.expectedCount, len(actualNames),
				"Expected %d experiments for filter '%s', but got %d",
				tc.expectedCount, tc.filterQuery, len(actualNames))

			// Verify the expected experiments are present
			assert.ElementsMatch(t, tc.expectedNames, actualNames,
				"Expected experiments %v for filter '%s', but got %v",
				tc.expectedNames, tc.filterQuery, actualNames)
		})
	}

	// Test error cases
	t.Run("Invalid filter syntax", func(t *testing.T) {
		pageSize := int32(10)
		invalidFilter := "invalid <<< syntax"
		listOptions := api.ListOptions{
			PageSize:    &pageSize,
			FilterQuery: &invalidFilter,
		}

		result, err := service.GetExperiments(listOptions)

		if assert.Error(t, err) {
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "invalid filter query")
		}
	})

	// Test combining filterQuery with pagination
	t.Run("Filter with pagination", func(t *testing.T) {
		pageSize := int32(1)
		filterQuery := "team = 'research'"
		listOptions := api.ListOptions{
			PageSize:    &pageSize,
			FilterQuery: &filterQuery,
		}

		// Get first page
		firstPage, err := service.GetExperiments(listOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(firstPage.Items))
		assert.NotEmpty(t, firstPage.NextPageToken)

		// Get second page
		listOptions.NextPageToken = &firstPage.NextPageToken
		secondPage, err := service.GetExperiments(listOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(secondPage.Items))

		// Ensure different items on each page
		assert.NotEqual(t, firstPage.Items[0].Id, secondPage.Items[0].Id)

		// Get third page
		listOptions.NextPageToken = &secondPage.NextPageToken
		thirdPage, err := service.GetExperiments(listOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(thirdPage.Items))

		// Ensure all pages have different items
		assert.NotEqual(t, firstPage.Items[0].Id, thirdPage.Items[0].Id)
		assert.NotEqual(t, secondPage.Items[0].Id, thirdPage.Items[0].Id)
	})

	// Test empty results
	t.Run("Filter with no matches", func(t *testing.T) {
		pageSize := int32(10)
		filterQuery := "project = 'nonexistent'"
		listOptions := api.ListOptions{
			PageSize:    &pageSize,
			FilterQuery: &filterQuery,
		}

		result, err := service.GetExperiments(listOptions)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, len(result.Items))
		assert.Equal(t, int32(0), result.Size)
	})

	// Test filtering with ordering
	t.Run("Filter with ordering", func(t *testing.T) {
		pageSize := int32(10)
		filterQuery := "project = 'nlp'"
		orderBy := "name"
		sortOrder := "DESC"
		listOptions := api.ListOptions{
			PageSize:    &pageSize,
			FilterQuery: &filterQuery,
			OrderBy:     &orderBy,
			SortOrder:   &sortOrder,
		}

		result, err := service.GetExperiments(listOptions)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items))

		// Verify descending order
		assert.Equal(t, "nlp-experiment-2", result.Items[0].Name)
		assert.Equal(t, "nlp-experiment-1", result.Items[1].Name)
	})
}
