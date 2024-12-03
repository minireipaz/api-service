package brokerclient

import (
	"encoding/json"
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/domain/models"
	"time"
)

type ActionsKafkaRepository struct {
	client KafkaClient
}

// const (
// 	CommandTypeCreate = "create"
// 	CommandTypeUpdate = "update"
// 	CommandTypeDelete = "delete"
// )

type ActionsCommand struct {
	Actions   *models.RequestGoogleAction `json:"Actions,omitempty"`
	Type      *string                     `json:"type,omitempty"`
	Timestamp *time.Time                  `json:"timestamp,omitempty"`
}

func NewActionsKafkaRepository(client KafkaClient) *ActionsKafkaRepository {
	return &ActionsKafkaRepository{
		client: client,
	}
}

func (a *ActionsKafkaRepository) Create(newAction *models.RequestGoogleAction) (sended bool) {
	// payload, err := a.workflowToPayload(workflow, CommandTypeCreate)
	// if err != nil {
	// 	log.Printf("ERROR | Cannot convert workflow to payload: %v", err)
	// 	return false
	// }
	command := ActionsCommand{
		Actions: newAction,
		// Type:      CommandTypeCreate,
		// Timestamp: time.Now(),
	}
	sended = a.PublishCommand(command, newAction.ActionID)
	return sended
}

func (a *ActionsKafkaRepository) PublishCommand(payload ActionsCommand, key string) bool {
	command, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ERROR | Cannot transform to JSON %v", err)
		return false
	}

	for i := 1; i < models.MaxAttempts; i++ {
		err = a.client.Produce("workflows.command", []byte(key), command)
		if err == nil {
			return true
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to Broker, attempt %d: %v. Retrying in %v", i, err, waitTime)
		time.After(waitTime)
	}

	return false
}
