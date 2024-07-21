package models

import (
	"time"

	"github.com/google/uuid"
)

type TypeErrors int64

const (
	LayoutTimestamp  = "2006-01-02T15:04:05Z07:00"
	MaxAttempts      = 10
	MinSleepDuration = 100 * time.Millisecond // min time wait
	MaxSleepDuration = 500 * time.Millisecond // max time wait
	SleepOffset      = 50 * time.Millisecond  // offset

)

const (
	UUIDExist                  = "UUID already exists"
	WorkflowNameExist          = "workflow name already exists for this user"
	UUIDCannotGenerate         = "error checking UUID existence"
	WorkflowNameCannotGenerate = "error checking workflow name existence"
  WorkflowNameNotGenerate = "cannot create new workflow"
)

type Workflow struct {
	Sub             string    `json:"sub"`
	UUID            uuid.UUID `json:"uuid"`
	WorkflowName    string    `json:"workflowname"`
	DirectoryToSave string    `json:"directorytosave"`
	CreatedAt       string    `json:"createdat"`
	UpdatedAt       string    `json:"updatedat"`
}
