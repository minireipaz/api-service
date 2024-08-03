package redisclient

import (
	"fmt"
	"log"
	"minireipaz/pkg/domain/models"
	"time"
)

type UserRedisRepository struct {
	redisClient *RedisClient
}

func NewUserRedisRepository(newRedisClient *RedisClient) *UserRedisRepository {
	return &UserRedisRepository{
		redisClient: newRedisClient,
	}
}

func (u *UserRedisRepository) CheckExist(user *models.Users) (exist bool, err error) {
	key := fmt.Sprintf("users:%s", user.Sub)
	countKeys, err := u.redisClient.Exists(key)
	return countKeys == 1, err
}

func (u *UserRedisRepository) AddLock(user *models.Users) (locked bool, err error) {
	key := fmt.Sprintf("lock:user:%s", user.Sub)
	duration := time.Duration(20 * time.Second)
	locked, err = u.redisClient.acquireLock(key, "dummy", duration)
	if err != nil {
		log.Printf("ERROR | Cannot create lock for user %s", user.Sub)
		return false, err
	}
	return locked, err
}

func (u *UserRedisRepository) RemoveLock(user *models.Users) (removedLock bool) {
	key := fmt.Sprintf("lock:user:%s", user.Sub)
	countRemoved, err := u.redisClient.removeLock(key)
	if err != nil {
		log.Printf("ERROR | Cannot connect to redis for user %s", user.Sub)
		return false
	}
	if countRemoved == 0 {
		log.Printf("WARNING | Key already removed, previuous process take more than 20 seconds")
	}
	return countRemoved == 1
}
