package core_test

import (
	"testing"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/core"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPartialUpdateBehavior verifies that partial updates only modify provided fields
func TestPartialUpdateBehavior(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("Experiment partial update preserves existing fields", func(t *testing.T) {
		testExperimentPartialUpdate(t, service)
	})

	t.Run("RegisteredModel partial update preserves existing fields", func(t *testing.T) {
		testRegisteredModelPartialUpdate(t, service)
	})

	t.Run("ModelVersion partial update preserves existing fields", func(t *testing.T) {
		testModelVersionPartialUpdate(t, service)
	})

	t.Run("Artifact partial update preserves existing fields", func(t *testing.T) {
		testArtifactPartialUpdate(t, service)
	})
}

func testExperimentPartialUpdate(t *testing.T, service *core.ModelRegistryService) {
	// Create experiment with all fields populated
	experiment := &openapi.Experiment{
		Name:        "test-partial-experiment",
		Description: apiutils.Of("Original description"),
		Owner:       apiutils.Of("original-owner"),
		ExternalId:  apiutils.Of("original-ext-id"),
		State:       (*openapi.ExperimentState)(apiutils.Of(string(openapi.EXPERIMENTSTATE_LIVE))),
		CustomProperties: &map[string]openapi.MetadataValue{
			"team": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue:  "data-science",
					MetadataType: "MetadataStringValue",
				},
			},
			"priority": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue:  "high",
					MetadataType: "MetadataStringValue",
				},
			},
		},
	}

	created, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	t.Run("update only description - preserves all other fields", func(t *testing.T) {
		// Update only description
		partialUpdate := &openapi.Experiment{
			Id:          created.Id,
			Description: apiutils.Of("Updated description only"),
		}

		updated, err := service.UpsertExperiment(partialUpdate)
		require.NoError(t, err)

		// Verify only description changed
		assert.Equal(t, "Updated description only", *updated.Description)

		// Verify all other fields preserved
		assert.Equal(t, "original-owner", *updated.Owner, "Owner should be preserved")
		assert.Equal(t, "original-ext-id", *updated.ExternalId, "ExternalId should be preserved")
		assert.Equal(t, openapi.EXPERIMENTSTATE_LIVE, *updated.State, "State should be preserved")

		// Verify custom properties preserved
		assert.Contains(t, *updated.CustomProperties, "team", "Custom property 'team' should be preserved")
		assert.Contains(t, *updated.CustomProperties, "priority", "Custom property 'priority' should be preserved")
		assert.Equal(t, "data-science", (*updated.CustomProperties)["team"].MetadataStringValue.StringValue)
		assert.Equal(t, "high", (*updated.CustomProperties)["priority"].MetadataStringValue.StringValue)
	})

	t.Run("update only custom properties - preserves all other fields", func(t *testing.T) {
		// Update only custom properties
		partialUpdate := &openapi.Experiment{
			Id: created.Id,
			CustomProperties: &map[string]openapi.MetadataValue{
				"environment": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "production",
						MetadataType: "MetadataStringValue",
					},
				},
			},
		}

		updated, err := service.UpsertExperiment(partialUpdate)
		require.NoError(t, err)

		// Verify custom properties completely replaced (this is expected behavior for maps)
		assert.Contains(t, *updated.CustomProperties, "environment")
		assert.Equal(t, "production", (*updated.CustomProperties)["environment"].MetadataStringValue.StringValue)

		// For complete replacement of CustomProperties map, previous properties should be gone
		// This is the expected behavior when you provide a CustomProperties field

		// Verify all other fields preserved
		assert.Equal(t, "Updated description only", *updated.Description, "Description should be preserved from previous update")
		assert.Equal(t, "original-owner", *updated.Owner, "Owner should be preserved")
		assert.Equal(t, "original-ext-id", *updated.ExternalId, "ExternalId should be preserved")
	})

	t.Run("nil/empty fields are ignored", func(t *testing.T) {
		// Update with some nil fields - they should be ignored
		partialUpdate := &openapi.Experiment{
			Id:          created.Id,
			Owner:       apiutils.Of("new-owner"),
			Description: nil, // This should be ignored
			ExternalId:  nil, // This should be ignored
		}

		updated, err := service.UpsertExperiment(partialUpdate)
		require.NoError(t, err)

		// Verify only non-nil fields updated
		assert.Equal(t, "new-owner", *updated.Owner, "Owner should be updated")

		// Verify nil fields were ignored (preserved from previous state)
		assert.Equal(t, "Updated description only", *updated.Description, "Description should be preserved (nil ignored)")
		assert.Equal(t, "original-ext-id", *updated.ExternalId, "ExternalId should be preserved (nil ignored)")
	})
}

func testRegisteredModelPartialUpdate(t *testing.T, service *core.ModelRegistryService) {
	// Create registered model with all fields
	model := &openapi.RegisteredModel{
		Name:        "test-partial-model",
		Description: apiutils.Of("Original model description"),
		Owner:       apiutils.Of("original-model-owner"),
		ExternalId:  apiutils.Of("original-model-ext-id"),
		State:       (*openapi.RegisteredModelState)(apiutils.Of(string(openapi.REGISTEREDMODELSTATE_LIVE))),
		CustomProperties: &map[string]openapi.MetadataValue{
			"framework": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue:  "tensorflow",
					MetadataType: "MetadataStringValue",
				},
			},
		},
	}

	created, err := service.UpsertRegisteredModel(model)
	require.NoError(t, err)

	// Test partial update - only description
	partialUpdate := &openapi.RegisteredModel{
		Id:          created.Id,
		Description: apiutils.Of("Updated model description"),
	}

	updated, err := service.UpsertRegisteredModel(partialUpdate)
	require.NoError(t, err)

	// Verify only description changed
	assert.Equal(t, "Updated model description", *updated.Description)

	// Verify all other fields preserved
	assert.Equal(t, "original-model-owner", *updated.Owner, "Owner should be preserved")
	assert.Equal(t, "original-model-ext-id", *updated.ExternalId, "ExternalId should be preserved")
	assert.Equal(t, openapi.REGISTEREDMODELSTATE_LIVE, *updated.State, "State should be preserved")
	assert.Contains(t, *updated.CustomProperties, "framework", "CustomProperties should be preserved")
}

func testModelVersionPartialUpdate(t *testing.T, service *core.ModelRegistryService) {
	// Create parent registered model
	rm := &openapi.RegisteredModel{Name: "test-mv-partial-parent"}
	createdRM, err := service.UpsertRegisteredModel(rm)
	require.NoError(t, err)

	// Create model version with all fields
	mv := &openapi.ModelVersion{
		Name:        "v1.0",
		Description: apiutils.Of("Original version description"),
		Author:      apiutils.Of("original-author"),
		ExternalId:  apiutils.Of("original-mv-ext-id"),
		State:       (*openapi.ModelVersionState)(apiutils.Of(string(openapi.MODELVERSIONSTATE_LIVE))),
		CustomProperties: &map[string]openapi.MetadataValue{
			"accuracy": {
				MetadataDoubleValue: &openapi.MetadataDoubleValue{
					DoubleValue:  0.95,
					MetadataType: "MetadataDoubleValue",
				},
			},
		},
	}

	created, err := service.UpsertModelVersion(mv, createdRM.Id)
	require.NoError(t, err)

	// Test partial update - only author
	partialUpdate := &openapi.ModelVersion{
		Id:     created.Id,
		Author: apiutils.Of("updated-author"),
	}

	updated, err := service.UpsertModelVersion(partialUpdate, nil)
	require.NoError(t, err)

	// Verify only author changed
	assert.Equal(t, "updated-author", *updated.Author)

	// Verify all other fields preserved
	assert.Equal(t, "Original version description", *updated.Description, "Description should be preserved")
	assert.Equal(t, "original-mv-ext-id", *updated.ExternalId, "ExternalId should be preserved")
	assert.Equal(t, openapi.MODELVERSIONSTATE_LIVE, *updated.State, "State should be preserved")
	assert.Contains(t, *updated.CustomProperties, "accuracy", "CustomProperties should be preserved")
}

func testArtifactPartialUpdate(t *testing.T, service *core.ModelRegistryService) {
	// Create model artifact with all fields
	artifact := &openapi.Artifact{
		ModelArtifact: &openapi.ModelArtifact{
			Name:        apiutils.Of("test-partial-artifact"),
			Description: apiutils.Of("Original artifact description"),
			Uri:         apiutils.Of("s3://original/path"),
			ExternalId:  apiutils.Of("original-artifact-ext-id"),
			State:       (*openapi.ArtifactState)(apiutils.Of(string(openapi.ARTIFACTSTATE_LIVE))),
			CustomProperties: &map[string]openapi.MetadataValue{
				"size": {
					MetadataIntValue: &openapi.MetadataIntValue{
						IntValue:     "1024",
						MetadataType: "MetadataIntValue",
					},
				},
			},
		},
	}

	created, err := service.UpsertArtifact(artifact)
	require.NoError(t, err)

	// Test partial update - only URI
	partialUpdate := &openapi.Artifact{
		ModelArtifact: &openapi.ModelArtifact{
			Id:  created.ModelArtifact.Id,
			Uri: apiutils.Of("s3://updated/path"),
		},
	}

	updated, err := service.UpsertArtifact(partialUpdate)
	require.NoError(t, err)

	// Verify only URI changed
	assert.Equal(t, "s3://updated/path", *updated.ModelArtifact.Uri)

	// Verify all other fields preserved
	assert.Equal(t, "Original artifact description", *updated.ModelArtifact.Description, "Description should be preserved")
	assert.Equal(t, "original-artifact-ext-id", *updated.ModelArtifact.ExternalId, "ExternalId should be preserved")
	assert.Equal(t, openapi.ARTIFACTSTATE_LIVE, *updated.ModelArtifact.State, "State should be preserved")
	assert.Contains(t, *updated.ModelArtifact.CustomProperties, "size", "CustomProperties should be preserved")
	assert.Equal(t, "1024", (*updated.ModelArtifact.CustomProperties)["size"].MetadataIntValue.IntValue)
}

// TestEmptyStringVsNilBehavior verifies the difference between empty strings and nil pointers
func TestEmptyStringVsNilBehavior(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment with description
	experiment := &openapi.Experiment{
		Name:        "test-empty-vs-nil",
		Description: apiutils.Of("Original description"),
		Owner:       apiutils.Of("original-owner"),
	}

	created, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	t.Run("nil description field is ignored", func(t *testing.T) {
		// Update with nil description - should be ignored
		updateWithNil := &openapi.Experiment{
			Id:          created.Id,
			Description: nil, // This should be ignored
			Owner:       apiutils.Of("updated-owner"),
		}

		updated, err := service.UpsertExperiment(updateWithNil)
		require.NoError(t, err)

		assert.Equal(t, "Original description", *updated.Description, "Nil description should be ignored")
		assert.Equal(t, "updated-owner", *updated.Owner, "Owner should be updated")
	})

	t.Run("empty string description field updates to empty", func(t *testing.T) {
		// Update with empty string - should actually update to empty
		updateWithEmpty := &openapi.Experiment{
			Id:          created.Id,
			Description: apiutils.Of(""), // This should update to empty string
		}

		updated, err := service.UpsertExperiment(updateWithEmpty)
		require.NoError(t, err)

		assert.Equal(t, "", *updated.Description, "Empty string should update the field")
		assert.Equal(t, "updated-owner", *updated.Owner, "Owner should be preserved")
	})
}
