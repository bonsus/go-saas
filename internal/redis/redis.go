package myredis

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/bonsus/go-saas/internal/config"
	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func initRedis() {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr: config.Cfg.Redis.Host + ":" + config.Cfg.Redis.Port,
		})
	})
}

func SetData(key string, data interface{}, expiration time.Duration) error {
	initRedis()
	if redisClient == nil {
		return errors.New("Redis client is not initialized")
	}

	// Encode struct ke JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.New("failed to marshal data")
	}

	err = redisClient.Set(key, jsonData, expiration).Err()
	if err != nil {
		return errors.New("failed to save data to Redis")
	}
	return nil
}

func GetData(key string, dest interface{}) error {
	initRedis()
	if redisClient == nil {
		return errors.New("Redis client is not initialized")
	}

	// Ambil data dari Redis
	data, err := redisClient.Get(key).Result()
	if err == redis.Nil {
		return errors.New("data not found")
	} else if err != nil {
		return errors.New("failed to fetch data from Redis")
	}
	// Decode JSON ke struct yang diberikan
	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		return errors.New("failed to unmarshal data")
	}
	return nil
}

func RemoveData(key string) error {
	initRedis()
	if redisClient == nil {
		return errors.New("Redis client is not initialized")
	}

	// Hapus data dari Redis
	err := redisClient.Del(key).Err()
	if err != nil {
		return errors.New("failed to delete data from Redis")
	}
	return nil
}
