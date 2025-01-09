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
		name                   string
		inputPolicy            *entity.Policy
		expectResult           []*entity.Agent
		expectError            error
		setMockAgentRepository func(context.Context, *mockRepository.MockAgentRepository)
	}{
		{
			name:         "success",
			inputPolicy:  policy,
			expectResult: []*entity.Agent{agent},
			expectError:  nil,
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any()).
					Return([]*entity.Agent{agent}, nil).
					Times(1)
			},
		},
		{
			name:         "not found",
			inputPolicy:  policy,
			expectResult: nil,
			expectError:  nil,
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:         "find error",
			inputPolicy:  policy,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindByIDsAndUserIDAndNotDeleted(ctx, gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:                   "no policy",
			inputPolicy:            nil,
			expectResult:           nil,
			expectError:            service.ErrRequiredPolicy,
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ar := mockRepository.NewMockAgentRepository(ctrl)

			ctx := context.Background()

			tt.setMockAgentRepository(ctx, ar)

			as := service.NewPolicyService(ar)
			result, err := as.GetAgents(ctx, tt.inputPolicy)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}
