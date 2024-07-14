package models

import "time"

type Workflow struct {
	ID              int       `json:"id"`
	WorkflowName    string    `json:"workflowname"`
	DirectoryToSave string    `json:"directorytosave"`
	CreatedAt       time.Time `json:"createdat"`
	UpdatedAt       time.Time `json:"updatedat"`
}
