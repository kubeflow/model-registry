package service

import (
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/internal/datastore"
)

const (
	CatalogModelTypeName           = "kf.CatalogModel"
	CatalogModelArtifactTypeName   = "kf.CatalogModelArtifact"
	CatalogMetricsArtifactTypeName = "kf.CatalogMetricsArtifact"
	CatalogSourceTypeName          = "kf.CatalogSource"
	McpServerTypeName              = "kf.McpServer"
	McpServerToolTypeName          = "kf.McpServerTool"
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
			AddString("logo").
			AddString("maturity").
			AddString("provider").
			AddString("readme").
			AddStruct("tasks"),
		).
		AddContext(CatalogSourceTypeName, datastore.NewSpecType(NewCatalogSourceRepository).
			AddString("status").
			AddString("error"),
		).
		AddContext(McpServerTypeName, datastore.NewSpecType(NewMcpServerRepository).
			AddString("source_id").
			AddString("description").
			AddString("logo").
			AddString("provider").
			AddString("version").
			AddString("status").
			AddString("transport").
			AddString("category").
			AddStruct("tags").
			AddString("endpoint").
			AddString("documentationUrl").
			AddString("repositoryUrl").
			AddString("sourceCode").
			AddString("readme").
			AddString("publishedDate").
			AddBoolean("verifiedSource").
			AddBoolean("secureEndpoint").
			AddBoolean("sast").
			AddBoolean("readOnlyTools").
			AddStruct("tools"),
		).
		AddArtifact(CatalogModelArtifactTypeName, datastore.NewSpecType(NewCatalogModelArtifactRepository).
			AddString("uri"),
		).
		AddArtifact(CatalogMetricsArtifactTypeName, datastore.NewSpecType(NewCatalogMetricsArtifactRepository).
			AddString("metricsType"),
		).
		AddOther(NewCatalogArtifactRepository).
		AddOther(NewPropertyOptionsRepository)
}

type Services struct {
	CatalogModelRepository           models.CatalogModelRepository
	CatalogArtifactRepository        models.CatalogArtifactRepository
	CatalogModelArtifactRepository   models.CatalogModelArtifactRepository
	CatalogMetricsArtifactRepository models.CatalogMetricsArtifactRepository
	CatalogSourceRepository          models.CatalogSourceRepository
	PropertyOptionsRepository        models.PropertyOptionsRepository
	McpServerRepository              models.McpServerRepository
}

func NewServices(
	catalogModelRepository models.CatalogModelRepository,
	catalogArtifactRepository models.CatalogArtifactRepository,
	catalogModelArtifactRepository models.CatalogModelArtifactRepository,
	catalogMetricsArtifactRepository models.CatalogMetricsArtifactRepository,
	catalogSourceRepository models.CatalogSourceRepository,
	propertyOptionsRepository models.PropertyOptionsRepository,
	mcpServerRepository models.McpServerRepository,
) Services {
	return Services{
		CatalogModelRepository:           catalogModelRepository,
		CatalogArtifactRepository:        catalogArtifactRepository,
		CatalogModelArtifactRepository:   catalogModelArtifactRepository,
		CatalogMetricsArtifactRepository: catalogMetricsArtifactRepository,
		CatalogSourceRepository:          catalogSourceRepository,
		PropertyOptionsRepository:        propertyOptionsRepository,
		McpServerRepository:              mcpServerRepository,
	}
}
