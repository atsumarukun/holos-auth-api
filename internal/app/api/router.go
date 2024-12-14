package api

import "github.com/gin-gonic/gin"

func getRoutes(r *gin.Engine) {
	users := r.Group("users")
	{
		users.POST("/", userHandler.Create)
		users.DELETE("/", authMiddleware.Authenticate, userHandler.Delete)
		users.PUT("/password", authMiddleware.Authenticate, userHandler.UpdatePassword)
	}

	agents := r.Group("agents")
	{
		agents.Use(authMiddleware.Authenticate)
		agents.GET("/", agentHandler.Gets)
		agents.POST("/", agentHandler.Create)
		agents.PUT("/:id", agentHandler.Update)
		agents.DELETE("/:id", agentHandler.Delete)
		agents.GET("/:id/policies", agentHandler.GetPolicies)
		agents.PUT("/:id/policies", agentHandler.UpdatePolicies)
	}

	policies := r.Group("policies")
	{
		policies.Use(authMiddleware.Authenticate)
		policies.GET("/", policyHandler.Gets)
		policies.POST("/", policyHandler.Create)
		policies.PUT("/:id", policyHandler.Update)
		policies.DELETE("/:id", policyHandler.Delete)
		policies.GET("/:id/agents", policyHandler.GetAgents)
	}

	auth := r.Group("auth")
	{
		auth.GET("/user_id", authHandler.GetUserID)
		auth.POST("/signin", authHandler.Signin)
		auth.DELETE("/signout", authHandler.Signout)
	}
}
