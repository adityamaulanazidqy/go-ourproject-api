package masterpiece_controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	identity "go-ourproject/models/identities"
	"go-ourproject/models/jwt_models"
	"go-ourproject/repositories/masterpiece_repository"
	"gorm.io/gorm"
	"path/filepath"
	"strconv"
	"time"
)

type MasterpieceController struct {
	db        *gorm.DB
	logLogrus *logrus.Logger
	rdb       *redis.Client

	masterRepo *masterpiece_repository.MasterpieceRepository
}

func NewMasterpieceController(db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) *MasterpieceController {
	return &MasterpieceController{
		db:         db,
		logLogrus:  logLogrus,
		rdb:        rdb,
		masterRepo: masterpiece_repository.NewMasterpieceRepository(db, logLogrus, rdb),
	}
}

func (c *MasterpieceController) PostMasterpiece(ctx *fiber.Ctx) error {
	const op = "masterpiece.controller.PostMasterpiece"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	statusID, err := strconv.Atoi(ctx.FormValue("status_id"))
	if err != nil {
		c.logError(claims.Email, err, "invalid status_id")
		return c.handleError(ctx, fiber.StatusBadRequest, op, err, "Invalid status ID")
	}

	classID, err := strconv.Atoi(ctx.FormValue("class_id"))
	if err != nil {
		c.logError(claims.Email, err, "invalid class_id")
		return c.handleError(ctx, fiber.StatusBadRequest, op, err, "Invalid class ID")
	}

	semesterID, err := strconv.Atoi(ctx.FormValue("semester_id"))
	if err != nil {
		c.logError(claims.Email, err, "invalid semester_id")
		return c.handleError(ctx, fiber.StatusBadRequest, op, err, "Invalid semester ID")
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		c.logError(claims.Email, err, "failed to parse form")
		return c.handleError(ctx, fiber.StatusBadRequest, op, err, "Invalid form data")
	}

	files := form.File["photos"]
	if len(files) == 0 {
		return c.handleError(ctx, fiber.StatusBadRequest, op,
			errors.New("no files uploaded"), "At least one photo is required")
	}

	var savedFiles []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.logError(claims.Email, err, "failed to open file")
			return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "failed to open file")
		}

		filename := filepath.Base(fileHeader.Filename)
		file.Close()

		savedFiles = append(savedFiles, filename)
	}

	if len(savedFiles) == 0 {
		c.logError(claims.Email, err, "no files uploaded")
		return c.handleError(ctx, fiber.StatusBadRequest, op,
			errors.New("no valid files"), "Could not save any uploaded files")
	}

	masterpiece := identity.Masterpiece{
		UserID:          claims.UserID,
		StatusID:        statusID,
		ClassID:         classID,
		SemesterID:      semesterID,
		PublicationDate: time.Now(),
		LinkGithub:      ctx.FormValue("link_github"),
	}

	responseRepo, code, opRepo, msg, err := c.masterRepo.CreateMasterpieceWithFiles(&masterpiece, savedFiles)
	if err != nil {
		c.logError(claims.Email, err, msg)
		return c.handleError(ctx, code, opRepo, err, msg)
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.logError(claims.Email, err, "failed to open file")
			return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "failed to open file")
		}

		_, err = helpers.SaveImages().Masterpiece(file, fileHeader, "_")
		if err != nil {
			c.logError(claims.Email, err, "failed to save masterpiece in storage")
			return c.handleError(ctx, code, opRepo, err, "failed to save masterpiece in storage")
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(helpers.ApiResponse{
		Message: "Masterpiece created successfully",
		Data:    responseRepo,
	})
}

func (c *MasterpieceController) logError(email string, err error, message string) {
	fields := logrus.Fields{"email": email, "message": message}
	if err != nil {
		fields["error"] = err.Error()
	}
	c.logLogrus.WithFields(fields).Error(message)
}

func (c *MasterpieceController) handleError(ctx *fiber.Ctx, status int, op string, err error, message string) error {
	fields := logrus.Fields{
		"operation": op,
		"error":     err,
	}
	c.logLogrus.WithFields(fields).Error(message)

	return ctx.Status(status).JSON(helpers.ApiResponse{
		Message: message,
	})
}
