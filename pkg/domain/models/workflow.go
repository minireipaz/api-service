package models

import (
	"time"

	"github.com/google/uuid"
)

type TypeErrors int64

const (
	UUIDExist                     = "UUID already exists"
	WorkflowNameRequired          = "Workflow name is required"
	WorkflowNameInvalid           = "Workflow name must be alphanumeric with max length of 255"
	WorkflowNameExist             = "workflow name already exists for this user"
	WorkflowNameNotExist          = "workflow name not exists for this user"
	UUIDCannotGenerate            = "error checking UUID existence"
	WorkflowNameCannotGenerate    = "error checking workflow name existence"
	WorkflowNameNotGenerate       = "cannot create new workflow"
	WorkflowDirectorySaveRequired = "Directory to save is required"
	WorkflowDirectoryInvalid      = "Directory to save must be alphanumeric with max length of 255"
	UUIDInvalid                   = "UUID must be a valid UUID"
	WorkflowDateInvalid           = "Invalid date"
)

type Workflow struct {
	Sub             string    `json:"sub,omitempty"`
	UUID            uuid.UUID `json:"uuid,omitempty"`
	WorkflowName    string    `json:"workflowname" binding:"required,alphanum,max=255"`
	DirectoryToSave string    `json:"directorytosave" binding:"required,alphanum,max=255"`
	CreatedAt       string    `json:"createdat,omitempty"`
	UpdatedAt       string    `json:"updatedat,omitempty"`
}

type WorkflowDetail struct {
	WorkflowID          string     `json:"workflow_id"`
	WorkflowName        string     `json:"workflow_name"`
	WorkflowDescription *string    `json:"workflow_description,omitempty"`
	WorkflowStatus      *int       `json:"workflow_status,omitempty"`
	ExecutionStatus     *int       `json:"execution_status,omitempty"`
	StartTime           *time.Time `json:"start_time,omitempty"`
	Duration            *int       `json:"duration,omitempty"`
}

type WorkflowCounts struct {
	TotalWorkflows      int `json:"total_workflows"`
	SuccessfulWorkflows int `json:"successful_workflows"`
	FailedWorkflows     int `json:"failed_workflows"`
	PendingWorkflows    int `json:"pending_workflows"`
}

type RecentWorkflow struct {
	WorkflowName        string     `json:"workflow_name"`
	WorkflowDescription *string    `json:"workflow_description,omitempty"`
	ExecutionStatus     *int       `json:"execution_status,omitempty"`
	StartTime           *time.Time `json:"start_time,omitempty"`
	Duration            *int       `json:"duration,omitempty"`
}
