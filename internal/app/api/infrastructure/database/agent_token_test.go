package database_test

import (
	"context"
	"database/sql"
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/infrastructure/database"
	"holos-auth-api/test"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestAgentToken_Save(t *testing.T) {
	agentToken, err := entity.NewAgentToken(uuid.New())
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name            string
		inputAgentToken *entity.AgentToken
		expectError     error
		setMockDB       func(sqlmock.Sqlmock)
	}{
		{
			name:            "success",
			inputAgentToken: agentToken,
			expectError:     nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("REPLACE agent_tokens (agent_id, token, generated_at) VALUES (?, ?, ?);")).
					WithArgs(agentToken.AgentID, agentToken.Token, agentToken.GeneratedAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:            "save error",
			inputAgentToken: agentToken,
			expectError:     sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("REPLACE agent_tokens (agent_id, token, generated_at) VALUES (?, ?, ?);")).
					WithArgs(agentToken.AgentID, agentToken.Token, agentToken.GeneratedAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:            "no agent token",
			inputAgentToken: nil,
			expectError:     database.ErrRequiredAgentToken,
			setMockDB:       func(mock sqlmock.Sqlmock) {},
		},
	}
	for _, tt := range tests {
		db, mock := test.NewMockDB(t)
		defer db.Close()

		ctx := context.Background()

		tt.setMockDB(mock)

		r := database.NewAgentTokenDBRepository(db)
		if err := r.Save(ctx, tt.inputAgentToken); !errors.Is(err, tt.expectError) {
			t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err.Error())
		}
	}
}

func TestAgentToken_Delete(t *testing.T) {
	agentToken, err := entity.NewAgentToken(uuid.New())
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name            string
		inputAgentToken *entity.AgentToken
		expectError     error
		setMockDB       func(sqlmock.Sqlmock)
	}{
		{
			name:            "success",
			inputAgentToken: agentToken,
			expectError:     nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM agent_tokens WHERE agent_id = ?;")).
					WithArgs(agentToken.AgentID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:            "delete error",
			inputAgentToken: agentToken,
			expectError:     sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM agent_tokens WHERE agent_id = ?;")).
					WithArgs(agentToken.AgentID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:            "no agent token",
			inputAgentToken: nil,
			expectError:     database.ErrRequiredAgentToken,
			setMockDB:       func(mock sqlmock.Sqlmock) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewAgentTokenDBRepository(db)
			if err := r.Delete(ctx, tt.inputAgentToken); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestAgentToken_FindOneByToken(t *testing.T) {
	agentToken, err := entity.NewAgentToken(uuid.New())
	if err != nil {
		t.Error(err.Error())
	}
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputAgentID uuid.UUID
		inputUserID  uuid.UUID
		expectResult *entity.AgentToken
		expectError  error
		setMockDB    func(sqlmock.Sqlmock)
	}{
		{
			name:         "found",
			inputAgentID: agentToken.AgentID,
			inputUserID:  agent.UserID,
			expectResult: agentToken,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT
						agent_tokens.agent_id,
						agent_tokens.token,
						agent_tokens.generated_at
					FROM
						agent_tokens
						INNER JOIN agents ON agent_tokens.agent_id = agents.id
					WHERE
						agent_tokens.agent_id = ?
						AND agents.user_id = ?
					LIMIT 1;`,
				)).
					WithArgs(agentToken.AgentID, agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"agent_id", "token", "generated_at"}).
							AddRow(agentToken.AgentID, agentToken.Token, agentToken.GeneratedAt),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputAgentID: agentToken.AgentID,
			inputUserID:  agent.UserID,
			expectResult: nil,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT
						agent_tokens.agent_id,
						agent_tokens.token,
						agent_tokens.generated_at
					FROM
						agent_tokens
						INNER JOIN agents ON agent_tokens.agent_id = agents.id
					WHERE
						agent_tokens.agent_id = ?
						AND agents.user_id = ?
					LIMIT 1;`,
				)).
					WithArgs(agentToken.AgentID, agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"agent_id", "token", "generated_at"}).
							AddRow(agentToken.AgentID, agentToken.Token, agentToken.GeneratedAt),
					).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:         "find error",
			inputAgentID: agentToken.AgentID,
			inputUserID:  agent.UserID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`SELECT
						agent_tokens.agent_id,
						agent_tokens.token,
						agent_tokens.generated_at
					FROM
						agent_tokens
						INNER JOIN agents ON agent_tokens.agent_id = agents.id
					WHERE
						agent_tokens.agent_id = ?
						AND agents.user_id = ?
					LIMIT 1;`,
				)).
					WithArgs(agentToken.AgentID, agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"agent_id", "token", "generated_at"}).
							AddRow(agentToken.AgentID, agentToken.Token, agentToken.GeneratedAt),
					).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewAgentTokenDBRepository(db)
			result, err := r.FindOneByAgentIDAndUserID(ctx, tt.inputAgentID, tt.inputUserID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}
