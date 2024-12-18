package httpclient

import (
	"encoding/json"
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/models"
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
	now := time.Now().UTC()
	typeCommand := models.CommandTypeCreate
	command := models.ActionsCommand{
		Actions:   newAction,
		Type:      &typeCommand,
		Timestamp: &now,
	}

	response := a.PublishCommand(&command, actionUserToken)
	if response == nil {
		return false
	}
	return *response != ""
}

func (a *ActionsHTTPRepository) PublishCommand(data *models.ActionsCommand, serviceUser *string) *string {
	url, err := getActionsURL("/api/actions/google/sheets")
	if err != nil {
		return nil
	}

	body, err := a.client.DoRequest("POST", url, *serviceUser, data)
	if err != nil {
		return nil
	}

	var response string
	if err := json.Unmarshal(body, &response); err != nil {
		return nil
	}

	return &response
}
