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

	dto, err := h.policyUsecase.Create(ctx, userID, req.Name, req.Service, req.Path, req.Methods)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.JSON(http.StatusCreated, builder.ToPolicyResponse(dto))
}

func (h *policyHandler) Update(c *gin.Context) {
	var req request.UpdatePolicyRequest
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

	dto, err := h.policyUsecase.Update(ctx, id, userID, req.Name, req.Service, req.Path, req.Methods)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.JSON(http.StatusOK, builder.ToPolicyResponse(dto))
}

func (h *policyHandler) Delete(c *gin.Context) {
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

	if err := h.policyUsecase.Delete(ctx, id, userID); err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *policyHandler) Gets(c *gin.Context) {
	userID, err := parameter.GetContextParameter[uuid.UUID](c, "userID")
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	ctx := c.Request.Context()

	dtos, err := h.policyUsecase.Gets(ctx, userID)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.JSON(http.StatusOK, builder.ToPolicyResponses(dtos))
}

func (h *policyHandler) UpdateAgents(c *gin.Context) {
	var req request.UpdatePolicyAgentsRequest
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

	dtos, err := h.policyUsecase.UpdateAgents(ctx, id, userID, req.AgentIDs)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.JSON(http.StatusOK, builder.ToAgentResponses(dtos))
}

func (h *policyHandler) GetAgents(c *gin.Context) {
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

	dtos, err := h.policyUsecase.GetAgents(ctx, id, userID)
	if err != nil {
		status := errors.HandleError(err)
		log.Println(status.Message())
		c.String(status.Code(), status.Message())
		return
	}

	c.JSON(http.StatusOK, builder.ToAgentResponses(dtos))
}
