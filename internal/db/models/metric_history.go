package models

import "github.com/kubeflow/model-registry/internal/db/filter"

const MetricHistoryType = "metric-history"

type MetricHistoryListOptions struct {
	Pagination
	Name            *string
	ExternalID      *string
	ExperimentRunID *int32
	StepIds         *string
}

// GetRestEntityType implements the FilterApplier interface
func (m *MetricHistoryListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityMetric // Metric history uses the same filtering rules as metrics
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
