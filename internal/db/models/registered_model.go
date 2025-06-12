package models

type RegisteredModelListOptions struct {
	Pagination
	Name       *string
	ExternalID *string
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
	GetID() *int32
	GetTypeID() *int32
	GetAttributes() *RegisteredModelAttributes
	GetProperties() *[]Properties
	GetCustomProperties() *[]Properties
}

type RegisteredModelImpl struct {
	ID               *int32
	TypeID           *int32
	Attributes       *RegisteredModelAttributes
	Properties       *[]Properties
	CustomProperties *[]Properties
}

func (r *RegisteredModelImpl) GetID() *int32 {
	return r.ID
}

func (r *RegisteredModelImpl) GetTypeID() *int32 {
	return r.TypeID
}

func (r *RegisteredModelImpl) GetAttributes() *RegisteredModelAttributes {
	return r.Attributes
}

func (r *RegisteredModelImpl) GetProperties() *[]Properties {
	return r.Properties
}

func (r *RegisteredModelImpl) GetCustomProperties() *[]Properties {
	return r.CustomProperties
}

type RegisteredModelRepository interface {
	GetByID(id int32) (RegisteredModel, error)
	List(listOptions RegisteredModelListOptions) (*ListWrapper[RegisteredModel], error)
	Save(model RegisteredModel) (RegisteredModel, error)
}
