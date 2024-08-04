package models

import (
	"github.com/google/uuid"
)

type TypeErrors int64

const (
	UUIDExist                  = "UUID already exists"
	WorkflowNameExist          = "workflow name already exists for this user"
	UUIDCannotGenerate         = "error checking UUID existence"
	WorkflowNameCannotGenerate = "error checking workflow name existence"
	WorkflowNameNotGenerate    = "cannot create new workflow"
)

type Workflow struct {
	Sub             string    `json:"sub,omitempty"`
	UUID            uuid.UUID `json:"uuid,omitempty"`
	WorkflowName    string    `json:"workflowname"`
	DirectoryToSave string    `json:"directorytosave,omitempty"`
	CreatedAt       string    `json:"createdat,omitempty"`
	UpdatedAt       string    `json:"updatedat,omitempty"`
}
