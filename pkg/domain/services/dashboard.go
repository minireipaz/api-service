package services

import (
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
)

type DashboardService struct {
	dashboardHTTPRepo repos.DashboardHTTPRepository
}

func NewDashboardService(repo repos.DashboardHTTPRepository) *DashboardService {
	return &DashboardService{dashboardHTTPRepo: repo}
}

func (d *DashboardService) QueryDashboardInfo(userID string) (models.InfoDashboard, error) {
	infoDashboard, err := d.dashboardHTTPRepo.GetLastWorkflowData(userID, 5)
	return infoDashboard, err
}
