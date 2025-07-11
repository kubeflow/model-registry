package models

const MetricHistoryType = "metric-history"

type MetricHistoryListOptions struct {
	Pagination
	Name            *string
	ExternalID      *string
	ExperimentRunID *int32
	StepIds         *string
}

type MetricHistoryAttributes struct {
	Name                     *string
	URI                      *string
	State                    *string
	ArtifactType             *string
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type MetricHistory interface {
	Entity[MetricHistoryAttributes]
}

type MetricHistoryImpl = BaseEntity[MetricHistoryAttributes]

type MetricHistoryRepository interface {
	GetByID(id int32) (MetricHistory, error)
	List(listOptions MetricHistoryListOptions) (*ListWrapper[MetricHistory], error)
	Save(metricHistory MetricHistory, experimentRunID *int32) (MetricHistory, error)
}
