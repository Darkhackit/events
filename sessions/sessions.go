package sessions

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	rc *redis.Client
}

func (r *RedisClient) CreateSession(ctx context.Context, sessionID string, userData interface{}, ttl time.Duration) error {
	if err := r.rc.Set(ctx, sessionID, userData, ttl).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) GetSession(ctx context.Context, sessionID string) (string, error) {
	result, err := r.rc.Get(ctx, sessionID).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r *RedisClient) DeleteSession(ctx context.Context, sessionID string) error {
	if err := r.rc.Del(ctx, sessionID).Err(); err != nil {
		return err
	}
	return nil
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
	return &RedisClient{rc: client}
}
