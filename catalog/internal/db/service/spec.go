package service

import (
	mcpcatalogmodels "github.com/kubeflow/model-registry/catalog/internal/catalog/mcpcatalog/models"
	mcpcatalogservice "github.com/kubeflow/model-registry/catalog/internal/catalog/mcpcatalog/service"
	modelcatalogmodels "github.com/kubeflow/model-registry/catalog/internal/catalog/modelcatalog/models"
	modelcatalogservice "github.com/kubeflow/model-registry/catalog/internal/catalog/modelcatalog/service"
	sharedmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
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
		AddContext(CatalogModelTypeName, datastore.NewSpecType(modelcatalogservice.NewCatalogModelRepository).
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
		AddContext(MCPServerTypeName, datastore.NewSpecType(mcpcatalogservice.NewMCPServerRepository).
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
		AddExecution(MCPServerToolTypeName, datastore.NewSpecType(mcpcatalogservice.NewMCPServerToolRepository).
			AddString("accessType").
			AddString("description").
			AddString("externalId").
			AddString("parameters"),
		).
		AddArtifact(CatalogModelArtifactTypeName, datastore.NewSpecType(modelcatalogservice.NewCatalogModelArtifactRepository).
			AddString("uri"),
		).
		AddArtifact(CatalogMetricsArtifactTypeName, datastore.NewSpecType(modelcatalogservice.NewCatalogMetricsArtifactRepository).
			AddString("metricsType"),
		).
		AddOther(NewCatalogArtifactRepository).
		AddOther(NewPropertyOptionsRepository)
}

type Services struct {
	CatalogModelRepository           modelcatalogmodels.CatalogModelRepository
	CatalogArtifactRepository        sharedmodels.CatalogArtifactRepository
	CatalogModelArtifactRepository   modelcatalogmodels.CatalogModelArtifactRepository
	CatalogMetricsArtifactRepository modelcatalogmodels.CatalogMetricsArtifactRepository
	CatalogSourceRepository          sharedmodels.CatalogSourceRepository
	PropertyOptionsRepository        sharedmodels.PropertyOptionsRepository
	MCPServerRepository              mcpcatalogmodels.MCPServerRepository
	MCPServerToolRepository          mcpcatalogmodels.MCPServerToolRepository
}

func NewServices(
	catalogModelRepository modelcatalogmodels.CatalogModelRepository,
	catalogArtifactRepository sharedmodels.CatalogArtifactRepository,
	catalogModelArtifactRepository modelcatalogmodels.CatalogModelArtifactRepository,
	catalogMetricsArtifactRepository modelcatalogmodels.CatalogMetricsArtifactRepository,
	catalogSourceRepository sharedmodels.CatalogSourceRepository,
	propertyOptionsRepository sharedmodels.PropertyOptionsRepository,
	mcpServerRepository mcpcatalogmodels.MCPServerRepository,
	mcpServerToolRepository mcpcatalogmodels.MCPServerToolRepository,
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
