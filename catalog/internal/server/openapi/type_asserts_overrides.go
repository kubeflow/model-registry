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

// AssertFilterOptionRequired checks if the required fields are not zero-ed
func AssertFilterOptionRequired(obj model.FilterOption) error {
	elements := map[string]interface{}{
		"type": obj.Type,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if obj.Range != nil {
		if err := AssertFilterOptionRangeRequired(*obj.Range); err != nil {
			return err
		}
	}
	return nil
}
