package api

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
)

func ErrToStatus(err error) int {
	// If the error is a gRPC error, we can extract the status code.
	if status, ok := status.FromError(err); ok {
		switch status.Code() {
		case codes.InvalidArgument:
			return http.StatusBadRequest
		case codes.AlreadyExists:
			return http.StatusConflict
		case codes.Unavailable:
			return http.StatusServiceUnavailable
		}
	}

	switch errors.Unwrap(err) {
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
