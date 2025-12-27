package catalog

import (
	"github.com/kubeflow/model-registry/catalog/internal/common"
)

// DetectYamlAssetType reads a YAML file and detects what type of assets it contains.
// It returns an error if the file has multiple asset types or no recognized asset types.
// This is a convenience wrapper around common.DetectYamlAssetType for Source types.
func DetectYamlAssetType(source *Source, reldir string) (AssetType, error) {
	return common.DetectYamlAssetType(source, reldir)
}
