package openapi

import (
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

// AssertCatalogArtifactRequired checks if the required fields are not zero-ed
func AssertCatalogArtifactRequired(obj model.CatalogArtifact) error {
	// CatalogArtifact has no required fields but the openapi code gen
	// checks the fields from CatalogModelArtifact, which doesn't compile.
	return nil
}
