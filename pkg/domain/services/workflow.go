package services

import (
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"time"
)

type WorkflowService struct {
	redisRepo   repos.WorkflowRedisRepoInterface
	brokerRepo  repos.WorkflowBrokerRepository
	idGenerator IDService
}

func NewWorkflowService(repoRedis repos.WorkflowRedisRepoInterface, repoBroker repos.WorkflowBrokerRepository, idGenerator IDService) *WorkflowService {
	return &WorkflowService{
		redisRepo:   repoRedis,
		brokerRepo:  repoBroker,
		idGenerator: idGenerator,
	}
}

func (s *WorkflowService) retriesWorkflow(workflow *models.Workflow) (bool, bool) {
	workflow.UUID = s.idGenerator.GenerateWorkflowID()
	existWorkflowUUID := s.ValidateUUID(workflow)
	if existWorkflowUUID {
		return false, true
	}

	lockKey := "lock:" + workflow.UUID
	acquired, err := s.redisRepo.AcquireLock(lockKey, "", 30*time.Second)
	if err != nil {
		log.Printf("ERROR | acquiring lock: %v", err)
		return false, false
	}
	if !acquired {
		// Cannot acquire lock, maybe already exist UUID
		return false, true
	}

	defer s.redisRepo.RemoveLock(lockKey) // in case

	workflow.CreatedAt = time.Now().Format(models.LayoutTimestamp) // right now not controlled by db
	workflow.UpdatedAt = workflow.CreatedAt                        // right now not controlled by db
	workflow.IsActive = models.Active                              // right now not controlled by db
	workflow.Status = models.Initial                               // right now not controlled by db
	workflow.WorkflowInit = time.Time{}
	workflow.WorkflowCompleted = time.Time{}

	createdRedis, existRedis := s.redisRepo.Create(workflow)
	if existRedis {
		return false, true
	}
	if !createdRedis {
		return false, false
	}

	sended := s.brokerRepo.Create(workflow)
	if !sended {
		log.Printf("ERROR | Failed to publish workflow event: %v", err)
		s.redisRepo.Remove(workflow)
		return false, false
	}
	return createdRedis, existRedis
}

func (s *WorkflowService) CreateWorkflow(workflowFrontend *models.WorkflowFrontend) (created bool, exist bool, workflow *models.Workflow) {
	workflow = s.fromWorkflowFrontendToBackend(workflowFrontend)
	for i := 1; i < models.MaxAttempts; i++ {
		created, exist = s.retriesWorkflow(workflow)
		if !exist && created {
			return created, exist, workflow
		}
		if exist {
			return false, true, workflow
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("WARNING | Failed to create workflow, attempt %d:. Retrying in %v", i, waitTime)
		time.Sleep(waitTime)
	}
	log.Print("ERROR | Needs to add to Dead Letter. Cannot create workflow")
	// TODO: dead letter
	return false, false, workflow
}

func (s *WorkflowService) fromWorkflowFrontendToBackend(fw *models.WorkflowFrontend) *models.Workflow {
	return &models.Workflow{
		Name:            fw.WorkflowName,
		Description:     fw.Description,
		DirectoryToSave: fw.DirectoryToSave,
		UserID:          fw.Sub,
	}
}

func (s *WorkflowService) ValidateUUID(workflow *models.Workflow) bool {
	return s.redisRepo.ValidateUUID(workflow)
}
