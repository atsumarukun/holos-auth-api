package handler_test

import (
	"bytes"
	"fmt"
	"holos-auth-api/internal/app/api/interface/handler"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"holos-auth-api/internal/app/api/usecase/dto"
	mock_usecase "holos-auth-api/test/mock/usecase"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestPolicy_Create(t *testing.T) {
	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		requestJSON          string
		resultDTO            *dto.PolicyDTO
		resultError          apierr.ApiError
		expect               int
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			requestJSON:          `{"name": "name", "service": "STORAGE", "path": "/", "allowed_methods": ["GET"]}`,
			resultDTO:            dto.NewPolicyDTO(uuid.New(), uuid.New(), "name", "STORAGE", "/", []string{"GET"}, time.Now(), time.Now()),
			resultError:          nil,
			expect:               http.StatusCreated,
		},
		{
			name:                 "invalid_request",
			isSetUserIDToContext: true,
			requestJSON:          "",
			resultDTO:            dto.NewPolicyDTO(uuid.New(), uuid.New(), "name", "STORAGE", "/", []string{"GET"}, time.Now(), time.Now()),
			resultError:          nil,
			expect:               http.StatusBadRequest,
		},
		{
			name:                 "context_does_not_have_user_id",
			isSetUserIDToContext: false,
			requestJSON:          `{"name": "name", "service": "STORAGE", "path": "/", "allowed_methods": ["GET"]}`,
			resultDTO:            dto.NewPolicyDTO(uuid.New(), uuid.New(), "name", "STORAGE", "/", []string{"GET"}, time.Now(), time.Now()),
			resultError:          nil,
			expect:               http.StatusInternalServerError,
		},
		{
			name:                 "result_error",
			isSetUserIDToContext: true,
			requestJSON:          `{"name": "name", "service": "STORAGE", "path": "/", "allowed_methods": ["GET"]}`,
			resultDTO:            nil,
			resultError:          apierr.NewApiError(http.StatusInternalServerError, "test error"),
			expect:               http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/policies", bytes.NewBuffer([]byte(tt.requestJSON)))
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

			pu := mock_usecase.NewMockPolicyUsecase(ctrl)
			pu.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.resultDTO, tt.resultError).AnyTimes()

			ph := handler.NewPolicyHandler(pu)
			ph.Create(ctx)

			if w.Code != tt.expect {
				t.Errorf("expect %d but got %d", tt.expect, w.Code)
			}
		})
	}
}

func TestPolicy_Update(t *testing.T) {
	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		requestJSON          string
		resultDTO            *dto.PolicyDTO
		resultError          apierr.ApiError
		expect               int
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			requestJSON:          `{"name": "name", "service": "STORAGE", "path": "/", "allowed_methods": ["GET"]}`,
			resultDTO:            dto.NewPolicyDTO(uuid.New(), uuid.New(), "name", "STORAGE", "/", []string{"GET"}, time.Now(), time.Now()),
			resultError:          nil,
			expect:               http.StatusOK,
		},
		{
			name:                 "invalid_request",
			isSetUserIDToContext: true,
			requestJSON:          "",
			resultDTO:            dto.NewPolicyDTO(uuid.New(), uuid.New(), "name", "STORAGE", "/", []string{"GET"}, time.Now(), time.Now()),
			resultError:          nil,
			expect:               http.StatusBadRequest,
		},
		{
			name:                 "context_does_not_have_user_id",
			isSetUserIDToContext: false,
			requestJSON:          `{"name": "name", "service": "STORAGE", "path": "/", "allowed_methods": ["GET"]}`,
			resultDTO:            dto.NewPolicyDTO(uuid.New(), uuid.New(), "name", "STORAGE", "/", []string{"GET"}, time.Now(), time.Now()),
			resultError:          nil,
			expect:               http.StatusInternalServerError,
		},
		{
			name:                 "result_error",
			isSetUserIDToContext: true,
			requestJSON:          `{"name": "name", "service": "STORAGE", "path": "/", "allowed_methods": ["GET"]}`,
			resultDTO:            nil,
			resultError:          apierr.NewApiError(http.StatusInternalServerError, "test error"),
			expect:               http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/policies/:id", bytes.NewBuffer([]byte(tt.requestJSON)))
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

			pu := mock_usecase.NewMockPolicyUsecase(ctrl)
			pu.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.resultDTO, tt.resultError).AnyTimes()

			ph := handler.NewPolicyHandler(pu)
			ph.Update(ctx)

			if w.Code != tt.expect {
				t.Errorf("expect %d but got %d", tt.expect, w.Code)
			}
		})
	}
}

func TestPolicy_Delete(t *testing.T) {
	tests := []struct {
		name                 string
		isSetUserIDToContext bool
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
			req, err := http.NewRequest("DELETE", "/policies/:id", nil)
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

			pu := mock_usecase.NewMockPolicyUsecase(ctrl)
			pu.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.resultError).AnyTimes()

			ph := handler.NewPolicyHandler(pu)
			ph.Delete(ctx)

			if w.Code != tt.expect {
				t.Errorf("expect %d but got %d", tt.expect, w.Code)
			}
		})
	}
}

func TestPolicy_Gets(t *testing.T) {
	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		resultDTO            []*dto.PolicyDTO
		resultError          apierr.ApiError
		expect               int
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			resultDTO:            []*dto.PolicyDTO{dto.NewPolicyDTO(uuid.New(), uuid.New(), "name", "STORAGE", "/", []string{"GET"}, time.Now(), time.Now())},
			resultError:          nil,
			expect:               http.StatusOK,
		},
		{
			name:                 "context_does_not_have_user_id",
			isSetUserIDToContext: false,
			resultDTO:            []*dto.PolicyDTO{dto.NewPolicyDTO(uuid.New(), uuid.New(), "name", "STORAGE", "/", []string{"GET"}, time.Now(), time.Now())},
			resultError:          nil,
			expect:               http.StatusInternalServerError,
		},
		{
			name:                 "result_error",
			isSetUserIDToContext: true,
			resultDTO:            nil,
			resultError:          apierr.NewApiError(http.StatusInternalServerError, "test error"),
			expect:               http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/policies", nil)
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

			pu := mock_usecase.NewMockPolicyUsecase(ctrl)
			pu.EXPECT().Gets(gomock.Any(), gomock.Any()).Return(tt.resultDTO, tt.resultError).AnyTimes()

			ph := handler.NewPolicyHandler(pu)
			ph.Gets(ctx)

			if w.Code != tt.expect {
				t.Errorf("expect %d but got %d", tt.expect, w.Code)
			}
		})
	}
}

func TestPolicy_UpdateAgents(t *testing.T) {
	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		requestJSON          string
		resultDTO            []*dto.AgentDTO
		resultError          apierr.ApiError
		expect               int
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			requestJSON:          fmt.Sprintf(`{"agent_ids": ["%s"]}`, uuid.New()),
			resultDTO:            []*dto.AgentDTO{dto.NewAgentDTO(uuid.New(), uuid.New(), "name", time.Now(), time.Now())},
			resultError:          nil,
			expect:               http.StatusOK,
		},
		{
			name:                 "invalid_request",
			isSetUserIDToContext: true,
			requestJSON:          "",
			resultDTO:            []*dto.AgentDTO{dto.NewAgentDTO(uuid.New(), uuid.New(), "name", time.Now(), time.Now())},
			resultError:          nil,
			expect:               http.StatusBadRequest,
		},
		{
			name:                 "context_does_not_have_user_id",
			isSetUserIDToContext: false,
			requestJSON:          fmt.Sprintf(`{"agent_ids": ["%s"]}`, uuid.New()),
			resultDTO:            []*dto.AgentDTO{dto.NewAgentDTO(uuid.New(), uuid.New(), "name", time.Now(), time.Now())},
			resultError:          nil,
			expect:               http.StatusInternalServerError,
		},
		{
			name:                 "result_error",
			isSetUserIDToContext: true,
			requestJSON:          fmt.Sprintf(`{"agent_ids": ["%s"]}`, uuid.New()),
			resultDTO:            nil,
			resultError:          apierr.NewApiError(http.StatusInternalServerError, "test error"),
			expect:               http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/policies/:id/agents", bytes.NewBuffer([]byte(tt.requestJSON)))
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

			pu := mock_usecase.NewMockPolicyUsecase(ctrl)
			pu.EXPECT().UpdateAgents(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.resultDTO, tt.resultError).AnyTimes()

			ah := handler.NewPolicyHandler(pu)
			ah.UpdateAgents(ctx)

			if w.Code != tt.expect {
				t.Errorf("expect %d but got %d", tt.expect, w.Code)
			}
		})
	}
}

func TestPolicy_GetAgents(t *testing.T) {
	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		resultDTO            []*dto.AgentDTO
		resultError          apierr.ApiError
		expect               int
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			resultDTO:            []*dto.AgentDTO{dto.NewAgentDTO(uuid.New(), uuid.New(), "name", time.Now(), time.Now())},
			resultError:          nil,
			expect:               http.StatusOK,
		},
		{
			name:                 "context_does_not_have_user_id",
			isSetUserIDToContext: false,
			resultDTO:            []*dto.AgentDTO{dto.NewAgentDTO(uuid.New(), uuid.New(), "name", time.Now(), time.Now())},
			resultError:          nil,
			expect:               http.StatusInternalServerError,
		},
		{
			name:                 "result_error",
			isSetUserIDToContext: true,
			resultDTO:            nil,
			resultError:          apierr.NewApiError(http.StatusInternalServerError, "test error"),
			expect:               http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/policies/:id/agents", nil)
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

			pu := mock_usecase.NewMockPolicyUsecase(ctrl)
			pu.EXPECT().GetAgents(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.resultDTO, tt.resultError).AnyTimes()

			ph := handler.NewPolicyHandler(pu)
			ph.GetAgents(ctx)

			if w.Code != tt.expect {
				t.Errorf("expect %d but got %d", tt.expect, w.Code)
			}
		})
	}
}
