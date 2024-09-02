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
	infoDashboard, err := d.dashboardHTTPRepo.GetWorkflowData(userID)
	return infoDashboard, err
}
