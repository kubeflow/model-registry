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
	MCPServerTypeName              = "kf.MCPServer"
	MCPServerToolTypeName          = "kf.MCPServerTool"
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
		AddContext(MCPServerTypeName, datastore.NewSpecType(NewMCPServerRepository).
			AddString("source_id").
			AddString("base_name").
			AddString("description").
			AddString("provider").
			AddString("license").
			AddString("license_link").
			AddString("logo").
			AddString("readme").
			AddString("version").
			AddStruct("tags").
			AddStruct("transports").
			AddString("deploymentMode").
			AddBoolean("verifiedSource").
			AddBoolean("secureEndpoint").
			AddBoolean("sast").
			AddBoolean("readOnlyTools"),
		).
		AddExecution(MCPServerToolTypeName, datastore.NewSpecType(NewMCPServerToolRepository).
			AddString("accessType").
			AddString("description").
			AddString("externalId").
			AddString("parameters"),
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
	MCPServerRepository              models.MCPServerRepository
	MCPServerToolRepository          models.MCPServerToolRepository
}

func NewServices(
	catalogModelRepository models.CatalogModelRepository,
	catalogArtifactRepository models.CatalogArtifactRepository,
	catalogModelArtifactRepository models.CatalogModelArtifactRepository,
	catalogMetricsArtifactRepository models.CatalogMetricsArtifactRepository,
	catalogSourceRepository models.CatalogSourceRepository,
	propertyOptionsRepository models.PropertyOptionsRepository,
	mcpServerRepository models.MCPServerRepository,
	mcpServerToolRepository models.MCPServerToolRepository,
) Services {
	return Services{
		CatalogModelRepository:           catalogModelRepository,
		CatalogArtifactRepository:        catalogArtifactRepository,
		CatalogModelArtifactRepository:   catalogModelArtifactRepository,
		CatalogMetricsArtifactRepository: catalogMetricsArtifactRepository,
		CatalogSourceRepository:          catalogSourceRepository,
		PropertyOptionsRepository:        propertyOptionsRepository,
		MCPServerRepository:              mcpServerRepository,
		MCPServerToolRepository:          mcpServerToolRepository,
	}
}
