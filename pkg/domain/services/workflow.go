package services

import (
	"log"
	"math/rand"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"time"

	"github.com/google/uuid"
)

type WorkflowService struct {
	repo        repos.WorkflowRepository
	idGenerator IDService
}

func NewWorkflowService(repo repos.WorkflowRepository, idGenerator IDService) *WorkflowService {
	return &WorkflowService{
		repo:        repo,
		idGenerator: idGenerator,
	}
}

func (s *WorkflowService) CreateWorkflow(workflow *models.Workflow) (created bool, exist bool) {
	for i := 0; i < models.MaxAttempts; i++ {
		workflow.UUID = s.idGenerator.GenerateWorkflowID()
		created, exist := s.repo.Create(workflow)
		if created {
			workflow.CreatedAt = time.Now().Format(models.LayoutTimestamp)
			workflow.UpdatedAt = time.Now().Format(models.LayoutTimestamp)
			return true, false
		}
    if (exist) {
      return false, true
    }
		waitTime := time.Duration(rand.Int63n(int64(models.MaxSleepDuration-models.MinSleepDuration))) + models.MinSleepDuration + models.SleepOffset
		log.Printf("WARNING | Failed to create workflow, attempt %d:. Retrying in %v", i, waitTime)
		time.Sleep(waitTime)
	}
	log.Print("ERROR | Needs to add to Dead Letter. Cannot create workflow")
	// TODO: dead letter
	return false, false
}

func (s *WorkflowService) GetWorkflow(uuid uuid.UUID) (*models.Workflow, error) {
	return s.repo.GetByUUID(uuid)
}
