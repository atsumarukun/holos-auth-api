package handler

import (
	"context"
	"holos-auth-api/internal/app/api/interface/request"
	"holos-auth-api/internal/app/api/interface/response"
	"holos-auth-api/internal/app/api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler interface {
	Create(*gin.Context)
	UpdateName(*gin.Context)
	UpdatePassword(*gin.Context)
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
		c.String(err.Error())
		return
	}

	c.JSON(http.StatusOK, response.NewUserResponse(dto.Name, dto.CreatedAt, dto.UpdatedAt))
}

func (uh *userHandler) UpdateName(c *gin.Context) {
	var req request.UpdateUserNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.String(http.StatusInternalServerError, "context does not have user id")
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		c.String(http.StatusInternalServerError, "Invalid user id type")
	}

	ctx := context.Background()

	dto, err := uh.userUsecase.UpdateName(ctx, id, req.Name)
	if err != nil {
		c.String(err.Error())
		return
	}

	c.JSON(http.StatusOK, response.NewUserResponse(dto.Name, dto.CreatedAt, dto.UpdatedAt))
}

func (uh *userHandler) UpdatePassword(c *gin.Context) {
	var req request.UpdateUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.String(http.StatusInternalServerError, "context does not have user id")
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		c.String(http.StatusInternalServerError, "Invalid user id type")
	}

	ctx := context.Background()

	dto, err := uh.userUsecase.UpdatePassword(ctx, id, req.CurrentPassword, req.NewPassword, req.ConfirmNewPassword)
	if err != nil {
		c.String(err.Error())
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

	userID, exists := c.Get("userID")
	if !exists {
		c.String(http.StatusInternalServerError, "context does not have user id")
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		c.String(http.StatusInternalServerError, "Invalid user id type")
	}

	ctx := context.Background()

	if err := uh.userUsecase.Delete(ctx, id, req.Password); err != nil {
		c.String(err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
