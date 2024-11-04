package api

import "github.com/gin-gonic/gin"

func getRoutes(r *gin.Engine) {
	users := r.Group("users")
	{
		users.POST("/", userHandler.Create)
		users.PUT("/:name", userHandler.Update)
		users.DELETE("/:name", userHandler.Delete)
	}
}
