package handler

import (
	"context"
	"holos-auth-api/internal/app/api/interface/pkg/parameter"
	"holos-auth-api/internal/app/api/interface/request"
	"holos-auth-api/internal/app/api/interface/response"
	"holos-auth-api/internal/app/api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PolicyHandler interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

type policyHandler struct {
	policyUsecase usecase.PolicyUsecase
}

func NewPolicyHandler(policyUsecase usecase.PolicyUsecase) PolicyHandler {
	return &policyHandler{
		policyUsecase: policyUsecase,
	}
}

func (h *policyHandler) Create(c *gin.Context) {
	var req request.CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		c.String(err.Error())
		return
	}

	ctx := context.Background()

	dto, err := h.policyUsecase.Create(ctx, userID, req.Name, req.Service, req.Path, req.AllowedMethods)
	if err != nil {
		c.String(err.Error())
		return
	}

	c.JSON(http.StatusCreated, response.NewPolicyResponse(dto.ID, dto.Name, dto.Service, dto.Path, dto.AllowedMethods, dto.CreatedAt, dto.UpdatedAt))
}

func (h *policyHandler) Update(c *gin.Context) {
	var req request.UpdatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		c.String(err.Error())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		c.String(err.Error())
		return
	}

	ctx := context.Background()

	dto, err := h.policyUsecase.Update(ctx, id, userID, req.Name, req.Service, req.Path, req.AllowedMethods)
	if err != nil {
		c.String(err.Error())
		return
	}

	c.JSON(http.StatusOK, response.NewPolicyResponse(dto.ID, dto.Name, dto.Service, dto.Path, dto.AllowedMethods, dto.CreatedAt, dto.UpdatedAt))
}

func (h *policyHandler) Delete(c *gin.Context) {
	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		c.String(err.Error())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		c.String(err.Error())
		return
	}

	ctx := context.Background()

	if err := h.policyUsecase.Delete(ctx, id, userID); err != nil {
		c.String(err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
