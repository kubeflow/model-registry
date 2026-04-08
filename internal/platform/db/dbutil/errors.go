package dbutil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/jackc/pgx/v5/pgconn"
	platformerrors "github.com/kubeflow/model-registry/internal/platform/errors"
)

func IsDatabaseTypeConversionError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

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
// and converts it to a BadRequest error.
func SanitizeDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	if IsDatabaseTypeConversionError(err) {
		glog.Warningf("Database type conversion error: %v", err)
		return fmt.Errorf("invalid filter query, type mismatch or invalid value in comparison: %w", platformerrors.ErrBadRequest)
	}
	return err
}

// EnhanceFilterQueryError provides user-friendly error messages for common filter query parsing mistakes.
func EnhanceFilterQueryError(err error, filterQuery string) error {
	errMsg := err.Error()

	if strings.Contains(errMsg, "expected Value") {
		return fmt.Errorf("invalid filter query syntax: %v. Hint: String values must be quoted (e.g., name=\"value\"). Numbers and booleans don't need quotes (e.g., id>5, active=true)", err)
	}

	if strings.Contains(errMsg, "unexpected token") {
		return fmt.Errorf("invalid filter query syntax: %v. Check that operators and values are correctly formatted", err)
	}

	return fmt.Errorf("invalid filter query syntax: %v", err)
}

// IsDuplicateKeyError checks if a database error is caused by a unique constraint violation.
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}

	return false
}
