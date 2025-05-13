package semesters_controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	identity "go-ourproject/models/identities"
	"go-ourproject/models/jwt_models"
	"go-ourproject/repositories/semesters_repository"
	"gorm.io/gorm"
	"strconv"
)

type SemestersController struct {
	db            *gorm.DB
	logLogrus     *logrus.Logger
	rdb           *redis.Client
	semestersRepo *semesters_repository.SemestersRepository
}

func NewSemestersController(db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) *SemestersController {
	return &SemestersController{
		db:            db,
		logLogrus:     logLogrus,
		rdb:           rdb,
		semestersRepo: semesters_repository.NewSemestersRepository(db, logLogrus, rdb),
	}
}

func (c *SemestersController) FindAllSemesters(ctx *fiber.Ctx) error {
	const op = "semesters.controller.FindAllSemesters"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	responseRepo, code, err := c.semestersRepo.FindAllSemestersRepository()
	if err != nil {
		c.logError(claims.Email, err, "Failed to get semesters")
		return c.handleError(ctx, code, op, err, "Failed to get semesters")
	}

	return ctx.Status(code).JSON(responseRepo)
}

func (c *SemestersController) FindSemestersById(ctx *fiber.Ctx) error {
	const op = "masterpiece.controller.FindSemestersById"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	semesterIDStr := ctx.Params("id", "")
	if semesterIDStr == "" {
		err := errors.New("missing semesterID")
		c.logError(claims.Email, err, "Failed to get semesterID")
		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get semesterID")
	}

	semesterID, err := strconv.Atoi(semesterIDStr)
	if err != nil {
		c.logError(claims.Email, err, "Failed to convert semesterID to int")
		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to convert semesterID to int")
	}

	responseRepo, code, err := c.semestersRepo.FindSemestersByIdRepository(semesterID)
	if err != nil {
		c.logError(claims.Email, err, "Failed to get semesters by id")
		return c.handleError(ctx, code, op, err, "Failed to get semesters by id")
	}

	return ctx.Status(code).JSON(responseRepo)
}

func (c *SemestersController) CreateSemesters(ctx *fiber.Ctx) error {
	const op = "masterpiece.controller.CreateSemesters"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	var semester identity.Semesters
	err := ctx.BodyParser(&semester)
	if err != nil {
		c.logError(claims.Email, err, "Failed to parse body semesters")
		return c.handleError(ctx, fiber.StatusBadRequest, op, err, "Failed to parse body")
	}

	responseRepo, code, err := c.semestersRepo.CreateSemestersRepository(&semester)
	if err != nil {
		c.logError(claims.Email, err, "Failed to add semesters")
		return c.handleError(ctx, code, op, err, "Failed to add semesters")
	}

	return ctx.Status(code).JSON(responseRepo)
}

func (c *SemestersController) logError(email string, err error, message string) {
	fields := logrus.Fields{"email": email, "message": message}
	if err != nil {
		fields["error"] = err.Error()
	}
	c.logLogrus.WithFields(fields).Error(message)
}

func (c *SemestersController) handleError(ctx *fiber.Ctx, status int, op string, err error, message string) error {
	fields := logrus.Fields{
		"operation": op,
		"error":     err,
	}
	c.logLogrus.WithFields(fields).Error(message)

	return ctx.Status(status).JSON(helpers.ApiResponse{
		Message: message,
	})
}
