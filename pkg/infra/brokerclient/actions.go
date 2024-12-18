package brokerclient

import (
	"encoding/json"
	"log"
	"minireipaz/pkg/domain/models"
)

type ActionsKafkaRepository struct {
	client KafkaClient
}

const (
	// CommandTypeCreate = "create"
	// CommandTypeUpdate = "update"
	// CommandTypeDelete = "delete"
	TopicName = "actions.command"
)

func NewActionsKafkaRepository(client KafkaClient) *ActionsKafkaRepository {
	return &ActionsKafkaRepository{
		client: client,
	}
}

func (a *ActionsKafkaRepository) Create(newAction *models.RequestGoogleAction) (sended bool) {
	command := models.ActionsCommand{
		Actions: newAction,
		// Type:      CommandTypeCreate,
		// Timestamp: time.Now(),
	}
	sended = a.PublishCommand(command, newAction.ActionID)
	return sended
}

func (a *ActionsKafkaRepository) PublishCommand(payload models.ActionsCommand, key string) bool {
	command, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ERROR | Cannot transform to JSON %v", err)
		return false
	}

	err = a.client.Produce(TopicName, []byte(key), command)
	return err == nil
}
