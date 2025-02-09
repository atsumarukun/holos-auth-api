package handler_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/interface/handler"
	"holos-auth-api/internal/app/api/usecase/dto"
	"holos-auth-api/internal/app/api/usecase/mapper"
	mockUsecase "holos-auth-api/test/mock/usecase"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestPolicy_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		requestJSON          string
		expectStatusCode     int
		setMockUsecase       func(*mockUsecase.MockPolicyUsecase)
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			requestJSON:          `{"name": "name", "effect": "ALLOW", "service": "STORAGE", "path": "/", "allowed_methods": ["GET"]}`,
			expectStatusCode:     http.StatusCreated,
			setMockUsecase: func(u *mockUsecase.MockPolicyUsecase) {
				u.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mapper.ToPolicyDTO(policy), nil).
					Times(1)
			},
		},
		{
			name:                 "no user id in context",
			isSetUserIDToContext: false,
			requestJSON:          `{"name": "name", "effect": "ALLOW", "service": "STORAGE", "path": "/", "allowed_methods": ["GET"]}`,
			expectStatusCode:     http.StatusInternalServerError,
			setMockUsecase:       func(u *mockUsecase.MockPolicyUsecase) {},
		},
		{
			name:                 "invalid request",
			isSetUserIDToContext: true,
			requestJSON:          "",
			expectStatusCode:     http.StatusBadRequest,
			setMockUsecase:       func(u *mockUsecase.MockPolicyUsecase) {},
		},
		{
			name:                 "create error",
			isSetUserIDToContext: true,
			requestJSON:          `{"name": "name", "effect": "ALLOW", "service": "STORAGE", "path": "/", "allowed_methods": ["GET"]}`,
			expectStatusCode:     http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockPolicyUsecase) {
				u.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
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
				ctx.Set("userID", policy.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockPolicyUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewPolicyHandler(u)
			h.Create(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestPolicy_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                   string
		isSetIDToPathParameter bool
		isSetUserIDToContext   bool
		requestJSON            string
		expectStatusCode       int
		setMockUsecase         func(*mockUsecase.MockPolicyUsecase)
	}{
		{
			name:                   "success",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			requestJSON:            `{"name": "name", "effect": "ALLOW", "service": "STORAGE", "path": "/", "allowed_methods": ["GET"]}`,
			expectStatusCode:       http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockPolicyUsecase) {
				u.EXPECT().
					Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mapper.ToPolicyDTO(policy), nil).
					Times(1)
			},
		},
		{
			name:                   "no id in path parameter",
			isSetIDToPathParameter: false,
			isSetUserIDToContext:   true,
			requestJSON:            `{"name": "name", "effect": "ALLOW", "service": "STORAGE", "path": "/", "allowed_methods": ["GET"]}`,
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockPolicyUsecase) {},
		},
		{
			name:                   "no user id in context",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   false,
			requestJSON:            `{"name": "name", "effect": "ALLOW", "service": "STORAGE", "path": "/", "allowed_methods": ["GET"]}`,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase:         func(u *mockUsecase.MockPolicyUsecase) {},
		},
		{
			name:                   "invalid request",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			requestJSON:            "",
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockPolicyUsecase) {},
		},
		{
			name:                   "update error",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			requestJSON:            `{"name": "name", "effect": "ALLOW", "service": "STORAGE", "path": "/", "allowed_methods": ["GET"]}`,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockPolicyUsecase) {
				u.EXPECT().
					Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
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
			if tt.isSetIDToPathParameter {
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: policy.ID.String()})
			}
			if tt.isSetUserIDToContext {
				ctx.Set("userID", policy.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockPolicyUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewPolicyHandler(u)
			h.Update(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestPolicy_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                   string
		isSetIDToPathParameter bool
		isSetUserIDToContext   bool
		expectStatusCode       int
		setMockUsecase         func(*mockUsecase.MockPolicyUsecase)
	}{
		{
			name:                   "success",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockPolicyUsecase) {
				u.EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:                   "no id in path parameter",
			isSetIDToPathParameter: false,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockPolicyUsecase) {},
		},
		{
			name:                   "no user id in context",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   false,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase:         func(u *mockUsecase.MockPolicyUsecase) {},
		},
		{
			name:                   "delete error",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockPolicyUsecase) {
				u.EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
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
			if tt.isSetIDToPathParameter {
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: policy.ID.String()})
			}
			if tt.isSetUserIDToContext {
				ctx.Set("userID", policy.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockPolicyUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewPolicyHandler(u)
			h.Delete(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestPolicy_Gets(t *testing.T) {
	gin.SetMode(gin.TestMode)

	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		expectStatusCode     int
		setMockUsecase       func(*mockUsecase.MockPolicyUsecase)
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			expectStatusCode:     http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockPolicyUsecase) {
				u.EXPECT().
					Gets(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*dto.PolicyDTO{mapper.ToPolicyDTO(policy)}, nil).
					Times(1)
			},
		},
		{
			name:                 "no user id in context",
			isSetUserIDToContext: false,
			expectStatusCode:     http.StatusInternalServerError,
			setMockUsecase:       func(u *mockUsecase.MockPolicyUsecase) {},
		},
		{
			name:                 "get error",
			isSetUserIDToContext: true,
			expectStatusCode:     http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockPolicyUsecase) {
				u.EXPECT().
					Gets(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
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
				ctx.Set("userID", policy.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockPolicyUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewPolicyHandler(u)
			h.Gets(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestPolicy_UpdateAgents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                   string
		isSetIDToPathParameter bool
		isSetUserIDToContext   bool
		requestJSON            string
		expectStatusCode       int
		setMockUsecase         func(*mockUsecase.MockPolicyUsecase)
	}{
		{
			name:                   "success",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			requestJSON:            fmt.Sprintf(`{"agent_ids": ["%s"]}`, agent.ID),
			expectStatusCode:       http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockPolicyUsecase) {
				u.EXPECT().
					UpdateAgents(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*dto.AgentDTO{mapper.ToAgentDTO(agent)}, nil).
					Times(1)
			},
		},
		{
			name:                   "no id in path parameter",
			isSetIDToPathParameter: false,
			isSetUserIDToContext:   true,
			requestJSON:            fmt.Sprintf(`{"agent_ids": ["%s"]}`, agent.ID),
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockPolicyUsecase) {},
		},
		{
			name:                   "no user id in context",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   false,
			requestJSON:            fmt.Sprintf(`{"agent_ids": ["%s"]}`, agent.ID),
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase:         func(u *mockUsecase.MockPolicyUsecase) {},
		},
		{
			name:                   "invalid request",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			requestJSON:            "",
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockPolicyUsecase) {},
		},
		{
			name:                   "update agents error",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			requestJSON:            fmt.Sprintf(`{"agent_ids": ["%s"]}`, agent.ID),
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockPolicyUsecase) {
				u.EXPECT().
					UpdateAgents(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
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
			if tt.isSetIDToPathParameter {
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: policy.ID.String()})
			}
			if tt.isSetUserIDToContext {
				ctx.Set("userID", policy.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockPolicyUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewPolicyHandler(u)
			h.UpdateAgents(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestPolicy_GetAgents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                   string
		isSetIDToPathParameter bool
		isSetUserIDToContext   bool
		expectStatusCode       int
		setMockUsecase         func(*mockUsecase.MockPolicyUsecase)
	}{
		{
			name:                   "success",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockPolicyUsecase) {
				u.EXPECT().
					GetAgents(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*dto.AgentDTO{mapper.ToAgentDTO(agent)}, nil).
					Times(1)
			},
		},
		{
			name:                   "no id in path parameter",
			isSetIDToPathParameter: false,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockPolicyUsecase) {},
		},
		{
			name:                   "no user id in context",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   false,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase:         func(u *mockUsecase.MockPolicyUsecase) {},
		},
		{
			name:                   "aget agents error",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockPolicyUsecase) {
				u.EXPECT().
					GetAgents(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
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
			if tt.isSetIDToPathParameter {
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: policy.ID.String()})
			}
			if tt.isSetUserIDToContext {
				ctx.Set("userID", policy.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockPolicyUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewPolicyHandler(u)
			h.GetAgents(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}
