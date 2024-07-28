package redisclient

import (
	"context"
	"fmt"
	"log"
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/models"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisClient() *RedisClient {
	opt, err := redis.ParseURL(config.GetEnv("REDIS_URI", ""))
	if err != nil {
		log.Panicf("ERROR | Not connected to Redis")
	}

	rdb := redis.NewClient(opt)
	// _, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	if rdb.Ping(context.Background()).Val() != "PONG" {
		log.Panicf("ERROR | Not connected to Redis")
	}

	return &RedisClient{
		Client: rdb,
		Ctx:    context.Background(),
	}
}

func (r *RedisClient) Set(key string, value interface{}) error {
	return r.Client.Set(r.Ctx, key, value, 0).Err()
}

func (r *RedisClient) Hset(key string, field string, values interface{}) bool {
	inserted := r.Client.HSet(r.Ctx, key, field, values).Val()
	return inserted != 0
}

func (r *RedisClient) Hget(key string, field string) error {
	return r.Client.HGet(r.Ctx, key, field).Err()
}

func (r *RedisClient) Hexists(key string, field string) bool {
	return r.Client.HExists(r.Ctx, key, field).Val()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.Client.Get(r.Ctx, key).Result()
}

func (r *RedisClient) WatchWorkflow(workflow *models.Workflow) error {
	return r.Client.Watch(r.Ctx, func(tx *redis.Tx) error {
		return r.CheckAndSetWorkflow(r.Ctx, tx, workflow)
	})
}

func (r *RedisClient) CheckAndSetWorkflow(ctx context.Context, tx *redis.Tx, workflow *models.Workflow) error {
	uuidExists, err := tx.HExists(ctx, "workflows:all", workflow.UUID.String()).Result()
	if err != nil {
		log.Printf("error checking UUID existence: %v", err)
		return fmt.Errorf(models.UUIDCannotGenerate)
	}
	if uuidExists {
		return fmt.Errorf(models.UUIDExist)
	}

	nameExists, err := tx.HExists(ctx, fmt.Sprintf("users:%s", workflow.Sub), workflow.WorkflowName).Result()
	if err != nil {
		log.Printf("error checking workflow name existence: %v", err)
		return fmt.Errorf(models.WorkflowNameCannotGenerate)
	}
	if nameExists {
		return fmt.Errorf(models.WorkflowNameExist)
	}

	_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, fmt.Sprintf("users:%s", workflow.Sub), workflow.WorkflowName, workflow.UUID.String())
		pipe.HSet(ctx, "workflows:all", workflow.UUID.String(), workflow.Sub)
		return nil
	})

	return err
}

func (r *RedisClient) WatchToken(data string, key string, expiresInSeconds time.Duration) error {
	err := r.Client.Watch(r.Ctx, func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(r.Ctx, func(pipe redis.Pipeliner) error {
			pipe.SetNX(r.Ctx, key, data, expiresInSeconds)
			return nil
		})
		return err
	}, key)

	return err
}