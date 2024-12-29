package handler_test

import (
	"bytes"
	"database/sql"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/interface/handler"
	"holos-auth-api/internal/app/api/usecase/mapper"
	mockUsecase "holos-auth-api/test/mock/usecase"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestUser_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name             string
		requestJSON      string
		expectStatusCode int
		setMockUsecase   func(*mockUsecase.MockUserUsecase)
	}{
		{
			name:             "success",
			requestJSON:      `{"name": "name", "password": "password", "confirm_password": "password"}`,
			expectStatusCode: http.StatusCreated,
			setMockUsecase: func(u *mockUsecase.MockUserUsecase) {
				u.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mapper.ToUserDTO(user), nil).
					Times(1)
			},
		},
		{
			name:             "invalid request",
			requestJSON:      "",
			expectStatusCode: http.StatusBadRequest,
			setMockUsecase:   func(u *mockUsecase.MockUserUsecase) {},
		},
		{
			name:             "create error",
			requestJSON:      `{"name": "name", "password": "password", "confirm_password": "password"}`,
			expectStatusCode: http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockUserUsecase) {
				u.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/user", bytes.NewBuffer([]byte(tt.requestJSON)))
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockUserUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewUserHandler(u)
			h.Create(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestUser_UpdateName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		requestJSON          string
		expectStatusCode     int
		setMockUsecase       func(*mockUsecase.MockUserUsecase)
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			requestJSON:          `{"name": "name"}`,
			expectStatusCode:     http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockUserUsecase) {
				u.EXPECT().
					UpdateName(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mapper.ToUserDTO(user), nil).
					Times(1)
			},
		},
		{
			name:                 "no user id in context",
			isSetUserIDToContext: false,
			requestJSON:          `{"name": "name"}`,
			expectStatusCode:     http.StatusInternalServerError,
			setMockUsecase:       func(u *mockUsecase.MockUserUsecase) {},
		},
		{
			name:                 "invalid request",
			isSetUserIDToContext: true,
			requestJSON:          "",
			expectStatusCode:     http.StatusBadRequest,
			setMockUsecase:       func(u *mockUsecase.MockUserUsecase) {},
		},
		{
			name:                 "update error",
			isSetUserIDToContext: true,
			requestJSON:          `{"name": "name"}`,
			expectStatusCode:     http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockUserUsecase) {
				u.EXPECT().
					UpdateName(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/user/name", bytes.NewBuffer([]byte(tt.requestJSON)))
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			if tt.isSetUserIDToContext {
				ctx.Set("userID", user.ID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockUserUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewUserHandler(u)
			h.UpdateName(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestUser_UpdatePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		requestJSON          string
		expectStatusCode     int
		setMockUsecase       func(*mockUsecase.MockUserUsecase)
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			requestJSON:          `{"current_password": "password", "new_password": "new_password", "confirm_new_password": "new_password"}`,
			expectStatusCode:     http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockUserUsecase) {
				u.EXPECT().
					UpdatePassword(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mapper.ToUserDTO(user), nil).
					Times(1)
			},
		},
		{
			name:                 "no user id in context",
			isSetUserIDToContext: false,
			requestJSON:          `{"current_password": "password", "new_password": "new_password", "confirm_new_password": "new_password"}`,
			expectStatusCode:     http.StatusInternalServerError,
			setMockUsecase:       func(u *mockUsecase.MockUserUsecase) {},
		},
		{
			name:                 "invalid_request",
			isSetUserIDToContext: true,
			requestJSON:          "",
			expectStatusCode:     http.StatusBadRequest,
			setMockUsecase:       func(u *mockUsecase.MockUserUsecase) {},
		},
		{
			name:                 "update error",
			isSetUserIDToContext: true,
			requestJSON:          `{"current_password": "password", "new_password": "new_password", "confirm_new_password": "new_password"}`,
			expectStatusCode:     http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockUserUsecase) {
				u.EXPECT().
					UpdatePassword(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/user/password", bytes.NewBuffer([]byte(tt.requestJSON)))
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			if tt.isSetUserIDToContext {
				ctx.Set("userID", uuid.New())
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockUserUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewUserHandler(u)
			h.UpdatePassword(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestUser_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		requestJSON          string
		expectStatusCode     int
		setMockUsecase       func(*mockUsecase.MockUserUsecase)
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			requestJSON:          `{"password": "password"}`,
			expectStatusCode:     http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockUserUsecase) {
				u.EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:                 "no user id in context",
			isSetUserIDToContext: false,
			requestJSON:          `{"password": "new_password"}`,
			expectStatusCode:     http.StatusInternalServerError,
			setMockUsecase:       func(u *mockUsecase.MockUserUsecase) {},
		},
		{
			name:                 "invalid_request",
			isSetUserIDToContext: true,
			requestJSON:          "",
			expectStatusCode:     http.StatusBadRequest,
			setMockUsecase:       func(u *mockUsecase.MockUserUsecase) {},
		},
		{
			name:                 "result_error",
			isSetUserIDToContext: true,
			requestJSON:          `{"password": "password"}`,
			expectStatusCode:     http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockUserUsecase) {
				u.EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/user", bytes.NewBuffer([]byte(tt.requestJSON)))
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			if tt.isSetUserIDToContext {
				ctx.Set("userID", uuid.New())
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockUserUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewUserHandler(u)
			h.Delete(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("expect: %d but got: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}
