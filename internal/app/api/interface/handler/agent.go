package handler

import (
	"holos-auth-api/internal/app/api/interface/pkg/parameter"
	"holos-auth-api/internal/app/api/interface/request"
	"holos-auth-api/internal/app/api/interface/response"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"holos-auth-api/internal/app/api/usecase"
	"holos-auth-api/internal/app/api/usecase/dto"
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
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	ctx := c.Request.Context()

	dto, err := h.agentUsecase.Create(ctx, userID, req.Name)
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	c.JSON(http.StatusCreated, h.convertToResponse(dto))
}

func (h *agentHandler) Update(c *gin.Context) {
	var req request.UpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	ctx := c.Request.Context()

	dto, err := h.agentUsecase.Update(ctx, id, userID, req.Name)
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	c.JSON(http.StatusOK, h.convertToResponse(dto))
}

func (h *agentHandler) Delete(c *gin.Context) {
	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	ctx := c.Request.Context()

	if err := h.agentUsecase.Delete(ctx, id, userID); err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *agentHandler) Gets(c *gin.Context) {
	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	ctx := c.Request.Context()

	dtos, err := h.agentUsecase.Gets(ctx, userID)
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	c.JSON(http.StatusOK, h.convertToResponses(dtos))
}

func (h *agentHandler) UpdatePolicies(c *gin.Context) {
	var req request.UpdateAgentPoliciesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	ctx := c.Request.Context()

	dtos, err := h.agentUsecase.UpdatePolicies(ctx, id, userID, req.PolicyIDs)
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	responses := make([]*response.PolicyResponse, len(dtos))
	for i, dto := range dtos {
		responses[i] = response.NewPolicyResponse(dto.ID, dto.Name, dto.Service, dto.Path, dto.Methods, dto.CreatedAt, dto.UpdatedAt)
	}
	c.JSON(http.StatusOK, responses)
}

func (h *agentHandler) GetPolicies(c *gin.Context) {
	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	ctx := c.Request.Context()

	dtos, err := h.agentUsecase.GetPolicies(ctx, id, userID)
	if err != nil {
		e := apierr.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	responses := make([]*response.PolicyResponse, len(dtos))
	for i, dto := range dtos {
		responses[i] = response.NewPolicyResponse(dto.ID, dto.Name, dto.Service, dto.Path, dto.Methods, dto.CreatedAt, dto.UpdatedAt)
	}
	c.JSON(http.StatusOK, responses)
}

func (h *agentHandler) convertToResponse(agent *dto.AgentDTO) *response.AgentResponse {
	return response.NewAgentResponse(agent.ID, agent.Name, agent.CreatedAt, agent.UpdatedAt)
}

func (h *agentHandler) convertToResponses(agents []*dto.AgentDTO) []*response.AgentResponse {
	responses := make([]*response.AgentResponse, len(agents))
	for i, agent := range agents {
		responses[i] = h.convertToResponse(agent)
	}
	return responses
}
