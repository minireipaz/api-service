package models

import "time"

type RequestGoogleAction struct {
	ActionID       string  `json:"actionid"`
	RequestID      string  `json:"requestid"`
	Pollmode       string  `json:"pollmode"`
	Selectdocument string  `json:"selectdocument"`
	Document       string  `json:"document"`
	NameDocument   string  `json:"namedocument"`
	ResourceID     string  `json:"resourceid"` // document id for example
	Operation      string  `json:"operation"`
	Data           string  `json:"data"`
	CredentialID   string  `json:"credentialid"`
	Sub            string  `json:"sub"`
	Type           string  `json:"type"`
	WorkflowID     string  `json:"workflowid"`
	NodeID         string  `json:"nodeid"`
	RedirectURL    string  `json:"redirecturl"`
	Status         string  `json:"status"`        // Default: 'pending'
	ErrorMessage   *string `json:"error_message"` // Nullable
	CreatedAt      string  `json:"createdat"`
}

type ActionData struct {
	ActionID string `json:"actioid"`
}

type ResponseGetGoogleSheetByID struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
	Data   string `json:"data"`
}

type ActionsCommand struct {
	Actions   *RequestGoogleAction `json:"actions"`
	Type      *string              `json:"type,omitempty"`
	Timestamp *time.Time           `json:"timestamp,omitempty"`
}
