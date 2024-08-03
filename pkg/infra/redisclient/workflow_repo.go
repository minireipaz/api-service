package redisclient

import (
	"encoding/json"
	"fmt"
	"minireipaz/pkg/domain/models"

	"github.com/google/uuid"
)

type WorkflowRepository struct {
	redisClient *RedisClient
}

func NewWorkflowRepository(redisClient *RedisClient) *WorkflowRepository {
	return &WorkflowRepository{redisClient: redisClient}
}

func (r *WorkflowRepository) Create(workflow *models.Workflow) (created bool, exist bool) {
	err := r.redisClient.WatchWorkflow(workflow)
	if err.Error() == models.WorkflowNameExist {
		return false, true
	}

	return false, false
}

func (r *WorkflowRepository) ValidateUUID(workflow *models.Workflow) bool {
	exist := r.redisClient.Hexists("workflows:all", workflow.UUID.String())
	return exist
}

func (r *WorkflowRepository) ValidateWorkflowName(workflow *models.Workflow) bool {
	exist := r.redisClient.Hexists(fmt.Sprintf("users:%s", workflow.Sub), workflow.WorkflowName)
	return exist
}

func (r *WorkflowRepository) GetByUUID(id uuid.UUID) (*models.Workflow, error) {
	workflowJSON, err := r.redisClient.Get(fmt.Sprintf("workflow:%d", id))
	if err != nil {
		return nil, err
	}

	var workflow models.Workflow
	err = json.Unmarshal([]byte(workflowJSON), &workflow)
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}
