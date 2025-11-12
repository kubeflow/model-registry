package dbutil

import (
	"errors"
	"strings"
	"testing"

	"github.com/kubeflow/model-registry/pkg/api"
)

func TestIsDatabaseTypeConversionError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "PostgreSQL invalid input syntax",
			err:      errors.New(`ERROR: invalid input syntax for type double precision: "x" (SQLSTATE 22P02)`),
			expected: true,
		},
		{
			name:     "PostgreSQL invalid input syntax - uppercase",
			err:      errors.New(`ERROR: INVALID INPUT SYNTAX for type integer: "abc"`),
			expected: true,
		},
		{
			name:     "PostgreSQL SQLSTATE 22P02",
			err:      errors.New(`pq: invalid input syntax for type numeric: "test" (SQLSTATE 22P02)`),
			expected: true,
		},
		{
			name:     "PostgreSQL SQLSTATE 22003",
			err:      errors.New(`pq: numeric field overflow (SQLSTATE 22003)`),
			expected: true,
		},
		{
			name:     "PostgreSQL SQLSTATE in different case",
			err:      errors.New(`some error with SqlState 22P02`),
			expected: true,
		},
		{
			name:     "MySQL incorrect double value",
			err:      errors.New(`Error 1366: Incorrect double value: 'x' for column 'value'`),
			expected: true,
		},
		{
			name:     "MySQL incorrect double value - mixed case",
			err:      errors.New(`Error 1366: INCORRECT DOUBLE VALUE: 'text' for column`),
			expected: true,
		},
		{
			name:     "MySQL incorrect integer value",
			err:      errors.New(`Error 1366: Incorrect integer value: 'abc' for column 'count'`),
			expected: true,
		},
		{
			name:     "MySQL truncated incorrect",
			err:      errors.New(`Error 1265: Data truncated incorrect DOUBLE value: 'notanumber'`),
			expected: true,
		},
		{
			name:     "PostgreSQL unable to encode numeric to text",
			err:      errors.New("failed to encode args[4]: unable to encode 5 into text format for text (OID 25): cannot find encode plan"),
			expected: true,
		},
		{
			name:     "PostgreSQL failed to encode",
			err:      errors.New("failed to encode args[2]: type mismatch"),
			expected: true,
		},
		{
			name:     "PostgreSQL cannot find encode plan - mixed case",
			err:      errors.New("Unable to Encode value: Cannot Find Encode Plan"),
			expected: true,
		},
		{
			name:     "non-matching error - record not found",
			err:      errors.New("record not found"),
			expected: false,
		},
		{
			name:     "non-matching error - constraint violation",
			err:      errors.New("UNIQUE constraint failed"),
			expected: false,
		},
		{
			name:     "non-matching error - connection error",
			err:      errors.New("connection refused"),
			expected: false,
		},
		{
			name:     "non-matching error - generic SQL error",
			err:      errors.New("SQL syntax error near 'SELECT'"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDatabaseTypeConversionError(tt.err)
			if result != tt.expected {
				t.Errorf("isDatabaseTypeConversionError() = %v, want %v for error: %v", result, tt.expected, tt.err)
			}
		})
	}
}

func TestSanitizeDatabaseError(t *testing.T) {
	tests := []struct {
		name           string
		inputErr       error
		expectModified bool
		expectBadReq   bool
		expectMessage  string
	}{
		{
			name:           "nil error returns nil",
			inputErr:       nil,
			expectModified: false,
			expectBadReq:   false,
		},
		{
			name:           "PostgreSQL type conversion error is sanitized",
			inputErr:       errors.New(`ERROR: invalid input syntax for type double precision: "x" (SQLSTATE 22P02)`),
			expectModified: true,
			expectBadReq:   true,
			expectMessage:  "invalid filter query, type mismatch or invalid value in comparison",
		},
		{
			name:           "MySQL type conversion error is sanitized",
			inputErr:       errors.New(`Error 1366: Incorrect double value: 'text' for column 'price'`),
			expectModified: true,
			expectBadReq:   true,
			expectMessage:  "invalid filter query, type mismatch or invalid value in comparison",
		},
		{
			name:           "non-type-conversion error passes through unchanged",
			inputErr:       errors.New("record not found"),
			expectModified: false,
			expectBadReq:   false,
		},
		{
			name:           "connection error passes through unchanged",
			inputErr:       errors.New("connection refused"),
			expectModified: false,
			expectBadReq:   false,
		},
		{
			name:           "generic database error passes through unchanged",
			inputErr:       errors.New("UNIQUE constraint violation"),
			expectModified: false,
			expectBadReq:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeDatabaseError(tt.inputErr)

			// Check if error is nil when expected
			if tt.inputErr == nil {
				if result != nil {
					t.Errorf("SanitizeDatabaseError() returned non-nil for nil input: %v", result)
				}
				return
			}

			// Check if error was modified
			if tt.expectModified {
				if result == tt.inputErr {
					t.Errorf("SanitizeDatabaseError() should have modified the error but returned same instance")
				}

				// Check if message contains expected text
				if !strings.Contains(result.Error(), tt.expectMessage) {
					t.Errorf("SanitizeDatabaseError() error message = %q, want to contain %q", result.Error(), tt.expectMessage)
				}

				// Check if it wraps api.ErrBadRequest
				if tt.expectBadReq {
					if !errors.Is(result, api.ErrBadRequest) {
						t.Errorf("SanitizeDatabaseError() should wrap api.ErrBadRequest but does not")
					}
				}
			} else {
				// Error should pass through unchanged
				if result != tt.inputErr {
					t.Errorf("SanitizeDatabaseError() modified non-type-conversion error: got %v, want %v", result, tt.inputErr)
				}
			}
		})
	}
}

func TestSanitizeDatabaseError_WrapsCorrectly(t *testing.T) {
	// Test that the sanitized error properly wraps api.ErrBadRequest
	inputErr := errors.New(`ERROR: invalid input syntax for type integer: "abc"`)
	result := SanitizeDatabaseError(inputErr)

	// Should wrap api.ErrBadRequest
	if !errors.Is(result, api.ErrBadRequest) {
		t.Error("SanitizeDatabaseError() should wrap api.ErrBadRequest")
	}

	// Original database error should NOT be in the chain (we don't want to expose it)
	if errors.Is(result, inputErr) {
		t.Error("SanitizeDatabaseError() should not wrap the original database error")
	}

	// Message should be generic
	if strings.Contains(result.Error(), "SQLSTATE") {
		t.Error("SanitizeDatabaseError() should not expose SQLSTATE in error message")
	}
	if strings.Contains(result.Error(), "invalid input syntax") {
		t.Error("SanitizeDatabaseError() should not expose database-specific error details")
	}
}

func TestSanitizeDatabaseError_CaseInsensitive(t *testing.T) {
	// Test various case combinations to ensure case-insensitive matching
	testCases := []string{
		"INVALID INPUT SYNTAX for type double",
		"Invalid Input Syntax for type integer",
		"invalid input syntax for type numeric",
		"SQLSTATE 22P02",
		"sqlstate 22003",
		"SqlState 22P03",
		"INCORRECT DOUBLE VALUE: 'x'",
		"Incorrect Double Value: 'y'",
		"incorrect integer value: 'z'",
	}

	for _, errMsg := range testCases {
		t.Run(errMsg, func(t *testing.T) {
			inputErr := errors.New(errMsg)
			result := SanitizeDatabaseError(inputErr)

			if !errors.Is(result, api.ErrBadRequest) {
				t.Errorf("Case-insensitive matching failed for: %s", errMsg)
			}
		})
	}
}

func TestEnhanceFilterQueryError(t *testing.T) {
	tests := []struct {
		name                string
		inputError          error
		filterQuery         string
		expectedErrContains []string
	}{
		{
			name:        "expected Value error - unquoted string",
			inputError:  errors.New("error parsing filter query: 1:5: unexpected token \"bar\" (expected Value)"),
			filterQuery: "foo=bar",
			expectedErrContains: []string{
				"invalid filter query syntax",
				"expected Value",
				"Hint: String values must be quoted",
				"name=\"value\"",
			},
		},
		{
			name:        "unexpected token error",
			inputError:  errors.New("error parsing filter query: 1:10: unexpected token \"@\""),
			filterQuery: "name=test@",
			expectedErrContains: []string{
				"invalid filter query syntax",
				"unexpected token",
				"Check that operators and values are correctly formatted",
			},
		},
		{
			name:        "generic parse error",
			inputError:  errors.New("error parsing filter query: unknown operator"),
			filterQuery: "name ?? value",
			expectedErrContains: []string{
				"invalid filter query syntax",
				"unknown operator",
			},
		},
		{
			name:        "another expected Value error",
			inputError:  errors.New("1:15: expected Value, got something else"),
			filterQuery: "id>5 AND name=test",
			expectedErrContains: []string{
				"invalid filter query syntax",
				"expected Value",
				"Hint: String values must be quoted",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EnhanceFilterQueryError(tt.inputError, tt.filterQuery)

			if result == nil {
				t.Error("expected non-nil error, got nil")
				return
			}

			resultMsg := result.Error()

			// Check that all expected strings are in the error message
			for _, expected := range tt.expectedErrContains {
				if !strings.Contains(resultMsg, expected) {
					t.Errorf("expected error message to contain %q, got: %q", expected, resultMsg)
				}
			}
		})
	}
}
