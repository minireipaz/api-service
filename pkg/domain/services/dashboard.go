package services

import (
	// "minireipaz/pkg/domain/models"
	"minireipaz/pkg/infra/httpclient"
)

type DashboardService struct {
	dashboardRepo *httpclient.DashboardRepository
}

func NewDashboardService(repo *httpclient.DashboardRepository) *DashboardService {
	return &DashboardService{dashboardRepo: repo}
}
