package model

import "github.com/kubeflow/model-registry/pkg/catalog/plugin"

func init() {
	plugin.Register(&ModelCatalogPlugin{})
}
