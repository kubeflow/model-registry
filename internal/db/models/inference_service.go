package models

type InferenceServiceListOptions struct {
	Pagination
	Name             *string
	ExternalID       *string
	ParentResourceID *int32
	Runtime          *string
}

type InferenceServiceAttributes struct {
	Name                     *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type InferenceService interface {
	GetID() *int32
	GetTypeID() *int32
	GetAttributes() *InferenceServiceAttributes
	GetProperties() *[]Properties
	GetCustomProperties() *[]Properties
}

type InferenceServiceImpl struct {
	ID               *int32
	TypeID           *int32
	Attributes       *InferenceServiceAttributes
	Properties       *[]Properties
	CustomProperties *[]Properties
}

func (r *InferenceServiceImpl) GetID() *int32 {
	return r.ID
}

func (r *InferenceServiceImpl) GetTypeID() *int32 {
	return r.TypeID
}

func (r *InferenceServiceImpl) GetAttributes() *InferenceServiceAttributes {
	return r.Attributes
}

func (r *InferenceServiceImpl) GetProperties() *[]Properties {
	return r.Properties
}

func (r *InferenceServiceImpl) GetCustomProperties() *[]Properties {
	return r.CustomProperties
}

type InferenceServiceRepository interface {
	GetByID(id int32) (InferenceService, error)
	List(listOptions InferenceServiceListOptions) (*ListWrapper[InferenceService], error)
	Save(model InferenceService) (InferenceService, error)
}
