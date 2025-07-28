package api

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
)

func ErrToStatus(err error) int {
	if errors.Is(err, ErrBadRequest) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound
	}

	if errors.Is(err, ErrConflict) {
		return http.StatusConflict
	}

	// Default error to return
	return http.StatusInternalServerError
}

func IgnoreNotFound(err error) error {
	if errors.Is(err, ErrNotFound) {
		return nil
	}

	return err
}
