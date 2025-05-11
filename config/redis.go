package config

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	ctx = context.Background()
)

func ConnRedis(logLogrus *logrus.Logger) *redis.Client {
	err := godotenv.Load(".env")
	if err != nil {
		logLogrus.WithFields(logrus.Fields{
			"error":   err,
			"message": "Error loading .env file",
		}).Error("Error loading .env file for redis")

		return nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		logLogrus.WithFields(logrus.Fields{
			"error":   err,
			"message": "Error connecting to redis",
		}).Error("Error connecting to redis")
		return nil
	}

	return rdb
}
