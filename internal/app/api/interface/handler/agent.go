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

type AgentHandler interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

type agentHandler struct {
	agentUsecase usecase.AgentUsecase
}

func NewAgentHandler(agentUsecase usecase.AgentUsecase) AgentHandler {
	return &agentHandler{
		agentUsecase: agentUsecase,
	}
}

func (h *agentHandler) Create(c *gin.Context) {
	var req request.CreateAgentRequest
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

	dto, err := h.agentUsecase.Create(ctx, userID, req.Name)
	if err != nil {
		c.String(err.Error())
		return
	}

	c.JSON(http.StatusOK, response.NewAgentResponse(dto.ID, dto.Name, dto.CreatedAt, dto.UpdatedAt))
}

func (h *agentHandler) Update(c *gin.Context) {
	var req request.UpdateAgentRequest
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

	dto, err := h.agentUsecase.Update(ctx, id, userID, req.Name)
	if err != nil {
		c.String(err.Error())
		return
	}

	c.JSON(http.StatusOK, response.NewAgentResponse(dto.ID, dto.Name, dto.CreatedAt, dto.UpdatedAt))
}

func (h *agentHandler) Delete(c *gin.Context) {
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

	if err := h.agentUsecase.Delete(ctx, id, userID); err != nil {
		c.String(err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
