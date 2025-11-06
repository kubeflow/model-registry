package models

import (
	"github.com/kubeflow/model-registry/catalog/internal/db/filter"
	dbfilter "github.com/kubeflow/model-registry/internal/db/filter"
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
func (c *CatalogArtifactListOptions) GetRestEntityType() dbfilter.RestEntityType {
	return dbfilter.RestEntityType(filter.RestEntityCatalogArtifact)
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
