package api

import (
	"holos-auth-api/internal/app/api/domain/service"
	"holos-auth-api/internal/app/api/infrastructure"
	"holos-auth-api/internal/app/api/interface/handler"
	"holos-auth-api/internal/app/api/usecase"

	"github.com/jmoiron/sqlx"
)

var (
	userHandler handler.UserHandler
	authHandler handler.AuthHandler
)

func inject(db *sqlx.DB) {
	transactionObject := infrastructure.NewSqlxTransactionObject(db)

	userInfrastructure := infrastructure.NewUserInfrastructure(db)
	userTokenInfrastructure := infrastructure.NewUserTokenInfrastructure(db)

	userService := service.NewUserService(userInfrastructure)

	userUsecase := usecase.NewUserUsecase(transactionObject, userInfrastructure, userService)
	authUsecase := usecase.NewAuthUsecase(transactionObject, userInfrastructure, userTokenInfrastructure)

	userHandler = handler.NewUserHandler(userUsecase)
	authHandler = handler.NewAuthHandler(authUsecase)
}
