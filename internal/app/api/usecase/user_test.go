package usecase_test

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/usecase"
	"holos-auth-api/internal/pkg/apierr"
	"holos-auth-api/test"
	mock_repository "holos-auth-api/test/mock/domain/repository"
	mock_service "holos-auth-api/test/mock/domain/service"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestUser_Create(t *testing.T) {
	tests := []struct {
		name     string
		password string
		exists   bool
		expect   apierr.ApiError
	}{
		{
			name:     "exists",
			password: "password",
			exists:   false,
			expect:   nil,
		},
		{
			name:     "not_exists",
			password: "password",
			exists:   true,
			expect:   usecase.ErrUserAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			to := test.NewTestTransactionObject()

			ur := mock_repository.NewMockUserRepository(ctrl)
			ur.EXPECT().Create(ctx, gomock.Any()).Return(nil).AnyTimes()

			us := mock_service.NewMockUserService(ctrl)
			us.EXPECT().Exists(ctx, gomock.Any()).Return(tt.exists, nil)

			uu := usecase.NewUserUsecase(to, ur, us)
			dto, err := uu.Create(ctx, tt.name, tt.password, tt.password)
			if err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
			if reflect.TypeOf(dto).Elem().Name() != "UserDTO" {
				t.Errorf("expect UserDTO but got %s", reflect.TypeOf(dto).Elem().Name())
			}
		})
	}
}

func TestUser_Update(t *testing.T) {
	tests := []struct {
		id          uuid.UUID
		name        string
		password    string
		isReturnNil bool
		expect      apierr.ApiError
	}{
		{
			id:          uuid.New(),
			name:        "exists",
			password:    "password",
			isReturnNil: false,
			expect:      nil,
		},
		{
			id:          uuid.New(),
			name:        "not_exists",
			password:    "password",
			isReturnNil: true,
			expect:      usecase.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entity.NewUser(tt.name, tt.password, tt.password)
			if err != nil {
				t.Error(err.Error())
			}

			res := user
			if tt.isReturnNil {
				res = nil
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			to := test.NewTestTransactionObject()

			ur := mock_repository.NewMockUserRepository(ctrl)
			ur.EXPECT().FindOneByID(ctx, tt.id).Return(res, nil)
			ur.EXPECT().Update(ctx, gomock.Any()).Return(nil).AnyTimes()

			us := mock_service.NewMockUserService(ctrl)

			uu := usecase.NewUserUsecase(to, ur, us)
			dto, err := uu.Update(ctx, tt.id, tt.password, tt.password, tt.password)
			if err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
			if reflect.TypeOf(dto).Elem().Name() != "UserDTO" {
				t.Errorf("expect UserDTO but got %s", reflect.TypeOf(dto).Elem().Name())
			}
		})
	}
}

func TestUser_Delete(t *testing.T) {
	tests := []struct {
		id          uuid.UUID
		name        string
		password    string
		isReturnNil bool
		expect      apierr.ApiError
	}{
		{
			id:          uuid.New(),
			name:        "exists",
			password:    "password",
			isReturnNil: false,
			expect:      nil,
		},
		{
			id:          uuid.New(),
			name:        "not_exists",
			password:    "password",
			isReturnNil: true,
			expect:      usecase.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entity.NewUser(tt.name, tt.password, tt.password)
			if err != nil {
				t.Error(err.Error())
			}

			res := user
			if tt.isReturnNil {
				res = nil
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			to := test.NewTestTransactionObject()

			ur := mock_repository.NewMockUserRepository(ctrl)
			ur.EXPECT().FindOneByID(ctx, tt.id).Return(res, nil)
			ur.EXPECT().Delete(ctx, gomock.Any()).Return(nil).AnyTimes()

			us := mock_service.NewMockUserService(ctrl)

			uu := usecase.NewUserUsecase(to, ur, us)
			err = uu.Delete(ctx, tt.id, tt.password)
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
