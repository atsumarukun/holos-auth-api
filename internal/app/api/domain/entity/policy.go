package entity

import (
	"holos-auth-api/internal/app/api/pkg/apierr"
	"net/http"
	"regexp"
	"slices"
	"time"

	"github.com/google/uuid"
)

var (
	services = []string{"STORAGE", "CONTENT"}
	methods  = []string{"GET", "POST", "PUT", "DELETE"}
)

var (
	ErrPolicyNameTooShort           = apierr.NewApiError(http.StatusBadRequest, "policy name must be 3 characters or more")
	ErrPolicyNameTooLong            = apierr.NewApiError(http.StatusBadRequest, "policy name must be 255 characters or less")
	ErrInvalidPolicyName            = apierr.NewApiError(http.StatusBadRequest, "invalid policy name")
	ErrInvalidPolicyService         = apierr.NewApiError(http.StatusBadRequest, "invalid policy service")
	ErrRequiredPolicyPath           = apierr.NewApiError(http.StatusBadRequest, "policy path is required")
	ErrPolicyPathTooLong            = apierr.NewApiError(http.StatusBadRequest, "policy path must be 255 characters or less")
	ErrInvalidPolicyPath            = apierr.NewApiError(http.StatusBadRequest, "invalid policy path")
	ErrRequiredPolicyAllowedMethods = apierr.NewApiError(http.StatusBadRequest, "policy allowed methods is required")
	ErrInvalidPolicyAllowedMethods  = apierr.NewApiError(http.StatusBadRequest, "invalid policy allowed methods")
)

type Policy struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Name           string
	Service        string
	Path           string
	AllowedMethods []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewPolicy(userID uuid.UUID, name string, service string, path string, allowedMethods []string) (*Policy, apierr.ApiError) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}

	policy := &Policy{
		ID:     id,
		UserID: userID,
	}

	if err := policy.SetName(name); err != nil {
		return nil, err
	}
	if err := policy.SetService(service); err != nil {
		return nil, err
	}
	if err := policy.SetPath(path); err != nil {
		return nil, err
	}
	if err := policy.SetAllowedMethods(allowedMethods); err != nil {
		return nil, err
	}

	now := time.Now()
	policy.CreatedAt = now
	policy.UpdatedAt = now

	return policy, nil
}

func RestorePolicy(id uuid.UUID, userID uuid.UUID, name string, service string, path string, allowedMethods []string, createdAt time.Time, updatedAt time.Time) *Policy {
	return &Policy{
		ID:             id,
		UserID:         userID,
		Name:           name,
		Service:        service,
		Path:           path,
		AllowedMethods: allowedMethods,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

func (p *Policy) SetName(name string) apierr.ApiError {
	if len(name) < 3 {
		return ErrPolicyNameTooShort
	}
	if 255 < len(name) {
		return ErrPolicyNameTooLong
	}
	matched, err := regexp.MatchString(`^[A-Za-z0-9_]*$`, name)
	if err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	if !matched {
		return ErrInvalidPolicyName
	}
	p.Name = name
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Policy) SetService(service string) apierr.ApiError {
	if !slices.Contains(services, service) {
		return ErrInvalidPolicyService
	}

	p.Service = service
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Policy) SetPath(path string) apierr.ApiError {
	if len(path) == 0 {
		return ErrRequiredPolicyPath
	}
	if 255 < len(path) {
		return ErrPolicyPathTooLong
	}
	if path[0] != '/' {
		return ErrInvalidPolicyPath
	}
	if path[len(path)-1:] == "/" && 1 < len(path) {
		path = path[:len(path)-1]
	}

	p.Path = path
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Policy) SetAllowedMethods(allowedMethods []string) apierr.ApiError {
	if len(allowedMethods) == 0 {
		return ErrRequiredPolicyAllowedMethods
	}
	for _, v := range allowedMethods {
		if !slices.Contains(methods, v) {
			return ErrInvalidPolicyAllowedMethods
		}
	}
	slices.Sort(allowedMethods)

	p.AllowedMethods = slices.Compact(allowedMethods)
	p.UpdatedAt = time.Now()
	return nil
}
