package handler

import (
	"holos-auth-api/internal/app/api/interface/pkg/errors"
	"holos-auth-api/internal/app/api/interface/pkg/parameter"
	"holos-auth-api/internal/app/api/interface/request"
	"holos-auth-api/internal/app/api/interface/response"
	"holos-auth-api/internal/app/api/usecase"
	"log"
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

func (h *userHandler) Create(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status := errors.StatusBadRequest
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	ctx := c.Request.Context()

	dto, err := h.userUsecase.Create(ctx, req.Name, req.Password, req.ConfirmPassword)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.JSON(http.StatusCreated, response.NewUserResponse(dto.Name, dto.CreatedAt, dto.UpdatedAt))
}

func (h *userHandler) UpdateName(c *gin.Context) {
	var req request.UpdateUserNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status := errors.StatusBadRequest
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	id, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	ctx := c.Request.Context()

	dto, err := h.userUsecase.UpdateName(ctx, id, req.Name)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.JSON(http.StatusOK, response.NewUserResponse(dto.Name, dto.CreatedAt, dto.UpdatedAt))
}

func (h *userHandler) UpdatePassword(c *gin.Context) {
	var req request.UpdateUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status := errors.StatusBadRequest
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	id, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	ctx := c.Request.Context()

	dto, err := h.userUsecase.UpdatePassword(ctx, id, req.CurrentPassword, req.NewPassword, req.ConfirmNewPassword)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.JSON(http.StatusOK, response.NewUserResponse(dto.Name, dto.CreatedAt, dto.UpdatedAt))
}

func (h *userHandler) Delete(c *gin.Context) {
	var req request.DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status := errors.StatusBadRequest
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	id, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	ctx := c.Request.Context()

	if err := h.userUsecase.Delete(ctx, id, req.Password); err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.Status(http.StatusNoContent)
}
