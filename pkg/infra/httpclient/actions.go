package httpclient

import (
	"encoding/json"
	"log"
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/models"
	"net/http"
	"time"
)

type ActionsHTTPRepository struct {
	databaseHTTPURL string
	token           string
	client          HTTPClient
}

func NewActionsClientHTTP(client HTTPClient, clickhouseConfig config.ClickhouseConfig) *ActionsHTTPRepository {
	return &ActionsHTTPRepository{
		client:          client,
		databaseHTTPURL: clickhouseConfig.GetClickhouseURI(),
		token:           clickhouseConfig.GetClickhouseToken(),
	}
}

func (a *ActionsHTTPRepository) SendAction(newAction *models.RequestGoogleAction, actionUserToken *string) (sended bool) {
	command := models.ActionsCommand{
		Actions:   newAction,
		Type:      models.CommandTypeCreate,
		Timestamp: time.Now().UTC(),
	}

	response := a.PublishCommand(&command, actionUserToken)
	if response == nil {
		return false
	}
	// from action return accepdted 202
	return response.Status == http.StatusAccepted
}

func (a *ActionsHTTPRepository) PublishCommand(data *models.ActionsCommand, serviceUser *string) *models.ResponseGetGoogleSheetByID {
	url, err := getActionsURL("/api/actions/google/sheets")
	if err != nil {
		return nil
	}
	// Request-Response btw response not get data
	// Btw testmode
	body, err := a.client.DoRequest("POST", url, *serviceUser, data)
	if err != nil {
		log.Printf("ERROR | Cannot send to action service %v", err)
		return nil
	}

	var response models.ResponseGetGoogleSheetByID
	if err := json.Unmarshal(body, &response); err != nil {
		return nil
	}
	// response can contain data, btw right now client polling data
	// bool option devtest can be used to response directly to request
	return &response
}

// func (a *ActionsHTTPRepository) GetActionByID(actionID *string, userID *string, limitCount uint64) (data *string, err error) {
// 	u, err := url.Parse(a.databaseHTTPURL + "/action_workflow_data.json")
// 	if err != nil {
// 		return nil, err
// 	}

// 	q := u.Query()
// 	q.Set("token", a.token)
// 	q.Set("action_id", *actionID)
// 	q.Set("user_id", *userID)
// 	q.Set("limit_count", fmt.Sprintf("%d", limitCount))
// 	u.RawQuery = q.Encode()

// 	req, err := http.NewRequest("GET", u.String(), nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp, err := a.client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		bodyBytes, _ := io.ReadAll(resp.Body)
// 		return nil, fmt.Errorf("ERROR | response: %d, body: %s", resp.StatusCode, string(bodyBytes))
// 	}

// 	var result *models.InfoWorkflow
// 	// if err := json.Unmarshal(bodyBytes, &result); err != nil {
// 	// 	log.Printf("ERROR | cannot decode body: %s %v", string(bodyBytes), err)
// 	// 	return models.InfoDashboard{}, fmt.Errorf("ERROR | cannot decode token: %v", err)
// 	// }
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		bodyBytes, _ := io.ReadAll(resp.Body)
// 		log.Printf("ERROR | cannot decode body: %s %v", string(bodyBytes), err)
// 		return nil, fmt.Errorf("ERROR | cannot decode token: %v", err)
// 	}

// 	return result, nil
// }
