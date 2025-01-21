package repos

import (
	"minireipaz/pkg/domain/models"
	"time"
)

type ActionsService interface {
	CreateActionsGoogleSheet(newAction models.RequestGoogleAction, actionUserToken *string) (sendedBroker bool, sendedToService bool, actionID *string)
	CreateActionsNotion(newAction models.RequestGoogleAction, actionUserToken *string) (sendedBroker bool, sendedToService bool, actionID *string)
}

type ActionsHTTPRepository interface {
	SendAction(newAction *models.RequestGoogleAction, actionUserToken *string) (sended bool)
	PublishCommand(data *models.ActionsCommand, serviceUser *string) *models.ResponseGetGoogleSheetByID
}

type ActionsRedisRepoInterface interface {
	Create(newAction *models.RequestGoogleAction) (created bool, exist bool, err error)
	Remove(newAction *models.RequestGoogleAction) (removed bool)
	ValidateActionGlobalUUID(field *string) (bool, error)
	AcquireLock(key, value string, expiration time.Duration) (locked bool, err error)
	RemoveLock(key string) bool
	SetNX(hashKey, actionID string, expiration time.Duration) (bool, error)
}

type ActionsBrokerRepository interface {
	Create(newAction *models.RequestGoogleAction) bool
}
