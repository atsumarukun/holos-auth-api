package api

import (
	"holos-auth-api/internal/app/api/domain/service"
	dbrepository "holos-auth-api/internal/app/api/infrastructure/db"
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
	transactionObject := dbrepository.NewSqlxTransactionObject(db)

	userDBRepository := dbrepository.NewUserDBRepository(db)
	userTokenDBRepository := dbrepository.NewUserTokenDBRepository(db)
	agentDBRepository := dbrepository.NewAgentDBRepository(db)
	policyDBRepository := dbrepository.NewPolicyDBRepository(db)

	userService := service.NewUserService(userDBRepository)

	userUsecase := usecase.NewUserUsecase(transactionObject, userDBRepository, userService)
	agentUsecase := usecase.NewAgentUsecase(transactionObject, agentDBRepository)
	policyUsecase := usecase.NewPolicyUsecase(transactionObject, policyDBRepository)
	authUsecase := usecase.NewAuthUsecase(transactionObject, userDBRepository, userTokenDBRepository)

	authMiddleware = middleware.NewAuthMiddleware(authUsecase)

	userHandler = handler.NewUserHandler(userUsecase)
	agentHandler = handler.NewAgentHandler(agentUsecase)
	policyHandler = handler.NewPolicyHandler(policyUsecase)
	authHandler = handler.NewAuthHandler(authUsecase)
}
