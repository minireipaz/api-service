package httpclient

import (
// "encoding/json"
// "fmt"
// "log"
// "minireipaz/pkg/domain/models"
)

type DashboardRepository struct {
	client HTTPClient
}

func NewDashboardRepository(client HTTPClient) *DashboardRepository {
	return &DashboardRepository{client: client}
}
