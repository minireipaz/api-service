package repos

import (
	"minireipaz/pkg/domain/models"
	"time"
)

type ActionsService interface {
	GetGoogleSheetByID(newAction models.RequestGoogleAction, actionUserToken *string) (created bool, exist bool, actionID *string)
}

type ActionsHTTPRepository interface {
	SendAction(newAction *models.RequestGoogleAction, actionUserToken *string) (sended bool)
	PublishCommand(data *models.ActionsCommand, serviceUser *string) *string
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
