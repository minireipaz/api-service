package models

// ResponseInfoDashboard represent response from getting dashboard info
type ResponseInfoDashboard struct {
	Data   InfoDashboard
	Status int
	Error  string
}

// InfoDashboard contains workflow counts
// also a list of the most recent workflows
type InfoDashboard struct {
	WorkflowCounts  WorkflowCounts   `json:"workflow_counts"`
	RecentWorkflows []RecentWorkflow `json:"recent_workflows"`
}
