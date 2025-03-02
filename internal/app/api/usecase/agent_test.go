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
	mockService "holos-auth-api/test/mock/domain/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestAgent_Create(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                   string
		inputUserID            uuid.UUID
		inputName              string
		expectResult           *dto.AgentDTO
		expectError            error
		setMockAgentRepository func(context.Context, *mockRepository.MockAgentRepository)
	}{
		{
			name:         "success",
			inputUserID:  agent.UserID,
			inputName:    "name",
			expectResult: &dto.AgentDTO{ID: agent.ID, UserID: agent.UserID, Name: agent.Name, Policies: []uuid.UUID{}, CreatedAt: agent.CreatedAt, UpdatedAt: agent.UpdatedAt},
			expectError:  nil,
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:                   "invalid name",
			inputUserID:            agent.UserID,
			inputName:              "なまえ",
			expectResult:           nil,
			expectError:            entity.ErrInvalidAgentName,
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {},
		},
		{
			name:         "create error",
			inputUserID:  agent.UserID,
			inputName:    "name",
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
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

			ar := mockRepository.NewMockAgentRepository(ctrl)

			ctx := context.Background()

			tt.setMockAgentRepository(ctx, ar)

			au := usecase.NewAgentUsecase(nil, ar, nil, nil, nil)
			result, err := au.Create(ctx, tt.inputUserID, tt.inputName)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.AgentDTO{}, "ID", "CreatedAt", "UpdatedAt"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestAgent_Update(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                     string
		inputID                  uuid.UUID
		inputUserID              uuid.UUID
		inputName                string
		expectResult             *dto.AgentDTO
		expectError              error
		setMockTransactionObject func(context.Context, *mockDomain.MockTransactionObject)
		setMockAgentRepository   func(context.Context, *mockRepository.MockAgentRepository)
	}{
		{
			name:         "success",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			inputName:    "update",
			expectResult: &dto.AgentDTO{ID: agent.ID, UserID: agent.UserID, Name: "update", Policies: []uuid.UUID{}, CreatedAt: agent.CreatedAt, UpdatedAt: agent.UpdatedAt},
			expectError:  nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt), nil).
					Times(1)
				ar.EXPECT().
					Update(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:         "invalid name",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			inputName:    "なまえ",
			expectResult: nil,
			expectError:  entity.ErrInvalidAgentName,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt), nil).
					Times(1)
			},
		},
		{
			name:         "agent not found",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			inputName:    "update",
			expectResult: nil,
			expectError:  usecase.ErrAgentNotFound,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:         "find error",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			inputName:    "update",
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
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:         "update error",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			inputName:    "update",
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
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt), nil).
					Times(1)
				ar.EXPECT().
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
			ar := mockRepository.NewMockAgentRepository(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockAgentRepository(ctx, ar)

			au := usecase.NewAgentUsecase(to, ar, nil, nil, nil)
			result, err := au.Update(ctx, tt.inputID, tt.inputUserID, tt.inputName)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.AgentDTO{}, "UpdatedAt"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestAgent_Delete(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                     string
		inputID                  uuid.UUID
		inputUserID              uuid.UUID
		expectError              error
		setMockTransactionObject func(context.Context, *mockDomain.MockTransactionObject)
		setMockAgentRepository   func(context.Context, *mockRepository.MockAgentRepository)
	}{
		{
			name:        "success",
			inputID:     agent.ID,
			inputUserID: agent.UserID,
			expectError: nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt), nil).
					Times(1)
				ar.EXPECT().
					Delete(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:        "agent not found",
			inputID:     agent.ID,
			inputUserID: agent.UserID,
			expectError: usecase.ErrAgentNotFound,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:        "find error",
			inputID:     agent.ID,
			inputUserID: agent.UserID,
			expectError: sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:        "delete error",
			inputID:     agent.ID,
			inputUserID: agent.UserID,
			expectError: sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt), nil).
					Times(1)
				ar.EXPECT().
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
			ar := mockRepository.NewMockAgentRepository(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockAgentRepository(ctx, ar)

			au := usecase.NewAgentUsecase(to, ar, nil, nil, nil)
			if err := au.Delete(ctx, tt.inputID, tt.inputUserID); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestAgent_Get(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                   string
		inputID                uuid.UUID
		inputUserID            uuid.UUID
		expectResult           *dto.AgentDTO
		expectError            error
		setMockAgentRepository func(context.Context, *mockRepository.MockAgentRepository)
	}{
		{
			name:         "success",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			expectResult: &dto.AgentDTO{ID: agent.ID, UserID: agent.UserID, Name: agent.Name, Policies: []uuid.UUID{}, CreatedAt: agent.CreatedAt, UpdatedAt: agent.UpdatedAt},
			expectError:  nil,
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt), nil).
					Times(1)
			},
		},
		{
			name:         "agent not found",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			expectResult: nil,
			expectError:  usecase.ErrAgentNotFound,
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:         "find error",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ar := mockRepository.NewMockAgentRepository(ctrl)

			ctx := context.Background()

			tt.setMockAgentRepository(ctx, ar)

			au := usecase.NewAgentUsecase(nil, ar, nil, nil, nil)
			result, err := au.Get(ctx, tt.inputID, tt.inputUserID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.AgentDTO{}, "UpdatedAt"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestAgent_Gets(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                   string
		inputKeyword           string
		inputUserID            uuid.UUID
		expectResult           []*dto.AgentDTO
		expectError            error
		setMockAgentRepository func(context.Context, *mockRepository.MockAgentRepository)
	}{
		{
			name:         "found",
			inputKeyword: "name",
			inputUserID:  agent.UserID,
			expectResult: []*dto.AgentDTO{{ID: agent.ID, UserID: agent.UserID, Name: agent.Name, Policies: []uuid.UUID{}, CreatedAt: agent.CreatedAt, UpdatedAt: agent.UpdatedAt}},
			expectError:  nil,
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindByNamePrefixAndUserIDAndNotDeleted(ctx, gomock.Any(), agent.UserID).
					Return([]*entity.Agent{entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt)}, nil).
					Times(1)
			},
		},
		{
			name:         "not found",
			inputKeyword: "name",
			inputUserID:  agent.UserID,
			expectResult: []*dto.AgentDTO{},
			expectError:  nil,
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindByNamePrefixAndUserIDAndNotDeleted(ctx, gomock.Any(), agent.UserID).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:         "find error",
			inputKeyword: "name",
			inputUserID:  agent.UserID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindByNamePrefixAndUserIDAndNotDeleted(ctx, gomock.Any(), agent.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ar := mockRepository.NewMockAgentRepository(ctrl)

			ctx := context.Background()

			tt.setMockAgentRepository(ctx, ar)

			au := usecase.NewAgentUsecase(nil, ar, nil, nil, nil)
			result, err := au.Gets(ctx, tt.inputKeyword, tt.inputUserID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestAgent_UpdatePolicies(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}
	policy, err := entity.NewPolicy(agent.UserID, "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                     string
		inputID                  uuid.UUID
		inputUserID              uuid.UUID
		inputPolicyIDs           []uuid.UUID
		expectResult             []*dto.PolicyDTO
		expectError              error
		setMockTransactionObject func(context.Context, *mockDomain.MockTransactionObject)
		setMockAgentRepository   func(context.Context, *mockRepository.MockAgentRepository)
		setMockPolicyRepository  func(context.Context, *mockRepository.MockPolicyRepository)
	}{
		{
			name:           "success",
			inputID:        agent.ID,
			inputUserID:    agent.UserID,
			inputPolicyIDs: []uuid.UUID{policy.ID},
			expectResult:   []*dto.PolicyDTO{{ID: policy.ID, UserID: policy.UserID, Name: policy.Name, Effect: policy.Effect, Service: policy.Service, Path: policy.Path, Methods: policy.Methods, Agents: []uuid.UUID{}, CreatedAt: policy.CreatedAt, UpdatedAt: policy.UpdatedAt}},
			expectError:    nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt), nil).
					Times(1)
				ar.EXPECT().
					Update(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any()).
					Return([]*entity.Policy{entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt)}, nil).
					Times(1)
			},
		},
		{
			name:           "agent not found",
			inputID:        agent.ID,
			inputUserID:    agent.UserID,
			inputPolicyIDs: []uuid.UUID{policy.ID},
			expectResult:   nil,
			expectError:    usecase.ErrAgentNotFound,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(nil, nil).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {},
		},
		{
			name:           "find agent error",
			inputID:        agent.ID,
			inputUserID:    agent.UserID,
			inputPolicyIDs: []uuid.UUID{policy.ID},
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {},
		},
		{
			name:           "find policies error",
			inputID:        agent.ID,
			inputUserID:    agent.UserID,
			inputPolicyIDs: []uuid.UUID{policy.ID},
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt), nil).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:           "update agent error",
			inputID:        agent.ID,
			inputUserID:    agent.UserID,
			inputPolicyIDs: []uuid.UUID{policy.ID},
			expectResult:   nil,
			expectError:    sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt), nil).
					Times(1)
				ar.EXPECT().
					Update(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any()).
					Return([]*entity.Policy{entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt)}, nil).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			to := mockDomain.NewMockTransactionObject(ctrl)
			ar := mockRepository.NewMockAgentRepository(ctrl)
			pr := mockRepository.NewMockPolicyRepository(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockAgentRepository(ctx, ar)
			tt.setMockPolicyRepository(ctx, pr)

			au := usecase.NewAgentUsecase(to, ar, nil, pr, nil)
			result, err := au.UpdatePolicies(ctx, tt.inputID, tt.inputUserID, tt.inputPolicyIDs)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestAgent_GetPolicies(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}
	policy, err := entity.NewPolicy(agent.UserID, "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                     string
		inputID                  uuid.UUID
		inputUserID              uuid.UUID
		inputKeyword             string
		expectResult             []*dto.PolicyDTO
		expectError              error
		setMockTransactionObject func(context.Context, *mockDomain.MockTransactionObject)
		setMockAgentRepository   func(context.Context, *mockRepository.MockAgentRepository)
		setMockAgentService      func(context.Context, *mockService.MockAgentService)
	}{
		{
			name:         "success",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			inputKeyword: "name",
			expectResult: []*dto.PolicyDTO{{ID: policy.ID, UserID: policy.UserID, Name: policy.Name, Effect: policy.Effect, Service: policy.Service, Path: policy.Path, Methods: policy.Methods, Agents: []uuid.UUID{}, CreatedAt: policy.CreatedAt, UpdatedAt: policy.UpdatedAt}},
			expectError:  nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt), nil).
					Times(1)
			},
			setMockAgentService: func(ctx context.Context, as *mockService.MockAgentService) {
				as.EXPECT().
					GetPolicies(ctx, gomock.Any(), gomock.Any()).
					Return([]*entity.Policy{entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Effect, policy.Service, policy.Path, policy.Methods, policy.Agents, policy.CreatedAt, policy.UpdatedAt)}, nil).
					Times(1)
			},
		},
		{
			name:         "agent not found",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			inputKeyword: "name",
			expectResult: nil,
			expectError:  usecase.ErrAgentNotFound,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(nil, nil).
					Times(1)
			},
			setMockAgentService: func(ctx context.Context, as *mockService.MockAgentService) {},
		},
		{
			name:         "policies not found",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			inputKeyword: "keyword",
			expectResult: []*dto.PolicyDTO{},
			expectError:  nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt), nil).
					Times(1)
			},
			setMockAgentService: func(ctx context.Context, as *mockService.MockAgentService) {
				as.EXPECT().
					GetPolicies(ctx, gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:         "find agent error",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			inputKeyword: "name",
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
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockAgentService: func(ctx context.Context, as *mockService.MockAgentService) {},
		},
		{
			name:         "find policies error",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			inputKeyword: "name",
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
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt), nil).
					Times(1)
			},
			setMockAgentService: func(ctx context.Context, as *mockService.MockAgentService) {
				as.EXPECT().
					GetPolicies(ctx, gomock.Any(), gomock.Any()).
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
			ar := mockRepository.NewMockAgentRepository(ctrl)
			as := mockService.NewMockAgentService(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockAgentRepository(ctx, ar)
			tt.setMockAgentService(ctx, as)

			au := usecase.NewAgentUsecase(to, ar, nil, nil, as)
			result, err := au.GetPolicies(ctx, tt.inputID, tt.inputUserID, tt.inputKeyword)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestAgent_GenerateToken(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                        string
		inputID                     uuid.UUID
		inputUserID                 uuid.UUID
		expectError                 error
		setMockTransactionObject    func(context.Context, *mockDomain.MockTransactionObject)
		setMockAgentRepository      func(context.Context, *mockRepository.MockAgentRepository)
		setMockAgentTokenRepository func(context.Context, *mockRepository.MockAgentTokenRepository)
	}{
		{
			name:        "success",
			inputID:     agent.ID,
			inputUserID: agent.UserID,
			expectError: nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt), nil).
					Times(1)
			},
			setMockAgentTokenRepository: func(ctx context.Context, pr *mockRepository.MockAgentTokenRepository) {
				pr.EXPECT().
					Save(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:        "agent not found",
			inputID:     agent.ID,
			inputUserID: agent.UserID,
			expectError: usecase.ErrAgentNotFound,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(nil, nil).
					Times(1)
			},
			setMockAgentTokenRepository: func(ctx context.Context, pr *mockRepository.MockAgentTokenRepository) {},
		},
		{
			name:        "find agent error",
			inputID:     agent.ID,
			inputUserID: agent.UserID,
			expectError: sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockAgentTokenRepository: func(ctx context.Context, pr *mockRepository.MockAgentTokenRepository) {},
		},
		{
			name:        "save agent token error",
			inputID:     agent.ID,
			inputUserID: agent.UserID,
			expectError: sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByIDAndUserIDAndNotDeleted(ctx, agent.ID, agent.UserID).
					Return(entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.Policies, agent.CreatedAt, agent.UpdatedAt), nil).
					Times(1)
			},
			setMockAgentTokenRepository: func(ctx context.Context, pr *mockRepository.MockAgentTokenRepository) {
				pr.EXPECT().
					Save(ctx, gomock.Any()).
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
			ar := mockRepository.NewMockAgentRepository(ctrl)
			atr := mockRepository.NewMockAgentTokenRepository(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockAgentRepository(ctx, ar)
			tt.setMockAgentTokenRepository(ctx, atr)

			au := usecase.NewAgentUsecase(to, ar, atr, nil, nil)
			_, err := au.GenerateToken(ctx, tt.inputID, tt.inputUserID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestAgent_DeleteToken(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}
	agentToken, err := entity.NewAgentToken(agent.ID)
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                        string
		inputID                     uuid.UUID
		inputUserID                 uuid.UUID
		expectError                 error
		setMockTransactionObject    func(context.Context, *mockDomain.MockTransactionObject)
		setMockAgentTokenRepository func(context.Context, *mockRepository.MockAgentTokenRepository)
	}{
		{
			name:        "success",
			inputID:     agent.ID,
			inputUserID: agent.UserID,
			expectError: nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentTokenRepository: func(ctx context.Context, pr *mockRepository.MockAgentTokenRepository) {
				pr.EXPECT().
					FindOneByAgentIDAndUserID(ctx, agent.ID, agent.UserID).
					Return(agentToken, nil).
					Times(1)
				pr.EXPECT().
					Delete(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:        "agent token not found",
			inputID:     agent.ID,
			inputUserID: agent.UserID,
			expectError: usecase.ErrAgentTokenNotFound,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentTokenRepository: func(ctx context.Context, pr *mockRepository.MockAgentTokenRepository) {
				pr.EXPECT().
					FindOneByAgentIDAndUserID(ctx, agent.ID, agent.UserID).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:        "find agent token error",
			inputID:     agent.ID,
			inputUserID: agent.UserID,
			expectError: sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentTokenRepository: func(ctx context.Context, pr *mockRepository.MockAgentTokenRepository) {
				pr.EXPECT().
					FindOneByAgentIDAndUserID(ctx, agent.ID, agent.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:        "delete agent token error",
			inputID:     agent.ID,
			inputUserID: agent.UserID,
			expectError: sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockAgentTokenRepository: func(ctx context.Context, pr *mockRepository.MockAgentTokenRepository) {
				pr.EXPECT().
					FindOneByAgentIDAndUserID(ctx, agent.ID, agent.UserID).
					Return(agentToken, nil).
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
			atr := mockRepository.NewMockAgentTokenRepository(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockAgentTokenRepository(ctx, atr)

			au := usecase.NewAgentUsecase(to, nil, atr, nil, nil)
			if err := au.DeleteToken(ctx, tt.inputID, tt.inputUserID); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestAgent_GetToken(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}
	agentToken, err := entity.NewAgentToken(agent.ID)
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                        string
		inputID                     uuid.UUID
		inputUserID                 uuid.UUID
		expectResult                *dto.AgentTokenDTO
		expectError                 error
		setMockAgentTokenRepository func(context.Context, *mockRepository.MockAgentTokenRepository)
	}{
		{
			name:         "success",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			expectResult: &dto.AgentTokenDTO{AgentID: agentToken.AgentID, Token: agentToken.Token, GeneratedAt: agentToken.GeneratedAt},
			expectError:  nil,
			setMockAgentTokenRepository: func(ctx context.Context, pr *mockRepository.MockAgentTokenRepository) {
				pr.EXPECT().
					FindOneByAgentIDAndUserID(ctx, agent.ID, agent.UserID).
					Return(agentToken, nil).
					Times(1)
			},
		},
		{
			name:         "agent token not found",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			expectResult: nil,
			expectError:  nil,
			setMockAgentTokenRepository: func(ctx context.Context, pr *mockRepository.MockAgentTokenRepository) {
				pr.EXPECT().
					FindOneByAgentIDAndUserID(ctx, agent.ID, agent.UserID).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:         "find agent token error",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockAgentTokenRepository: func(ctx context.Context, pr *mockRepository.MockAgentTokenRepository) {
				pr.EXPECT().
					FindOneByAgentIDAndUserID(ctx, agent.ID, agent.UserID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			atr := mockRepository.NewMockAgentTokenRepository(ctrl)

			ctx := context.Background()

			tt.setMockAgentTokenRepository(ctx, atr)

			au := usecase.NewAgentUsecase(nil, nil, atr, nil, nil)
			result, err := au.GetToken(ctx, tt.inputID, tt.inputUserID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}
