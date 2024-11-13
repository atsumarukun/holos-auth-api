package service_test

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/service"
	mock_repository "holos-auth-api/test/mock/domain/repository"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUser_Exists(t *testing.T) {
	tests := []struct {
		name        string
		isReturnNil bool
		expect      bool
	}{
		{
			name:        "exists",
			isReturnNil: false,
			expect:      true,
		},
		{
			name:        "not_exists",
			isReturnNil: true,
			expect:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entity.NewUser(tt.name, "password", "password")
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

			ur := mock_repository.NewMockUserRepository(ctrl)
			ur.EXPECT().FindOneByName(ctx, tt.name).Return(res, nil)

			us := service.NewUserService(ur)
			exists, err := us.Exists(ctx, user)
			if err != nil {
				t.Error(err.Error())
			}
			if exists != tt.expect {
				t.Errorf("expect %t but got %t", tt.expect, exists)
			}
		})
	}
}
