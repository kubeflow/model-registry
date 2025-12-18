package catalog

import (
	"github.com/kubeflow/model-registry/catalog/internal/common"
)

// AssetType is an alias to common.AssetType for convenience
type AssetType = common.AssetType

// Re-export constants from common package
const (
	AssetTypeModels     = common.AssetTypeModels
	AssetTypeMcpServers = common.AssetTypeMcpServers
)

// SourceProperties is an alias to common.SourceProperties
type SourceProperties = common.SourceProperties
