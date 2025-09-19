package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStepIds(t *testing.T) {
	testCases := []struct {
		name        string
		stepIds     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty string should be valid",
			stepIds:     "",
			expectError: false,
		},
		{
			name:        "single valid step ID",
			stepIds:     "1",
			expectError: false,
		},
		{
			name:        "multiple valid step IDs",
			stepIds:     "1,2,3",
			expectError: false,
		},
		{
			name:        "valid step IDs with whitespace",
			stepIds:     " 1 , 2 , 3 ",
			expectError: false,
		},
		{
			name:        "valid step IDs with empty parts",
			stepIds:     "1,,3",
			expectError: false,
		},
		{
			name:        "invalid Unicode character",
			stepIds:     "耪",
			expectError: true,
			errorMsg:    "invalid step ID '耪': must be a valid integer",
		},
		{
			name:        "mixed valid and invalid step IDs",
			stepIds:     "1,invalid,3",
			expectError: true,
			errorMsg:    "invalid step ID 'invalid': must be a valid integer",
		},
		{
			name:        "non-numeric characters",
			stepIds:     "abc",
			expectError: true,
			errorMsg:    "invalid step ID 'abc': must be a valid integer",
		},
		{
			name:        "numeric with non-numeric suffix",
			stepIds:     "1a",
			expectError: true,
			errorMsg:    "invalid step ID '1a': must be a valid integer",
		},
		{
			name:        "negative numbers should be valid",
			stepIds:     "-1,0,1",
			expectError: false,
		},
		{
			name:        "large numbers should be valid",
			stepIds:     "2147483647",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateStepIds(tc.stepIds)

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
