package usecase_test

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"holos-auth-api/internal/app/api/usecase"
	"holos-auth-api/test"
	mock_repository "holos-auth-api/test/mock/domain/repository"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestPolicy_Create(t *testing.T) {
	tests := []struct {
		name   string
		expect apierr.ApiError
	}{
		{
			name:   "valid",
			expect: nil,
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		to := test.NewTestTransactionObject()

		pr := mock_repository.NewMockPolicyRepository(ctrl)
		pr.EXPECT().Create(ctx, gomock.Any()).Return(nil)

		pu := usecase.NewPolicyUsecase(to, pr, nil)
		dto, err := pu.Create(ctx, uuid.New(), "name", "STORAGE", "/", []string{"GET"})
		if err != tt.expect {
			if err == nil {
				t.Error("expect err but got nil")
			} else {
				t.Error(err.Error())
			}
		}
		if reflect.TypeOf(dto).Elem().Name() != "PolicyDTO" {
			t.Errorf("expect PolicyDTO but got %s", reflect.TypeOf(dto).Elem().Name())
		}
	}
}

func TestPolicy_Update(t *testing.T) {
	tests := []struct {
		id          uuid.UUID
		userID      uuid.UUID
		name        string
		isReturnNil bool
		expect      apierr.ApiError
	}{
		{
			id:          uuid.New(),
			userID:      uuid.New(),
			name:        "success",
			isReturnNil: false,
			expect:      nil,
		},
		{
			id:          uuid.New(),
			userID:      uuid.New(),
			name:        "not_found",
			isReturnNil: true,
			expect:      usecase.ErrPolicyNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
			if err != nil {
				t.Error(err.Error())
			}

			res := policy
			if tt.isReturnNil {
				res = nil
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			to := test.NewTestTransactionObject()

			pr := mock_repository.NewMockPolicyRepository(ctrl)
			pr.EXPECT().FindOneByIDAndUserIDAndNotDeleted(ctx, tt.id, tt.userID).Return(res, nil)
			pr.EXPECT().Update(ctx, gomock.Any()).Return(nil).AnyTimes()

			pu := usecase.NewPolicyUsecase(to, pr, nil)
			dto, err := pu.Update(ctx, tt.id, tt.userID, tt.name, "STORAGE", "/", []string{"GET"})
			if err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
			if reflect.TypeOf(dto).Elem().Name() != "PolicyDTO" {
				t.Errorf("expect PolicyDTO but got %s", reflect.TypeOf(dto).Elem().Name())
			}
		})
	}
}

func TestPolicy_Delete(t *testing.T) {
	tests := []struct {
		id          uuid.UUID
		userID      uuid.UUID
		name        string
		isReturnNil bool
		expect      apierr.ApiError
	}{
		{
			id:          uuid.New(),
			userID:      uuid.New(),
			name:        "success",
			isReturnNil: false,
			expect:      nil,
		},
		{
			id:          uuid.New(),
			userID:      uuid.New(),
			name:        "not_found",
			isReturnNil: true,
			expect:      usecase.ErrPolicyNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
			if err != nil {
				t.Error(err.Error())
			}

			res := policy
			if tt.isReturnNil {
				res = nil
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			to := test.NewTestTransactionObject()

			pr := mock_repository.NewMockPolicyRepository(ctrl)
			pr.EXPECT().FindOneByIDAndUserIDAndNotDeleted(ctx, tt.id, tt.userID).Return(res, nil)
			pr.EXPECT().Delete(ctx, gomock.Any()).Return(nil).AnyTimes()

			pu := usecase.NewPolicyUsecase(to, pr, nil)
			err = pu.Delete(ctx, tt.id, tt.userID)
			if err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestPolicy_Gets(t *testing.T) {
	tests := []struct {
		id          uuid.UUID
		userID      uuid.UUID
		name        string
		isReturnNil bool
		expect      apierr.ApiError
	}{
		{
			id:          uuid.New(),
			userID:      uuid.New(),
			name:        "success",
			isReturnNil: false,
			expect:      nil,
		},
		{
			id:          uuid.New(),
			userID:      uuid.New(),
			name:        "not_found",
			isReturnNil: true,
			expect:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
			if err != nil {
				t.Error(err.Error())
			}

			res := []*entity.Policy{policy}
			if tt.isReturnNil {
				res = nil
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			to := test.NewTestTransactionObject()

			pr := mock_repository.NewMockPolicyRepository(ctrl)
			pr.EXPECT().FindByUserIDAndNotDeleted(ctx, tt.userID).Return(res, nil)

			pu := usecase.NewPolicyUsecase(to, pr, nil)
			_, err = pu.Gets(ctx, tt.userID)
			if err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestPolicy_UpdateAgents(t *testing.T) {
	tests := []struct {
		id     uuid.UUID
		userID uuid.UUID
		name   string
		expect apierr.ApiError
	}{
		{
			id:     uuid.New(),
			userID: uuid.New(),
			name:   "success",
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := entity.NewAgent(uuid.New(), "name")
			if err != nil {
				t.Error(err.Error())
			}
			policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
			if err != nil {
				t.Error(err.Error())
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			to := test.NewTestTransactionObject()

			pr := mock_repository.NewMockPolicyRepository(ctrl)
			pr.EXPECT().FindOneByIDAndUserIDAndNotDeleted(ctx, tt.id, tt.userID).Return(policy, nil)
			pr.EXPECT().UpdateAgents(ctx, gomock.Any(), gomock.Any()).Return(nil)

			ar := mock_repository.NewMockAgentRepository(ctrl)
			ar.EXPECT().FindByIDsAndUserIDAndNotDeleted(ctx, []uuid.UUID{agent.ID}, tt.userID)

			pu := usecase.NewPolicyUsecase(to, pr, ar)
			_, err = pu.UpdateAgents(ctx, tt.id, tt.userID, []uuid.UUID{agent.ID})
			if err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestPolicy_GetAgents(t *testing.T) {
	tests := []struct {
		id          uuid.UUID
		userID      uuid.UUID
		name        string
		isReturnNil bool
		expect      apierr.ApiError
	}{
		{
			id:          uuid.New(),
			userID:      uuid.New(),
			name:        "success",
			isReturnNil: false,
			expect:      nil,
		},
		{
			id:          uuid.New(),
			userID:      uuid.New(),
			name:        "not_found",
			isReturnNil: true,
			expect:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := entity.NewAgent(uuid.New(), "name")
			if err != nil {
				t.Error(err.Error())
			}

			res := []*entity.Agent{agent}
			if tt.isReturnNil {
				res = nil
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			to := test.NewTestTransactionObject()

			pr := mock_repository.NewMockPolicyRepository(ctrl)
			pr.EXPECT().GetAgents(ctx, tt.id, tt.userID).Return(res, nil)

			pu := usecase.NewPolicyUsecase(to, pr, nil)
			_, err = pu.GetAgents(ctx, tt.id, tt.userID)
			if err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
		})
	}
}
