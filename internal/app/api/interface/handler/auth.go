package handler

import (
	"holos-auth-api/internal/app/api/interface/pkg/errors"
	"holos-auth-api/internal/app/api/interface/request"
	"holos-auth-api/internal/app/api/usecase"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Signin(*gin.Context)
	Signout(*gin.Context)
	Authorize(*gin.Context)
}

type authHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) AuthHandler {
	return &authHandler{
		authUsecase: authUsecase,
	}
}

func (h *authHandler) Signin(c *gin.Context) {
	var req request.SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status := errors.StatusBadRequest
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	ctx := c.Request.Context()

	token, err := h.authUsecase.Signin(ctx, req.UserName, req.Password)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.String(http.StatusCreated, token)
}

func (h *authHandler) Signout(c *gin.Context) {
	bearerToken := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		status := errors.StatusUnauthorized
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	ctx := c.Request.Context()

	if err := h.authUsecase.Signout(ctx, bearerToken[1]); err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *authHandler) Authorize(c *gin.Context) {
	bearerToken := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		status := errors.StatusUnauthorized
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	operatorType := c.Request.Header.Get("Holos-Operator-Type")

	service := c.Query("service")
	path := c.Query("path")
	method := c.Query("method")

	ctx := c.Request.Context()

	userID, err := h.authUsecase.Authorize(ctx, bearerToken[1], operatorType, service, path, method)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.String(http.StatusOK, userID.String())
}
