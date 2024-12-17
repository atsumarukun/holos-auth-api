package errors

import (
	"holos-auth-api/internal/app/api/pkg/status"
	"net/http"
)

var (
	StatusOK                  = status.New(http.StatusOK, "success")
	StatusBadRequest          = status.New(http.StatusBadRequest, "invalid request")
	StatusUnauthorized        = status.New(http.StatusUnauthorized, "unauthorized")
	StatusNotFound            = status.New(http.StatusNotFound, "resource not found")
	StatusInternalServerError = status.New(http.StatusInternalServerError, "internal server error")
)

func HandleError(err error) *status.Status {
	if err == nil {
		return StatusOK
	}

	s := status.FromError(err)

	switch s.Code() {
	case http.StatusBadRequest:
		// bad requestの時のみstatusをそのまま返却する.
		return s
	case http.StatusUnauthorized:
		return StatusUnauthorized
	case http.StatusNotFound:
		return StatusNotFound
	default:
		return StatusInternalServerError
	}
}
