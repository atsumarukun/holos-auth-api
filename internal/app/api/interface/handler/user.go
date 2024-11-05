package handler

import (
	"context"
	"holos-auth-api/internal/app/api/interface/request"
	"holos-auth-api/internal/app/api/interface/response"
	"holos-auth-api/internal/app/api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) UserHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}

func (uh *userHandler) Create(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()

	dto, err := uh.userUsecase.Create(ctx, req.Name, req.Password, req.ConfirmPassword)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response.NewUserResponse(dto.Name, dto.CreatedAt, dto.UpdatedAt))
}

func (uh *userHandler) Update(c *gin.Context) {
	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()

	dto, err := uh.userUsecase.Update(ctx, c.Param("name"), req.CurrentPassword, req.NewPassword, req.ConfirmNewPassword)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response.NewUserResponse(dto.Name, dto.CreatedAt, dto.UpdatedAt))
}

func (uh *userHandler) Delete(c *gin.Context) {
	var req request.DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()

	if err := uh.userUsecase.Delete(ctx, c.Param("name"), req.Password); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
