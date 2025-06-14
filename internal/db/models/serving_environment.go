package models

type ServingEnvironmentListOptions struct {
	Pagination
	Name       *string
	ExternalID *string
}

type ServingEnvironmentAttributes struct {
	Name                     *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type ServingEnvironment interface {
	GetID() *int32
	GetTypeID() *int32
	GetAttributes() *ServingEnvironmentAttributes
	GetProperties() *[]Properties
	GetCustomProperties() *[]Properties
}

type ServingEnvironmentImpl struct {
	ID               *int32
	TypeID           *int32
	Attributes       *ServingEnvironmentAttributes
	Properties       *[]Properties
	CustomProperties *[]Properties
}

func (r *ServingEnvironmentImpl) GetID() *int32 {
	return r.ID
}

func (r *ServingEnvironmentImpl) GetTypeID() *int32 {
	return r.TypeID
}

func (r *ServingEnvironmentImpl) GetAttributes() *ServingEnvironmentAttributes {
	return r.Attributes
}

func (r *ServingEnvironmentImpl) GetProperties() *[]Properties {
	return r.Properties
}

func (r *ServingEnvironmentImpl) GetCustomProperties() *[]Properties {
	return r.CustomProperties
}

type ServingEnvironmentRepository interface {
	GetByID(id int32) (ServingEnvironment, error)
	List(listOptions ServingEnvironmentListOptions) (*ListWrapper[ServingEnvironment], error)
	Save(model ServingEnvironment) (ServingEnvironment, error)
}
