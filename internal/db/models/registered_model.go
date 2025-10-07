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

// Constructor functions for Properties

// NewStringProperty creates a string property
func NewStringProperty(name string, value string, isCustom bool) Properties {
	return Properties{
		Name:             name,
		IsCustomProperty: isCustom,
		StringValue:      &value,
	}
}

// NewIntProperty creates an int property
func NewIntProperty(name string, value int32, isCustom bool) Properties {
	return Properties{
		Name:             name,
		IsCustomProperty: isCustom,
		IntValue:         &value,
	}
}

// NewDoubleProperty creates a double property
func NewDoubleProperty(name string, value float64, isCustom bool) Properties {
	return Properties{
		Name:             name,
		IsCustomProperty: isCustom,
		DoubleValue:      &value,
	}
}

// NewBoolProperty creates a bool property
func NewBoolProperty(name string, value bool, isCustom bool) Properties {
	return Properties{
		Name:             name,
		IsCustomProperty: isCustom,
		BoolValue:        &value,
	}
}

// NewByteProperty creates a byte property
func NewByteProperty(name string, value []byte, isCustom bool) Properties {
	return Properties{
		Name:             name,
		IsCustomProperty: isCustom,
		ByteValue:        &value,
	}
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
