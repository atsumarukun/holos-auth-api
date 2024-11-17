package parameter

import (
	"fmt"
	"holos-auth-api/internal/pkg/apierr"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetPathParameter[T any](c *gin.Context, name string) (T, apierr.ApiError) {
	var zero T
	param := c.Param(name)

	switch any(zero).(type) {
	case uuid.UUID:
		v, err := uuid.Parse(param)
		if err != nil {
			return zero, apierr.NewApiError(http.StatusBadRequest, err.Error())
		}
		return any(v).(T), nil
	default:
		return zero, apierr.NewApiError(http.StatusInternalServerError, "invalid path parameter type")
	}
}

func GetContextParameter[T any](c *gin.Context, name string) (T, apierr.ApiError) {
	var zero T

	param, exists := c.Get(name)
	if !exists {
		return zero, apierr.NewApiError(http.StatusInternalServerError, fmt.Sprintf("context does not have %s", name))
	}

	v, ok := param.(T)
	if !ok {
		return zero, apierr.NewApiError(http.StatusInternalServerError, "invalid context parameter type")
	}

	return v, nil
}