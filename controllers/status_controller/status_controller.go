package status_controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	"go-ourproject/models/jwt_models"
	"go-ourproject/repositories/status_repository"
	"gorm.io/gorm"
)

type StatusController struct {
	db         *gorm.DB
	logLogrus  *logrus.Logger
	rdb        *redis.Client
	statusRepo *status_repository.StatusRepository
}

func NewStatusController(db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) *StatusController {
	return &StatusController{
		db:         db,
		logLogrus:  logLogrus,
		rdb:        rdb,
		statusRepo: status_repository.NewStatusRepository(db, logLogrus, rdb),
	}
}

func (c *StatusController) StatusMasterpiece(ctx *fiber.Ctx) error {
	const op = "masterpiece.controller.StatusMasterpiece"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	responseRepo, code, err := c.statusRepo.StatusMasterpieceRepository()
	if err != nil {
		c.logError(claims.Email, err, "Failed to get status master piece")
		return c.handleError(ctx, code, op, err, "Failed to get status master piece")
	}

	return ctx.Status(code).JSON(responseRepo)
}

func (c *StatusController) StatusThesis(ctx *fiber.Ctx) error {
	const op = "masterpiece.controller.StatusThesis"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	responseRepo, code, err := c.statusRepo.StatusThesisRepository()
	if err != nil {
		c.logError(claims.Email, err, "Failed to get status thesis")
		return c.handleError(ctx, code, op, err, "Failed to get status thesis")
	}

	return ctx.Status(code).JSON(responseRepo)
}

func (c *StatusController) logError(email string, err error, message string) {
	fields := logrus.Fields{"email": email, "message": message}
	if err != nil {
		fields["error"] = err.Error()
	}
	c.logLogrus.WithFields(fields).Error(message)
}

func (c *StatusController) handleError(ctx *fiber.Ctx, status int, op string, err error, message string) error {
	fields := logrus.Fields{
		"operation": op,
		"error":     err,
	}
	c.logLogrus.WithFields(fields).Error(message)

	return ctx.Status(status).JSON(helpers.ApiResponse{
		Message: message,
	})
}
