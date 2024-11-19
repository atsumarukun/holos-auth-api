package middleware_test

import (
	"holos-auth-api/internal/app/api/interface/middleware"
	"holos-auth-api/internal/app/api/pkg/apierr"
	mock_usecase "holos-auth-api/test/mock/usecase"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAuth_Authenticate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		id                  uuid.UUID
		name                string
		authorizationHeader string
		resultError         apierr.ApiError
		expectCode          int
	}{
		{
			id:                  uuid.New(),
			name:                "success",
			authorizationHeader: "Bearer token",
			resultError:         nil,
			expectCode:          http.StatusOK,
		},
		{
			id:                  uuid.Nil,
			name:                "header_is_not_set",
			authorizationHeader: "",
			resultError:         nil,
			expectCode:          http.StatusUnauthorized,
		},
		{
			id:                  uuid.Nil,
			name:                "invalid_header",
			authorizationHeader: "token",
			resultError:         nil,
			expectCode:          http.StatusUnauthorized,
		},
		{
			id:                  uuid.Nil,
			name:                "result_error",
			authorizationHeader: "Bearer token",
			resultError:         apierr.NewApiError(http.StatusInternalServerError, "test error"),
			expectCode:          http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/users/name", nil)
			if err != nil {
				t.Error(err.Error())
			}
			if tt.authorizationHeader != "" {
				req.Header.Add("Authorization", tt.authorizationHeader)
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			au := mock_usecase.NewMockAuthUsecase(ctrl)
			au.EXPECT().GetUserID(gomock.Any(), gomock.Any()).Return(tt.id, tt.resultError).AnyTimes()

			am := middleware.NewAuthMiddleware(au)
			am.Authenticate(ctx)

			userID, _ := ctx.Get("userID")
			id, _ := userID.(uuid.UUID)
			if id != tt.id {
				t.Errorf("expect %s but got %s", tt.id.String(), id.String())
			}

			if w.Code != tt.expectCode {
				t.Errorf("expect %d but got %d", tt.expectCode, w.Code)
			}
		})
	}
}
