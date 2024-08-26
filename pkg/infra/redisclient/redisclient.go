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

func (r *RedisClient) Exists(key string) (int64, error) {
	return r.Client.Exists(r.Ctx, key).Result()
}

func (r *RedisClient) Get(key string) (string, error) {
	result, err := r.Client.Get(r.Ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r *RedisClient) WatchWorkflow(workflow *models.Workflow, operation string) error {
	return r.Client.Watch(r.Ctx, func(tx *redis.Tx) error {
		return r.CheckAndModifyWorkflow(r.Ctx, tx, workflow, operation)
	})
}

func (r *RedisClient) CheckAndModifyWorkflow(ctx context.Context, tx *redis.Tx, workflow *models.Workflow, operation string) error {
	uuidExists, err := tx.HExists(ctx, "workflows:all", workflow.UUID.String()).Result()
	if err != nil {
		log.Printf("ERROR | checking UUID existence: %v", err)
		return fmt.Errorf(models.UUIDCannotGenerate)
	}

	// not necesary
	// nameExists, err := tx.HExists(ctx, fmt.Sprintf("users:%s", workflow.Sub), workflow.WorkflowName).Result()
	// if err != nil {
	// 	log.Printf("ERROR | checking workflow name existence: %v", err)
	// 	return fmt.Errorf(models.WorkflowNameCannotGenerate)
	// }

	switch operation {
	case "set":
		if uuidExists {
			return fmt.Errorf(models.UUIDExist)
		}
		// if nameExists {
		// 	return fmt.Errorf(models.WorkflowNameExist)
		// }
		return r.setWorkflow(ctx, tx, workflow)
	case "remove":
		// if !uuidExists {
		// 	return fmt.Errorf(models.UUIDNotExist)
		// }
		// if !nameExists {
		// 	return fmt.Errorf(models.WorkflowNameNotExist)
		// }
		return r.removeWorkflow(ctx, tx, workflow)
	default:
		return fmt.Errorf("unknown operation: %s", operation)
	}
}

func (r *RedisClient) setWorkflow(ctx context.Context, tx *redis.Tx, workflow *models.Workflow) error {
	_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, fmt.Sprintf("users:%s", workflow.Sub), workflow.WorkflowName, workflow.UUID.String())
		pipe.HSet(ctx, "workflows:all", workflow.UUID.String(), workflow.Sub)
		return nil
	})
	return err
}

func (r *RedisClient) removeWorkflow(ctx context.Context, tx *redis.Tx, workflow *models.Workflow) error {
	_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HDel(ctx, fmt.Sprintf("users:%s", workflow.Sub), workflow.WorkflowName)
		pipe.HDel(ctx, "workflows:all", workflow.UUID.String())
		return nil
	})
	return err
}

func (r *RedisClient) SetWorkflow(workflow *models.Workflow) error {
	return r.WatchWorkflow(workflow, "set")
}

func (r *RedisClient) RemoveWorkflow(workflow *models.Workflow) error {
	return r.WatchWorkflow(workflow, "remove")
}

func (r *RedisClient) WatchUser(user *models.SyncUserRequest, lockKey, userKey string, duration time.Duration) (inserted bool, lockExists bool, userExists bool, err error) {
	err = r.Client.Watch(r.Ctx, func(tx *redis.Tx) error {
		lockExists, err = checkLockExists(r.Ctx, tx, lockKey, user.Sub) // Quizas no es necesario
		if err != nil || lockExists {
			return err
		}

		userExists, err = checkUserExists(r.Ctx, tx, userKey, user.Sub)
		if err != nil || userExists {
			return err
		}

		inserted, err = executePipeline(r.Ctx, tx, lockKey, userKey, duration, user.Sub)
		return err
	}, lockKey)

	return inserted, lockExists, userExists, err
}

func checkLockExists(ctx context.Context, tx *redis.Tx, lockKey, userSub string) (bool, error) {
	lockExistsVal, err := tx.Exists(ctx, lockKey).Result()
	if err != nil {
		return false, fmt.Errorf("ERROR | Failed to check lock existence: %v", err)
	}
	if lockExistsVal == 1 {
		return true, fmt.Errorf("ERROR | Lock already exists for user %s", userSub)
	}
	return false, nil
}

func checkUserExists(ctx context.Context, tx *redis.Tx, userKey, userSub string) (bool, error) {
	userExistsVal, err := tx.Exists(ctx, userKey).Result()
	if err != nil {
		return false, fmt.Errorf("ERROR | Failed to check user existence: %v", err)
	}
	if userExistsVal == 1 {
		return true, fmt.Errorf("ERROR | User %s already exists", userSub)
	}
	return false, nil
}

func executePipeline(ctx context.Context, tx *redis.Tx, lockKey, userKey string, duration time.Duration, userSub string) (bool, error) {
	cmds, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SetNX(ctx, lockKey, "_", duration)
		pipe.HSet(ctx, userKey, "_", "_") // dummy value
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("ERROR | Transaction failed: %v", err)
	}

	for i, cmd := range cmds {
		if cmd.Err() != nil {
			return false, fmt.Errorf("ERROR | Command failed: %v for user %s", cmd.Err(), userSub)
		}

		switch i {
		case 0: // SetNX command
			lockSet, err := cmd.(*redis.BoolCmd).Result()
			if err != nil {
				return false, fmt.Errorf("ERROR | Failed to get SetNX result: %v", err)
			}
			if !lockSet {
				return false, fmt.Errorf("ERROR | Lock already exists for user %s", userSub)
			}
		case 1: // HSet command
			fieldsCreated, err := cmd.(*redis.IntCmd).Result()
			if err != nil {
				return false, fmt.Errorf("ERROR | Failed to get HSet result: %v", err)
			}
			if fieldsCreated == 0 {
				return false, fmt.Errorf("ERROR | User %s already exists", userSub)
			}
		}
	}

	return true, nil
}

func (r *RedisClient) WatchToken(data string, key string, expires time.Duration) error {
	err := r.Client.Watch(r.Ctx, func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(r.Ctx, func(pipe redis.Pipeliner) error {
			pipe.SetNX(r.Ctx, key, data, expires)
			return nil
		})
		return err
	}, key)

	return err
}

func (r *RedisClient) AcquireLock(key, value string, expiration time.Duration) (bool, error) {
	return r.Client.SetNX(r.Ctx, key, value, expiration).Result()
}

func (r *RedisClient) RemoveLock(key string) (int64, error) {
	result, err := r.Client.Del(r.Ctx, key).Result()
	return result, err
}
