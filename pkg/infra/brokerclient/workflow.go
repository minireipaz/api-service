package brokerclient

import (
	"encoding/json"
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/domain/models"
	"time"
)

type WorkflowKafkaRepository struct {
	client KafkaClient
}

const (
	CommandTypeCreate = "CREATE"
	CommandTypeUpdate = "UPDATE"
	CommandTypeDelete = "DELETE"
)

type WorkflowCommand struct {
	Type      string           `json:"type"`
	Workflow  *models.Workflow `json:"workflow"`
	Timestamp time.Time        `json:"timestamp"`
}

func NewWorkflowKafkaRepository(client KafkaClient) *WorkflowKafkaRepository {
	return &WorkflowKafkaRepository{
		client: client,
	}
}

func (w *WorkflowKafkaRepository) Create(workflow *models.Workflow) (sended bool) {
	command := WorkflowCommand{
		Type:      CommandTypeCreate,
		Workflow:  workflow,
		Timestamp: time.Now(),
	}
	sended = w.PublishCommand(command, workflow.UUID)
	return sended
}

func (w *WorkflowKafkaRepository) PublishCommand(workflowCommand WorkflowCommand, key string) bool {
	command, err := json.Marshal(workflowCommand)
	if err != nil {
		log.Printf("ERROR | Cannot transform to JSON %v", err)
		return false
	}

	for i := 1; i < models.MaxAttempts; i++ {
		err = w.client.Produce("workflows.command", []byte(key), command)
		if err == nil {
			return true
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to Broker, attempt %d: %v. Retrying in %v", i, err, waitTime)
		time.Sleep(waitTime)
	}

	return false
}
