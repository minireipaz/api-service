package redisclient

import (
	"encoding/json"
	"fmt"
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/domain/models"
	"time"

	"github.com/google/uuid"
)

type WorkflowService interface {
}

type WorkflowRepository struct {
	redisClient *RedisClient
}

func NewWorkflowRepository(redisClient *RedisClient) *WorkflowRepository {
	return &WorkflowRepository{redisClient: redisClient}
}

func (r *WorkflowRepository) Create(workflow *models.Workflow) (created bool, exist bool) {
	err := r.redisClient.SetWorkflow(workflow)
	if err != nil {
		return false, true
	}

	return true, false
}

func (r *WorkflowRepository) Update(workflow *models.Workflow) (updated bool, exist bool) {
	err := r.redisClient.UpdateWorkflow(workflow)
	if err != nil {
		return false, true
	}

	return true, false
}

func (r *WorkflowRepository) Remove(workflow *models.Workflow) bool {
	err := r.redisClient.RemoveWorkflow(workflow)
	return err != nil
}

func (r *WorkflowRepository) ValidateWorkflowGlobalUUID(uuid *string) bool {
	exist := r.redisClient.Hexists("workflows:all", *uuid)
	return exist
}

func (r *WorkflowRepository) ValidateUserWorkflowUUID(userID, name *string) bool {
	exist := r.redisClient.Hexists(fmt.Sprintf("users:%s", *userID), *name)
	return exist
}

func (r *WorkflowRepository) GetByUUID(id uuid.UUID) (*models.Workflow, error) {
	workflowJSON, err := r.redisClient.Get(fmt.Sprintf("workflow:%d", id))
	if err != nil {
		return nil, err
	}

	if workflowJSON == "" { // not exist key
		// TODO: better
		return nil, nil
	}

	var workflow models.Workflow
	err = json.Unmarshal([]byte(workflowJSON), &workflow)
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (r *WorkflowRepository) AcquireLock(key, value string, expiration time.Duration) (locked bool, err error) {
	for i := 1; i < models.MaxAttempts; i++ {
		locked, err = r.redisClient.AcquireLock(key, value, expiration)
		if err == nil {
			return locked, err
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to redis for key %s, attempt %d: %v. Retrying in %v", key, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return false, fmt.Errorf("ERROR | Cannot create lock for key %s. More than 10 intents", key)
}

func (r *WorkflowRepository) RemoveLock(key string) bool {
	for i := 1; i < models.MaxAttempts; i++ {
		countRemoved, err := r.redisClient.RemoveLock(key)
		if countRemoved == 0 {
			log.Printf("WARNING | Key already removed, previuous process take more than 20 seconds")
		}
		if err == nil && countRemoved <= 1 {
			return true
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to redis for key %s, attempt %d: %v. Retrying in %v", key, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return false
}
