package models

import "github.com/kubeflow/model-registry/internal/db/filter"

type ExperimentListOptions struct {
	Pagination
	Name       *string
	ExternalID *string
}

// GetRestEntityType implements the FilterApplier interface
func (e *ExperimentListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityExperiment
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
