package httpclient

import (
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/models"
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

func (a *ActionsHTTPRepository) GetGoogleSheetByID(_ models.RequestGoogleAction) string {
	return ""
}
