package dbutil

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/pkg/api"
)

// isDatabaseTypeConversionError checks if a database error is caused by type conversion issues
// (e.g., trying to convert string "x" to a number) that should return 400 Bad Request instead of 500
func isDatabaseTypeConversionError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	// type conversion error patterns (in lowercase for case-insensitive matching)
	errorPatterns := []string{
		"invalid input syntax",
		"sqlstate",
		"incorrect double value",
		"incorrect integer value",
		"truncated incorrect",
		"unable to encode",
		"failed to encode",
		"cannot find encode plan",
	}

	for _, pattern := range errorPatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// SanitizeDatabaseError checks if a database error is a type conversion error
// and converts it to a BadRequest error. Otherwise returns the original error.
// This prevents exposing internal database error details to users.
// Should be called after query execution to sanitize errors.
func SanitizeDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	if isDatabaseTypeConversionError(err) {
		// Log the actual database error internally for debugging
		glog.Warningf("Database type conversion error: %v", err)

		return fmt.Errorf("invalid filter query, type mismatch or invalid value in comparison: %w", api.ErrBadRequest)
	}
	return err
}

// EnhanceFilterQueryError provides user-friendly error messages for common filter query parsing mistakes.
// Should be called when filter query parsing fails to give users helpful hints.
func EnhanceFilterQueryError(err error, filterQuery string) error {
	errMsg := err.Error()

	// Check for common mistakes and provide helpful hints
	if strings.Contains(errMsg, "expected Value") {
		return fmt.Errorf("invalid filter query syntax: %v. Hint: String values must be quoted (e.g., name=\"value\"). Numbers and booleans don't need quotes (e.g., id>5, active=true)", err)
	}

	if strings.Contains(errMsg, "unexpected token") {
		return fmt.Errorf("invalid filter query syntax: %v. Check that operators and values are correctly formatted", err)
	}

	// Default: return the original error with a generic message
	return fmt.Errorf("invalid filter query syntax: %v", err)
}
