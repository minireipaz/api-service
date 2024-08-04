package redisclient

import (
	"fmt"
	"log"
	"minireipaz/pkg/domain/models"
	"time"
)

const (
	offset    = 1 * time.Second
	timedrift = 500 * time.Millisecond
  maxIntents = 10
)

type UserRedisRepository struct {
	redisClient *RedisClient
}

func NewUserRedisRepository(newRedisClient *RedisClient) *UserRedisRepository {
	return &UserRedisRepository{
		redisClient: newRedisClient,
	}
}

func (u *UserRedisRepository) CheckUserExist(user *models.Users) (exist bool, err error) {
	key := fmt.Sprintf("users:%s", user.Sub)
  for i := 0; i < maxIntents; i++ {
	countKeys, err := u.redisClient.Exists(key)
  if err == nil && countKeys == 1{
      return true, err
    }
    waitTime := offset + time.Duration(i)*timedrift
		log.Printf("ERROR | Cannot check if exist lock for user %s, attempt %d: %v. Retrying in %v", user.Sub, i, err, waitTime)
		time.Sleep(waitTime)
  }
	return false, fmt.Errorf("ERROR | Cannot check if user exist %s. More than 10 intents", user.Sub)
}

func (u *UserRedisRepository) CheckLockExist(user *models.Users) (exist bool, err error) {
	key := fmt.Sprintf("lock:users:%s", user.Sub)
  for i := 0; i < maxIntents; i++ {
    countKeys, err := u.redisClient.Exists(key)
    if err == nil && countKeys == 1{
      return true, err
    }
    waitTime := offset + time.Duration(i)*timedrift
		log.Printf("ERROR | Cannot check if exist lock for user %s, attempt %d: %v. Retrying in %v", user.Sub, i, err, waitTime)
		time.Sleep(waitTime)
  }
	return false, fmt.Errorf("ERROR | Cannot check if exist lock for user %s. More than 10 intents", user.Sub)
}

func (u *UserRedisRepository) AddLock(user *models.Users) (locked bool, err error) {
	key := fmt.Sprintf("lock:user:%s", user.Sub)
	duration := time.Duration(20 * time.Second)

	for i := 0; i < maxIntents; i++ {
		locked, err = u.redisClient.acquireLock(key, "dummy", duration)
		if err == nil {
			return locked, err
		}

		waitTime := offset + time.Duration(i)*timedrift
		log.Printf("ERROR | Cannot create lock for user %s, attempt %d: %v. Retrying in %v", user.Sub, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return false, fmt.Errorf("ERROR | Cannot create lock for user %s. More than 10 intents", user.Sub)
}

func (u *UserRedisRepository) RemoveLock(user *models.Users) (removedLock bool) {
	key := fmt.Sprintf("lock:user:%s", user.Sub)
	for i := 0; i < maxIntents; i++ {
		countRemoved, err := u.redisClient.removeLock(key)
		if countRemoved == 0 {
			log.Printf("WARNING | Key already removed, previuous process take more than 20 seconds")
		}
		if err == nil && countRemoved <= 1 {
			return true
		}

		waitTime := offset + time.Duration(i)*timedrift
		log.Printf("ERROR | Cannot connect to redis for user %s, attempt %d: %v. Retrying in %v", user.Sub, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return false
}
