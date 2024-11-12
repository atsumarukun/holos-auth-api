package api

import "github.com/gin-gonic/gin"

func getRoutes(r *gin.Engine) {
	users := r.Group("users")
	{
		users.POST("/", userHandler.Create)
		users.PUT("/", authMiddleware.Authenticate, userHandler.Update)
		users.DELETE("/", authMiddleware.Authenticate, userHandler.Delete)
	}

	auth := r.Group("auth")
	{
		auth.GET("/user_id", authHandler.GetUserID)
		auth.POST("/signin", authHandler.Signin)
		auth.DELETE("signout", authHandler.Signout)
	}
}
