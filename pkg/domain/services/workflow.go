package services

import (
	"context"
	"fmt"
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"time"
)

type WorkflowServiceImpl struct {
	redisRepo   repos.WorkflowRedisRepoInterface
	brokerRepo  repos.WorkflowBrokerRepository
	idGenerator IDService
	httpRepo    repos.WorkflowHTTPRepository
}

func NewWorkflowService(repoRedis repos.WorkflowRedisRepoInterface, repoBroker repos.WorkflowBrokerRepository, idGenerator IDService, repoHTTP repos.WorkflowHTTPRepository) repos.WorkflowService {
	return &WorkflowServiceImpl{
		redisRepo:   repoRedis,
		brokerRepo:  repoBroker,
		idGenerator: idGenerator,
		httpRepo:    repoHTTP,
	}
}

func (s *WorkflowServiceImpl) CreateWorkflow(workflowFrontend *models.WorkflowFrontend) (created bool, exist bool, workflow *models.Workflow) {
	// setup initial data
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

func (s *WorkflowServiceImpl) retriesCreateWorkflow(workflow *models.Workflow) (bool, bool) {
	workflow.UUID = s.idGenerator.GenerateWorkflowID()
	existWorkflowUUID := s.ValidateUserWorkflowUUID(&workflow.UUID, &workflow.Name) // check individual
	if existWorkflowUUID {
		return false, true
	}

	lockKey := "lock:" + workflow.UUID
	acquired, err := s.redisRepo.AcquireLock(lockKey, "_", models.MaxTimeForLocks)
	if err != nil {
		log.Printf("ERROR | acquiring lock: %v", err)
		return false, false
	}
	if !acquired {
		// Cannot acquire lock, maybe already exist UUID
		return false, true
	}

	defer s.redisRepo.RemoveLock(lockKey) // in case

	workflow.CreatedAt = time.Now().UTC().Format(models.LayoutTimestamp) // right now not controlled by db
	workflow.UpdatedAt = workflow.CreatedAt                              // right now not controlled by db
	workflow.IsActive = models.Active                                    // right now not controlled by db
	workflow.Status = models.Initial                                     // right now not controlled by db
	workflow.WorkflowInit = models.CustomTime{Time: models.TimeDefault}
	workflow.WorkflowCompleted = models.CustomTime{Time: models.TimeDefault}

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

func (s *WorkflowServiceImpl) GetWorkflow(userID, workflowID *string) (newWorkflow *models.Workflow, exist bool) {
	ctx, cancel := context.WithTimeout(context.Background(), models.MaxTimeoutContext)
	defer cancel()

	exist = s.ValidateWorkflowGlobalUUID(workflowID)
	if !exist {
		return nil, false
	}

	exist = s.retryTemplateWithBool(ctx, models.MaxAttempts, func() bool {
		newWorkflow, exist = s.retriesGetWorkflow(userID, workflowID)
		return exist
	})
	if exist {
		return newWorkflow, exist
	}

	// TODO: dead letter
	return nil, false
}

func (s *WorkflowServiceImpl) GetAllWorkflows(userID *string) (allWorkflows []models.Workflow, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), models.MaxTimeoutContext)
	defer cancel()

	err = s.retryTemplateWithError(ctx, models.MaxAttempts, func() error {
		var lastError error
		allWorkflows, lastError = s.retriesGetAllWorkflows(ctx, userID)
		return lastError
	})

	// TODO: Dead Letter Queue
	return allWorkflows, err
}

func (s *WorkflowServiceImpl) retryTemplateWithError(ctx context.Context, maxAttempts int, operation func() error) error {
	for i := 1; i < maxAttempts; i++ {
		if err := operation(); err == nil {
			return nil
		}
		// not necesaary last wait time
		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
		}
	}
	return fmt.Errorf("ERROR | operation failed after %d attempts", maxAttempts)
}

func (s *WorkflowServiceImpl) retryTemplateWithBool(ctx context.Context, maxAttempts int, operation func() bool) bool {
	for i := 1; i < maxAttempts; i++ {
		if exist := operation(); exist {
			return exist
		}
		// not necesaary last wait time
		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		select {
		case <-ctx.Done():
			return false
		case <-time.After(waitTime):
		}
	}
	log.Printf("ERROR | operation failed after %d attempts", maxAttempts)
	return false
}

func (s *WorkflowServiceImpl) UpdateWorkflow(workflow *models.Workflow) (updated bool, exist bool) {
	ctx, cancel := context.WithTimeout(context.Background(), models.MaxTimeoutContext)
	defer cancel()

	exist = s.ValidateWorkflowGlobalUUID(&workflow.UUID)
	if !exist {
		return false, exist
	}

	exist = s.retryTemplateWithBool(ctx, models.MaxAttempts, func() bool {
		updated, exist = s.retriesUpdateWorkflow(workflow)
		return exist
	})
	return updated, exist
}

func (s *WorkflowServiceImpl) fromWorkflowFrontendToBackend(fw *models.WorkflowFrontend) *models.Workflow {
	return &models.Workflow{
		Name:            fw.WorkflowName,
		Description:     fw.Description,
		DirectoryToSave: fw.DirectoryToSave,
		UserID:          fw.UserID,
		Nodes: []models.Node{
			s.createInitialNode(),
		},
		Viewport: &models.Viewport{
			X:    234.5,
			Y:    534.5,
			Zoom: 1,
		},
	}
}

func (s *WorkflowServiceImpl) ValidateWorkflowGlobalUUID(uuid *string) bool {
	return s.redisRepo.ValidateWorkflowGlobalUUID(uuid)
}

func (s *WorkflowServiceImpl) ValidateUserWorkflowUUID(worklfowID, name *string) bool {
	return s.redisRepo.ValidateUserWorkflowUUID(worklfowID, name)
}

func (s *WorkflowServiceImpl) retriesUpdateWorkflow(workflow *models.Workflow) (updated, exist bool) {
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

	workflow.UpdatedAt = time.Now().UTC().Format(models.LayoutTimestamp) // right now not controlled by db

	updated = s.brokerRepo.Update(workflow)
	if !updated {
		log.Printf("ERROR | Failed to publish workflow event")
		s.redisRepo.Remove(workflow)
		return false, false
	}
	return updated, exist
}

func (s *WorkflowServiceImpl) retriesGetWorkflow(userID, workflowID *string) (newWorkflow *models.Workflow, exist bool) {
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

func (s *WorkflowServiceImpl) retriesGetAllWorkflows(_ context.Context, userID *string) ([]models.Workflow, error) {
	// by default return array all workflows limited to 999
	allWorkflows, err := s.httpRepo.GetAllWorkflows(userID, models.MaxRowsFromDB) // limited to 999 row
	if err != nil {
		return []models.Workflow{}, err
	}
	return allWorkflows.Data, err
}

func (s *WorkflowServiceImpl) createInitialNode() models.Node {
	return models.Node{
		ID:   "initial-node",
		Type: "wrapperNode",
		Position: &models.Position{
			X: 2,
			Y: 0,
		},
		Data: &models.DataNode{
			ID:             "initial-node",
			Label:          "Start Point",
			Options:        "Initial Options",
			Description:    "This is the starting point of your workflow",
			WorkflowID:     "", // asigned later
			NodeID:         "initial-node",
			CredentialData: models.RequestCreateCredential{},
		},
		Measured: &models.Measured{
			Width:  50,
			Height: 50,
		},
	}
}
