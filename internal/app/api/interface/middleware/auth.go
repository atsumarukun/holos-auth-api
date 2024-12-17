package middleware

import (
	"context"
	"holos-auth-api/internal/app/api/interface/pkg/errors"
	"holos-auth-api/internal/app/api/usecase"
	"log"
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

func (m *authMiddleware) Authenticate(c *gin.Context) {
	bearerToken := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		status := errors.StatusUnauthorized
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		c.Abort()
		return
	}

	ctx := context.Background()

	userID, err := m.authUsecase.GetUserID(ctx, bearerToken[1])
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		c.Abort()
		return
	}

	c.Set("userID", userID)
	c.Next()
}
