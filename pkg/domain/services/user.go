package services

import (
	"log"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
)

// type UserServiceInterface interface {
// 	SynUser(user *models.SyncUserRequest) (created, exist bool)
// }

// var _ UserServiceInterface = (*UserService)(nil)

type UserService struct {
	userHTTPRepo   repos.UserHTTPRepository
	userRedisRepo  repos.UserRedisRepository
	userBrokerRepo repos.UserBrokerRepository
}

func NewUserService(repoHTTP repos.UserHTTPRepository, repoRedis repos.UserRedisRepository, repoBroker repos.UserBrokerRepository) *UserService {
	return &UserService{
		userHTTPRepo:   repoHTTP,
		userRedisRepo:  repoRedis,
		userBrokerRepo: repoBroker,
	}
}

// type UserService struct {
// 	repoHTTP   *httpclient.UserRepository
// 	repoRedis  *redisclient.UserRedisRepository
// 	repoBroker *brokerclient.UserKafkaRepository
// }

// func NewUserService(newRepoHTTP *httpclient.UserRepository, newRepoRedis *redisclient.UserRedisRepository, newRepoBroker *brokerclient.UserKafkaRepository) *UserService {
// 	return &UserService{
// 		repoHTTP:   newRepoHTTP,
// 		repoRedis:  newRepoRedis,
// 		repoBroker: newRepoBroker,
// 	}
// }

func (u *UserService) SynUser(user *models.SyncUserRequest) (created, exist bool) {
	exist, err := u.userRedisRepo.CheckUserExist(user)
	if err != nil {
		log.Printf("ERROR | Cannot access to repo redis %v", err)
		return false, false
	}
	if exist {
		return false, true
	}

	// new user
	exist, err = u.userRedisRepo.CheckLockExist(user)
	if err != nil {
		log.Printf("ERROR | Cannot access to repo redis %v", err)
		return false, false
	}
	if exist {
		return false, false
	}
	// InsertUser insert user and generate a lock for about 20 seconds
	locked, lockExists, userExists, err := u.userRedisRepo.InsertUser(user)
	if err != nil {
		// TODO: Dead letter
		log.Printf("ERROR | Needs to Added to dead letter %s", user.Sub)
	}

	if !locked {
		log.Printf("WARN | Cannot created lock for user %s", user.Sub)
	}

	if userExists {
		return false, true
	}

	if userExists && !lockExists {
		return false, true
	}

	setUserDefaults(user) // default roleID
	sended := u.userBrokerRepo.CreateUser(user)

	// in case cannot get autoremoved
	u.userRedisRepo.RemoveLock(user)

	return sended, false
}

func setUserDefaults(user *models.SyncUserRequest) {
	if user.Status == 0 {
		user.Status = models.StatusActive
	}

	if user.RoleID == 0 {
		user.RoleID = generateDefaultUserRoleID()
	}
}

func generateDefaultUserRoleID() models.UserRoleID {
	return models.RoleUser
}
