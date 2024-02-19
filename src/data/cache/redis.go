package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Arshia-Izadyar/Go-Ecommerce/src/config"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis(cfg *config.Config) error {

	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		// Password:     cfg.Redis.Password,
		DB:           0,
		DialTimeout:  cfg.Redis.DialTimeout * time.Second,
		ReadTimeout:  cfg.Redis.ReadTimeout * time.Second,
		WriteTimeout: cfg.Redis.WriteTimeout * time.Second,
		PoolSize:     cfg.Redis.PoolSize,
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() {
	err := redisClient.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func Set[T any](key string, value T, duration time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = redisClient.Set(context.Background(), key, v, duration).Result()
	if err != nil {
		return err
	}
	return nil
}

func Get[T any](key string) (*T, error) {
	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	dest := new(T)
	err = json.Unmarshal([]byte(result), &dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
}
