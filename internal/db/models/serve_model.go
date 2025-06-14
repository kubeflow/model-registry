package models

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

type ServeModelAttributes struct {
	Name                     *string
	ExternalID               *string
	LastKnownState           *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type ServeModel interface {
	GetID() *int32
	GetTypeID() *int32
	GetAttributes() *ServeModelAttributes
	GetProperties() *[]Properties
	GetCustomProperties() *[]Properties
}

type ServeModelImpl struct {
	ID               *int32
	TypeID           *int32
	Attributes       *ServeModelAttributes
	Properties       *[]Properties
	CustomProperties *[]Properties
}

func (r *ServeModelImpl) GetID() *int32 {
	return r.ID
}

func (r *ServeModelImpl) GetTypeID() *int32 {
	return r.TypeID
}

func (r *ServeModelImpl) GetAttributes() *ServeModelAttributes {
	return r.Attributes
}

func (r *ServeModelImpl) GetProperties() *[]Properties {
	return r.Properties
}

func (r *ServeModelImpl) GetCustomProperties() *[]Properties {
	return r.CustomProperties
}

type ServeModelRepository interface {
	GetByID(id int32) (ServeModel, error)
	List(listOptions ServeModelListOptions) (*ListWrapper[ServeModel], error)
	Save(serveModel ServeModel, inferenceServiceID *int32) (ServeModel, error)
}
