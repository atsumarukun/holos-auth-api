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
)

func TestUser_Exists(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                  string
		expectResult          bool
		expectError           error
		setMockUserRepository func(context.Context, *mockRepository.MockUserRepository)
	}{
		{
			name:         "exists",
			expectResult: true,
			setMockUserRepository: func(ctx context.Context, ur *mockRepository.MockUserRepository) {
				ur.EXPECT().
					FindOneByName(ctx, user.Name).
					Return(user, nil)
			},
		},
		{
			name:         "not exists",
			expectResult: false,
			expectError:  nil,
			setMockUserRepository: func(ctx context.Context, ur *mockRepository.MockUserRepository) {
				ur.EXPECT().
					FindOneByName(ctx, user.Name).
					Return(nil, nil)
			},
		},
		{
			name:         "mock return error",
			expectResult: false,
			expectError:  sql.ErrConnDone,
			setMockUserRepository: func(ctx context.Context, ur *mockRepository.MockUserRepository) {
				ur.EXPECT().
					FindOneByName(ctx, user.Name).
					Return(nil, sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ur := mockRepository.NewMockUserRepository(ctrl)

			ctx := context.Background()

			tt.setMockUserRepository(ctx, ur)

			s := service.NewUserService(ur)
			exists, err := s.Exists(ctx, user)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if exists != tt.expectResult {
				t.Errorf("\nexpect %t \ngot %t", tt.expectResult, exists)
			}
		})
	}
}
