package models

import (
	"github.com/kubeflow/model-registry/internal/db/filter"
	"github.com/kubeflow/model-registry/internal/db/models"
)

const CatalogModelArtifactType = "model-artifact"

type CatalogModelArtifactListOptions struct {
	models.Pagination
	Name             *string
	ExternalID       *string
	ParentResourceID *int32
}

// GetRestEntityType implements the FilterApplier interface
func (c *CatalogModelArtifactListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityModelArtifact // Reusing existing filter type
}

type CatalogModelArtifactAttributes struct {
	Name                     *string
	URI                      *string
	ArtifactType             *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type CatalogModelArtifact interface {
	models.Entity[CatalogModelArtifactAttributes]
}

type CatalogModelArtifactImpl = models.BaseEntity[CatalogModelArtifactAttributes]

type CatalogModelArtifactRepository interface {
	GetByID(id int32) (CatalogModelArtifact, error)
	List(listOptions CatalogModelArtifactListOptions) (*models.ListWrapper[CatalogModelArtifact], error)
	Save(modelArtifact CatalogModelArtifact, parentResourceID *int32) (CatalogModelArtifact, error)
}
