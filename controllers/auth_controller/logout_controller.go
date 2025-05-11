package auth_controller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	"go-ourproject/middlewares"
	"time"
)

type LogoutController struct {
	rdb       *redis.Client
	logLogrus *logrus.Logger
}

func NewLogoutController(rdb *redis.Client, logLogrus *logrus.Logger) *LogoutController {
	return &LogoutController{
		rdb:       rdb,
		logLogrus: logLogrus,
	}
}

func (controller *LogoutController) Logout(c *fiber.Ctx) error {
	token, err := middlewares.ExtractTokenFromHeader(c)
	if err != nil || token == "" {
		controller.logLogrus.WithFields(logrus.Fields{
			"error":   err,
			"token":   token,
			"message": "Token extraction failed",
		}).Warn("Token missing or invalid")

		return c.Status(fiber.StatusUnauthorized).JSON(helpers.ApiResponse{
			Message: "Missing or invalid token.",
			Data:    nil,
		})
	}

	claims, err := middlewares.VerifyToken(token)
	if err != nil {
		controller.logLogrus.WithFields(logrus.Fields{
			"error":   err,
			"token":   token,
			"message": "Token verification failed",
		}).Warn("Unauthorized token")

		return c.Status(fiber.StatusUnauthorized).JSON(helpers.ApiResponse{
			Message: "Unauthorized.",
			Data:    nil,
		})
	}

	expDuration := time.Until(claims.ExpiresAt.Time)
	if expDuration <= 0 {
		expDuration = time.Minute * 1
	}

	ctx := context.Background()
	err = controller.rdb.Set(ctx, "blacklist:"+token, "true", expDuration).Err()
	if err != nil {
		controller.logLogrus.WithFields(logrus.Fields{
			"error":   err,
			"token":   token,
			"message": "Redis blacklist set failed",
		}).Error("Redis error")

		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ApiResponse{
			Message: "Failed to logout.",
			Data:    nil,
		})
	}

	controller.logLogrus.Info("User successfully logged out.")

	return c.Status(fiber.StatusOK).JSON(helpers.ApiResponse{
		Message: "You have successfully logged out.",
		Data:    nil,
	})
}
