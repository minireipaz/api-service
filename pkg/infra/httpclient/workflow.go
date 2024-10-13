package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/models"
	"net/http"
	"net/url"
)

type WorkflowHTTPRepository struct {
	databaseHTTPURL string
	token           string
	client          HTTPClient
}

func NewWorkflowClientHTTP(client HTTPClient, clickhouseConfig config.ClickhouseConfig) *WorkflowHTTPRepository {
	return &WorkflowHTTPRepository{
		client:          client,
		databaseHTTPURL: clickhouseConfig.GetClickhouseURI(),
		token:           clickhouseConfig.GetClickhouseToken(),
	}
}

func (w *WorkflowHTTPRepository) GetWorkflowDataByID(userID, workflowID *string, limitCount uint64) (*models.InfoWorkflow, error) {
	u, err := url.Parse(w.databaseHTTPURL + "/workflow_data.json")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("token", w.token)
	q.Set("workflow_id", *workflowID)
	q.Set("user_id", *userID)
	q.Set("limit_count", fmt.Sprintf("%d", limitCount))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ERROR | response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result *models.InfoWorkflow
	// if err := json.Unmarshal(bodyBytes, &result); err != nil {
	// 	log.Printf("ERROR | cannot decode body: %s %v", string(bodyBytes), err)
	// 	return models.InfoDashboard{}, fmt.Errorf("ERROR | cannot decode token: %v", err)
	// }
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("ERROR | cannot decode body: %s %v", string(bodyBytes), err)
		return nil, fmt.Errorf("ERROR | cannot decode token: %v", err)
	}

	return result, nil
}
