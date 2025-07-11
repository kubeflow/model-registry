package models

import "github.com/kubeflow/model-registry/internal/db/filter"

type ExperimentRunListOptions struct {
	Pagination
	Name         *string
	ExternalID   *string
	ExperimentID *int32
}

// GetRestEntityType implements the FilterApplier interface
func (e *ExperimentRunListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityExperimentRun
}

type ExperimentRunAttributes struct {
	Name                     *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type ExperimentRun interface {
	Entity[ExperimentRunAttributes]
}

type ExperimentRunImpl = BaseEntity[ExperimentRunAttributes]

type ExperimentRunRepository interface {
	GetByID(id int32) (ExperimentRun, error)
	List(listOptions ExperimentRunListOptions) (*ListWrapper[ExperimentRun], error)
	Save(experimentRun ExperimentRun, experimentID *int32) (ExperimentRun, error)
}
