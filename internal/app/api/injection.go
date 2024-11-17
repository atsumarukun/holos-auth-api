package api

import (
	"holos-auth-api/internal/app/api/domain/service"
	"holos-auth-api/internal/app/api/infrastructure"
	"holos-auth-api/internal/app/api/interface/handler"
	"holos-auth-api/internal/app/api/interface/middleware"
	"holos-auth-api/internal/app/api/usecase"

	"github.com/jmoiron/sqlx"
)

var (
	authMiddleware middleware.AuthMiddleware

	userHandler  handler.UserHandler
	agentHandler handler.AgentHandler
	authHandler  handler.AuthHandler
)

func inject(db *sqlx.DB) {
	transactionObject := infrastructure.NewSqlxTransactionObject(db)

	userInfrastructure := infrastructure.NewUserInfrastructure(db)
	userTokenInfrastructure := infrastructure.NewUserTokenInfrastructure(db)
	agentInfrastructure := infrastructure.NewAgentInfrastructure(db)

	userService := service.NewUserService(userInfrastructure)
	agentService := service.NewAgentService(agentInfrastructure)

	userUsecase := usecase.NewUserUsecase(transactionObject, userInfrastructure, userService)
	agentUsecase := usecase.NewAgentUsecase(transactionObject, agentInfrastructure, agentService)
	authUsecase := usecase.NewAuthUsecase(transactionObject, userInfrastructure, userTokenInfrastructure)

	authMiddleware = middleware.NewAuthMiddleware(authUsecase)

	userHandler = handler.NewUserHandler(userUsecase)
	agentHandler = handler.NewAgentHandler(agentUsecase)
	authHandler = handler.NewAuthHandler(authUsecase)
}
