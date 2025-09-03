package models

const ParameterType = "parameter"

type ParameterListOptions struct {
	Pagination
	Name             *string
	ExternalID       *string
	ParentResourceID *int32
}

type ParameterAttributes struct {
	Name                     *string
	URI                      *string
	State                    *string
	ArtifactType             *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type Parameter interface {
	Entity[ParameterAttributes]
}

type ParameterImpl = BaseEntity[ParameterAttributes]

type ParameterRepository interface {
	GetByID(id int32) (Parameter, error)
	List(listOptions ParameterListOptions) (*ListWrapper[Parameter], error)
	Save(parameter Parameter, parentResourceID *int32) (Parameter, error)
}
