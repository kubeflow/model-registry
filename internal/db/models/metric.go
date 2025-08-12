package models

const MetricType = "metric"

type MetricListOptions struct {
	Pagination
	Name             *string
	ExternalID       *string
	ParentResourceID *int32
}

type MetricAttributes struct {
	Name                     *string
	URI                      *string
	State                    *string
	ArtifactType             *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type Metric interface {
	Entity[MetricAttributes]
}

type MetricImpl = BaseEntity[MetricAttributes]

type MetricRepository interface {
	GetByID(id int32) (Metric, error)
	List(listOptions MetricListOptions) (*ListWrapper[Metric], error)
	Save(metric Metric, parentResourceID *int32) (Metric, error)
}
