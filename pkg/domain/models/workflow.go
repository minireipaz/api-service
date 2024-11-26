package models

import (
	"encoding/json"
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
	RateLimitUpdate               = 10 * time.Second
)

type WorkflowFrontend struct {
	UserID          string `json:"user_id,omitempty"`
	UUID            string `json:"id,omitempty"`
	WorkflowName    string `json:"name" binding:"required,alphanum,max=255"`
	Description     string `json:"description,omitempty"`
	DirectoryToSave string `json:"directory_to_save" binding:"required,alphanum,max=255"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	//
	UserToken         string     `json:"access_token,omitempty"`
	IsActive          IsActive   `json:"is_active,omitempty"` // Enum8('active' = 1, 'draft' = 2, 'paused' = 3) DEFAULT 'active'
	WorkflowInit      CustomTime `json:"workflow_init,omitempty"`
	WorkflowCompleted CustomTime `json:"workflow_completed,omitempty"`
	Status            *Status    `json:"status,omitempty"`
	Duration          *int64     `json:"duration,omitempty"`
	Nodes             []Node     `json:"nodes,omitempty"`
	Edges             []Edge     `json:"edges,omitempty"`
	Viewport          *Viewport  `json:"viewport,omitempty"`
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
	UserToken         string    `json:"access_token,omitempty"`
	Duration          *int64    `json:"duration,omitempty"`
	Nodes             []Node    `json:"nodes,omitempty"`
	Edges             []Edge    `json:"edges,omitempty"`
	Viewport          *Viewport `json:"viewport,omitempty"`
}

type WorkflowCommand struct {
	Workflow WorkflowPayload `json:"workflow"`
}

type WorkflowPayload struct {
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
	UserToken         string    `json:"access_token,omitempty"`
	Duration          *int64    `json:"duration,omitempty"`
	Nodes             *string   `json:"nodes,omitempty"`
	Edges             *string   `json:"edges,omitempty"`
	Viewport          *string   `json:"viewport,omitempty"`
	Version           *uint32   `json:"version,omitempty"`
	TypeCommand       string    `json:"typecommand,omitempty"`
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
}

type RecentWorkflow struct {
	WorkflowName        string       `json:"workflow_name"`
	WorkflowDescription *string      `json:"workflow_description,omitempty"`
	ExecutionStatus     *int         `json:"execution_status,omitempty"`
	StartTime           *CustomTime  `json:"start_time,omitempty"`
	Duration            *interface{} `json:"duration,omitempty"`
}

type Edge struct {
	ID       *string `json:"id,omitempty"`
	Source   *string `json:"source,omitempty"`
	Target   *string `json:"target,omitempty"`
	Type     *string `json:"type,omitempty"`
	Animated *bool   `json:"animated,omitempty"`
	Style    *Style  `json:"style,omitempty"`
}

type Style struct {
	Stroke *string `json:"stroke,omitempty"`
}

type Node struct {
	ID       string    `json:"id,omitempty"`
	Type     string    `json:"type,omitempty"`
	Position *Position `json:"position,omitempty"`
	Data     *DataNode `json:"data,omitempty"`
	Measured *Measured `json:"measured,omitempty"`
}

type DataNode struct {
	ID             string                  `json:"id,omitempty"`
	Label          string                  `json:"label,omitempty"`
	Options        string                  `json:"options,omitempty"`
	Description    string                  `json:"description,omitempty"`
	WorkflowID     string                  `json:"workflowid,omitempty"`
	NodeID         string                  `json:"nodeid,omitempty"`
	CredentialData RequestCreateCredential `json:"credential"`
	Type           string                  `json:"type"`
}

type Measured struct {
	Width  float64 `json:"width,omitempty"`
	Height float64 `json:"height,omitempty"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Viewport struct {
	X    float64 `json:"x,omitempty"`
	Y    float64 `json:"y,omitempty"`
	Zoom float64 `json:"zoom,omitempty"`
}

type InfoWorkflow struct {
	Meta                   []Meta      `json:"meta,omitempty"`
	Data                   []Workflow  `json:"data,omitempty"`
	Rows                   *int64      `json:"rows,omitempty"`
	RowsBeforeLimitAtLeast *int64      `json:"rows_before_limit_at_least,omitempty"`
	Statistics             *Statistics `json:"statistics,omitempty"`
}

func (w *Workflow) UnmarshalJSON(data []byte) error {
	type Alias Workflow
	aux := &struct {
		*Alias
		Nodes    json.RawMessage `json:"nodes"`
		Edges    json.RawMessage `json:"edges"`
		Viewport json.RawMessage `json:"viewport"`
	}{
		Alias: (*Alias)(w),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if err := unmarshalNodes(aux.Nodes, &w.Nodes); err != nil {
		return err
	}

	if err := unmarshalEdges(aux.Edges, &w.Edges); err != nil {
		return err
	}

	if err := unmarshalViewport(aux.Viewport, &w.Viewport); err != nil {
		return err
	}

	return nil
}

func unmarshalNodes(raw json.RawMessage, nodes *[]Node) error {
	if raw == nil {
		return nil
	}
	// Try to deserialize as []Node
	if err := json.Unmarshal(raw, nodes); err != nil {
		// If fails, try to deserialize it as a string containing a JSON array
		var nodesStr string
		if err := json.Unmarshal(raw, &nodesStr); err != nil {
			return err
		}
		// Deserialize string as a JSON []Node
		if err := json.Unmarshal([]byte(nodesStr), nodes); err != nil {
			return err
		}
	}
	return nil
}

func unmarshalEdges(raw json.RawMessage, edges *[]Edge) error {
	if raw == nil {
		return nil
	}
	// Try to deserialize as []Edge
	if err := json.Unmarshal(raw, edges); err != nil {
		// If fails, try to deserialize it as a string containing a JSON array
		var edgesStr string
		if err := json.Unmarshal(raw, &edgesStr); err != nil {
			return err
		}
		// Deserialize string as a JSON []Edge
		if err := json.Unmarshal([]byte(edgesStr), edges); err != nil {
			return err
		}
	}
	return nil
}

func unmarshalViewport(raw json.RawMessage, viewport **Viewport) error {
	if raw == nil {
		return nil
	}
	// Try to deserialize as Viewport
	if err := json.Unmarshal(raw, viewport); err != nil {
		// If fails, try to deserialize it as a string containing a JSON object
		var viewportStr string
		if err := json.Unmarshal(raw, &viewportStr); err != nil {
			return err
		}
		// Deserialize string as a JSON Viewport
		if err := json.Unmarshal([]byte(viewportStr), viewport); err != nil {
			return err
		}
	}
	return nil
}
