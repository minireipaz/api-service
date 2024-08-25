package repos

import (
	"minireipaz/pkg/domain/models"
	"time"

	"github.com/google/uuid"
)

type WorkflowRedisRepoInterface interface {
	Create(workflow *models.Workflow) (created bool, exist bool)
	Remove(workflow *models.Workflow) (removed bool)
	ValidateUUID(workflow *models.Workflow) bool
	ValidateWorkflowName(workflow *models.Workflow) bool
	GetByUUID(id uuid.UUID) (*models.Workflow, error)
	AcquireLock(key, value string, expiration time.Duration) (locked bool, err error)
	RemoveLock(key string) bool
}

type WorkflowBrokerRepository interface {
	Create(workflow *models.Workflow) (sended bool)
}
