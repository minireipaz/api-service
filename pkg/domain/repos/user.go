package repos

import "minireipaz/pkg/domain/models"

type UserService interface {
	SynUser(user *models.SyncUserRequest) (created, exist bool)
}

type UserRedisRepository interface {
	CheckUserExist(user *models.SyncUserRequest) (bool, error)
	CheckLockExist(user *models.SyncUserRequest) (bool, error)
	InsertUser(user *models.SyncUserRequest) (locked bool, lockExists bool, userExists bool, err error)
	RemoveLock(user *models.SyncUserRequest) bool
}

type UserBrokerRepository interface {
	CreateUser(user *models.SyncUserRequest) bool
}

type UserHTTPRepository interface {
}
