package handler_test

import (
	"bytes"
	"errors"
	"holos-auth-api/internal/app/api/interface/handler"
	"holos-auth-api/internal/app/api/usecase/dto"
	mock_usecase "holos-auth-api/test/mock/domain/usecase"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestUser_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name        string
		requestJSON string
		resultDTO   *dto.UserDTO
		resultError error
		expect      int
	}{
		{
			name:        "success",
			requestJSON: `{"name": "name", "password": "password", "confirm_password": "password"}`,
			resultDTO:   dto.NewUserDTO(uuid.New(), "name", "password", time.Now(), time.Now()),
			resultError: nil,
			expect:      http.StatusOK,
		},
		{
			name:        "invalid_request",
			requestJSON: "",
			resultDTO:   dto.NewUserDTO(uuid.New(), "name", "password", time.Now(), time.Now()),
			resultError: nil,
			expect:      http.StatusBadRequest,
		},
		{
			name:        "result_error",
			requestJSON: `{"name": "name", "password": "password", "confirm_password": "password"}`,
			resultDTO:   nil,
			resultError: errors.New("test error"),
			expect:      http.StatusInternalServerError,
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

			uu := mock_usecase.NewMockUserUsecase(ctrl)
			uu.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.resultDTO, tt.resultError).AnyTimes()

			uh := handler.NewUserHandler(uu)
			uh.Create(ctx)

			if w.Code != tt.expect {
				t.Errorf("expect %d but got %d", tt.expect, w.Code)
			}
		})
	}
}

func TestUser_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name        string
		requestJSON string
		resultDTO   *dto.UserDTO
		resultError error
		expect      int
	}{
		{
			name:        "success",
			requestJSON: `{"current_password": "password", "new_password": "new_password", "confirm_new_password": "new_password"}`,
			resultDTO:   dto.NewUserDTO(uuid.New(), "name", "password", time.Now(), time.Now()),
			resultError: nil,
			expect:      http.StatusOK,
		},
		{
			name:        "invalid_request",
			requestJSON: "",
			resultDTO:   dto.NewUserDTO(uuid.New(), "name", "password", time.Now(), time.Now()),
			resultError: nil,
			expect:      http.StatusBadRequest,
		},
		{
			name:        "result_error",
			requestJSON: `{"current_password": "password", "new_password": "new_password", "confirm_new_password": "new_password"}`,
			resultDTO:   nil,
			resultError: errors.New("test error"),
			expect:      http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/user/:name", bytes.NewBuffer([]byte(tt.requestJSON)))
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			ctx.Params = append(ctx.Params, gin.Param{Key: "name", Value: tt.name})

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uu := mock_usecase.NewMockUserUsecase(ctrl)
			uu.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.resultDTO, tt.resultError).AnyTimes()

			uh := handler.NewUserHandler(uu)
			uh.Update(ctx)

			if w.Code != tt.expect {
				t.Errorf("expect %d but got %d", tt.expect, w.Code)
			}
		})
	}
}

func TestUser_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name        string
		requestJSON string
		resultDTO   *dto.UserDTO
		resultError error
		expect      int
	}{
		{
			name:        "success",
			requestJSON: `{"password": "password"}`,
			resultError: nil,
			expect:      http.StatusOK,
		},
		{
			name:        "invalid_request",
			requestJSON: "",
			resultError: nil,
			expect:      http.StatusBadRequest,
		},
		{
			name:        "result_error",
			requestJSON: `{"password": "password"}`,
			resultError: errors.New("test error"),
			expect:      http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/user/:name", bytes.NewBuffer([]byte(tt.requestJSON)))
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			ctx.Params = append(ctx.Params, gin.Param{Key: "name", Value: tt.name})

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uu := mock_usecase.NewMockUserUsecase(ctrl)
			uu.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.resultError).AnyTimes()

			uh := handler.NewUserHandler(uu)
			uh.Delete(ctx)

			if w.Code != tt.expect {
				t.Errorf("expect %d but got %d", tt.expect, w.Code)
			}
		})
	}
}
