package models

import (
	"github.com/kubeflow/model-registry/internal/db/filter"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
)

type MetricsType string

const (
	MetricsTypePerformance     MetricsType = "performance-metrics"
	MetricsTypeAccuracy        MetricsType = "accuracy-metrics"
	CatalogMetricsArtifactType             = "metrics-artifact"
)

type CatalogMetricsArtifactListOptions struct {
	dbmodels.Pagination
	Name             *string
	ExternalID       *string
	ParentResourceID *int32
}

// GetRestEntityType implements the FilterApplier interface
func (c *CatalogMetricsArtifactListOptions) GetRestEntityType() filter.RestEntityType {
	return filter.RestEntityModelArtifact // Reusing existing filter type
}

type CatalogMetricsArtifactAttributes struct {
	Name                     *string
	ArtifactType             *string
	MetricsType              MetricsType
	ExternalID               *string
	CreateTimeSinceEpoch     *int64
	LastUpdateTimeSinceEpoch *int64
}

type CatalogMetricsArtifact interface {
	dbmodels.Entity[CatalogMetricsArtifactAttributes]
}

type CatalogMetricsArtifactImpl = dbmodels.BaseEntity[CatalogMetricsArtifactAttributes]

type CatalogMetricsArtifactRepository interface {
	GetByID(id int32) (CatalogMetricsArtifact, error)
	List(listOptions CatalogMetricsArtifactListOptions) (*dbmodels.ListWrapper[CatalogMetricsArtifact], error)
	Save(metricsArtifact CatalogMetricsArtifact, parentResourceID *int32) (CatalogMetricsArtifact, error)
	// BatchSave inserts multiple metrics artifacts in a single batch operation
	BatchSave(metricsArtifacts []CatalogMetricsArtifact, parentResourceID *int32) ([]CatalogMetricsArtifact, error)
}
