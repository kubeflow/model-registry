package models

type ModelVersionListOptions struct {
	Pagination
	Name             *string
	ExternalID       *string
	ParentResourceID *int32
}

type ModelVersionAttributes struct {
	Name                     *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type ModelVersion interface {
	GetID() *int32
	GetTypeID() *int32
	GetAttributes() *ModelVersionAttributes
	GetProperties() *[]Properties
	GetCustomProperties() *[]Properties
}

type ModelVersionImpl struct {
	ID               *int32
	TypeID           *int32
	Attributes       *ModelVersionAttributes
	Properties       *[]Properties
	CustomProperties *[]Properties
}

func (r *ModelVersionImpl) GetID() *int32 {
	return r.ID
}

func (r *ModelVersionImpl) GetTypeID() *int32 {
	return r.TypeID
}

func (r *ModelVersionImpl) GetAttributes() *ModelVersionAttributes {
	return r.Attributes
}

func (r *ModelVersionImpl) GetProperties() *[]Properties {
	return r.Properties
}

func (r *ModelVersionImpl) GetCustomProperties() *[]Properties {
	return r.CustomProperties
}

type ModelVersionRepository interface {
	GetByID(id int32) (ModelVersion, error)
	List(listOptions ModelVersionListOptions) (*ListWrapper[ModelVersion], error)
	Save(model ModelVersion) (ModelVersion, error)
}
