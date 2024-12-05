package entity

import (
	"holos-auth-api/internal/app/api/pkg/apierr"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPermissionEffect = apierr.NewApiError(http.StatusBadRequest, "invalid permission effect")
)

type Permission struct {
	AgentID   uuid.UUID
	PolicyID  uuid.UUID
	Effect    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPermission(agentID uuid.UUID, policyID uuid.UUID, effect string) (*Permission, apierr.ApiError) {
	if !slices.Contains([]string{"ALLOW", "DENY"}, effect) {
		return nil, ErrInvalidPermissionEffect
	}

	now := time.Now()
	return &Permission{
		AgentID:   agentID,
		PolicyID:  policyID,
		Effect:    effect,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
