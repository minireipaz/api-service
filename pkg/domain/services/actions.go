package services

import (
	"fmt"
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"time"

	"github.com/google/uuid"
)

type ActionsServiceImpl struct {
	redisRepo  repos.ActionsRedisRepoInterface
	brokerRepo repos.ActionsBrokerRepository
	httpRepo   repos.ActionsHTTPRepository
}

func NewActionsService(repoRedis repos.ActionsRedisRepoInterface, repoBroker repos.ActionsBrokerRepository, repoHTTP repos.ActionsHTTPRepository) repos.ActionsService {
	return &ActionsServiceImpl{
		redisRepo:  repoRedis,
		brokerRepo: repoBroker,
		httpRepo:   repoHTTP,
	}
}

func (a *ActionsServiceImpl) GetGoogleSheetByID(newAction models.RequestGoogleAction) (created bool, exist bool, action *models.ActionData) {
	now := time.Now().UTC().Format(models.LayoutTimestamp)
	newAction.ActionID, newAction.RequestID = a.generateActionID(&now)
	exist, err := a.ValidateActionGlobalUUID(&newAction.ActionID) // check individual
	if err != nil {
		created = false
		exist = false
		return created, exist, nil
	}
	if exist {
		created = false
		return created, exist, nil
	}

	for i := 1; i < models.MaxAttempts; i++ {
		created, exist = a.retriesCreateAction(&newAction, now)
		if !exist && created {
			return created, exist, nil
		}
		if exist {
			return false, true, nil
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("WARNING | Failed to create workflow, attempt %d:. Retrying in %v", i, waitTime)
		time.Sleep(waitTime)
	}
	log.Print("ERROR | Needs to add to Dead Letter. Cannot create workflow")
	// TODO: dead letter
	return false, false, nil
}

func (a *ActionsServiceImpl) ValidateActionGlobalUUID(field *string) (bool, error) {
	return a.redisRepo.ValidateActionGlobalUUID(field)
}

func (a *ActionsServiceImpl) retriesCreateAction(newAction *models.RequestGoogleAction, now string) (created bool, exist bool) {
	newAction.CreatedAt = now
	created, exist, err := a.redisRepo.Create(newAction)
	if err != nil {
		log.Printf("ERROR | acquiring lock: %v", err)
		return false, false
	}
	if exist || !created {
		return created, exist
	}
	// if !createdRedis {
	// 	return false, false
	// }

	sended := a.brokerRepo.Create(newAction)
	if !sended {
		log.Printf("ERROR | Failed to publish action event %v", newAction)
		a.redisRepo.Remove(newAction)
		return false, false
	}
	return created, exist
}

// TODO: maybe can make general function to create requestID and IDs
func (a *ActionsServiceImpl) generateActionID(now *string) (string, string) {
	actionID := uuid.New().String()
	if actionID == "" { // in case fail
		actionID = uuid.New().String()
	}
	requestID := fmt.Sprintf("%s_%s", actionID, *now)
	return actionID, requestID
}
