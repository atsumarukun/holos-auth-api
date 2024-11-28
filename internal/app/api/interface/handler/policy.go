package handler

import (
	"context"
	"holos-auth-api/internal/app/api/interface/pkg/parameter"
	"holos-auth-api/internal/app/api/interface/request"
	"holos-auth-api/internal/app/api/interface/response"
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

	c.JSON(http.StatusOK, h.convertToResponse(dto))
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

func (h *policyHandler) Gets(c *gin.Context) {
	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		c.String(err.Error())
		return
	}

	ctx := context.Background()

	dtos, err := h.policyUsecase.Gets(ctx, userID)
	if err != nil {
		c.String(err.Error())
		return
	}

	c.JSON(http.StatusOK, h.convertToResponses(dtos))
}

func (h *policyHandler) convertToResponse(policy *dto.PolicyDTO) *response.PolicyResponse {
	return response.NewPolicyResponse(policy.ID, policy.Name, policy.Service, policy.Path, policy.AllowedMethods, policy.CreatedAt, policy.UpdatedAt)
}

func (h *policyHandler) convertToResponses(policies []*dto.PolicyDTO) []*response.PolicyResponse {
	responses := make([]*response.PolicyResponse, len(policies))
	for i, policy := range policies {
		responses[i] = h.convertToResponse(policy)
	}
	return responses
}
