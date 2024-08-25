package repos

import (
	"minireipaz/pkg/domain/models"

	"github.com/google/uuid"
)

type WorkflowRedisRepoInterface interface {
	Create(workflow *models.Workflow) (created bool, exist bool)
	ValidateUUID(workflow *models.Workflow) bool
	ValidateWorkflowName(workflow *models.Workflow) bool
	GetByUUID(id uuid.UUID) (*models.Workflow, error)
}
