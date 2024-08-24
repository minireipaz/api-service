package services

import (
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/infra/httpclient"
)

type DashboardService struct {
	dashboardRepo *httpclient.DashboardRepository
}

func NewDashboardService(repo *httpclient.DashboardRepository) *DashboardService {
	return &DashboardService{dashboardRepo: repo}
}

func (d *DashboardService) QueryDashboardInfo(userID string) (*models.InfoDashboard, error) {
	infoDashboard, err := d.dashboardRepo.GetWorkflowData(userID)
	return infoDashboard, err
}
