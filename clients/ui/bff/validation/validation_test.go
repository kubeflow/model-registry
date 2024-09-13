package validation

import (
	"github.com/kubeflow/model-registry/pkg/openapi"
	"testing"
)

func TestValidateRegisteredModel(t *testing.T) {
	specs := []TestSpec[openapi.RegisteredModel]{
		{
			name:    "Empty name",
			input:   openapi.RegisteredModel{Name: ""},
			wantErr: true,
		},
		{
			name:    "Valid name",
			input:   openapi.RegisteredModel{Name: "ValidName"},
			wantErr: false,
		},
	}

	ValidateTestSpecs(t, specs, ValidateRegisteredModel)
}

func TestValidateModelVersion(t *testing.T) {
	specs := []TestSpec[openapi.ModelVersion]{
		{
			name:    "Empty name",
			input:   openapi.ModelVersion{Name: ""},
			wantErr: true,
		},
		{
			name:    "Valid name",
			input:   openapi.ModelVersion{Name: "ValidName"},
			wantErr: false,
		},
	}

	ValidateTestSpecs(t, specs, ValidateModelVersion)
}
