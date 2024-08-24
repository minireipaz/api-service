package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/models"
	"net/http"
	"net/url"
)

type DashboardRepository struct {
	client          HTTPClient
	databaseHTTPURL string
	token           string
}

func NewDashboardRepository(client HTTPClient, clickhouseConfig config.ClickhouseConfig) *DashboardRepository {
	return &DashboardRepository{
		client:          client,
		databaseHTTPURL: clickhouseConfig.GetClickhouseURI(),
		token:           clickhouseConfig.GetClickhouseToken(),
	}
}

func (d *DashboardRepository) GetWorkflowData(userID string) (*models.InfoDashboard, error) {
	u, err := url.Parse(d.databaseHTTPURL + "/user_workflow_stats.json")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("token", d.token)
	q.Set("user_id", userID)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// d.setHeaders(req, d.token)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ERROR | response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.InfoDashboard
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ERROR | cannot decode token: %v", err)
	}

	return &result, nil
}

// func (c *DashboardRepository) setHeaders(req *http.Request, token string) {
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
// }
