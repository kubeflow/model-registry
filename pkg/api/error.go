package api

import (
	platformerrors "github.com/kubeflow/model-registry/internal/platform/errors"
)

// Re-export from platform/errors for backward compatibility
var (
	ErrBadRequest = platformerrors.ErrBadRequest
	ErrNotFound   = platformerrors.ErrNotFound
	ErrConflict   = platformerrors.ErrConflict
)

var ErrToStatus = platformerrors.ErrToStatus
var IgnoreNotFound = platformerrors.IgnoreNotFound
