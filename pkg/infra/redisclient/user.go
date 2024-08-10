package redisclient

import (
	"fmt"
	"log"
	"minireipaz/pkg/common"
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

func (u *UserRedisRepository) CheckUserExist(user *models.SyncUserRequest) (exist bool, err error) {
	key := fmt.Sprintf("users:%s", user.Sub)
	for i := 1; i < models.MaxAttempts; i++ {
		countKeys, err := u.redisClient.Exists(key)
		if err == nil && countKeys == 1 {
			return true, err
		}
		waitTime := common.RandomDuration(models.MaxSleepDuration, models.MinSleepDuration, i)
		log.Printf("ERROR | Cannot check if exist lock for user %s, attempt %d: %v. Retrying in %v", user.Sub, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return false, fmt.Errorf("ERROR | Cannot check if user exist %s. More than 10 intents", user.Sub)
}

func (u *UserRedisRepository) CheckLockExist(user *models.SyncUserRequest) (exist bool, err error) {
	key := fmt.Sprintf("lock:users:%s", user.Sub)
	for i := 1; i < models.MaxAttempts; i++ {
		countKeys, err := u.redisClient.Exists(key)
		if err == nil && countKeys == 1 {
			return true, err
		}
		waitTime := common.RandomDuration(models.MaxSleepDuration, models.MinSleepDuration, i)
		log.Printf("ERROR | Cannot check if exist lock for user %s, attempt %d: %v. Retrying in %v", user.Sub, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return false, fmt.Errorf("ERROR | Cannot check if exist lock for user %s. More than 10 intents", user.Sub)
}

func (u *UserRedisRepository) InsertUser(user *models.SyncUserRequest) (inserted, lockExists, userExists bool, err error) {
	lockKey := fmt.Sprintf("lock:user:%s", user.Sub)
	userKey := fmt.Sprintf("users:%s", user.Sub)
	duration := 20 * time.Second

	for i := 1; i < models.MaxAttempts; i++ {
		inserted, lockExists, userExists, err = u.redisClient.WatchUser(
			user,
			lockKey,
			userKey,
			duration,
		) // Transacction
		if err == nil && inserted {
			return inserted, lockExists, userExists, err
		}
		if lockExists || userExists {
			return inserted, lockExists, userExists, err
		}
		waitTime := common.RandomDuration(models.MaxSleepDuration, models.MinSleepDuration, i)
		log.Printf("ERROR | Cannot check if exist lock for user %s, attempt %d: %v. Retrying in %v", user.Sub, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return inserted, lockExists, userExists, fmt.Errorf("ERROR | Cannot insert user %s. More than 10 intents", user.Sub)
}

func (u *UserRedisRepository) AddLock(user *models.SyncUserRequest) (locked bool, err error) {
	key := fmt.Sprintf("lock:user:%s", user.Sub)
	duration := 20 * time.Second

	for i := 1; i < models.MaxAttempts; i++ {
		locked, err = u.redisClient.acquireLock(key, "dummy", duration)
		if err == nil {
			return locked, err
		}

		waitTime := common.RandomDuration(models.MaxSleepDuration, models.MinSleepDuration, i)
		log.Printf("ERROR | Cannot create lock for user %s, attempt %d: %v. Retrying in %v", user.Sub, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return false, fmt.Errorf("ERROR | Cannot create lock for user %s. More than 10 intents", user.Sub)
}

func (u *UserRedisRepository) RemoveLock(user *models.SyncUserRequest) (removedLock bool) {
	key := fmt.Sprintf("lock:user:%s", user.Sub)
	for i := 1; i < models.MaxAttempts; i++ {
		countRemoved, err := u.redisClient.removeLock(key)
		if countRemoved == 0 {
			log.Printf("WARNING | Key already removed, previuous process take more than 20 seconds")
		}
		if err == nil && countRemoved <= 1 {
			return true
		}

		waitTime := common.RandomDuration(models.MaxSleepDuration, models.MinSleepDuration, i)
		log.Printf("ERROR | Cannot connect to redis for user %s, attempt %d: %v. Retrying in %v", user.Sub, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return false
}
