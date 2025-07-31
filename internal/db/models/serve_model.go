package models

import "github.com/kubeflow/model-registry/internal/db/filter"

var (
	Execution_State_name = map[int32]string{
		0: "UNKNOWN",
		1: "NEW",
		2: "RUNNING",
		3: "COMPLETE",
		4: "FAILED",
		5: "CACHED",
		6: "CANCELED",
	}
	Execution_State_value = map[string]int32{
		"UNKNOWN":  0,
		"NEW":      1,
		"RUNNING":  2,
		"COMPLETE": 3,
		"FAILED":   4,
		"CACHED":   5,
		"CANCELED": 6,
	}
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
