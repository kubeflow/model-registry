package models

import "github.com/kubeflow/model-registry/pkg/catalog/plugin"

func init() {
	plugin.Register(&ModelCatalogPlugin{})
}
