package brokerclient

import (
	"encoding/json"
	"fmt"
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/domain/models"
	"time"
)

type WorkflowKafkaRepository struct {
	client KafkaClient
}

const (
	CommandTypeCreate = "create"
	CommandTypeUpdate = "update"
	CommandTypeDelete = "delete"
)

type WorkflowCommand struct {
	Workflow  *models.WorkflowPayload `json:"workflow,omitempty"`
	Type      *string                 `json:"type,omitempty"`
	Timestamp *time.Time              `json:"timestamp,omitempty"`
}

func NewWorkflowKafkaRepository(client KafkaClient) *WorkflowKafkaRepository {
	return &WorkflowKafkaRepository{
		client: client,
	}
}

func (w *WorkflowKafkaRepository) Create(workflow *models.Workflow) (sended bool) {
	payload, err := w.workflowToPayload(workflow, CommandTypeCreate)
	if err != nil {
		log.Printf("ERROR | Cannot convert workflow to payload: %v", err)
		return false
	}
	command := WorkflowCommand{
		Workflow: payload,
		// Type:      CommandTypeCreate,
		// Timestamp: time.Now(),
	}
	sended = w.PublishCommand(command, workflow.UUID)
	return sended
}

func (w *WorkflowKafkaRepository) Update(workflow *models.Workflow) (sended bool) {
	payload, err := w.workflowToPayload(workflow, CommandTypeUpdate)
	if err != nil {
		log.Printf("ERROR | Cannot convert workflow to payload: %v", err)
		return false
	}
	command := WorkflowCommand{
		Workflow: payload,
		// Type:     CommandTypeUpdate,
		// Timestamp: time.Now(),
	}
	sended = w.PublishCommand(command, workflow.UUID)
	return sended
}

// use sync.pool in serverless not necessary
func (w *WorkflowKafkaRepository) workflowToPayload(workflow *models.Workflow, commandType string) (*models.WorkflowPayload, error) {
	nodesJSON, err := w.serializeToJSON(workflow.Nodes)
	if err != nil {
		return nil, err
	}

	edgesJSON, err := w.serializeToJSON(workflow.Edges)
	if err != nil {
		return nil, err
	}

	viewportJSON, err := w.serializeToJSON(workflow.Viewport)
	if err != nil {
		return nil, err
	}

	return &models.WorkflowPayload{
		UUID:              workflow.UUID,
		UserID:            workflow.UserID,
		Name:              workflow.Name,
		Description:       workflow.Description,
		IsActive:          workflow.IsActive,
		CreatedAt:         workflow.CreatedAt,
		UpdatedAt:         workflow.UpdatedAt,
		WorkflowInit:      workflow.WorkflowInit,
		WorkflowCompleted: workflow.WorkflowCompleted,
		Status:            workflow.Status,
		DirectoryToSave:   workflow.DirectoryToSave,
		Nodes:             nodesJSON,
		Edges:             edgesJSON,
		Viewport:          viewportJSON,
		TypeCommand:       commandType,
	}, nil
}

// use sync.pool in serverless not necessary
func (w *WorkflowKafkaRepository) serializeToJSON(params interface{}) (*string, error) {
	if params == nil {
		return nil, fmt.Errorf("ERROR | Cannot transform to JSON %v", params)
	}
	bytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	str := string(bytes)
	return &str, nil
}

func (w *WorkflowKafkaRepository) PublishCommand(payload WorkflowCommand, key string) bool {
	command, err := json.Marshal(payload)
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
		time.After(waitTime)
	}

	return false
}
