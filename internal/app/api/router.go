package api

import "github.com/gin-gonic/gin"

func registerRouter(r *gin.Engine) {
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
		agents.POST("/:id/token", agentHandler.GenerateToken)
		agents.DELETE("/:id/token", agentHandler.DeleteToken)
	}

	policies := r.Group("policies")
	{
		policies.Use(authMiddleware.Authenticate)
		policies.GET("/", policyHandler.Gets)
		policies.POST("/", policyHandler.Create)
		policies.PUT("/:id", policyHandler.Update)
		policies.DELETE("/:id", policyHandler.Delete)
		policies.GET("/:id/agents", policyHandler.GetAgents)
		policies.PUT("/:id/agents", policyHandler.UpdateAgents)
	}

	auth := r.Group("auth")
	{
		auth.GET("/authorization", authHandler.Authorize)
		auth.POST("/signin", authHandler.Signin)
		auth.DELETE("/signout", authHandler.Signout)
	}
}
