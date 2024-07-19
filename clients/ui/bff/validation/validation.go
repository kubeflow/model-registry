package validation

import (
	"errors"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

func ValidateRegisteredModel(input openapi.RegisteredModel) error {
	if input.Name != nil && *input.Name == "" {
		return errors.New("name cannot be empty")
	}
	// Add more field validations as required
	return nil
}
