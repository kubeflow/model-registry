package models

type ExperimentListOptions struct {
	Pagination
	Name       *string
	ExternalID *string
}

type ExperimentAttributes struct {
	Name                     *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type Experiment interface {
	Entity[ExperimentAttributes]
}

type ExperimentImpl = BaseEntity[ExperimentAttributes]

type ExperimentRepository interface {
	GetByID(id int32) (Experiment, error)
	List(listOptions ExperimentListOptions) (*ListWrapper[Experiment], error)
	Save(experiment Experiment) (Experiment, error)
}
