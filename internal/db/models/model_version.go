package models

import "github.com/kubeflow/model-registry/internal/db/filter"

type ModelVersionListOptions struct {
	Pagination
	Name             *string
	ExternalID       *string
	ParentResourceID *int32
}

// GetRestEntityType implements the FilterApplier interface
func (m *ModelVersionListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityModelVersion
}

type ModelVersionAttributes struct {
	Name                     *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type ModelVersion interface {
	Entity[ModelVersionAttributes]
}

type ModelVersionImpl = BaseEntity[ModelVersionAttributes]

type ModelVersionRepository interface {
	GetByID(id int32) (ModelVersion, error)
	List(listOptions ModelVersionListOptions) (*ListWrapper[ModelVersion], error)
	Save(model ModelVersion) (ModelVersion, error)
}
