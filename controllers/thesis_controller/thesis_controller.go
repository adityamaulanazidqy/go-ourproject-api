package thesis_controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	identity "go-ourproject/models/identities"
	"go-ourproject/models/jwt_models"
	"go-ourproject/repositories/thesis_repository"
	"gorm.io/gorm"
)

type ThesisController struct {
	db         *gorm.DB
	logLogrus  *logrus.Logger
	rdb        *redis.Client
	thesisRepo *thesis_repository.ThesisRepository
}

func NewThesisController(db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) *ThesisController {
	return &ThesisController{
		db:         db,
		logLogrus:  logLogrus,
		rdb:        rdb,
		thesisRepo: thesis_repository.NewThesisRepository(db, logLogrus, rdb),
	}
}

func (c *ThesisController) CreateThesisTitle(ctx *fiber.Ctx) error {
	const op = "controllers.ThesisController.CreateThesisTitle"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	var thesisReq identity.ThesisRequest
	err := ctx.BodyParser(&thesisReq)
	if err != nil {
		c.logError(claims.Email, err, "Failed to parse thesis")
		return c.handleError(ctx, fiber.StatusBadRequest, op, err, "Failed to parse thesis")
	}

	thesisReq.UserID = claims.UserID

	responseRepo, code, opRepo, msg, err := c.thesisRepo.CreateThesisRepo(&thesisReq)
	if err != nil {
		c.logError(claims.Email, err, "Failed to create thesis")
		return c.handleError(ctx, code, opRepo, err, msg)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data": fiber.Map{
			"thesis": responseRepo,
		},
	})
}

func (c *ThesisController) CreateSupervision(ctx *fiber.Ctx) error {
	const op = "controllers.ThesisController.CreateSupervision"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	var SupervisionReq identity.SupervisionRequest
	err := ctx.BodyParser(&SupervisionReq)
	if err != nil {
		c.logError(claims.Email, err, "Failed to parse thesis")
		return c.handleError(ctx, fiber.StatusBadRequest, op, err, "Failed to parse request")
	}

	SupervisionReq.TeacherID = claims.UserID

	responseRepo, code, opRepo, msg, err := c.thesisRepo.CreateSupervisionRepo(&SupervisionReq)
	if err != nil {
		c.logError(claims.Email, err, "Failed to create thesis")
		return c.handleError(ctx, code, opRepo, err, msg)
	}

	return ctx.Status(code).JSON(fiber.Map{
		"message": msg,
		"data": fiber.Map{
			"supervision": responseRepo,
		},
	})
}

func (c *ThesisController) logError(email string, err error, message string) {
	fields := logrus.Fields{"email": email, "message": message}
	if err != nil {
		fields["error"] = err.Error()
	}
	c.logLogrus.WithFields(fields).Error(message)
}

func (c *ThesisController) handleError(ctx *fiber.Ctx, status int, op string, err error, message string) error {
	fields := logrus.Fields{
		"operation": op,
		"error":     err,
	}
	c.logLogrus.WithFields(fields).Error(message)

	return ctx.Status(status).JSON(helpers.ApiResponse{
		Message: message,
	})
}
