package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/usecase"
	mockDomain "holos-auth-api/test/mock/domain"
	mockRepository "holos-auth-api/test/mock/domain/repository"
	mockService "holos-auth-api/test/mock/domain/service"
	"testing"

	"github.com/golang/mock/gomock"
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
			if result != tt.expectResult {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectResult, result)
			}
		})
	}
}

func TestAuth_Authorize(t *testing.T) {
	userToken, err := entity.NewUserToken(uuid.New())
	if err != nil {
		t.Error(err.Error())
	}
	agent, err := entity.NewAgent(userToken.UserID, "name")
	if err != nil {
		t.Error(err.Error())
	}
	agentToken, err := entity.NewAgentToken(agent.ID)
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                       string
		inputToken                 string
		inputOperatorType          string
		inputService               string
		inputPath                  string
		inputMethod                string
		expectResult               uuid.UUID
		expectError                error
		setMockTransactionObject   func(context.Context, *mockDomain.MockTransactionObject)
		setMockUserTokenRepository func(context.Context, *mockRepository.MockUserTokenRepository)
		setMockAgentRepository     func(context.Context, *mockRepository.MockAgentRepository)
		setMockAgentService        func(context.Context, *mockService.MockAgentService)
	}{
		{
			name:                     "successful authentication of user access",
			inputToken:               userToken.Token,
			inputOperatorType:        "USER",
			inputService:             "STORAGE",
			inputPath:                "/",
			inputMethod:              "GET",
			expectResult:             userToken.UserID,
			expectError:              nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {
				utr.EXPECT().
					FindOneByTokenAndNotExpired(ctx, userToken.Token).
					Return(entity.RestoreUserToken(userToken.UserID, userToken.Token, userToken.ExpiresAt), nil).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {},
			setMockAgentService:    func(ctx context.Context, as *mockService.MockAgentService) {},
		},
		{
			name:                     "failure to authenticate user access",
			inputToken:               userToken.Token,
			inputOperatorType:        "USER",
			inputService:             "STORAGE",
			inputPath:                "/",
			inputMethod:              "GET",
			expectResult:             uuid.Nil,
			expectError:              usecase.ErrAuthenticationFailed,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {
				utr.EXPECT().
					FindOneByTokenAndNotExpired(ctx, userToken.Token).
					Return(nil, nil).
					Times(1)
			},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {},
			setMockAgentService:    func(ctx context.Context, as *mockService.MockAgentService) {},
		},
		{
			name:              "successful authentication of agent access",
			inputToken:        agentToken.Token,
			inputOperatorType: "AGENT",
			inputService:      "STORAGE",
			inputPath:         "/",
			inputMethod:       "GET",
			expectResult:      userToken.UserID,
			expectError:       nil,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByTokenAndNotDeleted(ctx, gomock.Any()).
					Return(agent, nil).
					Times(1)
			},
			setMockAgentService: func(ctx context.Context, as *mockService.MockAgentService) {
				as.EXPECT().
					HasPermission(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(true, nil).
					Times(1)
			},
		},
		{
			name:              "failure to authenticate agent access",
			inputToken:        agentToken.Token,
			inputOperatorType: "AGENT",
			inputService:      "STORAGE",
			inputPath:         "/",
			inputMethod:       "GET",
			expectResult:      uuid.Nil,
			expectError:       usecase.ErrAuthorizationFaild,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByTokenAndNotDeleted(ctx, gomock.Any()).
					Return(agent, nil).
					Times(1)
			},
			setMockAgentService: func(ctx context.Context, as *mockService.MockAgentService) {
				as.EXPECT().
					HasPermission(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(false, nil).
					Times(1)
			},
		},
		{
			name:              "agent not found",
			inputToken:        agentToken.Token,
			inputOperatorType: "AGENT",
			inputService:      "STORAGE",
			inputPath:         "/",
			inputMethod:       "GET",
			expectResult:      uuid.Nil,
			expectError:       usecase.ErrAuthenticationFailed,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByTokenAndNotDeleted(ctx, gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
			setMockAgentService: func(ctx context.Context, as *mockService.MockAgentService) {},
		},
		{
			name:              "find agent error",
			inputToken:        agentToken.Token,
			inputOperatorType: "AGENT",
			inputService:      "STORAGE",
			inputPath:         "/",
			inputMethod:       "GET",
			expectResult:      uuid.Nil,
			expectError:       sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByTokenAndNotDeleted(ctx, gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockAgentService: func(ctx context.Context, as *mockService.MockAgentService) {},
		},
		{
			name:              "permission decision error",
			inputToken:        agentToken.Token,
			inputOperatorType: "AGENT",
			inputService:      "STORAGE",
			inputPath:         "/",
			inputMethod:       "GET",
			expectResult:      uuid.Nil,
			expectError:       sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {},
			setMockAgentRepository: func(ctx context.Context, ar *mockRepository.MockAgentRepository) {
				ar.EXPECT().
					FindOneByTokenAndNotDeleted(ctx, gomock.Any()).
					Return(agent, nil).
					Times(1)
			},
			setMockAgentService: func(ctx context.Context, as *mockService.MockAgentService) {
				as.EXPECT().
					HasPermission(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(false, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:                       "invalid operator type",
			inputToken:                 userToken.Token,
			inputOperatorType:          "OPERATOR",
			inputService:               "STORAGE",
			inputPath:                  "/",
			inputMethod:                "GET",
			expectResult:               uuid.Nil,
			expectError:                usecase.ErrAuthenticationFailed,
			setMockTransactionObject:   func(ctx context.Context, to *mockDomain.MockTransactionObject) {},
			setMockUserTokenRepository: func(ctx context.Context, utr *mockRepository.MockUserTokenRepository) {},
			setMockAgentRepository:     func(ctx context.Context, ar *mockRepository.MockAgentRepository) {},
			setMockAgentService:        func(ctx context.Context, as *mockService.MockAgentService) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			to := mockDomain.NewMockTransactionObject(ctrl)
			utr := mockRepository.NewMockUserTokenRepository(ctrl)
			ar := mockRepository.NewMockAgentRepository(ctrl)
			as := mockService.NewMockAgentService(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockUserTokenRepository(ctx, utr)
			tt.setMockAgentRepository(ctx, ar)
			tt.setMockAgentService(ctx, as)

			au := usecase.NewAuthUsecase(to, nil, utr, ar, as)
			result, err := au.Authorize(ctx, tt.inputToken, tt.inputOperatorType, tt.inputService, tt.inputPath, tt.inputMethod)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if result != tt.expectResult {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectResult, result)
			}
		})
	}
}
