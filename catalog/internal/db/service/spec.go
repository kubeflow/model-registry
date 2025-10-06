package service

import (
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/internal/datastore"
)

const (
	CatalogModelTypeName           = "kf.CatalogModel"
	CatalogModelArtifactTypeName   = "kf.CatalogModelArtifact"
	CatalogMetricsArtifactTypeName = "kf.CatalogMetricsArtifact"
)

func DatastoreSpec() *datastore.Spec {
	return datastore.NewSpec().
		AddContext(CatalogModelTypeName, datastore.NewSpecType(NewCatalogModelRepository).
			AddString("source_id").
			AddString("description").
			AddString("owner").
			AddString("state").
			AddStruct("language").
			AddString("library_name").
			AddString("license_link").
			AddString("license").
			AddStruct("logo").
			AddString("maturity").
			AddString("provider").
			AddStruct("readme").
			AddStruct("tasks"),
		).
		AddArtifact(CatalogModelArtifactTypeName, datastore.NewSpecType(NewCatalogModelArtifactRepository).
			AddString("uri"),
		).
		AddArtifact(CatalogMetricsArtifactTypeName, datastore.NewSpecType(NewCatalogMetricsArtifactRepository).
			AddString("metricsType"),
		).
		AddOther(NewCatalogArtifactRepository)
}

type Services struct {
	CatalogModelRepository           models.CatalogModelRepository
	CatalogArtifactRepository        models.CatalogArtifactRepository
	CatalogModelArtifactRepository   models.CatalogModelArtifactRepository
	CatalogMetricsArtifactRepository models.CatalogMetricsArtifactRepository
}

func NewServices(
	catalogModelRepository models.CatalogModelRepository,
	catalogArtifactRepository models.CatalogArtifactRepository,
	catalogModelArtifactRepository models.CatalogModelArtifactRepository,
	catalogMetricsArtifactRepository models.CatalogMetricsArtifactRepository,
) Services {
	return Services{
		CatalogModelRepository:           catalogModelRepository,
		CatalogArtifactRepository:        catalogArtifactRepository,
		CatalogModelArtifactRepository:   catalogModelArtifactRepository,
		CatalogMetricsArtifactRepository: catalogMetricsArtifactRepository,
	}
}
