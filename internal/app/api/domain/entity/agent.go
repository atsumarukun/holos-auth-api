package entity

import (
	"holos-auth-api/internal/app/api/pkg/status"
	"net/http"
	"regexp"
	"slices"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAgentNameTooShort = status.Error(http.StatusBadRequest, "agent name must be 3 characters or more")
	ErrAgentNameTooLong  = status.Error(http.StatusBadRequest, "agent name must be 255 characters or less")
	ErrInvalidAgentName  = status.Error(http.StatusBadRequest, "invalid agent name")
)

type Agent struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Policies  []uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAgent(userID uuid.UUID, name string) (*Agent, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	agent := &Agent{
		ID:       id,
		UserID:   userID,
		Policies: []uuid.UUID{},
	}

	if err := agent.SetName(name); err != nil {
		return nil, err
	}

	now := time.Now()
	agent.CreatedAt = now
	agent.UpdatedAt = now

	return agent, nil
}

func RestoreAgent(id uuid.UUID, userID uuid.UUID, name string, policies []uuid.UUID, createdAt time.Time, updatedAt time.Time) *Agent {
	return &Agent{
		ID:        id,
		UserID:    userID,
		Name:      name,
		Policies:  policies,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (a *Agent) SetName(name string) error {
	if len(name) < 3 {
		return ErrAgentNameTooShort
	}
	if 255 < len(name) {
		return ErrAgentNameTooLong
	}
	matched, err := regexp.MatchString(`^[A-Za-z0-9_]*$`, name)
	if err != nil {
		return status.Error(http.StatusInternalServerError, err.Error())
	}
	if !matched {
		return ErrInvalidAgentName
	}
	a.Name = name
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Agent) SetPolicies(policies []*Policy) {
	ids := make([]uuid.UUID, len(policies))
	for i, policy := range policies {
		ids[i] = policy.ID
	}
	a.Policies = slices.Compact(ids)
	a.UpdatedAt = time.Now()
}
