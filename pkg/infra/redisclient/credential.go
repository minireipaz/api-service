package redisclient

import (
	"fmt"
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/domain/models"
	"time"
)

type CredentialRedisRepository struct {
	redisClient *RedisClient
}

func NewCredentialRedisRepository(newRedisClient *RedisClient) *CredentialRedisRepository {
	return &CredentialRedisRepository{
		redisClient: newRedisClient,
	}
}

func (c *CredentialRedisRepository) SaveTemporalAuthURLData(currentCredential *models.RequestCreateCredential) (inserted bool, err error) {
	userKey := fmt.Sprintf("temp:code:%s", currentCredential.Data.RedirectURL)

	for i := 1; i < models.MaxAttempts; i++ {
		saved, err := c.redisClient.SetEx(userKey, currentCredential.Sub, 12*time.Hour)
		if err == nil && saved {
			return true, nil
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot check if exist lock for user %s, attempt %d: %v. Retrying in %v", currentCredential.Sub, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return inserted, fmt.Errorf("ERROR | Cannot insert user %s. More than 10 intents", currentCredential.Sub)
}

func (c *CredentialRedisRepository) AddLock(sub *string) (locked bool, err error) {
	key := fmt.Sprintf("lock:credential:user:%s", *sub)
	duration := 5 * time.Second

	for i := 1; i < models.MaxAttempts; i++ {
		locked, err = c.redisClient.AcquireLock(key, "dummy", duration)
		if err == nil {
			return locked, err
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot create lock for credential user %s, attempt %d: %v. Retrying in %v", *sub, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return false, fmt.Errorf("ERROR | Cannot create lock for credential user %s. More than 10 intents", *sub)
}

func (c *CredentialRedisRepository) RemoveLock(sub *string) (removedLock bool) {
	key := fmt.Sprintf("lock:credential:user:%s", *sub)
	for i := 1; i < models.MaxAttempts; i++ {
		countRemoved, err := c.redisClient.RemoveLock(key)
		if countRemoved == 0 {
			log.Printf("WARNING | Key already removed, previuous process take more than 20 seconds")
		}
		if err == nil && countRemoved <= 1 {
			return true
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to redis for user %s, attempt %d: %v. Retrying in %v", *sub, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return false
}
