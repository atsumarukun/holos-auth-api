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

func TestAgent_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		requestJSON          string
		expectStatusCode     int
		setMockUsecase       func(*mockUsecase.MockAgentUsecase)
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			requestJSON:          `{"name": "name"}`,
			expectStatusCode:     http.StatusCreated,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mapper.ToAgentDTO(agent), nil).
					Times(1)
			},
		},
		{
			name:                 "no user id in context",
			isSetUserIDToContext: false,
			requestJSON:          `{"name": "name"}`,
			expectStatusCode:     http.StatusInternalServerError,
			setMockUsecase:       func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                 "invalid request",
			isSetUserIDToContext: true,
			requestJSON:          "",
			expectStatusCode:     http.StatusBadRequest,
			setMockUsecase:       func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                 "create error",
			isSetUserIDToContext: true,
			requestJSON:          `{"name": "name"}`,
			expectStatusCode:     http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/agent", bytes.NewBuffer([]byte(tt.requestJSON)))
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			if tt.isSetUserIDToContext {
				ctx.Set("userID", agent.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockAgentUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewAgentHandler(u)
			h.Create(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestAgent_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

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
		setMockUsecase         func(*mockUsecase.MockAgentUsecase)
	}{
		{
			name:                   "success",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			requestJSON:            `{"name": "name"}`,
			expectStatusCode:       http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mapper.ToAgentDTO(agent), nil).
					Times(1)
			},
		},
		{
			name:                   "no id in path parameter",
			isSetIDToPathParameter: false,
			isSetUserIDToContext:   true,
			requestJSON:            `{"name": "name"}`,
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "no user id in context",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   false,
			requestJSON:            `{"name": "name"}`,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "invalid request",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			requestJSON:            "",
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "update error",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			requestJSON:            `{"name": "name"}`,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
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
			if tt.isSetIDToPathParameter {
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: agent.ID.String()})
			}
			if tt.isSetUserIDToContext {
				ctx.Set("userID", agent.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockAgentUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewAgentHandler(u)
			h.Update(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestAgent_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                   string
		isSetIDToPathParameter bool
		isSetUserIDToContext   bool
		expectStatusCode       int
		setMockUsecase         func(*mockUsecase.MockAgentUsecase)
	}{
		{
			name:                   "success",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
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
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "no user id in context",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   false,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "delete error",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/agents/:id", nil)
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			if tt.isSetIDToPathParameter {
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: agent.ID.String()})
			}
			if tt.isSetUserIDToContext {
				ctx.Set("userID", agent.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockAgentUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewAgentHandler(u)
			h.Delete(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestAgent_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                   string
		isSetIDToPathParameter bool
		isSetUserIDToContext   bool
		expectStatusCode       int
		setMockUsecase         func(*mockUsecase.MockAgentUsecase)
	}{
		{
			name:                   "success",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mapper.ToAgentDTO(agent), nil).
					Times(1)
			},
		},
		{
			name:                   "no id in path parameter",
			isSetIDToPathParameter: false,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "no user id in context",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   false,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "get error",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/agents/:id", nil)
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			if tt.isSetIDToPathParameter {
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: agent.ID.String()})
			}
			if tt.isSetUserIDToContext {
				ctx.Set("userID", agent.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockAgentUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewAgentHandler(u)
			h.Get(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestAgent_Gets(t *testing.T) {
	gin.SetMode(gin.TestMode)

	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                 string
		isSetUserIDToContext bool
		expectStatusCode     int
		setMockUsecase       func(*mockUsecase.MockAgentUsecase)
	}{
		{
			name:                 "success",
			isSetUserIDToContext: true,
			expectStatusCode:     http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					Gets(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*dto.AgentDTO{mapper.ToAgentDTO(agent)}, nil).
					Times(1)
			},
		},
		{
			name:                 "no user id in context",
			isSetUserIDToContext: false,
			expectStatusCode:     http.StatusInternalServerError,
			setMockUsecase:       func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                 "get error",
			isSetUserIDToContext: true,
			expectStatusCode:     http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					Gets(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/agents", nil)
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			if tt.isSetUserIDToContext {
				ctx.Set("userID", agent.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockAgentUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewAgentHandler(u)
			h.Gets(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestAgent_UpdatePolicies(t *testing.T) {
	gin.SetMode(gin.TestMode)

	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}
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
		setMockUsecase         func(*mockUsecase.MockAgentUsecase)
	}{
		{
			name:                   "success",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			requestJSON:            fmt.Sprintf(`{"policy_ids": ["%s"]}`, policy.ID),
			expectStatusCode:       http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					UpdatePolicies(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*dto.PolicyDTO{mapper.ToPolicyDTO(policy)}, nil).
					Times(1)
			},
		},
		{
			name:                   "no id in path parameter",
			isSetIDToPathParameter: false,
			isSetUserIDToContext:   true,
			requestJSON:            fmt.Sprintf(`{"policy_ids": ["%s"]}`, policy.ID),
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "no user id in context",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   false,
			requestJSON:            fmt.Sprintf(`{"policy_ids": ["%s"]}`, policy.ID),
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "invalid request",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			requestJSON:            "",
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "update policies error",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			requestJSON:            fmt.Sprintf(`{"policy_ids": ["%s"]}`, policy.ID),
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					UpdatePolicies(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/agents/:id/policies", bytes.NewBuffer([]byte(tt.requestJSON)))
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			if tt.isSetIDToPathParameter {
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: agent.ID.String()})
			}
			if tt.isSetUserIDToContext {
				ctx.Set("userID", agent.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockAgentUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewAgentHandler(u)
			h.UpdatePolicies(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestAgent_GetPolicies(t *testing.T) {
	gin.SetMode(gin.TestMode)

	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}
	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                   string
		isSetIDToPathParameter bool
		isSetUserIDToContext   bool
		expectStatusCode       int
		setMockUsecase         func(*mockUsecase.MockAgentUsecase)
	}{
		{
			name:                   "success",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					GetPolicies(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*dto.PolicyDTO{mapper.ToPolicyDTO(policy)}, nil).
					Times(1)
			},
		},
		{
			name:                   "no id in path parameter",
			isSetIDToPathParameter: false,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "no user id in context",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   false,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "get policies error",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					GetPolicies(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/agents/:id/policies", nil)
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			if tt.isSetIDToPathParameter {
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: agent.ID.String()})
			}
			if tt.isSetUserIDToContext {
				ctx.Set("userID", agent.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockAgentUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewAgentHandler(u)
			h.GetPolicies(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestAgent_GenerateToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}
	agentToken, err := entity.NewAgentToken(agent.ID)
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                   string
		isSetIDToPathParameter bool
		isSetUserIDToContext   bool
		expectStatusCode       int
		setMockUsecase         func(*mockUsecase.MockAgentUsecase)
	}{
		{
			name:                   "success",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					GenerateToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(agentToken.Token, nil).
					Times(1)
			},
		},
		{
			name:                   "no id in path parameter",
			isSetIDToPathParameter: false,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "no user id in context",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   false,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "generate token error",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					GenerateToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/agents/:id/token", nil)
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			if tt.isSetIDToPathParameter {
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: agent.ID.String()})
			}
			if tt.isSetUserIDToContext {
				ctx.Set("userID", agent.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockAgentUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewAgentHandler(u)
			h.GenerateToken(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestAgent_DeleteToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                   string
		isSetIDToPathParameter bool
		isSetUserIDToContext   bool
		expectStatusCode       int
		setMockUsecase         func(*mockUsecase.MockAgentUsecase)
	}{
		{
			name:                   "success",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					DeleteToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:                   "no id in path parameter",
			isSetIDToPathParameter: false,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "no user id in context",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   false,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "delete token error",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					DeleteToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/agents/:id/token", nil)
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			if tt.isSetIDToPathParameter {
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: agent.ID.String()})
			}
			if tt.isSetUserIDToContext {
				ctx.Set("userID", agent.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockAgentUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewAgentHandler(u)
			h.DeleteToken(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}

func TestAgent_GetToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}
	agentToken, err := entity.NewAgentToken(agent.ID)
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                   string
		isSetIDToPathParameter bool
		isSetUserIDToContext   bool
		expectStatusCode       int
		setMockUsecase         func(*mockUsecase.MockAgentUsecase)
	}{
		{
			name:                   "success",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusOK,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					GetToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(mapper.ToAgentTokenDTO(agentToken), nil).
					Times(1)
			},
		},
		{
			name:                   "no id in path parameter",
			isSetIDToPathParameter: false,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusBadRequest,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "no user id in context",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   false,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase:         func(u *mockUsecase.MockAgentUsecase) {},
		},
		{
			name:                   "get token error",
			isSetIDToPathParameter: true,
			isSetUserIDToContext:   true,
			expectStatusCode:       http.StatusInternalServerError,
			setMockUsecase: func(u *mockUsecase.MockAgentUsecase) {
				u.EXPECT().
					GetToken(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/agents/:id/token", nil)
			if err != nil {
				t.Error(err.Error())
			}
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			if tt.isSetIDToPathParameter {
				ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: agent.ID.String()})
			}
			if tt.isSetUserIDToContext {
				ctx.Set("userID", agent.UserID)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := mockUsecase.NewMockAgentUsecase(ctrl)
			tt.setMockUsecase(u)

			h := handler.NewAgentHandler(u)
			h.GetToken(ctx)

			if w.Code != tt.expectStatusCode {
				t.Errorf("\nexpect: %d \ngot: %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}
