package models

import (
	"github.com/google/uuid"
)

type TypeErrors int64

const (
	UUIDExist                     = "UUID already exists"
	WorkflowNameRequired          = "Workflow name is required"
	WorkflowNameInvalid           = "Workflow name must be alphanumeric with max length of 255"
	WorkflowNameExist             = "workflow name already exists for this user"
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
	WorkflowName    string    `json:"workflow_name" binding:"required,alphanum,max=255"`
	DirectoryToSave string    `json:"directory_to_save" binding:"required,alphanum,max=255"`
	CreatedAt       string    `json:"createdat,omitempty"`
	UpdatedAt       string    `json:"updatedat,omitempty"`
}
