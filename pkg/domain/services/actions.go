package services

import (
	"fmt"
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/config"
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

func (a *ActionsServiceImpl) GetGoogleSheetByID(newAction models.RequestGoogleAction, actionUserToken *string) (created bool, exist bool, actionID *string) {
	for i := 1; i < models.MaxAttempts; i++ {
		now := time.Now().UTC().Format(models.LayoutTimestamp)
		// looped 10 times with time.sleep in case uuid collisions
		// locked with same uuid for 30 seconds
		locked, err := a.setActionID(&newAction, &now)
		if err != nil || !locked { //  in case that 10 loops cannot get new UUID just return because cannot get new uuid
			return false, locked, nil
		}
		// TODO: simplify
		created, exist = a.retriesCreateAction(&newAction, now, actionUserToken)
		if !exist && created { // happy path
			// remove lock in case not passed 30 seconds
			a.removeLockActionID(&newAction.ActionID)
			return created, exist, &newAction.ActionID
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

func (a *ActionsServiceImpl) retriesCreateAction(newAction *models.RequestGoogleAction, now string, actionUserToken *string) (created bool, exist bool) {
	newAction.CreatedAt = now
	created, exist, err := a.redisRepo.Create(newAction)
	if err != nil {
		log.Printf("ERROR | acquiring lock: %v", err)
		return false, false
	}
	if exist || !created {
		return created, exist
	}

	// case you dont have a connection of type http (http sink, http connect, ...) activated/created in the broker,
	//  enable sending to the service that processes the action. sending is done via http
	sended := a.brokerRepo.Create(newAction)
	if !sended {
		log.Printf("ERROR | Failed to publish action event %v", newAction)
		a.redisRepo.Remove(newAction)
		return false, false
	}
	// this section is necessary if not set http sink from kafka connect, connector, etc...
	// posible HTTP_SINK_ENABLED values: y/n
	if config.GetEnv("CONNECTOR_HTTP_SINK_ENABLED", "n") == "n" {
		if newAction.Pollmode == "none" {
			sended = a.httpRepo.SendAction(newAction, actionUserToken)
			if !sended {
				log.Printf("ERROR | Failed to publish action event %v", newAction)
				return false, false
			}
		} else {
			// TODO: scheduler
			log.Printf("TODO: implement scheduler")
		}
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

func (a *ActionsServiceImpl) lockActionID(actionID *string) (locked bool, err error) {
	locked, err = a.redisRepo.AcquireLock("lock:"+*actionID, "1", models.MaxTimeForLocks)
	return locked, err
}

func (a *ActionsServiceImpl) removeLockActionID(actionID *string) (removed bool) {
	return a.redisRepo.RemoveLock("lock:" + *actionID)
}

func (a *ActionsServiceImpl) setActionID(newAction *models.RequestGoogleAction, now *string) (locked bool, err error) {
	var actionID, requestID string

	for i := 0; i < models.MaxAttempts; i++ {
		actionID, requestID = a.generateActionID(now)
		locked, err = a.lockActionID(&actionID)
		if err != nil {
			locked = false
			return locked, err
		}
		if locked {
			break
		}
		// not used time.sleep
	}

	// All attempts failed to find a unique UUID
	if !locked {
		return locked, fmt.Errorf("all attempts to generate a unique UUID failed")
	}

	newAction.ActionID = actionID
	newAction.RequestID = requestID
	return locked, nil
}
