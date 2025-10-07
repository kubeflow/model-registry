package models

import (
	catalogfilter "github.com/kubeflow/model-registry/catalog/internal/db/filter"
	"github.com/kubeflow/model-registry/internal/db/filter"
	"github.com/kubeflow/model-registry/internal/db/models"
)

type CatalogModelListOptions struct {
	models.Pagination
	Name       *string
	ExternalID *string
	SourceIDs  *[]string
	Query      *string
}

// GetRestEntityType implements the FilterApplier interface
func (c *CatalogModelListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityType(catalogfilter.RestEntityCatalogModel)
}

type CatalogModelAttributes struct {
	Name                     *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type CatalogModel interface {
	models.Entity[CatalogModelAttributes]
}

type CatalogModelImpl = models.BaseEntity[CatalogModelAttributes]

type CatalogModelRepository interface {
	GetByID(id int32) (CatalogModel, error)
	GetByName(name string) (CatalogModel, error)
	List(listOptions CatalogModelListOptions) (*models.ListWrapper[CatalogModel], error)
	Save(model CatalogModel) (CatalogModel, error)
}
