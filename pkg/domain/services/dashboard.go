package services

import (
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
)

type DashboardServiceImpl struct {
	dashboardHTTPRepo repos.DashboardHTTPRepository
}

func NewDashboardService(repo repos.DashboardHTTPRepository) repos.DashboardService {
	return &DashboardServiceImpl{dashboardHTTPRepo: repo}
}

func (d *DashboardServiceImpl) QueryDashboardInfo(userID string) (models.InfoDashboard, error) {
	infoDashboard, err := d.dashboardHTTPRepo.GetLastWorkflowData(userID, 5)
	return infoDashboard, err
}
