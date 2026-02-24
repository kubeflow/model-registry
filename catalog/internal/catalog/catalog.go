package catalog

import (
	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog/modelcatalog"
)

// Re-export types that cmd needs from basecatalog
type (
	Source       = basecatalog.Source
	FieldFilter  = basecatalog.FieldFilter
	SourceConfig = basecatalog.SourceConfig
)

// Re-export types from modelcatalog
type (
	Loader                         = modelcatalog.Loader
	SourceCollection               = modelcatalog.SourceCollection
	LabelCollection                = modelcatalog.LabelCollection
	ModelProviderRecord            = modelcatalog.ModelProviderRecord
	PerformanceMetricsLoader       = modelcatalog.PerformanceMetricsLoader
	APIProvider                    = modelcatalog.APIProvider
	ListModelsParams               = modelcatalog.ListModelsParams
	ListArtifactsParams            = modelcatalog.ListArtifactsParams
	ListPerformanceArtifactsParams = modelcatalog.ListPerformanceArtifactsParams
	LoaderEventHandler             = modelcatalog.LoaderEventHandler
	PreviewConfig                  = modelcatalog.PreviewConfig
)

// Re-export factory functions from modelcatalog
var (
	NewLoader                   = modelcatalog.NewLoader
	NewSourceCollection         = modelcatalog.NewSourceCollection
	NewLabelCollection          = modelcatalog.NewLabelCollection
	NewDBCatalog                = modelcatalog.NewDBCatalog
	NewPerformanceMetricsLoader = modelcatalog.NewPerformanceMetricsLoader
	RegisterModelProvider       = modelcatalog.RegisterModelProvider
)

// Re-export validation functions from basecatalog
var (
	ValidateNamedQueries = basecatalog.ValidateNamedQueries
)

// Re-export validation and preview functions from modelcatalog
var (
	ValidateSourceFilters = modelcatalog.ValidateSourceFilters
	ParsePreviewConfig    = modelcatalog.ParsePreviewConfig
	PreviewSourceModels   = modelcatalog.PreviewSourceModels
)
