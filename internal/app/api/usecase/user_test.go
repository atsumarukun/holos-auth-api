package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/usecase"
	"holos-auth-api/internal/app/api/usecase/dto"
	"holos-auth-api/internal/app/api/usecase/mapper"
	mockDomain "holos-auth-api/test/mock/domain"
	mockRepository "holos-auth-api/test/mock/domain/repository"
	mockService "holos-auth-api/test/mock/domain/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestUser_Create(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                     string
		inputName                string
		inputPassword            string
		inputConfirmPassword     string
		expectResult             *dto.UserDTO
		expectError              error
		setMockTransactionObject func(context.Context, *mockDomain.MockTransactionObject)
		setMockUserRepository    func(context.Context, *mockRepository.MockUserRepository)
		setMockUserService       func(context.Context, *mockService.MockUserService)
	}{
		{
			name:                 "success",
			inputName:            "name",
			inputPassword:        "password",
			inputConfirmPassword: "password",
			expectResult:         mapper.ToUserDTO(user),
			expectError:          nil,
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
					Create(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
			setMockUserService: func(ctx context.Context, us *mockService.MockUserService) {
				us.EXPECT().
					Exists(ctx, gomock.Any()).
					Return(false, nil).
					Times(1)
			},
		},
		{
			name:                     "invalid name",
			inputName:                "なまえ",
			inputPassword:            "password",
			inputConfirmPassword:     "password",
			expectResult:             nil,
			expectError:              entity.ErrInvalidUserName,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {},
			setMockUserRepository:    func(ctx context.Context, ur *mockRepository.MockUserRepository) {},
			setMockUserService:       func(ctx context.Context, us *mockService.MockUserService) {},
		},
		{
			name:                 "user already exists",
			inputName:            "name",
			inputPassword:        "password",
			inputConfirmPassword: "password",
			expectResult:         nil,
			expectError:          usecase.ErrUserAlreadyExists,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserRepository: func(ctx context.Context, ur *mockRepository.MockUserRepository) {},
			setMockUserService: func(ctx context.Context, us *mockService.MockUserService) {
				us.EXPECT().
					Exists(ctx, gomock.Any()).
					Return(true, nil).
					Times(1)
			},
		},
		{
			name:                 "existence check error",
			inputName:            "name",
			inputPassword:        "password",
			inputConfirmPassword: "password",
			expectResult:         nil,
			expectError:          sql.ErrConnDone,
			setMockTransactionObject: func(ctx context.Context, to *mockDomain.MockTransactionObject) {
				to.EXPECT().
					Transaction(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Times(1)
			},
			setMockUserRepository: func(ctx context.Context, ur *mockRepository.MockUserRepository) {},
			setMockUserService: func(ctx context.Context, us *mockService.MockUserService) {
				us.EXPECT().
					Exists(ctx, gomock.Any()).
					Return(false, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:                 "create error",
			inputName:            "name",
			inputPassword:        "password",
			inputConfirmPassword: "password",
			expectResult:         nil,
			expectError:          sql.ErrConnDone,
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
					Create(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
			setMockUserService: func(ctx context.Context, us *mockService.MockUserService) {
				us.EXPECT().
					Exists(ctx, gomock.Any()).
					Return(false, nil).
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
			us := mockService.NewMockUserService(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockUserRepository(ctx, ur)
			tt.setMockUserService(ctx, us)

			uu := usecase.NewUserUsecase(to, ur, us)
			result, err := uu.Create(ctx, tt.inputName, tt.inputPassword, tt.inputConfirmPassword)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.UserDTO{}, "ID", "Password", "CreatedAt", "UpdatedAt"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestUser_UpdateName(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                     string
		inputID                  uuid.UUID
		inputName                string
		expectResult             *dto.UserDTO
		expectError              error
		setMockTransactionObject func(context.Context, *mockDomain.MockTransactionObject)
		setMockUserRepository    func(context.Context, *mockRepository.MockUserRepository)
		setMockUserService       func(context.Context, *mockService.MockUserService)
	}{
		{
			name:         "success",
			inputID:      user.ID,
			inputName:    "update",
			expectResult: &dto.UserDTO{ID: user.ID, Name: "update", Password: user.Password, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt},
			expectError:  nil,
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
					Update(ctx, gomock.Any()).
					Return(nil).
					Times(1)
				ur.EXPECT().
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
					Times(1)
			},
			setMockUserService: func(ctx context.Context, us *mockService.MockUserService) {
				us.EXPECT().
					Exists(ctx, gomock.Any()).
					Return(false, nil).
					Times(1)
			},
		},
		{
			name:         "invalid name",
			inputID:      user.ID,
			inputName:    "なまえ",
			expectResult: nil,
			expectError:  entity.ErrInvalidUserName,
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
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
					Times(1)
			},
			setMockUserService: func(ctx context.Context, us *mockService.MockUserService) {},
		},
		{
			name:         "user not found",
			inputID:      user.ID,
			inputName:    "update",
			expectResult: nil,
			expectError:  usecase.ErrUserNotFound,
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
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(nil, nil).
					Times(1)
			},
			setMockUserService: func(ctx context.Context, us *mockService.MockUserService) {},
		},
		{
			name:         "user already exists",
			inputID:      user.ID,
			inputName:    "update",
			expectResult: nil,
			expectError:  usecase.ErrUserAlreadyExists,
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
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
					Times(1)
			},
			setMockUserService: func(ctx context.Context, us *mockService.MockUserService) {
				us.EXPECT().
					Exists(ctx, gomock.Any()).
					Return(true, nil).
					Times(1)
			},
		},
		{
			name:         "find user error",
			inputID:      user.ID,
			inputName:    "update",
			expectResult: nil,
			expectError:  sql.ErrConnDone,
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
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
			setMockUserService: func(ctx context.Context, us *mockService.MockUserService) {},
		},
		{
			name:         "existence check error",
			inputID:      user.ID,
			inputName:    "update",
			expectResult: nil,
			expectError:  sql.ErrConnDone,
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
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
					Times(1)
			},
			setMockUserService: func(ctx context.Context, us *mockService.MockUserService) {
				us.EXPECT().
					Exists(ctx, gomock.Any()).
					Return(false, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:         "update error",
			inputID:      user.ID,
			inputName:    "update",
			expectResult: nil,
			expectError:  sql.ErrConnDone,
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
					Update(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
				ur.EXPECT().
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
					Times(1)
			},
			setMockUserService: func(ctx context.Context, us *mockService.MockUserService) {
				us.EXPECT().
					Exists(ctx, gomock.Any()).
					Return(false, nil).
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
			us := mockService.NewMockUserService(ctrl)

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockUserRepository(ctx, ur)
			tt.setMockUserService(ctx, us)

			uu := usecase.NewUserUsecase(to, ur, us)
			result, err := uu.UpdateName(ctx, tt.inputID, tt.inputName)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.UserDTO{}, "UpdatedAt"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestUser_UpdatePassword(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                     string
		inputID                  uuid.UUID
		inputCurrentPassword     string
		inputNewPassword         string
		inputConfirmNewPassword  string
		expectResult             *dto.UserDTO
		expectError              error
		setMockTransactionObject func(context.Context, *mockDomain.MockTransactionObject)
		setMockUserRepository    func(context.Context, *mockRepository.MockUserRepository)
	}{
		{
			name:                    "success",
			inputID:                 user.ID,
			inputCurrentPassword:    "password",
			inputNewPassword:        "update_password",
			inputConfirmNewPassword: "update_password",
			expectResult:            mapper.ToUserDTO(user),
			expectError:             nil,
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
					Update(ctx, gomock.Any()).
					Return(nil).
					Times(1)
				ur.EXPECT().
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
					Times(1)
			},
		},
		{
			name:                    "invalid password",
			inputID:                 user.ID,
			inputCurrentPassword:    "password",
			inputNewPassword:        "ぱすわーど",
			inputConfirmNewPassword: "ぱすわーど",
			expectResult:            nil,
			expectError:             entity.ErrInvalidUserPassword,
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
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
					Times(1)
			},
		},
		{
			name:                    "verification failed",
			inputID:                 user.ID,
			inputCurrentPassword:    "PASSWORD",
			inputNewPassword:        "password",
			inputConfirmNewPassword: "password",
			expectResult:            nil,
			expectError:             entity.ErrAuthenticationFailed,
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
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
					Times(1)
			},
		},
		{
			name:                    "new password does not match",
			inputID:                 user.ID,
			inputCurrentPassword:    "password",
			inputNewPassword:        "password",
			inputConfirmNewPassword: "confirm_update",
			expectResult:            nil,
			expectError:             entity.ErrUserPasswordDoesNotMatch,
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
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
					Times(1)
			},
		},
		{
			name:                    "user not found",
			inputID:                 user.ID,
			inputCurrentPassword:    "password",
			inputNewPassword:        "password",
			inputConfirmNewPassword: "password",
			expectResult:            nil,
			expectError:             usecase.ErrUserNotFound,
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
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:                    "find user error",
			inputID:                 user.ID,
			inputCurrentPassword:    "password",
			inputNewPassword:        "password",
			inputConfirmNewPassword: "password",
			expectResult:            nil,
			expectError:             sql.ErrConnDone,
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
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:                    "update error",
			inputID:                 user.ID,
			inputCurrentPassword:    "password",
			inputNewPassword:        "password",
			inputConfirmNewPassword: "password",
			expectResult:            nil,
			expectError:             sql.ErrConnDone,
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
					Update(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
				ur.EXPECT().
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
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

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockUserRepository(ctx, ur)

			uu := usecase.NewUserUsecase(to, ur, nil)
			result, err := uu.UpdatePassword(
				ctx,
				tt.inputID,
				tt.inputCurrentPassword,
				tt.inputNewPassword,
				tt.inputConfirmNewPassword,
			)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			opts := cmp.Options{
				cmpopts.IgnoreFields(dto.UserDTO{}, "Password", "UpdatedAt"),
			}
			if diff := cmp.Diff(result, tt.expectResult, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestUser_Delete(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                     string
		inputID                  uuid.UUID
		inputPassword            string
		expectError              error
		setMockTransactionObject func(context.Context, *mockDomain.MockTransactionObject)
		setMockUserRepository    func(context.Context, *mockRepository.MockUserRepository)
	}{
		{
			name:          "success",
			inputID:       user.ID,
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
					Delete(ctx, gomock.Any()).
					Return(nil).
					Times(1)
				ur.EXPECT().
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
					Times(1)
			},
		},
		{
			name:          "user not found",
			inputID:       user.ID,
			inputPassword: "password",
			expectError:   usecase.ErrUserNotFound,
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
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name:          "verification failed",
			inputID:       user.ID,
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
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
					Times(1)
			},
		},
		{
			name:          "find user error",
			inputID:       user.ID,
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
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
		{
			name:          "delete error",
			inputID:       user.ID,
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
					Delete(ctx, gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
				ur.EXPECT().
					FindOneByIDAndNotDeleted(ctx, user.ID).
					Return(entity.RestoreUser(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil).
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

			ctx := context.Background()

			tt.setMockTransactionObject(ctx, to)
			tt.setMockUserRepository(ctx, ur)

			uu := usecase.NewUserUsecase(to, ur, nil)
			err := uu.Delete(ctx, tt.inputID, tt.inputPassword)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}
