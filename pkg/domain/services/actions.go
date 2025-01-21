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

func (a *ActionsServiceImpl) CreateActionsGoogleSheet(newAction models.RequestGoogleAction, actionUserToken *string) (sendedBroker bool, sendedToService bool, actionID *string) {
	for i := 1; i < models.MaxAttempts; i++ {
		now := time.Now().UTC().Format(models.LayoutTimestamp)
		// looped 10 times with time.sleep in case uuid collisions
		// locked with same uuid for 30 seconds
		locked, err := a.setActionID(&newAction, &now)
		if err != nil || !locked { //  in case that 10 loops cannot get new UUID just return because cannot get new uuid
			return false, locked, nil
		}
		// dont create lock, just check if exist lock and in case not exist lock return false
		sendedBroker, sendedToService = a.retriesCreateAction(&newAction, now, actionUserToken, sendedBroker, sendedToService)
		if sendedBroker && sendedToService { // happy path
			// remove lock in case not passed 30 seconds
			a.removeLockActionID(&newAction.ActionID)
			return sendedBroker, sendedToService, &newAction.ActionID
		}
		// if !sendedBroker && !sendedToService {
		// 	return true, true, nil
		// }

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("WARNING | Failed to create action %s for user %s , attempt %d:. Retrying in %v", newAction.ActionID, newAction.Sub, i, waitTime)
		time.Sleep(waitTime)
	}
	log.Print("ERROR | Needs to add to Dead Letter. Cannot send action to broker or service")
	// TODO: dead letter
	return false, false, nil
}

func (a *ActionsServiceImpl) ValidateActionGlobalUUID(field *string) (bool, error) {
	return a.redisRepo.ValidateActionGlobalUUID(field)
}

// sendedBroker passed to Function CDC problema
func (a *ActionsServiceImpl) retriesCreateAction(newAction *models.RequestGoogleAction, now string, actionUserToken *string, sendedBroker bool, sendedToService bool) (bool, bool) {
	newAction.CreatedAt = now
	created, exist, err := a.redisRepo.Create(newAction)
	if err != nil {
		log.Printf("ERROR | acquiring lock for retriesCreateAction: %v", err)
		return false, false
	}
	if exist || !created {
		return false, false
	}

	if !sendedBroker {
		// case you dont have a connection of type http (http sink, http connect, ...) activated/created in the broker,
		//  enable sending to the service that processes the action. sending is done via http
		sendedBroker = a.brokerRepo.Create(newAction)
		if !sendedBroker {
			log.Printf("ERROR | SendBroker Failed to publish action event %v", newAction)
			a.redisRepo.Remove(newAction)
			return false, false
		}
	}

	// this section is necessary if not set http sink from kafka connect, connector, etc...
	// posible HTTP_SINK_ENABLED values: y/n
	if config.GetEnv("CONNECTOR_HTTP_SINK_ENABLED", "n") == "n" {
		if !sendedToService {
			// TODO: mode testmode for testing pourpose Pollmode maybe can be omitted
			if newAction.Testmode || newAction.Pollmode == models.NopollNode {
				sendedToService = a.httpRepo.SendAction(newAction, actionUserToken)
				if !sendedToService {
					log.Printf("ERROR | SendedtoService Failed to publish action event %v", newAction)
					return false, false
				}
			} else {
				// TODO: scheduler
				log.Printf("TODO: implement scheduler")
			}
		}
	}
	return sendedBroker, sendedToService
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

// maybe can be merged
func (a *ActionsServiceImpl) CreateActionsNotion(newAction models.RequestGoogleAction, actionUserToken *string) (sendedBroker bool, sendedToService bool, actionID *string) {
	for i := 1; i < models.MaxAttempts; i++ {
		now := time.Now().UTC().Format(models.LayoutTimestamp)
		// looped 10 times with time.sleep in case uuid collisions
		// locked with same uuid for 30 seconds
		locked, err := a.setActionID(&newAction, &now)
		if err != nil || !locked { //  in case that 10 loops cannot get new UUID just return because cannot get new uuid
			return false, locked, nil
		}
		// dont create lock, just check if exist lock and in case not exist lock return false
		sendedBroker, sendedToService = a.retriesCreateAction(&newAction, now, actionUserToken, sendedBroker, sendedToService)
		if sendedBroker && sendedToService { // happy path
			// remove lock in case not passed 30 seconds
			a.removeLockActionID(&newAction.ActionID)
			return sendedBroker, sendedToService, &newAction.ActionID
		}
		// if !sendedBroker && !sendedToService {
		// 	return true, true, nil
		// }

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("WARNING | Failed to create action %s for user %s , attempt %d:. Retrying in %v", newAction.ActionID, newAction.Sub, i, waitTime)
		time.Sleep(waitTime)
	}
	log.Print("ERROR | Needs to add to Dead Letter. Cannot send action to broker or service")
	// TODO: dead letter
	return false, false, nil
}
