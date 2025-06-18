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
	Entity[InferenceServiceAttributes]
}

type InferenceServiceImpl = BaseEntity[InferenceServiceAttributes]

type InferenceServiceRepository interface {
	GetByID(id int32) (InferenceService, error)
	List(listOptions InferenceServiceListOptions) (*ListWrapper[InferenceService], error)
	Save(model InferenceService) (InferenceService, error)
}
