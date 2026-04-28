package models

import (
	catalogmodels "github.com/kubeflow/hub/catalog/internal/db/models"
	dbmodels "github.com/kubeflow/hub/internal/platform/db/entity"
	"github.com/kubeflow/hub/internal/platform/db/filter"
)

const CatalogModelArtifactType = "model-artifact"

type CatalogModelArtifactListOptions struct {
	dbmodels.Pagination
	Name             *string
	ExternalID       *string
	ParentResourceID *int32
}

// GetRestEntityType implements the FilterApplier interface
func (c *CatalogModelArtifactListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityType(catalogmodels.RestEntityCatalogArtifact)
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
	dbmodels.Entity[CatalogModelArtifactAttributes]
}

type CatalogModelArtifactImpl = dbmodels.BaseEntity[CatalogModelArtifactAttributes]

type CatalogModelArtifactRepository interface {
	GetByID(id int32) (CatalogModelArtifact, error)
	List(listOptions CatalogModelArtifactListOptions) (*dbmodels.ListWrapper[CatalogModelArtifact], error)
	Save(modelArtifact CatalogModelArtifact, parentResourceID *int32) (CatalogModelArtifact, error)
}
