package repos

import (
	"minireipaz/pkg/domain/models"
	"time"
)

type ActionsService interface {
	GetGoogleSheetByID(newAction models.RequestGoogleAction) (created bool, exist bool, workflow *models.ActionData)
}

type ActionsHTTPRepository interface {
	GetGoogleSheetByID(newAction models.RequestGoogleAction) string
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
