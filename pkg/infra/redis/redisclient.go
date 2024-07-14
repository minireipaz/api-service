package redis

// import (
// 	"context"
// 	"minireipaz/pkg/config"

// 	"github.com/go-redis/redis/v8"
// )

// type RedisClient struct {
// 	Client *redis.Client
// }

// func NewRedisClient() *RedisClient {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr: config.GetEnv("VAULT_URI", ""),
// 	})

// 	return &RedisClient{
// 		Client: rdb,
// 	}
// }

// func (r *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
// 	return r.Client.Set(ctx, key, value, 0).Err()
// }

// func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
// 	return r.Client.Get(ctx, key).Result()
// }
