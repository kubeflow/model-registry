package models

import (
	"github.com/kubeflow/model-registry/internal/db/constants"
	"github.com/kubeflow/model-registry/internal/db/filter"
)

var (
	// Use centralized state mappings from constants package
	Execution_State_name  = constants.ExecutionStateNames
	Execution_State_value = constants.ExecutionStateMapping
)

type ServeModelListOptions struct {
	Pagination
	Name               *string
	ExternalID         *string
	InferenceServiceID *int32
}

// GetRestEntityType implements the FilterApplier interface
func (s *ServeModelListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityServeModel
}

type ServeModelAttributes struct {
	Name                     *string
	ExternalID               *string
	LastKnownState           *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type ServeModel interface {
	Entity[ServeModelAttributes]
}

type ServeModelImpl = BaseEntity[ServeModelAttributes]

type ServeModelRepository interface {
	GetByID(id int32) (ServeModel, error)
	List(listOptions ServeModelListOptions) (*ListWrapper[ServeModel], error)
	Save(serveModel ServeModel, inferenceServiceID *int32) (ServeModel, error)
}
