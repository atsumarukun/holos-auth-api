package entity

import (
	"holos-auth-api/internal/app/api/pkg/apierr"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAgentNameTooShort = apierr.NewApiError(http.StatusBadRequest, "agent name must be 3 characters or more")
	ErrAgentNameTooLong  = apierr.NewApiError(http.StatusBadRequest, "agent name must be 255 characters or less")
	ErrInvalidAgentName  = apierr.NewApiError(http.StatusBadRequest, "invalid agent name")
)

type Agent struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Permissions []*Permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewAgent(userID uuid.UUID, name string) (*Agent, apierr.ApiError) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}

	agent := &Agent{
		ID:          id,
		UserID:      userID,
		Permissions: []*Permission{},
	}

	if err := agent.SetName(name); err != nil {
		return nil, err
	}

	now := time.Now()
	agent.CreatedAt = now
	agent.UpdatedAt = now

	return agent, nil
}

func RestoreAgent(id uuid.UUID, userID uuid.UUID, name string, createdAt time.Time, updatedAt time.Time) *Agent {
	return &Agent{
		ID:        id,
		UserID:    userID,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (a *Agent) SetName(name string) apierr.ApiError {
	if len(name) < 3 {
		return ErrAgentNameTooShort
	}
	if 255 < len(name) {
		return ErrAgentNameTooLong
	}
	matched, err := regexp.MatchString(`^[A-Za-z0-9_]*$`, name)
	if err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	if !matched {
		return ErrInvalidAgentName
	}
	a.Name = name
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Agent) SetPermissions(policies []*Policy, effect string) apierr.ApiError {
	var err apierr.ApiError
	a.Permissions = make([]*Permission, len(policies))
	for i, policy := range policies {
		a.Permissions[i], err = NewPermission(a.ID, policy.ID, effect)
		if err != nil {
			return err
		}
	}
	return nil
}
