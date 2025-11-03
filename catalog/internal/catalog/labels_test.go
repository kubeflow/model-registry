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

	labels1 := []map[string]string{
		{"name": "labelNameOne", "displayName": "Label Name One"},
		{"name": "labelNameTwo", "displayName": "Label Name Two"},
	}

	lc.Merge("source1", labels1)

	all := lc.All()
	assert.Len(t, all, 2)
	assert.Contains(t, all, labels1[0])
	assert.Contains(t, all, labels1[1])
}

func TestLabelCollection_Merge_ReplaceLabels(t *testing.T) {
	lc := NewLabelCollection()

	// Add initial labels from source1
	labels1 := []map[string]string{
		{"name": "labelNameOne", "displayName": "Label Name One"},
		{"name": "labelNameTwo", "displayName": "Label Name Two"},
	}
	lc.Merge("source1", labels1)

	// Replace labels from source1
	labels2 := []map[string]string{
		{"name": "labelNameThree", "displayName": "Label Name Three"},
	}
	lc.Merge("source1", labels2)

	all := lc.All()
	assert.Len(t, all, 1, "Should only have the new label from source1")
	assert.Equal(t, "labelNameThree", all[0]["name"])
}

func TestLabelCollection_Merge_MultipleOrigins(t *testing.T) {
	lc := NewLabelCollection()

	labels1 := []map[string]string{
		{"name": "labelNameOne", "displayName": "Label Name One"},
	}
	labels2 := []map[string]string{
		{"name": "labelNameTwo", "displayName": "Label Name Two"},
	}

	lc.Merge("source1", labels1)
	lc.Merge("source2", labels2)

	all := lc.All()
	assert.Len(t, all, 2)
	assert.Contains(t, all, labels1[0])
	assert.Contains(t, all, labels2[0])
}

func TestLabelCollection_Merge_Deduplicate(t *testing.T) {
	lc := NewLabelCollection()

	// Same label from two sources
	label := map[string]string{"name": "labelNameOne", "displayName": "Label Name One"}

	lc.Merge("source1", []map[string]string{label})
	lc.Merge("source2", []map[string]string{label})

	all := lc.All()
	assert.Len(t, all, 1, "Duplicate labels should be deduplicated")
	assert.Equal(t, label, all[0])
}

func TestLabelCollection_Merge_EmptyLabels(t *testing.T) {
	lc := NewLabelCollection()

	labels1 := []map[string]string{
		{"name": "labelNameOne", "displayName": "Label Name One"},
	}
	lc.Merge("source1", labels1)

	// Clear labels from source1 by merging empty slice
	lc.Merge("source1", []map[string]string{})

	all := lc.All()
	assert.Empty(t, all, "Labels from source1 should be removed")
}

func TestLabelCollection_Merge_UpdateOrigin(t *testing.T) {
	lc := NewLabelCollection()

	// Add labels from source1
	labels1 := []map[string]string{
		{"name": "labelNameOne", "displayName": "Label Name One"},
		{"name": "labelNameTwo", "displayName": "Label Name Two"},
	}
	lc.Merge("source1", labels1)

	// Add labels from source2
	labels2 := []map[string]string{
		{"name": "labelNameThree", "displayName": "Label Name Three"},
	}
	lc.Merge("source2", labels2)

	// Update source1 with different labels
	labels3 := []map[string]string{
		{"name": "labelNameFour", "displayName": "Label Name Four"},
	}
	lc.Merge("source1", labels3)

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

func TestMapsEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        map[string]string
		b        map[string]string
		expected bool
	}{
		{
			name:     "equal maps",
			a:        map[string]string{"name": "test", "value": "123"},
			b:        map[string]string{"name": "test", "value": "123"},
			expected: true,
		},
		{
			name:     "different values",
			a:        map[string]string{"name": "test", "value": "123"},
			b:        map[string]string{"name": "test", "value": "456"},
			expected: false,
		},
		{
			name:     "different keys",
			a:        map[string]string{"name": "test"},
			b:        map[string]string{"other": "test"},
			expected: false,
		},
		{
			name:     "different lengths",
			a:        map[string]string{"name": "test"},
			b:        map[string]string{"name": "test", "value": "123"},
			expected: false,
		},
		{
			name:     "both empty",
			a:        map[string]string{},
			b:        map[string]string{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapsEqual(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}
