package services

import (
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"time"

	"github.com/google/uuid"
)

type WorkflowService struct {
	redisRepo   repos.WorkflowRedisRepoInterface
	idGenerator IDService
}

func NewWorkflowService(repoRedis repos.WorkflowRedisRepoInterface, idGenerator IDService) *WorkflowService {
	return &WorkflowService{
		redisRepo:   repoRedis,
		idGenerator: idGenerator,
	}
}

func (s *WorkflowService) CreateWorkflow(workflow *models.Workflow) (created bool, exist bool) {
	for i := 1; i < models.MaxAttempts; i++ {
		workflow.UUID = s.idGenerator.GenerateWorkflowID()
		created, exist := s.redisRepo.Create(workflow)
		if created {
			workflow.CreatedAt = time.Now().Format(models.LayoutTimestamp)
			workflow.UpdatedAt = time.Now().Format(models.LayoutTimestamp)
			return true, false
		}
		if exist {
			return false, true
		}
		waitTime := common.RandomDuration(models.MaxSleepDuration, models.MinSleepDuration, i)
		log.Printf("WARNING | Failed to create workflow, attempt %d:. Retrying in %v", i, waitTime)
		time.Sleep(waitTime)
	}
	log.Print("ERROR | Needs to add to Dead Letter. Cannot create workflow")
	// TODO: dead letter
	return false, false
}

func (s *WorkflowService) GetWorkflow(uuid uuid.UUID) (*models.Workflow, error) {
	return s.redisRepo.GetByUUID(uuid)
}
