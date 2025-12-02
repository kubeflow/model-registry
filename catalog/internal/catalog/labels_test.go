package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabelCollection_NewLabelCollection(t *testing.T) {
	lc := NewLabelCollection()
	assert.NotNil(t, lc)
	assert.Empty(t, lc.All())
}

func TestLabelCollection_Merge_AddLabels(t *testing.T) {
	lc := NewLabelCollection()

	labels1 := []map[string]any{
		{"name": "labelNameOne", "displayName": "Label Name One"},
		{"name": "labelNameTwo", "displayName": "Label Name Two"},
		{"name": nil, "displayName": "Null Label"},
	}

	err := lc.Merge("source1", labels1)
	assert.NoError(t, err)

	all := lc.All()
	assert.Len(t, all, 3)
	assert.Contains(t, all, labels1[0])
	assert.Contains(t, all, labels1[1])
	assert.Contains(t, all, labels1[2])
}

func TestLabelCollection_Merge_ReplaceLabels(t *testing.T) {
	lc := NewLabelCollection()

	// Add initial labels from source1
	labels1 := []map[string]any{
		{"name": "labelNameOne", "displayName": "Label Name One"},
		{"name": "labelNameTwo", "displayName": "Label Name Two"},
	}
	err := lc.Merge("source1", labels1)
	assert.NoError(t, err)

	// Replace labels from source1
	labels2 := []map[string]any{
		{"name": "labelNameThree", "displayName": "Label Name Three"},
	}
	err = lc.Merge("source1", labels2)
	assert.NoError(t, err)

	all := lc.All()
	assert.Len(t, all, 1, "Should only have the new label from source1")
	assert.Equal(t, "labelNameThree", all[0]["name"])
}

func TestLabelCollection_Merge_MultipleOrigins(t *testing.T) {
	lc := NewLabelCollection()

	labels1 := []map[string]any{
		{"name": "labelNameOne", "displayName": "Label Name One"},
	}
	labels2 := []map[string]any{
		{"name": "labelNameTwo", "displayName": "Label Name Two"},
	}

	err := lc.Merge("source1", labels1)
	assert.NoError(t, err)
	err = lc.Merge("source2", labels2)
	assert.NoError(t, err)

	all := lc.All()
	assert.Len(t, all, 2)
	assert.Contains(t, all, labels1[0])
	assert.Contains(t, all, labels2[0])
}

func TestLabelCollection_Merge_DuplicateNameWithinSameOrigin(t *testing.T) {
	lc := NewLabelCollection()

	// Try to add labels with duplicate names in the same batch
	labels := []map[string]any{
		{"name": "labelNameOne", "displayName": "Label Name One"},
		{"name": "labelNameTwo", "displayName": "Label Name Two"},
		{"name": "labelNameOne", "displayName": "Label Name One"}, // Duplicate!
	}
	err := lc.Merge("source1", labels)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate label name 'labelNameOne' within the same origin")

	// Verify no labels were added (transaction-like behavior)
	all := lc.All()
	assert.Len(t, all, 0)
}

func TestLabelCollection_Merge_DuplicateNullNameWithinSameOrigin(t *testing.T) {
	lc := NewLabelCollection()

	// Try to add labels with duplicate names in the same batch
	labels := []map[string]any{
		{"name": nil, "displayName": "Label Name Null"},
		{"name": "labelNameTwo", "displayName": "Label Name Two"},
		{"name": nil, "displayName": "Label Name Null"}, // Duplicate!
	}
	err := lc.Merge("source1", labels)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate label name '<nil>' within the same origin")

	// Verify no labels were added (transaction-like behavior)
	all := lc.All()
	assert.Len(t, all, 0)
}

func TestLabelCollection_Merge_DuplicateNameFromDifferentOrigins(t *testing.T) {
	lc := NewLabelCollection()

	// Add label from source1
	labels1 := []map[string]any{
		{"name": "labelNameOne", "displayName": "Label Name One from source1"},
	}
	err := lc.Merge("source1", labels1)
	assert.NoError(t, err)

	// Try to add a label with the same name from source2
	labels2 := []map[string]any{
		{"name": "labelNameOne", "displayName": "Label Name One from source2"},
	}
	err = lc.Merge("source2", labels2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "label with name 'labelNameOne' already exists from another origin")

	// Verify only the first label exists
	all := lc.All()
	assert.Len(t, all, 1)
	assert.Equal(t, "Label Name One from source1", all[0]["displayName"])
}

func TestLabelCollection_Merge_SameOriginCanReplaceSameName(t *testing.T) {
	lc := NewLabelCollection()

	// Add label from source1
	labels1 := []map[string]any{
		{"name": "labelNameOne", "displayName": "Label Name One Version 1"},
	}
	err := lc.Merge("source1", labels1)
	assert.NoError(t, err)

	// Replace with same name from same source - should work
	labels2 := []map[string]any{
		{"name": "labelNameOne", "displayName": "Label Name One Version 2"},
	}
	err = lc.Merge("source1", labels2)
	assert.NoError(t, err)

	// Verify the label was updated
	all := lc.All()
	assert.Len(t, all, 1)
	assert.Equal(t, "Label Name One Version 2", all[0]["displayName"])
}

func TestLabelCollection_Merge_SameOriginWithIdenticalLabels(t *testing.T) {
	lc := NewLabelCollection()

	// Add labels from source1
	labels := []map[string]any{
		{"name": "labelNameOne", "displayName": "Label Name One"},
		{"name": "labelNameTwo", "displayName": "Label Name Two"},
	}
	err := lc.Merge("source1", labels)
	assert.NoError(t, err)

	// Verify labels exist
	all := lc.All()
	assert.Len(t, all, 2)

	// Merge again with IDENTICAL labels (no changes)
	// This should work - labels should be replaced with identical copies
	err = lc.Merge("source1", labels)
	assert.NoError(t, err)

	// Verify labels still exist (not removed due to deep equality check)
	all = lc.All()
	assert.Len(t, all, 2, "Labels should still exist after re-merging identical data")
	assert.Contains(t, all, labels[0])
	assert.Contains(t, all, labels[1])
}

func TestLabelCollection_Merge_RollbackOnValidationFailure(t *testing.T) {
	lc := NewLabelCollection()

	// Add initial labels from source1
	labels1 := []map[string]any{
		{"name": "labelNameOne", "displayName": "Label Name One"},
		{"name": "labelNameTwo", "displayName": "Label Name Two"},
	}
	err := lc.Merge("source1", labels1)
	assert.NoError(t, err)

	// Verify initial state
	all := lc.All()
	assert.Len(t, all, 2)

	// Try to update source1 with invalid labels (duplicates within same batch)
	invalidLabels := []map[string]any{
		{"name": "enterprise", "displayName": "Enterprise"},
		{"name": "enterprise", "displayName": "Enterprise Duplicate"}, // Duplicate!
	}
	err = lc.Merge("source1", invalidLabels)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate label name 'enterprise' within the same origin")

	// Verify that the original labels are STILL THERE (rollback behavior)
	all = lc.All()
	assert.Len(t, all, 2, "Original labels should remain after failed validation")
	assert.Contains(t, all, labels1[0])
	assert.Contains(t, all, labels1[1])
	assert.Equal(t, "Label Name One", all[0]["displayName"])
	assert.Equal(t, "Label Name Two", all[1]["displayName"])
}

func TestLabelCollection_Merge_RollbackOnCrossOriginConflict(t *testing.T) {
	lc := NewLabelCollection()

	// Add labels from source1
	labels1 := []map[string]any{
		{"name": "source1-label", "displayName": "Source 1 Label"},
	}
	err := lc.Merge("source1", labels1)
	assert.NoError(t, err)

	// Add labels from source2
	labels2 := []map[string]any{
		{"name": "source2-label", "displayName": "Source 2 Label"},
	}
	err = lc.Merge("source2", labels2)
	assert.NoError(t, err)

	// Verify both sources have their labels
	all := lc.All()
	assert.Len(t, all, 2)

	// Try to update source1 with a label that conflicts with source2
	conflictingLabels := []map[string]any{
		{"name": "source2-label", "displayName": "Trying to steal source2 label"},
	}
	err = lc.Merge("source1", conflictingLabels)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "label with name 'source2-label' already exists from another origin")

	// Verify that source1's original labels are STILL THERE (rollback behavior)
	all = lc.All()
	assert.Len(t, all, 2, "Both original labels should remain after failed validation")

	// Find source1's label
	var source1Label map[string]any
	for _, label := range all {
		if label["name"] == "source1-label" {
			source1Label = label
			break
		}
	}
	assert.NotNil(t, source1Label, "Source1 label should still exist")
	assert.Equal(t, "Source 1 Label", source1Label["displayName"])
}

func TestLabelCollection_Merge_Deduplicate(t *testing.T) {
	lc := NewLabelCollection()

	// Same label from two sources - should fail because of name conflict
	label := map[string]any{"name": "labelNameOne", "displayName": "Label Name One"}

	err := lc.Merge("source1", []map[string]any{label})
	assert.NoError(t, err)

	// This should fail because a label with the same name already exists
	err = lc.Merge("source2", []map[string]any{label})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "label with name 'labelNameOne' already exists from another origin")
}

func TestLabelCollection_Merge_EmptyLabels(t *testing.T) {
	lc := NewLabelCollection()

	labels1 := []map[string]any{
		{"name": "labelNameOne", "displayName": "Label Name One"},
	}
	err := lc.Merge("source1", labels1)
	assert.NoError(t, err)

	// Clear labels from source1 by merging empty slice
	err = lc.Merge("source1", []map[string]any{})
	assert.NoError(t, err)

	all := lc.All()
	assert.Empty(t, all, "Labels from source1 should be removed")
}

func TestLabelCollection_Merge_UpdateOrigin(t *testing.T) {
	lc := NewLabelCollection()

	// Add labels from source1
	labels1 := []map[string]any{
		{"name": "labelNameOne", "displayName": "Label Name One"},
		{"name": "labelNameTwo", "displayName": "Label Name Two"},
	}
	err := lc.Merge("source1", labels1)
	assert.NoError(t, err)

	// Add labels from source2
	labels2 := []map[string]any{
		{"name": "labelNameThree", "displayName": "Label Name Three"},
	}
	err = lc.Merge("source2", labels2)
	assert.NoError(t, err)

	// Update source1 with different labels
	labels3 := []map[string]any{
		{"name": "labelNameFour", "displayName": "Label Name Four"},
	}
	err = lc.Merge("source1", labels3)
	assert.NoError(t, err)

	all := lc.All()
	assert.Len(t, all, 2, "Should have labelNameFour from source1 and labelNameThree from source2")

	// Verify source2 labels are still there
	assert.Contains(t, all, labels2[0])
	// Verify new source1 labels are there
	assert.Contains(t, all, labels3[0])
	// Verify old source1 labels are gone
	assert.NotContains(t, all, labels1[0])
	assert.NotContains(t, all, labels1[1])
}
