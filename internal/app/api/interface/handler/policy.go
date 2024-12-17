package handler

import (
	"holos-auth-api/internal/app/api/interface/pkg/parameter"
	"holos-auth-api/internal/app/api/interface/request"
	"holos-auth-api/internal/app/api/interface/response"
	"holos-auth-api/internal/app/api/pkg/status"
	"holos-auth-api/internal/app/api/usecase"
	"holos-auth-api/internal/app/api/usecase/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PolicyHandler interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	Gets(*gin.Context)
	UpdateAgents(*gin.Context)
	GetAgents(*gin.Context)
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
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	ctx := c.Request.Context()

	dto, err := h.policyUsecase.Create(ctx, userID, req.Name, req.Service, req.Path, req.Methods)
	if err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	c.JSON(http.StatusCreated, h.convertToResponse(dto))
}

func (h *policyHandler) Update(c *gin.Context) {
	var req request.UpdatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	ctx := c.Request.Context()

	dto, err := h.policyUsecase.Update(ctx, id, userID, req.Name, req.Service, req.Path, req.Methods)
	if err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	c.JSON(http.StatusOK, h.convertToResponse(dto))
}

func (h *policyHandler) Delete(c *gin.Context) {
	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	ctx := c.Request.Context()

	if err := h.policyUsecase.Delete(ctx, id, userID); err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *policyHandler) Gets(c *gin.Context) {
	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	ctx := c.Request.Context()

	dtos, err := h.policyUsecase.Gets(ctx, userID)
	if err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	c.JSON(http.StatusOK, h.convertToResponses(dtos))
}

func (h *policyHandler) UpdateAgents(c *gin.Context) {
	var req request.UpdatePolicyAgentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	ctx := c.Request.Context()

	dtos, err := h.policyUsecase.UpdateAgents(ctx, id, userID, req.AgentIDs)
	if err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	responses := make([]*response.AgentResponse, len(dtos))
	for i, dto := range dtos {
		responses[i] = response.NewAgentResponse(dto.ID, dto.Name, dto.CreatedAt, dto.UpdatedAt)
	}
	c.JSON(http.StatusOK, responses)
}

func (h *policyHandler) GetAgents(c *gin.Context) {
	id, err := parameter.GetPathParameter[uuid.UUID](c, "id")
	if err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	ctx := c.Request.Context()

	dtos, err := h.policyUsecase.GetAgents(ctx, id, userID)
	if err != nil {
		e := status.FromError(err)
		c.String(e.Code(), e.Message())
		return
	}

	responses := make([]*response.AgentResponse, len(dtos))
	for i, dto := range dtos {
		responses[i] = response.NewAgentResponse(dto.ID, dto.Name, dto.CreatedAt, dto.UpdatedAt)
	}
	c.JSON(http.StatusOK, responses)
}

func (h *policyHandler) convertToResponse(policy *dto.PolicyDTO) *response.PolicyResponse {
	return response.NewPolicyResponse(policy.ID, policy.Name, policy.Service, policy.Path, policy.Methods, policy.CreatedAt, policy.UpdatedAt)
}

func (h *policyHandler) convertToResponses(policies []*dto.PolicyDTO) []*response.PolicyResponse {
	responses := make([]*response.PolicyResponse, len(policies))
	for i, policy := range policies {
		responses[i] = h.convertToResponse(policy)
	}
	return responses
}
