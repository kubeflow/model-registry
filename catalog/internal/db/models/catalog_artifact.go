package models

import (
	"github.com/kubeflow/model-registry/internal/db/filter"
	"github.com/kubeflow/model-registry/internal/db/models"
)

type CatalogArtifactListOptions struct {
	models.Pagination
	Name                *string
	ExternalID          *string
	ParentResourceID    *int32
	ArtifactType        *string
	ArtifactTypesFilter []string
}

// GetRestEntityType implements the FilterApplier interface
// This enables advanced filtering support for catalog artifacts
func (c *CatalogArtifactListOptions) GetRestEntityType() filter.RestEntityType {
	// Determine the appropriate REST entity type based on artifact type
	if c.ArtifactType != nil {
		switch *c.ArtifactType {
		case "model-artifact":
			return filter.RestEntityModelArtifact
		case "metrics-artifact":
			return filter.RestEntityModelArtifact // Reusing existing filter type
		}
	}
	// Default to ModelArtifact if no specific type is provided
	return filter.RestEntityModelArtifact
}

// CatalogArtifact is a discriminated union that can hold different catalog artifact types
type CatalogArtifact struct {
	CatalogModelArtifact   CatalogModelArtifact
	CatalogMetricsArtifact CatalogMetricsArtifact
}

type CatalogArtifactRepository interface {
	GetByID(id int32) (CatalogArtifact, error)
	List(listOptions CatalogArtifactListOptions) (*models.ListWrapper[CatalogArtifact], error)
	DeleteByParentID(artifactType string, parentResourceID int32) error
}
