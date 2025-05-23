package masterpiece_controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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

func (c *MasterpieceController) GetMasterpieces(ctx *fiber.Ctx) error {
	const op = "controller.Masterpieces.masterpieceController.GetMasterpieces"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	responseRepo, code, opResp, msg, err := c.masterRepo.GetMasterpiecesRepository()
	if err != nil {
		c.logError(claims.Email, err, msg)
		return c.handleError(ctx, code, opResp, err, msg)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": msg,
		"data": fiber.Map{
			"masterpieces": responseRepo,
		},
	})
}

func (c *MasterpieceController) GetMasterpieceById(ctx *fiber.Ctx) error {
	const op = "controller.Masterpieces.GetMasterpieceById"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	masterpieceID := ctx.Params("id")
	if masterpieceID == "" {
		err := errors.New("invalid id")
		c.logError(claims.Email, err, "Invalid id parameter")
		return c.handleError(ctx, fiber.StatusBadRequest, op, err, "Invalid masterpiece id parameter")
	}

	responseRepo, code, opRepo, msg, err := c.masterRepo.GetMasterpieceById(masterpieceID)
	if err != nil {
		c.logError(claims.Email, err, msg)
		return c.handleError(ctx, code, opRepo, err, msg)
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": msg,
		"data": fiber.Map{
			"masterpiece": responseRepo,
		},
	})
}

func (c *MasterpieceController) GetMasterpiecesByStatusId(ctx *fiber.Ctx) error {
	const op = "controller.Masterpieces.GetMasterpiecesByStatusId"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	statusID := ctx.Params("status_id")
	if statusID == "" {
		err := errors.New("invalid masterpieces statusId")
		c.logError(claims.Email, err, "Invalid masterpieces statusId")
		return c.handleError(ctx, fiber.StatusBadRequest, op, err, "Invalid masterpieces statusId")
	}

	responseRepo, code, opRepo, msg, err := c.masterRepo.GetMasterpiecesByStatusId(statusID)
	if err != nil {
		c.logError(claims.Email, err, msg)
		return c.handleError(ctx, code, opRepo, err, msg)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": msg,
		"data": fiber.Map{
			"masterpieces": responseRepo,
		},
	})
}

func (c *MasterpieceController) SearchMasterpiecesSocket(conn *websocket.Conn) error {
	c.masterRepo.SearchMasterpiecesSocket(conn)
	return nil
}

func (c *MasterpieceController) CreateComment(ctx *fiber.Ctx) error {
	const op = "controller.Masterpieces.CreateComment"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	var comment identity.Comments

	if err := ctx.BodyParser(&comment); err != nil {
		c.logError(claims.Email, err, "Failed to parse body")
		return c.handleError(ctx, fiber.StatusBadRequest, op, err, "Failed to parse body")
	}

	comment.UserId = claims.UserID

	responseRepo, code, opRepo, msg, err := c.masterRepo.CreateCommentRepository(comment)
	if err != nil {
		c.logError(claims.Email, err, msg)
		return c.handleError(ctx, code, opRepo, err, msg)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": msg,
		"data": fiber.Map{
			"masterpieces": responseRepo,
		},
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
