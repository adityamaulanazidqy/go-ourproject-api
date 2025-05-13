package classes_controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	identity "go-ourproject/models/identities"
	"go-ourproject/models/jwt_models"
	"go-ourproject/repositories/classes_repository"
	"gorm.io/gorm"
	"strconv"
)

type ClassesController struct {
	db          *gorm.DB
	logLogrus   *logrus.Logger
	rdb         *redis.Client
	classesRepo *classes_repository.ClassesRepository
}

func NewClassesController(db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) *ClassesController {
	return &ClassesController{
		db:          db,
		logLogrus:   logLogrus,
		rdb:         rdb,
		classesRepo: classes_repository.NewClassesRepository(db, logLogrus, rdb),
	}
}

func (c *ClassesController) FindAllClasses(ctx *fiber.Ctx) error {
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

	responseRepo, code, err := c.classesRepo.FindAllClassesRepository()
	if err != nil {
		c.logError(claims.Email, err, "Failed to get classes")
		return c.handleError(ctx, code, op, err, "Failed to get classes")
	}

	return ctx.Status(code).JSON(responseRepo)
}

func (c *ClassesController) FindClassesById(ctx *fiber.Ctx) error {
	const op = "masterpiece.controller.FindClassesById"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	classIDStr := ctx.Params("id", "")
	if classIDStr == "" {
		err := errors.New("missing classID")
		c.logError(claims.Email, err, "Failed to get classID")
		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get classID")
	}

	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		c.logError(claims.Email, err, "Failed to convert classID to int")
		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to convert classID to int")
	}

	responseRepo, code, err := c.classesRepo.FindClassesByIdRepository(classID)
	if err != nil {
		c.logError(claims.Email, err, "Failed to get classes by id")
		return c.handleError(ctx, code, op, err, "Failed to get classes by id")
	}

	return ctx.Status(code).JSON(responseRepo)
}

func (c *ClassesController) CreateClasses(ctx *fiber.Ctx) error {
	const op = "masterpiece.controller.CreateClasses"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	var class identity.Classes
	err := ctx.BodyParser(&class)
	if err != nil {
		c.logError(claims.Email, err, "Failed to parse body classes")
		return c.handleError(ctx, fiber.StatusBadRequest, op, err, "Failed to parse body")
	}

	responseRepo, code, err := c.classesRepo.CreateClassesRepository(&class)
	if err != nil {
		c.logError(claims.Email, err, "Failed to add classes")
		return c.handleError(ctx, code, op, err, "Failed to add classes")
	}

	return ctx.Status(code).JSON(responseRepo)
}

func (c *ClassesController) logError(email string, err error, message string) {
	fields := logrus.Fields{"email": email, "message": message}
	if err != nil {
		fields["error"] = err.Error()
	}
	c.logLogrus.WithFields(fields).Error(message)
}

func (c *ClassesController) handleError(ctx *fiber.Ctx, status int, op string, err error, message string) error {
	fields := logrus.Fields{
		"operation": op,
		"error":     err,
	}
	c.logLogrus.WithFields(fields).Error(message)

	return ctx.Status(status).JSON(helpers.ApiResponse{
		Message: message,
	})
}
