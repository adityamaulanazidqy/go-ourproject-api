package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	"os"
	"strconv"
	"time"
)

func RateLimiter() fiber.Handler {
	err := godotenv.Load(".env")
	if err != nil {
		logLogrus.WithFields(logrus.Fields{
			"error":   err,
			"message": "Error loading .env file",
		}).Error("Error loading .env file for redis")

		return nil
	}

	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		logLogrus.WithFields(logrus.Fields{
			"error":   err,
			"message": "Error converting REDIS_PORT to int",
		}).Error("Error converting REDIS_PORT to int")

		return nil
	}

	store := redis.New(redis.Config{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     port,
		Password: os.Getenv("REDIS_PASSWORD"),
		Database: 0,
		Reset:    false,
	})

	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		Storage:    store,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(helpers.ApiResponse{
				Message: "Too many requests, please try again later.",
			})
		},
	})
}
