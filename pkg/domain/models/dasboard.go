package models

// ResponseInfoDashboard represent response from getting dashboard info
type ResponseInfoDashboard struct {
	Error  string             `json:"error"`
	Data   []DashboardDetails `json:"data"`
	Status int                `json:"status"`
}

type InfoDashboard struct {
	Rows                   *int64             `json:"rows,omitempty"`
	RowsBeforeLimitAtLeast *int64             `json:"rows_before_limit_at_least,omitempty"`
	Statistics             *Statistics        `json:"statistics,omitempty"`
	Meta                   []Meta             `json:"meta,omitempty"`
	Data                   []DashboardDetails `json:"data,omitempty"`
}

type DashboardDetails struct {
	TotalWorkflows      *int64      `json:"total_workflows,omitempty"`
	SuccessfulWorkflows *int64      `json:"successful_workflows,omitempty"`
	FailedWorkflows     *int64      `json:"failed_workflows,omitempty"`
	PendingWorkflows    *int64      `json:"pending_workflows,omitempty"`
	RecentWorkflows     [][]*string `json:"recent_workflows,omitempty"`
}

type Meta struct {
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

type Statistics struct {
	Elapsed   *float64 `json:"elapsed,omitempty"`
	RowsRead  *int64   `json:"rows_read,omitempty"`
	BytesRead *int64   `json:"bytes_read,omitempty"`
}
