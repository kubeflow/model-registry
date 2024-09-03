package validation

import (
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateRegisteredModel(t *testing.T) {
	tests := []struct {
		name    string
		input   openapi.RegisteredModel
		wantErr bool
	}{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRegisteredModel(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
