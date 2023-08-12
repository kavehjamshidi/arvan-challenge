package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
)

func Setup() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_ADDRESS"),
		Password: viper.GetString("REDIS_PASSWORD"),
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Panicf("could not connect to redis: %v\n", err)
	}

	return redisClient
}

func TestSetup() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("TEST_REDIS_ADDRESS"),
		Password: viper.GetString("TEST_REDIS_PASSWORD"),
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Panicf("could not connect to redis: %v\n", err)
	}

	return redisClient
}
