package handler_test

import (
	"bytes"
	"holos-auth-api/internal/app/api/interface/handler"
	"holos-auth-api/internal/app/api/usecase/dto"
	"holos-auth-api/internal/pkg/apierr"
	mock_usecase "holos-auth-api/test/mock/usecase"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAgent_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		requestJSON          string
		resultDTO            *dto.AgentDTO
		resultError          apierr.ApiError
		expect               int
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			requestJSON:          `{"name": "name"}`,
			resultDTO:            dto.NewAgentDTO(uuid.New(), uuid.New(), "name", time.Now(), time.Now()),
			resultError:          nil,
			expect:               http.StatusOK,
		},
		{
			name:                 "invalid_request",
			isSetUserIDToContext: true,
			requestJSON:          "",
			resultDTO:            dto.NewAgentDTO(uuid.New(), uuid.New(), "name", time.Now(), time.Now()),
			resultError:          nil,
			expect:               http.StatusBadRequest,
		},
		{
			name:                 "context_does_not_have_user_id",
			isSetUserIDToContext: false,
			requestJSON:          `{"name": "name"}`,
			resultDTO:            dto.NewAgentDTO(uuid.New(), uuid.New(), "name", time.Now(), time.Now()),
			resultError:          nil,
			expect:               http.StatusInternalServerError,
		},
		{
			name:                 "result_error",
			isSetUserIDToContext: true,
			requestJSON:          `{"name": "name"}`,
			resultDTO:            nil,
			resultError:          apierr.NewApiError(http.StatusInternalServerError, "test error"),
			expect:               http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		req, err := http.NewRequest("POST", "/agent", bytes.NewBuffer([]byte(tt.requestJSON)))
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

		au := mock_usecase.NewMockAgentUsecase(ctrl)
		au.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.resultDTO, tt.resultError).AnyTimes()

		ah := handler.NewAgentHandler(au)
		ah.Create(ctx)

		if w.Code != tt.expect {
			t.Errorf("expect %d but got %d", tt.expect, w.Code)
		}
	}
}

func TestAgent_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		requestJSON          string
		resultDTO            *dto.AgentDTO
		resultError          apierr.ApiError
		expect               int
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			requestJSON:          `{"name": "name"}`,
			resultDTO:            dto.NewAgentDTO(uuid.New(), uuid.New(), "name", time.Now(), time.Now()),
			resultError:          nil,
			expect:               http.StatusOK,
		},
		{
			name:                 "invalid_request",
			isSetUserIDToContext: true,
			requestJSON:          "",
			resultDTO:            dto.NewAgentDTO(uuid.New(), uuid.New(), "name", time.Now(), time.Now()),
			resultError:          nil,
			expect:               http.StatusBadRequest,
		},
		{
			name:                 "context_does_not_have_user_id",
			isSetUserIDToContext: false,
			requestJSON:          `{"name": "name"}`,
			resultDTO:            dto.NewAgentDTO(uuid.New(), uuid.New(), "name", time.Now(), time.Now()),
			resultError:          nil,
			expect:               http.StatusInternalServerError,
		},
		{
			name:                 "result_error",
			isSetUserIDToContext: true,
			requestJSON:          `{"name": "name"}`,
			resultDTO:            nil,
			resultError:          apierr.NewApiError(http.StatusInternalServerError, "test error"),
			expect:               http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/agent/:id", bytes.NewBuffer([]byte(tt.requestJSON)))
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: uuid.New().String()})
			if tt.isSetUserIDToContext {
				ctx.Set("userID", uuid.New())
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			au := mock_usecase.NewMockAgentUsecase(ctrl)
			au.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.resultDTO, tt.resultError).AnyTimes()

			ah := handler.NewAgentHandler(au)
			ah.Update(ctx)

			if w.Code != tt.expect {
				t.Errorf("expect %d but got %d", tt.expect, w.Code)
			}
		})
	}
}

func TestAgent_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		resultDTO            *dto.AgentDTO
		resultError          apierr.ApiError
		expect               int
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			resultError:          nil,
			expect:               http.StatusOK,
		},
		{
			name:                 "context_does_not_have_user_id",
			isSetUserIDToContext: false,
			resultError:          nil,
			expect:               http.StatusInternalServerError,
		},
		{
			name:                 "result_error",
			isSetUserIDToContext: true,
			resultError:          apierr.NewApiError(http.StatusInternalServerError, "test error"),
			expect:               http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/agents/", nil)
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: uuid.New().String()})
			if tt.isSetUserIDToContext {
				ctx.Set("userID", uuid.New())
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			au := mock_usecase.NewMockAgentUsecase(ctrl)
			au.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.resultError).AnyTimes()

			ah := handler.NewAgentHandler(au)
			ah.Delete(ctx)

			if w.Code != tt.expect {
				t.Errorf("expect %d but got %d", tt.expect, w.Code)
			}
		})
	}
}