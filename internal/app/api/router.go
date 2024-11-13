package api

import "github.com/gin-gonic/gin"

func getRoutes(r *gin.Engine) {
	users := r.Group("users")
	{
		users.POST("/", userHandler.Create)
		users.DELETE("/", authMiddleware.Authenticate, userHandler.Delete)
		users.PUT("/password", authMiddleware.Authenticate, userHandler.UpdatePassword)
	}

	auth := r.Group("auth")
	{
		auth.GET("/user_id", authHandler.GetUserID)
		auth.POST("/signin", authHandler.Signin)
		auth.DELETE("signout", authHandler.Signout)
	}
}
