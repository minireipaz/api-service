package repos

import "minireipaz/pkg/domain/models"

type DashboardHTTPRepository interface {
	GetLastWorkflowData(userID string, limitCount uint64) (models.InfoDashboard, error)
}
