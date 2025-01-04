package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/usecase"
	"holos-auth-api/internal/app/api/usecase/dto"
	mockDomain "holos-auth-api/test/mock/domain"
	mockRepository "holos-auth-api/test/mock/domain/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestPolicy_Create(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                    string
		inputUserID             uuid.UUID
		inputName               string
		inputEffect             string
		inputService            string
		inputPath               string
		inputMethods            []string
		expectResult            *dto.PolicyDTO
		expectError             error
		setMockPolicyRepository func(context.Context, *mockRepository.MockPolicyRepository)
	}{
		{
			name:         "success",
			inputUserID:  policy.UserID,
			inputName:    "name",
			inputEffect:  "ALLOW",
			inputService: "STORAGE",
			inputPath:    "/",
			inputMethods: []string{"GET"},
			expectResult: &dto.PolicyDTO{ID: policy.ID, UserID: policy.UserID, Name: policy.Name, Effect: policy.Effect, Service: policy.Service, Path: policy.Path, Methods: policy.Methods, Agents: []uuid.UUID{}, CreatedAt: policy.CreatedAt, UpdatedAt: policy.UpdatedAt},
			expectError:  nil,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:                    "invalid name",
			inputUserID:             policy.UserID,
			inputName:               "なまえ",
			inputEffect:             "ALLOW",
			inputService:            "STORAGE",
			inputPath:               "/",
			inputMethods:            []string{"GET"},
			expectResult:            nil,
			expectError:             entity.ErrInvalidPolicyName,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {},
		},
		{
			name:                    "invalid effect",
			inputUserID:             policy.UserID,
			inputName:               "name",
			inputEffect:             "EFFECT",
			inputService:            "STORAGE",
			inputPath:               "/",
			inputMethods:            []string{"GET"},
			expectResult:            nil,
			expectError:             entity.ErrInvalidPolicyEffect,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {},
		},
		{
			name:                    "invalid service",
			inputUserID:             policy.UserID,
			inputName:               "name",
			inputEffect:             "ALLOW",
			inputService:            "SERVICE",
			inputPath:               "/",
			inputMethods:            []string{"GET"},
			expectResult:            nil,
			expectError:             entity.ErrInvalidPolicyService,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {},
		},
		{
			name:                    "invalid path",
			inputUserID:             policy.UserID,
			inputName:               "name",
			inputEffect:             "ALLOW",
			inputService:            "STORAGE",
			inputPath:               "path",
			inputMethods:            []string{"GET"},
			expectResult:            nil,
			expectError:             entity.ErrInvalidPolicyPath,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {},
		},
		{
			name:                    "invalid methods",
			inputUserID:             policy.UserID,
			inputName:               "name",
			inputEffect:             "ALLOW",
			inputService:            "STORAGE",
			inputPath:               "/",
			inputMethods:            []string{"PATCH"},
			expectResult:            nil,
			expectError:             entity.ErrInvalidPolicyMethods,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {},
		},
		{
			name:         "create error",
			inputUserID:  policy.UserID,
			inputName:    "name",
			inputEffect:  "ALLOW",
			inputService: "STORAGE",
			inputPath:    "/",
			inputMethods: []string{"GET"},
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					Create(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pr := mockRepository.NewMockPolicyRepository(ctrl)

			ctx := context.Background()

			tt.setMockPolicyRepository(ctx, pr)

			pu := usecase.NewPolicyUsecase(nil, pr, nil)
			result, err := pu.Create(ctx, tt.inputUserID, tt.inputName, tt.inputEffect, tt.inputService, tt.inputPath, tt.inputMethods)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.PolicyDTO{}, "ID", "CreatedAt", "UpdatedAt"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestPolicy_Update(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                     string
		inputID                  uuid.UUID
		inputUserID              uuid.UUID
		inputName                string
		inputEffect              string
		inputService             string
		inputPath                string
		inputMethods             []string
		expectResult             *dto.PolicyDTO
		expectError              error
		setMockTransactionObject func(context.Context, *mockDomain.MockTransactionObject)
		setMockPolicyRepository  func(context.Context, *mockRepository.MockPolicyRepository)
	}{
		{
			name:         "success",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			inputName:    "update",
			inputEffect:  "DENY",
			inputService: "CONTENT",
			inputPath:    "/path",
			inputMethods: []string{"PUT"},
			expectResult: &dto.PolicyDTO{ID: policy.ID, UserID: policy.UserID, Name: "update", Effect: "DENY", Service: "CONTENT", Path: "/path", Methods: []string{"PUT"}, Agents: []uuid.UUID{}, CreatedAt: policy.CreatedAt, UpdatedAt: policy.UpdatedAt},
			expectError:  nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
				pr.EXPECT().
					Update(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:         "invalid name",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			inputName:    "なまえ",
			inputEffect:  "DENY",
			inputService: "CONTENT",
			inputPath:    "/path",
			inputMethods: []string{"PUT"},
			expectResult: nil,
			expectError:  entity.ErrInvalidPolicyName,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
			},
		},
		{
			name:         "invalid effect",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			inputName:    "update",
			inputEffect:  "EFFECT",
			inputService: "CONTENT",
			inputPath:    "/path",
			inputMethods: []string{"PUT"},
			expectResult: nil,
			expectError:  entity.ErrInvalidPolicyEffect,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
			},
		},
		{
			name:         "invalid service",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			inputName:    "update",
			inputEffect:  "DENY",
			inputService: "SERVICE",
			inputPath:    "/path",
			inputMethods: []string{"PUT"},
			expectResult: nil,
			expectError:  entity.ErrInvalidPolicyService,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
			},
		},
		{
			name:         "invalid path",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			inputName:    "update",
			inputEffect:  "DENY",
			inputService: "CONTENT",
			inputPath:    "path",
			inputMethods: []string{"PUT"},
			expectResult: nil,
			expectError:  entity.ErrInvalidPolicyPath,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
			},
		},
		{
			name:         "invalid methods",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			inputName:    "update",
			inputEffect:  "DENY",
			inputService: "CONTENT",
			inputPath:    "/path",
			inputMethods: []string{"PATCH"},
			expectResult: nil,
			expectError:  entity.ErrInvalidPolicyMethods,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
			},
		},
		{
			name:         "policy not found",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			inputName:    "update",
			inputEffect:  "DENY",
			inputService: "CONTENT",
			inputPath:    "/path",
			inputMethods: []string{"PUT"},
			expectResult: nil,
			expectError:  usecase.ErrPolicyNotFound,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:         "find error",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			inputName:    "update",
			inputEffect:  "DENY",
			inputService: "CONTENT",
			inputPath:    "/path",
			inputMethods: []string{"PUT"},
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:         "update error",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			inputName:    "update",
			inputEffect:  "DENY",
			inputService: "CONTENT",
			inputPath:    "/path",
			inputMethods: []string{"PUT"},
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
				pr.EXPECT().
					Update(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			to := mockDomain.NewMockTransactionObject(ctrl)
			pr := mockRepository.NewMockPolicyRepository(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockPolicyRepository(ctx, pr)

			pu := usecase.NewPolicyUsecase(to, pr, nil)
			result, err := pu.Update(ctx, tt.inputID, tt.inputUserID, tt.inputName, tt.inputEffect, tt.inputService, tt.inputPath, tt.inputMethods)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.PolicyDTO{}, "UpdatedAt"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestPolicy_Delete(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                     string
		inputID                  uuid.UUID
		inputUserID              uuid.UUID
		expectError              error
		setMockTransactionObject func(context.Context, *mockDomain.MockTransactionObject)
		setMockPolicyRepository  func(context.Context, *mockRepository.MockPolicyRepository)
	}{
		{
			name:        "success",
			inputID:     policy.ID,
			inputUserID: policy.UserID,
			expectError: nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
				pr.EXPECT().
					Delete(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:        "policy not found",
			inputID:     policy.ID,
			inputUserID: policy.UserID,
			expectError: usecase.ErrPolicyNotFound,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:        "find error",
			inputID:     policy.ID,
			inputUserID: policy.UserID,
			expectError: sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:        "delete error",
			inputID:     policy.ID,
			inputUserID: policy.UserID,
			expectError: sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
				pr.EXPECT().
					Delete(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			to := mockDomain.NewMockTransactionObject(ctrl)
			pr := mockRepository.NewMockPolicyRepository(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockPolicyRepository(ctx, pr)

			pu := usecase.NewPolicyUsecase(to, pr, nil)
			if err := pu.Delete(ctx, tt.inputID, tt.inputUserID); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestPolicy_Gets(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                    string
		inputUserID             uuid.UUID
		expectResult            []*dto.PolicyDTO
		expectError             error
		setMockPolicyRepository func(context.Context, *mockRepository.MockPolicyRepository)
	}{
		{
			name:         "found",
			inputUserID:  policy.UserID,
			expectResult: []*dto.PolicyDTO{{ID: policy.ID, UserID: policy.UserID, Name: policy.Name, Effect: policy.Effect, Service: policy.Service, Path: policy.Path, Methods: policy.Methods, Agents: policy.Agents, CreatedAt: policy.CreatedAt, UpdatedAt: policy.UpdatedAt}},
			expectError:  nil,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByUserIDAndNotDeleted(ctx, policy.UserID).
					Return([]*entity.Policy{entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt)}, nil).
					Times(1)
			},
		},
		{
			name:         "not found",
			inputUserID:  policy.UserID,
			expectResult: []*dto.PolicyDTO{},
			expectError:  nil,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByUserIDAndNotDeleted(ctx, policy.UserID).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:         "find error",
			inputUserID:  policy.UserID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByUserIDAndNotDeleted(ctx, policy.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pr := mockRepository.NewMockPolicyRepository(ctrl)

			ctx := context.Background()

			tt.setMockPolicyRepository(ctx, pr)

			pu := usecase.NewPolicyUsecase(nil, pr, nil)
			result, err := pu.Gets(ctx, tt.inputUserID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestPolicy_UpdateAgents(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}
	agent, err := entity.NewAgent(policy.UserID, "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                     string
		inputID                  uuid.UUID
		inputUserID              uuid.UUID
		inputAgentIDs            []uuid.UUID
		expectResult             []*dto.AgentDTO
		expectError              error
		setMockTransactionObject func(context.Context, *mockDomain.MockTransactionObject)
		setMockPolicyRepository  func(context.Context, *mockRepository.MockPolicyRepository)
		setMockAgentRepository   func(context.Context, *mockRepository.MockAgentRepository)
	}{
		{
			name:          "success",
			inputID:       policy.ID,
			inputUserID:   policy.UserID,
			inputAgentIDs: []uuid.UUID{agent.ID},
			expectResult:  []*dto.AgentDTO{{ID: agent.ID, UserID: agent.UserID, Name: agent.Name, Policies: agent.Policies, CreatedAt: agent.CreatedAt, UpdatedAt: agent.UpdatedAt}},
			expectError:   nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
				pr.EXPECT().
					Update(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), agent.UserID).
					Return([]*entity.Agent{entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt)}, nil).
					Times(1)
			},
		},
		{
			name:          "policy not found",
			inputID:       policy.ID,
			inputUserID:   policy.UserID,
			inputAgentIDs: []uuid.UUID{agent.ID},
			expectResult:  nil,
			expectError:   usecase.ErrPolicyNotFound,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(nil, nil).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {},
		},
		{
			name:          "find policy error",
			inputID:       policy.ID,
			inputUserID:   policy.UserID,
			inputAgentIDs: []uuid.UUID{agent.ID},
			expectResult:  nil,
			expectError:   sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {},
		},
		{
			name:          "find agents error",
			inputID:       policy.ID,
			inputUserID:   policy.UserID,
			inputAgentIDs: []uuid.UUID{agent.ID},
			expectResult:  nil,
			expectError:   sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), agent.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:          "update policy error",
			inputID:       policy.ID,
			inputUserID:   policy.UserID,
			inputAgentIDs: []uuid.UUID{agent.ID},
			expectResult:  nil,
			expectError:   sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
				pr.EXPECT().
					Update(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), agent.UserID).
					Return([]*entity.Agent{entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt)}, nil).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			to := mockDomain.NewMockTransactionObject(ctrl)
			pr := mockRepository.NewMockPolicyRepository(ctrl)
			ar := mockRepository.NewMockAgentRepository(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockPolicyRepository(ctx, pr)
			tt.setMockAgentRepository(ctx, ar)

			pu := usecase.NewPolicyUsecase(to, pr, ar)
			result, err := pu.UpdateAgents(ctx, tt.inputID, tt.inputUserID, tt.inputAgentIDs)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestPolicy_GetAgents(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}
	agent, err := entity.NewAgent(policy.UserID, "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                     string
		inputID                  uuid.UUID
		inputUserID              uuid.UUID
		expectResult             []*dto.AgentDTO
		expectError              error
		setMockTransactionObject func(context.Context, *mockDomain.MockTransactionObject)
		setMockPolicyRepository  func(context.Context, *mockRepository.MockPolicyRepository)
		setMockAgentRepository   func(context.Context, *mockRepository.MockAgentRepository)
	}{
		{
			name:         "success",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			expectResult: []*dto.AgentDTO{{ID: agent.ID, UserID: agent.UserID, Name: agent.Name, Policies: agent.Policies, CreatedAt: agent.CreatedAt, UpdatedAt: agent.UpdatedAt}},
			expectError:  nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), agent.UserID).
					Return([]*entity.Agent{entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt)}, nil).
					Times(1)
			},
		},
		{
			name:         "policy not found",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			expectResult: nil,
			expectError:  usecase.ErrPolicyNotFound,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(nil, nil).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {},
		},
		{
			name:         "agents not found",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			expectResult: []*dto.AgentDTO{},
			expectError:  nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), agent.UserID).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:         "find policy error",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {},
		},
		{
			name:         "find agents error",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, policy.ID, policy.UserID).
					Return(entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt), nil).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), agent.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			to := mockDomain.NewMockTransactionObject(ctrl)
			pr := mockRepository.NewMockPolicyRepository(ctrl)
			ar := mockRepository.NewMockAgentRepository(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockPolicyRepository(ctx, pr)
			tt.setMockAgentRepository(ctx, ar)

			pu := usecase.NewPolicyUsecase(to, pr, ar)
			result, err := pu.GetAgents(ctx, tt.inputID, tt.inputUserID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}
