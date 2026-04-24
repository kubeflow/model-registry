package modelcatalog

import (
	"testing"

	"github.com/kubeflow/hub/catalog/internal/catalog/basecatalog"
	apimodels "github.com/kubeflow/hub/catalog/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModelFilterFromSourceMergesLegacy(t *testing.T) {
	source := &basecatalog.ModelSource{
		CatalogSource: apimodels.CatalogSource{
			Id:             "test",
			Name:           "Test source",
			Labels:         []string{},
			IncludedModels: []string{"Granite/*"},
		},
	}

	filter, err := NewModelFilterFromSource(source, nil, []string{"Legacy/*"})
	require.NoError(t, err)

	assert.True(t, filter.Allows("Granite/model"))
	assert.False(t, filter.Allows("Legacy/model"))
}
