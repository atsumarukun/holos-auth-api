package api

import (
	"holos-auth-api/internal/app/api/domain/service"
	"holos-auth-api/internal/app/api/infrastructure/database"
	"holos-auth-api/internal/app/api/interface/handler"
	"holos-auth-api/internal/app/api/interface/middleware"
	"holos-auth-api/internal/app/api/usecase"

	"github.com/jmoiron/sqlx"
)

var (
	authMiddleware middleware.AuthMiddleware

	userHandler   handler.UserHandler
	agentHandler  handler.AgentHandler
	policyHandler handler.PolicyHandler
	authHandler   handler.AuthHandler
)

func inject(db *sqlx.DB) {
	transactionObject := database.NewDBTransactionObject(db)

	userDBRepository := database.NewUserDBRepository(db)
	userTokenDBRepository := database.NewUserTokenDBRepository(db)
	agentDBRepository := database.NewAgentDBRepository(db)
	policyDBRepository := database.NewPolicyDBRepository(db)

	userService := service.NewUserService(userDBRepository)

	userUsecase := usecase.NewUserUsecase(transactionObject, userDBRepository, userService)
	agentUsecase := usecase.NewAgentUsecase(transactionObject, agentDBRepository, policyDBRepository)
	policyUsecase := usecase.NewPolicyUsecase(transactionObject, policyDBRepository, agentDBRepository)
	authUsecase := usecase.NewAuthUsecase(transactionObject, userDBRepository, userTokenDBRepository)

	authMiddleware = middleware.NewAuthMiddleware(authUsecase)

	userHandler = handler.NewUserHandler(userUsecase)
	agentHandler = handler.NewAgentHandler(agentUsecase)
	policyHandler = handler.NewPolicyHandler(policyUsecase)
	authHandler = handler.NewAuthHandler(authUsecase)
}
