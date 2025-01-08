package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/usecase"
	mockDomain "holos-auth-api/test/mock/domain"
	mockRepository "holos-auth-api/test/mock/domain/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestAuth_Signin(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                       string
		inputUserName              string
		inputPassword              string
		expectError                error
		setMockTransactionObject   func(context.Context, *mockDomain.MockTransactionObject)
		setMockUserRepository      func(context.Context, *mockRepository.MockUserRepository)
		setMockUserTokenRepository func(context.Context, *mockRepository.MockUserTokenRepository)
	}{
		{
			name:          "success",
			inputUserName: "name",
			inputPassword: "password",
			expectError:   nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserRepository: func(ctx context.Context, ur *mockRepository.MockUserRepository) {
				ur.EXPECT().
					FindOneByName(ctx, user.Name).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
					Times(1)
			},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {
				utr.EXPECT().
					Save(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:          "user not found",
			inputUserName: "name",
			inputPassword: "password",
			expectError:   usecase.ErrAuthenticationFailed,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserRepository: func(ctx context.Context, ur *mockRepository.MockUserRepository) {
				ur.EXPECT().
					FindOneByName(ctx, user.Name).
					Return(nil, nil).
					Times(1)
			},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {},
		},
		{
			name:          "verification failed",
			inputUserName: "name",
			inputPassword: "PASSWORD",
			expectError:   entity.ErrAuthenticationFailed,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserRepository: func(ctx context.Context, ur *mockRepository.MockUserRepository) {
				ur.EXPECT().
					FindOneByName(ctx, user.Name).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
					Times(1)
			},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {},
		},
		{
			name:          "find user error",
			inputUserName: "name",
			inputPassword: "password",
			expectError:   sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserRepository: func(ctx context.Context, ur *mockRepository.MockUserRepository) {
				ur.EXPECT().
					FindOneByName(ctx, user.Name).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {},
		},
		{
			name:          "save user token error",
			inputUserName: "name",
			inputPassword: "password",
			expectError:   sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserRepository: func(ctx context.Context, ur *mockRepository.MockUserRepository) {
				ur.EXPECT().
					FindOneByName(ctx, user.Name).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
					Times(1)
			},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {
				utr.EXPECT().
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
			ur := mockRepository.NewMockUserRepository(ctrl)
			utr := mockRepository.NewMockUserTokenRepository(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockUserRepository(ctx, ur)
			tt.setMockUserTokenRepository(ctx, utr)

			au := usecase.NewAuthUsecase(to, ur, utr, nil, nil)
			_, err = au.Signin(ctx, tt.inputUserName, tt.inputPassword)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestAuth_Signout(t *testing.T) {
	userToken, err := entity.NewUserToken(uuid.New())
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                       string
		inputToken                 string
		expectError                error
		setMockTransactionObject   func(context.Context, *mockDomain.MockTransactionObject)
		setMockUserTokenRepository func(context.Context, *mockRepository.MockUserTokenRepository)
	}{
		{
			name:        "success",
			inputToken:  userToken.Token,
			expectError: nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {
				utr.EXPECT().
					Delete(ctx, gomock.Any()).
					Return(nil).
					Times(1)
				utr.EXPECT().
					FindOneByTokenAndNotExpired(ctx, userToken.Token).
					Return(entity.RestoreUserToken(userToken.UserID, userToken.Token, userToken.ExpiresAt), nil).
					Times(1)
			},
		},
		{
			name:        "user token not found",
			inputToken:  userToken.Token,
			expectError: usecase.ErrAuthenticationFailed,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {
				utr.EXPECT().
					FindOneByTokenAndNotExpired(ctx, userToken.Token).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:        "find user token error",
			inputToken:  userToken.Token,
			expectError: sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {
				utr.EXPECT().
					FindOneByTokenAndNotExpired(ctx, userToken.Token).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:        "delete user token error",
			inputToken:  userToken.Token,
			expectError: sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {
				utr.EXPECT().
					Delete(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
				utr.EXPECT().
					FindOneByTokenAndNotExpired(ctx, userToken.Token).
					Return(entity.RestoreUserToken(userToken.UserID, userToken.Token, userToken.ExpiresAt), nil).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			to := mockDomain.NewMockTransactionObject(ctrl)
			utr := mockRepository.NewMockUserTokenRepository(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockUserTokenRepository(ctx, utr)

			au := usecase.NewAuthUsecase(to, nil, utr, nil, nil)
			if err := au.Signout(ctx, tt.inputToken); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}

func TestAuth_Authenticate(t *testing.T) {
	userToken, err := entity.NewUserToken(uuid.New())
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                       string
		inputToken                 string
		expectResult               uuid.UUID
		expectError                error
		setMockUserTokenRepository func(context.Context, *mockRepository.MockUserTokenRepository)
	}{
		{
			name:         "success",
			inputToken:   userToken.Token,
			expectResult: userToken.UserID,
			expectError:  nil,
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {
				utr.EXPECT().
					FindOneByTokenAndNotExpired(ctx, userToken.Token).
					Return(entity.RestoreUserToken(userToken.UserID, userToken.Token, userToken.ExpiresAt), nil).
					Times(1)
			},
		},
		{
			name:         "user token not found",
			inputToken:   userToken.Token,
			expectResult: uuid.Nil,
			expectError:  usecase.ErrAuthenticationFailed,
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {
				utr.EXPECT().
					FindOneByTokenAndNotExpired(ctx, userToken.Token).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:         "find user token error",
			inputToken:   userToken.Token,
			expectResult: uuid.Nil,
			expectError:  sql.ErrConnDone,
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {
				utr.EXPECT().
					FindOneByTokenAndNotExpired(ctx, userToken.Token).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			utr := mockRepository.NewMockUserTokenRepository(ctrl)

			ctx := context.Background()

			tt.setMockUserTokenRepository(ctx, utr)

			au := usecase.NewAuthUsecase(nil, nil, utr, nil, nil)
			result, err := au.Authenticate(ctx, tt.inputToken)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}
		})
	}
}
