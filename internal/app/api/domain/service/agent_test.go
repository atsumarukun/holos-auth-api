package service_test

import (
	"context"
	"database/sql"
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/service"
	mockRepository "holos-auth-api/test/mock/domain/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

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
		name                    string
		inputAgent              *entity.Agent
		inputKeyword            string
		expectResult            []*entity.Policy
		expectError             error
		setMockPolicyRepository func(context.Context, *mockRepository.MockPolicyRepository)
	}{
		{
			name:         "success",
			inputAgent:   agent,
			inputKeyword: "name",
			expectResult: []*entity.Policy{policy},
			expectError:  nil,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByIDsAndNamePrefixAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*entity.Policy{policy}, nil).
					Times(1)
			},
		},
		{
			name:         "not found",
			inputAgent:   agent,
			inputKeyword: "keyword",
			expectResult: nil,
			expectError:  nil,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByIDsAndNamePrefixAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:         "find error",
			inputAgent:   agent,
			inputKeyword: "name",
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByIDsAndNamePrefixAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:                    "no agent",
			inputAgent:              nil,
			inputKeyword:            "name",
			expectResult:            nil,
			expectError:             service.ErrRequiredAgent,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pr := mockRepository.NewMockPolicyRepository(ctrl)

			ctx := context.Background()

			tt.setMockPolicyRepository(ctx, pr)

			as := service.NewAgentService(pr)
			result, err := as.GetPolicies(ctx, tt.inputAgent, tt.inputKeyword)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestAgent_HasPermission(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}
	rootAllowPolicy, err := entity.NewPolicy(agent.UserID, "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}
	rootDenyPolicy, err := entity.NewPolicy(agent.UserID, "name", "DENY", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}
	allowPolicy, err := entity.NewPolicy(agent.UserID, "name", "ALLOW", "STORAGE", "/path/:id", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}
	denyPolicy, err := entity.NewPolicy(agent.UserID, "name", "DENY", "STORAGE", "/path/:id", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                    string
		inputAgent              *entity.Agent
		inputService            string
		inputPath               string
		inputMethod             string
		expectResult            bool
		expectError             error
		setMockPolicyRepository func(context.Context, *mockRepository.MockPolicyRepository)
	}{
		{
			name:         "has permission",
			inputAgent:   agent,
			inputService: "STORAGE",
			inputPath:    "/path/1",
			inputMethod:  "GET",
			expectResult: true,
			expectError:  nil,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any()).
					Return([]*entity.Policy{rootDenyPolicy, allowPolicy}, nil).
					Times(1)
			},
		},
		{
			name:         "has root permission",
			inputAgent:   agent,
			inputService: "STORAGE",
			inputPath:    "/path/1",
			inputMethod:  "GET",
			expectResult: true,
			expectError:  nil,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any()).
					Return([]*entity.Policy{rootAllowPolicy}, nil).
					Times(1)
			},
		},
		{
			name:         "does not have permission",
			inputAgent:   agent,
			inputService: "STORAGE",
			inputPath:    "/path/1",
			inputMethod:  "GET",
			expectResult: false,
			expectError:  nil,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any()).
					Return([]*entity.Policy{rootAllowPolicy, denyPolicy}, nil).
					Times(1)
			},
		},
		{
			name:         "does not have root permission",
			inputAgent:   agent,
			inputService: "STORAGE",
			inputPath:    "/path/1",
			inputMethod:  "GET",
			expectResult: false,
			expectError:  nil,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any()).
					Return([]*entity.Policy{rootDenyPolicy}, nil).
					Times(1)
			},
		},
		{
			name:         "does not have policies",
			inputAgent:   agent,
			inputService: "STORAGE",
			inputPath:    "/path/1",
			inputMethod:  "GET",
			expectResult: false,
			expectError:  nil,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any()).
					Return([]*entity.Policy{}, nil).
					Times(1)
			},
		},
		{
			name:         "service not matched",
			inputAgent:   agent,
			inputService: "CONTENT",
			inputPath:    "/path/1",
			inputMethod:  "GET",
			expectResult: false,
			expectError:  nil,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any()).
					Return([]*entity.Policy{rootAllowPolicy}, nil).
					Times(1)
			},
		},
		{
			name:         "method not matched",
			inputAgent:   agent,
			inputService: "STORAGE",
			inputPath:    "/path/1",
			inputMethod:  "POST",
			expectResult: false,
			expectError:  nil,
			setMockPolicyRepository: func(ctx context.Context, pr *mockRepository.MockPolicyRepository) {
				pr.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any()).
					Return([]*entity.Policy{}, nil).
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

			as := service.NewAgentService(pr)
			result, err := as.HasPermission(ctx, tt.inputAgent, tt.inputService, tt.inputPath, tt.inputMethod)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if result != tt.expectResult {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectResult, result)
			}
		})
	}
}
