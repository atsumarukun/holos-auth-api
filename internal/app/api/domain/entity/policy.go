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
	ErrPolicyNameTooShort    = status.Error(http.StatusBadRequest, "policy name must be 3 characters or more")
	ErrPolicyNameTooLong     = status.Error(http.StatusBadRequest, "policy name must be 255 characters or less")
	ErrInvalidPolicyName     = status.Error(http.StatusBadRequest, "invalid policy name")
	ErrInvalidPolicyEffect   = status.Error(http.StatusBadRequest, "invalid policy effect")
	ErrInvalidPolicyService  = status.Error(http.StatusBadRequest, "invalid policy service")
	ErrRequiredPolicyPath    = status.Error(http.StatusBadRequest, "policy path is required")
	ErrPolicyPathTooLong     = status.Error(http.StatusBadRequest, "policy path must be 255 characters or less")
	ErrInvalidPolicyPath     = status.Error(http.StatusBadRequest, "invalid policy path")
	ErrRequiredPolicyMethods = status.Error(http.StatusBadRequest, "policy methods is required")
	ErrInvalidPolicyMethods  = status.Error(http.StatusBadRequest, "invalid policy methods")
)

type Policy struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Effect    string
	Service   string
	Path      string
	Methods   []string
	Agents    []uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPolicy(userID uuid.UUID, name string, effect string, service string, path string, methods []string) (*Policy, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	policy := &Policy{
		ID:     id,
		UserID: userID,
		Agents: []uuid.UUID{},
	}

	if err := policy.SetName(name); err != nil {
		return nil, err
	}
	if err := policy.SetEffect(effect); err != nil {
		return nil, err
	}
	if err := policy.SetService(service); err != nil {
		return nil, err
	}
	if err := policy.SetPath(path); err != nil {
		return nil, err
	}
	if err := policy.SetMethods(methods); err != nil {
		return nil, err
	}

	now := time.Now()
	policy.CreatedAt = now
	policy.UpdatedAt = now

	return policy, nil
}

func RestorePolicy(id uuid.UUID, userID uuid.UUID, name string, effect string, service string, path string, methods []string, agents []uuid.UUID, createdAt time.Time, updatedAt time.Time) *Policy {
	return &Policy{
		ID:        id,
		UserID:    userID,
		Name:      name,
		Effect:    effect,
		Service:   service,
		Path:      path,
		Methods:   methods,
		Agents:    agents,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (p *Policy) SetName(name string) error {
	if len(name) < 3 {
		return ErrPolicyNameTooShort
	}
	if 255 < len(name) {
		return ErrPolicyNameTooLong
	}
	matched, err := regexp.MatchString(`^[A-Za-z0-9_]*$`, name)
	if err != nil {
		return err
	}
	if !matched {
		return ErrInvalidPolicyName
	}
	p.Name = name
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Policy) SetEffect(effect string) error {
	if !slices.Contains([]string{"ALLOW", "DENY"}, effect) {
		return ErrInvalidPolicyEffect
	}

	p.Effect = effect
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Policy) SetService(service string) error {
	if !slices.Contains([]string{"STORAGE", "CONTENT"}, service) {
		return ErrInvalidPolicyService
	}

	p.Service = service
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Policy) SetPath(path string) error {
	if len(path) == 0 {
		return ErrRequiredPolicyPath
	}
	if 255 < len(path) {
		return ErrPolicyPathTooLong
	}
	if path[0] != '/' || path[len(path)-1:] == "/" && 1 < len(path) {
		return ErrInvalidPolicyPath
	}
	matched, err := regexp.MatchString(`^[a-z\-/]*$`, path)
	if err != nil {
		return err
	}
	if !matched {
		return ErrInvalidPolicyPath
	}

	p.Path = path
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Policy) SetMethods(methods []string) error {
	if len(methods) == 0 {
		return ErrRequiredPolicyMethods
	}
	for _, v := range methods {
		if !slices.Contains([]string{"GET", "POST", "PUT", "DELETE"}, v) {
			return ErrInvalidPolicyMethods
		}
	}
	slices.Sort(methods)

	p.Methods = slices.Compact(methods)
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Policy) SetAgents(agents []*Agent) {
	ids := make([]uuid.UUID, len(agents))
	for i, agent := range agents {
		ids[i] = agent.ID
	}
	p.Agents = slices.Compact(ids)
	p.UpdatedAt = time.Now()
}
