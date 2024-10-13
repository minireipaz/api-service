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
	CommandTypeCreate = "create"
	CommandTypeUpdate = "update"
	CommandTypeDelete = "delete"
)

type WorkflowCommand struct {
	Type      string                  `json:"type"`
	Workflow  *models.WorkflowPayload `json:"workflow"`
	Timestamp time.Time               `json:"timestamp"`
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
		Type:      CommandTypeCreate,
		Workflow:  payload,
		Timestamp: time.Now(),
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
		Type:      CommandTypeUpdate,
		Workflow:  payload,
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

func (w *WorkflowKafkaRepository) workflowToPayload(workflow *models.Workflow, commandType string) (*models.WorkflowPayload, error) {
	var nodesJSON, edgesJSON, viewportJSON *string

	if workflow.Nodes != nil {
		nodeBytes, err := json.Marshal(workflow.Nodes)
		if err != nil {
			return &models.WorkflowPayload{}, err
		}
		nodeStr := string(nodeBytes)
		nodesJSON = &nodeStr
	}
	if workflow.Edges != nil {
		edgeBytes, err := json.Marshal(workflow.Edges)
		if err != nil {
			return &models.WorkflowPayload{}, err
		}
		edgeStr := string(edgeBytes)
		edgesJSON = &edgeStr
	}
	if workflow.Viewport != nil {
		viewportBytes, err := json.Marshal(workflow.Viewport)
		if err != nil {
			return &models.WorkflowPayload{}, err
		}
		viewportStr := string(viewportBytes)
		viewportJSON = &viewportStr
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
		Type:              commandType,
	}, nil
}
