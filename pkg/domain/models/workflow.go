package models

import (
	"time"
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

type WorkflowFrontend struct {
	Sub             string `json:"sub,omitempty"`
	UUID            string `json:"id,omitempty"`
	WorkflowName    string `json:"name" binding:"required,alphanum,max=255"`
	Description     string `json:"description,omitempty"`
	DirectoryToSave string `json:"directory_to_save" binding:"required,alphanum,max=255"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
}

type IsActive uint8

const (
	Active IsActive = iota + 1 // Active = 1
	Draft                      // Draft = 2
	Paused                     // Paused = 3
)

type Status uint8

const (
	Initial    Status = iota + 1 // Initial = 1
	Pending                      // Pending = 2
	Completed                    // Completed = 3
	Processing                   // Processing = 4
	Failed                       // Failed = 5
)

type Workflow struct {
	UUID              string    `json:"id,omitempty"`
	UserID            string    `json:"user_id,omitempty"`
	Name              string    `json:"name,omitempty"`
	Description       string    `json:"description,omitempty"`
	IsActive          IsActive  `json:"is_active,omitempty"` // Enum8('active' = 1, 'draft' = 2, 'paused' = 3) DEFAULT 'active'
	CreatedAt         string    `json:"created_at,omitempty"`
	UpdatedAt         string    `json:"updated_at,omitempty"`
	WorkflowInit      time.Time `json:"workflow_init,omitempty"`
	WorkflowCompleted time.Time `json:"workflow_completed,omitempty"`
	Status            Status    `json:"status,omitempty"` // Enum8('initial' = 1, 'pending' = 2, 'completed' = 3, 'processing' = 4, 'failed' = 5) DEFAULT 'initial'
	DirectoryToSave   string    `json:"directory_to_save"`
}

type WorkflowDetail struct {
	WorkflowID          string      `json:"workflow_id"`
	WorkflowName        string      `json:"workflow_name"`
	WorkflowDescription *string     `json:"workflow_description,omitempty"`
	WorkflowStatus      *int        `json:"workflow_status,omitempty"`
	ExecutionStatus     *int        `json:"execution_status,omitempty"`
	StartTime           *CustomTime `json:"start_time,omitempty"`
	Duration            *int        `json:"duration,omitempty"`
}

type WorkflowsCount struct {
	TotalWorkflows      *int64 `json:"total_workflows,omitempty"`
	SuccessfulWorkflows *int64 `json:"successful_workflows,omitempty"`
	FailedWorkflows     *int64 `json:"failed_workflows,omitempty"`
	PendingWorkflows    *int64 `json:"pending_workflows,omitempty"`
	// RecentWorkflow      []RecentWorkflow
}

type RecentWorkflow struct {
	WorkflowName        string       `json:"workflow_name"`
	WorkflowDescription *string      `json:"workflow_description,omitempty"`
	ExecutionStatus     *int         `json:"execution_status,omitempty"`
	StartTime           *CustomTime  `json:"start_time,omitempty"`
	Duration            *interface{} `json:"duration,omitempty"`
}
