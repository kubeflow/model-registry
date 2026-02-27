package models

import (
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/filter"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
)

type CatalogModelListOptions struct {
	dbmodels.Pagination
	Name       *string
	ExternalID *string
	SourceIDs  *[]string
	Query      *string
}

// GetRestEntityType implements the FilterApplier interface
func (c *CatalogModelListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityType(models.RestEntityCatalogModel)
}

type CatalogModelAttributes struct {
	Name                     *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type CatalogModel interface {
	dbmodels.Entity[CatalogModelAttributes]
}

type CatalogModelImpl = dbmodels.BaseEntity[CatalogModelAttributes]

type CatalogModelRepository interface {
	GetByID(id int32) (CatalogModel, error)
	GetByName(name string) (CatalogModel, error)
	List(listOptions CatalogModelListOptions) (*dbmodels.ListWrapper[CatalogModel], error)
	Save(model CatalogModel) (CatalogModel, error)
	DeleteBySource(sourceID string) error
	DeleteByID(id int32) error
	GetDistinctSourceIDs() ([]string, error)
}