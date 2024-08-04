package redisclient_test

import (
	"context"
	"fmt"
	"testing"

	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/infra/redisclient"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRedisClient_checkAndSetWorkflow(t *testing.T) {
	ctx := context.Background()
	redisClient, mock := redismock.NewClientMock()

	r := &redisclient.RedisClient{
		Client: redisClient,
		Ctx:    ctx,
	}

	workflow := &models.Workflow{
		Sub:          "user123",
		WorkflowName: "test-workflow",
		UUID:         uuid.New(),
	}

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful transaction",
			setup: func() {
				mock.ExpectHExists("workflows:all", workflow.UUID.String()).SetVal(false)
				mock.ExpectHExists(fmt.Sprintf("users:%s", workflow.Sub), workflow.WorkflowName).SetVal(false)
				mock.ExpectTxPipeline()
				mock.ExpectHSet(fmt.Sprintf("users:%s", workflow.Sub), workflow.WorkflowName, workflow.UUID.String()).SetVal(1)
				mock.ExpectHSet("workflows:all", workflow.UUID.String(), workflow.Sub).SetVal(1)
				mock.ExpectTxPipelineExec()
			},
			wantErr: false,
		},
		{
			name: "UUID already exists",
			setup: func() {
				mock.ExpectHExists("workflows:all", workflow.UUID.String()).SetVal(true)
			},
			wantErr: true,
			errMsg:  "UUID already exists",
		},
		{
			name: "workflow name already exists for this user",
			setup: func() {
				mock.ExpectHExists("workflows:all", workflow.UUID.String()).SetVal(false)
				mock.ExpectHExists(fmt.Sprintf("users:%s", workflow.Sub), workflow.WorkflowName).SetVal(true)
			},
			wantErr: true,
			errMsg:  "workflow name already exists for this user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := r.Client.Watch(r.Ctx, func(tx *redis.Tx) error {
				return r.CheckAndSetWorkflow(ctx, tx, workflow)
			})
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
