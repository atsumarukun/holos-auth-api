package usecase_test

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/usecase"
	"holos-auth-api/internal/pkg/apierr"
	"holos-auth-api/test"
	mock_repository "holos-auth-api/test/mock/domain/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAuth_Signin(t *testing.T) {
	tests := []struct {
		name           string
		userPassword   string
		signinPassword string
		isReturnNil    bool
		expect         apierr.ApiError
	}{
		{
			name:           "success",
			userPassword:   "password",
			signinPassword: "password",
			isReturnNil:    false,
			expect:         nil,
		},
		{
			name:           "user_not_found",
			userPassword:   "password",
			signinPassword: "password",
			isReturnNil:    true,
			expect:         usecase.ErrAuthenticationFailed,
		},
		{
			name:           "incorrect_password",
			userPassword:   "password",
			signinPassword: "incorrect_password",
			isReturnNil:    false,
			expect:         entity.ErrAuthenticationFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entity.NewUser(tt.name, tt.userPassword, tt.userPassword)
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
			ur.EXPECT().FindOneByName(ctx, tt.name).Return(res, nil)

			utr := mock_repository.NewMockUserTokenRepository(ctrl)
			utr.EXPECT().Save(ctx, gomock.Any()).Return(nil).AnyTimes()

			au := usecase.NewAuthUsecase(to, ur, utr)
			_, err = au.Signin(ctx, tt.name, tt.signinPassword)
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

func TestAuth_Signout(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		isReturnNil bool
		expect      apierr.ApiError
	}{
		{
			name:        "success",
			token:       "token",
			isReturnNil: false,
			expect:      nil,
		},
		{
			name:        "user_token_not_found",
			token:       "token",
			isReturnNil: true,
			expect:      usecase.ErrAuthenticationFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userToken, err := entity.NewUserToken(uuid.New())
			if err != nil {
				t.Error(err.Error())
			}

			res := userToken
			if tt.isReturnNil {
				res = nil
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			to := test.NewTestTransactionObject()

			ur := mock_repository.NewMockUserRepository(ctrl)

			utr := mock_repository.NewMockUserTokenRepository(ctrl)
			utr.EXPECT().FindOneByTokenAndNotExpired(ctx, tt.token).Return(res, nil)
			utr.EXPECT().Delete(ctx, gomock.Any()).Return(nil).AnyTimes()

			au := usecase.NewAuthUsecase(to, ur, utr)
			if err := au.Signout(ctx, tt.token); err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
		})
	}
}
