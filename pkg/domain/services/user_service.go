package services

import (
	"log"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/infra/brokerclient"
	"minireipaz/pkg/infra/httpclient"
	"minireipaz/pkg/infra/redisclient"
	"time"
)

type UserServiceInterface interface {
	SynUser(user *models.Users) (created, exist bool)
}

var _ UserServiceInterface = (*UserService)(nil)

type UserService struct {
	repoHTTP   *httpclient.UserRepository
	repoRedis  *redisclient.UserRedisRepository
	repoBroker *brokerclient.UserKafkaRepository
}

func NewUserService(newRepoHTTP *httpclient.UserRepository, newRepoRedis *redisclient.UserRedisRepository, newRepoBroker *brokerclient.UserKafkaRepository) *UserService {
	return &UserService{
		repoHTTP:   newRepoHTTP,
		repoRedis:  newRepoRedis,
		repoBroker: newRepoBroker,
	}
}

func (u *UserService) SynUser(user *models.Users) (created, exist bool) {
	exist, err := u.repoRedis.CheckExist(user)
  if err != nil {
    log.Printf("ERROR | Cannot access to repo redis %v", err)
    return false, false
  }
	if exist {
		return false, true
	}
	// not exist generate lock for about 20 seconds
	// and 10 retries
	for i := 0; i < 10; i++ {
		locked, _ := u.repoRedis.AddLock(user)
		if locked {
			break
		}
		time.Sleep(1 * time.Second)
	}

  setDefaults(user)

  created, exist = u.repoBroker.Create(user)
	if exist {
		log.Printf("WARN | Already exist in Result DB")
	}
	// in case cannot get autoremoved
	u.repoRedis.RemoveLock(user)

	return created, exist
}

func setDefaults(user *models.Users) {
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
