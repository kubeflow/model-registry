package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testSpec[T any] struct {
	name    string
	input   T
	wantErr bool
}

func validateTestSpecs[T any](t *testing.T, specs []testSpec[T], validator func(input T) error) {
	for _, tt := range specs {
		t.Run(tt.name, func(t *testing.T) {
			err := validator(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
