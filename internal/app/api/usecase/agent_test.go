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

func TestAgent_Create(t *testing.T) {
	tests := []struct {
		name   string
		expect apierr.ApiError
	}{
		{
			name:   "exists",
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			to := test.NewTestTransactionObject()

			ar := mock_repository.NewMockAgentRepository(ctrl)
			ar.EXPECT().Create(ctx, gomock.Any()).Return(nil).AnyTimes()

			au := usecase.NewAgentUsecase(to, ar)
			dto, err := au.Create(ctx, uuid.New(), tt.name)
			if err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
			if reflect.TypeOf(dto).Elem().Name() != "AgentDTO" {
				t.Errorf("expect AgentDTO but got %s", reflect.TypeOf(dto).Elem().Name())
			}
		})
	}
}

func TestAgent_Update(t *testing.T) {
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
			expect:      usecase.ErrAgentNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := entity.NewAgent(uuid.New(), "name")
			if err != nil {
				t.Error(err.Error())
			}

			res := agent
			if tt.isReturnNil {
				res = nil
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			to := test.NewTestTransactionObject()

			ar := mock_repository.NewMockAgentRepository(ctrl)
			ar.EXPECT().FindOneByIDAndUserIDAndNotDeleted(ctx, tt.id, tt.userID).Return(res, nil)
			ar.EXPECT().Update(ctx, gomock.Any()).Return(nil).AnyTimes()

			au := usecase.NewAgentUsecase(to, ar)
			dto, err := au.Update(ctx, tt.id, tt.userID, tt.name)
			if err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
			if reflect.TypeOf(dto).Elem().Name() != "AgentDTO" {
				t.Errorf("expect AgentDTO but got %s", reflect.TypeOf(dto).Elem().Name())
			}
		})
	}
}

func TestAgent_Delete(t *testing.T) {
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
			expect:      usecase.ErrAgentNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := entity.NewAgent(uuid.New(), "name")
			if err != nil {
				t.Error(err.Error())
			}

			res := agent
			if tt.isReturnNil {
				res = nil
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			to := test.NewTestTransactionObject()

			ar := mock_repository.NewMockAgentRepository(ctrl)
			ar.EXPECT().FindOneByIDAndUserIDAndNotDeleted(ctx, tt.id, tt.userID).Return(res, nil)
			ar.EXPECT().Delete(ctx, gomock.Any()).Return(nil).AnyTimes()

			au := usecase.NewAgentUsecase(to, ar)
			err = au.Delete(ctx, tt.id, tt.userID)
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
