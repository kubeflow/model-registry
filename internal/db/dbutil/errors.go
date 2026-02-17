package dbutil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/jackc/pgx/v5/pgconn"
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

// IsDuplicateKeyError checks if a database error is caused by a unique constraint violation.
// Uses database-specific error codes for reliable detection across versions and locales.
// Returns true for both MySQL and PostgreSQL unique constraint errors.
//
// Error codes used:
//   - PostgreSQL: 23505 (unique_violation)
//   - MySQL: 1062 (ER_DUP_ENTRY - duplicate entry for key)
//
// References:
//   - PostgreSQL: https://www.postgresql.org/docs/current/errcodes-appendix.html
//   - MySQL: https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
//   - pgx: https://github.com/jackc/pgx/wiki/Error-Handling
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	// PostgreSQL: Check for unique_violation error code (23505)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}

	// MySQL: Check for ER_DUP_ENTRY error code (1062)
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062 // ER_DUP_ENTRY
	}

	return false
}
