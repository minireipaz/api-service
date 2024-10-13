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
	httpRepo    repos.WorkflowHTTPRepository
}

func NewWorkflowService(repoRedis repos.WorkflowRedisRepoInterface, repoBroker repos.WorkflowBrokerRepository, idGenerator IDService, repoHTTP repos.WorkflowHTTPRepository) *WorkflowService {
	return &WorkflowService{
		redisRepo:   repoRedis,
		brokerRepo:  repoBroker,
		idGenerator: idGenerator,
		httpRepo:    repoHTTP,
	}
}

func (s *WorkflowService) retriesCreateWorkflow(workflow *models.Workflow) (bool, bool) {
	workflow.UUID = s.idGenerator.GenerateWorkflowID()
	existWorkflowUUID := s.ValidateUserWorkflowUUID(&workflow.UUID, &workflow.Name) // check individual
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
		log.Printf("ERROR | Failed to publish workflow event")
		s.redisRepo.Remove(workflow)
		return false, false
	}
	return createdRedis, existRedis
}

func (s *WorkflowService) CreateWorkflow(workflowFrontend *models.WorkflowFrontend) (created bool, exist bool, workflow *models.Workflow) {
	workflow = s.fromWorkflowFrontendToBackend(workflowFrontend)
	// validate if exist globally uuid
	exist = s.ValidateWorkflowGlobalUUID(&workflow.UUID)
	if exist {
		return false, exist, workflow
	}
	for i := 1; i < models.MaxAttempts; i++ {
		created, exist = s.retriesCreateWorkflow(workflow)
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

func (s *WorkflowService) GetWorkflow(userID, workflowID *string) (newWorkflow *models.Workflow, exist bool) {
	exist = s.ValidateWorkflowGlobalUUID(workflowID)
	if !exist {
		return nil, false
	}
	for i := 1; i < models.MaxAttempts; i++ {
		newWorkflow, exist = s.retriesGetWorkflow(userID, workflowID)
		if exist {
			return newWorkflow, exist
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("WARNING | Failed to get workflow, attempt %d:. Retrying in %v", i, waitTime)
		time.Sleep(waitTime)
	}
	log.Print("ERROR | Needs to add to Dead Letter. Cannot get workflow")
	// TODO: dead letter
	return nil, false
}

func (s *WorkflowService) UpdateWorkflow(workflow *models.Workflow) (updated bool, exist bool) {
	exist = s.ValidateWorkflowGlobalUUID(&workflow.UUID)
	if !exist {
		return false, exist
	}
	for i := 1; i < models.MaxAttempts; i++ {
		updated, exist = s.retriesUpdateWorkflow(workflow)
		if !exist {
			return updated, exist
		}
		if exist && updated {
			return updated, exist
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("WARNING | Failed to update workflow, attempt %d:. Retrying in %v", i, waitTime)
		time.Sleep(waitTime)
	}
	log.Print("ERROR | Needs to add to Dead Letter. Cannot update workflow")
	// TODO: dead letter
	return false, exist // exist but not updated
}

func (s *WorkflowService) fromWorkflowFrontendToBackend(fw *models.WorkflowFrontend) *models.Workflow {
	return &models.Workflow{
		Name:            fw.WorkflowName,
		Description:     fw.Description,
		DirectoryToSave: fw.DirectoryToSave,
		UserID:          fw.UserID,
	}
}

func (s *WorkflowService) ValidateWorkflowGlobalUUID(uuid *string) bool {
	return s.redisRepo.ValidateWorkflowGlobalUUID(uuid)
}

func (s *WorkflowService) ValidateUserWorkflowUUID(worklfowID, name *string) bool {
	return s.redisRepo.ValidateUserWorkflowUUID(worklfowID, name)
}

func (s *WorkflowService) retriesUpdateWorkflow(workflow *models.Workflow) (updated, exist bool) {
	exist = s.ValidateWorkflowGlobalUUID(&workflow.UUID)
	if !exist {
		return false, exist
	}

	lockKey := "lock:" + workflow.UUID
	// case lockkey cannot be created, it is assumed here the "10 seconds rate-limit" time has not passed
	acquired, err := s.redisRepo.AcquireLock(lockKey, "", models.RateLimitUpdate)
	if err != nil {
		log.Printf("ERROR | acquiring lock: %v", err)
		return false, false
	}
	if !acquired {
		// Cannot acquire lock, maybe already exist UUID
		return false, true
	}

	defer s.redisRepo.RemoveLock(lockKey) // in case

	workflow.UpdatedAt = time.Now().Format(models.LayoutTimestamp) // right now not controlled by db

	updated = s.brokerRepo.Update(workflow)
	if !updated {
		log.Printf("ERROR | Failed to publish workflow event")
		s.redisRepo.Remove(workflow)
		return false, false
	}
	return updated, exist
}

func (s *WorkflowService) retriesGetWorkflow(userID, workflowID *string) (newWorkflow *models.Workflow, exist bool) {
	exist = s.ValidateWorkflowGlobalUUID(workflowID) // not necessary validate global, better local validation
	if !exist {
		return nil, exist
	}

	lockKey := "lock:" + *workflowID
	// case lockkey cannot be created, it is assumed here the "10 seconds rate-limit" time has not passed
	acquired, err := s.redisRepo.AcquireLock(lockKey, "", models.RateLimitUpdate)
	if err != nil {
		log.Printf("ERROR | acquiring lock: %v", err)
		return nil, false
	}
	if !acquired {
		// Cannot acquire lock, maybe already exist UUID
		return nil, false
	}

	defer s.redisRepo.RemoveLock(lockKey) // in case

	// by default return array workflows but only need 1 row
	reponseWorkflow, err := s.httpRepo.GetWorkflowDataByID(userID, workflowID, 1) // limited to 1 row
	if err != nil {
		return nil, false
	}
	if len(reponseWorkflow.Data) > 1 || len(reponseWorkflow.Data) == 0 {
		return nil, false
	}

	// only need one row
	return &reponseWorkflow.Data[0], exist
}
