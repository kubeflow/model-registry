package apiutils

import (
	"fmt"
	"strconv"

	platformerrors "github.com/kubeflow/model-registry/internal/platform/errors"
)

// ZeroIfNil return the zeroed value if input is a nil pointer
func ZeroIfNil[T any](input *T) T {
	if input != nil {
		return *input
	}
	return *new(T)
}

// Of returns a pointer to the provided literal/const input
func Of[E any](e E) *E {
	return &e
}

func StrPtr(notEmpty string) *string {
	if notEmpty == "" {
		return nil
	}
	return &notEmpty
}

// SafeString returns a string representation of a string pointer.
func SafeString(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}

// ValidateIDAsInt32 validates and converts a string ID to int32
func ValidateIDAsInt32(id string, entityName string) (int32, error) {
	convertedId, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s ID: %v: %w", entityName, err, platformerrors.ErrBadRequest)
	}
	return int32(convertedId), nil
}

// ValidateIDAsInt32Ptr validates and converts a string pointer ID to int32 pointer
func ValidateIDAsInt32Ptr(id *string, entityName string) (*int32, error) {
	if id == nil {
		return nil, nil
	}
	convertedId, err := strconv.ParseInt(*id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid %s ID: %v: %w", entityName, err, platformerrors.ErrBadRequest)
	}
	result := int32(convertedId)
	return &result, nil
}

// ValidateIDAsInt64 validates and converts a string ID to int64
func ValidateIDAsInt64(id string, entityName string) (int64, error) {
	convertedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s ID: %v: %w", entityName, err, platformerrors.ErrBadRequest)
	}
	return convertedId, nil
}

// ValidateIDAsInt64Ptr validates and converts a string pointer ID to int64 pointer
func ValidateIDAsInt64Ptr(id *string, entityName string) (*int64, error) {
	if id == nil {
		return nil, nil
	}
	convertedId, err := strconv.ParseInt(*id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid %s ID: %v: %w", entityName, err, platformerrors.ErrBadRequest)
	}
	return &convertedId, nil
}
