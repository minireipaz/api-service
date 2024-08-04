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
	Sub             string    `json:"sub"`
	UUID            uuid.UUID `json:"uuid"`
	WorkflowName    string    `json:"workflowname"`
	DirectoryToSave string    `json:"directorytosave"`
	CreatedAt       string    `json:"createdat"`
	UpdatedAt       string    `json:"updatedat"`
}
