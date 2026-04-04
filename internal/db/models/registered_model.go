package models

import (
	"github.com/kubeflow/model-registry/internal/db/filter"
	"github.com/kubeflow/model-registry/internal/platform/db/entity"
)

type RegisteredModelListOptions struct {
	Pagination
	Name       *string
	ExternalID *string
}

func (r *RegisteredModelListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityRegisteredModel
}

// Properties is a type alias for entity.Properties
type Properties = entity.Properties

// Re-export constructor functions
var NewStringProperty = entity.NewStringProperty
var NewIntProperty = entity.NewIntProperty
var NewDoubleProperty = entity.NewDoubleProperty
var NewBoolProperty = entity.NewBoolProperty
var NewByteProperty = entity.NewByteProperty

type RegisteredModelAttributes struct {
	Name                     *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type RegisteredModel interface {
	Entity[RegisteredModelAttributes]
}

type RegisteredModelImpl = BaseEntity[RegisteredModelAttributes]

type RegisteredModelRepository interface {
	GetByID(id int32) (RegisteredModel, error)
	List(listOptions RegisteredModelListOptions) (*ListWrapper[RegisteredModel], error)
	Save(model RegisteredModel) (RegisteredModel, error)
}
