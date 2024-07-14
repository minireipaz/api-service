package redis

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"minireipaz/pkg/domain/models"
// )

// type WorkflowRepository struct {
// 	redisClient *RedisClient
// }

// func NewWorkflowRepository(redisClient *RedisClient) *WorkflowRepository {
// 	return &WorkflowRepository{redisClient: redisClient}
// }

// func (r *WorkflowRepository) CreateWorkflow(workflow models.Workflow) models.Workflow {
// 	ctx := context.Background()
// 	workflow.ID = r.generateID()
// 	workflowJSON, _ := json.Marshal(workflow)
// 	r.redisClient.Set(ctx, fmt.Sprintf("workflow:%d", workflow.ID), workflowJSON)
// 	return workflow
// }

// func (r *WorkflowRepository) GetWorkflowByID(id int) (models.Workflow, error) {
// 	ctx := context.Background()
// 	workflowJSON, err := r.redisClient.Get(ctx, fmt.Sprintf("workflow:%d", id))
// 	if err != nil {
// 		return models.Workflow{}, err
// 	}

// 	var workflow models.Workflow
// 	json.Unmarshal([]byte(workflowJSON), &workflow)
// 	return workflow, nil
// }

// func (r *WorkflowRepository) generateID() int {
// 	// Lógica para generar un ID único
// 	// Esto puede variar dependiendo de cómo deseas manejar los IDs
// 	return 1
// }
