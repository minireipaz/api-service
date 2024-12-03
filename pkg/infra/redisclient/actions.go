package redisclient

import (
	"context"
	"fmt"
	"log"
	"minireipaz/pkg/common"
	"minireipaz/pkg/domain/models"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	ActionsGlobalAll = "actions:all"
	EmptyValue       = "_"
)

type ActionsService interface {
}

type ActionsRepository struct {
	redisClient *RedisClient
}

func NewActionsRepository(redisClient *RedisClient) *ActionsRepository {
	return &ActionsRepository{redisClient: redisClient}
}

func (a *ActionsRepository) GetActionsGlobalAll() string {
	return ActionsGlobalAll
}

func (a *ActionsRepository) ValidateActionGlobalUUID(uuid *string) (bool, error) {
	exists, err := a.redisClient.Hexists(ActionsGlobalAll, *uuid)
	if err != nil {
		log.Printf("ERROR | Redis HExists error: %v", err)
		return true, err
	}
	return exists, err
}

func (a *ActionsRepository) Create(newAction *models.RequestGoogleAction) (created bool, existed bool, err error) {
	ctx := context.Background()

	txf := func(tx *redis.Tx) error {
		lockKey := fmt.Sprintf("lock:%s", newAction.ActionID)
		key := ActionsGlobalAll
		field := newAction.ActionID
		value := newAction.Sub

		// Check if actionid already exists
		exists, err := tx.HExists(ctx, key, field).Result()
		if err != nil {
			return err
		}
		if exists {
			return ErrActionExists
		}

		// SETNX only insert to set the lock key if cannot insert lock exist
		// dummy value
		setnxRes, err := tx.SetNX(ctx, lockKey, "1", models.MaxTimeForLocks).Result()
		if err != nil {
			return err
		}
		if !setnxRes {
			// Lock already exists
			return ErrActionExists
		}

		// Proceed to set the action
		_, err = tx.HSet(ctx, key, field, value).Result()
		if err != nil {
			// Clean up the lock if setting the action fails
			tx.Del(ctx, lockKey)
			return err
		}

		return nil
	}

	created, existed, err = a.redisClient.SetAction(ctx, newAction, txf)
	return created, existed, err
}

func (a *ActionsRepository) Remove(newAction *models.RequestGoogleAction) bool {
	countRemoved, err := a.redisClient.Hdel(newAction.ActionID, newAction.Sub)
	if countRemoved == 1 && err == nil {
		return true
	}
	return false
}

func (a *ActionsRepository) AcquireLock(key, value string, expiration time.Duration) (locked bool, err error) {
	for i := 1; i < models.MaxAttempts; i++ {
		locked, err = a.redisClient.AcquireLock(key, value, expiration)
		if err == nil {
			return locked, err
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to redis for key %s, attempt %d: %v. Retrying in %v", key, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return false, fmt.Errorf("ERROR | Cannot create lock for key %s. More than 10 intents", key)
}

func (a *ActionsRepository) RemoveLock(key string) bool {
	for i := 1; i < models.MaxAttempts; i++ {
		countRemoved, err := a.redisClient.RemoveLock(key)
		if countRemoved == 0 {
			log.Printf("WARNING | Key already removed, previuous process take more than 20 seconds")
		}
		if err == nil && countRemoved <= 1 {
			return true
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to redis for key %s, attempt %d: %v. Retrying in %v", key, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return false
}

func (a *ActionsRepository) SetNX(hashKey, actionID string, expiration time.Duration) (bool, error) {
	inserted, err := a.redisClient.SetEx(hashKey, actionID, expiration)
	return inserted, err
}
