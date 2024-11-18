package service_test

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/service"
	mock_repository "holos-auth-api/test/mock/domain/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAgent_Exists(t *testing.T) {
	tests := []struct {
		name        string
		isReturnNil bool
		expect      bool
	}{
		{
			name:        "exists",
			isReturnNil: false,
			expect:      true,
		},
		{
			name:        "not_exists",
			isReturnNil: true,
			expect:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := entity.NewAgent(uuid.New(), tt.name)
			if err != nil {
				t.Error(err.Error())
			}

			res := agent
			if tt.isReturnNil {
				res = nil
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			ar := mock_repository.NewMockAgentRepository(ctrl)
			ar.EXPECT().FindOneByUserIDAndName(ctx, gomock.Any(), tt.name).Return(res, nil)

			as := service.NewAgentService(ar)
			exists, err := as.Exists(ctx, agent)
			if err != nil {
				t.Error(err.Error())
			}
			if exists != tt.expect {
				t.Errorf("expect %t but got %t", tt.expect, exists)
			}
		})
	}
}
