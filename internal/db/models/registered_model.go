package models

import "github.com/kubeflow/model-registry/internal/db/filter"

type RegisteredModelListOptions struct {
	Pagination
	Name       *string
	ExternalID *string
}

// GetRestEntityType implements the FilterApplier interface
func (r *RegisteredModelListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityRegisteredModel
}

type Properties struct {
	Name             string
	IsCustomProperty bool
	IntValue         *int32
	DoubleValue      *float64
	StringValue      *string
	BoolValue        *bool
	ByteValue        *[]byte
	ProtoValue       *[]byte
}

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
