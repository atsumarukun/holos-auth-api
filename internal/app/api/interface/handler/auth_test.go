package handler_test

import (
	"bytes"
	"holos-auth-api/internal/app/api/interface/handler"
	"holos-auth-api/internal/pkg/apierr"
	mock_usecase "holos-auth-api/test/mock/usecase"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestAuth_Signin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name         string
		requestJSON  string
		resultString string
		resultError  apierr.ApiError
		expect       int
	}{
		{
			name:         "success",
			requestJSON:  `{"user_name": "user_name", "password": "password"}`,
			resultString: "token",
			resultError:  nil,
			expect:       http.StatusOK,
		},
		{
			name:         "invalid_request",
			requestJSON:  "",
			resultString: "token",
			resultError:  nil,
			expect:       http.StatusBadRequest,
		},
		{
			name:         "result_error",
			requestJSON:  `{"name": "name", "password": "password", "confirm_password": "password"}`,
			resultString: "",
			resultError:  apierr.NewApiError(http.StatusInternalServerError, "test error"),
			expect:       http.StatusInternalServerError,
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

			au := mock_usecase.NewMockAuthUsecase(ctrl)
			au.EXPECT().Signin(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.resultString, tt.resultError).AnyTimes()

			ah := handler.NewAuthHandler(au)
			ah.Signin(ctx)

			if w.Code != tt.expect {
				t.Errorf("expect %d but got %d", tt.expect, w.Code)
			}
		})
	}
}

func TestAuth_Signout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name                string
		authorizationHeader string
		resultError         apierr.ApiError
		expect              int
	}{
		{
			name:                "success",
			authorizationHeader: "Bearer token",
			resultError:         nil,
			expect:              http.StatusOK,
		},
		{
			name:                "invalid_header",
			authorizationHeader: "",
			resultError:         nil,
			expect:              http.StatusUnauthorized,
		},
		{
			name:                "result_error",
			authorizationHeader: "Bearer token",
			resultError:         apierr.NewApiError(http.StatusInternalServerError, "test error"),
			expect:              http.StatusInternalServerError,
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

			au := mock_usecase.NewMockAuthUsecase(ctrl)
			au.EXPECT().Signout(gomock.Any(), gomock.Any()).Return(tt.resultError).AnyTimes()

			ah := handler.NewAuthHandler(au)
			ah.Signout(ctx)

			if w.Code != tt.expect {
				t.Errorf("expect %d but got %d", tt.expect, w.Code)
			}
		})
	}
}
