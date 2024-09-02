package repos

import "minireipaz/pkg/domain/models"

type DashboardHTTPRepository interface {
	GetWorkflowData(userID string) (models.InfoDashboard, error)
}
