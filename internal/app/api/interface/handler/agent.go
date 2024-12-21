package handler

import (
	"holos-auth-api/internal/app/api/interface/builder"
	"holos-auth-api/internal/app/api/interface/pkg/errors"
	"holos-auth-api/internal/app/api/interface/pkg/parameter"
	"holos-auth-api/internal/app/api/interface/request"
	"holos-auth-api/internal/app/api/usecase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AgentHandler interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	Gets(*gin.Context)
	UpdatePolicies(*gin.Context)
	GetPolicies(*gin.Context)
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
		status := errors.StatusBadRequest
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	ctx := c.Request.Context()

	dto, err := h.agentUsecase.Create(ctx, userID, req.Name)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.JSON(http.StatusCreated, builder.ToAgentResponse(dto))
}

func (h *agentHandler) Update(c *gin.Context) {
	var req request.UpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status := errors.StatusBadRequest
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	ctx := c.Request.Context()

	dto, err := h.agentUsecase.Update(ctx, id, userID, req.Name)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.JSON(http.StatusOK, builder.ToAgentResponse(dto))
}

func (h *agentHandler) Delete(c *gin.Context) {
	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	ctx := c.Request.Context()

	if err := h.agentUsecase.Delete(ctx, id, userID); err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *agentHandler) Gets(c *gin.Context) {
	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	ctx := c.Request.Context()

	dtos, err := h.agentUsecase.Gets(ctx, userID)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.JSON(http.StatusOK, builder.ToAgentResponses(dtos))
}

func (h *agentHandler) UpdatePolicies(c *gin.Context) {
	var req request.UpdateAgentPoliciesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status := errors.StatusBadRequest
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	ctx := c.Request.Context()

	dtos, err := h.agentUsecase.UpdatePolicies(ctx, id, userID, req.PolicyIDs)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.JSON(http.StatusOK, builder.ToPolicyResponses(dtos))
}

func (h *agentHandler) GetPolicies(c *gin.Context) {
	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	ctx := c.Request.Context()

	dtos, err := h.agentUsecase.GetPolicies(ctx, id, userID)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.JSON(http.StatusOK, builder.ToPolicyResponses(dtos))
}
