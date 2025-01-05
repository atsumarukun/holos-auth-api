package handler_test

import (
	"bytes"
	"database/sql"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/interface/handler"
	mockUsecase "holos-auth-api/test/mock/usecase"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAuth_Signin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userToken, err := entity.NewUserToken(uuid.New())
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name             string
		requestJSON      string
		expectStatusCode int
		setMockUsecase   func(*mockUsecase.MockAuthUsecase)
	}{
		{
			name:             "success",
			requestJSON:      `{"user_name": "user_name", "password": "password"}`,
			expectStatusCode: http.StatusCreated,
			setMockUsecase: func(u *mockUsecase.MockAuthUsecase) {
				u.EXPECT().
					Signin(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(userToken.Token, nil).
					Times(1)
			},
		},
		{
			name:             "invalid request",
			requestJSON:      "",
			expectStatusCode: http.StatusBadRequest,
			setMockUsecase:   func(u *mockUsecase.MockAuthUsecase) {},
		},
		{
			name:             "signin error",
			requestJSON:      `{"user_name": "user_name", "password": "password"}`,
			expectStatusCode: http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockAuthUsecase) {
				u.EXPECT().
					Signin(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/auth/signin", bytes.NewBuffer([]byte(tt.requestJSON)))
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockAuthUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewAuthHandler(u)
			h.Signin(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestAuth_Signout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userToken, err := entity.NewUserToken(uuid.New())
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                string
		authorizationHeader string
		expectStatusCode    int
		setMockUsecase      func(*mockUsecase.MockAuthUsecase)
	}{
		{
			name:                "success",
			authorizationHeader: "Bearer " + userToken.Token,
			expectStatusCode:    http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockAuthUsecase) {
				u.EXPECT().
					Signout(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:                "invalid header",
			authorizationHeader: "",
			expectStatusCode:    http.StatusUnauthorized,
			setMockUsecase:      func(u *mockUsecase.MockAuthUsecase) {},
		},
		{
			name:                "signout error",
			authorizationHeader: "Bearer " + userToken.Token,
			expectStatusCode:    http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockAuthUsecase) {
				u.EXPECT().
					Signout(gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/auth/signout", nil)
			if err != nil {
				t.Error(err.Error())
			}
			req.Header.Add("Authorization", tt.authorizationHeader)
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockAuthUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewAuthHandler(u)
			h.Signout(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestAuth_Authorize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userToken, err := entity.NewUserToken(uuid.New())
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                string
		authorizationHeader string
		expectStatusCode    int
		setMockUsecase      func(*mockUsecase.MockAuthUsecase)
	}{
		{
			name:                "success",
			authorizationHeader: "Bearer " + userToken.Token,
			expectStatusCode:    http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockAuthUsecase) {
				u.EXPECT().
					Authorize(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(userToken.UserID, nil).
					Times(1)
			},
		},
		{
			name:                "invalid header",
			authorizationHeader: "",
			expectStatusCode:    http.StatusUnauthorized,
			setMockUsecase:      func(u *mockUsecase.MockAuthUsecase) {},
		},
		{
			name:                "signout error",
			authorizationHeader: "Bearer " + userToken.Token,
			expectStatusCode:    http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockAuthUsecase) {
				u.EXPECT().
					Authorize(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(uuid.Nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/auth/user_id", nil)
			if err != nil {
				t.Error(err.Error())
			}
			req.Header.Add("Authorization", tt.authorizationHeader)
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockAuthUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewAuthHandler(u)
			h.Authorize(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}
