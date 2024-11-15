package handler

import (
	"context"
	"holos-auth-api/internal/app/api/interface/request"
	"holos-auth-api/internal/app/api/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Signin(*gin.Context)
	Signout(*gin.Context)
	GetUserID(*gin.Context)
}

type authHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) AuthHandler {
	return &authHandler{
		authUsecase: authUsecase,
	}
}

func (ah *authHandler) Signin(c *gin.Context) {
	var req request.SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()

	token, err := ah.authUsecase.Signin(ctx, req.UserName, req.Password)
	if err != nil {
		c.String(err.Error())
		return
	}

	c.String(http.StatusOK, token)
}

func (ah *authHandler) Signout(c *gin.Context) {
	bearerToken := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		c.String(http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx := context.Background()

	if err := ah.authUsecase.Signout(ctx, bearerToken[1]); err != nil {
		c.String(err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (ah *authHandler) GetUserID(c *gin.Context) {
	bearerToken := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		c.String(http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx := context.Background()

	userID, err := ah.authUsecase.GetUserID(ctx, bearerToken[1])
	if err != nil {
		c.String(err.Error())
		return
	}

	c.String(http.StatusOK, userID.String())
}
