package middleware_test

import (
	"database/sql"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/interface/middleware"
	mockUsecase "holos-auth-api/test/mock/usecase"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAuth_Authenticate(t *testing.T) {
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
					Authenticate(gomock.Any(), gomock.Any()).
					Return(userToken.UserID, nil).
					Times(1)
			},
		},
		{
			name:                "unset header",
			authorizationHeader: "",
			expectStatusCode:    http.StatusUnauthorized,
			setMockUsecase:      func(u *mockUsecase.MockAuthUsecase) {},
		},
		{
			name:                "invalid header",
			authorizationHeader: userToken.Token,
			expectStatusCode:    http.StatusUnauthorized,
			setMockUsecase:      func(u *mockUsecase.MockAuthUsecase) {},
		},
		{
			name:                "authenticate error",
			authorizationHeader: "Bearer " + userToken.Token,
			expectStatusCode:    http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockAuthUsecase) {
				u.EXPECT().
					Authenticate(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, sql.ErrConnDone).
					Times(1)
			},
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

			u := mockUsecase.NewMockAuthUsecase(ctrl)
			tt.setMockUsecase(u)

			m := middleware.NewAuthMiddleware(u)
			m.Authenticate(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}
