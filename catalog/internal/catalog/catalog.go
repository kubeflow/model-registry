package catalog

import (
	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog/mcpcatalog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog/modelcatalog"
)

type (
	ModelSource = basecatalog.ModelSource
)

type (
	SourceCollection               = modelcatalog.SourceCollection
	LabelCollection                = modelcatalog.LabelCollection
	ModelProviderRecord            = modelcatalog.ModelProviderRecord
	APIProvider                    = modelcatalog.APIProvider
	ListModelsParams               = modelcatalog.ListModelsParams
	ListArtifactsParams            = modelcatalog.ListArtifactsParams
	ListPerformanceArtifactsParams = modelcatalog.ListPerformanceArtifactsParams

	// MCP catalog types
	MCPProvider                   = mcpcatalog.MCPCatalogProvider
	ListMCPServersParams          = mcpcatalog.ListMCPServersParams
	ListMCPServerToolsParams      = mcpcatalog.ListMCPServerToolsParams
)

var (
	NewSourceCollection         = modelcatalog.NewSourceCollection
	NewLabelCollection          = modelcatalog.NewLabelCollection
	NewDBCatalog                = modelcatalog.NewDBCatalog
	NewPerformanceMetricsLoader = modelcatalog.NewPerformanceMetricsLoader

	// MCP catalog constructors
	NewDBMCPCatalog = mcpcatalog.NewDBMCPCatalog
)

var (
	ParsePreviewConfig  = modelcatalog.ParsePreviewConfig
	PreviewSourceModels = modelcatalog.PreviewSourceModels
)
