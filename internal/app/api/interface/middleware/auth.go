package middleware

import (
	"context"
	"holos-auth-api/internal/app/api/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	Authenticate(*gin.Context)
}

type authMiddleware struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthMiddleware(authUsecase usecase.AuthUsecase) AuthMiddleware {
	return &authMiddleware{
		authUsecase: authUsecase,
	}
}

func (am *authMiddleware) Authenticate(c *gin.Context) {
	bearerToken := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		c.String(http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	ctx := context.Background()

	userID, err := am.authUsecase.GetUserID(ctx, bearerToken[1])
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		c.String(err.Error())
		c.Abort()
		return
	}

	c.Set("userID", userID)
	c.Next()
}
