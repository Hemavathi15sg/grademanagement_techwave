package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"techwave/models"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

func (rc *RedisCache) GetGrade(id int) (*models.Grade, error) {
	key := fmt.Sprintf("grade:%d", id)
	val, err := rc.client.Get(rc.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, err
	}

	var grade models.Grade
	if err := json.Unmarshal([]byte(val), &grade); err != nil {
		return nil, err
	}

	return &grade, nil
}

func (rc *RedisCache) SetGrade(grade models.Grade, ttl time.Duration) error {
	key := fmt.Sprintf("grade:%d", grade.ID)
	data, err := json.Marshal(grade)
	if err != nil {
		return err
	}

	return rc.client.Set(rc.ctx, key, data, ttl).Err()
}

func (rc *RedisCache) InvalidateGrade(id int) error {
	key := fmt.Sprintf("grade:%d", id)
	return rc.client.Del(rc.ctx, key).Err()
}

func (rc *RedisCache) Close() error {
	return rc.client.Close()
}
